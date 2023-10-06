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
	return kubernetes.ListAllContexts()
}

func AddContext(ctx dtos.PunqContext) (*dtos.PunqContext, error) {
	secret := kubernetes.SecretFor(utils.CONFIG.Kubernetes.OwnNamespace, utils.CONTEXTSSECRET, nil)
	if secret == nil {
		msg := fmt.Sprintf("failed to get '%s/%s' secret", utils.CONFIG.Kubernetes.OwnNamespace, utils.CONTEXTSSECRET)
		logger.Log.Error(msg)
		return nil, errors.New(msg)
	}

	// check if context already exists
	currentCtxs := ListContexts()
	for _, aCtx := range currentCtxs {
		if aCtx.ContextHash == ctx.ContextHash {
			return nil, fmt.Errorf("context '%s' already exists", ctx.Name)
		}
	}

	rawData, err := json.Marshal(ctx)
	if err != nil {
		msg := fmt.Sprintf("failed to Marshal context '%s'", ctx.Id)
		logger.Log.Error(msg)
		return nil, errors.New(msg)
	}
	secret.Data[ctx.Id] = rawData

	workloadResult := kubernetes.UpdateK8sSecret(*secret, nil)
	if workloadResult.Error != nil {
		return nil, err
	}

	// Update LocalContextArray
	kubernetes.ContextAddOne(ctx)

	return &ctx, nil
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
	kubernetes.ContextAddMany(ListContexts())

	return nil, errors.New(fmt.Sprintf("%v", workloadResult.Error))
}

func DeleteContext(id string) (interface{}, error) {
	secret := kubernetes.SecretFor(utils.CONFIG.Kubernetes.OwnNamespace, utils.CONTEXTSSECRET, nil)
	if secret == nil {
		msg := fmt.Sprintf("failed to get '%s/%s' secret", utils.CONFIG.Kubernetes.OwnNamespace, utils.CONTEXTSSECRET)
		return nil, errors.New(msg)
	}

	if id == utils.CONTEXTOWN {
		return nil, errors.New("own-context cannot be deleted")
	}

	if secret.Data[id] != nil {
		delete(secret.Data, id)
	} else {
		msg := fmt.Sprintf("Context '%s' not found.", id)
		return nil, errors.New(msg)
	}

	workloadResult := kubernetes.UpdateK8sSecret(*secret, nil)
	if workloadResult.Error == nil && workloadResult.Result != nil {
		// success
		workloadResult.Result = fmt.Sprintf("Context %s successfully deleted.", id)

		// Update LocalContextArray
		kubernetes.ContextAddMany(ListContexts())

		return workloadResult.Result, nil
	}

	return nil, fmt.Errorf("%v", workloadResult.Error)
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
	if contextId := c.GetHeader("X-Context-Id"); contextId != "" {
		return &contextId
	}
	return nil
}

func GetGinNamespace(c *gin.Context) *string {
	if namespace := c.GetHeader("X-Namespace"); namespace != "" {
		return &namespace
	}
	return nil
}

func GetGinPodname(c *gin.Context) *string {
	if podname := c.GetHeader("X-Podname"); podname != "" {
		return &podname
	}
	return nil
}

func GetGinContainername(c *gin.Context) *string {
	if container := c.GetHeader("X-Container"); container != "" {
		return &container
	}
	return nil
}

// func GetGinContextContexts(c *gin.Context) *[]dtos.PunqContext {
// 	if contextArray, exists := c.Get("contexts"); exists {
// 		contexts, ok := contextArray.([]dtos.PunqContext)
// 		if !ok {
// 			utils.MalformedMessage(c, "Type Assertion failed. Expected Array of PunqContext but received something different.")
// 			return nil
// 		}
// 		return &contexts
// 	}
// 	return nil
// }
