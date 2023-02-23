package ansible

import (
	"log"
	"os"
	"os/exec"
	"path"
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
	playbookInvocation := exec.Command(PlaybookRunnerProg, "--extra-vars", path.Join("@", playbook.RenderedTemplatePath), playbook.Path)
	output, err := playbookInvocation.Output()

	if err != nil {
		log.Fatalf("Failed to invoke playbook %s , Detailed Error : %s", playbook.Name, err)
		return PlayBookResults{}
	}
	log.Printf("The output of the playbook run is : \n %s",string(output))
	// get admin adminPassword into variable
	// calculate apiserver address
	finalClusterName := os.Getenv("CLUSTER_NAME")
	finalClusterBaseDomain := os.Getenv("CLUSTER_BASE_DOMAIN")
	addressParts := []string{"api", finalClusterName, finalClusterBaseDomain}
	apiServerAddress := strings.Join(addressParts, ".")
	transportApiServerAddressPort := "https://" + apiServerAddress + ":6443"
	adminPasswordFileLocation := filepath.Join(playbook.Path, "build", finalClusterName, finalClusterBaseDomain, "auth","kubeadmin-adminPassword")
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
		AdditionalFields: nil,
	}
	return playbookResults
}
