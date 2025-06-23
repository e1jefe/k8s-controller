# Build stage
FROM golang:1.21-alpine AS builder

# Install git and ca-certificates
RUN apk add --no-cache git ca-certificates tzdata

# Create a non-root user for the build
RUN adduser -D -g '' appuser

# Set the working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download
RUN go mod verify

# Copy the source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo -o k8s-controller .

# Final stage - distroless
FROM gcr.io/distroless/static:nonroot

# Add labels for the image
LABEL org.opencontainers.image.title="k8s-controller"
LABEL org.opencontainers.image.description="Kubernetes resource management tool for listing deployments and managing cluster resources"
LABEL org.opencontainers.image.vendor="e1jefe"
LABEL org.opencontainers.image.source="https://github.com/e1jefe/k8s-controller"

# Import from builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/passwd /etc/passwd

# Copy the binary
COPY --from=builder /app/k8s-controller /k8s-controller

# Use non-root user (distroless nonroot user ID)
USER 65532:65532

# Run the binary
ENTRYPOINT ["/k8s-controller"]
CMD ["--help"]

# Usage examples:
# Show help:
#   docker run k8s-controller:latest
# List deployments (with kubeconfig):
#   docker run -v ~/.kube/config:/tmp/kubeconfig:ro k8s-controller:latest list deployments --kubeconfig /tmp/kubeconfig
# List deployments (in-cluster):
#   kubectl run k8s-controller --image=k8s-controller:latest --restart=Never -- list deployments 