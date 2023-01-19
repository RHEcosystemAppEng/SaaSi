package packager

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/RHEcosystemAppEng/SaaSi/deployer/pkg/config"
	"github.com/RHEcosystemAppEng/SaaSi/deployer/pkg/utils"
	"github.com/google/uuid"
)

const (
	SOURCE_KUSTOMIZE_DIR = "../exporter/output/Infinity/installer/kustomize"

	OUTPUT_DIR     = "output"
	KUSTOMIZE_DIR  = "kustomize"
	DEPLOYMENT_DIR = "deploy"
	TEMPLATE_DIR   = "template"
	CONFIGMAPS_DIR = "params"
	SECRETS_DIR    = "secrets"

	PARAM_FILE_EXT = ".env"

	EMPTY_PLACEHOLDER     = "__DEFAULT__"
	MANDATORY_PLACEHOLDER = "__MANDATORY__"

	COMMON_ANNOTATION_KEY = "app.kubernetes.io/saasi-pkg-uuid:"
)

var (
	err       error
	nsTmplDir string
)

type ApplicationPkg struct {
	Uuid                 uuid.UUID
	AppConfig            config.Application
	AppDir               string
	KustomizeDir         string
	DeloymentDir         string
	UnsetMandatoryParams map[string][]string
}

func NewApplicationPkg(appConfig config.Application) *ApplicationPkg {

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
	pkg.AppDir = filepath.Join(pwd, OUTPUT_DIR, appConfig.Name)
	utils.CreateDir(pkg.AppDir)
	// kustomize directory for namespace artifacts
	pkg.KustomizeDir = filepath.Join(pkg.AppDir, KUSTOMIZE_DIR)
	utils.CreateDir(pkg.KustomizeDir)
	// deployment directory for deployment packages
	pkg.DeloymentDir = filepath.Join(pkg.AppDir, DEPLOYMENT_DIR)
	utils.CreateDir(pkg.DeloymentDir)

	// init UnsetMandatoryParams to empty dict
	pkg.UnsetMandatoryParams = map[string][]string{}

	// in compliance with application config: generate a kustomize-able artifact, invoke cutomizations and build package
	pkg.generateApplicationPkg()

	return &pkg
}

func (pkg *ApplicationPkg) generateApplicationPkg() {
	for _, ns := range pkg.AppConfig.Namespaces {

		// define path to namespace template directory
		nsTmplDir = filepath.Join(pkg.KustomizeDir, ns.Name, TEMPLATE_DIR)

		// generate artifact for current namespace
		pkg.generateNsArtifact(ns)

		// invoke customizations onto artifact
		pkg.invokeNsCustomizations(ns)

		// pkg artifact
		pkg.buildNsDeployment(ns)
	}
}

func (pkg *ApplicationPkg) generateNsArtifact(ns config.Namespaces) {

	source := filepath.Join(SOURCE_KUSTOMIZE_DIR, ns.Name)
	// create pkg template at pkg template path
	cmd := exec.Command("cp", "-r", source, pkg.KustomizeDir)
	if err := cmd.Run(); err != nil {
		log.Fatalf("Failed to generate kustomize template for namespace: %s, Error: %s", ns.Name, err)
	}
}

func (pkg *ApplicationPkg) buildNsDeployment(ns config.Namespaces) {

	// define path to pkg deployment file
	nsDeploymentFilepath := filepath.Join(pkg.DeloymentDir, ns.Name+".yaml")

	// create pkg deployment file
	nsDeploymentFile, err := os.Create(nsDeploymentFilepath)
	if err != nil {
		log.Fatalf("Failed to create deployment file for namespace: %s, Error: %s", ns.Name, err)
	}
	defer nsDeploymentFile.Close()

	// build pkg with kustomize configuration
	cmd := exec.Command("kustomize", "build", nsTmplDir)
	cmd.Stdout = nsDeploymentFile
	if err = cmd.Run(); err != nil {
		log.Fatalf("Failed to build deployment file with kustomize for namespace: %s, Error: %s", ns.Name, err)
	}
}
