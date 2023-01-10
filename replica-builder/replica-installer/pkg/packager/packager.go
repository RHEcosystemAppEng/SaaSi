package packager

import (
	"os/exec"
	"os"
	"path/filepath"
	"errors"

	"github.com/google/uuid"
)

const (
	TEMPLATE_DIR = "template"
	COMMON_ANNOTATION_KEY = "app.kubernetes.io/saasi-pkg-uuid:"
	DEPLOYMENT_FILE_NAME = "deployment.yaml"
)

var (
	err error
)

type Package struct {
	Uuid uuid.UUID
	TemplatePath string
	Namespace string
}

func NewPkg(pkgNs string, kustomizePath string) (*Package, error) {

	// generate pkg uuid
	pkgUuid := uuid.New()
	
	// init pkg template
	pkgTmplPath, err := generatePkgTemplate(pkgUuid, pkgNs, kustomizePath)
	if err != nil {
		return nil, err
	}

	// struct new pkg
	Package := Package{
		Uuid: pkgUuid,
		TemplatePath: pkgTmplPath,
		Namespace: pkgNs,
	}

	return &Package, nil
}

func generatePkgTemplate(pkgUuid uuid.UUID, pkgNs string, kustomizePath string) (string, error) {

	// init path to pkg
	pkgTmplPath := filepath.Join(kustomizePath, pkgNs, pkgUuid.String())

	// init path to general template
	baseTmplPath := filepath.Join(kustomizePath, pkgNs, TEMPLATE_DIR)

	// init pkg template at pkg path
	cmd := exec.Command("cp", "-r", baseTmplPath, pkgTmplPath)
	err = cmd.Run()
	if err != nil {
		return "", err
	}

	return pkgTmplPath, nil
}

func (p *Package) ConfigurePkgKustomize() error {

	// validate kustomize cli
	err = validateRequirements()
	if err != nil {
		return errors.New("kustomize command not found")
	}

	// set namespace resource in pkg kustomize template
	cmd := exec.Command("kustomize", "edit", "set", "namespace", p.Namespace)
	cmd.Dir = p.TemplatePath
	err = cmd.Run()
	if err != nil {
		return errors.New("Failed to set namesapce by kustomize command")
	}

	// set common annotation resource in pkg kustomize template
	cmd = exec.Command("kustomize", "edit", "set", "annotation", COMMON_ANNOTATION_KEY + p.Uuid.String())
	cmd.Dir = p.TemplatePath
	err = cmd.Run()
	if err != nil {
		return errors.New("Failed to set common annotations by kustomize command")
	}

	return nil
}

func validateRequirements() error {

	// validate kustomize CLI
	_, err := exec.LookPath("kustomize")
	if err != nil {
		return errors.New("kustomize command not found")
	}
	return nil
}

func (p *Package) BuildPkg() error {

	// build pkg with kustomize configuration
	cmd := exec.Command("kustomize", "build", p.TemplatePath)

	// init path to deployment pkg
	pkgDeploymentPath := filepath.Join(p.TemplatePath, DEPLOYMENT_FILE_NAME)
	
	// create deployment file
	pkgDeploymentFile, err := os.Create(pkgDeploymentPath)
	if err != nil {
		return errors.New("Failed to create file for deployment file")
	}
	defer pkgDeploymentFile.Close()

	// write cmd output to deployment file
	cmd.Stdout = pkgDeploymentFile

	if err := cmd.Run(); err != nil {
		return errors.New("Failed to run kustomize CLI command")
	}

	return nil
}

// func (p *Package) DeployPkg() {

// }