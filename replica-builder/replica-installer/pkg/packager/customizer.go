package packager

import (
	"log"
	"os/exec"
	"path/filepath"

	"github.com/RHEcosystemAppEng/SaaSi/replica-builder/replica-installer/pkg/config"
	"github.com/RHEcosystemAppEng/SaaSi/replica-builder/replica-installer/pkg/utils"
)

func (pkg *DeploymentPkg)invokePkgCustomizations(ns config.SourceNamespace) { 

	// validate kustomize cli
	utils.ValidateRequirements()
	
	// define path to namespace template directory
	nsTmplDir := filepath.Join(pkg.KustomizeFolder, ns.Name, TEMPLATE_DIR)

	// set the namespace resource to target namespace
	cmd := exec.Command("kustomize", "edit", "set", "namespace", ns.Target)
	cmd.Dir = nsTmplDir
	if err = cmd.Run(); err != nil {
		log.Fatalf("Failed to set namespace resource in %s template", ns.Name)
	}

	// set a common annotation to uuid
	cmd = exec.Command("kustomize", "edit", "set", "annotation", COMMON_ANNOTATION_KEY + pkg.Uuid.String())
	cmd.Dir = nsTmplDir
	if err = cmd.Run(); err != nil {
		log.Fatalf("Failed to set uuid common annotations in %s template", ns.Name)
	}
}