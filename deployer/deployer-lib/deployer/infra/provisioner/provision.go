package provisioner

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	filepath "path/filepath"
	"strings"

	"github.com/RHEcosystemAppEng/SaaSi/deployer/deployer-lib/config"
	"github.com/RHEcosystemAppEng/SaaSi/deployer/deployer-lib/context"
	"github.com/RHEcosystemAppEng/SaaSi/deployer/deployer-lib/deployer/infra/ansible"
)

func ProvisionCluster(infraCtx *context.InfraContext, customParams *config.ClusterParams, awsCredentials config.AwsSettings) ansible.PlayBookResults {
	playbook := &ansible.Playbook{
		Name:                   "ocp_lab_provisioner",
		Path:                   path.Join(infraCtx.AnsiblePlaybookPath, "site.yaml"),
		OverrideParametersPath: "",
		RenderedTemplatePath:   "",
	}
	// get clusters directory for fetching them assembled default configuration from source cluster.
	inputClusterDirectory := infraCtx.SourceClustersDir
	//find env file in input clusters directory
	envFilePath := findClusterEnvironmentFile(inputClusterDirectory)

	playbook.ParseDefaultEnvFile(envFilePath)

	customParametersPath := playbook.BuildCustomParameters(*customParams, infraCtx.InfraRootDir)
	playbook.OverrideParametersPath = customParametersPath

	playbook.OverrideParametersWithCustoms(awsCredentials)

	//Render template according to environment variables that were set.
	playbook.RenderTemplate(infraCtx.ScriptPath, envFilePath, customParametersPath, infraCtx)
	//Need full path for rendered Template
	playbook.RenderedTemplatePath = filepath.Join(infraCtx.InfraRootDir, playbook.RenderedTemplatePath)
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
