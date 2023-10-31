// Taken from https://github.com/gianarb/kube-port-forward
// Thanks for the wonderfull work @gianarb and the great blog entry

package kubernetes

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httputil"
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

func StartPortForward(localPort int, podPort int, readyChannel chan struct{}, stopChannel chan struct{}, contextId *string) {
	for {
		pod := GetFirstPodForLabelName(utils.CONFIG.Kubernetes.OwnNamespace, "app=punq", contextId)
		if pod == nil {
			return
		}

		fmt.Printf("Starting PortForward for %s(%d:%d) ...\n", pod.Name, localPort, podPort)

		var wg sync.WaitGroup
		wg.Add(1)

		out, errOut := new(bytes.Buffer), new(bytes.Buffer)

		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		go func() {
			<-sigs
			fmt.Printf("Port-Forward to punq (%d:%d) closed!\n", localPort, podPort)
			close(stopChannel)
			wg.Done()
		}()

		go func() {
			err := portForwardAPod(PortForwardAPodRequest{
				Pod:       *pod,
				LocalPort: localPort,
				PodPort:   podPort,
				Out:       out,
				ErrOut:    errOut,
				StopCh:    stopChannel,
				ReadyCh:   readyChannel,
			}, contextId)
			if err != nil {
				logger.Log.Warning("ERROR DURING PORTFORWARD!")
				panic(err)
			}
		}()

		<-stopChannel
		fmt.Printf("PortForward for %s is stopped!\n", pod.Name)

		wg.Wait()

		logger.Log.Warning("TUNNEL CLOSED!")
		time.Sleep(1 * time.Second) // wait a sec before retrying
	}
}

func portForwardAPod(req PortForwardAPodRequest, contextId *string) error {
	provider, err := NewKubeProvider(contextId)
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/api/v1/namespaces/%s/pods/%s/portforward", req.Pod.Namespace, req.Pod.Name)
	hostIP := strings.TrimLeft(provider.ClientConfig.Host, "htps:/")

	transport, upgrader, err := spdy.RoundTripperFor(&provider.ClientConfig)
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

func Proxy(backendUrl string, frontendUrl string, websocketUrl string) {
	localPort := fmt.Sprintf(":%d", utils.CONFIG.Misc.ProxyPort)

	backendURL, err := url.Parse(backendUrl)
	if err != nil {
		utils.FatalError(fmt.Sprintf("Error parsing backend url: %s", err.Error()))
	}
	frontendURL, err := url.Parse(frontendUrl)
	if err != nil {
		utils.FatalError(fmt.Sprintf("Error parsing frontend url: %s", err.Error()))
	}
	websocketURL, err := url.Parse(websocketUrl)
	if err != nil {
		utils.FatalError(fmt.Sprintf("Error parsing websocket url: %s", err.Error()))
	}

	backendProxy := httputil.NewSingleHostReverseProxy(backendURL)
	frontendProxy := httputil.NewSingleHostReverseProxy(frontendURL)
	websocketProxy := httputil.NewSingleHostReverseProxy(websocketURL)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		orig := r.URL
		fmt.Printf("FRONTEND: localhost%s%s -> %s%s\n", localPort, orig, frontendURL, r.URL)
		frontendProxy.ServeHTTP(w, r)
	})

	http.HandleFunc("/backend/", func(w http.ResponseWriter, r *http.Request) {
		orig := r.URL
		r.URL.Path = strings.Replace(r.URL.Path, "/backend", "", 1)
		fmt.Printf("BACKEND:  localhost%s%s -> %s%s\n", localPort, orig, backendURL, r.URL)
		backendProxy.ServeHTTP(w, r)
	})

	http.HandleFunc("/websocket/", func(w http.ResponseWriter, r *http.Request) {
		orig := r.URL
		r.URL.Path = strings.Replace(r.URL.Path, "/websocket", "", 1)
		r.Proto = "http"
		fmt.Printf("WEBSOCKET:localhost%s%s -> %s%s\n", localPort, orig, websocketURL, r.URL)
		websocketProxy.ServeHTTP(w, r)
	})

	if err = http.ListenAndServe(localPort, nil); err != nil {
		utils.FatalError(fmt.Sprintf("Error starting proxy: %s", err.Error()))
	}
}
