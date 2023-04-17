package deployer

import (
	"os/exec"

	"github.com/RHEcosystemAppEng/SaaSi/deployer/deployer-lib/deployer/app/packager"
	"github.com/RHEcosystemAppEng/SaaSi/deployer/deployer-lib/utils"
)

var err error

func DeployApplication(pkg *packager.ApplicationPkg) error {

	// validate oc cli
	if err = utils.ValidateRequirements(utils.OC); err != nil {
		pkg.DeployerContext.GetLogger().Errorf("%s command not found", utils.OC)
		return err
	}

	// deploy application target namespaces using oc cli
	cmd := exec.Command("oc", "apply", "-f", ".", "--kubeconfig", pkg.DeployerContext.GetKubeConfigPath())
	cmd.Dir = pkg.TargetNamespaceDir
	if err = cmd.Run(); err != nil {
		pkg.DeployerContext.GetLogger().Errorf("Failed to deploy files from target namespace directory: %s", pkg.TargetNamespaceDir)
		return err
	}
	pkg.DeployerContext.GetLogger().Infof("Successfully deployed all files from target namespace directory: %s", pkg.TargetNamespaceDir)

	// deploy application package using oc cli
	cmd = exec.Command("oc", "apply", "-f", ".", "--kubeconfig", pkg.DeployerContext.GetKubeConfigPath())
	cmd.Dir = pkg.DeloymentDir
	if err = cmd.Run(); err != nil {
		pkg.DeployerContext.GetLogger().Errorf("Failed to deploy files from deployment directory: %s", pkg.DeloymentDir)
		return err
	}
	pkg.DeployerContext.GetLogger().Infof("Successfully deployed all files from deployment directory: %s", pkg.DeloymentDir)

	return nil
}
