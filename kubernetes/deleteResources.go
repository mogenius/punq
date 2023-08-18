package kubernetes

import (
	"context"

	"github.com/mogenius/punq/utils"
	"github.com/mogenius/punq/version"

	"github.com/mogenius/punq/logger"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Remove() {
	provider, err := NewKubeProviderLocal()
	if err != nil {
		panic(err)
	}

	// namespace is not deleted on purpose
	removeRbac(provider)
	removeDeployment(provider)
	removeSecret(provider)
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

func removeSecret(kubeProvider *KubeProvider) {
	secretClient := kubeProvider.ClientSet.CoreV1().Secrets(utils.CONFIG.Kubernetes.OwnNamespace)

	// DELETE Secret
	logger.Log.Infof("Deleting %s secret ...", version.Name)
	deletePolicy := metav1.DeletePropagationForeground
	err := secretClient.Delete(context.TODO(), version.Name, metav1.DeleteOptions{PropagationPolicy: &deletePolicy})
	if err != nil {
		logger.Log.Error(err)
		return
	}
	logger.Log.Infof("Deleted %s secret.", version.Name)
}
