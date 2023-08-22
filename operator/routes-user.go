package operator

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitUserRoutes(router *gin.Engine) {
	router.GET("/user", user)
}

func user(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "user",
	})
}
