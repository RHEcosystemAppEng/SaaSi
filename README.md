# SaaSi
One Stop shop for hybrid cloud native application packing and deployment tools

## Software architecture
A tool to replicate an existing OpenShift/Kubernetes environment on multiple clusters and namespaces. Consists of two components:
* [exporter](./exporter/README.md) a Golang CLI tool to extract and manipulate the configured resources
  from a running environment and generate a reusable installer
* [deployer](./deployer/README.md) a Golang CLI tool to replicate the initial environment on different
  clusters and namespaces

![](./images/architecture.jpg)

## Features
* Manage multiple namespaces
  * Allow mapping of per-namespace target
  * Provides predefines mapping policies
* Customize deployment parameters
  * Override only what is actually needed
  * Ensure that mandatory options are provided
  * Encrypt secrets
    **TODO**

## ADRs placeholders
- reasons to select Konveyor tools
  - collateral goal is to extend the tools with custom transformers for our need
- reasons to use Golang with Konveyor tools script
  - future goal is to convert it to full Golang app
- extraction of ConfigMap keys
  - how they simplify the customization of selected properties 
- all SaasI components should expose API
  - during the first phases some SaasI components might be triggerd from CLI however the CLI will call an API exposed by the SaaSi components  
  - future goal is to turn SaaSi components to services 
- develop UI application using react with typescript instead of using backstage.io
  - backstage.io customization requires more time to understand the framework
  - it is possible that at some point we may have to face backstage software catelog limitations for our use case
