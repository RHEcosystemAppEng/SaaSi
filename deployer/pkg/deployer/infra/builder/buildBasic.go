package builder

import (
	"github.com/RHEcosystemAppEng/SaaSi/deployer/pkg/config"
	"github.com/RHEcosystemAppEng/SaaSi/deployer/pkg/deployer/infra/ansible"
)



func (ansible.Playbook) BuildCustomParameters(customParams config.ClusterParams) string {
	panic("implement me")
}

func (ansible.Playbook) RenderTemplate(pathToScript string, pathToEnvironmentFile string, pathToCustomEnvFile string) string {
	panic("implement me")
}

func (ansible.Playbook) OverrideParametersWithCustoms(params config.ClusterParams, params2 config.ClusterParams) (config.ClusterParams, bool) {
	panic("implement me")
}
