package operator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	gin "github.com/gin-gonic/gin"
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
