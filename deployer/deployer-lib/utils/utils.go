package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/sirupsen/logrus"
)

type ResponseStatus int64

const (
	KUSTOMIZE                = "kustomize"
	OC                       = "oc"
	Ok        ResponseStatus = iota
	Failed
)

var (
	err error
)

func (s ResponseStatus) String() string {
	switch s {
	case Ok:
		return "ok"
	case Failed:
		return "failed"
	}
	return "unknown"
}

func ValidateRequirements(prog string) error {

	// verify program exists
	if _, err = exec.LookPath(prog); err != nil {
		return err
	}

	return nil
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

func GetLogger(debug bool) *logrus.Logger {
	log := logrus.New()
	if debug {
		log.SetLevel(logrus.DebugLevel)
	}
	return log
}

func PrettyPrint(logger *logrus.Logger, message string, obj any) {
	b, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		logger.Infof(message, obj)
	}
	logger.Infof(message, string(b))
}
