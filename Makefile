# Makefile for k8s-controller

# Application name
APP_NAME := k8s-controller
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Go parameters
GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
GOGET := $(GOCMD) get
GOMOD := $(GOCMD) mod
GOFMT := $(GOCMD) fmt
GOLINT := golangci-lint

# Build parameters
BINARY_NAME := $(APP_NAME)
BINARY_UNIX := $(BINARY_NAME)_unix
BUILD_DIR := ./bin
DOCKER_REGISTRY ?= ghcr.io
DOCKER_IMAGE := $(DOCKER_REGISTRY)/e1jefe/$(APP_NAME)

# Ldflags
LDFLAGS := -ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.buildDate=$(BUILD_DATE)"

# Test variables
ENVTEST_K8S_VERSION = 1.29.0
ENVTEST_BIN_DIR = bin/k8s
ENVTEST_PLATFORM = $(shell go env GOOS)-$(shell go env GOARCH)

.PHONY: all build clean test test-coverage deps fmt lint vet check docker-build docker-push help envtest

# Default target
all: clean deps fmt test build

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .

# Build for current OS
build-local:
	@echo "Building $(BINARY_NAME) for local OS..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .

# Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	@chmod -R +w $(BUILD_DIR) 2>/dev/null || true
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out cover.out cmd/coverage.out

# Setup envtest
envtest:
	@echo "Setting up envtest..."
	@mkdir -p $(ENVTEST_BIN_DIR)
	@test -f $(ENVTEST_BIN_DIR)/setup-envtest || \
	GOBIN=$(PWD)/$(ENVTEST_BIN_DIR) go install sigs.k8s.io/controller-runtime/tools/setup-envtest@latest
	@$(ENVTEST_BIN_DIR)/setup-envtest use $(ENVTEST_K8S_VERSION) --bin-dir $(ENVTEST_BIN_DIR) >/dev/null

# Run tests with envtest
test: envtest
	@echo "Running tests with envtest..."
	cd cmd && KUBEBUILDER_ASSETS="../$(ENVTEST_BIN_DIR)/k8s/$(ENVTEST_K8S_VERSION)-$(ENVTEST_PLATFORM)" $(GOTEST) -v -timeout 60s

# Run tests with coverage
test-coverage: envtest
	@echo "Running tests with coverage..."
	cd cmd && KUBEBUILDER_ASSETS="../$(ENVTEST_BIN_DIR)/k8s/$(ENVTEST_K8S_VERSION)-$(ENVTEST_PLATFORM)" $(GOTEST) -v -race -coverprofile=coverage.out -timeout 60s
	cd cmd && $(GOCMD) tool cover -html=coverage.out -o coverage.html

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

# Format code
fmt:
	@echo "Formatting code..."
	$(GOFMT) ./...

# Vet code
vet:
	@echo "Vetting code..."
	$(GOCMD) vet ./...

# Run all checks
check: fmt vet test

# Build Docker image
docker-build:
	@echo "Building Docker image $(DOCKER_IMAGE):$(VERSION)..."
	docker build -t $(DOCKER_IMAGE):$(VERSION) -t $(DOCKER_IMAGE):latest .

# Push Docker image
docker-push:
	@echo "Pushing Docker image $(DOCKER_IMAGE):$(VERSION)..."
	docker push $(DOCKER_IMAGE):$(VERSION)
	docker push $(DOCKER_IMAGE):latest

# Run the application locally
run:
	@echo "Running $(BINARY_NAME)..."
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) . && $(BUILD_DIR)/$(BINARY_NAME) --help

# Run the informer
run-informer:
	@echo "Running $(BINARY_NAME) informer..."
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) . && $(BUILD_DIR)/$(BINARY_NAME) informer

# Run the list deployments command (requires Kubernetes cluster access)
run-list-deployments:
	@echo "Running $(BINARY_NAME) list deployments..."
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) . && $(BUILD_DIR)/$(BINARY_NAME) list deployments

# Test Kubernetes connectivity (requires cluster access)
test-k8s:
	@echo "Testing Kubernetes connectivity..."
	@if $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) . && $(BUILD_DIR)/$(BINARY_NAME) list deployments >/dev/null 2>&1; then \
		echo "✓ Kubernetes connectivity test passed"; \
	else \
		echo "✗ Kubernetes connectivity test failed - ensure you have a valid kubeconfig and cluster access"; \
	fi

# Install development tools
install-tools:
	@echo "Installing development tools..."
	@echo "No additional tools required for this simplified setup"

# Security scan
security:
	@echo "Running security scan..."
	@if command -v gosec >/dev/null 2>&1; then \
		gosec ./...; \
	else \
		echo "gosec not found. Install it with: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest"; \
	fi

# Install binary
install:
	@echo "Installing $(BINARY_NAME)..."
	$(GOCMD) install .

# Show help
help:
	@echo "Available targets:"
	@echo "  all                    - Run clean, deps, fmt, test, and build"
	@echo "  build                  - Build binary for Linux (production)"
	@echo "  build-local            - Build binary for current OS"
	@echo "  clean                  - Clean build artifacts"
	@echo "  envtest                - Setup envtest binaries"
	@echo "  test                   - Run tests with envtest"
	@echo "  test-coverage          - Run tests with coverage report"
	@echo "  deps                   - Download and tidy dependencies"
	@echo "  fmt                    - Format code"
	@echo "  vet                    - Vet code"
	@echo "  check                  - Run all checks (fmt, vet, test)"
	@echo "  docker-build           - Build Docker image"
	@echo "  docker-push            - Push Docker image"
	@echo "  run                    - Build and show help"
	@echo "  run-informer           - Build and run informer"
	@echo "  run-list-deployments   - Build and run list deployments command"
	@echo "  test-k8s               - Test Kubernetes connectivity"
	@echo "  install                - Install binary"
	@echo "  install-tools          - Install development tools"
	@echo "  security               - Run security scan"
	@echo "  help                   - Show this help message" 