package kubernetes

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/mogenius/punq/dtos"
	"github.com/mogenius/punq/version"

	"github.com/mogenius/punq/utils"

	"github.com/mogenius/punq/logger"

	core "k8s.io/api/core/v1"
	rbac "k8s.io/api/rbac/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	applyconfapp "k8s.io/client-go/applyconfigurations/apps/v1"
	applyconfcore "k8s.io/client-go/applyconfigurations/core/v1"
	applyconfmeta "k8s.io/client-go/applyconfigurations/meta/v1"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

func Deploy(clusterName string, ingressHostname string) {
	provider := NewKubeProvider(nil)
	if provider == nil {
		logger.Log.Fatal("Failed to load kubeprovider.")
	}

	applyNamespace(provider)
	addRbac(provider)
	addDeployment(provider)

	_, err := CreateClusterSecretIfNotExist(provider)
	if err != nil {
		logger.Log.Fatalf("Error creating cluster secret. Aborting: %s.", err.Error())
	}

	// adminUser, err := CreateUserSecretIfNotExist(provider)
	// if err != nil {
	// 	logger.Log.Fatalf("Error creating user secret. Aborting: %s.", err.Error())
	// }
	// if adminUser != nil {
	// 	structs.PrettyPrint(adminUser)
	// }

	_, err = CreateContextSecretIfNotExist(provider)
	if err != nil {
		logger.Log.Fatalf("Error creating context secret. Aborting: %s.", err.Error())
	}
	// if adminUser != nil {
	// 	fmt.Printf("Contexts saved (%d bytes).\n", len(ownContext.Context))
	// }

	if ingressHostname != "" {
		addService(provider)
		addIngress(provider, ingressHostname)
	}

	fmt.Printf("\n🚀🚀🚀 Successfuly installed punq in '%s'.\n\n", clusterName)
}

func addService(provider *KubeProvider) {
	fmt.Println("Creating punq service ...")

	punqService := utils.InitPunqService()
	punqService.ObjectMeta.Name = SERVICENAME
	punqService.ObjectMeta.Namespace = utils.CONFIG.Kubernetes.OwnNamespace
	punqService.Spec.Ports[0].Name = fmt.Sprintf("%d-%s-backend", utils.CONFIG.Backend.Port, SERVICENAME)
	punqService.Spec.Ports[0].Protocol = core.ProtocolTCP
	punqService.Spec.Ports[0].Port = int32(utils.CONFIG.Backend.Port)
	punqService.Spec.Ports[0].TargetPort = intstr.Parse(fmt.Sprint(utils.CONFIG.Backend.Port))
	punqService.Spec.Ports[1].Name = fmt.Sprintf("%d-%s-frontend", utils.CONFIG.Frontend.Port, SERVICENAME)
	punqService.Spec.Ports[1].Protocol = core.ProtocolTCP
	punqService.Spec.Ports[1].Port = int32(utils.CONFIG.Frontend.Port)
	punqService.Spec.Ports[1].TargetPort = intstr.Parse(fmt.Sprint(utils.CONFIG.Frontend.Port))
	punqService.Spec.Selector["app"] = version.Name

	serviceClient := provider.ClientSet.CoreV1().Services(utils.CONFIG.Kubernetes.OwnNamespace)
	_, err := serviceClient.Create(context.TODO(), &punqService, metav1.CreateOptions{})
	if err != nil {
		logger.Log.Fatalf("Service Creation Err: %s", err.Error())
	}

	fmt.Println("Created punq service. ✅")
}

func addIngress(provider *KubeProvider, ingressHostname string) {
	fmt.Printf("Creating punq ingress (%s) ...\n", ingressHostname)

	punqIngress := utils.InitPunqIngress()
	punqIngress.ObjectMeta.Name = INGRESSNAME
	punqIngress.ObjectMeta.Namespace = utils.CONFIG.Kubernetes.OwnNamespace
	punqIngress.Spec.Rules[0].Host = ingressHostname
	punqIngress.Spec.Rules[0].HTTP.Paths[0].Backend.Service.Name = SERVICENAME

	ingressClient := provider.ClientSet.NetworkingV1().Ingresses(utils.CONFIG.Kubernetes.OwnNamespace)
	_, err := ingressClient.Create(context.TODO(), &punqIngress, metav1.CreateOptions{})
	if err != nil {
		logger.Log.Fatalf("Ingress Creation Err: %s", err.Error())
	}

	fmt.Printf("Created punq ingress (%s). ✅\n", ingressHostname)
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
				APIGroups: []string{"", "*"},
				Resources: RBACRESOURCES,
				Verbs:     []string{"*"},
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
	fmt.Println("Creating punq RBAC ...")
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
	fmt.Println("Created punq RBAC. ✅")
	return nil
}

func applyNamespace(kubeProvider *KubeProvider) {
	serviceClient := kubeProvider.ClientSet.CoreV1().Namespaces()

	namespace := applyconfcore.Namespace(utils.CONFIG.Kubernetes.OwnNamespace)

	applyOptions := metav1.ApplyOptions{
		Force:        true,
		FieldManager: version.Name,
	}

	fmt.Println("Creating punq namespace ...")
	_, err := serviceClient.Apply(context.TODO(), namespace, applyOptions)
	if err != nil {
		logger.Log.Error(err)
	}
	fmt.Println("Created punq namespace. ✅")
}

func CreateClusterSecretIfNotExist(kubeProvider *KubeProvider) (utils.ClusterSecret, error) {
	secretClient := kubeProvider.ClientSet.CoreV1().Secrets(utils.CONFIG.Kubernetes.OwnNamespace)

	existingSecret, getErr := secretClient.Get(context.TODO(), utils.CONFIG.Kubernetes.OwnNamespace, metav1.GetOptions{})
	return writePunqSecret(secretClient, existingSecret, getErr)
}

func writePunqSecret(secretClient v1.SecretInterface, existingSecret *core.Secret, getErr error) (utils.ClusterSecret, error) {
	clusterSecret := utils.ClusterSecret{
		ApiKey:       utils.NanoIdExtraLong(),
		ClusterMfaId: utils.NanoIdExtraLong(),
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
		_, err := secretClient.Create(context.TODO(), &secret, MoCreateOptions())
		if err != nil {
			logger.Log.Error(err)
			return clusterSecret, err
		}
		fmt.Println("Created new punq secret. ✅")
	} else {
		if string(existingSecret.Data["api-key"]) != clusterSecret.ApiKey ||
			string(existingSecret.Data["cluster-name"]) != clusterSecret.ClusterName {
			fmt.Println("Updating existing punq secret ...")
			// keep existing mfa-id if possible
			if string(existingSecret.Data["cluster-mfa-id"]) != "" {
				clusterSecret.ClusterMfaId = string(existingSecret.Data["cluster-mfa-id"])
				secret.StringData["cluster-mfa-id"] = clusterSecret.ClusterMfaId
			}
			_, err := secretClient.Update(context.TODO(), &secret, MoUpdateOptions())
			if err != nil {
				logger.Log.Error(err)
				return clusterSecret, err
			}
			fmt.Println("Updated punq secret. ✅")
		} else {
			clusterSecret.ClusterMfaId = string(existingSecret.Data["cluster-mfa-id"])
			fmt.Println("Using existing punq secret. ✅")
		}
	}

	return clusterSecret, nil
}

// func CreateUserSecretIfNotExist(kubeProvider *KubeProvider) (*dtos.PunqUser, error) {
// 	secretClient := kubeProvider.ClientSet.CoreV1().Secrets(utils.CONFIG.Kubernetes.OwnNamespace)
//
// 	existingSecret, getErr := secretClient.Get(context.TODO(), utils.USERSSECRET, metav1.GetOptions{})
// 	return writeUserSecret(secretClient, existingSecret, getErr)
// }

// func writeUserSecret(secretClient v1.SecretInterface, existingSecret *core.Secret, getErr error) (*dtos.PunqUser, error) {
// 	password := utils.NanoId()
// 	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	adminUser := dtos.PunqUser{
// 		Id:          utils.NanoId(),
// 		Email:       "your-email@mogenius.com",
// 		Password:    string(hashedPassword),
// 		DisplayName: "Admin User",
// 		AccessLevel: dtos.ADMIN,
// 		Created:     time.Now().Format(time.RFC3339),
// 	}
//
// 	rawAdmin, err := json.Marshal(adminUser)
// 	if err != nil {
// 		logger.Log.Errorf("Error marshaling %s", err)
// 	}
//
// 	secret := utils.InitSecret()
// 	secret.ObjectMeta.Name = utils.USERSSECRET
// 	secret.ObjectMeta.Namespace = utils.CONFIG.Kubernetes.OwnNamespace
// 	delete(secret.StringData, "exampleData") // delete example data
// 	secret.StringData[adminUser.Id] = string(rawAdmin)
//
// 	if existingSecret == nil || getErr != nil {
// 		fmt.Println("Creating new punq-user secret ...")
// 		_, err := secretClient.Create(context.TODO(), &secret, MoCreateOptions())
// 		if err != nil {
// 			logger.Log.Error(err)
// 			return nil, err
// 		}
// 		fmt.Println("Created new punq-user secret. ✅")
// 		adminUser.Password = password
// 		return &adminUser, nil
// 	}
// 	return nil, nil
// }

func CreateContextSecretIfNotExist(kubeProvider *KubeProvider) (*dtos.PunqContext, error) {
	secretClient := kubeProvider.ClientSet.CoreV1().Secrets(utils.CONFIG.Kubernetes.OwnNamespace)

	existingSecret, getErr := secretClient.Get(context.TODO(), utils.CONTEXTSSECRET, metav1.GetOptions{})
	return writeContextSecret(secretClient, existingSecret, getErr)
}

func writeContextSecret(secretClient v1.SecretInterface, existingSecret *core.Secret, getErr error) (*dtos.PunqContext, error) {
	kubeconfigEnvVar := os.Getenv("KUBECONFIG")
	if kubeconfigEnvVar != "" {
		kubeconfigData, err := os.ReadFile(kubeconfigEnvVar)
		if err != nil {
			logger.Log.Fatalf("error reading kubeconfig: %s", err.Error())
		}

		ownContext := dtos.PunqContext{
			Id:      utils.CONTEXTOWN,
			Name:    utils.CONTEXTOWN,
			Context: string(kubeconfigData),
		}

		rawAdmin, err := json.Marshal(ownContext)
		if err != nil {
			logger.Log.Errorf("Error marshaling %s", err)
		}

		secret := utils.InitSecret()
		secret.ObjectMeta.Name = utils.CONTEXTSSECRET
		secret.ObjectMeta.Namespace = utils.CONFIG.Kubernetes.OwnNamespace
		delete(secret.StringData, "exampleData") // delete example data
		secret.StringData[utils.CONTEXTOWN] = string(rawAdmin)

		if existingSecret == nil || getErr != nil {
			fmt.Println("Creating new punq-context secret ...")
			_, err := secretClient.Create(context.TODO(), &secret, MoCreateOptions())
			if err != nil {
				logger.Log.Error(err)
				return nil, err
			}
			fmt.Println("Created new punq-context secret. ✅")
			return &ownContext, nil
		}
		return nil, nil
	} else {
		logger.Log.Fatal("$KUBECONFIG is empty. Cannot locate your context. Please use the -c flag to define the location of your kube-config.")
	}
	return nil, nil
}

func addDeployment(kubeProvider *KubeProvider) {
	deploymentClient := kubeProvider.ClientSet.AppsV1().Deployments(utils.CONFIG.Kubernetes.OwnNamespace)

	deploymentContainer := applyconfcore.Container()
	deploymentContainer.WithImagePullPolicy(core.PullAlways)
	deploymentContainer.WithName(version.Name)
	deploymentContainer.WithImage(DEPLOYMENTNAME())

	deploymentContainer.WithPorts(applyconfcore.ContainerPort().WithContainerPort(int32(utils.CONFIG.Backend.Port)).WithContainerPort(int32(utils.CONFIG.Frontend.Port)))

	// envVars := []applyconfcore.EnvVarApplyConfiguration{}
	// envVars = append(envVars, applyconfcore.EnvVarApplyConfiguration{
	// 	Name:  utils.Pointer("cluster_name"),
	// 	Value: utils.Pointer("TestClusterFromCode"),
	// })
	// envVars = append(envVars, applyconfcore.EnvVarApplyConfiguration{
	// 	Name:  utils.Pointer("api_key"),
	// 	Value: utils.Pointer("94E23575-A689-4F88-8D67-215A274F4E6E"), // dont worry. this is a test key
	// })
	// deploymentContainer.Env = envVars
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
	fmt.Println("Creating punq deployment ...")
	_, err := deploymentClient.Apply(context.TODO(), deployment, applyOptions)
	if err != nil {
		logger.Log.Error(err)
	}
	fmt.Println("Created punq deployment. ✅")
}
