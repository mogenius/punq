package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/charmbracelet/glamour"
	"golang.org/x/term"
	"k8s.io/client-go/util/homedir"

	"github.com/mogenius/punq/version"

	"github.com/mogenius/punq/logger"

	"github.com/ilyakaznacheev/cleanenv"
)

const CONFIGVERSION = 2
const USERSSECRET = "punq-users"
const JWTSECRET = "punq-jwt"
const USERADMIN = "admin"
const CONTEXTSSECRET = "punq-contexts"
const CONTEXTOWN = "own-context"

// This object will initially created in secrets when the software is installed into the cluster for the first time (resource: secret -> mogenius/mogenius)
type ClusterSecret struct {
	ApiKey       string
	ClusterMfaId string
	ClusterName  string
}

type Config struct {
	Config struct {
		Version int `yaml:"version" env-description:"Version of the configuration yaml."`
	} `yaml:"config"`
	Frontend struct {
		Host string `yaml:"host" env:"frontend_host" env-description:"Host of the frontend server."`
		Port int    `yaml:"port" env:"frontend_port" env-description:"Port of the frontend server."`
	} `yaml:"frontend"`
	Backend struct {
		Host string `yaml:"host" env:"backend_host" env-description:"Host of the backend server."`
		Port int    `yaml:"port" env:"backend_port" env-description:"Port of the backend server."`
	} `yaml:"backend"`
	Websocket struct {
		Host string `yaml:"host" env:"websocket_host" env-description:"Host of the websocket server."`
		Port int    `yaml:"port" env:"websocket_port" env-description:"Port of the websocket server."`
	} `yaml:"websocket"`
	Kubernetes struct {
		ClusterName  string `yaml:"cluster_name" env:"cluster_name" env-description:"The Name of the Kubernetes Cluster"`
		OwnNamespace string `yaml:"own_namespace" env:"OWN_NAMESPACE" env-description:"The Namespace of mogenius platform"`
		RunInCluster bool   `yaml:"run_in_cluster" env:"run_in_cluster" env-description:"If set to true, the application will run in the cluster (using the service account token). Otherwise it will try to load your local default context." env-default:"false"`
	} `yaml:"kubernetes"`
	Misc struct {
		Stage              string   `yaml:"stage" env:"stage" env-description:"Stage to run in" env-default:"prod"`
		Debug              bool     `yaml:"debug" env:"debug" env-description:"If set to true, debug features will be enabled." env-default:"false"`
		CheckForUpdates    int      `yaml:"check_for_updates" env:"check_for_updates" env-description:"Time interval between update checks." env-default:"86400"`
		ProxyPort          int      `yaml:"proxy_port" env:"proxy_port" env-description:"Default port for proxy releated stuff." env-default:"8888"`
		IgnoreNamespaces   []string `yaml:"ignore_namespaces" env:"ignore_namespaces" env-description:"List of all ignored namespaces." env-default:""`
		ForbidCountryCheck bool     `yaml:"forbid_country_check" env:"forbid_country_check" env-description:"Check clusters location" env-default:"false"`
	} `yaml:"misc"`
}

var DefaultConfigLocalFile string
var DefaultConfigFileOperator string
var DefaultConfigFileProd string
var ChangeLog string
var WelcomeMessage string
var CONFIG Config
var ConfigPath string

func InitConfigYaml(showDebug bool, customConfigName string, stage string) {
	_, ConfigPath = GetDirectories(customConfigName)

	// create default config if not exists
	// if stage is set, then we overwrite the config
	if stage == "" {
		if _, err := os.Stat(ConfigPath); err == nil || os.IsExist(err) {
			// do nothing, file exists
		} else {
			WriteDefaultConfig(stage)
		}
	} else {
		WriteDefaultConfig(stage)
	}

	// read configuration from the file and environment variables
	if err := cleanenv.ReadConfig(ConfigPath, &CONFIG); err != nil {
		if strings.HasPrefix(err.Error(), "config file parsing error:") {
			logger.Log.Notice("Config file is corrupted. Creating a new one by using -r flag.")
		}
		logger.Log.Fatal(err)
	}

	if CONFIG.Kubernetes.RunInCluster {
		ConfigPath = "RUNS_IN_CLUSTER_NO_CONFIG_NEEDED"
	}

	if showDebug {
		PrintSettings()
	}

	if CONFIGVERSION > CONFIG.Config.Version {
		FatalError(fmt.Sprintf("Config version is outdated. Please delete your config file 'punq system reset-config' and restart the application. (Your Config version: %d, Needed: %d)", CONFIG.Config.Version, CONFIGVERSION))
	}

	// if CONFIG.Misc.Debug {
	// 	logger.Log.Notice("Starting serice for pprof in localhost:6060")
	// 	go func() {
	// 		logger.Log.Info(http.ListenAndServe("localhost:6060", nil))
	// 		logger.Log.Info("1. Portforward punq to 6060")
	// 		logger.Log.Info("2. wget http://localhost:6060/debug/pprof/profile?seconds=60 -O cpu.pprof")
	// 		logger.Log.Info("3. wget http://localhost:6060/debug/pprof/heap -O mem.pprof")
	// 		logger.Log.Info("4. go tool pprof -http=localhost:8081 cpu.pprof")
	// 		logger.Log.Info("5. go tool pprof -http=localhost:8081 mem.pprof")
	// 		logger.Log.Info("OR: go tool pprof mem.pprof -> Then type in commands like top, top --cum, list")
	// 		logger.Log.Info("http://localhost:6060/debug/pprof/ This is the index page that lists all available profiles.")
	// 		logger.Log.Info("http://localhost:6060/debug/pprof/profile This serves a CPU profile. You can set the profiling duration through the seconds parameter. For example, ?seconds=30 would profile your CPU for 30 seconds.")
	// 		logger.Log.Info("http://localhost:6060/debug/pprof/heap This serves a snapshot of the current heap memory usage.")
	// 		logger.Log.Info("http://localhost:6060/debug/pprof/goroutine This serves a snapshot of the current goroutines stack traces.")
	// 		logger.Log.Info("http://localhost:6060/debug/pprof/block This serves a snapshot of stack traces that led to blocking on synchronization primitives.")
	// 		logger.Log.Info("http://localhost:6060/debug/pprof/threadcreate This serves a snapshot of all OS thread creation stack traces.")
	// 		logger.Log.Info("http://localhost:6060/debug/pprof/cmdline This returns the command line invocation of the current program.")
	// 		logger.Log.Info("http://localhost:6060/debug/pprof/symbol This is used to look up the program counters listed in a pprof profile.")
	// 		logger.Log.Info("http://localhost:6060/debug/pprof/trace This serves a trace of execution of the current program. You can set the trace duration through the seconds parameter.")
	// 	}()
	// }
}

func PrintSettings() {
	fmt.Printf("Config\n")
	fmt.Printf("Version:                  %d\n", CONFIG.Config.Version)

	fmt.Printf("\nFrontend\n")
	fmt.Printf("Host:                     %s\n", CONFIG.Frontend.Host)
	fmt.Printf("Port:                     %d\n", CONFIG.Frontend.Port)

	fmt.Printf("\nBackend\n")
	fmt.Printf("Host:                     %s\n", CONFIG.Backend.Host)
	fmt.Printf("Port:                     %d\n", CONFIG.Backend.Port)

	fmt.Printf("\nKUBERNETES\n")
	fmt.Printf("ClusterName:              %s\n", CONFIG.Kubernetes.ClusterName)
	fmt.Printf("OwnNamespace:             %s\n", CONFIG.Kubernetes.OwnNamespace)
	fmt.Printf("RunInCluster:             %t\n", CONFIG.Kubernetes.RunInCluster)

	fmt.Printf("\nMISC\n")
	fmt.Printf("Stage:                    %s\n", CONFIG.Misc.Stage)
	fmt.Printf("Debug:                    %t\n", CONFIG.Misc.Debug)
	fmt.Printf("CheckForUpdates:          %d\n", CONFIG.Misc.CheckForUpdates)
	fmt.Printf("ProxyPort:                %d\n", CONFIG.Misc.ProxyPort)
	fmt.Printf("IgnoreNamespaces:         %s\n", strings.Join(CONFIG.Misc.IgnoreNamespaces, ","))
	fmt.Printf("CountryCheck:             %t\n\n", CONFIG.Misc.ForbidCountryCheck)

	fmt.Printf("Config:                   %s\n\n", ConfigPath)
}

func PrintVersionInfo() {
	fmt.Println("")
	logger.Log.Infof("Version:     %s", version.Ver)
	logger.Log.Infof("Operator:    %s", version.OperatorImage)
	logger.Log.Infof("Branch:      %s", version.Branch)
	logger.Log.Infof("Hash:        %s", version.GitCommitHash)
	logger.Log.Infof("BuildAt:     %s", version.BuildTimestamp)
}

func PrintChangeLog() {
	r, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(getTerminalSize()),
	)
	out, _ := r.Render(ChangeLog)
	fmt.Println(out)
}

func PrintWelcomeMessage() {
	r, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(getTerminalSize()),
	)
	out, _ := r.Render(WelcomeMessage)
	fmt.Println(out)
}

func getTerminalSize() int {
	fd := 0
	if runtime.GOOS == "windows" {
		fd = int(os.Stdout.Fd())
	}
	width, _, err := term.GetSize(fd)
	if err != nil {
		logger.Log.Errorf("Failed getting terminal size: %v", err)
	}
	return width
}

func GetDirectories(customConfigPath string) (configDir string, configPath string) {
	homeDirName, err := os.UserHomeDir()
	if err != nil {
		logger.Log.Error(err)
	}

	if customConfigPath != "" {
		if _, err := os.Stat(configPath); err == nil || os.IsExist(err) {
			configPath = customConfigPath
			configDir = filepath.Dir(customConfigPath)
		} else {
			logger.Log.Errorf("Custom config not found '%s'.", customConfigPath)
		}
	} else {
		configDir = homeDirName + "/.punq/"
		configPath = configDir + "config.yaml"
	}

	return configDir, configPath
}

func DeleteCurrentConfig() {
	_, configPath := GetDirectories("")
	err := os.Remove(configPath)
	if err != nil {
		PrintError(fmt.Sprintf("Error removing config file. '%s'.", err.Error()))
	} else {
		PrintInfo(fmt.Sprintf("%s succesfully deleted.", configPath))
		os.Exit(0)
	}
}

func WriteDefaultConfig(stage string) {
	configDir, configPath := GetDirectories("")

	// write it to default location
	err := os.Mkdir(configDir, 0755)
	if err != nil && err.Error() != "mkdir "+configDir+": file exists" {
		logger.Log.Warning("Error creating folder " + configDir)
		logger.Log.Warning(err)
	}

	// check if stage is set via env variable
	envVarStage := strings.ToLower(os.Getenv("stage"))
	if envVarStage != "" {
		stage = envVarStage
	} else {
		// default stage is prod
		if stage == "" {
			stage = "prod"
		}
	}

	if stage == "operator" {
		err = os.WriteFile(configPath, []byte(DefaultConfigFileOperator), 0755)
	} else if stage == "prod" {
		err = os.WriteFile(configPath, []byte(DefaultConfigFileProd), 0755)
	} else if stage == "local" {
		err = os.WriteFile(configPath, []byte(DefaultConfigLocalFile), 0755)
	} else {
		fmt.Println("No stage set. Using local config.")
		err = os.WriteFile(configPath, []byte(DefaultConfigLocalFile), 0755)
	}
	if err != nil {
		logger.Log.Error("Error writing " + configPath + " file")
		logger.Log.Fatal(err.Error())
	}
}

func GetDefaultKubeConfig() string {
	var kubeconfig string = os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		if home := homedir.HomeDir(); home != "" {
			kubeconfig = filepath.Join(home, ".kube", "config")
		}
	}
	// check if file exists in kubeconfig
	if _, err := os.Stat(kubeconfig); os.IsNotExist(err) {
		logger.Log.Fatalf("$KUBECONFIG is empty and default context cannot be loaded. Please set $KUBECONFIG or use --context flag to proceed.")
	}
	return kubeconfig
}
