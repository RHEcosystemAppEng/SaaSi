package context

import (
	"github.com/RHEcosystemAppEng/SaaSi/deployer/deployer-lib/config"
	"github.com/RHEcosystemAppEng/SaaSi/deployer/deployer-lib/connect"
	"k8s.io/client-go/rest"
)

// type Context interface {
// 	RootFolder() string
// }

type DeployerContext struct {
	KubeConnection *connect.KubeConnection
	RootOutputDir  string
	RootSourceDir  string
}

func InitDeployerContext(flagArgs config.FlagArgs, kubeConnection *connect.KubeConnection) *DeployerContext {
	dc := DeployerContext{
		KubeConnection: kubeConnection,
		RootOutputDir:  flagArgs.RootOutputDir,
		RootSourceDir:  flagArgs.RootSourceDir,
	}

	return &dc
}

// func (dc *DeployerContext) RootFolder() string {
// 	log.Fatal("Not implemented: RootFolder()")
// 	return ""
// }

func (dc *DeployerContext) GetKubeConfig() *rest.Config {
	return dc.KubeConnection.KubeConfig
}

func (dc *DeployerContext) GetKubeConfigPath() string {
	return dc.KubeConnection.KubeConfigPath
}

func (dc *DeployerContext) GetRootOutputDir() string {
	return dc.RootOutputDir
}

func (dc *DeployerContext) GetRootSourceDir() string {
	return dc.RootSourceDir
}
