package operator

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mogenius/punq/dtos"
	"github.com/mogenius/punq/structs"
	punqVersion "github.com/mogenius/punq/version"
)

func InitGeneralRoutes(router *gin.Engine) {
	router.GET("/version", versionData)
	router.GET("/providers", allProviders)
}

// @Tags Misc
// @Produce json
// @Success 200 {object} structs.Version
// @Router /backend/version [get]
func versionData(c *gin.Context) {
	c.JSON(http.StatusOK, structs.VersionFrom(punqVersion.Name, punqVersion.Ver, punqVersion.Branch, punqVersion.BuildTimestamp, punqVersion.GitCommitHash, punqVersion.OperatorImage))
}

// @Tags Misc
// @Produce json
// @Success 200 {array} string
// @Router /backend/providers [get]
// @Security Bearer
func allProviders(c *gin.Context) {
	c.JSON(http.StatusOK, dtos.ALL_PROVIDER)
}
