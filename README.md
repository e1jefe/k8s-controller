# Kubernetes Controller

A Kubernetes controller built with [fasthttp](https://github.com/valyala/fasthttp) and [Cobra CLI](https://github.com/spf13/cobra), providing a high-performance HTTP server for controller operations.

## Features

- High-performance HTTP server using FastHTTP
- RESTful API endpoints for health checks and server information
- Graceful shutdown with configurable timeout
- Docker support with multi-stage builds using distroless images
- Comprehensive build system with Makefile
- Request logging and monitoring capabilities

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

# Build and run locally
make run

# Build and run server directly
make run-server

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

# Show server command help
./bin/k8s-controller server --help
```

### Starting the Server

```bash
# Start server on default port 8080
./bin/k8s-controller server

# Start server on custom port
./bin/k8s-controller server --port 3000
./bin/k8s-controller server -p 3000

# Start server on custom host and port
./bin/k8s-controller server --host 0.0.0.0 --port 8080
./bin/k8s-controller server -H 0.0.0.0 -p 8080
```

### API Endpoints

The server provides the following endpoints:

- `GET /` - Welcome message
- `GET /health` - Health check endpoint
- `GET /info` - Server information including request details

### Server Flags

- `--port, -p`: Port to run the server on (default: 8080)
- `--host, -H`: Host address to bind the server to (default: localhost)

## Development

### Prerequisites

- Go 1.21 or later
- Docker (for containerized builds)
- Make

### Available Make Targets

```bash
make all             # Run clean, deps, fmt, test, and build
make build           # Build binary for Linux (production)
make build-local     # Build binary for current OS
make clean           # Clean build artifacts
make test            # Run tests
make test-coverage   # Run tests with coverage report
make deps            # Download and tidy dependencies
make fmt             # Format code
make check           # Run all checks (fmt, test)
make docker-build    # Build Docker image
make docker-push     # Push Docker image
make run             # Build and run application
make run-server      # Build and run server command
make security        # Run security scan (requires gosec)
```

### Testing

```bash
# Run tests
make test

# Run tests with coverage
make test-coverage
```

## Docker Deployment

The project includes a multi-stage Dockerfile using distroless images for security:

```bash
# Build image
docker build -t k8s-controller:latest .

# Run container
docker run -p 8080:8080 k8s-controller:latest
```

The Docker image:
- Uses Go 1.21 Alpine for building
- Creates a statically linked binary
- Uses distroless/static:nonroot for the final image
- Includes health checks
- Runs as non-root user
- Exposes port 8080

## License

This project is a Kubernetes controller with FastHTTP server integration. 