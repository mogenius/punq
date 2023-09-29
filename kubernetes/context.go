package kubernetes

import (
	"errors"
	"fmt"
	"os"

	"github.com/mogenius/punq/dtos"
	"github.com/mogenius/punq/utils"
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
