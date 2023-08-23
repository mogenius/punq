package operator

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mogenius/punq/dtos"
	"github.com/mogenius/punq/services"
)

func InitContextRoutes(router *gin.Engine) {
	router.GET("/context/all", allContexts)
	router.GET("/context/:ctxId", getContext)
}

func allContexts(c *gin.Context) {
	c.JSON(http.StatusOK, services.ListContexts())
}

func getContext(c *gin.Context) {
	isAuthorized, err := HasSufficientAccess(c, dtos.ADMIN)
	if err != nil {
		return
	}

	if isAuthorized {
		ctxId := c.Param("ctxId")
		result := services.GetContext(ctxId)
		if result == nil {
			NotFound(c, "Context not found.")
			return
		}
		c.JSON(http.StatusOK, services.GetContext(ctxId))
	}
}
