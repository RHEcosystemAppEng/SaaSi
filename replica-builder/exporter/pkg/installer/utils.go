package installer

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

func RunCommand(name string, options ...string) {
	cmd := exec.Command(name, options...)
	log.Printf("Running %s", cmd)
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
func RunCommandAndLogStderr(name string, options ...string) {
	cmd := exec.Command(name, options...)
	log.Printf("Running %s", cmd)
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	slurp, _ := io.ReadAll(stderr)
	log.Printf("%s\n", slurp)
	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}
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

func AppendToFile(file string, text string, args ...any) {
	f, err := os.OpenFile(file, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	defer f.Close()
	if len(args) > 0 {
		fmt.Fprintf(f, text, args...)
	} else {
		fmt.Fprint(f, text)
	}
}

func SystemNameForSA(namespace string, serviceAccount string) string {
	return fmt.Sprintf("system:serviceaccount:%s:%s", namespace, serviceAccount)
}
