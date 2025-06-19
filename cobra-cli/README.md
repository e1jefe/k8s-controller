# Kubernetes CLI Tool

A simple and efficient command-line interface for managing Kubernetes clusters, built with Go and Cobra. This tool provides essential Kubernetes operations in a clean, kubectl-like interface.

## Features

- **List Pods** - View pods with status, readiness, and restart information
- **List Deployments** - Check deployment status with replica information
- **List Services** - Display services with networking details
- **Namespace Support** - Work with specific namespaces or all namespaces
- **Kubeconfig Integration** - Automatic detection of cluster configuration
- **Cross-Platform** - Works on Linux, macOS, and Windows

## Prerequisites

- **Go 1.24.4+** - [Download Go](https://golang.org/dl/)
- **Kubernetes Cluster Access** - Valid kubeconfig file
- **kubectl** (optional) - For cluster setup and verification

## Installation

### Option 1: Build from Source

1. **Clone or download the project:**
   ```bash
   git clone <repository-url>
   cd go-k8s-controller/cobra-cli
   ```

2. **Install dependencies:**
   ```bash
   go mod tidy
   ```

3. **Build the executable:**
   ```bash
   go build -o k8s-cli
   ```

4. **Make it executable (Linux/macOS):**
   ```bash
   chmod +x k8s-cli
   ```

### Option 2: Direct Build

If you have the source code in a directory:

```bash
cd cobra-cli
go build -o k8s-cli
```

## Configuration

The CLI automatically detects your Kubernetes configuration using:

1. **In-cluster config** - If running inside a Kubernetes pod
2. **Kubeconfig file** - `~/.kube/config` (default location)
3. **KUBECONFIG environment variable** - Custom kubeconfig path

### Setting Custom Kubeconfig

```bash
export KUBECONFIG=/path/to/your/kubeconfig
./k8s-cli pods
```

## Usage

### Basic Commands

```bash
# Show help
./k8s-cli --help

# List available commands
./k8s-cli
```

### List Pods

```bash
# List pods in current namespace
./k8s-cli pods

# List pods in specific namespace
./k8s-cli pods --namespace kube-system
./k8s-cli pods -n default

# List pods in all namespaces
./k8s-cli pods --all-namespaces
./k8s-cli pods -A
```

**Example Output:**
```
NAME                          READY   STATUS    RESTARTS   AGE
nginx-deployment-66b6c48dd5   1/1     Running   0          2d
redis-master-6b5f7c9c8d      1/1     Running   0          1d
```

### List Deployments

```bash
# List deployments in current namespace
./k8s-cli deployments

# List deployments in specific namespace
./k8s-cli deployments --namespace kube-system
./k8s-cli deployments -n default

# List deployments in all namespaces
./k8s-cli deployments --all-namespaces
./k8s-cli deployments -A
```

**Example Output:**
```
NAME               READY   UP-TO-DATE   AVAILABLE   AGE
nginx-deployment   3/3     3            3           2d
redis-master       1/1     1            1           1d
```

### List Services

```bash
# List services in current namespace
./k8s-cli services

# List services in specific namespace
./k8s-cli services --namespace kube-system
./k8s-cli services -n default

# List services in all namespaces
./k8s-cli services --all-namespaces
./k8s-cli services -A
```

**Example Output:**
```
NAME         TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)   AGE
kubernetes   ClusterIP   10.96.0.1       <none>        443/TCP   5d
nginx-svc    NodePort    10.96.100.100   <none>        80:30080/TCP   2d
```

## Command Reference

### Global Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--namespace` | `-n` | Specify target namespace |
| `--all-namespaces` | `-A` | List resources from all namespaces |
| `--help` | `-h` | Show help information |

### Available Commands

| Command | Description |
|---------|-------------|
| `pods` | List pods in the cluster |
| `deployments` | List deployments in the cluster |
| `services` | List services in the cluster |
| `help` | Help about any command |

## Examples

### Common Workflows

**Check cluster status:**
```bash
./k8s-cli pods -A
./k8s-cli deployments -A
./k8s-cli services -A
```

**Monitor specific namespace:**
```bash
./k8s-cli pods -n production
./k8s-cli deployments -n production
./k8s-cli services -n production
```

**Quick health check:**
```bash
# Check if all deployments are ready
./k8s-cli deployments

# Look for failing pods
./k8s-cli pods
```

## Development

### Project Structure

```
cobra-cli/
├── go.mod                 # Go module file
├── go.sum                 # Dependencies checksum
├── main.go                # Application entry point
├── README.md              # This file
└── cmd/
    ├── root.go            # Root command definition
    ├── k8s.go             # Kubernetes client utilities
    ├── pods.go            # Pods command implementation
    ├── deployments.go     # Deployments command implementation
    └── services.go        # Services command implementation
```

### Adding New Commands

1. Create a new file in the `cmd/` directory
2. Define your command using Cobra patterns
3. Register it with the root command in the `init()` function
4. Rebuild the application

### Dependencies

- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [Kubernetes Client-Go](https://github.com/kubernetes/client-go) - Kubernetes API client
- [Kubernetes API](https://github.com/kubernetes/api) - Kubernetes API types

## Troubleshooting

### Common Issues

**"Error creating Kubernetes client"**
- Ensure your kubeconfig file is valid
- Check if you have access to the cluster: `kubectl cluster-info`
- Verify the kubeconfig path: `echo $KUBECONFIG`

**"No pods/deployments/services found"**
- Check if you're in the correct namespace
- Verify resources exist: `kubectl get pods -A`
- Ensure you have proper RBAC permissions

**"Permission denied"**
- Make the binary executable: `chmod +x k8s-cli`
- Check if you have cluster access permissions

### Getting Help

```bash
# Show general help
./k8s-cli --help

# Show help for specific command
./k8s-cli pods --help
./k8s-cli deployments --help
./k8s-cli services --help
```

## License

This project is open source and available under the [MIT License](LICENSE).

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

---

**Note:** This CLI tool is designed for basic Kubernetes operations. For advanced cluster management, consider using `kubectl` or other specialized tools. 