package dtos

type KubernetesProvider string

const (
	UNKNOWN                      KubernetesProvider = "UNKNOWN"
	BRING_YOUR_OWN               KubernetesProvider = "BRING_YOUR_OWN"
	DOCKER_ENTERPRISE            KubernetesProvider = "DOCKER_ENTERPRISE"            // Docker
	DOCKER_DESKTOP               KubernetesProvider = "DOCKER_DESKTOP"               // Docker
	AKS                          KubernetesProvider = "AKS"                          // Azure Kubernetes Service
	GKE                          KubernetesProvider = "GKE"                          // Google Kubernetes Engine
	EKS                          KubernetesProvider = "EKS"                          // Amazon Elastic Kubernetes Service
	K3S                          KubernetesProvider = "K3S"                          // K3S
	K3D                          KubernetesProvider = "K3D"                          // K3D
	MINIKUBE                     KubernetesProvider = "MINIKUBE"                     // Minikube
	KIND                         KubernetesProvider = "KIND"                         // Kind
	KUBERNETES                   KubernetesProvider = "KUBERNETES"                   // Kubernetes
	SELF_HOSTED                  KubernetesProvider = "SELF_HOSTED"                  // Self Hosted
	DOKS                         KubernetesProvider = "DOKS"                         // Digital Ocean Kubernetes
	LINODE                       KubernetesProvider = "LINODE"                       // Linode Kubernetes
	IBM                          KubernetesProvider = "IBM"                          // IBM Kubernetes
	ACK                          KubernetesProvider = "ACK"                          // Alibaba Cloud Kubernetes
	OKE                          KubernetesProvider = "OKE"                          // Oracle Cloud Kubernetes
	OTC                          KubernetesProvider = "OTC"                          // Telekom cloud
	OPEN_SHIFT                   KubernetesProvider = "OPEN_SHIFT"                   // RED HAT OpenShift
	GKE_ON_PREM                  KubernetesProvider = "GKE_ON_PREM"                  // Google Kubernetes Engine On-Prem
	RKE                          KubernetesProvider = "RKE"                          // Rancher Kubernetes Engine
	KUBEADM                      KubernetesProvider = "KUBEADM"                      // Kubeadm
	KUBEADM_ON_PREM              KubernetesProvider = "KUBEADM_ON_PREM"              // Kubeadm On-Prem
	KUBEADM_ON_PREM_HETZNER      KubernetesProvider = "KUBEADM_ON_PREM_HETZNER"      // Kubeadm On-Prem Hetzner
	KUBEADM_ON_PREM_DIGITALOCEAN KubernetesProvider = "KUBEADM_ON_PREM_DIGITALOCEAN" // Kubeadm On-Prem Digital Ocean
	KUBEADM_ON_PREM_LINODE       KubernetesProvider = "KUBEADM_ON_PREM_LINODE"       // Kubeadm On-Prem Linode
	KUBEADM_ON_PREM_AWS          KubernetesProvider = "KUBEADM_ON_PREM_AWS"          // Kubeadm On-Prem AWS
	KUBEADM_ON_PREM_AZURE        KubernetesProvider = "KUBEADM_ON_PREM_AZURE"        // Kubeadm On-Prem Azure
	KUBEADM_ON_PREM_GCP          KubernetesProvider = "KUBEADM_ON_PREM_GCP"          // Kubeadm On-Prem GCP
	SYS_ELEVEN                   KubernetesProvider = "SYS_ELEVEN"                   // Managed Kubernetes by SysEleven
	STACKIT                      KubernetesProvider = "SKE"                          // STACKIT Kubernetes Engine (SKE)
	IONOS                        KubernetesProvider = "IONOS"                        // IONOS Cloud Managed
	SCALEWAY                     KubernetesProvider = "SCALEWAY"                     // scaleway
	VMWARE                       KubernetesProvider = "VMWARE"                       // VMware Tanzu Kubernetes Grid Integrated Edition
	MICROK8S                     KubernetesProvider = "MICROK8S"                     // MicroK8s
	CIVO                         KubernetesProvider = "CIVO"                         // Civo Kubernetes
	GIANTSWARM                   KubernetesProvider = "GIANTSWARM"                   // Giant Swarm Kubernetes
	OVHCLOUD                     KubernetesProvider = "OVHCLOUD"                     // OVHCloud Kubernetes
	GARDENER                     KubernetesProvider = "GARDENER"                     // SAP Gardener Kubernetes
	HUAWEI                       KubernetesProvider = "HUAWEI"                       // Huawei Cloud Kubernetes
	NIRMATA                      KubernetesProvider = "NIRMATA"                      // Nirmata Kubernetes
	PF9                          KubernetesProvider = "PF9"                          // Platform9 Kubernetes
	NKS                          KubernetesProvider = "NKS"                          // Netapp Kubernetes Service
	APPSCODE                     KubernetesProvider = "APPSCODE"                     // AppsCode Kubernetes
	LOFT                         KubernetesProvider = "LOFT"                         // Loft Kubernetes
	SPECTROCLOUD                 KubernetesProvider = "SPECTROCLOUD"                 // Spectro Cloud Kubernetes
	DIAMANTI                     KubernetesProvider = "DIAMANTI"                     // Diamanti Kubernetes
)

var ALL_PROVIDER []string = []string{
	string(BRING_YOUR_OWN),
	string(DOCKER_ENTERPRISE),
	string(DOCKER_DESKTOP),
	string(AKS),
	string(GKE),
	string(EKS),
	string(K3S),
	string(K3D),
	string(MINIKUBE),
	string(KIND),
	string(KUBERNETES),
	string(SELF_HOSTED),
	string(DOKS),
	string(LINODE),
	string(IBM),
	string(ACK),
	string(OKE),
	string(OTC),
	string(OPEN_SHIFT),
	string(GKE_ON_PREM),
	string(RKE),
	string(KUBEADM),
	string(KUBEADM_ON_PREM),
	string(KUBEADM_ON_PREM_HETZNER),
	string(KUBEADM_ON_PREM_DIGITALOCEAN),
	string(KUBEADM_ON_PREM_LINODE),
	string(KUBEADM_ON_PREM_AWS),
	string(KUBEADM_ON_PREM_AZURE),
	string(KUBEADM_ON_PREM_GCP),
	string(SYS_ELEVEN),
	string(STACKIT),
	string(IONOS),
	string(SCALEWAY),
	string(VMWARE),
	string(MICROK8S),
	string(CIVO),
	string(GIANTSWARM),
	string(OVHCLOUD),
	string(GARDENER),
	string(HUAWEI),
	string(NIRMATA),
	string(PF9),
	string(NKS),
	string(APPSCODE),
	string(LOFT),
	string(SPECTROCLOUD),
	string(DIAMANTI),
}
