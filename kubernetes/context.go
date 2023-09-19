package kubernetes

import (
	"errors"
	"fmt"
	"os"

	"github.com/mogenius/punq/dtos"
	"github.com/mogenius/punq/logger"
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
		logger.Log.Error("context already exists")
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

func ContextUpdateLocalCache(ctxs []dtos.PunqContext) {
	if len(ctxs) > 0 {
		allContexts = ctxs
	}
	for _, ctx := range ctxs {
		contextWrite(ctx.Id)
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
