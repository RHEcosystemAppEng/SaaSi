package context

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/RHEcosystemAppEng/SaaSi/exporter/exporter-lib/config"
	"github.com/RHEcosystemAppEng/SaaSi/exporter/exporter-lib/connect"
	"k8s.io/client-go/rest"
)

type Context interface {
	RootFolder() string
}

type ExporterContext struct {
	OutputFolder     string
	ConnectionStatus *connect.ConnectionStatus
}

func (c *ExporterContext) InitFromConfig(config *config.Config, connectionStatus *connect.ConnectionStatus) {
	c.ConnectionStatus = connectionStatus
	c.OutputFolder = config.RootOutputFolder
}

func (c *ExporterContext) RootFolder() string {
	log.Fatal("Not implemented: RootFolder()")
	return ""
}

func (c *ExporterContext) KubeConfig() *rest.Config {
	return c.ConnectionStatus.KubeConfig
}

func (c *ExporterContext) KubeConfigPath() string {
	return c.ConnectionStatus.KubeConfigPath
}

func LookupOrCreateFolder(c Context, path ...string) string {
	fullPath := ""
	if strings.HasPrefix(path[0], c.RootFolder()) {
		fullPath = filepath.Join(path...)
	} else {
		fullPath = filepath.Join(append([]string{c.RootFolder()}, path...)...)
	}
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		if err := os.MkdirAll(fullPath, os.ModePerm); err != nil {
			log.Fatalf("Cannot create %v folder: %v", fullPath, err)
		}
	}
	return fullPath
}
