package services

import (
	"encoding/json"
	"fmt"

	"github.com/mogenius/punq/dtos"
	"github.com/mogenius/punq/kubernetes"
	"github.com/mogenius/punq/logger"
	"github.com/mogenius/punq/utils"
)

func ListUsers() []dtos.PunqUser {
	users := []dtos.PunqUser{}

	secret := kubernetes.SecretFor(utils.CONFIG.Kubernetes.OwnNamespace, utils.USERSSECRET)
	if secret == nil {
		logger.Log.Errorf("Failed to get '%s/%s' secret.", utils.CONFIG.Kubernetes.OwnNamespace, utils.USERSSECRET)
		return users
	}

	for userId, userRaw := range secret.Data {
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

	rawData, err := json.Marshal(user)
	if err != nil {
		logger.Log.Error("failed to Marshal user '%s'", user.Id)
	}
	secret.Data[user.Id] = rawData

	return kubernetes.UpdateK8sSecret(*secret)
}

func DeleteUser(id string) kubernetes.K8sWorkloadResult {
	secret := kubernetes.SecretFor(utils.CONFIG.Kubernetes.OwnNamespace, utils.USERSSECRET)
	if secret == nil {
		return kubernetes.WorkloadResultError(fmt.Sprintf("failed to get '%s/%s' secret", utils.CONFIG.Kubernetes.OwnNamespace, utils.USERSSECRET))
	}

	if id == "admin" {
		return kubernetes.WorkloadResultError("admin user cannot be deleted")
	}

	delete(secret.Data, id)

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

	for userId, userRaw := range secret.Data {
		user := dtos.PunqUser{}
		err := json.Unmarshal(userRaw, &user)
		if err != nil {
			logger.Log.Error("Failed to Unmarshal user '%s'.", userId)
		}
		if user.Id == id {
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

	for userId, userRaw := range secret.Data {
		if userId == "admin" {
			admin := dtos.PunqUser{}
			err := json.Unmarshal([]byte(userRaw), &admin)
			if err != nil {
				logger.Log.Error("Failed to Unmarshal user '%s'.", userId)
			}
			return &admin, nil
		}
	}
	return nil, fmt.Errorf("admin user not found")
}
