package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type K8sWorkloadResult struct {
	Result interface{} `json:"result,omitempty"`
	Error  interface{} `json:"error,omitempty"`
}

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

func Unauthorized(c *gin.Context, msg string) {
	c.JSON(http.StatusUnauthorized, gin.H{
		"err": msg,
	})
}

func MissingHeader(c *gin.Context, header string) {
	c.JSON(http.StatusBadRequest, gin.H{
		"err": fmt.Sprintf("Missing header '%s'.", header),
	})
	c.AbortWithError(http.StatusBadRequest, fmt.Errorf("%s header is required", header))
}

func MalformedMessage(c *gin.Context, msg string) {
	c.JSON(http.StatusBadRequest, gin.H{
		"err": msg,
	})
}

func HttpRespondForWorkloadResult(c *gin.Context, workloadResult K8sWorkloadResult) {
	if workloadResult.Error == nil {
		c.JSON(http.StatusOK, workloadResult)
	} else {
		c.JSON(http.StatusBadRequest, workloadResult)
	}
}
