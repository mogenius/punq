package services

import (
	"encoding/json"
	"fmt"

	"github.com/mogenius/punq/dtos"
	"github.com/mogenius/punq/kubernetes"
	"github.com/mogenius/punq/logger"
	"github.com/mogenius/punq/utils"
)

func ListContexts() []dtos.PunqContext {
	contexts := []dtos.PunqContext{}

	secret := kubernetes.SecretFor(utils.CONFIG.Kubernetes.OwnNamespace, utils.CONTEXTSSECRET)
	if secret == nil {
		logger.Log.Errorf("Failed to get '%s/%s' secret.", utils.CONFIG.Kubernetes.OwnNamespace, utils.CONTEXTSSECRET)
		return contexts
	}

	for ctxId, contextRaw := range secret.Data {
		ctx := dtos.PunqContext{}
		err := json.Unmarshal(contextRaw, &ctx)
		if err != nil {
			logger.Log.Error("Failed to Unmarshal context '%s'.", ctxId)
		}
		contexts = append(contexts, ctx)
	}

	return contexts
}

func AddContext(ctx dtos.PunqContext) kubernetes.K8sWorkloadResult {
	secret := kubernetes.SecretFor(utils.CONFIG.Kubernetes.OwnNamespace, utils.CONTEXTSSECRET)
	if secret == nil {
		return kubernetes.WorkloadResultError(fmt.Sprintf("failed to get '%s/%s' secret", utils.CONFIG.Kubernetes.OwnNamespace, utils.CONTEXTSSECRET))
	}

	rawData, err := json.Marshal(ctx)
	if err != nil {
		logger.Log.Error("failed to Marshal context '%s'", ctx.Id)
	}
	secret.Data[ctx.Id] = rawData

	return kubernetes.UpdateK8sSecret(*secret)
}

func DeleteContext(id string) kubernetes.K8sWorkloadResult {
	secret := kubernetes.SecretFor(utils.CONFIG.Kubernetes.OwnNamespace, utils.CONTEXTSSECRET)
	if secret == nil {
		return kubernetes.WorkloadResultError(fmt.Sprintf("failed to get '%s/%s' secret", utils.CONFIG.Kubernetes.OwnNamespace, utils.CONTEXTSSECRET))
	}

	if id == utils.CONTEXTOWN {
		return kubernetes.WorkloadResultError("own context cannot be deleted")
	}

	if secret.Data[id] != nil {
		delete(secret.Data, id)
	} else {
		return kubernetes.WorkloadResultError(fmt.Sprintf("Context '%s' not found.", id))
	}

	result := kubernetes.UpdateK8sSecret(*secret)
	if result.Error == nil && result.Result == nil {
		// success
		result.Result = fmt.Sprintf("Context %s successfuly deleted.", id)
	}
	return result
}

func GetContext(id string) *dtos.PunqContext {
	secret := kubernetes.SecretFor(utils.CONFIG.Kubernetes.OwnNamespace, utils.CONTEXTSSECRET)
	if secret == nil {
		logger.Log.Errorf("Failed to get '%s/%s' secret.", utils.CONFIG.Kubernetes.OwnNamespace, utils.CONTEXTSSECRET)
		return nil
	}

	for ctxId, ctxRaw := range secret.Data {
		ctx := dtos.PunqContext{}
		err := json.Unmarshal(ctxRaw, &ctx)
		if err != nil {
			logger.Log.Error("Failed to Unmarshal context '%s'.", ctxId)
		}
		if ctx.Id == id {
			return &ctx
		}
	}

	return nil
}

func GetOwnContext() (*dtos.PunqContext, error) {
	secret := kubernetes.SecretFor(utils.CONFIG.Kubernetes.OwnNamespace, utils.CONTEXTSSECRET)
	if secret == nil {
		err := fmt.Errorf("failed to get '%s/%s' secret", utils.CONFIG.Kubernetes.OwnNamespace, utils.CONTEXTSSECRET)
		logger.Log.Error(err)
		return nil, err
	}

	for ctxId, contextRaw := range secret.Data {
		if ctxId == utils.CONTEXTOWN {
			ownContext := dtos.PunqContext{}
			err := json.Unmarshal([]byte(contextRaw), &ownContext)
			if err != nil {
				logger.Log.Error("Failed to Unmarshal context '%s'.", ctxId)
			}
			return &ownContext, nil
		}
	}
	return nil, fmt.Errorf("%s not found", utils.CONTEXTOWN)
}
