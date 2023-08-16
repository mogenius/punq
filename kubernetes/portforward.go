// Taken from https://github.com/gianarb/kube-port-forward
// Thanks for the wonderfull work @gianarb and the great blog entry

package kubernetes

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"punq/logger"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

func StartPortForward(kubeProvider *KubeProvider, useLocalKubeConfig bool) {
	for {
		pod, err := getPodRedisPodname(kubeProvider, useLocalKubeConfig)
		if err != nil {
			logger.Log.Error(err)
		}

		logger.Log.Infof("Starting PortForward for %s ...", pod.Name)

		var wg sync.WaitGroup
		wg.Add(1)

		// stopCh control the port forwarding lifecycle. When it gets closed the
		// port forward will terminate
		stopCh := make(chan struct{}, 1)
		// readyCh communicate when the port forward is ready to get traffic
		readyCh := make(chan struct{})
		// stream is used to tell the port forwarder where to place its output or
		// where to expect input if needed. For the port forwarding we just need
		// the output eventually
		out, errOut := new(bytes.Buffer), new(bytes.Buffer)

		// managing termination signal from the terminal. As you can see the stopCh
		// gets closed to gracefully handle its termination.
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		go func() {
			<-sigs
			fmt.Println("Bye...")
			close(stopCh)
			wg.Done()
		}()

		go func() {
			// PortForward the pod specified from its port 9090 to the local port
			// 8080
			err := portForwardAPod(kubeProvider, useLocalKubeConfig, PortForwardAPodRequest{
				Pod:       *pod,
				LocalPort: int(REDISPORT),
				PodPort:   int(REDISPORT),
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
			logger.Log.Infof("PortForward for %s is ready!", pod.Name)
			break
		case <-stopCh:
			logger.Log.Infof("PortForward for %s is stopped!", pod.Name)
			break
		}

		wg.Wait()

		logger.Log.Warning("TUNNEL CLOSED!")
		time.Sleep(1 * time.Second) // wait a sec before retrying
	}
}

func portForwardAPod(kubeProvider *KubeProvider, useLocalKubeConfig bool, req PortForwardAPodRequest) error {
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

func getPodRedisPodname(kubeProvider *KubeProvider, useLocalKubeConfig bool) (*v1.Pod, error) {
	labelSelector := fmt.Sprintf("app=%s", REDISNAME)
	pods, err := kubeProvider.ClientSet.CoreV1().Pods(NAMESPACE).List(context.TODO(), metav1.ListOptions{LabelSelector: labelSelector})

	for _, pod := range pods.Items {
		return &pod, nil
	}

	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	return nil, fmt.Errorf("Neither pod found nor error received.")
}
