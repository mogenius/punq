package utils

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/mogenius/punq/version"

	"github.com/mogenius/punq/logger"

	"github.com/ilyakaznacheev/cleanenv"
)

// This object will initially created in secrets when the software is installed into the cluster for the first time (resource: secret -> mogenius/mogenius)
type ClusterSecret struct {
	ApiKey       string
	ClusterMfaId string
	ClusterName  string
}

type Config struct {
	Browser struct {
		Host string `yaml:"host" env:"browser_host" env-description:"Host of the browser app."`
		Port string `yaml:"port" env:"browser_port" env-description:"Port of the browser app."`
	} `yaml:"browser"`
	Kubernetes struct {
		ClusterName  string `yaml:"cluster_name" env:"cluster_name" env-description:"The Name of the Kubernetes Cluster"`
		OwnNamespace string `yaml:"own_namespace" env:"OWN_NAMESPACE" env-description:"The Namespace of mogenius platform"`
		RunInCluster bool   `yaml:"run_in_cluster" env:"run_in_cluster" env-description:"If set to true, the application will run in the cluster (using the service account token). Otherwise it will try to load your local default context." env-default:"false"`
	} `yaml:"kubernetes"`
	Misc struct {
		Stage            string   `yaml:"stage" env:"stage" env-description:"mogenius k8s-manager stage" env-default:"prod"`
		Debug            bool     `yaml:"debug" env:"debug" env-description:"If set to true, debug features will be enabled." env-default:"false"`
		CheckForUpdates  int      `yaml:"check_for_updates" env:"check_for_updates" env-description:"Time interval between update checks." env-default:"86400"`
		IgnoreNamespaces []string `yaml:"ignore_namespaces" env:"ignore_namespaces" env-description:"List of all ignored namespaces." env-default:""`
	} `yaml:"misc"`
}

var DefaultConfigLocalFile string
var DefaultConfigClusterFileDev string
var DefaultConfigClusterFileProd string
var CONFIG Config

func InitConfigYaml(showDebug bool, customConfigName *string, loadClusterConfig bool) {
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

	if showDebug {
		PrintSettings()
	}

	if CONFIG.Misc.Debug {
		logger.Log.Notice("Starting serice for pprof in localhost:6060")
		go func() {
			logger.Log.Info(http.ListenAndServe("localhost:6060", nil))
			logger.Log.Info("1. Portforward punq to 6060")
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
	logger.Log.Infof("BROWSER")
	logger.Log.Infof("Host:                     %s", CONFIG.Browser.Host)
	logger.Log.Infof("Port:                     %s", CONFIG.Browser.Port)

	logger.Log.Infof("KUBERNETES")
	logger.Log.Infof("ClusterName:              %s", CONFIG.Kubernetes.ClusterName)
	logger.Log.Infof("OwnNamespace:             %s", CONFIG.Kubernetes.OwnNamespace)
	logger.Log.Infof("RunInCluster:             %t", CONFIG.Kubernetes.RunInCluster)

	logger.Log.Infof("MISC")
	logger.Log.Infof("Stage:                    %s", CONFIG.Misc.Stage)
	logger.Log.Infof("Debug:                    %t", CONFIG.Misc.Debug)
	logger.Log.Infof("CheckForUpdates:          %d", CONFIG.Misc.CheckForUpdates)
	logger.Log.Infof("IgnoreNamespaces:         %s", strings.Join(CONFIG.Misc.IgnoreNamespaces, ","))
}

func PrintVersionInfo() {
	fmt.Println("")
	logger.Log.Infof("Version:     %s", version.Ver)
	logger.Log.Infof("Operator:    %s", version.OperatorVersion)
	logger.Log.Infof("Branch:      %s", version.Branch)
	logger.Log.Infof("Hash:        %s", version.GitCommitHash)
	logger.Log.Infof("BuildAt:     %s", version.BuildTimestamp)
}

func GetDirectories(customConfigName *string) (configDir string, configPath string) {
	homeDirName, err := os.UserHomeDir()
	if err != nil {
		logger.Log.Error(err)
	}

	configDir = homeDirName + "/.punq/"
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
