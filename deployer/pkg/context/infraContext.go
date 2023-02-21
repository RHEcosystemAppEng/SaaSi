package context

import (
	"os"
	"path"
)

type InfraContext struct {

	RenderingScriptPath string
	AnsiblePlaybookPath string
}

func InitInfraContext() *InfraContext {
	rootDir, err := os.Getwd()
	if err != nil {
		return nil
	}
	ic := InfraContext{
		RenderingScriptPath: path.Join(rootDir,"infra") ,
		AnsiblePlaybookPath: path.Join(rootDir,"playbook"),
	}

	return &ic
}