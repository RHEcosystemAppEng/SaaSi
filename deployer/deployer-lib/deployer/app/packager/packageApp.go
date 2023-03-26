package packager

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	regex "regexp"

	"github.com/RHEcosystemAppEng/SaaSi/deployer/deployer-lib/config"
	"github.com/RHEcosystemAppEng/SaaSi/deployer/deployer-lib/context"
	"github.com/RHEcosystemAppEng/SaaSi/deployer/deployer-lib/utils"
	"github.com/google/uuid"
	"gopkg.in/yaml.v2"
)

const (
	APPLICATION_DIR      = "applications"
	KUSTOMIZE_DIR        = "kustomize"
	DEPLOYMENT_DIR       = "deploy"
	TEMPLATE_DIR         = "template"
	CONFIGMAPS_DIR       = "params"
	SECRETS_DIR          = "secrets"
	TARGET_NAMESPACE_DIR = "target-namespaces"

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
	AppConfig            config.ApplicationConfig
	DeployerContext      context.DeployerContext
	UuidDir              string
	KustomizeDir         string
	DeloymentDir         string
	TargetNamespaceDir   string
	UnsetMandatoryParams map[string][]string
	Error                error
}

func NewApplicationPkg(appConfig config.ApplicationConfig, deployerContext *context.DeployerContext) *ApplicationPkg {

	// init ApplicationPkg
	pkg := ApplicationPkg{}

	// generate uuid
	pkg.Uuid = uuid.New()

	// assign configuration
	pkg.AppConfig = appConfig

	// assign deployer context
	pkg.DeployerContext = *deployerContext

	// create application directories
	// unique application directory by uuid
	pkg.UuidDir = filepath.Join(pkg.DeployerContext.GetRootOutputDir(), APPLICATION_DIR, pkg.AppConfig.Name, pkg.Uuid.String())
	utils.CreateDir(pkg.UuidDir)
	// kustomize directory for namespace artifacts
	pkg.KustomizeDir = filepath.Join(pkg.UuidDir, KUSTOMIZE_DIR)
	utils.CreateDir(pkg.KustomizeDir)
	// deployment directory for deployment packages
	pkg.DeloymentDir = filepath.Join(pkg.UuidDir, DEPLOYMENT_DIR)
	utils.CreateDir(pkg.DeloymentDir)
	// target namespace directory for target namespace resources
	pkg.TargetNamespaceDir = filepath.Join(pkg.DeloymentDir, TARGET_NAMESPACE_DIR)
	utils.CreateDir(pkg.TargetNamespaceDir)

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

		pkg.establishTargetNsName(&ns)

		// generate artifact for current namespace
		pkg.generateNsArtifact(ns)
		if pkg.Error != nil {
			return
		}

		// invoke customizations onto artifact
		pkg.invokeNsCustomizations(ns)
		if pkg.Error != nil {
			return
		}

		// pkg artifact
		pkg.buildNsDeployment(ns)
		if pkg.Error != nil {
			return
		}

		// create target namespace resource
		pkg.buildTargetNsResource(ns)
	}
}

func (pkg *ApplicationPkg) generateNsArtifact(ns config.Namespaces) {

	source := filepath.Join(pkg.DeployerContext.GetRootSourceDir(), APPLICATION_DIR, pkg.AppConfig.Name, KUSTOMIZE_DIR, ns.Name)
	// create pkg template at pkg template path
	cmd := exec.Command("cp", "-r", source, pkg.KustomizeDir)
	if err := cmd.Run(); err != nil {
		pkg.DeployerContext.GetLogger().Errorf("Failed to generate kustomize template for namespace: %s", ns.Name)
		pkg.Error = err
	}
}

func (pkg *ApplicationPkg) buildNsDeployment(ns config.Namespaces) {

	// define path to pkg deployment file
	nsDeploymentFilepath := filepath.Join(pkg.DeloymentDir, ns.Name+".yaml")

	// create pkg deployment file
	nsDeploymentFile, err := os.Create(nsDeploymentFilepath)
	if err != nil {
		pkg.DeployerContext.GetLogger().Errorf("Failed to create deployment file for namespace: %s", ns.Name)
		pkg.Error = err
		return
	}
	defer nsDeploymentFile.Close()

	// build pkg with kustomize configuration
	cmd := exec.Command("kustomize", "build", nsTmplDir)
	cmd.Stdout = nsDeploymentFile
	if err = cmd.Run(); err != nil {
		pkg.DeployerContext.GetLogger().Errorf("Failed to build deployment file with kustomize for namespace: %s", ns.Name)
		pkg.Error = err
	}
}

func (pkg *ApplicationPkg) buildTargetNsResource(ns config.Namespaces) {

	// define path to pkg deployment file
	targetNsResourceFilepath := filepath.Join(pkg.TargetNamespaceDir, ns.Target+".yaml")

	// create pkg deployment file
	targetNsResourceFile, err := os.Create(targetNsResourceFilepath)
	if err != nil {
		pkg.DeployerContext.GetLogger().Errorf("Failed to create target namespace file %s for namespace: %s", targetNsResourceFilepath, ns.Name)
		pkg.Error = err
		return
	}
	defer targetNsResourceFile.Close()

	// marshal target namespace data to yaml
	targetNsResource := map[string]any{"apiVersion": "v1", "kind": "Namespace", "metadata": map[string]string{"name": ns.Target}}
	targetNsResourceData, err := yaml.Marshal(targetNsResource)
	if err != nil {
		pkg.DeployerContext.GetLogger().Errorf("Failed to marshal target namespace \"%s\" data to yaml for namespace: %s", ns.Target, ns.Name)
		pkg.Error = err
		return
	}

	// write target namespace data to yaml file
	err = ioutil.WriteFile(targetNsResourceFilepath, targetNsResourceData, 0)
	if err != nil {
		pkg.DeployerContext.GetLogger().Errorf("Failed to write target namespace \"%s\" data to file %s for namespace: %s", ns.Target, targetNsResourceFilepath, ns.Name)
		pkg.Error = err
	}
}

func (pkg *ApplicationPkg) establishTargetNsName(ns *config.Namespaces) {

	// if target namespace name is specified, keep it
	if ns.Target != "" {
		return
	}

	// if target namespace name is not specified,
	// check if namespace mapping format is specified, if it is, override target param
	if nsf := pkg.AppConfig.NamespaceMappingFormat; nsf != "" {
		if match, _ := regex.MatchString("(^%s\\S+)|(\\S+%s\\S+)|(\\S+%s$)", nsf); match {
			ns.Target = fmt.Sprintf(pkg.AppConfig.NamespaceMappingFormat, ns.Name)
			return
		} else {
			pkg.DeployerContext.GetLogger().Warningf("NamespaceMappingFormat does not match required format, using original namespace name: %s", ns.Name)
		}
	}

	// if target namespace name is not specified,
	// and namespace mapping format is not specified or does not match required format,
	// override target param with original namespace name
	ns.Target = ns.Name
}
