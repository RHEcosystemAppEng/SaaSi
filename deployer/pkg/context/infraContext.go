package context

import (
	"os"
	"path"
	"log"
)

const (
	ClustersFolder = "clusters"
)

type InfraContext struct {

	ScriptPath          string
	AnsiblePlaybookPath string
	SourceClustersDir   string
}

func InitInfraContext() *InfraContext {
	rootDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Cannot extract current working directory from environment, detailed error:  %s", err)
		return nil
	}
	ic := InfraContext{
		ScriptPath:          path.Join(rootDir,"infra/parser.sh") ,
		AnsiblePlaybookPath: path.Join(rootDir,"infra/playbook"),
		SourceClustersDir: ClustersFolder,
	}

	return &ic
}