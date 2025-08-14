# Jumpstarter Tekton Tasks

This repository provides a collection of reusable **Tekton Tasks** and **Pipeline** examples designed to work with [Jumpstarter](https://jumpstarter.dev). These tasks simplify the integration between Tekton and Jumpstarter, enabling secure and dynamic workload execution in Kubernetes.

## Available Tasks

### 1. `jumpstarter-get-lease`

Requests a lease from the Jumpstarter controller.  
Leases are used to authorize and manage time-bound or environment-bound access to resources.

### 2. `jumpstarter-release-lease`

Releases a previously acquired lease to ensure proper resource cleanup and avoid leaks.

### 3. `jumpstarter-run-cmd`

Executes a user-defined command inside a leased environment.  
Requires a valid lease to be available from a previous task.

### 4. `jumpstarter-setup-sa-client`

Generates a Jumpstarter `ClientConfig` dynamically using task parameters and the pod's Kubernetes ServiceAccount.  
This eliminates the need for a pre-provided Kubernetes Secret, allowing per-run authentication to be configured securely and on the fly.

## Use Cases

These tasks can be composed into Tekton Pipelines to support:
- Controlled access to configured machines with exporters environments via leases,
- Secure command execution within those machines,
- Dynamic service account–based authentication without manual secret management.

## Getting Started

You can include these tasks in your Tekton Pipelines by referencing them in your pipeline YAML files. Example pipeline definitions and usage scenarios will be provided in the `examples/` directory.

## License

[LICENSE](LICENSE)
