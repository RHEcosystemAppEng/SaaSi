# replica-builder
A tool to replicate an existing OpenShift/Kubernetes environment on multiple clusters and namespaces.

## Architecture
Consists of two components:
* [install-builder](./install-builder/README.md) a Golang CLI tool to extract and manipulate the configured resources 
  from a running environment and generate a reusable installer
* [replica-installer](./replica-installer/README.md) a Golang CLI tool to replicate the initial environment on different 
  clusters and namespaces

## Features
* Manage multiple namespaces
    * Allow mapping of per-namespace target
    * Provides predefines mapping policies
* Customize deployment parameters
    * Override only what is actually needed
    * Ensure that mandatory options are provided
    * Encrypt secrets
**TODO**

## Evaluating installer options
See [evaluation notes](./InstallerEvaluation.md)
