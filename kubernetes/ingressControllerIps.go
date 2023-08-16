// If you have an ingress controller which is processing the traffic from the load balancer
// most of the external traffic will be counted as local traffic because it is ingress-controller
// to pod communication. To identify this traffic we gather the ingress-controller internal ips
// to exclude this traffic from the local traffic counting.

package kubernetes

import (
	"context"
	"fmt"
	"net"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetIngressControllerIps(useLocalKubeConfig bool) []net.IP {
	var result []net.IP
	var kubeProvider *KubeProvider
	var err error
	if useLocalKubeConfig == true {
		kubeProvider, err = NewKubeProviderLocal()
	} else {
		kubeProvider, err = NewKubeProviderInCluster()
	}
	if err != nil {
		panic(err)
	}

	labelSelector := fmt.Sprintf("app.kubernetes.io/component=controller,app.kubernetes.io/instance=nginx-ingress,app.kubernetes.io/name=ingress-nginx")

	pods, err := kubeProvider.ClientSet.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{LabelSelector: labelSelector})

	for _, pod := range pods.Items {
		ip := net.ParseIP(pod.Status.PodIP)
		fmt.Println(pod.Name, ip)
		if ip != nil {
			result = append(result, ip)
		}
	}

	if err != nil {
		fmt.Println("Error:", err)
		return result
	}
	return result
}
