package kubernetes

import (
	"fmt"
	"io"
	"log"

	"github.com/gorilla/websocket"
	"github.com/mogenius/punq/utils"
	// corev1 "k8s.io/api/core/v1"
	//"k8s.io/client-go/tools/remotecommand"
	//"k8s.io/kubectl/pkg/scheme"
)

func SendData(cmdStdin io.WriteCloser, cmdStdout io.ReadCloser) {
	// Create a dialer
	dialer := websocket.DefaultDialer

	// Connect to the server
	conn, _, err := dialer.Dial("ws://127.0.0.1:8080/ws-shell", utils.HttpHeader(""))
	if err != nil {
		fmt.Println("Error connecting to WebSocket:", err)
		return
	}
	defer conn.Close()

	go func() {
		// Forward stdout to the WebSocket
		buf := make([]byte, 1024)
		for {
			n, err := cmdStdout.Read(buf)
			if err != nil {
				log.Println(err)
				return
			}

			err = conn.WriteMessage(websocket.TextMessage, buf[:n])
			if err != nil {
				log.Println(err)
				return
			}
		}
	}()

	for {
		// Forward the WebSocket to stdin
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}

		_, err = cmdStdin.Write(message)
		if err != nil {
			log.Println(err)
			break
		}
	}
}
