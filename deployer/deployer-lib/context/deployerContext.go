package context

import (
	"github.com/RHEcosystemAppEng/SaaSi/deployer/deployer-lib/config"
	"github.com/RHEcosystemAppEng/SaaSi/deployer/deployer-lib/connect"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/rest"
)

type DeployerContext struct {
	KubeConnection *connect.KubeConnection
	RootOutputDir  string
	RootSourceDir  string
	Logger         *logrus.Logger
}

func InitDeployerContext(args *config.Args, kubeConnection *connect.KubeConnection, logger *logrus.Logger) *DeployerContext {
	dc := DeployerContext{
		KubeConnection: kubeConnection,
		RootOutputDir:  args.RootOutputDir,
		RootSourceDir:  args.RootSourceDir,
		Logger:         logger,
	}

	return &dc
}

// getters

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

func (dc *DeployerContext) GetLogger() *logrus.Logger {
	return dc.Logger
}
