package config

import (
	"log"
	"os"
	"strings"
)

type EnvConfig struct {
	RootOutputDir string
	RootSourceDir string
	Debug         bool
}

func ParseEnvs() *EnvConfig {
	// init config for environment variable
	config := EnvConfig{}

	// get output directory variable
	rootOutputDir, ok := os.LookupEnv("OUTPUT_DIR")
	if !ok {
		log.Fatal("missing mandatory environment variable OUTPUT_DIR")
	}
	config.RootOutputDir = rootOutputDir

	// get source directory variable
	rootSourceDir, ok := os.LookupEnv("SOURCE_DIR")
	if !ok {
		log.Fatal("missing mandatory environment variable SOURCE_DIR")
	}
	config.RootSourceDir = rootSourceDir

	// get debug variable
	debug, ok := os.LookupEnv("DEBUG")
	config.Debug = ok && strings.ToLower(debug) == "true"

	return &config
}
