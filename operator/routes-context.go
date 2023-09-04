package operator

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mogenius/punq/dtos"
	"github.com/mogenius/punq/logger"
	"github.com/mogenius/punq/services"
)

func InitContextRoutes(router *gin.Engine) {
	router.GET("/context/all", Auth(dtos.ADMIN), allContexts)
	router.GET("/context/:ctxId", Auth(dtos.ADMIN), getContext)
	router.DELETE("/context/:ctxId", Auth(dtos.ADMIN), deleteContext)
	router.PATCH("/context/:ctxId", Auth(dtos.ADMIN), updateContext)
}

// @Tags Context
// @Produce json
// @Success 200 {array} dtos.PunqContext
// @Router /context/all [get]
func allContexts(c *gin.Context) {
	c.JSON(http.StatusOK, services.ListContexts())
}

func getContext(c *gin.Context) {
	ctxId := c.Param("ctxId")
	result := services.GetContext(ctxId)
	if result == nil {
		NotFound(c, "Context not found.")
		return
	}
	c.JSON(http.StatusOK, result)
}

func deleteContext(c *gin.Context) {
	ctxId := c.Param("ctxId")
	result := services.DeleteContext(ctxId)
	c.JSON(http.StatusOK, result)
}

// TODO -> This is crap. validator is shit bind is shit.
func updateContext(c *gin.Context) {
	ctxId := c.Param("ctxId")
	result := services.GetContext(ctxId)
	if result == nil {
		NotFound(c, "Context not found.")
		return
	}

	var context dtos.PunqContext
	err := c.Bind(&context)
	if err != nil {
		logger.Log.Error(err.Error())
	}

	c.JSON(http.StatusOK, result)
}
