# Deployment Guide

This guide provides comprehensive instructions for deploying the Home Lab API service using Docker and Kubernetes.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Docker Deployment](#docker-deployment)
  - [Building the Docker Image](#building-the-docker-image)
  - [Running with Docker](#running-with-docker)
  - [Docker Configuration](#docker-configuration)
- [Kubernetes Deployment](#kubernetes-deployment)
  - [Deployment Steps](#deployment-steps)
  - [Verifying Deployment](#verifying-deployment)
  - [Accessing the Service](#accessing-the-service)
- [Environment Variables](#environment-variables)
- [Configuration](#configuration)
- [Troubleshooting](#troubleshooting)
- [Monitoring and Health Checks](#monitoring-and-health-checks)

## Prerequisites

Before deploying, ensure you have:

- **Docker**: Version 20.10 or later
- **Kubernetes**: Version 1.25 or later (for K8s deployments)
- **kubectl**: Configured to access your cluster
- **Git**: To clone the repository
- **Go**: Version 1.25.0 or later (for building from source)

## Docker Deployment

### Dockerfile Options

The project provides two Dockerfile options:

1. **`Dockerfile`** (Default) - Uses Alpine Linux base image
   - **Image size**: ~35-40MB
   - **Base image**: `alpine:latest`
   - **Security**: Runs as non-root user (UID 1000)
   - **Features**: Includes curl for health checks, shell available for debugging
   - **Health checks**: Built-in Docker HEALTHCHECK using curl
   - **Best for**: Standard deployments with health check support

2. **`Dockerfile.distroless`** - Uses distroless base image
   - **Image size**: ~33MB
   - **Base image**: `gcr.io/distroless/static-debian12:nonroot`
   - **Security**: Runs as non-root user (UID 65532), minimal attack surface
   - **Features**: No shell, package manager, or unnecessary binaries
   - **Health checks**: Use Kubernetes liveness/readiness probes or external monitoring
   - **Best for**: Production deployments requiring maximum security and minimal size

### Building the Docker Image

The project uses a multi-stage Dockerfile for optimized image size (<50MB).

#### Build using Make

```bash
make docker
```

This builds the image as `homelab-api:latest`.

#### Build using Docker directly (default - Alpine)

```bash
docker build -t homelab-api:latest -f deployments/Dockerfile .
```

#### Build using distroless variant (smaller, more secure)

```bash
docker build -t homelab-api:distroless -f deployments/Dockerfile.distroless .
```

#### Build arguments

The Dockerfile supports build arguments for customization:

```bash
docker build \
  --build-arg GO_VERSION=1.25.0 \
  -t homelab-api:latest \
  -f deployments/Dockerfile .
```

### Running with Docker

#### Basic run

```bash
docker run -p 8080:8080 homelab-api:latest
```

#### Run with custom port

```bash
docker run -p 3000:3000 -e PORT=3000 homelab-api:latest
```

#### Run with environment variables

```bash
docker run -p 8080:8080 \
  -e PORT=8080 \
  -e LOG_LEVEL=debug \
  -e RATE_LIMIT=100 \
  homelab-api:latest
```

#### Run in detached mode

```bash
docker run -d -p 8080:8080 --name homelab-api homelab-api:latest
```

#### View logs

```bash
docker logs homelab-api
docker logs -f homelab-api  # Follow logs
```

### Docker Configuration

#### Resource Limits

Run with memory and CPU limits:

```bash
docker run -p 8080:8080 \
  --memory="100m" \
  --cpus="0.2" \
  homelab-api:latest
```

#### Health Check

The container includes health checks on the `/health` endpoint:

```bash
docker run -p 8080:8080 \
  --health-cmd="wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1" \
  --health-interval=10s \
  --health-timeout=5s \
  --health-retries=3 \
  homelab-api:latest
```

## Kubernetes Deployment

### Deployment Steps

#### 1. Deploy the application

Apply the deployment manifest:

```bash
kubectl apply -f deployments/k8s/deployment.yaml
```

This creates:
- A Deployment with 2 replicas
- Resource limits: 100Mi memory, 200m CPU
- Liveness and readiness probes on `/health`

#### 2. Create the service

Apply the service manifest:

```bash
kubectl apply -f deployments/k8s/service.yaml
```

This creates a ClusterIP service exposing port 80, routing to container port 8080.

#### 3. Apply configuration

Apply the ConfigMap for environment variables:

```bash
kubectl apply -f deployments/k8s/configmap.yaml
```

#### 4. Deploy all at once

You can apply all manifests together:

```bash
kubectl apply -f deployments/k8s/
```

### Verifying Deployment

#### Check pod status

```bash
kubectl get pods -l app=homelab-api
```

Expected output:
```
NAME                           READY   STATUS    RESTARTS   AGE
homelab-api-xxxxxxxxxx-xxxxx   1/1     Running   0          30s
homelab-api-xxxxxxxxxx-xxxxx   1/1     Running   0          30s
```

#### Check deployment

```bash
kubectl get deployment homelab-api
```

#### View logs

```bash
# View logs from all pods
kubectl logs -l app=homelab-api

# View logs from a specific pod
kubectl logs homelab-api-xxxxxxxxxx-xxxxx

# Follow logs
kubectl logs -f -l app=homelab-api
```

#### Check service

```bash
kubectl get service homelab-api
```

#### Describe resources for troubleshooting

```bash
kubectl describe deployment homelab-api
kubectl describe pod -l app=homelab-api
kubectl describe service homelab-api
```

### Accessing the Service

#### Port-forward for testing

```bash
kubectl port-forward service/homelab-api 8080:80
```

Then access:
- Health endpoint: http://localhost:8080/health
- API docs: http://localhost:8080/api/docs/index.html

#### Using an Ingress

Create an Ingress resource to expose the service externally:

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: homelab-api-ingress
  annotations:
    kubernetes.io/ingress.class: nginx
spec:
  rules:
  - host: api.example.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: homelab-api
            port:
              number: 80
```

Apply with:
```bash
kubectl apply -f ingress.yaml
```

## Environment Variables

The application supports the following environment variables:

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `PORT` | HTTP server port | `8080` | No |
| `LOG_LEVEL` | Logging level (debug, info, warn, error) | `info` | No |
| `RATE_LIMIT` | Requests per minute per IP | `100` | No |
| `CORS_ORIGINS` | Allowed CORS origins (comma-separated) | `*` | No |

### Setting Environment Variables

#### In Docker

```bash
docker run -e PORT=8080 -e LOG_LEVEL=debug homelab-api:latest
```

#### In Kubernetes

Environment variables are managed through the ConfigMap (`deployments/k8s/configmap.yaml`):

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: homelab-api-config
data:
  PORT: "8080"
  LOG_LEVEL: "info"
  RATE_LIMIT: "100"
  CORS_ORIGINS: "*"
```

The deployment references this ConfigMap:

```yaml
envFrom:
- configMapRef:
    name: homelab-api-config
```

## Configuration

### Scaling

#### Scale replicas

```bash
kubectl scale deployment homelab-api --replicas=3
```

#### Auto-scaling

Create a Horizontal Pod Autoscaler:

```bash
kubectl autoscale deployment homelab-api \
  --cpu-percent=80 \
  --min=2 \
  --max=10
```

### Resource Limits

Current resource configuration in `deployment.yaml`:

```yaml
resources:
  limits:
    memory: "100Mi"
    cpu: "200m"
  requests:
    memory: "50Mi"
    cpu: "100m"
```

Adjust these based on your workload requirements.

### Health Checks

#### Liveness Probe

Checks if the application is alive (restarts if failing):

```yaml
livenessProbe:
  httpGet:
    path: /health
    port: 8080
  initialDelaySeconds: 10
  periodSeconds: 10
  timeoutSeconds: 5
  failureThreshold: 3
```

#### Readiness Probe

Checks if the application is ready to serve traffic:

```yaml
readinessProbe:
  httpGet:
    path: /health
    port: 8080
  initialDelaySeconds: 5
  periodSeconds: 5
  timeoutSeconds: 3
  failureThreshold: 3
```

## Troubleshooting

### Common Issues

#### Pods not starting

**Problem**: Pods stuck in `Pending` or `ContainerCreating` state.

**Solution**:
```bash
# Check pod events
kubectl describe pod -l app=homelab-api

# Check node resources
kubectl top nodes

# Check pod resource requests
kubectl get pods -o yaml | grep -A 5 resources
```

Common causes:
- Insufficient cluster resources
- Image pull errors
- Volume mount issues

#### Image pull errors

**Problem**: `ErrImagePull` or `ImagePullBackOff` status.

**Solution**:
```bash
# Check image name in deployment
kubectl get deployment homelab-api -o yaml | grep image:

# Verify image exists
docker images | grep homelab-api

# Check pod events for details
kubectl describe pod -l app=homelab-api
```

For private registries, create an image pull secret:
```bash
kubectl create secret docker-registry regcred \
  --docker-server=<registry> \
  --docker-username=<username> \
  --docker-password=<password> \
  --docker-email=<email>
```

#### Health check failures

**Problem**: Pods restarting due to failed liveness probes.

**Solution**:
```bash
# Check pod logs before restart
kubectl logs -l app=homelab-api --previous

# Check health endpoint manually
kubectl port-forward service/homelab-api 8080:80
curl http://localhost:8080/health
```

Adjust probe timing if the application needs more time to start:
```yaml
livenessProbe:
  initialDelaySeconds: 30  # Increase from 10
  timeoutSeconds: 10        # Increase from 5
```

#### Service not accessible

**Problem**: Cannot reach service endpoints.

**Solution**:
```bash
# Verify service exists and has endpoints
kubectl get service homelab-api
kubectl get endpoints homelab-api

# Check if pods are ready
kubectl get pods -l app=homelab-api

# Test service from within cluster
kubectl run -it --rm debug --image=alpine --restart=Never -- sh
# Inside pod: wget -O- http://homelab-api/health

# Check service selector matches pod labels
kubectl get service homelab-api -o yaml | grep selector
kubectl get pods -l app=homelab-api --show-labels
```

#### High memory usage

**Problem**: Pods being OOMKilled (Out of Memory).

**Solution**:
```bash
# Check current memory usage
kubectl top pods -l app=homelab-api

# View pod events
kubectl get events --field-selector involvedObject.name=<pod-name>

# Increase memory limits in deployment.yaml
kubectl edit deployment homelab-api
```

#### Configuration not applied

**Problem**: Environment variables not taking effect.

**Solution**:
```bash
# Verify ConfigMap exists
kubectl get configmap homelab-api-config

# Check ConfigMap data
kubectl describe configmap homelab-api-config

# Verify deployment references ConfigMap
kubectl get deployment homelab-api -o yaml | grep -A 5 envFrom

# Restart pods to pick up changes
kubectl rollout restart deployment homelab-api
```

### Logging

#### View application logs

```bash
# All pods
kubectl logs -l app=homelab-api

# Specific pod
kubectl logs <pod-name>

# Follow logs
kubectl logs -f <pod-name>

# Previous container (after restart)
kubectl logs <pod-name> --previous

# Last N lines
kubectl logs <pod-name> --tail=100

# Since timestamp
kubectl logs <pod-name> --since=1h
```

#### Enable debug logging

Update ConfigMap and restart:

```bash
kubectl patch configmap homelab-api-config -p '{"data":{"LOG_LEVEL":"debug"}}'
kubectl rollout restart deployment homelab-api
```

### Performance Issues

#### Check resource usage

```bash
# Pod metrics
kubectl top pods -l app=homelab-api

# Node metrics
kubectl top nodes
```

#### Analyze bottlenecks

```bash
# Check pod resource limits
kubectl describe pod -l app=homelab-api | grep -A 5 Limits

# View pod events
kubectl get events --sort-by='.lastTimestamp' | grep homelab-api
```

### Rollback

If a deployment fails:

```bash
# View rollout history
kubectl rollout history deployment homelab-api

# Rollback to previous version
kubectl rollout undo deployment homelab-api

# Rollback to specific revision
kubectl rollout undo deployment homelab-api --to-revision=2
```

## Monitoring and Health Checks

### Health Endpoint

The `/health` endpoint returns application health status:

```bash
curl http://localhost:8080/health
```

Response:
```json
{
  "status": "healthy",
  "uptime": "2h 30m 15s"
}
```

### API Documentation

Access Swagger UI for API documentation:

```bash
# Local: http://localhost:8080/api/docs/index.html
# K8s: kubectl port-forward service/homelab-api 8080:80
# Browser: http://localhost:8080/api/docs/index.html
```

### Prometheus Metrics

If you have Prometheus in your cluster, you can add a ServiceMonitor:

```yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: homelab-api
spec:
  selector:
    matchLabels:
      app: homelab-api
  endpoints:
  - port: http
    interval: 30s
```

### Useful Commands

```bash
# Get all resources for the application
kubectl get all -l app=homelab-api

# Check cluster info
kubectl cluster-info

# View API server version
kubectl version

# Check namespace resources
kubectl get all -n default

# Export deployment configuration
kubectl get deployment homelab-api -o yaml > deployment-backup.yaml
```

---

## Quick Reference

### Build and Deploy (Docker)

```bash
# Build
make docker

# Run
docker run -p 8080:8080 homelab-api:latest

# Test
curl http://localhost:8080/health
```

### Deploy to Kubernetes

```bash
# Deploy
kubectl apply -f deployments/k8s/

# Verify
kubectl get pods -l app=homelab-api

# Access
kubectl port-forward service/homelab-api 8080:80

# Test
curl http://localhost:8080/health
```

### Common Operations

```bash
# View logs
kubectl logs -f -l app=homelab-api

# Scale
kubectl scale deployment homelab-api --replicas=3

# Update config
kubectl edit configmap homelab-api-config
kubectl rollout restart deployment homelab-api

# Rollback
kubectl rollout undo deployment homelab-api
```

---

For more information, see:
- [Project README](../README.md)
- [API Documentation](http://localhost:8080/api/docs/index.html)
- [Kubernetes Documentation](https://kubernetes.io/docs/)
- [Docker Documentation](https://docs.docker.com/)
