package provisioner

import (
	"github.com/RHEcosystemAppEng/SaaSi/deployer/pkg/config"
	"github.com/RHEcosystemAppEng/SaaSi/deployer/pkg/context"
	"github.com/RHEcosystemAppEng/SaaSi/deployer/pkg/deployer/infra/ansible"
	"io/ioutil"
	"log"
	"path"
	"path/filepath"
	"strings"
)

func ProvisionCluster(infraCtx *context.InfraContext, customParams *config.ClusterParams, sourceDirRoot string) ansible.PlayBookResults {
	playbook := &ansible.Playbook{
		Name:                   "test",
		Path:                   path.Join(infraCtx.AnsiblePlaybookPath,"site.yaml"),
		OverrideParametersPath: "",
		RenderedTemplatePath:   "",
	}
	inputClusterDirectory := filepath.Join(sourceDirRoot, infraCtx.SourceClustersDir)
	//find env file in input clusters directory
	envFilePath := findClusterEnvironmentFile(inputClusterDirectory)
	fullEnvFilePath := filepath.Join(inputClusterDirectory, envFilePath)

	customParametersPath := playbook.BuildCustomParameters(*customParams, infraCtx.ScriptPath)
	playbook.OverrideParametersPath = customParametersPath
	playbook.RenderTemplate(infraCtx.ScriptPath,fullEnvFilePath,customParametersPath)
	return playbook.Run()


}

func findClusterEnvironmentFile(inputClusterDirectory string) string {
	files, err := ioutil.ReadDir(inputClusterDirectory)
	if err != nil {
		log.Fatalf("Error reading the Source Clusters Directory %s, detailed error : \n %s", inputClusterDirectory, err)
		return ""
	}
	var pathToEnvFile string = ""
	var foundEnvFile bool = false
	for _, file := range files {
		if file.IsDir() {
			pathToEnvFile += path.Join(pathToEnvFile, file.Name())
			clusterRootPlusDir := filepath.Join(inputClusterDirectory, pathToEnvFile)
			innerFiles, err := ioutil.ReadDir(clusterRootPlusDir)
			if err != nil {
				log.Fatalf("Error Failed to read env file from path %s, detailed error : \n %s", clusterRootPlusDir, err)
				return ""
			}
			for _, innerFile := range innerFiles {
				if !innerFile.IsDir() && strings.HasSuffix(innerFile.Name(), "env") {
					pathToEnvFile += path.Join(pathToEnvFile, innerFile.Name())
					foundEnvFile = true
					break
				}
			}
		} else {
			if strings.HasSuffix(file.Name(), "env") && !foundEnvFile {
				pathToEnvFile = path.Join(pathToEnvFile, file.Name())
				foundEnvFile = true
			}
		}
	}
	return pathToEnvFile
}
