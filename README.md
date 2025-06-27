# Kubernetes Controller

[![CI Status](https://github.com/e1jefe/k8s-controller/workflows/K8S-controller/badge.svg)](https://github.com/e1jefe/k8s-controller/actions)
[![Release](https://img.shields.io/github/v/release/e1jefe/k8s-controller)](https://github.com/e1jefe/k8s-controller/releases)
[![Docker](https://img.shields.io/badge/docker-ghcr.io%2Fe1jefe%2Fk8s--controller-blue)](https://ghcr.io/e1jefe/k8s-controller)
[![Go Version](https://img.shields.io/badge/go-1.21-blue.svg)](https://golang.org/dl/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

A lightweight Kubernetes management tool built with [Cobra CLI](https://github.com/spf13/cobra), [client-go](https://github.com/kubernetes/client-go), and [controller-runtime](https://github.com/kubernetes-sigs/controller-runtime), providing deployment management and real-time event monitoring capabilities.

## Features

- **📋 Deployment Listing**: List all deployments in specified namespaces
- **👁️ Real-time Informer**: Watch deployment changes with live event logging
- **🎯 Controller-Runtime**: Advanced controller with detailed event logging
- **🌐 JSON API Server**: Fast HTTP API with informer cache for programmatic access
- **🔐 Flexible Authentication**: Kubeconfig and in-cluster authentication support
- **🚀 Simple CLI**: Clean, intuitive command interface
- **🧪 Comprehensive Testing**: EnvTest integration with real Kubernetes API
- **🐳 Container Ready**: Docker support with distroless images
- **⚙️ CI/CD Ready**: Automated testing and building

## Quick Start

### Installation

```bash
# Clone and build
git clone <repository-url>
cd go-k8s-controller
make build-local

# Or use Docker
make docker-build
```

### Basic Usage

```bash
# List deployments
./bin/k8s-controller list deployments

# Watch deployment changes with basic informer
./bin/k8s-controller informer

# Run controller-runtime based controller with detailed logging
./bin/k8s-controller controller

# Start JSON API server
./bin/k8s-controller api

# Watch specific namespace
./bin/k8s-controller informer --namespace=kube-system
```

## Commands

### 📋 List Deployments

Display all deployments in a namespace with status information.

```bash
# Default namespace
./bin/k8s-controller list deployments

# Custom kubeconfig
./bin/k8s-controller list deployments --kubeconfig /path/to/config
```

**Example Output:**
```
Found 2 deployments in default namespace:

NAME                           READY      UP-TO-DATE AVAILABLE  AGE            
----------------------------------------------------------------------------------------------
nginx                          2/2        2          2          5m             
web-app                        3/3        3          3          2h             
```

### 👁️ Deployment Informer

Watch for real-time deployment changes and log events as they happen using basic informers.

```bash
# Watch default namespace
./bin/k8s-controller informer

# Watch specific namespace  
./bin/k8s-controller informer --namespace=production

# Custom kubeconfig
./bin/k8s-controller informer --kubeconfig ~/.kube/config
```

**Example Output:**
```
Starting informer for deployments in namespace: default
Informer running! Press Ctrl+C to stop...
Deployment ADDED: default/nginx-deployment
Deployment UPDATED: default/nginx-deployment  
Deployment DELETED: default/old-deployment
```

### 🎯 Controller-Runtime Controller (NEW!)

Run an advanced controller using `sigs.k8s.io/controller-runtime` with detailed event logging and reconciliation.

```bash
# Run controller watching all namespaces
./bin/k8s-controller controller

# Custom kubeconfig
./bin/k8s-controller controller --kubeconfig ~/.kube/config
```

**Example Output:**
```
Starting controller - watching Deployment events...
Deployment event        name=coredns namespace=kube-system replicas=1 ready=1 
                       image=rancher/mirrored-coredns-coredns:1.10.1 
                       time=2025-06-27T15:53:41+02:00

Deployment event        name=nginx-app namespace=default replicas=3 ready=2
                       image=nginx:1.21 time=2025-06-27T15:54:15+02:00

Deployment deleted      name=old-app namespace=default time=2025-06-27T15:55:02+02:00
```

**Key Features:**
- ⚡ **Controller-Runtime**: Uses the standard Kubernetes controller pattern
- 🔄 **Real-time Events**: Logs every CREATE, UPDATE, DELETE event
- 📊 **Structured Logging**: Clean key-value format for easy parsing
- 🎯 **Reconciliation**: Proper Kubernetes controller reconciliation loop
- 🔧 **Event Details**: Logs name, namespace, replicas, ready count, image, and timestamp

**What Events Are Logged:**
- **CREATE**: When new Deployments are created
- **UPDATE**: When Deployments are modified (spec changes, status updates)
- **DELETE**: When Deployments are removed

### 🌐 JSON API Server

Start a fast HTTP API server that provides JSON access to deployment data from informer cache.

```bash
# Start API server on default port 8080
./bin/k8s-controller api

# Custom port and namespace
./bin/k8s-controller api --port=9090 --namespace=production

# Custom kubeconfig
./bin/k8s-controller api --kubeconfig ~/.kube/config
```

**API Endpoint:**
```bash
# Get all deployments as JSON
curl http://localhost:8080/deployments
```

**Example JSON Response:**
```json
[
  {
    "name": "nginx-deployment",
    "namespace": "default",
    "replicas": 3,
    "ready": 2
  },
  {
    "name": "web-app",
    "namespace": "default", 
    "replicas": 5,
    "ready": 5
  }
]
```

**Key Features:**
- ⚡ **Fast Response**: 1-5ms using informer cache (vs 50-200ms direct API)
- 🔄 **Real-time Data**: Cache automatically syncs with Kubernetes
- 📊 **Simple Format**: Clean JSON with essential deployment info
- 🔧 **Easy Integration**: Perfect for scripts, monitoring, and automation

**Integration Examples:**
```bash
# Monitor deployment count
watch -n 5 'curl -s localhost:8080/deployments | jq "length"'

# Check if all deployments are ready
curl -s localhost:8080/deployments | jq '.[] | select(.ready != .replicas) | .name'

# Get deployment names only
curl -s localhost:8080/deployments | jq -r '.[].name'

# Use in monitoring scripts
TOTAL=$(curl -s localhost:8080/deployments | jq 'length')
echo "Found $TOTAL deployments"
```

### 🔐 Authentication

All commands support flexible authentication:

- **Kubeconfig**: Uses `~/.kube/config` by default
- **Custom Path**: `--kubeconfig /path/to/config`
- **In-cluster**: Automatic when running in Kubernetes pods

## Controller Architecture

The project provides multiple ways to watch Kubernetes Deployments:

### Basic Informer (`informer` command)
```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   CLI Informer  │    │  Client-Go      │    │  Kubernetes API │
│                 │───▶│  SharedInformer │───▶│                 │
│ • Simple logs   │    │  • Event Handler│    │ • Deployments   │
│ • Basic events  │    │  • Local Cache  │    │ • Real-time     │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

### Controller-Runtime (`controller` command)
```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│ Controller-RT   │    │  Manager &      │    │  Kubernetes API │
│                 │───▶│  Reconciler     │───▶│                 │
│ • Reconcile     │    │  • Work Queue   │    │ • Deployments   │
│ • Detailed logs │    │  • Error Retry  │    │ • Events        │
│ • Event handling│    │  • Rate Limit   │    │ • Real-time     │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## Development

### Prerequisites

- Go 1.21+
- Make
- Docker (optional)

### Building & Testing

```bash
# Development workflow
make deps          # Download dependencies
make fmt           # Format code
make test          # Run tests with envtest
make build-local   # Build for current OS

# Production build
make build         # Build for Linux

# Docker
make docker-build  # Build container image
```

### Testing

The project uses [EnvTest](https://book.kubebuilder.io/reference/envtest.html) for integration testing with real Kubernetes APIs:

```bash
# Setup test environment
make envtest

# Run tests
make test

# Run with coverage
make test-coverage
```

### Available Make Targets

| Target | Description |
|--------|-------------|
| `make build` | Build Linux binary (production) |
| `make build-local` | Build for current OS |
| `make test` | Run envtest integration tests |
| `make envtest` | Setup Kubernetes test binaries |
| `make clean` | Clean build artifacts |
| `make fmt` | Format code |
| `make vet` | Run go vet |
| `make docker-build` | Build Docker image |
| `make help` | Show all targets |

## CI/CD

The project includes a comprehensive GitHub Actions pipeline:

### 🧪 **Test Pipeline**
- **Unit Tests**: EnvTest with real Kubernetes API
- **Integration Tests**: Full deployment lifecycle testing with Kind
- **Platform Support**: Automatic detection (Linux/macOS, ARM64/AMD64)

### 🔄 **Build Pipeline**  
- **Multi-platform**: Linux and macOS binaries
- **Docker Images**: Distroless containers for security
- **Automated Releases**: Tagged versions with artifacts

### ✅ **Quality Checks**
- **Code Coverage**: Comprehensive test coverage reporting
- **Security Scanning**: Automated vulnerability detection
- **Code Formatting**: Enforced Go formatting standards

## Docker Usage

### Pre-built Images

```bash
# Pull and run
docker pull ghcr.io/e1jefe/k8s-controller:latest

# List deployments
docker run --rm -v ~/.kube:/root/.kube ghcr.io/e1jefe/k8s-controller list deployments

# Run informer
docker run --rm -v ~/.kube:/root/.kube ghcr.io/e1jefe/k8s-controller informer

# Run controller-runtime controller
docker run --rm -v ~/.kube:/root/.kube ghcr.io/e1jefe/k8s-controller controller

# Start API server
docker run --rm -p 8080:8080 -v ~/.kube:/root/.kube ghcr.io/e1jefe/k8s-controller api
```

### Building Locally

```bash
# Build image
make docker-build

# Run locally built image
docker run --rm -v ~/.kube:/root/.kube k8s-controller:latest --help
```

## Requirements

- **Kubernetes**: v1.20+ cluster access
- **Authentication**: Valid kubeconfig or in-cluster permissions
- **Permissions**: Read access to deployments in target namespaces

## Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   CLI Commands  │    │  Client-Go API  │    │  Kubernetes API │
│                 │───▶│                 │───▶│                 │
│ • list          │    │ • REST Client   │    │ • Deployments   │
│ • informer      │    │ • Informers     │    │ • Events        │
│ • controller    │    │ • Controller-RT │    │ • Real-time     │
│ • api (HTTP)    │    │ • Cache Store   │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
        │                       │
        │                       ▼
        │              ┌─────────────────┐
        └─────────────▶│  HTTP API       │
                       │                 │
                       │ GET /deployments│ ◀── JSON Clients
                       │ (Cached Data)   │     (curl, scripts, apps)
                       └─────────────────┘
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make changes with tests
4. Run `make check` to verify
5. Submit a pull request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details. 