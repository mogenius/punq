package utils

import (
	"fmt"
	"net/http"
	"os"
	"punq/logger"
	"punq/version"
	"strings"

	"github.com/ilyakaznacheev/cleanenv"
)

// This object will initially created in secrets when the software is installed into the cluster for the first time (resource: secret -> mogenius/mogenius)
type ClusterSecret struct {
	ApiKey       string
	ClusterMfaId string
	ClusterName  string
}

type Config struct {
	Kubernetes struct {
		ApiKey                   string `yaml:"api_key" env:"api_key" env-description:"Api Key to access the server"`
		ClusterName              string `yaml:"cluster_name" env:"cluster_name" env-description:"The Name of the Kubernetes Cluster"`
		OwnNamespace             string `yaml:"own_namespace" env:"OWN_NAMESPACE" env-description:"The Namespace of mogenius platform"`
		ClusterMfaId             string `yaml:"cluster_mfa_id" env:"cluster_mfa_id" env-description:"UUID of the Kubernetes Cluster for MFA purpose"`
		RunInCluster             bool   `yaml:"run_in_cluster" env:"run_in_cluster" env-description:"If set to true, the application will run in the cluster (using the service account token). Otherwise it will try to load your local default context." env-default:"false"`
		DefaultContainerRegistry string `yaml:"default_container_registry" env:"default_container_registry" env-description:"Default Container Image Registry"`
		BboltDbPath              string `yaml:"bbolt_db_path" env:"bbolt_db_path" env-description:"Path to the bbolt database. This db stores build-related information."`
	} `yaml:"kubernetes"`
	ApiServer struct {
		Http_Server string `yaml:"http_server" env:"api_http_server" env-description:"Server host" env-default:"https://platform-api.mogenius.com"`
		Ws_Server   string `yaml:"ws_server" env:"api_ws_server" env-description:"Server host" env-default:"127.0.0.1:8080"`
		Ws_Scheme   string `yaml:"ws_server_scheme" env:"api_ws_scheme" env-description:"Server host scheme. (ws/wss)" env-default:"wss"`
		WS_Path     string `yaml:"ws_path" env:"api_ws_path" env-description:"Server Path" env-default:"/ws"`
	} `yaml:"api_server"`
	EventServer struct {
		Server string `yaml:"server" env:"event_server" env-description:"Server host" env-default:"127.0.0.1:8080"`
		Scheme string `yaml:"scheme" env:"event_scheme" env-description:"Server host scheme. (ws/wss)" env-default:"wss"`
		Path   string `yaml:"path" env:"event_path" env-description:"Server Path" env-default:"/ws-event"`
	} `yaml:"event_server"`
	ShellServer struct {
		Server string `yaml:"server" env:"shell_server" env-description:"Server host" env-default:"127.0.0.1:8080"`
		Scheme string `yaml:"scheme" env:"shell_scheme" env-description:"Server host scheme. (ws/wss)" env-default:"wss"`
		Path   string `yaml:"path" env:"shell_path" env-description:"Server Path" env-default:"/ws-shell"`
	} `yaml:"shell_server"`
	Misc struct {
		Stage                 string   `yaml:"stage" env:"stage" env-description:"mogenius k8s-manager stage" env-default:"prod"`
		Debug                 bool     `yaml:"debug" env:"debug" env-description:"If set to true, debug features will be enabled." env-default:"false"`
		LogKubernetesEvents   bool     `yaml:"log_kubernetes_events" env:"log_kubernetes_events" env-description:"If set to true, all kubernetes events will be logged to std-out." env-default:"false"`
		DefaultMountPath      string   `yaml:"default_mount_path" env:"default_mount_path" env-description:"All containers will have access to this mount point"`
		IgnoreNamespaces      []string `yaml:"ignore_namespaces" env:"ignore_namespaces" env-description:"List of all ignored namespaces." env-default:""`
		AutoMountNfs          bool     `yaml:"auto_mount_nfs" env:"auto_mount_nfs" env-description:"If set to true, nfs pvc will automatically be mounted." env-default:"true"`
		IgnoreResourcesBackup []string `yaml:"ignore_resources_backup" env:"ignore_resources_backup" env-description:"List of all ignored resources while backup." env-default:""`
		CheckForUpdates       int      `yaml:"check_for_updates" env:"check_for_updates" env-description:"Time interval between update checks." env-default:"86400"`
		HelmIndex             string   `yaml:"helm_index" env:"helm_index" env-description:"URL of the helm index file." env-default:"https://helm.mogenius.com/public/index.yaml"`
		ClusterProvider       string   `yaml:"cluster_provider" env:"cluster_provider" env-description:"Providers like AKS, EKS, GCP etc."`
		NfsPodPrefix          string   `yaml:"nfs_pod_prefix" env:"nfs_pod_prefix" env-description:"A prefix for the nfs-server pod. This will always be applied in order to detect the pod."`
	} `yaml:"misc"`
	Builder struct {
		BuildTimeout int `yaml:"max_build_time" env:"max_build_time" env-description:"Seconds until the build will be canceled." env-default:"3600"`
		ScanTimeout  int `yaml:"max_scan_time" env:"max_build_time" env-description:"Seconds until the vulnerability scan will be canceled." env-default:"200"`
	} `yaml:"builder"`
	Git struct {
		GitUserEmail      string `yaml:"git_user_email" env:"git_user_email" env-description:"Email address which is used when interacting with git." env-default:"git@mogenius.com"`
		GitUserName       string `yaml:"git_user_name" env:"git_user_name" env-description:"User name which is used when interacting with git." env-default:"mogenius git-user"`
		GitDefaultBranch  string `yaml:"git_default_branch" env:"git_default_branch" env-description:"Default branch name which is used when creating a repository." env-default:"main"`
		GitAddIgnoredFile string `yaml:"git_add_ignored_file" env:"git_add_ignored_file" env-description:"Gits behaviour when adding ignored files." env-default:"false"`
	} `yaml:"git"`
}

var DefaultConfigLocalFile string
var DefaultConfigClusterFileDev string
var DefaultConfigClusterFileProd string
var CONFIG Config

func InitConfigYaml(showDebug bool, customConfigName *string, clusterSecret ClusterSecret, loadClusterConfig bool) {
	_, configPath := GetDirectories(customConfigName)

	if _, err := os.Stat(configPath); err == nil || os.IsExist(err) {
		// file exists
		if err := cleanenv.ReadConfig(configPath, &CONFIG); err != nil {
			if strings.HasPrefix(err.Error(), "config file parsing error:") {
				logger.Log.Notice("Config file is corrupted. Creating a new one by using -r flag.")
			}
			logger.Log.Fatal(err)
		}
	} else {
		WriteDefaultConfig(loadClusterConfig)

		// read configuration from the file and environment variables
		if err := cleanenv.ReadConfig(configPath, &CONFIG); err != nil {
			logger.Log.Fatal(err)
		}
	}

	// LOCAL ALWAYS WINS
	if clusterSecret.ClusterMfaId != "" && CONFIG.Kubernetes.RunInCluster {
		CONFIG.Kubernetes.ClusterMfaId = clusterSecret.ClusterMfaId
	}
	if clusterSecret.ApiKey != "" && CONFIG.Kubernetes.RunInCluster {
		CONFIG.Kubernetes.ApiKey = clusterSecret.ApiKey
	}
	if clusterSecret.ClusterName != "" && CONFIG.Kubernetes.RunInCluster {
		CONFIG.Kubernetes.ClusterName = clusterSecret.ClusterName
	}

	if showDebug {
		PrintSettings()
	}

	// CHECKS FOR CLUSTER
	if loadClusterConfig {
		if CONFIG.Kubernetes.ClusterName == "your-cluster-name" || CONFIG.Kubernetes.ClusterName == "" {
			if !showDebug {
				PrintSettings()
			}
			logger.Log.Fatalf("Environment Variable 'cluster_name' not setup. TERMINATING.")
		}
		if CONFIG.Kubernetes.ApiKey == "YOUR_API_KEY" || CONFIG.Kubernetes.ApiKey == "" {
			if !showDebug {
				PrintSettings()
			}
			logger.Log.Fatalf("Environment Variable 'api_key' not setup or default value not overwritten. TERMINATING.")
		}
	}

	if CONFIG.Misc.Debug {
		logger.Log.Notice("Starting serice for pprof in localhost:6060")
		go func() {
			logger.Log.Info(http.ListenAndServe("localhost:6060", nil))
			logger.Log.Info("1. Portforward mogenius-k8s-manager to 6060")
			logger.Log.Info("2. wget http://localhost:6060/debug/pprof/profile?seconds=60 -O cpu.pprof")
			logger.Log.Info("3. wget http://localhost:6060/debug/pprof/heap -O mem.pprof")
			logger.Log.Info("4. go tool pprof -http=localhost:8081 cpu.pprof")
			logger.Log.Info("5. go tool pprof -http=localhost:8081 mem.pprof")
			logger.Log.Info("OR: go tool pprof mem.pprof -> Then type in commands like top, top --cum, list")
			logger.Log.Info("http://localhost:6060/debug/pprof/ This is the index page that lists all available profiles.")
			logger.Log.Info("http://localhost:6060/debug/pprof/profile This serves a CPU profile. You can set the profiling duration through the seconds parameter. For example, ?seconds=30 would profile your CPU for 30 seconds.")
			logger.Log.Info("http://localhost:6060/debug/pprof/heap This serves a snapshot of the current heap memory usage.")
			logger.Log.Info("http://localhost:6060/debug/pprof/goroutine This serves a snapshot of the current goroutines stack traces.")
			logger.Log.Info("http://localhost:6060/debug/pprof/block This serves a snapshot of stack traces that led to blocking on synchronization primitives.")
			logger.Log.Info("http://localhost:6060/debug/pprof/threadcreate This serves a snapshot of all OS thread creation stack traces.")
			logger.Log.Info("http://localhost:6060/debug/pprof/cmdline This returns the command line invocation of the current program.")
			logger.Log.Info("http://localhost:6060/debug/pprof/symbol This is used to look up the program counters listed in a pprof profile.")
			logger.Log.Info("http://localhost:6060/debug/pprof/trace This serves a trace of execution of the current program. You can set the trace duration through the seconds parameter.")
		}()
	}
}

func PrintSettings() {
	logger.Log.Infof("KUBERNETES")
	logger.Log.Infof("OwnNamespace:             %s", CONFIG.Kubernetes.OwnNamespace)
	logger.Log.Infof("ClusterName:              %s", CONFIG.Kubernetes.ClusterName)
	logger.Log.Infof("ClusterMfaId:             %s", CONFIG.Kubernetes.ClusterMfaId)
	logger.Log.Infof("RunInCluster:             %t", CONFIG.Kubernetes.RunInCluster)
	logger.Log.Infof("DefaultContainerRegistry: %s", CONFIG.Kubernetes.DefaultContainerRegistry)
	logger.Log.Infof("ApiKey:                   %s", CONFIG.Kubernetes.ApiKey)
	logger.Log.Infof("BboltDbPath:              %s\n\n", CONFIG.Kubernetes.BboltDbPath)

	logger.Log.Infof("API")
	logger.Log.Infof("HttpServer:               %s", CONFIG.ApiServer.Http_Server)
	logger.Log.Infof("WsServer:                 %s", CONFIG.ApiServer.Ws_Server)
	logger.Log.Infof("WsScheme:                 %s", CONFIG.ApiServer.Ws_Scheme)
	logger.Log.Infof("WsPath:                   %s", CONFIG.ApiServer.WS_Path)

	logger.Log.Infof("EVENTS")
	logger.Log.Infof("EventServer:              %s", CONFIG.EventServer.Server)
	logger.Log.Infof("EventScheme:              %s", CONFIG.EventServer.Scheme)
	logger.Log.Infof("EventPath:                %s\n\n", CONFIG.EventServer.Path)

	logger.Log.Infof("SHELL")
	logger.Log.Infof("ShellServer:              %s", CONFIG.ShellServer.Server)
	logger.Log.Infof("ShellScheme:              %s", CONFIG.ShellServer.Scheme)
	logger.Log.Infof("ShellPath:                %s\n\n", CONFIG.ShellServer.Path)

	logger.Log.Infof("MISC")
	logger.Log.Infof("Stage:                    %s", CONFIG.Misc.Stage)
	logger.Log.Infof("Debug:                    %t", CONFIG.Misc.Debug)
	logger.Log.Infof("AutoMountNfs:             %t", CONFIG.Misc.AutoMountNfs)
	logger.Log.Infof("LogKubernetesEvents:      %t", CONFIG.Misc.LogKubernetesEvents)
	logger.Log.Infof("DefaultMountPath:         %s", CONFIG.Misc.DefaultMountPath)
	logger.Log.Infof("IgnoreResourcesBackup:    %s", strings.Join(CONFIG.Misc.IgnoreResourcesBackup, ","))
	logger.Log.Infof("IgnoreNamespaces:         %s", strings.Join(CONFIG.Misc.IgnoreNamespaces, ","))
	logger.Log.Infof("CheckForUpdates:          %d", CONFIG.Misc.CheckForUpdates)
	logger.Log.Infof("HelmIndex:                %s", CONFIG.Misc.HelmIndex)
	logger.Log.Infof("ClusterProvider:          %s", CONFIG.Misc.ClusterProvider)
	logger.Log.Infof("NfsPodPrefix:             %s", CONFIG.Misc.NfsPodPrefix)
	logger.Log.Infof("ClusterProvider:          %d\n\n", CONFIG.Builder.BuildTimeout)

	logger.Log.Infof("GIT")
	logger.Log.Infof("GitUserEmail:             %s", CONFIG.Git.GitUserEmail)
	logger.Log.Infof("GitUserName:              %s", CONFIG.Git.GitUserName)
	logger.Log.Infof("GitDefaultBranch:         %s", CONFIG.Git.GitDefaultBranch)
	logger.Log.Infof("GitAddIgnoredFile:        %s\n\n", CONFIG.Git.GitAddIgnoredFile)
}

func PrintVersionInfo() {
	fmt.Println("")
	logger.Log.Infof("Version:     %s", version.Ver)
	logger.Log.Infof("Branch:      %s", version.Branch)
	logger.Log.Infof("Hash:        %s", version.GitCommitHash)
	logger.Log.Infof("BuildAt:     %s", version.BuildTimestamp)
}

func GetDirectories(customConfigName *string) (configDir string, configPath string) {
	homeDirName, err := os.UserHomeDir()
	if err != nil {
		logger.Log.Error(err)
	}

	configDir = homeDirName + "/.mogenius-k8s-manager/"
	if customConfigName != nil {
		newConfigName := *customConfigName
		if newConfigName != "" {
			configPath = configDir + newConfigName
		}
	} else {
		configPath = configDir + "config.yaml"
	}

	return configDir, configPath
}

func WriteDefaultConfig(loadClusterConfig bool) {
	configDir, configPath := GetDirectories(nil)

	// write it to default location
	err := os.Mkdir(configDir, 0755)
	if err != nil && err.Error() != "mkdir "+configDir+": file exists" {
		logger.Log.Warning("Error creating folder " + configDir)
		logger.Log.Warning(err)
	}

	stage := os.Getenv("STAGE")

	if loadClusterConfig {
		if stage == "prod" {
			err = os.WriteFile(configPath, []byte(DefaultConfigClusterFileProd), 0755)
		} else {
			err = os.WriteFile(configPath, []byte(DefaultConfigClusterFileDev), 0755)
		}
	} else {
		err = os.WriteFile(configPath, []byte(DefaultConfigLocalFile), 0755)
	}
	if err != nil {
		logger.Log.Error("Error writing " + configPath + " file")
		logger.Log.Fatal(err.Error())
	}
}
