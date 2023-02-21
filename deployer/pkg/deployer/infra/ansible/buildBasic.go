package ansible

import (
	"github.com/RHEcosystemAppEng/SaaSi/deployer/pkg/config"
	"os/exec"
)



func (playbook Playbook) BuildCustomParameters(customParams config.ClusterParams) string {

	//exec.Command

}

func (playbook Playbook ) RenderTemplate(pathToScript string, pathToEnvironmentFile string, pathToCustomEnvFile string) {


}

func (playbook Ansible.Playbook) OverrideParametersWithCustoms(params config.ClusterParams, params2 config.ClusterParams) (config.ClusterParams, bool) {

   return config.ClusterParams{},true
}
