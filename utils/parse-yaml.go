package utils

import (
	"embed"

	cmapi "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	v1 "k8s.io/api/apps/v1"
	v1job "k8s.io/api/batch/v1"
	core "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/kubectl/pkg/scheme"
)

//go:embed yaml-templates
var YamlTemplatesFolder embed.FS

func InitPersistentVolumeYaml() string {
	yaml, err := YamlTemplatesFolder.ReadFile("yaml-templates/volume-example.yaml")
	if err != nil {
		return err.Error()
	}
	return string(yaml)
}
func InitPersistentVolume() core.PersistentVolume {
	yaml, err := YamlTemplatesFolder.ReadFile("yaml-templates/volume-nfs-pv.yaml")
	if err != nil {
		panic(err.Error())
	}

	s := json.NewYAMLSerializer(json.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)

	var app core.PersistentVolume
	_, _, err = s.Decode(yaml, nil, &app)
	if err != nil {
		panic(err)
	}
	return app
}

func InitPersistentVolumeClaimYaml() string {
	yaml, err := YamlTemplatesFolder.ReadFile("yaml-templates/pvc-example.yaml")
	if err != nil {
		return err.Error()
	}
	return string(yaml)
}
func InitPersistentVolumeClaim() core.PersistentVolumeClaim {
	yaml, err := YamlTemplatesFolder.ReadFile("yaml-templates/volumeclaim-cephfs.yaml")
	if err != nil {
		panic(err.Error())
	}

	s := json.NewYAMLSerializer(json.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)

	var app core.PersistentVolumeClaim
	_, _, err = s.Decode(yaml, nil, &app)
	if err != nil {
		panic(err)
	}
	return app
}

func InitContainerSecretYaml() string {
	yaml, err := YamlTemplatesFolder.ReadFile("yaml-templates/container-secret.yaml")
	if err != nil {
		return err.Error()
	}
	return string(yaml)
}
func InitContainerSecret() core.Secret {
	yaml, err := YamlTemplatesFolder.ReadFile("yaml-templates/container-secret.yaml")
	if err != nil {
		panic(err.Error())
	}

	s := json.NewYAMLSerializer(json.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)

	var app core.Secret
	_, _, err = s.Decode(yaml, nil, &app)
	if err != nil {
		panic(err)
	}
	return app
}

func InitSecretYaml() string {
	yaml, err := YamlTemplatesFolder.ReadFile("yaml-templates/secret-sample.yaml")
	if err != nil {
		return err.Error()
	}
	return string(yaml)
}
func InitSecret() core.Secret {
	yaml, err := YamlTemplatesFolder.ReadFile("yaml-templates/secret.yaml")
	if err != nil {
		panic(err.Error())
	}

	s := json.NewYAMLSerializer(json.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)

	var app core.Secret
	_, _, err = s.Decode(yaml, nil, &app)
	if err != nil {
		panic(err)
	}
	return app
}

func InitConfigMapYaml() string {
	yaml, err := YamlTemplatesFolder.ReadFile("yaml-templates/configmap-example.yaml")
	if err != nil {
		return err.Error()
	}
	return string(yaml)
}
func InitConfigMap() core.ConfigMap {
	yaml, err := YamlTemplatesFolder.ReadFile("yaml-templates/configmap.yaml")
	if err != nil {
		panic(err.Error())
	}

	s := json.NewYAMLSerializer(json.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)

	var app core.ConfigMap
	_, _, err = s.Decode(yaml, nil, &app)
	if err != nil {
		panic(err)
	}
	return app
}

func InitDeploymentYaml() string {
	yaml, err := YamlTemplatesFolder.ReadFile("yaml-templates/deployment-example.yaml")
	if err != nil {
		return err.Error()
	}
	return string(yaml)
}
func InitDeployment() v1.Deployment {
	yaml, err := YamlTemplatesFolder.ReadFile("yaml-templates/deployment.yaml")
	if err != nil {
		panic(err.Error())
	}

	s := json.NewYAMLSerializer(json.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)

	var app v1.Deployment
	_, _, err = s.Decode(yaml, nil, &app)
	if err != nil {
		panic(err)
	}
	return app
}

func InitIngressYaml() string {
	yaml, err := YamlTemplatesFolder.ReadFile("yaml-templates/ingress.yaml")
	if err != nil {
		return err.Error()
	}
	return string(yaml)
}
func InitIngress() netv1.Ingress {
	yaml, err := YamlTemplatesFolder.ReadFile("yaml-templates/ingress.yaml")
	if err != nil {
		panic(err.Error())
	}

	s := json.NewYAMLSerializer(json.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)

	var app netv1.Ingress
	_, _, err = s.Decode(yaml, nil, &app)
	if err != nil {
		panic(err)
	}
	return app
}
func InitPunqIngress() netv1.Ingress {
	yaml, err := YamlTemplatesFolder.ReadFile("yaml-templates/punq-ingress.yaml")
	if err != nil {
		panic(err.Error())
	}

	s := json.NewYAMLSerializer(json.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)

	var app netv1.Ingress
	_, _, err = s.Decode(yaml, nil, &app)
	if err != nil {
		panic(err)
	}
	return app
}

func InitNetPolNamespaceYaml() string {
	yaml, err := YamlTemplatesFolder.ReadFile("yaml-templates/network-policy-namespace.yaml")
	if err != nil {
		return err.Error()
	}
	return string(yaml)
}
func InitNetPolNamespace() netv1.NetworkPolicy {
	yaml, err := YamlTemplatesFolder.ReadFile("yaml-templates/network-policy-namespace.yaml")
	if err != nil {
		panic(err.Error())
	}

	s := json.NewYAMLSerializer(json.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)

	var app netv1.NetworkPolicy
	_, _, err = s.Decode(yaml, nil, &app)
	if err != nil {
		panic(err)
	}
	return app
}

func InitNetPolServiceYaml() string {
	yaml, err := YamlTemplatesFolder.ReadFile("yaml-templates/network-policy-service.yaml")
	if err != nil {
		return err.Error()
	}
	return string(yaml)
}
func InitNetPolService() netv1.NetworkPolicy {
	yaml, err := YamlTemplatesFolder.ReadFile("yaml-templates/network-policy-service.yaml")
	if err != nil {
		panic(err.Error())
	}

	s := json.NewYAMLSerializer(json.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)

	var app netv1.NetworkPolicy
	_, _, err = s.Decode(yaml, nil, &app)
	if err != nil {
		panic(err)
	}
	return app
}

func InitCertificateYaml() string {
	yaml, err := YamlTemplatesFolder.ReadFile("yaml-templates/certificate-example.yaml")
	if err != nil {
		return err.Error()
	}
	return string(yaml)
}
func InitCertificate() cmapi.Certificate {
	yaml, err := YamlTemplatesFolder.ReadFile("yaml-templates/certificate.yaml")
	if err != nil {
		panic(err.Error())
	}

	s := json.NewYAMLSerializer(json.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)

	var app cmapi.Certificate
	_, _, err = s.Decode(yaml, nil, &app)
	if err != nil {
		panic(err)
	}
	return app
}

func InitServiceYaml() string {
	yaml, err := YamlTemplatesFolder.ReadFile("yaml-templates/service.yaml")
	if err != nil {
		return err.Error()
	}
	return string(yaml)
}
func InitService() core.Service {
	yaml, err := YamlTemplatesFolder.ReadFile("yaml-templates/service.yaml")
	if err != nil {
		panic(err.Error())
	}

	s := json.NewYAMLSerializer(json.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)

	var app core.Service
	_, _, err = s.Decode(yaml, nil, &app)
	if err != nil {
		panic(err)
	}
	return app
}
func InitPunqService() core.Service {
	yaml, err := YamlTemplatesFolder.ReadFile("yaml-templates/punq-service.yaml")
	if err != nil {
		panic(err.Error())
	}

	s := json.NewYAMLSerializer(json.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)

	var app core.Service
	_, _, err = s.Decode(yaml, nil, &app)
	if err != nil {
		panic(err)
	}
	return app
}

func InitUpgradeConfigMap() core.ConfigMap {
	yaml, err := YamlTemplatesFolder.ReadFile("yaml-templates/upgrade-configmap.yaml")
	if err != nil {
		panic(err.Error())
	}

	s := json.NewYAMLSerializer(json.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)

	var app core.ConfigMap
	_, _, err = s.Decode(yaml, nil, &app)
	if err != nil {
		panic(err)
	}
	return app
}

func InitUpgradeJob() v1job.Job {
	yaml, err := YamlTemplatesFolder.ReadFile("yaml-templates/upgrade-job.yaml")
	if err != nil {
		panic(err.Error())
	}

	s := json.NewYAMLSerializer(json.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)

	var app v1job.Job
	_, _, err = s.Decode(yaml, nil, &app)
	if err != nil {
		panic(err)
	}
	return app
}

func InitVolumeAttachmentYaml() string {
	yaml, err := YamlTemplatesFolder.ReadFile("yaml-templates/volumeattachment.yaml")
	if err != nil {
		return err.Error()
	}
	return string(yaml)
}

func InitStorageClassYaml() string {
	yaml, err := YamlTemplatesFolder.ReadFile("yaml-templates/storageclass.yaml")
	if err != nil {
		return err.Error()
	}
	return string(yaml)
}

func InitStatefulsetYaml() string {
	yaml, err := YamlTemplatesFolder.ReadFile("yaml-templates/statefulset.yaml")
	if err != nil {
		return err.Error()
	}
	return string(yaml)
}

func InitServiceExampleYaml() string {
	yaml, err := YamlTemplatesFolder.ReadFile("yaml-templates/service-example.yaml")
	if err != nil {
		return err.Error()
	}
	return string(yaml)
}

func InitServiceAccountExampleYaml() string {
	yaml, err := YamlTemplatesFolder.ReadFile("yaml-templates/serviceaccount-example.yaml")
	if err != nil {
		return err.Error()
	}
	return string(yaml)
}

func InitRoleYaml() string {
	yaml, err := YamlTemplatesFolder.ReadFile("yaml-templates/role-example.yaml")
	if err != nil {
		return err.Error()
	}
	return string(yaml)
}

func InitRoleBindingYaml() string {
	yaml, err := YamlTemplatesFolder.ReadFile("yaml-templates/rolebinding-example.yaml")
	if err != nil {
		return err.Error()
	}
	return string(yaml)
}

func InitReplicaSetYaml() string {
	yaml, err := YamlTemplatesFolder.ReadFile("yaml-templates/replicaset-example.yaml")
	if err != nil {
		return err.Error()
	}
	return string(yaml)
}

func InitPodYaml() string {
	yaml, err := YamlTemplatesFolder.ReadFile("yaml-templates/pod-example.yaml")
	if err != nil {
		return err.Error()
	}
	return string(yaml)
}

func InitOrderYaml() string {
	yaml, err := YamlTemplatesFolder.ReadFile("yaml-templates/order-example.yaml")
	if err != nil {
		return err.Error()
	}
	return string(yaml)
}

func InitNetPolYaml() string {
	yaml, err := YamlTemplatesFolder.ReadFile("yaml-templates/netpol-example.yaml")
	if err != nil {
		return err.Error()
	}
	return string(yaml)
}

func InitNamespaceYaml() string {
	yaml, err := YamlTemplatesFolder.ReadFile("yaml-templates/namespace-example.yaml")
	if err != nil {
		return err.Error()
	}
	return string(yaml)
}

func InitJobYaml() string {
	yaml, err := YamlTemplatesFolder.ReadFile("yaml-templates/job-example.yaml")
	if err != nil {
		return err.Error()
	}
	return string(yaml)
}

func InitIssuerYaml() string {
	yaml, err := YamlTemplatesFolder.ReadFile("yaml-templates/issuer-example.yaml")
	if err != nil {
		return err.Error()
	}
	return string(yaml)
}

func InitIngresYaml() string {
	yaml, err := YamlTemplatesFolder.ReadFile("yaml-templates/ingress-example.yaml")
	if err != nil {
		return err.Error()
	}
	return string(yaml)
}

func InitHpaYaml() string {
	yaml, err := YamlTemplatesFolder.ReadFile("yaml-templates/hpa-example.yaml")
	if err != nil {
		return err.Error()
	}
	return string(yaml)
}

func InitDaemonsetYaml() string {
	yaml, err := YamlTemplatesFolder.ReadFile("yaml-templates/daemonset-example.yaml")
	if err != nil {
		return err.Error()
	}
	return string(yaml)
}

func InitCertificateSigningRequestYaml() string {
	yaml, err := YamlTemplatesFolder.ReadFile("yaml-templates/csr-example.yaml")
	if err != nil {
		return err.Error()
	}
	return string(yaml)
}

func InitCronJobYaml() string {
	yaml, err := YamlTemplatesFolder.ReadFile("yaml-templates/cronjob-example.yaml")
	if err != nil {
		return err.Error()
	}
	return string(yaml)
}

func InitCronJob() v1job.CronJob {
	yaml, err := YamlTemplatesFolder.ReadFile("yaml-templates/cronjob-example.yaml")
	if err != nil {
		panic(err.Error())
	}

	s := json.NewYAMLSerializer(json.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)

	var app v1job.CronJob
	_, _, err = s.Decode(yaml, nil, &app)
	if err != nil {
		panic(err)
	}
	return app
}

func InitClusterIssuerYaml() string {
	yaml, err := YamlTemplatesFolder.ReadFile("yaml-templates/cluster-issuer-example.yaml")
	if err != nil {
		return err.Error()
	}
	return string(yaml)
}

func InitClusterRoleYaml() string {
	yaml, err := YamlTemplatesFolder.ReadFile("yaml-templates/cluster-role-example.yaml")
	if err != nil {
		return err.Error()
	}
	return string(yaml)
}

func InitClusterRoleBindingYaml() string {
	yaml, err := YamlTemplatesFolder.ReadFile("yaml-templates/cluster-role-binding-example.yaml")
	if err != nil {
		return err.Error()
	}
	return string(yaml)
}

func InitLeaseYaml() string {
	yaml, err := YamlTemplatesFolder.ReadFile("yaml-templates/lease-example.yaml")
	if err != nil {
		return err.Error()
	}
	return string(yaml)
}

func InitCustomResourceDefinitionYaml() string {
	yaml, err := YamlTemplatesFolder.ReadFile("yaml-templates/crd-example.yaml")
	if err != nil {
		return err.Error()
	}
	return string(yaml)
}

func InitPriorityClassYaml() string {
	yaml, err := YamlTemplatesFolder.ReadFile("yaml-templates/priorityclass-example.yaml")
	if err != nil {
		return err.Error()
	}
	return string(yaml)
}

func InitEndPointYaml() string {
	yaml, err := YamlTemplatesFolder.ReadFile("yaml-templates/endpoint-example.yaml")
	if err != nil {
		return err.Error()
	}
	return string(yaml)
}

func InitResourceQuotaYaml() string {
	yaml, err := YamlTemplatesFolder.ReadFile("yaml-templates/resourcequota-example.yaml")
	if err != nil {
		return err.Error()
	}
	return string(yaml)
}

func InitVolumeSnapshotYaml() string {
	yaml, err := YamlTemplatesFolder.ReadFile("yaml-templates/volumesnapshot-example.yaml")
	if err != nil {
		return err.Error()
	}
	return string(yaml)
}
