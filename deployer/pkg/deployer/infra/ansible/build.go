package ansible

import (
	"github.com/RHEcosystemAppEng/SaaSi/deployer/pkg/config"
	"github.com/RHEcosystemAppEng/SaaSi/deployer/pkg/context"
)

type ParametersBuilder interface {
	BuildCustomParameters(customParams config.ClusterParams, pathToBuild string) string
	RenderTemplate(pathToScript string , pathToEnvironmentFile string, pathToCustomEnvFile string, ctx *context.InfraContext) string
	OverrideParametersWithCustoms(awsCredentials config.AwsSettings)
	ParseDefaultEnvFile(pathToEnvironmentFile string)
}





