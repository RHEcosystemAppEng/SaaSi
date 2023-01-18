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

func DeployApplication(pkg *packager.ApplicationPkg) {

	// validate oc cli
	utils.ValidateRequirements(utils.OC)

	// deploy application package using oc cli
	cmd := exec.Command("oc", "apply", "-f", ".")
	cmd.Dir = pkg.DeloymentDir
	if err = cmd.Run(); err != nil {
		log.Fatalf("Failed to deploy files from deployment directory: %s, Error: %s", pkg.DeloymentDir, err)
	}
	log.Printf("Successfully deployed all files from deployment directory: %s", pkg.DeloymentDir)
}
