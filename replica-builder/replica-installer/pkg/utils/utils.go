package utils

import (
	"os"
	"os/exec"
	"log"
)

var (
	err error
)

func ValidateRequirements() {

	// validate kustomize CLI
	if _, err = exec.LookPath("kustomize"); err != nil {
		log.Fatalf("kustomize command not found")
	}
}

func CreateDir(filepath string) {

	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath, os.ModePerm); err != nil {
			log.Fatalf("Cannot create %v folder: %v", filepath, err)
		}
	}
}