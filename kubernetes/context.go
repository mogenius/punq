package kubernetes

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/mogenius/punq/dtos"
	"github.com/mogenius/punq/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var allContexts []dtos.PunqContext = []dtos.PunqContext{}

func ContextForId(id string) *dtos.PunqContext {
	for _, ctx := range allContexts {
		if ctx.Id == id {
			return &ctx
		}
	}
	return nil
}

func ContextAddOne(ctx dtos.PunqContext) {
	if ContextForId(ctx.Id) != nil {
		return
	}
	allContexts = append(allContexts, ctx)
	contextWrite(ctx.Id)
}

func ContextAddMany(ctxs []dtos.PunqContext) {
	for _, ctx := range ctxs {
		ContextAddOne(ctx)
	}
}

func ContextList() []dtos.PunqContext {
	return allContexts
}

func ContextFlag(id *string) string {
	if id == nil {
		return ""
	}
	return fmt.Sprintf("--kubeconfig=%s.yaml", *id)
}

func contextWrite(id string) error {
	if !utils.CONFIG.Kubernetes.RunInCluster {
		return nil
	}

	ctx := ContextForId(id)
	if ctx == nil {
		return errors.New("context not found")
	}

	// write ctx.context into file
	err := os.WriteFile(fmt.Sprintf("%s.yaml", ctx.Id), []byte(ctx.Context), 0644)
	if err != nil {
		_ = os.Remove(ctx.Id)
		return err
	}
	return nil
}

func CheckContext(ctx dtos.PunqContext) (bool, dtos.KubernetesProvider, error) {
	configFromString, err := clientcmd.NewClientConfigFromBytes([]byte(ctx.Context))
	if err != nil {
		return false, dtos.UNKNOWN, err
	}

	config, err := configFromString.ClientConfig()
	if err != nil {
		return false, dtos.UNKNOWN, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return false, dtos.UNKNOWN, err
	}

	nodeList, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return false, dtos.UNKNOWN, err
	}

	provider, err := GuessCluserProviderFromNodeList(nodeList)
	if err != nil {
		return false, dtos.UNKNOWN, err
	} else {
		return true, provider, nil
	}
}
