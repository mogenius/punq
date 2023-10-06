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

func RequireNamespace() gin.HandlerFunc {
	return func(c *gin.Context) {
		contextId := services.GetGinNamespace(c)
		if contextId == nil {
			utils.MissingHeader(c, "X-Namespace")
			c.Abort()
			return
		} else {
			c.Next()
		}
	}
}

func RequirePodName() gin.HandlerFunc {
	return func(c *gin.Context) {
		contextId := services.GetGinPodname(c)
		if contextId == nil {
			utils.MissingHeader(c, "X-Podname")
			c.Abort()
			return
		} else {
			c.Next()
		}
	}
}

func RequireContainerName() gin.HandlerFunc {
	return func(c *gin.Context) {
		contextId := services.GetGinContainername(c)
		if contextId == nil {
			utils.MissingHeader(c, "X-Container")
			c.Abort()
			return
		} else {
			c.Next()
		}
	}
}
