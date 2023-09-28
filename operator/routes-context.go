package operator

import (
	"fmt"
	"net/http"
	"os"

	"github.com/mogenius/punq/kubernetes"
	"github.com/mogenius/punq/logger"
	"github.com/mogenius/punq/utils"

	"github.com/gin-gonic/gin"
	"github.com/mogenius/punq/dtos"
	"github.com/mogenius/punq/services"
)

func InitContextRoutes(router *gin.Engine) {

	contextRoutes := router.Group("/context", Auth(dtos.ADMIN))
	{
		contextRoutes.GET("/all", Auth(dtos.ADMIN), allContexts)
		contextRoutes.GET("/info", Auth(dtos.ADMIN), RequireContextId(), getInfoContexts)
		contextRoutes.GET("", Auth(dtos.ADMIN), RequireContextId(), getContext)
		contextRoutes.DELETE("", Auth(dtos.ADMIN), RequireContextId(), deleteContext)
		contextRoutes.POST("/validate-config", Auth(dtos.ADMIN), validateConfig)
		contextRoutes.POST("", Auth(dtos.ADMIN), addContext)
	}
}

// @Tags Context
// @Produce json
// @Success 200 {array} dtos.PunqContext
// @Router /backend/context/all [get]
// @Security Bearer
func allContexts(c *gin.Context) {
	c.JSON(http.StatusOK, services.ListContexts())
}

// @Tags Context
// @Produce json
// @Success 200 {object} dtos.ClusterInfoDto
// @Router /backend/context/info [get]
// @Param string header string true "X-Context-Id"
// @Security Bearer
func getInfoContexts(c *gin.Context) {
	c.JSON(http.StatusOK, kubernetes.ClusterInfo(services.GetGinContextId(c)))
}

// @Tags Context
// @Produce json
// @Success 200 {object} dtos.PunqContext
// @Router /backend/context [get]
// @Param string header string true "X-Context-Id"
// @Security Bearer
func getContext(c *gin.Context) {
	ctxId := services.GetGinContextId(c)

	if ctxId != nil {
		result, _ := services.GetContext(*ctxId)
		if result == nil {
			utils.NotFound(c, "Context not found.")
			return
		}
		c.JSON(http.StatusOK, result)
	} else {
		utils.MalformedMessage(c, "No context-id found.")
		return
	}
}

// @Tags Context
// @Produce json
// @Success 200 {object} dtos.PunqContext
// @Router /backend/context [delete]
// @Param string header string true "X-Context-Id"
// @Security Bearer
func deleteContext(c *gin.Context) {
	ctxId := services.GetGinContextId(c)

	if ctxId != nil {
		result, err := services.DeleteContext(*ctxId)
		if err != nil {
			utils.MalformedMessage(c, err.Error())
			return
		}

		c.JSON(http.StatusOK, result)
	} else {
		utils.MalformedMessage(c, "No context-id found.")
		return
	}
}

// @Tags Context
// @Produce json
// @Success 200 {array} dtos.PunqContext
// @Router /backend/context/validate-config [post]
// @Security Bearer
func validateConfig(c *gin.Context) {
	tempFilename := fmt.Sprintf("%s.yaml", utils.NanoId())

	// SAVE temp file
	file, err := c.FormFile("file")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to get the file",
		})
		return
	}
	if err := c.SaveUploadedFile(file, tempFilename); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to save the file",
		})
		return
	}

	// READ
	dataBytes, err := os.ReadFile(tempFilename)
	if err != nil {
		logger.Log.Errorf("Error reading file '%s': %s", tempFilename, err.Error())
	}

	// PARSE
	contexts, err := services.ParseConfigToPunqContexts(dataBytes)
	if err != nil {
		logger.Log.Error(err.Error())
	}

	// CLEANUP
	err = os.Remove(tempFilename)
	if err != nil {
		logger.Log.Errorf("Failed to remove file '%s': %s", tempFilename, err.Error())
	}

	c.JSON(200, contexts)
}

// @Tags Context
// @Produce json
// @Success 200 {array} dtos.PunqContext
// @Router /backend/context [post]
// @Param body body dtos.PunqContext false "PunqContext"
// @Security Bearer
func addContext(c *gin.Context) {
	receivedContexts := []dtos.PunqContext{}
	if err := c.BindJSON(&receivedContexts); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	addedContexts := []dtos.PunqContext{}
	for _, ctx := range receivedContexts {
		addedCtx, err := services.AddContext(ctx)
		if err != nil {
			fmt.Println(err.Error())
			c.JSON(http.StatusInternalServerError, err)
		}
		fmt.Printf("Context '%s' added ✅.\n", addedCtx.Name)
		if addedCtx != nil {
			addedContexts = append(addedContexts, *addedCtx)
		}
	}

	c.JSON(200, addedContexts)
}
