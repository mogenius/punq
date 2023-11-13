package operator

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/creack/pty"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/mogenius/punq/dtos"
	"github.com/mogenius/punq/kubernetes"
	"github.com/mogenius/punq/utils"
)

type windowSize struct {
	Rows uint16 `json:"rows"`
	Cols uint16 `json:"cols"`
}

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

	contextId, contextOk := c.GetQuery("context")
	if !contextOk || contextId == "" {
		utils.MissingQueryParameter(c, "context")
		return
	}

	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade ws: %s", err.Error())
		return
	}
	defer func() {
		ws.Close()
	}()
	if err != nil {
		log.Printf("Unable to upgrade connection: %s", err.Error())
		return
	}

	selectedShell := FindValidShell(container, namespace, podName)

	cmd := utils.RunOnLocalShell(fmt.Sprintf("kubectl exec -it -c %s -n %s %s %s -- %s -c 'echo -e \"\033[1;34mConnected to %s/%s/%s using \"$(echo $0)\". Happy hacking!\033[0m 🚀 🚀 🚀\"; %s'", container, namespace, podName, kubernetes.ContextFlag(&contextId), selectedShell, namespace, podName, container, selectedShell))
	fmt.Println(cmd.String())
	cmd.Env = append(os.Environ(), "TERM=xterm-color")

	tty, err := pty.Start(cmd)
	if err != nil {
		log.Printf("Unable to start pty/cmd: %s", err.Error())
		ws.WriteMessage(websocket.TextMessage, []byte(err.Error()))
		return
	}

	defer func() {
		cmd.Process.Kill()
		cmd.Process.Wait()
		tty.Close()
		ws.Close()
	}()

	go func() {
		for {
			buf := make([]byte, 1024)
			read, err := tty.Read(buf)
			if err != nil {
				ws.WriteMessage(websocket.TextMessage, []byte(err.Error()))
				log.Printf("Unable to read from pty/cmd: %s", err.Error())
				return
			}
			ws.WriteMessage(websocket.BinaryMessage, buf[:read])
		}
	}()

	for {
		_, reader, err := ws.ReadMessage()
		if err != nil {
			log.Printf("Unable to grab next reader: %s", err.Error())
			return
		}

		if strings.HasPrefix(string(reader), "\x04") {
			str := strings.TrimPrefix(string(reader), "\x04")

			var resizeMessage windowSize
			err := json.Unmarshal([]byte(str), &resizeMessage)
			if err != nil {
				log.Printf("%s", err.Error())
				continue
			}

			if err := pty.Setsize(tty, &pty.Winsize{Rows: uint16(resizeMessage.Rows), Cols: uint16(resizeMessage.Cols)}); err != nil {
				log.Printf("Unable to resize: %s", err.Error())
				continue
			}
			continue
		}

		tty.Write(reader)
	}
}

func FindValidShell(container string, namespace string, podName string) string {
	availableShells := []string{"bash", "ash", "zsh", "sh", "ksh", "csh"}
	for _, shell := range availableShells {
		cmd := utils.RunOnLocalShell(fmt.Sprintf("kubectl exec -it -c %s -n %s %s -- sh -c '%s'", container, namespace, podName, shell))
		err := cmd.Run()
		if err == nil {
			return shell
		}
	}
	return "sh"
}
