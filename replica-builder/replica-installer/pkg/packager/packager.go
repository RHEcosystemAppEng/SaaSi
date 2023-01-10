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
	Namespace string
	TemplateDirpath string
	DeploymentFilepath string
}

func NewPkg(pkgNs string) (*Package, error) {

	// struct new pkg
	Package := Package{
		Uuid: uuid.New(),
		Namespace: pkgNs,
		TemplateDirpath: "",
		DeploymentFilepath: "",
	}

	return &Package, nil
}

func (p *Package) GeneratePkgTemplate(kustomizePath string) error {

	// init path to pkg template
	pkgTmplPath := filepath.Join(kustomizePath, p.Namespace, p.Uuid.String())

	// init path to base template
	baseTmplPath := filepath.Join(kustomizePath, p.Namespace, TEMPLATE_DIR)

	// create pkg template at pkg template path
	cmd := exec.Command("cp", "-r", baseTmplPath, pkgTmplPath)
	if err = cmd.Run(); err != nil {
		return errors.New("Failed to create pkg template")
	}

	// update pkg struct
	p.TemplateDirpath = pkgTmplPath

	return nil
}

func (p *Package) InvokePkgCustomizations() error {

	// validate kustomize cli
	if err = validateRequirements(); err != nil {
		return errors.New("kustomize command not found")
	}

	// set namespace resource in pkg kustomize template
	cmd := exec.Command("kustomize", "edit", "set", "namespace", p.Namespace)
	cmd.Dir = p.TemplateDirpath
	if err = cmd.Run(); err != nil {
		return errors.New("Failed to set namesapce by kustomize command")
	}

	// set common annotation resource in pkg kustomize template
	cmd = exec.Command("kustomize", "edit", "set", "annotation", COMMON_ANNOTATION_KEY + p.Uuid.String())
	cmd.Dir = p.TemplateDirpath
	if err = cmd.Run(); err != nil {
		return errors.New("Failed to set common annotations by kustomize command")
	}

	return nil
}

func (p *Package) BuildPkg() error {

	// build pkg with kustomize configuration
	cmd := exec.Command("kustomize", "build", p.TemplateDirpath)

	// init path to pkg deployment file
	pkgDeploymentPath := filepath.Join(p.TemplateDirpath, DEPLOYMENT_FILE_NAME)
	
	// create pkg deployment file
	pkgDeploymentFile, err := os.Create(pkgDeploymentPath)
	if err != nil {
		return errors.New("Failed to create file for deployment file")
	}
	defer pkgDeploymentFile.Close()

	// write cmd output to deployment file
	cmd.Stdout = pkgDeploymentFile

	if err = cmd.Run(); err != nil {
		return errors.New("Failed to run kustomize CLI command")
	}

	// update pkg struct
	p.DeploymentFilepath = pkgDeploymentPath

	return nil
}

// func (p *Package) DeployPkg() {


// }

// -----------------------
// --- Private Methods ---
// -----------------------

func validateRequirements() error {

	// validate kustomize CLI
	
	if _, err = exec.LookPath("kustomize"); err != nil {
		return errors.New("kustomize command not found")
	}
	return nil
}