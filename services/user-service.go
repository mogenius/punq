package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mogenius/punq/structs"
	"time"

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
	CreateAdminUser()
}

func CreateUserSecret() {
	provider := kubernetes.NewKubeProvider()
	if provider == nil {
		logger.Log.Fatal("Failed to load kubeprovider.")
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
		fmt.Println("Created new punq-users secret. ✅")
	}
}

func CreateAdminUser() {
	secret := kubernetes.SecretFor(utils.CONFIG.Kubernetes.OwnNamespace, utils.USERSSECRET)
	if secret == nil {
		return
	}

	if secret.Data[PunqAdminIdKey] != nil {
		return
	}

	password := utils.NanoId()

	adminUser := dtos.PunqUser{
		Id:          utils.NanoId(),
		Email:       "admin@punq.dev",
		Password:    password,
		DisplayName: "Admin User",
		AccessLevel: dtos.ADMIN,
		Created:     time.Now().Format(time.RFC3339),
	}

	AddUser(adminUser)
	structs.PrettyPrint(adminUser)
	secret = kubernetes.SecretFor(utils.CONFIG.Kubernetes.OwnNamespace, utils.USERSSECRET)
	strData := make(map[string]string)
	strData[PunqAdminIdKey] = adminUser.Id
	secret.StringData = strData
	kubernetes.UpdateK8sSecret(*secret)
}

func ListUsers() []dtos.PunqUser {
	users := []dtos.PunqUser{}

	secret := kubernetes.SecretFor(utils.CONFIG.Kubernetes.OwnNamespace, utils.USERSSECRET)
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

func AddUser(user dtos.PunqUser) kubernetes.K8sWorkloadResult {
	secret := kubernetes.SecretFor(utils.CONFIG.Kubernetes.OwnNamespace, utils.USERSSECRET)
	if secret == nil {
		return kubernetes.WorkloadResultError(fmt.Sprintf("failed to get '%s/%s' secret", utils.CONFIG.Kubernetes.OwnNamespace, utils.USERSSECRET))
	}

	// check for duplicates
	for _, data := range secret.Data {
		userDto := &dtos.PunqUser{}
		err := json.Unmarshal(data, userDto)
		if err == nil {
			if userDto.Email == user.Email {
				errStr := fmt.Sprintf("Duplicated email: '%s'", user.Email)
				logger.Log.Error(errStr)
				return kubernetes.WorkloadResultError(errStr)
			}
		}
	}

	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return kubernetes.WorkloadResultError(err.Error())
	}
	user.Password = string(hashedPassword)

	rawData, err := json.Marshal(user)
	if err != nil {
		logger.Log.Error("failed to Marshal user '%s'", user.Id)
	}

	if secret.StringData == nil {
		secret.StringData = make(map[string]string)
	}
	secret.StringData[user.Id] = string(rawData)

	return kubernetes.UpdateK8sSecret(*secret)
}

func UpdateUser(user dtos.PunqUser) kubernetes.K8sWorkloadResult {
	secret := kubernetes.SecretFor(utils.CONFIG.Kubernetes.OwnNamespace, utils.USERSSECRET)
	if secret == nil {
		return kubernetes.WorkloadResultError(fmt.Sprintf("failed to get '%s/%s' secret", utils.CONFIG.Kubernetes.OwnNamespace, utils.USERSSECRET))
	}

	rawData, err := json.Marshal(user)
	if err != nil {
		logger.Log.Error("failed to Marshal user '%s'", user.Id)
	}
	userObj := GetUser(user.Id)
	if userObj.Password != user.Password {
		// hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return kubernetes.WorkloadResultError(err.Error())
		}
		user.Password = string(hashedPassword)
	}
	secret.Data[user.Id] = rawData

	return kubernetes.UpdateK8sSecret(*secret)
}

func DeleteUser(id string) kubernetes.K8sWorkloadResult {
	secret := kubernetes.SecretFor(utils.CONFIG.Kubernetes.OwnNamespace, utils.USERSSECRET)
	if secret == nil {
		return kubernetes.WorkloadResultError(fmt.Sprintf("failed to get '%s/%s' secret", utils.CONFIG.Kubernetes.OwnNamespace, utils.USERSSECRET))
	}

	if id == utils.USERADMIN {
		return kubernetes.WorkloadResultError("admin user cannot be deleted")
	}

	if secret.Data[id] != nil {
		delete(secret.Data, id)
	} else {
		return kubernetes.WorkloadResultError(fmt.Sprintf("USer '%s' not found.", id))
	}

	result := kubernetes.UpdateK8sSecret(*secret)
	if result.Error == nil && result.Result == nil {
		// success
		result.Result = fmt.Sprintf("User %s successfuly deleted.", id)
	}
	return result
}

func GetUser(id string) *dtos.PunqUser {
	secret := kubernetes.SecretFor(utils.CONFIG.Kubernetes.OwnNamespace, utils.USERSSECRET)
	if secret == nil {
		logger.Log.Errorf("Failed to get '%s/%s' secret.", utils.CONFIG.Kubernetes.OwnNamespace, utils.USERSSECRET)
		return nil
	}

	if secret.Data[id] != nil {
		user := dtos.PunqUser{}
		err := json.Unmarshal(secret.Data[id], &user)
		if err != nil {
			logger.Log.Error("Failed to Unmarshal user '%s'.", id)
		}
		return &user
	}

	return nil
}

func GetUserByEmail(email string) *dtos.PunqUser {
	secret := kubernetes.SecretFor(utils.CONFIG.Kubernetes.OwnNamespace, utils.USERSSECRET)
	if secret == nil {
		logger.Log.Errorf("Failed to get '%s/%s' secret.", utils.CONFIG.Kubernetes.OwnNamespace, utils.USERSSECRET)
		return nil
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
		if user.Email == email {
			return &user
		}
	}

	return nil
}

func GetAdmin() (*dtos.PunqUser, error) {
	secret := kubernetes.SecretFor(utils.CONFIG.Kubernetes.OwnNamespace, utils.USERSSECRET)
	if secret == nil {
		err := fmt.Errorf("failed to get '%s/%s' secret", utils.CONFIG.Kubernetes.OwnNamespace, utils.USERSSECRET)
		logger.Log.Error(err)
		return nil, err
	}

	adminId := string(secret.Data[PunqAdminIdKey])

	adminUser := GetUser(adminId)
	if adminUser != nil {
		return adminUser, nil
	}

	return nil, fmt.Errorf("admin user not found")
}
