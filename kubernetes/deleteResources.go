package kubernetes

import (
	"context"
	"fmt"
	"strings"

	"github.com/mogenius/punq/utils"
	"github.com/mogenius/punq/version"

	"github.com/mogenius/punq/logger"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Remove(clusterName string) {
	provider, err := NewKubeProvider(nil)
	if provider == nil || err != nil {
		logger.Log.Fatal("Failed to load provider.")
	}

	// namespace is not deleted on purpose
	removeRbac(provider)
	removeDeployment(provider)
	removeContextsSecret(provider)
	removeUsersSecret(provider)
	removeService(provider)
	removeIngress(provider)
}

func removeDeployment(provider *KubeProvider) {
	deploymentClient := provider.ClientSet.AppsV1().Deployments(utils.CONFIG.Kubernetes.OwnNamespace)

	// DELETE Deployment
	fmt.Printf("Deleting %s deployment ...\n", version.Name)
	deletePolicy := metav1.DeletePropagationForeground
	err := deploymentClient.Delete(context.TODO(), version.Name, metav1.DeleteOptions{PropagationPolicy: &deletePolicy})
	if err != nil {
		if !apierrors.IsNotFound(err) {
			logger.Log.Error(err)
			return
		}
	}
	fmt.Printf("Deleted %s deployment. ✅\n", version.Name)
}

func removeService(provider *KubeProvider) {
	serviceClient := provider.ClientSet.CoreV1().Services(utils.CONFIG.Kubernetes.OwnNamespace)

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

func removeIngress(provider *KubeProvider) {
	ingressClient := provider.ClientSet.NetworkingV1().Ingresses(utils.CONFIG.Kubernetes.OwnNamespace)

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

	ingressControllerType, err := DetermineIngressControllerType(nil)
	if err != nil {
		if ingressControllerType != NONE && ingressControllerType != UNKNOWN {
			utils.FatalError(err.Error())
		}
	}
	if ingressControllerType == TRAEFIK {
		fmt.Printf("Deleting TRAEFIK middleware ...\n")
		cmd := utils.RunOnLocalShell(fmt.Sprintf("kubectl delete middleware mw-backend -n %s", utils.CONFIG.Kubernetes.OwnNamespace))

		output, err := cmd.CombinedOutput()
		if err != nil {
			if !strings.HasPrefix(string(output), "Error from server (NotFound): ") {
				utils.FatalError(fmt.Sprintf("failed to execute command (%s): %s\n%s", cmd.String(), err.Error(), string(output)))
			}
		}
		fmt.Printf("Deleted TRAEFIK middleware. ✅\n")
	}

}

func removeRbac(provider *KubeProvider) {
	// CREATE RBAC
	fmt.Printf("Deleting %s RBAC ...\n", version.Name)
	err := provider.ClientSet.CoreV1().ServiceAccounts(utils.CONFIG.Kubernetes.OwnNamespace).Delete(context.TODO(), SERVICEACCOUNTNAME, metav1.DeleteOptions{})
	if err != nil {
		if !apierrors.IsNotFound(err) {
			logger.Log.Error(err)
			return
		}
	}
	err = provider.ClientSet.RbacV1().ClusterRoles().Delete(context.TODO(), CLUSTERROLENAME, metav1.DeleteOptions{})
	if err != nil {
		if !apierrors.IsNotFound(err) {
			logger.Log.Error(err)
			return
		}
	}
	err = provider.ClientSet.RbacV1().ClusterRoleBindings().Delete(context.TODO(), CLUSTERROLEBINDINGNAME, metav1.DeleteOptions{})
	if err != nil {
		if !apierrors.IsNotFound(err) {
			logger.Log.Error(err)
			return
		}
	}
	fmt.Printf("Deleted %s RBAC. ✅\n", version.Name)
}

func removeUsersSecret(provider *KubeProvider) {
	secretClient := provider.ClientSet.CoreV1().Secrets(utils.CONFIG.Kubernetes.OwnNamespace)

	fmt.Printf("Deleting %s/%s secret ...\n", utils.CONFIG.Kubernetes.OwnNamespace, utils.USERSSECRET)
	deletePolicy := metav1.DeletePropagationForeground
	err := secretClient.Delete(context.TODO(), utils.USERSSECRET, metav1.DeleteOptions{PropagationPolicy: &deletePolicy})
	if err != nil {
		if !apierrors.IsNotFound(err) {
			logger.Log.Error(err)
			return
		}
	}
	fmt.Printf("Deleted %s/%s secret. ✅\n", utils.CONFIG.Kubernetes.OwnNamespace, utils.USERSSECRET)
}

func removeContextsSecret(provider *KubeProvider) {
	secretClient := provider.ClientSet.CoreV1().Secrets(utils.CONFIG.Kubernetes.OwnNamespace)

	fmt.Printf("Deleting %s/%s secret ...\n", utils.CONFIG.Kubernetes.OwnNamespace, utils.CONTEXTSSECRET)
	deletePolicy := metav1.DeletePropagationForeground
	err := secretClient.Delete(context.TODO(), utils.CONTEXTSSECRET, metav1.DeleteOptions{PropagationPolicy: &deletePolicy})
	if err != nil {
		if !apierrors.IsNotFound(err) {
			logger.Log.Error(err)
			return
		}
	}
	fmt.Printf("Deleted %s/%s secret. ✅\n", utils.CONFIG.Kubernetes.OwnNamespace, utils.CONTEXTSSECRET)
}
