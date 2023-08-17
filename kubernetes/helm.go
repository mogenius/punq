package kubernetes

// func InstallMogeniusNfsStorage(job *structs.Job, clusterProvider string, wg *sync.WaitGroup) []*structs.Command {
// 	cmds := []*structs.Command{}

// 	addRepoCmd := structs.CreateBashCommand("Install/Update helm repo.", job, "sleep 1", wg)
// 	cmds = append(cmds, addRepoCmd)

// 	nfsStorageClassStr := ""

// 	// "BRING_YOUR_OWN", "EKS", "AKS", "GKE", "DOCKER_ENTERPRISE", "DOKS", "LINODE", "IBM", "ACK", "OKE", "OTC", "OPEN_SHIFT"
// 	switch clusterProvider {
// 	case "EKS":
// 		nfsStorageClassStr = " --set-string nfsStorageClass.backendStorageClass=gp2"
// 	case "GKE":
// 		nfsStorageClassStr = " --set-string nfsStorageClass.backendStorageClass=standard-rwo"
// 	case "AKS":
// 		nfsStorageClassStr = " --set-string nfsStorageClass.backendStorageClass=default"
// 	case "OTC":
// 		nfsStorageClassStr = " --set-string nfsStorageClass.backendStorageClass=csi-disk"
// 	default:
// 		// nothing to do
// 		errMsg := fmt.Sprintf("CLUSTERPROVIDER '%s' HAS NOT BEEN TESTED YET!", clusterProvider)
// 		logger.Log.Errorf(errMsg)
// 		addRepoCmd.Fail(errMsg)
// 		return cmds
// 	}
// 	instRelCmd := structs.CreateBashCommand("Install helm release.", job, fmt.Sprintf("helm repo add mo-openebs-nfs https://openebs.github.io/dynamic-nfs-provisioner; helm repo update; helm install mogenius-nfs-storage mo-openebs-nfs/nfs-provisioner -n %s --set analytics.enabled=false%s", utils.CONFIG.Kubernetes.OwnNamespace, nfsStorageClassStr), wg)
// 	cmds = append(cmds, instRelCmd)

// 	return cmds
// }

// func UninstallMogeniusNfsStorage(job *structs.Job, wg *sync.WaitGroup) []*structs.Command {
// 	cmds := []*structs.Command{}

// 	uninstRelCmd := structs.CreateBashCommand("Uninstall helm release.", job, fmt.Sprintf("helm uninstall mogenius-nfs-storage -n %s", utils.CONFIG.Kubernetes.OwnNamespace), wg)
// 	cmds = append(cmds, uninstRelCmd)
// 	// storageClassCmd := DeleteMogeniusNfsStorageClass(job, c, wg)
// 	// cmds = append(cmds, storageClassCmd)

// 	return cmds
// }

// func AllHelmCharts(namespaceName string) []cmapi.CertificateRequest {
// 	result := []cmapi.CertificateRequest{}

// 	provider := NewKubeProvider()
// 	resources, err := provider.ClientSet.Discovery().ServerPreferredResources()
// 	if err != nil {
// 		logger.Log.Errorf("ServerPreferredResources Error: %s", err.Error())
// 		return result
// 	}

// 	for _, certificate := range certificatesList.Items {
// 		if !utils.Contains(utils.CONFIG.Misc.IgnoreNamespaces, certificate.ObjectMeta.Namespace) {
// 			result = append(result, certificate)
// 		}
// 	}
// 	return result
// }

// func UpdateHelmChart(data cmapi.CertificateRequest) K8sWorkloadResult {
// 	kubeProvider := NewKubeProviderCertManager()
// 	certificateClient := kubeProvider.ClientSet.CertmanagerV1().CertificateRequests(data.Namespace)
// 	_, err := certificateClient.Update(context.TODO(), &data, metav1.UpdateOptions{})
// 	if err != nil {
// 		return WorkloadResult(nil, err)
// 	}
// 	return WorkloadResult(nil, nil)
// }

// func DeleteK8sHelmChart(data cmapi.CertificateRequest) K8sWorkloadResult {
// 	kubeProvider := NewKubeProviderCertManager()
// 	certificateClient := kubeProvider.ClientSet.CertmanagerV1().CertificateRequests(data.Namespace)
// 	err := certificateClient.Delete(context.TODO(), data.Name, metav1.DeleteOptions{})
// 	if err != nil {
// 		return WorkloadResult(nil, err)
// 	}
// 	return WorkloadResult(nil, nil)
// }
