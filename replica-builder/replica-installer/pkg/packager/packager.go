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

type ApplicationPkg struct {
	Uuid 			uuid.UUID
	AppConfig 		*config.ApplicationConfig
	AppDir    		string
	KustomizeDir 	string
	DeloymentDir 	string
}

func NewApplicationPkg(appConfig *config.ApplicationConfig) *ApplicationPkg {
	
	// init ApplicationPkg
	pkg := ApplicationPkg{}

	// generate uuid
	pkg.Uuid = uuid.New()

	// assign configuration
	pkg.AppConfig = appConfig

	// create application directories
	// application directory
	pwd, _ := os.Getwd()
	log.Printf("Running from %v", pwd)
	pkg.AppDir = filepath.Join(pwd, OUTPUT_DIR, appConfig.Application.Name)
	utils.CreateDir(pkg.AppDir)
	// kustomize directory for namesapce artifacts
	pkg.KustomizeDir = filepath.Join(pkg.AppDir, KUSTOMIZE_DIR)
	utils.CreateDir(pkg.KustomizeDir)
	// deployment directory for deployment packages
	pkg.DeloymentDir = filepath.Join(pkg.AppDir, DEPLOYMENT_DIR)
	utils.CreateDir(pkg.DeloymentDir)

	// in compliance with application config: generate a kustomize-able artifact, invoke cutomizations and build package
	pkg.generateApplicationPkg()

	return &pkg
}

func (pkg *ApplicationPkg) generateApplicationPkg() {
	for _, ns := range pkg.AppConfig.Application.Namespaces {

		// generate artifact for current namespace
		pkg.generateNsArtifact(ns)
		
		// invoke customizations onto artifact
		pkg.invokeNsCustomizations(ns)

		// pkg artifact
		pkg.buildNsDeployment(ns)
	}
}

func (pkg *ApplicationPkg) generateNsArtifact(ns config.SourceNamespace) {
	
	source := filepath.Join(SOURCE_KUSTOMIZE_DIR, ns.Name)
	// create pkg template at pkg template path
	cmd := exec.Command("cp", "-r", source, pkg.KustomizeDir)
	if err := cmd.Run(); err != nil {
		log.Fatalf("Failed to generate kustomize template for namespace: %s", ns.Name)
	}
}

func (pkg *ApplicationPkg) buildNsDeployment(ns config.SourceNamespace) {

	// define path to namespace template directory
	nsTmplDir := filepath.Join(pkg.KustomizeDir, ns.Name, TEMPLATE_DIR)
	
	// define path to pkg deployment file
	nsDeploymentFilePath := filepath.Join(pkg.DeloymentDir, ns.Name + ".yaml")

	// create pkg deployment file
	nsDeploymentFile, err := os.Create(nsDeploymentFilePath)
	if err != nil {
		log.Fatalf("Failed to create deployment file for namespace: %s", ns.Name)
	}
	defer nsDeploymentFile.Close()

	// build pkg with kustomize configuration
	cmd := exec.Command("kustomize", "build", nsTmplDir)
	cmd.Stdout = nsDeploymentFile
	if err = cmd.Run(); err != nil {
		log.Fatalf("Failed to build deployment file with kustomize for namespace: %s", ns.Name)
	}
}