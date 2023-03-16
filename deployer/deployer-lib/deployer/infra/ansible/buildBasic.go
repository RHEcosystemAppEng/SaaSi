package ansible

import (
	"bufio"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/RHEcosystemAppEng/SaaSi/deployer/deployer-lib/config"
	"github.com/RHEcosystemAppEng/SaaSi/deployer/deployer-lib/context"
	"gopkg.in/yaml.v2"
)

// BuildCustomParameters This Function gets the input customParams , and set them as environment variables , and make an environment file , which each pair prefixed by "export" keyword.
func (playbook Playbook) BuildCustomParameters(customParams config.ClusterParams, pathToBuild string) string {

	fullFilePath := path.Join(pathToBuild, playbook.Name+"-"+"custom.env")
	customEnvFile, err := os.Create(fullFilePath)
	if err != nil {
		log.Fatalf("Error creating file %s, \n details: %s", fullFilePath, err)
		return ""
	}
	defer customEnvFile.Close()
	yamlFile, err := yaml.Marshal(customParams)
	var params map[string]string
	if err != nil {
		log.Fatalf("Error creating yaml out of cluster Params %s, \n details: %s", customParams, err)
		return ""
	} else {
		err := yaml.Unmarshal(yamlFile, &params)
		if err != nil {
			log.Fatalf("Error creating map out of cluster Params yaml %s, \n details: %s", params, err)
			return ""
		}
	}
	result := ""
	writer := bufio.NewWriter(customEnvFile)
	for key, value := range params {
		if value != "" {
			os.Setenv(key, value)
			result = result + "export " + key + "=" + value + "\n"
		}
	}
	_, err = writer.WriteString(result)
	if err != nil {
		log.Fatalf("Error creating env file out of cluster Params %s, \n details: %s", customParams, err)
		return ""
	}
	err = writer.Flush()
	if err != nil {
		return ""
	}

	return fullFilePath

}

func (playbook *Playbook) RenderTemplate(pathToScript string, pathToEnvironmentFile string, pathToCustomEnvFile string, ctx *context.InfraContext) {
	//Render template with Jinja template engine using shell script, according to environment variables that were set.
	log.Print("About to render template to get configuration file to be passed to ansible playbook...")
	command := exec.Command(pathToScript, "-s", pathToEnvironmentFile, "-c", pathToCustomEnvFile)
	command.Dir = ctx.InfraRootDir
	output, err := command.Output()
	if err != nil {
		log.Fatalf("Failed to render template of playbook, error : %s", err)
		return
	} else {
		log.Printf("Successfully Rendered template of playbook, result of invocation : %s", string(output))
	}

	colonIndex := strings.LastIndexAny(string(output), ":")
	playbook.RenderedTemplatePath = strings.Trim(string(output[colonIndex+1:]), " \n")

}

func (playbook Playbook) OverrideParametersWithCustoms(awsCredentials config.AwsSettings) {
	os.Setenv("aws_access_key_id", awsCredentials.AwsAccessKeyId)
	os.Setenv("aws_secret_access_key", awsCredentials.AwsSecretAccessKey)
	os.Setenv("aws_public_domain", awsCredentials.AwsPublicDomain)
	os.Setenv("aws_account_name", awsCredentials.AwsAccountName)
}

// ParseDefaultEnvFile This function parse (key,value) pairs from environment file prefixed with "export", in order to set them as environment variables in this application,
// and not only in script.
func (playbook Playbook) ParseDefaultEnvFile(pathToEnvironmentFile string) {
	readFile, err := os.Open(pathToEnvironmentFile)

	if err != nil {
		log.Fatal(err)
	}
	fileScanner := bufio.NewScanner(readFile)

	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		line := fileScanner.Text()
		if strings.Contains(line, "export") {
			omitExport := line[7:]
			split2Array := strings.Split(omitExport, "=")
			//Only if there is value need to set env
			if split2Array[1] != "" {
				os.Setenv(split2Array[0], split2Array[1])
			}
		}
	}

	readFile.Close()
}
