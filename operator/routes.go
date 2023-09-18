package operator

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mogenius/punq/structs"
	punqVersion "github.com/mogenius/punq/version"
)

func InitGeneralRoutes(router *gin.Engine) {
	router.GET("/version", versionData)
}

// @Tags Misc
// @Produce json
// @Success 200 {object} structs.Version
// @Router /version [get]
func versionData(c *gin.Context) {
	c.JSON(http.StatusOK, structs.VersionFrom(punqVersion.Name, punqVersion.Ver, punqVersion.Branch, punqVersion.BuildTimestamp, punqVersion.GitCommitHash, punqVersion.OperatorImage))
}
