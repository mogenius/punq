package operator

import (
	"github.com/mogenius/punq/kubernetes"
	"github.com/mogenius/punq/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mogenius/punq/dtos"
	"github.com/mogenius/punq/logger"
	"github.com/mogenius/punq/services"
)

func InitContextRoutes(router *gin.Engine) {

	contextRoutes := router.Group("/context", Auth(dtos.ADMIN))
	{
		contextRoutes.GET("/all", Auth(dtos.ADMIN), allContexts)
		contextRoutes.GET("/info", Auth(dtos.ADMIN), RequireContextId(), getInfoContexts)
		contextRoutes.GET("/:ctxId", Auth(dtos.ADMIN), getContext)
		contextRoutes.DELETE("/:ctxId", Auth(dtos.ADMIN), deleteContext)
		contextRoutes.PATCH("/:ctxId", Auth(dtos.ADMIN), updateContext)
	}

}

// @Tags Context
// @Produce json
// @Success 200 {array} dtos.PunqContext
// @Router /context/all [get]
// @Security Bearer
func allContexts(c *gin.Context) {
	c.JSON(http.StatusOK, services.ListContexts())
}

// @Tags Context
// @Produce json
// @Success 200 {object} dtos.ClusterInfoDto
// @Router /context/info [get]
// @Security Bearer
func getInfoContexts(c *gin.Context) {
	c.JSON(http.StatusOK, kubernetes.ClusterInfo(services.GetGinContextId(c)))
}

// @Tags Context
// @Produce json
// @Success 200 {object} dtos.PunqContext
// @Router /context/{ctxId} [get]
// @Param ctxId path string false  "ctxId of the context-id"
// @Security Bearer
func getContext(c *gin.Context) {
	ctxId := c.Param("ctxId")

	result, _ := services.GetContext(ctxId)
	if result == nil {
		utils.NotFound(c, "Context not found.")
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Tags Context
// @Produce json
// @Success 200 {object} dtos.PunqContext
// @Router /context/{ctxId} [delete]
// @Param ctxId path string false  "ID of the context"
// @Security Bearer
func deleteContext(c *gin.Context) {
	ctxId := c.Param("ctxId")

	result, err := services.DeleteContext(ctxId)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

// TODO -> This is crap. validator is shit bind is shit.
func updateContext(c *gin.Context) {
	ctxId := c.Param("ctxId")

	result, _ := services.GetContext(ctxId)
	if result == nil {
		utils.NotFound(c, "Context not found.")
		return
	}

	var context dtos.PunqContext
	err := c.Bind(&context)
	if err != nil {
		logger.Log.Error(err.Error())
		utils.MalformedMessage(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}
