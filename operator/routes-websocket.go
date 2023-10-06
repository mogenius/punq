package operator

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/mogenius/punq/dtos"
	"github.com/mogenius/punq/utils"
)

func InitWebsocketRoutes(router *gin.Engine) {
	router.GET("/exec-sh", AuthByParameter(dtos.ADMIN), connectWs)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // adjust to implement your origin validation logic
	},
}

func connectWs(c *gin.Context) {
	namespace, namespaceOk := c.GetQuery("namespace")
	if !namespaceOk || namespace == "" {
		utils.MissingQueryParameter(c, "namespace")
		return
	}

	container, containerOk := c.GetQuery("container")
	if !containerOk || container == "" {
		utils.MissingQueryParameter(c, "container")
		return
	}

	podName, podNameOk := c.GetQuery("podname")
	if !podNameOk || podName == "" {
		utils.MissingQueryParameter(c, "podname")
		return
	}

	log.Printf("exec-sh: %s %s %s\n", namespace, container, podName)

	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade ws: %+v", err)
		return
	}
	defer func() {
		ws.Close()
	}()

	cmd := exec.Command("sh", "-c", fmt.Sprintf("kubectl exec -i --tty -c %s -n %s %s -- /bin/sh", container, namespace, podName))
	cmd.Env = os.Environ()
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal("Error creating stdin pipe:", err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal("Error creating stdout pipe:", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal("Error creating stderr pipe:", err)
	}

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			data := scanner.Bytes()
			if utils.CONFIG.Misc.Debug {
				fmt.Printf("Response-Line: '%s'\n", string(data))
			}
			err = ws.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				log.Printf("Error writing to ws: %+v", err)
			}
		}
	}()
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			data := scanner.Bytes()
			if utils.CONFIG.Misc.Debug {
				fmt.Printf("Response-Line: '%s'\n", string(data))
			}
			err = ws.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				log.Printf("Error writing to ws: %+v", err)
			}
		}
	}()

	go func() {
		for {
			_, msg, err := ws.ReadMessage()
			if utils.CONFIG.Misc.Debug {
				fmt.Printf("Received Cmd: '%s'", string(msg))
			}
			msg = append(msg, '\n')
			if err != nil {
				log.Printf("Error reading from ws: %+v", err)
				log.Printf("CLOSE: exec-sh: %s %s %s\n", namespace, container, podName)
				break
			}
			_, err = stdin.Write(msg)
			if err != nil {
				log.Printf("Error writing to stdin: %+v", err)
			}
		}
	}()

	err = cmd.Start()
	if err != nil {
		log.Printf("Error starting cmd: %+v", err)
		return
	}

	err = cmd.Wait()
	if err != nil {
		log.Printf("Cmd returned error: %+v", err.Error())
		return
	}
}
