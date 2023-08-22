package operator

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitContextRoutes(router *gin.Engine) {
	router.GET("/context", context)
}

func context(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "context",
	})
}
