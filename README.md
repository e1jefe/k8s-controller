# Go Kubernetes Controller

A Go-based Kubernetes toolkit that provides both programmatic access to Kubernetes APIs and a command-line interface for cluster management.

## Project Structure

This repository contains:

- **`cobra-cli/`** - A command-line interface for managing Kubernetes clusters, built with Go and Cobra

## Components

### Cobra CLI Tool

A simple and efficient command-line interface for managing Kubernetes clusters. The CLI provides essential Kubernetes operations in a clean, kubectl-like interface.

**Features:**
- List Pods, Deployments, and Services
- Namespace support (specific namespace or all namespaces)
- Automatic kubeconfig detection
- Cross-platform support (Linux, macOS, Windows)

See [`cobra-cli/README.md`](cobra-cli/README.md) for detailed usage instructions.

## Quick Start

```bash
# Navigate to the CLI tool
cd cobra-cli

# Build the CLI
go build -o k8s-cli

# List pods in current namespace
./k8s-cli pods

# List all pods across namespaces
./k8s-cli pods --all-namespaces
```

## Prerequisites

- Go 1.24.4 or later
- Access to a Kubernetes cluster
- Valid kubeconfig file

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is open source and available under the [MIT License](LICENSE).