package kubernetes

import (
	"context"

	"github.com/mogenius/punq/utils"
	"github.com/mogenius/punq/version"

	"github.com/mogenius/punq/logger"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Remove(clusterName string) {
	provider, err := NewKubeProviderLocal()
	if err != nil {
		panic(err)
	}

	// namespace is not deleted on purpose
	removeRbac(provider)
	removeDeployment(provider)
	removeContextsSecret(provider)
	removeUsersSecret(provider)
	removeService(provider)
	removeIngress(provider)

	logger.Log.Noticef("🚀🚀🚀 Successfuly uninstalled punq from '%s'.", clusterName)
}

func removeDeployment(kubeProvider *KubeProvider) {
	deploymentClient := kubeProvider.ClientSet.AppsV1().Deployments(utils.CONFIG.Kubernetes.OwnNamespace)

	// DELETE Deployment
	logger.Log.Infof("Deleting %s deployment ...", version.Name)
	deletePolicy := metav1.DeletePropagationForeground
	err := deploymentClient.Delete(context.TODO(), version.Name, metav1.DeleteOptions{PropagationPolicy: &deletePolicy})
	if err != nil {
		logger.Log.Error(err)
		return
	}
	logger.Log.Infof("Deleted %s deployment.", version.Name)
}

func removeService(kubeProvider *KubeProvider) {
	serviceClient := kubeProvider.ClientSet.CoreV1().Services(utils.CONFIG.Kubernetes.OwnNamespace)

	logger.Log.Infof("Deleting %s service ...", SERVICENAME)
	deletePolicy := metav1.DeletePropagationForeground
	err := serviceClient.Delete(context.TODO(), SERVICENAME, metav1.DeleteOptions{PropagationPolicy: &deletePolicy})
	if err != nil {
		if !apierrors.IsNotFound(err) {
			logger.Log.Error(err)
			return
		}
	}
	logger.Log.Infof("Deleted %s service.", SERVICENAME)
}

func removeIngress(kubeProvider *KubeProvider) {
	ingressClient := kubeProvider.ClientSet.NetworkingV1().Ingresses(utils.CONFIG.Kubernetes.OwnNamespace)

	logger.Log.Infof("Deleting %s ingress ...", INGRESSNAME)
	deletePolicy := metav1.DeletePropagationForeground
	err := ingressClient.Delete(context.TODO(), INGRESSNAME, metav1.DeleteOptions{PropagationPolicy: &deletePolicy})
	if err != nil {
		if !apierrors.IsNotFound(err) {
			logger.Log.Error(err)
			return
		}
	}
	logger.Log.Infof("Deleted %s ingress.", INGRESSNAME)
}

func removeRbac(kubeProvider *KubeProvider) {
	// CREATE RBAC
	logger.Log.Infof("Deleting %s RBAC ...", version.Name)
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
	logger.Log.Infof("Deleted %s RBAC.", version.Name)
}

func removeUsersSecret(kubeProvider *KubeProvider) {
	secretClient := kubeProvider.ClientSet.CoreV1().Secrets(utils.CONFIG.Kubernetes.OwnNamespace)

	logger.Log.Infof("Deleting %s/%s secret ...", utils.CONFIG.Kubernetes.OwnNamespace, utils.USERSSECRET)
	deletePolicy := metav1.DeletePropagationForeground
	err := secretClient.Delete(context.TODO(), utils.USERSSECRET, metav1.DeleteOptions{PropagationPolicy: &deletePolicy})
	if err != nil {
		logger.Log.Error(err)
		return
	}
	logger.Log.Infof("Deleted %s/%s secret.", utils.CONFIG.Kubernetes.OwnNamespace, utils.USERSSECRET)
}

func removeContextsSecret(kubeProvider *KubeProvider) {
	secretClient := kubeProvider.ClientSet.CoreV1().Secrets(utils.CONFIG.Kubernetes.OwnNamespace)

	logger.Log.Infof("Deleting %s/%s secret ...", utils.CONFIG.Kubernetes.OwnNamespace, utils.CONTEXTSSECRET)
	deletePolicy := metav1.DeletePropagationForeground
	err := secretClient.Delete(context.TODO(), utils.CONTEXTSSECRET, metav1.DeleteOptions{PropagationPolicy: &deletePolicy})
	if err != nil {
		logger.Log.Error(err)
		return
	}
	logger.Log.Infof("Deleted %s/%s secret.", utils.CONFIG.Kubernetes.OwnNamespace, utils.CONTEXTSSECRET)
}
