package infra

import (
	"path/filepath"

	"github.com/RHEcosystemAppEng/SaaSi/exporter/pkg/config"
	"github.com/RHEcosystemAppEng/SaaSi/exporter/pkg/connect"
	"github.com/RHEcosystemAppEng/SaaSi/exporter/pkg/context"
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

func NewInfraContextFromConfig(config *config.Config, connectionStatus *connect.ConnectionStatus) *InfraContext {
	context := InfraContext{clusterConfig: &config.Exporter.Cluster}

	context.InitFromConfig(config, connectionStatus)
	context.scriptFolder = filepath.Join(config.RootInstallationFolder, "infra")
	context.ExportScript = filepath.Join(context.scriptFolder, "exporter.sh")
	context.ClusterFolder = filepath.Join(context.OutputFolder, ClustersFolder, config.Exporter.Cluster.ClusterId)
	return &context
}
