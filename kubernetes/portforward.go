// Taken from https://github.com/gianarb/kube-port-forward
// Thanks for the wonderfull work @gianarb and the great blog entry

package kubernetes

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/mogenius/punq/logger"
	"github.com/mogenius/punq/utils"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
)

type PortForwardAPodRequest struct {
	Pod       v1.Pod
	LocalPort int
	PodPort   int
	Out       *bytes.Buffer
	ErrOut    *bytes.Buffer
	// StopCh is the channel used to manage the port forward lifecycle
	StopCh <-chan struct{}
	// ReadyCh communicates when the tunnel is ready to receive traffic
	ReadyCh chan struct{}
}

func StartPortForward(localPort int, podPort int) {
	for {
		pod := GetFirstPodForLabelName(utils.CONFIG.Kubernetes.OwnNamespace, "app=punq")
		if pod == nil {
			return
		}

		fmt.Printf("Starting PortForward for %s(%d:%d) ...\n", pod.Name, localPort, podPort)

		var wg sync.WaitGroup
		wg.Add(1)

		stopCh := make(chan struct{}, 1)
		readyCh := make(chan struct{})
		out, errOut := new(bytes.Buffer), new(bytes.Buffer)

		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		go func() {
			<-sigs
			fmt.Println("Port-Forward to punq closed!")
			close(stopCh)
			wg.Done()
		}()

		go func() {
			err := portForwardAPod(PortForwardAPodRequest{
				Pod:       *pod,
				LocalPort: localPort,
				PodPort:   podPort,
				Out:       out,
				ErrOut:    errOut,
				StopCh:    stopCh,
				ReadyCh:   readyCh,
			})
			if err != nil {
				logger.Log.Warning("ERROR DURING PORTFORWARD!")
				panic(err)
			}
		}()

		select {
		case <-readyCh:
			fmt.Printf("PortForward for %s(%d:%d) is ready!\n", pod.Name, localPort, podPort)
			break
		case <-stopCh:
			fmt.Printf("PortForward for %s is stopped!\n", pod.Name)
			break
		}

		wg.Wait()

		logger.Log.Warning("TUNNEL CLOSED!")
		time.Sleep(1 * time.Second) // wait a sec before retrying
	}
}

func portForwardAPod(req PortForwardAPodRequest) error {
	kubeProvider := NewKubeProvider()

	path := fmt.Sprintf("/api/v1/namespaces/%s/pods/%s/portforward", req.Pod.Namespace, req.Pod.Name)
	hostIP := strings.TrimLeft(kubeProvider.ClientConfig.Host, "htps:/")

	transport, upgrader, err := spdy.RoundTripperFor(&kubeProvider.ClientConfig)
	if err != nil {
		logger.Log.Error(err)
		return err
	}

	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: transport}, http.MethodPost, &url.URL{Scheme: "https", Path: path, Host: hostIP})
	fw, err := portforward.New(dialer, []string{fmt.Sprintf("%d:%d", req.LocalPort, req.PodPort)}, req.StopCh, req.ReadyCh, req.Out, req.ErrOut)
	if err != nil {
		logger.Log.Error(err)
		return err
	}
	return fw.ForwardPorts()
}