package utils

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

const (
	KUSTOMIZE = "kustomize"
	OC        = "oc"
)

var (
	err error
)

func ValidateRequirements(prog string) {

	// verify program exists
	if _, err = exec.LookPath(prog); err != nil {
		log.Fatalf("%s command not found, Error: %s", prog, err)
	}
}

func CreateDir(filepath string) {

	// create file path if not exists
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath, os.ModePerm); err != nil {
			log.Fatalf("Cannot create %v folder: %v", filepath, err)
		}
	}
}

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func StringifyMap(mapSrc map[string][]string) string {

	var mapStr string

	for key, value := range mapSrc {
		mapStr = fmt.Sprintf("%s%s:\n%v.\n", mapStr, key, strings.Join(value, ", "))
	}
	return mapStr
}
