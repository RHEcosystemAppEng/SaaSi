package builder

import (
	"github.com/RHEcosystemAppEng/SaaSi/deployer/pkg/config"
)

type ParametersBuilder interface {
	BuildCustomParameters(customParams config.ClusterParams) string
	RenderTemplate(pathToScript string , pathToEnvironmentFile string, pathToCustomEnvFile string) string
	OverrideParametersWithCustoms(config.ClusterParams, config.ClusterParams) (config.ClusterParams, bool)
}





