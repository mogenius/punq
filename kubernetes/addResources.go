package kubernetes

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

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
	_ = utils.GetDefaultKubeConfig()

	provider, err := NewKubeProvider(nil)
	if provider == nil || err != nil {
		logger.Log.Fatal("Failed to load provider.")
	}

	applyNamespace(provider)
	addRbac(provider)
	addDeployment(provider)

	_, err = CreateContextSecretIfNotExist(provider)
	if err != nil {
		logger.Log.Fatalf("Error creating context secret. Aborting: %s.", err.Error())
	}

	if ingressHostname != "" {
		addService(provider)
		addIngress(provider, clusterName, ingressHostname)
	}

	fmt.Printf("\n🚀🚀🚀 Successfully installed punq in '%s'.\n\n", clusterName)
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
	punqService.Spec.Ports[2].Name = fmt.Sprintf("%d-%s-websocket", utils.CONFIG.Websocket.Port, SERVICENAME)
	punqService.Spec.Ports[2].Protocol = core.ProtocolTCP
	punqService.Spec.Ports[2].Port = int32(utils.CONFIG.Websocket.Port)
	punqService.Spec.Ports[2].TargetPort = intstr.Parse(fmt.Sprint(utils.CONFIG.Websocket.Port))
	punqService.Spec.Selector["app"] = version.Name

	serviceClient := provider.ClientSet.CoreV1().Services(utils.CONFIG.Kubernetes.OwnNamespace)
	_, err := serviceClient.Create(context.TODO(), &punqService, metav1.CreateOptions{})
	if err != nil {
		logger.Log.Fatalf("Service Creation Err: %s", err.Error())
	}

	fmt.Println("Created punq service. ✅")
}

func addIngress(provider *KubeProvider, clusterName string, ingressHostname string) {
	// 1. Determine IngressType
	controllerType, err := DetermineIngressControllerType(nil)
	switch controllerType {
	case NGINX:
		addNginxIngress(provider, ingressHostname)
	case TRAEFIK:
		addTraefikIngress(provider, ingressHostname)
		addTraefikMiddleware(provider, ingressHostname)
	case NONE:
		// clean everything
		Remove(clusterName)
		utils.FatalError("No ingress controller found.\nWe recomend installing TRAEFIK:\n  helm repo add traefik https://traefik.github.io/charts\n  helm install traefik traefik/traefik\nAfter installing TRAEFIK, you can retry installing punq.")
	case MULTIPLE:
		utils.FatalError(err.Error())
	}
}

func addTraefikIngress(provider *KubeProvider, ingressHostname string) {
	fmt.Printf("Creating TRAEFIK punq ingress (%s) ...\n", ingressHostname)
	punqIngress := utils.InitPunqIngressTraefik()
	punqIngress.ObjectMeta.Name = INGRESSNAME
	punqIngress.ObjectMeta.Namespace = utils.CONFIG.Kubernetes.OwnNamespace
	punqIngress.Spec.Rules[0].Host = ingressHostname
	punqIngress.Spec.Rules[0].HTTP.Paths[0].Backend.Service.Name = SERVICENAME
	punqIngress.Spec.Rules[0].HTTP.Paths[0].Backend.Service.Port.Number = int32(utils.CONFIG.Backend.Port)
	punqIngress.Spec.Rules[0].HTTP.Paths[1].Backend.Service.Name = SERVICENAME
	punqIngress.Spec.Rules[0].HTTP.Paths[1].Backend.Service.Port.Number = int32(utils.CONFIG.Websocket.Port)
	punqIngress.Spec.Rules[0].HTTP.Paths[2].Backend.Service.Name = SERVICENAME
	punqIngress.Spec.Rules[0].HTTP.Paths[2].Backend.Service.Port.Number = int32(utils.CONFIG.Frontend.Port)

	ingressClient := provider.ClientSet.NetworkingV1().Ingresses(utils.CONFIG.Kubernetes.OwnNamespace)
	_, err := ingressClient.Create(context.TODO(), &punqIngress, metav1.CreateOptions{})
	if err != nil {
		logger.Log.Fatalf("Ingress TRAEFIK Creation Err: %s", err.Error())
	}
	fmt.Printf("Created TRAEFIK punq ingress (%s). ✅\n", ingressHostname)
}

func addTraefikMiddleware(provider *KubeProvider, ingressHostname string) {
	fmt.Printf("Creating TRAEFIK middleware (%s) ...\n", ingressHostname)
	mwYaml := utils.InitPunqIngressTraefikMiddlewareYaml()

	cmd := exec.Command("bash", "-c", fmt.Sprintf("echo \"%s\" | kubectl %s apply -f -", mwYaml, ContextFlag(nil)))

	output, err := cmd.CombinedOutput()
	if err != nil {
		utils.FatalError(fmt.Sprintf("failed to execute command (%s): %s\n%s", cmd.String(), err.Error(), string(output)))

	}
	fmt.Printf("Created TRAEFIK middleware (%s). ✅\n", ingressHostname)
}

func addNginxIngress(provider *KubeProvider, ingressHostname string) {
	fmt.Printf("Creating NGINX punq ingress (%s) ...\n", ingressHostname)
	punqIngress := utils.InitPunqIngress()
	punqIngress.ObjectMeta.Name = INGRESSNAME
	punqIngress.ObjectMeta.Namespace = utils.CONFIG.Kubernetes.OwnNamespace
	punqIngress.Spec.Rules[0].Host = ingressHostname
	punqIngress.Spec.Rules[0].HTTP.Paths[0].Backend.Service.Name = SERVICENAME
	punqIngress.Spec.Rules[0].HTTP.Paths[0].Backend.Service.Port.Number = int32(utils.CONFIG.Backend.Port)
	punqIngress.Spec.Rules[0].HTTP.Paths[1].Backend.Service.Name = SERVICENAME
	punqIngress.Spec.Rules[0].HTTP.Paths[1].Backend.Service.Port.Number = int32(utils.CONFIG.Websocket.Port)
	punqIngress.Spec.Rules[0].HTTP.Paths[2].Backend.Service.Name = SERVICENAME
	punqIngress.Spec.Rules[0].HTTP.Paths[2].Backend.Service.Port.Number = int32(utils.CONFIG.Frontend.Port)

	ingressClient := provider.ClientSet.NetworkingV1().Ingresses(utils.CONFIG.Kubernetes.OwnNamespace)
	_, err := ingressClient.Create(context.TODO(), &punqIngress, metav1.CreateOptions{})
	if err != nil {
		logger.Log.Fatalf("Ingress Creation Err: %s", err.Error())
	}
	fmt.Printf("Created NGINX punq ingress (%s). ✅\n", ingressHostname)
}

func addRbac(provider *KubeProvider) error {
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
	_, err := provider.ClientSet.CoreV1().ServiceAccounts(utils.CONFIG.Kubernetes.OwnNamespace).Create(context.TODO(), serviceAccount, MoCreateOptions())
	if err != nil && !k8serrors.IsAlreadyExists(err) {
		return err
	}
	_, err = provider.ClientSet.RbacV1().ClusterRoles().Create(context.TODO(), clusterRole, MoCreateOptions())
	if err != nil && !k8serrors.IsAlreadyExists(err) {
		return err
	}
	_, err = provider.ClientSet.RbacV1().ClusterRoleBindings().Create(context.TODO(), clusterRoleBinding, MoCreateOptions())
	if err != nil && !k8serrors.IsAlreadyExists(err) {
		return err
	}
	fmt.Println("Created punq RBAC. ✅")
	return nil
}

func applyNamespace(provider *KubeProvider) {
	serviceClient := provider.ClientSet.CoreV1().Namespaces()

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

func CreateContextSecretIfNotExist(provider *KubeProvider) (*dtos.PunqContext, error) {
	secretClient := provider.ClientSet.CoreV1().Secrets(utils.CONFIG.Kubernetes.OwnNamespace)

	existingSecret, getErr := secretClient.Get(context.TODO(), utils.CONTEXTSSECRET, metav1.GetOptions{})
	return writeContextSecret(secretClient, existingSecret, getErr)
}

func writeContextSecret(secretClient v1.SecretInterface, existingSecret *core.Secret, getErr error) (*dtos.PunqContext, error) {
	kubeconfigEnvVar := utils.GetDefaultKubeConfig()

	kubeconfigData, err := os.ReadFile(kubeconfigEnvVar)
	if err != nil {
		logger.Log.Fatalf("error reading kubeconfig: %s", err.Error())
	}

	ownContext, err := dtos.ParseCurrentContextConfigToPunqContext(kubeconfigData)
	if err != nil {
		utils.FatalError(err.Error())
	}

	fmt.Println("Determining cluster provider ...")
	ownProvider, err := GuessClusterProvider(nil)
	if err == nil {
		ownContext.Provider = string(ownProvider)
		fmt.Printf("Determined cluster provider: '%s'. ✅\n", ownProvider)
	} else {
		fmt.Println("Determining cluster provider failed.")
	}

	ownContext.Id = utils.CONTEXTOWN
	ownContext.Name = utils.CONTEXTOWN

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
}

func addDeployment(provider *KubeProvider) {
	deploymentClient := provider.ClientSet.AppsV1().Deployments(utils.CONFIG.Kubernetes.OwnNamespace)

	deploymentContainer := applyconfcore.Container()
	deploymentContainer.WithImagePullPolicy(core.PullAlways)
	deploymentContainer.WithName(version.Name)
	deploymentContainer.WithImage(version.OperatorImage)

	deploymentContainer.WithPorts(applyconfcore.ContainerPort().WithContainerPort(int32(utils.CONFIG.Backend.Port)).WithContainerPort(int32(utils.CONFIG.Frontend.Port)).WithContainerPort(int32(utils.CONFIG.Websocket.Port)))

	envVars := []applyconfcore.EnvVarApplyConfiguration{}
	envVars = append(envVars, applyconfcore.EnvVarApplyConfiguration{
		Name:  utils.Pointer("stage"),
		Value: utils.Pointer("operator"),
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
	fmt.Println("Creating punq deployment ...")
	_, err := deploymentClient.Apply(context.TODO(), deployment, applyOptions)
	if err != nil {
		logger.Log.Error(err)
	}
	fmt.Println("Created punq deployment. ✅")
}
