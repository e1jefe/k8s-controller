# JSON API for Kubernetes Deployments

This document describes how to use the JSON API server that provides access to Kubernetes deployment resources using informer cache storage.

## Quick Start

### 1. Start the API Server

```bash
# Start on default port 8080, watching default namespace
./k8s-controller api

# Start on custom port, watching specific namespace  
./k8s-controller api --port=9090 --namespace=kube-system

# Use custom kubeconfig
./k8s-controller api --kubeconfig=~/.kube/config
```

### 2. Available Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/` | GET | API information and available endpoints |
| `/health` | GET | Health check and cache sync status |
| `/api/v1/deployments` | GET | List all deployments from cache |
| `/api/v1/deployments/{name}` | GET | Get specific deployment by name |

## API Usage Examples

### List All Deployments

```bash
curl http://localhost:8080/api/v1/deployments
```

**Response:**
```json
{
  "items": [
    {
      "name": "nginx-deployment",
      "namespace": "default",
      "labels": {
        "app": "nginx"
      },
      "annotations": {},
      "replicas": 3,
      "readyReplicas": 3,
      "availableReplicas": 3,
      "updatedReplicas": 3,
      "creationTimestamp": "2024-01-01T10:00:00Z",
      "conditions": ["Available", "Progressing"]
    }
  ],
  "total": 1
}
```

### Get Specific Deployment

```bash
curl http://localhost:8080/api/v1/deployments/nginx-deployment
```

**Response:**
```json
{
  "name": "nginx-deployment",
  "namespace": "default",
  "labels": {
    "app": "nginx"
  },
  "replicas": 3,
  "readyReplicas": 3,
  "availableReplicas": 3,
  "updatedReplicas": 3,
  "creationTimestamp": "2024-01-01T10:00:00Z",
  "conditions": ["Available", "Progressing"]
}
```

### Health Check

```bash
curl http://localhost:8080/health
```

**Response:**
```json
{
  "status": "healthy",
  "timestamp": "2024-01-01T10:30:00Z",
  "cache_synced": true,
  "namespace": "default"
}
```

### Error Response Format

When an error occurs, the API returns:
```json
{
  "error": "Not Found",
  "message": "Deployment 'unknown-app' not found in namespace 'default'"
}
```

## Key Features

### 🚀 **High Performance**
- Uses Kubernetes informer cache for lightning-fast responses
- No direct API server calls for read operations
- Real-time data updates as cache syncs with cluster

### 📊 **Rich Data Format**
- Complete deployment information including status
- Ready/Available/Updated replica counts
- Creation timestamps and labels
- Deployment conditions (Available, Progressing, etc.)

### 🔄 **Real-time Updates**
- Informer automatically syncs with Kubernetes API
- Cache stays up-to-date with cluster changes
- 30-second sync interval (configurable)

### 🛡️ **Error Handling**
- Proper HTTP status codes
- Structured JSON error responses
- Graceful handling of missing resources

## Advanced Usage

### Watch for Changes

```bash
# Monitor deployment count every 5 seconds
watch -n 5 'curl -s http://localhost:8080/api/v1/deployments | jq ".total"'

# Watch specific deployment status
watch -n 2 'curl -s http://localhost:8080/api/v1/deployments/my-app | jq ".readyReplicas"'
```

### Save Data to File

```bash
# Save all deployments to JSON file
curl -s http://localhost:8080/api/v1/deployments > deployments.json

# Pretty print with jq
curl -s http://localhost:8080/api/v1/deployments | jq '.' > deployments-pretty.json
```

### Check if Deployment Exists

```bash
# Returns 200 if exists, 404 if not found
STATUS=$(curl -s -o /dev/null -w '%{http_code}' http://localhost:8080/api/v1/deployments/my-app)
if [ "$STATUS" = "200" ]; then
  echo "Deployment exists"
else
  echo "Deployment not found"
fi
```

### Integration with Scripts

```bash
#!/bin/bash

# Get deployment count
TOTAL=$(curl -s http://localhost:8080/api/v1/deployments | jq '.total')
echo "Found $TOTAL deployments"

# Check if all deployments are ready
curl -s http://localhost:8080/api/v1/deployments | jq -r '.items[] | select(.readyReplicas != .replicas) | .name'
```

## Configuration Options

| Flag | Default | Description |
|------|---------|-------------|
| `--port` | `8080` | Port to run the API server on |
| `--namespace` | `default` | Kubernetes namespace to watch |
| `--kubeconfig` | `~/.kube/config` | Path to kubeconfig file |

## Performance Benefits vs Direct API

| Approach | Response Time | Resource Usage | Real-time Updates |
|----------|---------------|----------------|-------------------|
| **Informer Cache** | ~1-5ms | Very Low | ✅ Yes |
| Direct K8s API | ~50-200ms | High | ❌ No |

## Error Scenarios

| HTTP Code | Error | Cause |
|-----------|-------|-------|
| `200` | Success | Request completed successfully |
| `404` | Not Found | Deployment doesn't exist in namespace |
| `405` | Method Not Allowed | Using wrong HTTP method (only GET supported) |
| `500` | Internal Server Error | Cache error or server issue |

## Demo Script

Run the included demo script to see the API in action:

```bash
chmod +x examples/api_usage.sh
./examples/api_usage.sh
```

This script demonstrates all endpoints and provides examples of common usage patterns.

## Integration Examples

### Python

```python
import requests
import json

# Get all deployments
response = requests.get('http://localhost:8080/api/v1/deployments')
data = response.json()
print(f"Found {data['total']} deployments")

for deployment in data['items']:
    print(f"- {deployment['name']}: {deployment['readyReplicas']}/{deployment['replicas']} ready")
```

### JavaScript/Node.js

```javascript
const fetch = require('node-fetch');

async function getDeployments() {
  const response = await fetch('http://localhost:8080/api/v1/deployments');
  const data = await response.json();
  
  console.log(`Found ${data.total} deployments`);
  data.items.forEach(dep => {
    console.log(`- ${dep.name}: ${dep.readyReplicas}/${dep.replicas} ready`);
  });
}

getDeployments();
```

### Go

```go
package main

import (
    "encoding/json"
    "fmt"
    "net/http"
)

type DeploymentListResponse struct {
    Items []DeploymentResponse `json:"items"`
    Total int                  `json:"total"`
}

func main() {
    resp, err := http.Get("http://localhost:8080/api/v1/deployments")
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    var data DeploymentListResponse
    json.NewDecoder(resp.Body).Decode(&data)
    
    fmt.Printf("Found %d deployments\n", data.Total)
    for _, dep := range data.Items {
        fmt.Printf("- %s: %d/%d ready\n", dep.Name, dep.ReadyReplicas, dep.Replicas)
    }
}
```

---

**🎯 The JSON API provides fast, cached access to Kubernetes deployment data with real-time updates and a simple REST interface.** 