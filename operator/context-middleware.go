package operator

import (
	"github.com/gin-gonic/gin"
	"github.com/mogenius/punq/services"
	"github.com/mogenius/punq/utils"
)

func RequireContextId() gin.HandlerFunc {
	return func(c *gin.Context) {
		contextId := services.GetGinContextId(c)
		if contextId == nil {
			utils.MissingHeader(c, "X-Context-Id")
			c.Abort()
			return
		} else {
			c.Next()
		}
	}
}
