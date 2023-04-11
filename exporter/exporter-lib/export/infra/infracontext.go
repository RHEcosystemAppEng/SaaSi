package infra

import (
	"embed"
	_ "embed"
	"log"
	"path/filepath"

	"github.com/RHEcosystemAppEng/SaaSi/exporter/exporter-lib/config"
	"github.com/RHEcosystemAppEng/SaaSi/exporter/exporter-lib/connect"
	ctx "github.com/RHEcosystemAppEng/SaaSi/exporter/exporter-lib/context"
	"github.com/RHEcosystemAppEng/SaaSi/exporter/exporter-lib/export/utils"
	"github.com/sirupsen/logrus"
)

const (
	ClustersFolder = "clusters"
)

type InfraContext struct {
	ctx.ExporterContext

	clusterConfig *config.ClusterConfig
	ClusterFolder string
	scriptFolder  string
	ExportScript  string
}

//go:embed scripts/*
var scripts embed.FS

func NewInfraContextFromConfig(config *config.Config, clusterConfig *config.ClusterConfig, connectionStatus *connect.ConnectionStatus, logger *logrus.Logger) *InfraContext {
	context := InfraContext{clusterConfig: clusterConfig}

	context.InitFromConfig(config, connectionStatus, logger, config.Debug)

	var err error
	context.scriptFolder, err = utils.CopyEmbedderFolderToTempDir(scripts, "scripts")
	if err != nil {
		log.Fatalf("Cannot copy embedded scripts infra to temporary directory: %s", err)
	}
	context.ExportScript = filepath.Join(context.scriptFolder, "exporter.sh")
	context.ClusterFolder = filepath.Join(context.OutputFolder, ctx.ExportFolder, ClustersFolder, clusterConfig.Uid)
	return &context
}

func (c *InfraContext) InitClusterFolderForUid(uid string) string {
	c.ClusterFolder = filepath.Join(c.ClusterFolder, uid)
	return c.ClusterFolder
}
