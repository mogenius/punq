package operator

import (
	"net/http"

	"github.com/gin-gonic/gin"
	punqVersion "github.com/mogenius/punq/version"
)

func InitGeneralRoutes(router *gin.Engine) {
	router.GET("/version", version)
}

func version(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"Name":           punqVersion.Name,
		"Version":        punqVersion.Ver,
		"Branch":         punqVersion.Branch,
		"BuildTimestamp": punqVersion.BuildTimestamp,
		"GitCommitHash":  punqVersion.GitCommitHash,
	})
}
