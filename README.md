# SaaSi
One Stop shop for hybrid cloud native application packing and deployment tools

## Software architecture
Composed by two components:
* [replica-builder](./replica-builder/README.md): a tool to replicate an existing OpenShift/Kubernetes environment on 
multiple clusters and namespaces. Consists of two components:
  * [install-builder](./replica-builder/install-builder/README.md) a Golang CLI tool to extract and manipulate the configured resources
    from a running environment and generate a reusable installer
  * [replica-installer](./replica-builder/replica-installer/README.md) a Golang CLI tool to replicate the initial environment on different
    clusters and namespaces
* [saas-engine](./saas-engine/README.md): TBD
  * Uses the [replica-installer](./replica-builder/replica-installer) tool defined in the `replica-builder` product

![](./images/architecture.jpg)

## ADRs placeholders
- reasons to select Konveyor tools
  - collateral goal is to extend the tools with custom transformers for our need
- reasons to use Golang with Konveyor tools script
  - future goal is to convert it to full Golang app
- extraction of ConfigMap keys
  - how they simplify the customization of selected properties 
