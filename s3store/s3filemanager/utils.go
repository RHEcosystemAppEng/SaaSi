package s3filemanager

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

func countOpenFiles() int64 {
	out, err := exec.Command("/bin/sh", "-c", fmt.Sprintf("lsof -p %v", os.Getpid())).Output()
	if err != nil {
		fmt.Println(err.Error())
	}
	lines := strings.Split(string(out), "\n")
	return int64(len(lines) - 1)
}

func ValidateEnvVariables() {
	if _, ok := os.LookupEnv("AWS_ACCESS_KEY_ID"); !ok {
		log.Fatalf("Missing mandatory environment variable AWS_ACCESS_KEY_ID")
	}
	if _, ok := os.LookupEnv("AWS_SECRET_ACCESS_KEY"); !ok {
		log.Fatalf("Missing mandatory environment variable AWS_SECRET_ACCESS_KEY")
	}
	if _, ok := os.LookupEnv("S3_ENDPOINT"); !ok {
		log.Fatalf("Missing mandatory environment variable S3_ENDPOINT")
	}
	if _, ok := os.LookupEnv("S3_REGION"); !ok {
		log.Fatalf("Missing mandatoworkbench.action.debug.configurery environment variable S3_REGION")
	}
}

func ConnectWithEnvVariables() (*session.Session, error) {
	ValidateEnvVariables()
	return ConnectS3Session(os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY"), os.Getenv("S3_ENDPOINT"), os.Getenv("S3_REGION"))
}

func ConnectS3Session(accessKey string, secretKey string, endpoint string, region string) (*session.Session, error) {
	var S3ForcePathStyle = true
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	return session.NewSession(&aws.Config{
		HTTPClient:       client,
		S3ForcePathStyle: &S3ForcePathStyle,
		Endpoint:         &endpoint,
		Region:           aws.String(region),
		Credentials:      credentials.NewStaticCredentials(accessKey, secretKey, ""),
	})
}
