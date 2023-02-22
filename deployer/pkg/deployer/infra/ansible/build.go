package ansible

import (
	"github.com/RHEcosystemAppEng/SaaSi/deployer/pkg/config"
)

type ParametersBuilder interface {
	BuildCustomParameters(customParams config.ClusterParams, pathToBuild string) string
	RenderTemplate(pathToScript string , pathToEnvironmentFile string, pathToCustomEnvFile string) string
	OverrideParametersWithCustoms(config.ClusterParams, config.ClusterParams) (config.ClusterParams, bool)
}





