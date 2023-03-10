package provisioner

import (
	"github.com/RHEcosystemAppEng/SaaSi/deployer/pkg/config"
	"github.com/RHEcosystemAppEng/SaaSi/deployer/pkg/context"
	"github.com/RHEcosystemAppEng/SaaSi/deployer/pkg/deployer/infra/ansible"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	filepath "path/filepath"
	"strings"
)

func ProvisionCluster(infraCtx *context.InfraContext, customParams *config.ClusterParams, awsCredentials config.AwsSettings, sourceDirRoot string) ansible.PlayBookResults {
	playbook := &ansible.Playbook{
		Name:                   "ocp_lab_provisioner",
		Path:                   path.Join(infraCtx.AnsiblePlaybookPath,"site.yaml"),
		OverrideParametersPath: "",
		RenderedTemplatePath:   "",
	}
	// get clusters directory for fetching them assembled default configuration from source cluster.
	inputClusterDirectory := filepath.Join(sourceDirRoot, infraCtx.SourceClustersDir)
	//find env file in input clusters directory
	envFilePath := findClusterEnvironmentFile(inputClusterDirectory)
	workingDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current working directory in order to calculate full path of cluster environment file, Error : %s ", err)
		//return ansible.PlayBookResults{}, errors.New("Failed Openning current working directory")
	}
	//calculate full env file path
	fullEnvFilePath := filepath.Join(workingDir,envFilePath)


	playbook.ParseDefaultEnvFile(fullEnvFilePath)

	customParametersPath := playbook.BuildCustomParameters(*customParams, infraCtx.InfraRootDir)
	playbook.OverrideParametersPath = customParametersPath

	playbook.OverrideParametersWithCustoms(awsCredentials)

	//Render template according to environment variables that were set.
	playbook.RenderTemplate(infraCtx.ScriptPath,fullEnvFilePath,customParametersPath,infraCtx)
	//Need full path for rendered Template
	playbook.RenderedTemplatePath = filepath.Join(infraCtx.InfraRootDir,playbook.RenderedTemplatePath)
	// Copy rendered input file to playbook directory and update renderedTemplatePath to this new location
	playbook.RenderedTemplatePath = copyRenderedTemplateToPlaybookDir(playbook)
	return playbook.Run(infraCtx)

}

func copyRenderedTemplateToPlaybookDir(playbook *ansible.Playbook) string {
	src, err := os.Open(playbook.RenderedTemplatePath)
	if err != nil {
		log.Fatalf("Error -  Failed to open rendered template file from path %s, detailed error : \n %s", src, err)
	}
	defer src.Close()
	renderedTemplateDestPath := filepath.Join(filepath.Dir(playbook.Path), "deployment.yaml")
	dst, err := os.Create(renderedTemplateDestPath)
	if err != nil {
		log.Fatalf("Error -  Failed to create destination for Copying rendered template file from path %s to path %s, detailed error : \n %s", src, dst, err)
	}
	defer dst.Close()
	_, err = io.Copy(dst, src)
	if err != nil {
		log.Fatalf("Error - Failed to to Copy rendered template file from path %s to path %s, detailed error : \n %s", src, dst, err)
	}
	return renderedTemplateDestPath
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
