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

	// validate oc CLI
	if _, err = exec.LookPath("oc"); err != nil {
		log.Fatalf("oc command not found")
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