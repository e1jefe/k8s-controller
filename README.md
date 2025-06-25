# Kubernetes Controller

[![CI Status](https://github.com/e1jefe/k8s-controller/workflows/K8S-controller/badge.svg)](https://github.com/e1jefe/k8s-controller/actions)
[![Release](https://img.shields.io/github/v/release/e1jefe/k8s-controller)](https://github.com/e1jefe/k8s-controller/releases)
[![Docker](https://img.shields.io/badge/docker-ghcr.io%2Fe1jefe%2Fk8s--controller-blue)](https://ghcr.io/e1jefe/k8s-controller)
[![Go Version](https://img.shields.io/badge/go-1.21-blue.svg)](https://golang.org/dl/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

A Kubernetes resource management tool built with [Cobra CLI](https://github.com/spf13/cobra) and [client-go](https://github.com/kubernetes/client-go), providing comprehensive Kubernetes resource management capabilities.

## Features

- **Kubernetes deployment listing** with kubeconfig authentication
- **Multiple authentication methods** (kubeconfig file, custom path, in-cluster)
- **Clean CLI interface** with intuitive commands
- **Docker support** with distroless images for security
- **Comprehensive build system** with Makefile
- **Kubernetes client-go integration** for cluster operations
- **CI/CD pipeline** with automated testing
- **List Deployments**: List all deployments in a specified namespace
- **Deployment Informer**: Watch for deployment changes (add, update, delete) with real-time logging

## Installation

### Local Development

```bash
# Clone the repository
git clone <repository-url>
cd go-k8s-controller

# Download dependencies
make deps

# Build the application
make build-local
```

### Using Make

```bash
# Build for production (Linux)
make build

# Build and show help
make run

# Build and run list deployments
make run-list-deployments

# Run all checks (format, test)
make check

# Show all available targets
make help
```

### Using Docker

```bash
# Build Docker image
make docker-build

# Push to registry
make docker-push
```

## Usage

### Basic Commands

```bash
# Show help
./bin/k8s-controller --help

# Show list command help
./bin/k8s-controller list --help
```

### Kubernetes Operations

#### Listing Deployments

```bash
# List deployments in default namespace (uses default kubeconfig)
./bin/k8s-controller list deployments

# List deployments with custom kubeconfig
./bin/k8s-controller list deployments --kubeconfig /path/to/kubeconfig

# List deployments with kubeconfig from custom location
./bin/k8s-controller list deployments --kubeconfig ~/.kube/custom-config
```

**Output Example:**
```
Found 2 deployments in default namespace:

NAME                           READY      UP-TO-DATE AVAILABLE  AGE            
----------------------------------------------------------------------------------------------
nginx                          2/2        2          2          5m             
web-app                        3/3        3          3          2h             
```

#### Authentication

The application supports multiple authentication methods:

1. **Kubeconfig file** (default: `~/.kube/config`)
2. **Custom kubeconfig path** via `--kubeconfig` flag
3. **In-cluster configuration** when running inside a Kubernetes pod

### Command Flags

#### Global Flags
- `--help, -h`: Show help information

#### List Flags
- `--kubeconfig`: Path to kubeconfig file (default: `$HOME/.kube/config`)

### Deployment Informer

The informer watches for real-time changes to Kubernetes deployment resources and logs detailed information about each event.

```bash
# Watch deployments in default namespace
./bin/k8s-controller informer

# Watch deployments in specific namespace
./bin/k8s-controller informer --namespace=kube-system

# Use in-cluster authentication (when running inside a pod)
./bin/k8s-controller informer --in-cluster

# Custom kubeconfig and resync period
./bin/k8s-controller informer --kubeconfig ~/.kube/config --resync-period=60s
```

### Informer Features

- **Real-time Event Detection**: Monitors ADD, UPDATE, and DELETE events for deployments
- **Detailed Logging**: Logs comprehensive deployment information including:
  - Name, namespace, and labels
  - Replica counts (desired, ready, available, updated)
  - Creation timestamps
  - Deployment conditions and status
- **Graceful Shutdown**: Handles SIGINT/SIGTERM for clean shutdown
- **Authentication Options**: Supports both kubeconfig and in-cluster authentication
- **Configurable Namespace**: Watch deployments in any namespace
- **Resync Period**: Configurable resync interval for cache synchronization

## Development

### Prerequisites

- Go 1.21 or later
- Docker (for containerized builds)
- Make
- Kubernetes cluster access (for testing list functionality)

### Available Make Targets

```bash
make all                    # Run clean, deps, fmt, test, and build
make build                  # Build binary for Linux (production)
make build-local            # Build binary for current OS
make clean                  # Clean build artifacts
make test                   # Run tests
make test-coverage          # Run tests with coverage report
make deps                   # Download and tidy dependencies
make fmt                    # Format code
make check                  # Run all checks (fmt, test)
make docker-build           # Build Docker image
make docker-push            # Push Docker image
make run                    # Build and show help
make run-list-deployments   # Build and run list deployments command
make test-k8s               # Test Kubernetes connectivity
make security               # Run security scan (requires gosec)
```

### Testing

```bash
# Run unit tests
make test

# Run tests with coverage
make test-coverage

# Test Kubernetes connectivity (requires cluster access)
make test-k8s
```

### Kubernetes Integration Testing

The project includes integration tests that run against a real Kubernetes cluster:

```bash
# Start local cluster (minikube/kind) and run
make test-k8s

# Or manually test with actual deployments
kubectl create deployment nginx --image=nginx --replicas=2
./bin/k8s-controller list deployments
kubectl delete deployment nginx
```

## Docker Deployment

The project includes a multi-stage Dockerfile using distroless images for security:

### Basic Usage

```bash
# Build image
docker build -t k8s-controller:latest .

# Show help
docker run k8s-controller:latest
```

### Kubernetes Operations with Docker

```bash
# Run with mounted kubeconfig (external cluster)
docker run -v ~/.kube/config:/tmp/kubeconfig:ro \
  k8s-controller:latest list deployments --kubeconfig /tmp/kubeconfig

# Run in Kubernetes cluster (in-cluster auth)
kubectl run k8s-controller \
  --image=k8s-controller:latest \
  --restart=Never \
  -- list deployments
```

### Docker Image Features

The Docker image:
- Uses Go 1.21 Alpine for building
- Creates a statically linked binary
- Uses distroless/static:nonroot for the final image
- Runs as non-root user
- Supports volume mounting for kubeconfig
- Multi-architecture support (linux/amd64, linux/arm64)

## Deployment Examples

### Local Development

```bash
# Install and test locally
make build-local
./bin/k8s-controller list deployments
```

### Kubernetes Job

```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: k8s-controller-list
spec:
  template:
    spec:
      containers:
      - name: k8s-controller
        image: k8s-controller:latest
        command: ["k8s-controller", "list", "deployments"]
      restartPolicy: Never
  backoffLimit: 4
```

### Kubernetes CronJob (Regular Deployment Checks)

```yaml
apiVersion: batch/v1
kind: CronJob
metadata:
  name: deployment-monitor
spec:
  schedule: "*/5 * * * *"  # Every 5 minutes
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: k8s-controller
            image: k8s-controller:latest
            command: ["k8s-controller", "list", "deployments"]
          restartPolicy: OnFailure
```

## CI/CD Pipeline

The project includes a comprehensive GitHub Actions workflow that:

- **Tests** the application with unit tests and coverage reporting
- **Kubernetes Integration Tests** using minikube to verify deployment listing functionality
- **Builds** multi-architecture Docker images
- **Pushes** to GitHub Container Registry
- **Validates** the entire pipeline on every push and pull request

## Troubleshooting

### Common Issues

1. **Kubeconfig not found**
   ```bash
   # Verify kubeconfig location
   ls -la ~/.kube/config
   
   # Test cluster connectivity
   kubectl cluster-info
   
   # Use custom kubeconfig
   ./k8s-controller list deployments --kubeconfig /path/to/config
   ```

2. **Permission denied errors**
   ```bash
   # Check cluster permissions
   kubectl auth can-i list deployments
   
   # Verify current context
   kubectl config current-context
   ```

3. **Docker volume mounting issues**
   ```bash
   # Ensure proper permissions on kubeconfig
   chmod 644 ~/.kube/config
   
   # Use absolute paths in volume mounts
   docker run -v /absolute/path/to/kubeconfig:/tmp/kubeconfig:ro ...
   ```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Run `make check` to ensure code quality
6. Submit a pull request

## License

This project is a Kubernetes resource management tool with comprehensive deployment listing capabilities. 