package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"golang.org/x/crypto/bcrypt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/mogenius/punq/dtos"
	"github.com/mogenius/punq/kubernetes"
	"github.com/mogenius/punq/logger"
	"github.com/mogenius/punq/utils"
)

const PunqAdminIdKey = "admin_id"

func InitUserService() {
	CreateUserSecret()
}

func CreateUserSecret() {
	provider, err := kubernetes.NewKubeProvider(nil)
	if provider == nil || err != nil {
		logger.Log.Fatal(err.Error())
	}

	secretClient := provider.ClientSet.CoreV1().Secrets(utils.CONFIG.Kubernetes.OwnNamespace)
	existingSecret, getErr := secretClient.Get(context.TODO(), utils.USERSSECRET, metav1.GetOptions{})

	secret := utils.InitSecret()
	secret.ObjectMeta.Name = utils.USERSSECRET
	secret.ObjectMeta.Namespace = utils.CONFIG.Kubernetes.OwnNamespace
	delete(secret.StringData, "exampleData") // delete example data

	// if not exist
	if existingSecret == nil || getErr != nil {
		fmt.Println("Creating new punq-auth secret ...")
		_, err := secretClient.Create(context.TODO(), &secret, kubernetes.MoCreateOptions())
		if err != nil {
			logger.Log.Error(err)
			return
		}
		fmt.Println("Created new punq-auth secret. ✅")
	}
}

func CreateAdminUser() {
	secret := kubernetes.SecretFor(utils.CONFIG.Kubernetes.OwnNamespace, utils.USERSSECRET, nil)
	if secret == nil {
		return
	}

	if secret.Data[PunqAdminIdKey] != nil {
		return
	}

	password := utils.NanoId()

	adminUser, _ := AddUser(dtos.PunqUserCreateInput{
		Email:       fmt.Sprintf("%s-%s@punq.dev", strings.ToLower(utils.RandomFirstName()), strings.ToLower(utils.RandomLastName())),
		Password:    password,
		DisplayName: "ADMIN USER",
		AccessLevel: dtos.ADMIN,
	})

	// display admin user
	displayAdminUser := adminUser
	displayAdminUser.Password = password
	utils.PrintInfo("Please store following admin user credentials in a safe place:")
	displayAdminUser.PrintToTerminalWithPwd()

	secret = kubernetes.SecretFor(utils.CONFIG.Kubernetes.OwnNamespace, utils.USERSSECRET, nil)
	strData := make(map[string]string)
	strData[PunqAdminIdKey] = adminUser.Id
	secret.StringData = strData
	kubernetes.UpdateK8sSecret(*secret, nil)
}

func ListUsers() []dtos.PunqUser {
	users := []dtos.PunqUser{}

	secret := kubernetes.SecretFor(utils.CONFIG.Kubernetes.OwnNamespace, utils.USERSSECRET, nil)
	if secret == nil {
		logger.Log.Errorf("Failed to get '%s/%s' secret.", utils.CONFIG.Kubernetes.OwnNamespace, utils.USERSSECRET)
		return users
	}

	for userId, userRaw := range secret.Data {
		if userId == PunqAdminIdKey {
			continue
		}
		user := dtos.PunqUser{}
		err := json.Unmarshal(userRaw, &user)
		if err != nil {
			logger.Log.Error("Failed to Unmarshal user '%s'.", userId)
		}
		users = append(users, user)
	}

	return users
}

func AddUser(userCreateInput dtos.PunqUserCreateInput) (*dtos.PunqUser, error) {
	secret := kubernetes.SecretFor(utils.CONFIG.Kubernetes.OwnNamespace, utils.USERSSECRET, nil)
	if secret == nil {
		return nil, errors.New(fmt.Sprintf("failed to get '%s/%s' secret", utils.CONFIG.Kubernetes.OwnNamespace, utils.USERSSECRET))
	}

	// check for duplicates
	for _, data := range secret.Data {
		userDto := &dtos.PunqUser{}
		err := json.Unmarshal(data, userDto)
		if err == nil {
			if userDto.Email == userCreateInput.Email {
				errStr := fmt.Sprintf("Duplicated email: '%s'", userCreateInput.Email)
				logger.Log.Error(errStr)
				return nil, errors.New(errStr)
			}
		}
	}

	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userCreateInput.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	userCreateInput.Password = string(hashedPassword)

	jsonData, err := json.Marshal(userCreateInput)
	if err != nil {
		errStr := fmt.Sprintf("failed marshalling userCreateInput %v", err)
		logger.Log.Error(errStr)
		return nil, errors.New(errStr)
	}

	user := dtos.PunqUser{
		Id:      utils.NanoId(),
		Created: time.Now().Format(time.RFC3339),
	}

	err = json.Unmarshal(jsonData, &user)
	if err != nil {
		errStr := fmt.Sprintf("failed unmarshalling into user %v", err)
		logger.Log.Error(errStr)
		return nil, errors.New(errStr)
	}

	rawData, err := json.Marshal(user)
	if err != nil {
		errStr := fmt.Sprintf("failed to Marshal user '%s'", user.Id)
		logger.Log.Error(errStr)
		return nil, errors.New(errStr)
	}

	if secret.StringData == nil {
		secret.StringData = make(map[string]string)
	}
	secret.StringData[user.Id] = string(rawData)

	// add user to secret
	kubernetes.UpdateK8sSecret(*secret, nil)

	return &user, nil
}

func UpdateUser(userUpdateInput dtos.PunqUser) (*dtos.PunqUser, error) {
	secret := kubernetes.SecretFor(utils.CONFIG.Kubernetes.OwnNamespace, utils.USERSSECRET, nil)
	if secret == nil {
		return nil, errors.New(fmt.Sprintf("failed to get '%s/%s' secret", utils.CONFIG.Kubernetes.OwnNamespace, utils.USERSSECRET))
	}

	user, err := GetUser(userUpdateInput.Id)
	if err != nil {
		return nil, err
	}

	// check duplicated email
	if user.Email != userUpdateInput.Email {
		findByEmail, _ := GetUserByEmail(userUpdateInput.Email)
		if findByEmail != nil && findByEmail.Id != userUpdateInput.Id {
			errStr := fmt.Sprintf("Duplicated email: '%s'", userUpdateInput.Email)
			logger.Log.Error(errStr)
			return nil, errors.New(errStr)
		}
	}

	// hash new password
	if userUpdateInput.Password != "" && user.Password != userUpdateInput.Password {
		// hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userUpdateInput.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, errors.New(err.Error())
		}
		userUpdateInput.Password = string(hashedPassword)
	}

	jsonData, err := json.Marshal(userUpdateInput)
	if err != nil {
		errStr := fmt.Sprintf("failed marshalling userCreateInput %v", err)
		logger.Log.Error(errStr)
		return nil, errors.New(errStr)
	}

	err = json.Unmarshal(jsonData, &user)
	if err != nil {
		errStr := fmt.Sprintf("failed unmarshalling into user %v", err)
		logger.Log.Error(errStr)
		return nil, errors.New(errStr)
	}

	rawData, err := json.Marshal(user)
	secret.Data[userUpdateInput.Id] = rawData

	// update user
	kubernetes.UpdateK8sSecret(*secret, nil)

	return user, nil
}

func DeleteUser(id string) error {
	secret := kubernetes.SecretFor(utils.CONFIG.Kubernetes.OwnNamespace, utils.USERSSECRET, nil)
	if secret == nil {
		return errors.New(fmt.Sprintf("failed to get '%s/%s' secret", utils.CONFIG.Kubernetes.OwnNamespace, utils.USERSSECRET))
	}

	if id == utils.USERADMIN {
		return errors.New("admin user cannot be deleted")
	}

	if secret.Data[id] != nil {
		delete(secret.Data, id)
	} else {
		return errors.New(fmt.Sprintf("USer '%s' not found.", id))
	}

	result := kubernetes.UpdateK8sSecret(*secret, nil)
	if result.Error == nil && result.Result == nil {
		// success
		result.Result = fmt.Sprintf("User %s successfully deleted.", id)
	}
	return nil
}

func GetUser(id string) (*dtos.PunqUser, error) {
	secret := kubernetes.SecretFor(utils.CONFIG.Kubernetes.OwnNamespace, utils.USERSSECRET, nil)
	if secret == nil {
		msg := fmt.Sprintf("Failed to get '%s/%s' secret.", utils.CONFIG.Kubernetes.OwnNamespace, utils.USERSSECRET)
		logger.Log.Error(msg)
		return nil, errors.New(msg)
	}

	if secret.Data[id] != nil {
		user := dtos.PunqUser{}
		err := json.Unmarshal(secret.Data[id], &user)
		if err != nil {
			msg := fmt.Sprintf("Failed to Unmarshal user '%s'.", id)
			logger.Log.Error(msg)
			return nil, errors.New(msg)
		}
		return &user, nil
	}

	return nil, errors.New("user not found")
}

func GetUserByEmail(email string) (*dtos.PunqUser, error) {
	secret := kubernetes.SecretFor(utils.CONFIG.Kubernetes.OwnNamespace, utils.USERSSECRET, nil)
	if secret == nil {
		msg := fmt.Sprintf("Failed to get '%s/%s' secret.", utils.CONFIG.Kubernetes.OwnNamespace, utils.USERSSECRET)
		logger.Log.Error(msg)
		return nil, errors.New(msg)
	}

	for userId, userRaw := range secret.Data {
		if userId == PunqAdminIdKey {
			continue
		}
		user := dtos.PunqUser{}
		err := json.Unmarshal(userRaw, &user)
		if err != nil {
			msg := fmt.Sprintf("Failed to Unmarshal user '%s'.", userId)
			logger.Log.Error(msg)
			return nil, errors.New(msg)
		}
		if user.Email == email {
			return &user, nil
		}
	}

	return nil, errors.New("user not found")
}

func GetAdmin() (*dtos.PunqUser, error) {
	secret := kubernetes.SecretFor(utils.CONFIG.Kubernetes.OwnNamespace, utils.USERSSECRET, nil)
	if secret == nil {
		err := fmt.Errorf("failed to get '%s/%s' secret", utils.CONFIG.Kubernetes.OwnNamespace, utils.USERSSECRET)
		logger.Log.Error(err)
		return nil, err
	}

	adminId := string(secret.Data[PunqAdminIdKey])

	adminUser, err := GetUser(adminId)
	if err != nil {
		return nil, err
	}

	return adminUser, nil
}

func GetGinContextUser(c *gin.Context) *dtos.PunqUser {
	if temp, exists := c.Get("user"); exists {
		user, ok := temp.(dtos.PunqUser)
		if !ok {
			utils.MalformedMessage(c, "Type Assertion failed. Expected PunqUser but received something different.")
			return nil
		}
		return &user
	}
	return nil
}
