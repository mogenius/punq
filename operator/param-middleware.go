package operator

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mogenius/punq/utils"
)

func validateParam(params ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, param := range params {
			id := c.Param(param)

			if id == "" {
				utils.MalformedMessage(c, fmt.Sprintf("%s cannot be empty", param))
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
