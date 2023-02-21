package ansible

import (
	"log"
	"os/exec"
	"path"
)

const PlaybookRunnerProg = "ansible-playbook"

func (playbook Playbook) run(pathToPlaybook string, pathToParametersFile string) PlayBookResults {

	_, err := exec.LookPath(PlaybookRunnerProg)
	if err != nil {
		log.Fatalf("Failed to find %s binary on the system Path, Error : %s", PlaybookRunnerProg , err)
		return PlayBookResults{}
	}
	// run playbook
	playbookInvocation := exec.Command(PlaybookRunnerProg, "--extra-vars", path.Join("@", playbook.renderedTemplatePath), playbook.path)
	output, err := playbookInvocation.Output()

	if err != nil {
		log.Fatalf("Failed to invoke playbook %s , Detailed Error : %s", playbook.name , err)
		return PlayBookResults{}
	}
	log.Printf("The output of the playbook run is : \n %s",string(output))


}
