package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/jedib0t/go-pretty/table"
	"github.com/mogenius/punq/dtos"
	"github.com/mogenius/punq/kubernetes"
	"github.com/mogenius/punq/logger"
	"github.com/mogenius/punq/utils"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

func ListContexts() []dtos.PunqContext {
	return kubernetes.ListAllContexts()
}

func ExtractSingleConfigFromContext(config *api.Config, contextName string) (*api.Config, error) {
	context, contextExists := config.Contexts[contextName]
	if !contextExists {
		return nil, fmt.Errorf("Context %s not found in source kubeconfig\n", contextName)
	}
	cluster, clusterExists := config.Clusters[context.Cluster]
	if !clusterExists {
		return nil, fmt.Errorf("Cluster %s for context %s not found in source kubeconfig\n", context.Cluster, contextName)
	}
	authInfo, userExists := config.AuthInfos[context.AuthInfo]
	if !userExists {
		return nil, fmt.Errorf("User %s for context %s not found in source kubeconfig\n", context.AuthInfo, contextName)
	}

	singleConfig := api.NewConfig()
	singleConfig.APIVersion = config.APIVersion
	singleConfig.Kind = config.Kind
	singleConfig.CurrentContext = contextName
	singleConfig.Contexts = map[string]*api.Context{contextName: context}
	singleConfig.Clusters = map[string]*api.Cluster{context.Cluster: cluster}
	singleConfig.AuthInfos = map[string]*api.AuthInfo{context.AuthInfo: authInfo}

	return singleConfig, nil
}

func WriteSingleConfigFileFromContext(config *api.Config, contextName string) error {
	fileName := fmt.Sprintf("%s.yaml", contextName)

	newConfig, err := ExtractSingleConfigFromContext(config, contextName)
	if err != nil {
		return err
	}

	err = clientcmd.WriteToFile(*newConfig, fileName)
	if err != nil {
		return fmt.Errorf("Failed to write kubeconfig: %s - %s\n", fileName, err.Error())
	}

	fmt.Printf("Successfully extracted kubeconfig: %s\n", fileName)
	return nil
}

func ParseConfigToPunqContexts(data []byte) ([]dtos.PunqContext, error) {
	result := []dtos.PunqContext{}
	config, err := clientcmd.Load(data)
	if err != nil {
		fmt.Printf("Failed to load kubeconfig: %v\n", err)
		return result, err
	}
	for contextName := range config.Contexts {
		aConfig, err := ExtractSingleConfigFromContext(config, contextName)
		if err != nil {
			return result, err
		}
		configBytes, err := clientcmd.Write(*aConfig)
		if err != nil {
			return result, err
		}
		result = append(result, dtos.CreateContext("", contextName, string(configBytes), "", []dtos.PunqAccess{}))
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})

	return result, nil
}

func PrintAllContextFromConfig(config *api.Config) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetAutoIndex(true)
	t.SetAllowedColumnLengths([]int{30, 30, 30, 50})
	t.AppendHeader(table.Row{"Context", "Cluster", "User", "Server"})
	t.AppendRow(
		table.Row{"ALL CONTEXTS", "*", "*", "*"},
	)
	for contextName, context := range config.Contexts {
		dtos.CreateContext("", contextName, "", "", []dtos.PunqAccess{})
		t.AppendRow(
			table.Row{contextName, context.Cluster, context.AuthInfo, config.Clusters[context.Cluster].Server},
		)
	}
	t.Render()
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
	if workloadResult.Result != nil {
		return workloadResult.Result.(*dtos.PunqContext), nil
	}

	// Update LocalContextArray
	kubernetes.ContextUpdateLocalCache(ListContexts())

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
	kubernetes.ContextUpdateLocalCache(ListContexts())

	return nil, errors.New(fmt.Sprintf("%v", workloadResult.Error))
}

func DeleteContext(id string) (interface{}, error) {
	secret := kubernetes.SecretFor(utils.CONFIG.Kubernetes.OwnNamespace, utils.CONTEXTSSECRET, nil)
	if secret == nil {
		msg := fmt.Sprintf("failed to get '%s/%s' secret", utils.CONFIG.Kubernetes.OwnNamespace, utils.CONTEXTSSECRET)
		return nil, errors.New(msg)
	}

	if id == utils.CONTEXTOWN {
		msg := fmt.Sprintf("own context cannot be deleted")
		return nil, errors.New(msg)
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
		workloadResult.Result = fmt.Sprintf("Context %s successfuly deleted.", id)

		kubernetes.ContextUpdateLocalCache(ListContexts())

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
