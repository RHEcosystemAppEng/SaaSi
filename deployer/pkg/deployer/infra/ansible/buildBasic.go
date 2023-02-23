package ansible

import (
	"bufio"
	"github.com/RHEcosystemAppEng/SaaSi/deployer/pkg/config"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
)



func (playbook Playbook) BuildCustomParameters(customParams config.ClusterParams,pathToBuild string) string {

	fullFilePath := path.Join(pathToBuild, playbook.Name + "-" + "custom.env")
	customEnvFile, err := os.Create(fullFilePath)
	if err != nil {
		log.Fatalf("Error creating file %s, \n details: %s", fullFilePath,err)
		return ""
	}
	defer customEnvFile.Close()
	yamlFile, err := yaml.Marshal(customParams)
	var params map[string] string
	if err != nil {
		log.Fatalf("Error creating yaml out of cluster Params %s, \n details: %s", customParams,err)
		return ""
	} else {
		err := yaml.Unmarshal(yamlFile, params)
		if err != nil {
			log.Fatalf("Error creating map out of cluster Params yaml %s, \n details: %s", params,err)
			return ""
		}
	}
	result := ""
	writer := bufio.NewWriter(customEnvFile)
	for key, value := range params {
        result = result + "export " + key + "=" + value + "\n"
	}
	_, err = writer.WriteString(result)
	if err != nil {
		log.Fatalf("Error creating env file out of cluster Params %s, \n details: %s", customParams,err)
		return ""
	}
	err = writer.Flush()
	if err != nil {
		return ""
	}

	return fullFilePath

}

func (playbook Playbook ) RenderTemplate(pathToScript string, pathToEnvironmentFile string, pathToCustomEnvFile string) {
	command := exec.Command(pathToScript, "-s", pathToEnvironmentFile, "-c", pathToCustomEnvFile)
	output, err := command.Output()
	if err != nil {
		log.Fatalf("Failed to render template of playbook, error : %s",err)
		return
	} else{
		log.Printf("Successfully Rendered template of playbook, result of invocation : %s",string(output))
	}

	colonIndex := strings.LastIndexAny(string(output), ":")
	playbook.RenderedTemplatePath = strings.Trim(string(output)[colonIndex+1:]," ")

}

func (playbook Playbook) OverrideParametersWithCustoms(awsCredentials config.AwsSettings) () {
   os.Setenv("aws_access_key_id",awsCredentials.AwsAccessKeyId)
   os.Setenv("aws_secret_access_key",awsCredentials.AwsSecretAccessKey)
   os.Setenv("aws_public_domain",awsCredentials.AwsPublicDomain)
   os.Setenv("aws_account_name",awsCredentials.AwsAccountName)
}
