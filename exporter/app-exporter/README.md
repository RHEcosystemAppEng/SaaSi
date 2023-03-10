# app-exporter
An application to export the `Exporter` functionality as a REST POST service.

## Versioning
The version is printed at the application startup, as in:
```bash
Running /go/bin/exporter with version 0.1
```

## Run locally
```bash
OUTPUT_DIR=/tmp/output DEBUG=True go run main.go
```

## Build the container image
```bash
export BUILD_VERSION=0.1
docker build --build-arg BUILD_VERSION=${BUILD_VERSION} -f ./Dockerfile -t quay.io/ecosystem-appeng/saasi-app-exporter:0.1 .
```

## Running in OpenShift
```bash
oc apply -f resources.yaml
```

## REST APIs 
**TODO**
### GET /export/application
### POST /export/application