package kubernetes

import (
	"context"
	"fmt"

	"github.com/mogenius/punq/utils"
	"github.com/mogenius/punq/version"

	"github.com/mogenius/punq/logger"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Remove(clusterName string) {
	provider := NewKubeProvider(nil)
	if provider == nil {
		logger.Log.Fatal("Failed to load kubeprovider.")
	}

	// namespace is not deleted on purpose
	removeRbac(provider)
	removeDeployment(provider)
	removeContextsSecret(provider)
	removeUsersSecret(provider)
	removeService(provider)
	removeIngress(provider)

	fmt.Printf("\n🚀🚀🚀 Successfuly uninstalled punq from '%s'.\n\n", clusterName)
}

func removeDeployment(kubeProvider *KubeProvider) {
	deploymentClient := kubeProvider.ClientSet.AppsV1().Deployments(utils.CONFIG.Kubernetes.OwnNamespace)

	// DELETE Deployment
	fmt.Printf("Deleting %s deployment ...\n", version.Name)
	deletePolicy := metav1.DeletePropagationForeground
	err := deploymentClient.Delete(context.TODO(), version.Name, metav1.DeleteOptions{PropagationPolicy: &deletePolicy})
	if err != nil {
		logger.Log.Error(err)
		return
	}
	fmt.Printf("Deleted %s deployment. ✅\n", version.Name)
}

func removeService(kubeProvider *KubeProvider) {
	serviceClient := kubeProvider.ClientSet.CoreV1().Services(utils.CONFIG.Kubernetes.OwnNamespace)

	fmt.Printf("Deleting %s service ...\n", SERVICENAME)
	deletePolicy := metav1.DeletePropagationForeground
	err := serviceClient.Delete(context.TODO(), SERVICENAME, metav1.DeleteOptions{PropagationPolicy: &deletePolicy})
	if err != nil {
		if !apierrors.IsNotFound(err) {
			logger.Log.Error(err)
			return
		}
	}
	fmt.Printf("Deleted %s service. ✅\n", SERVICENAME)
}

func removeIngress(kubeProvider *KubeProvider) {
	ingressClient := kubeProvider.ClientSet.NetworkingV1().Ingresses(utils.CONFIG.Kubernetes.OwnNamespace)

	fmt.Printf("Deleting %s ingress ...\n", INGRESSNAME)
	deletePolicy := metav1.DeletePropagationForeground
	err := ingressClient.Delete(context.TODO(), INGRESSNAME, metav1.DeleteOptions{PropagationPolicy: &deletePolicy})
	if err != nil {
		if !apierrors.IsNotFound(err) {
			logger.Log.Error(err)
			return
		}
	}
	fmt.Printf("Deleted %s ingress. ✅\n", INGRESSNAME)
}

func removeRbac(kubeProvider *KubeProvider) {
	// CREATE RBAC
	fmt.Printf("Deleting %s RBAC ...\n", version.Name)
	err := kubeProvider.ClientSet.CoreV1().ServiceAccounts(utils.CONFIG.Kubernetes.OwnNamespace).Delete(context.TODO(), SERVICEACCOUNTNAME, metav1.DeleteOptions{})
	if err != nil {
		logger.Log.Error(err)
		return
	}
	err = kubeProvider.ClientSet.RbacV1().ClusterRoles().Delete(context.TODO(), CLUSTERROLENAME, metav1.DeleteOptions{})
	if err != nil {
		logger.Log.Error(err)
		return
	}
	err = kubeProvider.ClientSet.RbacV1().ClusterRoleBindings().Delete(context.TODO(), CLUSTERROLEBINDINGNAME, metav1.DeleteOptions{})
	if err != nil {
		logger.Log.Error(err)
		return
	}
	fmt.Printf("Deleted %s RBAC. ✅\n", version.Name)
}

func removeUsersSecret(kubeProvider *KubeProvider) {
	secretClient := kubeProvider.ClientSet.CoreV1().Secrets(utils.CONFIG.Kubernetes.OwnNamespace)

	fmt.Printf("Deleting %s/%s secret ...\n", utils.CONFIG.Kubernetes.OwnNamespace, utils.USERSSECRET)
	deletePolicy := metav1.DeletePropagationForeground
	err := secretClient.Delete(context.TODO(), utils.USERSSECRET, metav1.DeleteOptions{PropagationPolicy: &deletePolicy})
	if err != nil {
		logger.Log.Error(err)
		return
	}
	fmt.Printf("Deleted %s/%s secret. ✅\n", utils.CONFIG.Kubernetes.OwnNamespace, utils.USERSSECRET)
}

func removeContextsSecret(kubeProvider *KubeProvider) {
	secretClient := kubeProvider.ClientSet.CoreV1().Secrets(utils.CONFIG.Kubernetes.OwnNamespace)

	fmt.Printf("Deleting %s/%s secret ...\n", utils.CONFIG.Kubernetes.OwnNamespace, utils.CONTEXTSSECRET)
	deletePolicy := metav1.DeletePropagationForeground
	err := secretClient.Delete(context.TODO(), utils.CONTEXTSSECRET, metav1.DeleteOptions{PropagationPolicy: &deletePolicy})
	if err != nil {
		logger.Log.Error(err)
		return
	}
	fmt.Printf("Deleted %s/%s secret. ✅\n", utils.CONFIG.Kubernetes.OwnNamespace, utils.CONTEXTSSECRET)
}
