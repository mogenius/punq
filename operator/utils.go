package operator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/gin-gonic/gin"
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
