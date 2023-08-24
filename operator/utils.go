package operator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	gin "github.com/gin-gonic/gin"
	"github.com/mogenius/punq/kubernetes"
)

const (
	OPERATOR_PORT = 8080
)

func PrintPrettyPost(c *gin.Context) {
	var out bytes.Buffer
	body, _ := io.ReadAll(c.Request.Body)
	err := json.Indent(&out, []byte(body), "", "  ")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(out.String())
}

func GetRequiredHeader(c *gin.Context, headerField string) string {
	selectedField := c.Request.Header.Get(headerField)
	if selectedField == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": fmt.Sprintf("Missing header '%s'.", headerField),
		})
	}
	return selectedField
}

func NotFound(c *gin.Context, msg string) {
	c.JSON(http.StatusNotFound, gin.H{
		"err": msg,
	})
}

func MalformedMessage(c *gin.Context, msg string) {
	c.JSON(http.StatusBadRequest, gin.H{
		"err": msg,
	})
}

func RespondForWorkloadResult(c *gin.Context, workloadResult kubernetes.K8sWorkloadResult) {
	if workloadResult.Error == nil {
		c.JSON(http.StatusOK, workloadResult)
	} else {
		c.JSON(http.StatusBadRequest, workloadResult)
	}
}
