package ansible

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const PlaybookRunnerProg = "ansible-playbook"

func (playbook Playbook) Run() PlayBookResults {

	_, err := exec.LookPath(PlaybookRunnerProg)
	if err != nil {
		log.Fatalf("Failed to find %s binary on the system Path, Error : %s", PlaybookRunnerProg , err)
		return PlayBookResults{}
	}
	// run playbook
	playbookDir := filepath.Dir(playbook.Path)
	playbookInvocation := exec.Command(PlaybookRunnerProg, "--extra-vars=" + "deployment=" + filepath.Base(playbook.RenderedTemplatePath) , playbook.Path)
	////playbookInvocation := exec.Command(PlaybookRunnerProg, "--extra-vars", "deployment=" ,path.Join("@", playbook.RenderedTemplatePath), playbook.Path)
	playbookInvocation.Dir = playbookDir

	output, err := playbookInvocation.Output()

	if err != nil  {
		log.Fatalf("Failed to invoke playbook %s , Detailed Error : %s", playbook.Name, err.Error())
		return PlayBookResults{}
	}
	log.Printf("The output of the playbook run is : \n %s",string(output))

	// calculate apiserver address
	finalClusterName := os.Getenv("CLUSTER_NAME")
	finalClusterBaseDomain := os.Getenv("CLUSTER_BASE_DOMAIN")
	addressParts := []string{"api", finalClusterName, finalClusterBaseDomain}
	apiServerAddress := strings.Join(addressParts, ".")
	transportApiServerAddressPort := "https://" + apiServerAddress + ":6443"

	// get admin adminPassword into variable
	adminPasswordFileLocation := filepath.Join(playbookDir, "build", finalClusterName + "." + finalClusterBaseDomain, "auth","kubeadmin-password")
	// get kubeconfig file path location for authetication for deployer
	kubeConfigOutputLocation := filepath.Join(playbookDir, "build", finalClusterName + "." + finalClusterBaseDomain, "auth","kubeconfig")
	passwordCommand := exec.Command("cat", adminPasswordFileLocation)
	adminPassword, err := passwordCommand.Output()
	if err != nil {
		log.Fatalf("Failed to read adminPassword file , Detailed Error : %s",  err)
		return PlayBookResults{}
	}

	playbookResults := PlayBookResults{
		User:             "kubeadmin",
		Password:         string(adminPassword),
		ApiServer:        transportApiServerAddressPort,
		KubeConfigPath:   kubeConfigOutputLocation,
		AdditionalFields: nil,
	}
	return playbookResults
}
