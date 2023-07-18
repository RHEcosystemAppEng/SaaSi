package context

import (
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/RHEcosystemAppEng/SaaSi/deployer/deployer-lib/config"
	"github.com/sirupsen/logrus"
)

const (
	CLUSTER_DIR = "clusters"
)

type InfraContext struct {
	ScriptPath           string
	AnsiblePlaybookPath  string
	InfraRootDir         string
	SourceClustersDir    string
	OutputClustersFolder string
	Logger               *logrus.Logger
}

func InitInfraContext(args *config.Args, logger *logrus.Logger) *InfraContext {
	rootDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Cannot extract current working directory from environment, detailed error:  %s", err)
		return nil
	}

	ic := InfraContext{
		ScriptPath:           path.Join(filepath.Dir(rootDir), "deployer-lib/deployer/infra/scripts/parser.sh"),
		AnsiblePlaybookPath:  path.Join(filepath.Dir(rootDir), "deployer-lib/deployer/infra/scripts/playbook"),
		InfraRootDir:         path.Join(filepath.Dir(rootDir), "deployer-lib/deployer/infra/scripts"),
		SourceClustersDir:    path.Join(args.RootSourceDir, CLUSTER_DIR),
		OutputClustersFolder: path.Join(args.RootOutputDir, CLUSTER_DIR),
		Logger:               logger,
	}

	return &ic
}
