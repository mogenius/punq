package services

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/mogenius/punq/dtos"
	"github.com/mogenius/punq/kubernetes"
	"github.com/mogenius/punq/logger"
	"github.com/mogenius/punq/utils"
)

func ListContexts() []dtos.PunqContext {
	contexts := []dtos.PunqContext{}

	secret := kubernetes.SecretFor(utils.CONFIG.Kubernetes.OwnNamespace, utils.CONTEXTSSECRET, nil)
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

	kubernetes.AllContexts = contexts

	return contexts
}

func AddContext(ctx dtos.PunqContext) (interface{}, error) {
	secret := kubernetes.SecretFor(utils.CONFIG.Kubernetes.OwnNamespace, utils.CONTEXTSSECRET, nil)
	if secret == nil {
		msg := fmt.Sprintf("failed to get '%s/%s' secret", utils.CONFIG.Kubernetes.OwnNamespace, utils.CONTEXTSSECRET)
		logger.Log.Error(msg)
		return nil, errors.New(msg)
	}

	rawData, err := json.Marshal(ctx)
	if err != nil {
		msg := fmt.Sprintf("failed to Marshal context '%s'", ctx.Id)
		logger.Log.Error(msg)
		return nil, errors.New(msg)
	}
	secret.Data[ctx.Id] = rawData

	workloadResult := kubernetes.UpdateK8sSecret(*secret, nil)
	if workloadResult.Result != nil {
		return workloadResult.Result, nil
	}

	// Update LocalContextArray
	ListContexts()

	return nil, errors.New(fmt.Sprintf("%v", workloadResult.Error))
}

func UpdateContext(ctx dtos.PunqContext) (interface{}, error) {
	secret := kubernetes.SecretFor(utils.CONFIG.Kubernetes.OwnNamespace, utils.CONTEXTSSECRET, nil)
	if secret == nil {
		msg := fmt.Sprintf("failed to get '%s/%s' secret", utils.CONFIG.Kubernetes.OwnNamespace, utils.CONTEXTSSECRET)
		logger.Log.Error(msg)
		return nil, errors.New(msg)
	}

	rawData, err := json.Marshal(ctx)
	if err != nil {
		msg := fmt.Sprintf("failed to Marshal context '%s'", ctx.Id)
		logger.Log.Error(msg)
		return nil, errors.New(msg)
	}
	secret.Data[ctx.Id] = rawData

	workloadResult := kubernetes.UpdateK8sSecret(*secret, nil)
	if workloadResult.Result != nil {
		return workloadResult.Result, nil
	}

	// Update LocalContextArray
	ListContexts()

	return nil, errors.New(fmt.Sprintf("%v", workloadResult.Error))
}

func DeleteContext(id string) (interface{}, error) {
	secret := kubernetes.SecretFor(utils.CONFIG.Kubernetes.OwnNamespace, utils.CONTEXTSSECRET, nil)
	if secret == nil {
		msg := fmt.Sprintf("failed to get '%s/%s' secret", utils.CONFIG.Kubernetes.OwnNamespace, utils.CONTEXTSSECRET)
		logger.Log.Error(msg)
		return nil, errors.New(msg)
	}

	if id == utils.CONTEXTOWN {
		msg := fmt.Sprintf("own context cannot be deleted")
		logger.Log.Error(msg)
		return nil, errors.New(msg)
	}

	if secret.Data[id] != nil {
		delete(secret.Data, id)
	} else {
		msg := fmt.Sprintf("Context '%s' not found.", id)
		logger.Log.Error(msg)
		return nil, errors.New(msg)
	}

	workloadResult := kubernetes.UpdateK8sSecret(*secret, nil)
	if workloadResult.Error == nil && workloadResult.Result == nil {
		// success
		workloadResult.Result = fmt.Sprintf("Context %s successfuly deleted.", id)
		return workloadResult.Result, nil
	}

	// Update LocalContextArray
	ListContexts()

	return nil, errors.New(fmt.Sprintf("%v", workloadResult.Error))
}

func GetContext(id string) (*dtos.PunqContext, error) {
	secret := kubernetes.SecretFor(utils.CONFIG.Kubernetes.OwnNamespace, utils.CONTEXTSSECRET, nil)
	if secret == nil {
		msg := fmt.Sprintf("Failed to get '%s/%s' secret.", utils.CONFIG.Kubernetes.OwnNamespace, utils.CONTEXTSSECRET)
		logger.Log.Error(msg)
		return nil, errors.New(msg)
	}

	for ctxId, ctxRaw := range secret.Data {
		ctx := dtos.PunqContext{}
		err := json.Unmarshal(ctxRaw, &ctx)
		if err != nil {
			logger.Log.Error("Failed to Unmarshal context '%s'.", ctxId)
		}
		if ctx.Id == id {
			return &ctx, nil
		}
	}

	msg := fmt.Sprintf("context not found")
	logger.Log.Error(msg)
	return nil, errors.New(msg)
}

func GetOwnContext() (*dtos.PunqContext, error) {
	secret := kubernetes.SecretFor(utils.CONFIG.Kubernetes.OwnNamespace, utils.CONTEXTSSECRET, nil)
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

func GetGinContextId(c *gin.Context) *string {
	if contextId := c.GetString("context-id"); contextId != "" {
		return &contextId
	}
	return nil
}
