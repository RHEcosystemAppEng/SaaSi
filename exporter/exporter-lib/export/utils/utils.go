package utils

import (
	"bufio"
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

type ResponseStatus int64

const (
	Ok ResponseStatus = iota
	Failed
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

func RunCommand(logger *logrus.Logger, name string, options ...string) error {
	cmd := exec.Command(name, options...)
	logger.Infof("Running %s", cmd)
	return cmd.Run()
}
func RunCommandAndLog(logger *logrus.Logger, name string, options ...string) error {
	cmd := exec.Command(name, options...)
	logger.Infof("Running %s", cmd)
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}

	multi := io.MultiReader(stdout, stderr)
	in := bufio.NewScanner(multi)
	for in.Scan() {
		logger.Infof(in.Text())
	}
	return nil
}

func BackupFile(file string) string {
	return strings.Join([]string{file, "bak"}, ".")
}

func FileContains(file string, text string) bool {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	s := string(b)
	return strings.Contains(s, text)
}

func ReplaceInFile(file string, original string, replacement string) {
	input, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("Cannot read from %s: %s", file, err)
	}

	os.Rename(file, BackupFile(file))
	output := bytes.Replace(input, []byte(original), []byte(replacement), -1)
	if err = ioutil.WriteFile(file, output, 0644); err != nil {
		log.Fatalf("Cannot write %s: %s", file, err)
	}
}

func AppendToFile(file string, text string, args ...any) error {
	f, err := os.OpenFile(file, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	defer f.Close()
	if len(args) > 0 {
		fmt.Fprintf(f, text, args...)
	} else {
		fmt.Fprint(f, text)
	}
	return nil
}

func SystemNameForSA(namespace string, serviceAccount string) string {
	return fmt.Sprintf("system:serviceaccount:%s:%s", namespace, serviceAccount)
}

func CopyEmbedderFolderToTempDir(embedFolder embed.FS, folderName string) (string, error) {
	tempDir, err := os.MkdirTemp("", folderName)
	if err != nil {
		return "", err
	}

	err = recursiveCopyEmbeddedFolderToTempDir(embedFolder, folderName, tempDir)
	return tempDir, err
}

func recursiveCopyEmbeddedFolderToTempDir(embedFolder embed.FS, folderName string, tempDir string) error {
	log.Printf("Copying embedded folder %s to: %s", folderName, tempDir)

	files, err := embedFolder.ReadDir(folderName)
	if err != nil {
		return err
	}
	for _, f := range files {
		if f.IsDir() {
			log.Printf("Creating folder %s", f.Name())
			folderName = filepath.Join(folderName, f.Name())
			tempDir := filepath.Join(tempDir, f.Name())
			os.Mkdir(tempDir, 0755)

			recursiveCopyEmbeddedFolderToTempDir(embedFolder, folderName, tempDir)
		} else {
			log.Printf("Copying file %s", f.Name())
			filePath := filepath.Join(folderName, f.Name())
			content, err := embedFolder.ReadFile(filePath)
			if err != nil {
				log.Fatalf("Cannot read from %s: %s", filePath, err)
				return err
			}

			destFile := filepath.Join(tempDir, f.Name())
			if err = os.WriteFile(destFile, content, 0755); err != nil {
				log.Fatalf("Cannot write %s: %s", destFile, err)
				return err
			}
		}
	}

	return nil
}

func PrettyPrint(logger *logrus.Logger, message string, obj any) {
	b, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		logger.Infof(message, obj)
	}
	logger.Infof(message, string(b))
}
