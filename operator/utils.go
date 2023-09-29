package operator

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/mogenius/punq/utils"
)

func CreateLogger(loggerPrefix string) gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		var statusColor, methodColor, resetColor string
		if param.IsOutputColor() {
			statusColor = param.StatusCodeColor()
			methodColor = param.MethodColor()
			resetColor = param.ResetColor()
		}

		size := utils.BytesToHumanReadable(int64(param.BodySize))
		if param.BodySize < 0 {
			size = "0 B"
		}

		return fmt.Sprintf("[%s] %v |%s %3d %s| %13v | %15s |%s %-7s %s | %10s | %#v\n%s",
			loggerPrefix,
			param.TimeStamp.Format("2006/01/02 - 15:04:05"),
			statusColor, param.StatusCode, resetColor,
			param.Latency,
			param.ClientIP,
			methodColor, param.Method, resetColor,
			size,
			param.Path,
			param.ErrorMessage,
		)
	})
}
