package ansible

import "github.com/RHEcosystemAppEng/SaaSi/deployer/deployer-lib/context"

type PlayBookResults struct {
	User             string
	Password         string
	ApiServer        string
	KubeConfigPath   string
	AdditionalFields map[string]string
}

type PlaybookRunner interface {
	Run(*context.InfraContext) PlayBookResults
}
