package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/jessevdk/go-flags"
)

const PORT = "8080"

var err error

type Args struct {
	ConfigFile    string `short:"f" description:"Application configuration file for deployemnt" required:"true"`
	RootOutputDir string `short:"o" long:"output-dir" default:"output" description:"Root output folder"`
	RootSourceDir string `short:"s" long:"source-dir" description:"Root source folder" required:"true"`
	Port          int
	Debug         bool
}

func ParseEnvs() *Args {

	// init args config for environment variables
	args := Args{}

	// get output directory variable
	rootOutputDir, ok := os.LookupEnv("OUTPUT_DIR")
	if !ok {
		log.Fatal("missing mandatory environment variable OUTPUT_DIR")
	}
	args.RootOutputDir = rootOutputDir

	// get source directory variable
	rootSourceDir, ok := os.LookupEnv("SOURCE_DIR")
	if !ok {
		log.Fatal("missing mandatory environment variable SOURCE_DIR")
	}
	args.RootSourceDir = rootSourceDir

	// get port variable
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = PORT
	}
	args.Port, err = strconv.Atoi(port)
	if err != nil {
		log.Fatalf("Invalid port %s configured", port)
	}

	// get debug variable
	debug, ok := os.LookupEnv("DEBUG")
	args.Debug = ok && strings.ToLower(debug) == "true"

	return &args
}

func ParseFlags() *Args {

	// init args config for input arguments
	args := Args{}

	// parse input arguments from os into args config
	_, err := flags.Parse(&args)
	if err != nil {
		log.Fatal("Failed to parse os input arguments")
	}

	return &args
}
