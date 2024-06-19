package config

import (
	v1 "github.com/jd-opensource/joylive-injector/client-go/apis/injector/v1"
	"github.com/jd-opensource/joylive-injector/pkg/log"
	"github.com/jd-opensource/joylive-injector/pkg/resource"
	"go.uber.org/zap"
	"os"
)

const (
	AgentConfig           = "config.yaml"
	InjectorConfigName    = "injector.yaml"
	LogConfig             = "logback.xml"
	BootConfig            = "bootstrap.properties"
	ConfigMountPath       = "/joylive/config"
	EmptyDirMountPath     = "/joylive"
	InitEmptyDirMountPath = "/agent"
	InitContainerCmd      = "/bin/sh"
	InitContainerArgs     = "-c, cp -r /joylive/* /agent && chmod -R 777 /agent"
	ConfigMapEnvName      = "JOYLIVE_CONFIGMAP_NAME"
	NamespaceEnvName      = "JOYLIVE_NAMESPACE"
	DefaultNamespace      = "joylive"
	AgentVersionLabel     = "x-live-version"
	LiveSpaceIdLabel      = "x-live-space-id"
	LiveUnitLabel         = "x-live-unit"
	LiveCellLabel         = "x-live-cell"
)

var (
	Cert               string
	Key                string
	Addr               string
	ConfigMountSubPath string
	MatchLabel         string
)

// injection_deploy config
var (
	InitContainerName        string
	DefaultInjectorConfigMap map[string]string
	InjectorConfigMaps       map[string]map[string]string
	InjectorConfig           *AgentInjectorConfig
	InjectorAgentVersion     map[string]v1.AgentVersionSpec
)

func init() {
	// Start the ConfigMap listener and initialize the content
	cmWatcher := NewConfigMapWatcher(resource.GetResource().ClientSet)
	err := cmWatcher.Start()
	if err != nil {
		log.Fatal("start cmWatcher error", zap.Error(err))
	}
	err = cmWatcher.InitConfigMap(GetNamespace())
	if err != nil {
		log.Fatal("init cm error", zap.Error(err))
	}

	// Start the AgentVersion listener and initialize the content
	avWatcher := NewAgentVersionWatcher(resource.GetResource().RestConfig)
	err = avWatcher.Start()
	if err != nil {
		log.Fatal("start avWatcher error", zap.Error(err))
	}
	err = avWatcher.InitAgentVersion(GetNamespace())
	if err != nil {
		log.Fatal("init agentVersion error", zap.Error(err))
	}
}

func GetNamespace() string {
	namespace := os.Getenv(NamespaceEnvName)
	if len(namespace) == 0 {
		namespaceBytes, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
		if err != nil {
			log.Warnf("Failed to read namespace file: %v", err)
		}
		namespace = string(namespaceBytes)
	}
	if len(namespace) == 0 {
		return DefaultNamespace
	}
	return namespace
}
