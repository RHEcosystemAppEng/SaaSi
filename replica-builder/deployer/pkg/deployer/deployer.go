package deployer

import (
	"fmt"

	"github.com/RHEcosystemAppEng/SaaSi/replica-builder/deployer/pkg/packager"
)

func NewDeployment(pkg *packager.ApplicationPkg) {
	fmt.Println("deploying application....")
}