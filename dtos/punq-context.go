package dtos

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/jedib0t/go-pretty/table"
	"github.com/mogenius/punq/utils"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

type PunqContext struct {
	Id          string       `json:"id" validate:"required"`
	Name        string       `json:"name" validate:"required"`
	ContextHash string       `json:"contextHash" validate:"required"`
	Context     string       `json:"context" validate:"required"`
	Provider    string       `json:"provider" validate:"required"`
	Reachable   bool         `json:"reachable" validate:"required"`
	Access      []PunqAccess `json:"access" validate:"required"`
}

func CreateContext(id string, name string, context string, provider string, access []PunqAccess) PunqContext {
	ctx := PunqContext{}

	ctx.Name = name

	if id == "" {
		ctx.Id = utils.NanoId()
	} else {
		ctx.Id = id
	}

	ctx.Context = context
	ctx.ContextHash = utils.HashString(context)

	if provider == "" {
		ctx.Provider = "UNKNOWN"
	} else {
		ctx.Provider = provider
	}

	if len(access) > 0 {
		ctx.Access = access
	} else {
		ctx.Access = []PunqAccess{}
	}

	return ctx
}

func (c *PunqContext) AddAccess(userId string, accessLevel AccessLevel) {
	for _, access := range c.Access {
		if access.UserId == userId {
			// UPDATE EXISTING
			access.Level = accessLevel
			return
		}
	}
	// CREATE NEW
	c.Access = append(c.Access, PunqAccess{
		UserId: userId,
		Level:  accessLevel,
	})
}

func (c *PunqContext) RemoveAccess(userId string) {
	resultingArray := []PunqAccess{}
	for _, access := range c.Access {
		if access.UserId != userId {
			resultingArray = append(resultingArray, access)
		}
	}
	c.Access = resultingArray
}

func (c *PunqContext) PrintToTerminal() {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"ID", "Name", "Access", "Hash"})
	accessStr := "*"
	accessEntries := []string{}
	for _, access := range c.Access {
		accessEntries = append(accessEntries, fmt.Sprintf("%s (%d)", access.UserId, access.Level))
	}
	if len(accessEntries) > 0 {
		accessStr = strings.Join(accessEntries, ", ")
	}
	t.AppendRow(
		table.Row{c.Id, c.Name, accessStr, c.ContextHash},
	)
	t.Render()
}

func ListContextsToTerminal(contexts []PunqContext) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "ID", "Name", "Reachable", "Provider", "Access"})
	for index, context := range contexts {
		accessStr := "*"
		accessEntries := []string{}
		for _, access := range context.Access {
			accessEntries = append(accessEntries, fmt.Sprintf("%s (%d)", access.UserId, access.Level))
		}
		if len(accessEntries) > 0 {
			accessStr = strings.Join(accessEntries, ", ")
		}
		t.AppendRow(
			table.Row{index + 1, context.Id, context.Name, utils.StatusEmoji(context.Reachable), context.Provider, accessStr},
		)
	}
	t.Render()
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

func ParseConfigToPunqContexts(data []byte) ([]PunqContext, error) {
	result := []PunqContext{}
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
		result = append(result, CreateContext("", contextName, string(configBytes), "", []PunqAccess{}))
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})

	return result, nil
}

func ParseCurrentContextConfigToPunqContext(data []byte) (PunqContext, error) {
	result := PunqContext{}
	config, err := clientcmd.Load(data)
	if err != nil {
		fmt.Printf("Failed to load kubeconfig: %v\n", err)
		return result, err
	}
	for contextName := range config.Contexts {
		if contextName != config.CurrentContext {
			continue
		}

		aConfig, err := ExtractSingleConfigFromContext(config, contextName)
		if err != nil {
			return result, err
		}
		configBytes, err := clientcmd.Write(*aConfig)
		if err != nil {
			return result, err
		}
		result = CreateContext("", contextName, string(configBytes), "", []PunqAccess{})
	}

	if result.Id == "" {
		return result, fmt.Errorf("current context not found in source kubeconfig")
	}

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
		CreateContext("", contextName, "", "", []PunqAccess{})
		t.AppendRow(
			table.Row{contextName, context.Cluster, context.AuthInfo, config.Clusters[context.Cluster].Server},
		)
	}
	t.Render()
}
