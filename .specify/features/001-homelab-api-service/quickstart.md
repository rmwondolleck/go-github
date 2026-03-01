# Quickstart Guide - Home Lab API Service

Welcome to the Home Lab API Service! This guide will help you get up and running quickly with our REST API for managing your home lab services.

## Table of Contents
- [Prerequisites](#prerequisites)
- [Local Development Setup](#local-development-setup)
- [Example API Calls](#example-api-calls)
- [Docker Quickstart](#docker-quickstart)
- [Kubernetes Quickstart](#kubernetes-quickstart)
- [Next Steps](#next-steps)

## Prerequisites

Before you begin, ensure you have the following installed:

- **Go 1.25.0 or later** - [Download Go](https://go.dev/dl/)
- **Make** (optional, but recommended) - Usually pre-installed on Unix systems
- **Docker** (for Docker deployment) - [Install Docker](https://docs.docker.com/get-docker/)
- **kubectl** (for Kubernetes deployment) - [Install kubectl](https://kubernetes.io/docs/tasks/tools/)
- **curl** or **httpie** (for testing API endpoints)

## Local Development Setup

### Quick Start (3 steps)

1. **Clone the repository:**
   ```bash
   git clone https://github.com/rmwondolleck/go-github.git
   cd go-github
   ```

2. **Build the application:**
   ```bash
   make build
   ```
   Or without Make:
   ```bash
   go build -o bin/homelab-api ./cmd/api
   ```

3. **Run the application:**
   ```bash
   make run
   ```
   Or directly:
   ```bash
   ./bin/homelab-api
   ```

The API server will start on port 8080 by default. You should see output like:
```
2026/03/01 15:14:23 INFO server started port=8080
```

### Development with Hot Reload

For a better development experience with automatic reload on code changes:

1. **Install Air:**
   ```bash
   go install github.com/cosmtrek/air@latest
   ```

2. **Run with hot reload:**
   ```bash
   make dev
   ```

### Configuration

The API can be configured using environment variables:

- `PORT` - Server port (default: 8080)

Example:
```bash
PORT=3000 ./bin/homelab-api
```

## Example API Calls

Once the server is running, you can test the API endpoints:

### Health Check

Check if the API is running and healthy:

```bash
curl http://localhost:8080/health
```

**Response:**
```json
{
  "status": "ok"
}
```

**Response Codes:**
- `200 OK` - Service is healthy
- `503 Service Unavailable` - Service is unhealthy (if implemented)

### API Version Info

Get information about the API version:

```bash
curl http://localhost:8080/api/v1
```

**Response:**
```json
{
  "message": "API v1"
}
```

### Testing with HTTPie (Alternative)

If you prefer HTTPie for a more readable output:

```bash
# Install HTTPie
pip install httpie

# Test endpoints
http localhost:8080/health
http localhost:8080/api/v1
```

### Advanced Testing

You can test the API with more complex scenarios:

#### Check Response Headers
```bash
curl -i http://localhost:8080/health
```

Expected headers include:
- `Content-Type: application/json`
- Request ID header (for tracking requests)

#### Load Testing (using Apache Bench)
```bash
# Install Apache Bench (usually comes with Apache)
ab -n 1000 -c 10 http://localhost:8080/health
```

This sends 1000 requests with 10 concurrent connections to test API performance.

## Docker Quickstart

### Building the Docker Image

Currently, a Dockerfile is not yet included in the repository. Here's a recommended Dockerfile for building the image:

**Create `deployments/Dockerfile`:**
```dockerfile
# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o homelab-api ./cmd/api

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/homelab-api .

EXPOSE 8080

CMD ["./homelab-api"]
```

### Build and Run with Docker

1. **Build the Docker image:**
   ```bash
   make docker
   ```
   Or manually:
   ```bash
   docker build -t homelab-api:latest -f deployments/Dockerfile .
   ```

2. **Run the container:**
   ```bash
   docker run -d -p 8080:8080 --name homelab-api homelab-api:latest
   ```

3. **Test the containerized API:**
   ```bash
   curl http://localhost:8080/health
   ```

4. **View logs:**
   ```bash
   docker logs -f homelab-api
   ```

5. **Stop and remove the container:**
   ```bash
   docker stop homelab-api
   docker rm homelab-api
   ```

### Docker Compose (Optional)

For easier management, create a `docker-compose.yml`:

```yaml
version: '3.8'

services:
  homelab-api:
    build:
      context: .
      dockerfile: deployments/Dockerfile
    ports:
      - "8080:8080"
    environment:
      - PORT=8080
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "--quiet", "--tries=1", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 5s
      retries: 3
```

Run with:
```bash
docker-compose up -d
```

## Kubernetes Quickstart

### Prerequisites for K8s Deployment

- A running Kubernetes cluster (minikube, kind, or a cloud provider)
- kubectl configured to connect to your cluster
- Docker image built and available (either locally or in a registry)

### Quick Deploy to Kubernetes

1. **Ensure your Docker image is available:**
   
   If using minikube or kind (local development):
   ```bash
   # For minikube
   eval $(minikube docker-env)
   docker build -t homelab-api:latest -f deployments/Dockerfile .
   
   # For kind
   docker build -t homelab-api:latest -f deployments/Dockerfile .
   kind load docker-image homelab-api:latest
   ```
   
   For production, push to a container registry:
   ```bash
   docker tag homelab-api:latest your-registry/homelab-api:latest
   docker push your-registry/homelab-api:latest
   ```

2. **Deploy to Kubernetes:**
   ```bash
   kubectl apply -f deployments/k8s/deployment.yaml
   ```

3. **Create a Service to expose the deployment:**
   
   Create `deployments/k8s/service.yaml`:
   ```yaml
   apiVersion: v1
   kind: Service
   metadata:
     name: homelab-api-service
     labels:
       app: homelab-api
   spec:
     type: LoadBalancer
     ports:
       - port: 80
         targetPort: 8080
         protocol: TCP
         name: http
     selector:
       app: homelab-api
   ```
   
   Apply the service:
   ```bash
   kubectl apply -f deployments/k8s/service.yaml
   ```

4. **Verify the deployment:**
   ```bash
   # Check deployment status
   kubectl get deployments
   
   # Check pods
   kubectl get pods -l app=homelab-api
   
   # Check service
   kubectl get services homelab-api-service
   ```

5. **Access the API:**
   
   For minikube:
   ```bash
   minikube service homelab-api-service
   ```
   
   For LoadBalancer (cloud):
   ```bash
   # Get the external IP
   kubectl get service homelab-api-service
   
   # Test the endpoint
   curl http://<EXTERNAL-IP>/health
   ```
   
   For NodePort (alternative):
   ```bash
   # Get the NodePort
   kubectl get service homelab-api-service -o jsonpath='{.spec.ports[0].nodePort}'
   
   # Test the endpoint
   curl http://<NODE-IP>:<NODE-PORT>/health
   ```

### Port Forwarding (Quick Testing)

For quick testing without creating a Service:

```bash
# Forward local port 8080 to pod port 8080
kubectl port-forward deployment/homelab-api 8080:8080

# In another terminal, test the API
curl http://localhost:8080/health
```

### Monitoring the Application

```bash
# View logs from all pods
kubectl logs -l app=homelab-api --tail=100 -f

# View logs from a specific pod
kubectl logs <pod-name> -f

# Check pod events
kubectl describe pod <pod-name>

# Check deployment events
kubectl describe deployment homelab-api
```

### Scaling the Application

```bash
# Scale to 3 replicas
kubectl scale deployment homelab-api --replicas=3

# Verify scaling
kubectl get pods -l app=homelab-api
```

### Health Checks in Kubernetes

The deployment already includes liveness and readiness probes:

- **Liveness Probe**: Checks `/health` endpoint every 10 seconds
- **Readiness Probe**: Checks `/health` endpoint every 5 seconds

These ensure that:
- Unhealthy pods are automatically restarted
- Traffic is only sent to ready pods
- Zero-downtime deployments are possible

### Clean Up

To remove the deployment from Kubernetes:

```bash
kubectl delete -f deployments/k8s/deployment.yaml
kubectl delete -f deployments/k8s/service.yaml
```

Or delete by name:
```bash
kubectl delete deployment homelab-api
kubectl delete service homelab-api-service
```

## Next Steps

Now that you have the API running, here are some suggested next steps:

### Development

1. **Run Tests:**
   ```bash
   make test
   ```
   This runs the full test suite with coverage reporting.

2. **Run Linter:**
   ```bash
   make lint
   ```
   Ensure code quality with golangci-lint.

3. **Generate API Documentation:**
   ```bash
   make swagger
   ```
   Generate Swagger/OpenAPI documentation (if configured).

### Explore the Codebase

- `cmd/api/main.go` - Application entry point
- `internal/server/` - HTTP server and routing
- `internal/middleware/` - HTTP middleware (logging, request ID, recovery)
- `internal/models/` - Data models
- `internal/health/` - Health check implementation
- `tests/integration/` - Integration tests

### Learn More

- Read the [Feature Specification](spec.md) for detailed requirements
- Review the [Implementation Plan](plan.md) for architecture decisions
- Check out the [Tasks](tasks.md) for implementation details
- Explore the main [README.md](../../../README.md) for project overview

### Contribute

Interested in contributing? Check out:
- Coding conventions in the project
- Test coverage requirements
- Pull request process

### Monitoring and Observability

Future enhancements may include:
- Prometheus metrics endpoint
- Distributed tracing with OpenTelemetry
- Structured logging with correlation IDs
- Grafana dashboards

### Need Help?

- Check existing issues in the repository
- Review the documentation in `.specify/features/`
- Look at test files for usage examples

---

**Happy coding!** 🚀
