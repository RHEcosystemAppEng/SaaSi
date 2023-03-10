package infra

import (
	"embed"
	_ "embed"
	"log"
	"path/filepath"

	"github.com/RHEcosystemAppEng/SaaSi/exporter/exporter-lib/config"
	"github.com/RHEcosystemAppEng/SaaSi/exporter/exporter-lib/connect"
	"github.com/RHEcosystemAppEng/SaaSi/exporter/exporter-lib/context"
	"github.com/RHEcosystemAppEng/SaaSi/exporter/exporter-lib/export/utils"
	"github.com/sirupsen/logrus"
)

const (
	ClustersFolder = "clusters"
)

type InfraContext struct {
	context.ExporterContext

	clusterConfig *config.ClusterConfig
	ClusterFolder string
	scriptFolder  string
	ExportScript  string
}

//go:embed scripts/*
var scripts embed.FS

func NewInfraContextFromConfig(config *config.Config, exporterConfig *config.ExporterConfig, connectionStatus *connect.ConnectionStatus, logger *logrus.Logger) *InfraContext {
	context := InfraContext{clusterConfig: &exporterConfig.Cluster}

	context.InitFromConfig(config, connectionStatus, logger)

	var err error
	context.scriptFolder, err = utils.CopyEmbedderFolderToTempDir(scripts, "scripts")
	if err != nil {
		log.Fatalf("Cannot copy embedded scripts infra to temporary directory: %s", err)
	}
	context.ExportScript = filepath.Join(context.scriptFolder, "exporter.sh")
	context.ClusterFolder = filepath.Join(context.OutputFolder, ClustersFolder, exporterConfig.Cluster.ClusterId)
	return &context
}
