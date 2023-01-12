package packager

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/RHEcosystemAppEng/SaaSi/replica-builder/replica-installer/pkg/config"
	"github.com/RHEcosystemAppEng/SaaSi/replica-builder/replica-installer/pkg/utils"
)

const (
	OUTPUT_DIR = "output"
	KUSTOMIZE_DIR = "kustomize"
	DEPLOYMENT_DIR = "deploy"
	TEMPLATE_DIR = "template"

	COMMON_ANNOTATION_KEY = "app.kubernetes.io/saasi-pkg-uuid:"

	SOURCE_KUSTOMIZE_DIR = "../install-builder/output/Infinity/installer/kustomize"
)

var (
	err error
)

type DeploymentPkg struct {
	Uuid uuid.UUID
	AppFolder    string
	KustomizeFolder string
	DeloymentFolder string
}

func NewDeploymentPkg(appConfig *config.ApplicationConfig) *DeploymentPkg {
	
	// init DeploymentPkg
	pkg := DeploymentPkg{}

	// generate uuid
	pkg.Uuid = uuid.New()

	// create application directories
	// application directory
	pwd, _ := os.Getwd()
	log.Printf("Running from %v", pwd)
	pkg.AppFolder = filepath.Join(pwd, OUTPUT_DIR, appConfig.Application.Name)
	utils.CreateDir(pkg.AppFolder)
	// kustomize directory for namesapce artifacts
	pkg.KustomizeFolder = filepath.Join(pkg.AppFolder, KUSTOMIZE_DIR)
	utils.CreateDir(pkg.KustomizeFolder)
	// deployment directory for deployment packages
	pkg.DeloymentFolder = filepath.Join(pkg.AppFolder, DEPLOYMENT_DIR)
	utils.CreateDir(pkg.DeloymentFolder)

	// generate a kustomize-able artifact for each requested namespaces, invoke cutomizations from application config file and build package
	pkg.generatePkg(appConfig)

	// _ = pkg.buildPkgs()

	return &pkg
}

func (pkg *DeploymentPkg) generatePkg(appConfig *config.ApplicationConfig) {
	for _, ns := range appConfig.Application.Namespaces {

		// generate artifact for current namespace
		pkg.generatePkgArtifact(ns)
		
		// invoke customizations onto artifact
		pkg.invokePkgCustomizations(ns)

		// pkg artifact
		pkg.buildPkg(ns)
	}
}

func (pkg *DeploymentPkg) generatePkgArtifact(ns config.SourceNamespace) {
	
	source := filepath.Join(SOURCE_KUSTOMIZE_DIR, ns.Name)
	// create pkg template at pkg template path
	cmd := exec.Command("cp", "-r", source, pkg.KustomizeFolder)
	if err := cmd.Run(); err != nil {
		log.Fatalf("Failed to create pkg kustomize template for namespace: %s", ns.Name)
	}
}

func (pkg *DeploymentPkg) buildPkg(ns config.SourceNamespace) {

	// define path to namespace template directory
	nsTmplDir := filepath.Join(pkg.KustomizeFolder, ns.Name, TEMPLATE_DIR)

	// build pkg with kustomize configuration
	cmd := exec.Command("kustomize", "build", nsTmplDir)
	
	// define path to namespace template directory
	deploymentFile := filepath.Join(pkg.DeloymentFolder, ns.Name + ".yaml")

	// create pkg deployment file
	pkgDeploymentFile, err := os.Create(deploymentFile)
	if err != nil {
		log.Fatalf("Failed to create file for deployment file for namespace: %s", ns.Name)
	}
	defer pkgDeploymentFile.Close()

	// write cmd output to deployment file
	cmd.Stdout = pkgDeploymentFile

	if err = cmd.Run(); err != nil {
		log.Fatalf("Failed to run kustomize build for namespace: %s", ns.Name)
	}
}