package deployer

import (
	"log"
	"os/exec"

	"github.com/RHEcosystemAppEng/SaaSi/replica-builder/deployer/pkg/packager"
	"github.com/RHEcosystemAppEng/SaaSi/replica-builder/deployer/pkg/utils"
)

var (
	err error
)

func NewDeployment(pkg *packager.ApplicationPkg) {
	
	// validate kustomize cli
	utils.ValidateRequirements()

	// set the namespace resource to target namespace
	cmd := exec.Command("oc", "apply", "-f", ".")
	cmd.Dir = pkg.DeloymentDir
	if err = cmd.Run(); err != nil {
		log.Fatalf("Failed deploy files from deployment directory: %s", pkg.DeloymentDir)
	}
	log.Printf("Successfully deployed all files from deployment directory: %s", pkg.DeloymentDir)
}