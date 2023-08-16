package kubernetes

import (
	"context"
	"os"
	"time"

	"punq/logger"

	core "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	rbac "k8s.io/api/rbac/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	applyconfapp "k8s.io/client-go/applyconfigurations/apps/v1"
	applyconfcore "k8s.io/client-go/applyconfigurations/core/v1"
	applyconfmeta "k8s.io/client-go/applyconfigurations/meta/v1"
)

func Deploy() {
	provider, err := NewKubeProviderLocal()
	if err != nil {
		panic(err)
	}

	applyNamespace(provider)
	addRbac(provider)
	addRedisService(provider)
	addRedis(provider)
	addDaemonSet(provider)
	time.Sleep(3 * time.Second) // TODO: <-- this is realy dumb. find a better solution
	go StartPortForward(provider, false)
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
				Verbs:     []string{"list", "get", "watch"},
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
				Namespace: NAMESPACE,
			},
		},
	}

	// CREATE RBAC
	logger.Log.Info("Creating punq RBAC ...")
	_, err := kubeProvider.ClientSet.CoreV1().ServiceAccounts(NAMESPACE).Create(context.TODO(), serviceAccount, metav1.CreateOptions{})
	if err != nil && !k8serrors.IsAlreadyExists(err) {
		return err
	}
	_, err = kubeProvider.ClientSet.RbacV1().ClusterRoles().Create(context.TODO(), clusterRole, metav1.CreateOptions{})
	if err != nil && !k8serrors.IsAlreadyExists(err) {
		return err
	}
	_, err = kubeProvider.ClientSet.RbacV1().ClusterRoleBindings().Create(context.TODO(), clusterRoleBinding, metav1.CreateOptions{})
	if err != nil && !k8serrors.IsAlreadyExists(err) {
		return err
	}
	logger.Log.Info("Created punq RBAC.")
	return nil
}

func applyNamespace(kubeProvider *KubeProvider) {
	serviceClient := kubeProvider.ClientSet.CoreV1().Namespaces()

	namespace := applyconfcore.Namespace(NAMESPACE)

	applyOptions := metav1.ApplyOptions{
		Force:        true,
		FieldManager: REDISSERVICENAME,
	}

	logger.Log.Info("Creating punq namespace ...")
	result, err := serviceClient.Apply(context.TODO(), namespace, applyOptions)
	if err != nil {
		logger.Log.Error(err)
	}
	logger.Log.Info("Created punq namespace", result.GetObjectMeta().GetName(), ".")
}

func addDaemonSet(kubeProvider *KubeProvider) {
	daemonSetClient := kubeProvider.ClientSet.AppsV1().DaemonSets(NAMESPACE)

	daemonsetContainer := applyconfcore.Container()
	daemonsetContainer.WithName(DAEMONSETNAME)
	daemonsetContainer.WithImage(DAEMONSETIMAGE)
	daemonsetContainer.WithImagePullPolicy(core.PullAlways)
	daemonsetContainer.WithEnv(
		applyconfcore.EnvVar().WithName("STAGE").WithValue(os.Getenv("STAGE")),
		applyconfcore.EnvVar().WithName("REDIS_SERVICE_NAME").WithValue(os.Getenv("REDIS_SERVICE_NAME")),
		applyconfcore.EnvVar().WithName("REDIS_PORT").WithValue(os.Getenv("REDIS_PORT")),
		applyconfcore.EnvVar().WithName("API_HOST").WithValue(os.Getenv("API_HOST")),
		applyconfcore.EnvVar().WithName("API_PORT").WithValue(os.Getenv("API_PORT")),
		applyconfcore.EnvVar().WithName("OWN_NODE_NAME").WithValueFrom(
			applyconfcore.EnvVarSource().WithFieldRef(
				applyconfcore.ObjectFieldSelector().WithAPIVersion("v1").WithFieldPath("spec.nodeName"),
			),
		),
		applyconfcore.EnvVar().WithName("OWN_POD_NAME").WithValueFrom(
			applyconfcore.EnvVarSource().WithFieldRef(
				applyconfcore.ObjectFieldSelector().WithAPIVersion("v1").WithFieldPath("metadata.name"),
			),
		),
		applyconfcore.EnvVar().WithName("OWN_NAMESPACE").WithValueFrom(
			applyconfcore.EnvVarSource().WithFieldRef(
				applyconfcore.ObjectFieldSelector().WithAPIVersion("v1").WithFieldPath("metadata.namespace"),
			),
		),
	)

	caps := applyconfcore.Capabilities().WithDrop("ALL")

	caps = caps.WithAdd("NET_RAW")
	caps = caps.WithAdd("NET_ADMIN")
	caps = caps.WithAdd("SYS_ADMIN")
	caps = caps.WithAdd("SYS_PTRACE")
	caps = caps.WithAdd("DAC_OVERRIDE")
	caps = caps.WithAdd("SYS_RESOURCE")
	daemonsetContainer.WithSecurityContext(applyconfcore.SecurityContext().WithCapabilities(caps))

	agentResourceLimits := core.ResourceList{
		"cpu":               resource.MustParse("1000m"),
		"memory":            resource.MustParse("512Mi"),
		"ephemeral-storage": resource.MustParse("100Mi"),
	}
	agentResourceRequests := core.ResourceList{
		"cpu":               resource.MustParse("500m"),
		"memory":            resource.MustParse("128Mi"),
		"ephemeral-storage": resource.MustParse("10Mi"),
	}
	agentResources := applyconfcore.ResourceRequirements().WithRequests(agentResourceRequests).WithLimits(agentResourceLimits)
	daemonsetContainer.WithResources(agentResources)

	// Host procfs is needed inside the container because we need access to
	//	the network namespaces of processes on the machine.
	//
	procfsVolume := applyconfcore.Volume()
	procfsVolume.WithName(PROCFSVOLUMENAME).WithHostPath(applyconfcore.HostPathVolumeSource().WithPath("/proc"))
	procfsVolumeMount := applyconfcore.VolumeMount().WithName(PROCFSVOLUMENAME).WithMountPath(PROCFSMOUNTPATH).WithReadOnly(true)
	daemonsetContainer.WithVolumeMounts(procfsVolumeMount)

	// We need access to /sys in order to install certain eBPF tracepoints
	//
	sysfsVolume := applyconfcore.Volume()
	sysfsVolume.WithName(SYSFSVOLUMENAME).WithHostPath(applyconfcore.HostPathVolumeSource().WithPath("/sys"))
	sysfsVolumeMount := applyconfcore.VolumeMount().WithName(SYSFSVOLUMENAME).WithMountPath(SYSFSMOUNTPATH).WithReadOnly(true)
	daemonsetContainer.WithVolumeMounts(sysfsVolumeMount)

	podSpec := applyconfcore.PodSpec()
	podSpec.WithHostNetwork(true)
	podSpec.WithDNSPolicy(core.DNSClusterFirstWithHostNet)
	podSpec.WithTerminationGracePeriodSeconds(0)
	podSpec.WithServiceAccountName(SERVICEACCOUNTNAME)

	podSpec.WithContainers(daemonsetContainer)
	podSpec.WithVolumes(procfsVolume, sysfsVolume)

	applyOptions := metav1.ApplyOptions{
		Force:        true,
		FieldManager: DAEMONSETNAME,
	}

	labelSelector := applyconfmeta.LabelSelector()
	labelSelector.WithMatchLabels(map[string]string{"app": DAEMONSETNAME})

	podTemplate := applyconfcore.PodTemplateSpec()
	podTemplate.WithLabels(map[string]string{
		"app": DAEMONSETNAME,
	})
	podTemplate.WithSpec(podSpec)

	daemonSet := applyconfapp.DaemonSet(DAEMONSETNAME, NAMESPACE)
	daemonSet.WithSpec(applyconfapp.DaemonSetSpec().WithSelector(labelSelector).WithTemplate(podTemplate))

	// Create DaemonSet
	logger.Log.Info("Creating punq daemonset ...")
	result, err := daemonSetClient.Apply(context.TODO(), daemonSet, applyOptions)
	if err != nil {
		logger.Log.Error(err)
	}
	logger.Log.Info("Created punq daemonset.", result.GetObjectMeta().GetName(), ".")
}

func addRedis(kubeProvider *KubeProvider) {
	deploymentClient := kubeProvider.ClientSet.AppsV1().Deployments(NAMESPACE)

	deploymentContainer := applyconfcore.Container()
	deploymentContainer.WithName(REDISNAME)
	deploymentContainer.WithImage(REDISIMAGE)
	deploymentContainer.WithEnv(
		applyconfcore.EnvVar().WithName("STAGE").WithValue(os.Getenv("STAGE")),
	)
	agentResourceLimits := core.ResourceList{
		"cpu":               resource.MustParse("300m"),
		"memory":            resource.MustParse("256Mi"),
		"ephemeral-storage": resource.MustParse("100Mi"),
	}
	agentResourceRequests := core.ResourceList{
		"cpu":               resource.MustParse("100m"),
		"memory":            resource.MustParse("128Mi"),
		"ephemeral-storage": resource.MustParse("10Mi"),
	}
	agentResources := applyconfcore.ResourceRequirements().WithRequests(agentResourceRequests).WithLimits(agentResourceLimits)
	deploymentContainer.WithResources(agentResources)
	deploymentContainer.WithPorts(applyconfcore.ContainerPort().WithContainerPort(REDISPORT).WithProtocol(v1.ProtocolTCP).WithName("redis"))

	podSpec := applyconfcore.PodSpec()
	podSpec.WithTerminationGracePeriodSeconds(0)
	podSpec.WithServiceAccountName(SERVICEACCOUNTNAME)

	podSpec.WithContainers(deploymentContainer)

	applyOptions := metav1.ApplyOptions{
		Force:        true,
		FieldManager: REDISNAME,
	}

	labelSelector := applyconfmeta.LabelSelector()
	labelSelector.WithMatchLabels(map[string]string{"app": REDISNAME})

	podTemplate := applyconfcore.PodTemplateSpec()
	podTemplate.WithLabels(map[string]string{
		"app": REDISNAME,
	})
	podTemplate.WithSpec(podSpec)

	deployment := applyconfapp.Deployment(REDISNAME, NAMESPACE)
	deployment.WithSpec(applyconfapp.DeploymentSpec().WithSelector(labelSelector).WithTemplate(podTemplate))

	// Create Redis Deployment
	logger.Log.Info("Creating punq redis ...")
	result, err := deploymentClient.Apply(context.TODO(), deployment, applyOptions)
	if err != nil {
		logger.Log.Error(err)
	}
	logger.Log.Info("Created punq redis.", result.GetObjectMeta().GetName(), ".")
}

func addRedisService(kubeProvider *KubeProvider) {
	serviceClient := kubeProvider.ClientSet.CoreV1().Services(NAMESPACE)

	serviceSpec := applyconfcore.ServiceSpec()
	serviceSpec.WithSelector(map[string]string{"app": REDISNAME})
	serviceSpec.WithPorts(applyconfcore.ServicePort().WithPort(REDISPORT).WithTargetPort(intstr.FromString(REDISTARGETPORT)).WithProtocol(v1.ProtocolTCP))

	applyOptions := metav1.ApplyOptions{
		Force:        true,
		FieldManager: REDISSERVICENAME,
	}

	service := applyconfcore.Service(REDISSERVICENAME, NAMESPACE)
	service.WithSpec(serviceSpec)

	// Create Redis Deployment
	logger.Log.Info("Creating punq redis service ...")
	result, err := serviceClient.Apply(context.TODO(), service, applyOptions)
	if err != nil {
		logger.Log.Error(err)
	}
	logger.Log.Info("Created punqdis service", result.GetObjectMeta().GetName(), ".")
}
