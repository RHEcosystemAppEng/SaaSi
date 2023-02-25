package provisioner

import (
	"github.com/RHEcosystemAppEng/SaaSi/deployer/pkg/config"
	"github.com/RHEcosystemAppEng/SaaSi/deployer/pkg/context"
	"github.com/RHEcosystemAppEng/SaaSi/deployer/pkg/deployer/infra/ansible"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func ProvisionCluster(infraCtx *context.InfraContext, customParams *config.ClusterParams, awsCredentials config.AwsSettings, sourceDirRoot string) ansible.PlayBookResults {
	playbook := &ansible.Playbook{
		Name:                   "test",
		Path:                   path.Join(infraCtx.AnsiblePlaybookPath,"site.yaml"),
		OverrideParametersPath: "",
		RenderedTemplatePath:   "",
	}
	inputClusterDirectory := filepath.Join(sourceDirRoot, infraCtx.SourceClustersDir)
	//find env file in input clusters directory
	envFilePath := findClusterEnvironmentFile(inputClusterDirectory)
	workingDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current working directory in order to calculate full path of cluster environment file, Error : %s ", err)
		return ansible.PlayBookResults{}
	}
	fullEnvFilePath := filepath.Join(workingDir,envFilePath)
	playbook.ParseDefaultEnvFile(fullEnvFilePath)
	customParametersPath := playbook.BuildCustomParameters(*customParams, infraCtx.InfraRootDir)
	playbook.OverrideParametersPath = customParametersPath
	playbook.OverrideParametersWithCustoms(awsCredentials)
	playbook.RenderTemplate(infraCtx.ScriptPath,fullEnvFilePath,customParametersPath)
	//Need full path for rendered Template
	playbook.RenderedTemplatePath = filepath.Join(infraCtx.InfraRootDir,playbook.RenderedTemplatePath)
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
			clusterRootPlusDir := filepath.Join(inputClusterDirectory, file.Name())
			innerFiles, err := ioutil.ReadDir(clusterRootPlusDir)
			if err != nil {
				log.Fatalf("Error Failed to read env file from path %s, detailed error : \n %s", clusterRootPlusDir, err)
				return ""
			}
			for _, innerFile := range innerFiles {
				if !innerFile.IsDir() && strings.HasSuffix(innerFile.Name(), "env") {
					pathToEnvFile = path.Join(clusterRootPlusDir, innerFile.Name())
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
