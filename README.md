# Go Kubernetes Controller

A Go-based Kubernetes command-line interface tool that provides easy access to Kubernetes cluster resources. Built with Cobra CLI framework and the official Kubernetes Go client library.

## Features

- **List Pods** - View pod status, readiness, restarts, and age
- **List Deployments** - View deployment status, replicas, and availability  
- **List Services** - View service types, IPs, ports, and age
- **Namespace Support** - Target specific namespaces or view all namespaces
- **Automatic kubeconfig Detection** - Uses your existing Kubernetes configuration
- **Formatted Output** - Clean, tabular output similar to kubectl

## Project Structure

```
go-k8s-controller/
├── cmd/                    # Cobra CLI commands
│   ├── root.go            # Root command and CLI setup
│   ├── pods.go            # Pods listing command
│   ├── deployments.go     # Deployments listing command
│   ├── services.go        # Services listing command
│   └── k8s.go             # Kubernetes client utilities
├── main.go                # Application entry point
├── go.mod                 # Go module definition
└── go.sum                 # Go module checksums
```

## Quick Start

### Prerequisites

- Go 1.24.4 or later
- Access to a Kubernetes cluster
- Valid kubeconfig file (typically located at `~/.kube/config`)

### Installation

```bash
# Clone the repository
git clone <repository-url>
cd go-k8s-controller

# Build the CLI
go build -o k8s-cli

# Make it executable (optional)
chmod +x k8s-cli
```

### Usage Examples

```bash
# List pods in current namespace
./k8s-cli pods

# List all pods across all namespaces
./k8s-cli pods --all-namespaces
./k8s-cli pods -A

# List pods in specific namespace
./k8s-cli pods --namespace kube-system
./k8s-cli pods -n default

# List deployments
./k8s-cli deployments

# List deployments across all namespaces
./k8s-cli deployments --all-namespaces

# List services
./k8s-cli services

# List services in specific namespace
./k8s-cli services -n kube-system

# Get help for any command
./k8s-cli --help
./k8s-cli pods --help
```

## Available Commands

| Command | Description | Flags |
|---------|-------------|-------|
| `pods` | List pods in the cluster | `-n, --namespace`, `-A, --all-namespaces` |
| `deployments` | List deployments in the cluster | `-n, --namespace`, `-A, --all-namespaces` |
| `services` | List services in the cluster | `-n, --namespace`, `-A, --all-namespaces` |

## Dependencies

This project uses the following key dependencies:

- [Cobra](https://github.com/spf13/cobra) - Modern CLI framework for Go
- [Kubernetes Go Client](https://github.com/kubernetes/client-go) - Official Kubernetes Go client library

See `go.mod` for the complete list of dependencies.

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/new-feature`)
3. Make your changes
4. Add tests if applicable
5. Commit your changes (`git commit -am 'Add new feature'`)
6. Push to the branch (`git push origin feature/new-feature`)
7. Submit a pull request

## License

This project is open source and available under the [MIT License](LICENSE).