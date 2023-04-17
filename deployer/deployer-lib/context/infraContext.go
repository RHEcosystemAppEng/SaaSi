package context

import (
	"log"
	"os"
	"path"

	"github.com/RHEcosystemAppEng/SaaSi/deployer/deployer-lib/config"
)

const (
	ClustersFolder = "clusters"
)

type InfraContext struct {
	ScriptPath           string
	AnsiblePlaybookPath  string
	SourceClustersDir    string
	InfraRootDir         string
	OutputClustersFolder string
}

func InitInfraContext(flags *config.Args) *InfraContext {
	rootDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Cannot extract current working directory from environment, detailed error:  %s", err)
		return nil

	}
	ic := InfraContext{
		ScriptPath:           path.Join(rootDir, "infra/parser.sh"),
		AnsiblePlaybookPath:  path.Join(rootDir, "infra/playbook"),
		InfraRootDir:         path.Join(rootDir, "infra"),
		SourceClustersDir:    ClustersFolder,
		OutputClustersFolder: path.Join(rootDir, flags.RootOutputDir, ClustersFolder),
	}

	return &ic
}
