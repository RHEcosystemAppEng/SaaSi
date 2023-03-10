package context

import (
	"github.com/RHEcosystemAppEng/SaaSi/deployer/pkg/config"
	"log"
	"os"
	"path"
)

const (
	ClustersFolder = "clusters"
)

type InfraContext struct {
	ScriptPath          string
	AnsiblePlaybookPath string
	SourceClustersDir   string
	InfraRootDir        string
	OutputClustersFolder string
}

func InitInfraContext(flags config.FlagArgs) *InfraContext {
	rootDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Cannot extract current working directory from environment, detailed error:  %s", err)
		return nil

	}
	ic := InfraContext{
		ScriptPath:          path.Join(rootDir,"infra/parser.sh") ,
		AnsiblePlaybookPath: path.Join(rootDir,"infra/playbook"),
		InfraRootDir: path.Join(rootDir,"infra"),
		SourceClustersDir: ClustersFolder,
		OutputClustersFolder: path.Join(rootDir,flags.RootOutputDir,ClustersFolder),
	}

	return &ic
}