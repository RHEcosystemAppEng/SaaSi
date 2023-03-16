FROM golang:1.19-buster as builder

ARG APP_NAME infra-exporter
ARG BUILD_VERSION
WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN ls -lh /tmp/
RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly  -v -o exporter -ldflags "-X main.BuildVersion=$BUILD_VERSION"

FROM alpine:3 as runner
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/exporter /go/bin/exporter

ENV OUTPUT_DIR=/tmp/output
EXPOSE 8080
ENV DEBUG=false
ENTRYPOINT ["/go/bin/exporter"]
