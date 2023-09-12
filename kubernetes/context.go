package kubernetes

import "github.com/mogenius/punq/dtos"

var AllContexts []dtos.PunqContext = []dtos.PunqContext{}

func ContextForId(id string) *dtos.PunqContext {
	for _, ctx := range AllContexts {
		if ctx.Id == id {
			return &ctx
		}
	}
	return nil
}
