# infra-exporter
An application to expose the `Infra Exporter` functionality as a REST POST service.

## Versioning
The version is printed at the application startup, as in:
```bash
Running /go/bin/exporter with version 0.1
```

## Run locally
Two environmenmt variables are defined:
* `OUTPUT_DIR` locates the root folder of the export file system. It's a mandatory parameter
* `DEBUG` is used to activate logging at debug level (extended to `crane` commands). It's an optional flag

Launch the application locally as:
```bash
OUTPUT_DIR=/tmp/output DEBUG=True go run main.go
```

The server listens on port `8080` (not configurable).

## Build the container image
Run the following to build the container image with a given version and push it to the remote registry (if needed):
```bash
export BUILD_VERSION=0.1
docker build --build-arg BUILD_VERSION=${BUILD_VERSION} -f ./Dockerfile -t FULL_IMAGE_NAME:IMAGE_TAG .
docker push FULL_IMAGE_NAME:IMAGE_TAG .
```

**Note**: by default the image is published at `quay.io/ecosystem-appeng/saasi-infra-exporter` by the GitHub actions defined in the repository at every new release.

## Running in OpenShift
You can test the application by creating a `Pod` and the `Route` to access it on on OpenShift cluster using the following:
```bash
oc apply -f resources.yaml
```

To log the execution:
```bash
oc logs -f saasi-infra-exporter
```

To invoke the exposed service from a test `test.json` file:
```bash
ROUTE=$(oc get route saasi-infra-exporter -ojsonpath='{.spec.host}') && curl ${ROUTE}/export/infra
ROUTE=$(oc get route saasi-infra-exporter -ojsonpath='{.spec.host}') && curl -X POST ${ROUTE}/export/infra -d @test.json
```

## REST APIs 
Two APIs are defined to fetch the applicaiton status and execute an export request.

### GET /export/infra
Returns the status of the application:
```json
{"Name":"infra-exporter","Version":"0.1","Status":"up"}
```

### POST /export/infra
Executes an export request. It expects an input following the schema presented in the [examples](../../examples/) folder for the infra exporters:
```json
{
  "clusterId": CLUSTER_ID,
  "server": SERVER_HOST,
  "token": VALID_TOKEN
}
```

The output contains the export status and the information to locate the exported installer:
```json
{
  "clusterId": CLUSTER_ID,
  "status": "ok",
  "errorMessage":"",
  "location": OUTPUT_FOLDER
}
```

Errors are either managed by standard HTTP status codes (e.g. `404 Not Found` or `405 Method Not Allowed`) or reported as the `failed` values in the `status` field, as in:
```json
{
  "clusterId": "",
  "status": "failed",
  "errorMessage":
  "Invalid configuration: missing cluster configuration",
  "location": ""
}
```
(the above response will return HTTP status `400 Bad Request`)