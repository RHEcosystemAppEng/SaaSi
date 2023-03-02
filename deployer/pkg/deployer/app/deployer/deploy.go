package deployer

import (
	"log"
	"os"
	"os/exec"

	"github.com/RHEcosystemAppEng/SaaSi/deployer/pkg/deployer/app/packager"
	"github.com/RHEcosystemAppEng/SaaSi/deployer/pkg/utils"
)

var (
	err error
)

func DeployApplication(pkg *packager.ApplicationPkg) {

	// validate oc cli
	utils.ValidateRequirements(utils.OC)

	// deploy application target namespaces using oc cli
	cmd := exec.Command("oc", "apply", "-f", pkg.TargetNamespaceDir, "--kubeconfig", pkg.DeployerContext.GetKubeConfigPath())
	cmd.Dir, _ = os.Getwd()
	if err = cmd.Run(); err != nil {
		log.Fatalf("Failed to deploy files from target namespace directory: %s, Error: %s", pkg.TargetNamespaceDir, err)
	}
	log.Printf("Successfully deployed all files from target namespace directory: %s", pkg.TargetNamespaceDir)

	// deploy application package using oc cli
	cmd = exec.Command("oc", "apply", "-f", ".", "--kubeconfig", pkg.DeployerContext.GetKubeConfigPath())
	cmd.Dir = pkg.DeloymentDir
	if err = cmd.Run(); err != nil {
		log.Fatalf("Failed to deploy files from deployment directory: %s, Error: %s", pkg.DeloymentDir, err)
	}
	log.Printf("Successfully deployed all files from deployment directory: %s", pkg.DeloymentDir)
}
