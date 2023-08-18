package kubernetes

import (
	"context"

	"github.com/google/uuid"
	"github.com/mogenius/punq/version"

	"github.com/mogenius/punq/utils"

	"github.com/mogenius/punq/logger"

	core "k8s.io/api/core/v1"
	rbac "k8s.io/api/rbac/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	applyconfapp "k8s.io/client-go/applyconfigurations/apps/v1"
	applyconfcore "k8s.io/client-go/applyconfigurations/core/v1"
	applyconfmeta "k8s.io/client-go/applyconfigurations/meta/v1"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

func Deploy() {
	provider, err := NewKubeProviderLocal()
	if err != nil {
		panic(err)
	}

	applyNamespace(provider)
	addRbac(provider)
	addDeployment(provider)

	_, err = CreateClusterSecretIfNotExist(provider)
	if err != nil {
		logger.Log.Fatalf("Error Creating cluster secret. Aborting: %s.", err.Error())
	}
}

func addRbac(kubeProvider *KubeProvider) error {
	serviceAccount := &core.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name: SERVICEACCOUNTNAME,
		},
	}
	clusterRole := &rbac.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: CLUSTERROLENAME,
		},
		Rules: []rbac.PolicyRule{
			{
				APIGroups: []string{"", "extensions", "apps"},
				Resources: RBACRESOURCES,
				Verbs:     []string{"list", "get", "watch", "create", "update"},
			},
		},
	}
	clusterRoleBinding := &rbac.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: CLUSTERROLEBINDINGNAME,
		},
		RoleRef: rbac.RoleRef{
			Name:     CLUSTERROLENAME,
			Kind:     "ClusterRole",
			APIGroup: "rbac.authorization.k8s.io",
		},
		Subjects: []rbac.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      SERVICEACCOUNTNAME,
				Namespace: utils.CONFIG.Kubernetes.OwnNamespace,
			},
		},
	}

	// CREATE RBAC
	logger.Log.Info("Creating punq RBAC ...")
	_, err := kubeProvider.ClientSet.CoreV1().ServiceAccounts(utils.CONFIG.Kubernetes.OwnNamespace).Create(context.TODO(), serviceAccount, MoCreateOptions())
	if err != nil && !k8serrors.IsAlreadyExists(err) {
		return err
	}
	_, err = kubeProvider.ClientSet.RbacV1().ClusterRoles().Create(context.TODO(), clusterRole, MoCreateOptions())
	if err != nil && !k8serrors.IsAlreadyExists(err) {
		return err
	}
	_, err = kubeProvider.ClientSet.RbacV1().ClusterRoleBindings().Create(context.TODO(), clusterRoleBinding, MoCreateOptions())
	if err != nil && !k8serrors.IsAlreadyExists(err) {
		return err
	}
	logger.Log.Info("Created punq RBAC.")
	return nil
}

func applyNamespace(kubeProvider *KubeProvider) {
	serviceClient := kubeProvider.ClientSet.CoreV1().Namespaces()

	namespace := applyconfcore.Namespace(utils.CONFIG.Kubernetes.OwnNamespace)

	applyOptions := metav1.ApplyOptions{
		Force:        true,
		FieldManager: version.Name,
	}

	logger.Log.Info("Creating punq namespace ...")
	result, err := serviceClient.Apply(context.TODO(), namespace, applyOptions)
	if err != nil {
		logger.Log.Error(err)
	}
	logger.Log.Info("Created punq namespace", result.GetObjectMeta().GetName())
}

func CreateClusterSecretIfNotExist(kubeProvider *KubeProvider) (utils.ClusterSecret, error) {
	secretClient := kubeProvider.ClientSet.CoreV1().Secrets(utils.CONFIG.Kubernetes.OwnNamespace)

	existingSecret, getErr := secretClient.Get(context.TODO(), utils.CONFIG.Kubernetes.OwnNamespace, metav1.GetOptions{})
	return writePunqSecret(secretClient, existingSecret, getErr)
}

func writePunqSecret(secretClient v1.SecretInterface, existingSecret *core.Secret, getErr error) (utils.ClusterSecret, error) {
	clusterSecret := utils.ClusterSecret{
		ApiKey:       uuid.New().String(),
		ClusterMfaId: uuid.New().String(),
		ClusterName:  utils.CONFIG.Kubernetes.ClusterName,
	}

	secret := utils.InitSecret()
	secret.ObjectMeta.Name = utils.CONFIG.Kubernetes.OwnNamespace
	secret.ObjectMeta.Namespace = utils.CONFIG.Kubernetes.OwnNamespace
	delete(secret.StringData, "exampleData") // delete example data
	secret.StringData["cluster-mfa-id"] = clusterSecret.ClusterMfaId
	secret.StringData["api-key"] = clusterSecret.ApiKey
	secret.StringData["cluster-name"] = clusterSecret.ClusterName

	if existingSecret == nil || getErr != nil {
		logger.Log.Info("Creating new punq secret ...")
		result, err := secretClient.Create(context.TODO(), &secret, MoCreateOptions())
		if err != nil {
			logger.Log.Error(err)
			return clusterSecret, err
		}
		logger.Log.Info("Created new punq secret", result.GetObjectMeta().GetName())
	} else {
		if string(existingSecret.Data["api-key"]) != clusterSecret.ApiKey ||
			string(existingSecret.Data["cluster-name"]) != clusterSecret.ClusterName {
			logger.Log.Info("Updating existing punq secret ...")
			// keep existing mfa-id if possible
			if string(existingSecret.Data["cluster-mfa-id"]) != "" {
				clusterSecret.ClusterMfaId = string(existingSecret.Data["cluster-mfa-id"])
				secret.StringData["cluster-mfa-id"] = clusterSecret.ClusterMfaId
			}
			result, err := secretClient.Update(context.TODO(), &secret, MoUpdateOptions())
			if err != nil {
				logger.Log.Error(err)
				return clusterSecret, err
			}
			logger.Log.Info("Updated punq secret", result.GetObjectMeta().GetName())
		} else {
			clusterSecret.ClusterMfaId = string(existingSecret.Data["cluster-mfa-id"])
			logger.Log.Info("Using existing punq secret.")
		}
	}

	return clusterSecret, nil
}

func addDeployment(kubeProvider *KubeProvider) {
	deploymentClient := kubeProvider.ClientSet.AppsV1().Deployments(utils.CONFIG.Kubernetes.OwnNamespace)

	deploymentContainer := applyconfcore.Container()
	deploymentContainer.WithImagePullPolicy(core.PullAlways)
	deploymentContainer.WithName(version.Name)
	deploymentContainer.WithImage(DEPLOYMENTIMAGE)

	envVars := []applyconfcore.EnvVarApplyConfiguration{}
	envVars = append(envVars, applyconfcore.EnvVarApplyConfiguration{
		Name:  utils.Pointer("cluster_name"),
		Value: utils.Pointer("TestClusterFromCode"),
	})
	envVars = append(envVars, applyconfcore.EnvVarApplyConfiguration{
		Name:  utils.Pointer("api_key"),
		Value: utils.Pointer("94E23575-A689-4F88-8D67-215A274F4E6E"), // dont worry. this is a test key
	})
	deploymentContainer.Env = envVars
	// agentResourceLimits := core.ResourceList{
	// 	"cpu":               resource.MustParse("300m"),
	// 	"memory":            resource.MustParse("256Mi"),
	// 	"ephemeral-storage": resource.MustParse("100Mi"),
	// }
	// agentResourceRequests := core.ResourceList{
	// 	"cpu":               resource.MustParse("100m"),
	// 	"memory":            resource.MustParse("128Mi"),
	// 	"ephemeral-storage": resource.MustParse("10Mi"),
	// }
	// agentResources := applyconfcore.ResourceRequirements().WithRequests(agentResourceRequests).WithLimits(agentResourceLimits)
	// deploymentContainer.WithResources(agentResources)
	deploymentContainer.WithName(version.Name)

	podSpec := applyconfcore.PodSpec()
	podSpec.WithTerminationGracePeriodSeconds(0)
	podSpec.WithServiceAccountName(SERVICEACCOUNTNAME)

	podSpec.WithContainers(deploymentContainer)

	applyOptions := metav1.ApplyOptions{
		Force:        true,
		FieldManager: version.Name,
	}

	labelSelector := applyconfmeta.LabelSelector()
	labelSelector.WithMatchLabels(map[string]string{"app": version.Name})

	podTemplate := applyconfcore.PodTemplateSpec()
	podTemplate.WithLabels(map[string]string{
		"app": version.Name,
	})
	podTemplate.WithSpec(podSpec)

	deployment := applyconfapp.Deployment(version.Name, utils.CONFIG.Kubernetes.OwnNamespace)
	deployment.WithSpec(applyconfapp.DeploymentSpec().WithSelector(labelSelector).WithTemplate(podTemplate))

	// Create Deployment
	logger.Log.Info("Creating punq deployment ...")
	result, err := deploymentClient.Apply(context.TODO(), deployment, applyOptions)
	if err != nil {
		logger.Log.Error(err)
	}
	logger.Log.Info("Created punq deployment.", result.GetObjectMeta().GetName())
}
