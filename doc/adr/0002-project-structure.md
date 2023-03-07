# 2. project structure

Date: 2023-03-07

## Status

Proposed

## Context

We need to define how to structure the projects that define the application, considering that it has multiple
management options, each one with its own assumptions:
* From CLI we want to either export or deploy both a cluster (infra) and a customer application
  * Ideally, also either a cluster or an application
* From REST, since each step might have different dependencies, resulting in different container images to be built,
we want separate services each one exposing a precise functionality, like /export/app or /export/infra

## Decision

The change we want to propose is to structure the code as follows:
```bash
├── exporter
│   ├── exporter-lib (library with functions for both steps)
│   │   └── pkg
│   │       ├── config
│   │       ├── connect
│   │       ├── context
│   │       └── export
│   ├── app-exporter (REST service for exporting an app)
│   └── infra-exporter (REST service for exporting a cluster)
├── deployer
│   ├── deployer-lib (library with functions for both steps)
│   │   └── pkg
│   │       ├── config
│   │       ├── connect
│   │       ├── context
│   │       ├── deployer
│   │       └── utils
│   ├── app-deployer (REST service for deploying an app)
│   └── infra-deployer (REST service for deploying a cluster)
├── saasi (CLI binary with options to `--deploy` or `--export`, input is given in `--file config.yaml`)
    └── pkg
        ├── config
        └── cmd
```

## Consequences

* Functional code is all defined in the `xyz-lib` modules and reusable but modules exposing client interface
* There is a single CLI tool to be developed (less code)
  * This can be managed in a different issue in GitHub and developed according to the right priority (e.g., maybe
  it's less urgent than the REST services)
* REST services are exposed by different applications so we can build specific images for each of them, each one with
its own dependencies (Ansible, crane, kustomize and so on)