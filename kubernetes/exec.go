package kubernetes

import (
	"fmt"
	"io"
	"log"
	"os/exec"

	"github.com/gorilla/websocket"
	"github.com/mogenius/punq/logger"
	"github.com/mogenius/punq/utils"
	// corev1 "k8s.io/api/core/v1"
	//"k8s.io/client-go/tools/remotecommand"
	//"k8s.io/kubectl/pkg/scheme"
)

func ExecTest() error {
	namespace := "mogenius"
	podName := "mogenius-k8s-manager-6dcf5df696-8bsf4"
	container := "mogenius-k8s-manager"
	execCmd := []string{"exec", "--stdin", "--tty", "-n", namespace, "--container", container, podName, "--", "ifconfig"} // /bin/sh

	// Create an *exec.Cmd
	cmd := exec.Command("kubectl", execCmd...)
	cmdStdout, _ := cmd.StdoutPipe()
	cmdStdin, _ := cmd.StdinPipe()

	// Assign os.Stdin, os.Stdout, and os.Stderr
	// cmd.Stdin = os.Stdin
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr

	go sendData(cmdStdin, cmdStdout)

	// Run the command
	err := cmd.Run()

	if err != nil {
		logger.Log.Errorf("ExecTest ERR: %s", err.Error())
	}

	return err
}

// This method is the equivalent of:
// kubectl exec -n benegeilomat-dev-8umm0v --container nginx nginx-6b64bff7c9-p6vpp -- ls
// func ExecTest() error {
// 	provider,err := NewKubeProvider()

// 	namespace := "benegeilomat-dev-8umm0v"
// 	podName := "nginx-6b64bff7c9-p6vpp"
// 	container := "nginx"
// 	cmd := []string{"/bin/bash"} // ls

// 	req := provider.ClientSet.CoreV1().RESTClient().Post().
// 		Resource("pods").
// 		Namespace(namespace).
// 		Name(podName).
// 		SubResource("exec").
// 		VersionedParams(&corev1.PodExecOptions{
// 			Command:   cmd,
// 			Container: container,
// 			Stdin:     true,  // DO NOT CHANGE: ONLY WORKING CONFIGURATION
// 			Stdout:    true,  // DO NOT CHANGE: ONLY WORKING CONFIGURATION
// 			Stderr:    false, // DO NOT CHANGE: ONLY WORKING CONFIGURATION
// 			TTY:       true,  // DO NOT CHANGE: ONLY WORKING CONFIGURATION
// 		}, scheme.ParameterCodec)

// 	exec, err := remotecommand.NewSPDYExecutor(&provider.ClientConfig, "POST", req.URL())
// 	if err != nil {
// 		logger.Log.Errorf("ExecTest ERR: %s", err.Error())
// 	}

// 	err = exec.StreamWithContext(context.TODO(), remotecommand.StreamOptions{
// 		Stdin:  os.Stdin,  // DO NOT CHANGE: ONLY WORKING CONFIGURATION
// 		Stdout: os.Stdout, // DO NOT CHANGE: ONLY WORKING CONFIGURATION
// 		Stderr: nil,       // DO NOT CHANGE: ONLY WORKING CONFIGURATION
// 		Tty:    true,      // DO NOT CHANGE: ONLY WORKING CONFIGURATION
// 	})

// 	if err != nil {
// 		logger.Log.Errorf("ExecTest ERR: %s", err.Error())
// 	}
// 	return err
// }

func sendData(cmdStdin io.WriteCloser, cmdStdout io.ReadCloser) {
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
