# Project Setup and Deployment Guide

## Install Prerequisites

### Install Protocol Buffers Compiler (`protoc`)

1. **macOS**: Use Homebrew to install `protoc`.

    ```sh
    brew install protobuf
    ```

2. **Ubuntu/Debian**: Install using `apt-get`.

    ```sh
    sudo apt-get update
    sudo apt-get install -y protobuf-compiler
    ```

3. **Windows**: Download the latest release from the [Protobuf releases page](https://github.com/protocolbuffers/protobuf/releases) and add it to your system PATH.

### Install Go Protobuf Plugins

1. Install the Go Protobuf plugin:

    ```sh
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
    ```

2. Install the Go gRPC plugin:

    ```sh
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
    ```

    Ensure the Go binaries are in the PATH. For example, add the following to the `.bashrc` or `.zshrc`:

    ```sh
    export PATH="$PATH:$(go env GOPATH)/bin"
    ```

## Generate Protobuf Code

Run the following commands to generate Go code from the Protobuf definitions:

```sh
protoc --go_out=services/user-service/api/gen --go-grpc_out=services/user-service/api/gen proto/user.proto
protoc --go_out=services/course-service/api/gen --go-grpc_out=services/course-service/api/gen proto/course.proto
```

## Option 1: Deploying with Kubernetes

### Build Docker Images

Build the Docker image for the services:

```sh
docker build -t amrimuf/<service-name>:latest .
```

## Push the Updated Image

```sh
docker push amrimuf/<service-name>:latest
```

```sh
kubectl rollout restart deployment/user-service
```

### Create Docker Registry Secret

Create a Kubernetes secret for accessing the Docker registry:

```sh
kubectl create secret docker-registry my-registry-secret \
  --docker-server=<the_REGISTRY_SERVER> \
  --docker-username=<the_USERNAME> \
  --docker-password=<the_PASSWORD> \
  --docker-email=<the_EMAIL>
```

### Apply Kubernetes Manifests

Deploy the services and deployments to Kubernetes:

```sh
kubectl apply -f deploy/k8s/base/deployments/
kubectl apply -f deploy/k8s/base/services/
```

### Verify Deployment

Check the status of the Pods, Deployments, Services, and Endpoints:

```sh
kubectl get pods
kubectl get deployments
kubectl get services
kubectl get endpoints
```

### Port Forward to Access HTTP Service

If your service exposes an HTTP endpoint on port `8083`, you can access it using `kubectl port-forward`. This maps the service port in the Kubernetes cluster to your local machine.

1. Forward the HTTP port (e.g., `8083`):

    ```sh
    kubectl port-forward svc/user-service 8083:8083
    ```

    This command forwards the port `8083` from the `user-service` in the Kubernetes cluster to your local machine on the same port.

2. Verify the service is accessible by sending a request using `curl`:

    ```sh
    curl http://localhost:8083/users
    ```

    If everything is working, you should see a response with the list of users from the service.

## Option 2: Deploying with Docker Compose

### Build and Start Services

#### Development Environment

To start your services in a development environment, use the docker-compose.dev.yml file:

```sh
docker-compose -f docker-compose.dev.yml up --build
```

#### Production Environment

To start your services in a production environment, use the docker-compose.prod.yml file:

```sh
docker-compose -f docker-compose.prod.yml up --build
```

### Stopping the Services

To stop the running services, use:

```sh
docker-compose down
```

### Cleanup Docker Images

If you need to remove all stopped services and their images, you can use:

```sh
docker-compose down --rmi all
```

---

## Troubleshooting

To view logs and get details for a specific Pod:

```sh
kubectl logs <pod-name> --previous
kubectl describe pod <pod-name>
```

## Cleanup

Remove all deployments, services, and pods:

```sh
kubectl delete deployments --all
kubectl delete services --all
kubectl delete pods --all
```

Alternatively, delete all resources in the base directory:

```sh
kubectl delete -f deploy/k8s/base/
```

Remove all resources in the namespace:

```sh
kubectl delete all --all
```

---

Replace `<the_REGISTRY_SERVER>`, `<the_USERNAME>`, `<the_PASSWORD>`, and `<the_EMAIL>` with the actual Docker registry credentials.

## Project Structure

### Current Structure
```
backend/
├── gateway/
│   ├── middleware/
│   ├── routes/
│   ├── main.go
│   └── Dockerfile
├── services/
│   ├── auth-service/
│   │   ├── api/
│   │   │   ├── gen/
│   │   │   │   └── authpb/
│   │   │   └── handlers/
│   │   ├── cmd/
│   │   │   └── main.go
│   │   ├── internal/
│   │   │   └── services/
│   │   ├── Dockerfile
│   │   ├── go.mod
│   │   └── go.sum
│   ├── user-service/
│   │   ├── api/
│   │   │   ├── gen/
│   │   │   │   └── userpb/
│   │   │   └── handlers/
│   │   ├── cmd/
│   │   │   └── main.go
│   │   ├── internal/
│   │   ├── Dockerfile
│   │   ├── go.mod
│   │   └── go.sum
│   └── course-service/
│       ├── api/
│       │   ├── gen/
│       │   │   ├── coursepb/
│       │   │   └── userpb/
│       │   └── handlers/
│       ├── cmd/
│       │   └── main.go
│       ├── internal/
│       │   └── service/
│       ├── Dockerfile
│       ├── go.mod
│       └── go.sum
├── proto/
│   ├── auth.proto
│   ├── user.proto
│   └── course.proto
└── README.md
```

### Recommended Structure
```
backend/
├── config/                    # Centralized configuration
│   ├── development.yaml
│   └── production.yaml
├── deploy/                    # Deployment configurations
│   ├── kubernetes/
│   │   ├── base/
│   │   └── overlays/
│   └── docker-compose/
├── docs/                      # Documentation
│   ├── api/
│   │   ├── swagger/
│   │   └── postman/
│   └── architecture/
│       ├── diagrams/
│       └── decisions/
├── gateway/
│   ├── middleware/
│   ├── routes/
│   ├── main.go
│   └── Dockerfile
├── pkg/                       # Shared packages
│   ├── logger/
│   │   └── zap.go
│   ├── database/
│   │   └── postgres.go
│   └── middleware/
│       └── auth.go
├── proto/                     # Protocol Buffers
│   ├── auth.proto
│   ├── user.proto
│   └── course.proto
├── scripts/                   # Utility scripts
│   ├── build.sh
│   ├── test.sh
│   └── deploy.sh
├── services/                  # Microservices
│   ├── auth-service/
│   │   ├── api/
│   │   ├── cmd/
│   │   ├── internal/
│   │   ├── tests/            # Tests directory
│   │   │   ├── unit/
│   │   │   ├── integration/
│   │   │   └── e2e/
│   │   ├── metrics/          # Monitoring
│   │   │   └── prometheus/
│   │   ├── health/           # Health checks
│   │   │   ├── liveness.go
│   │   │   └── readiness.go
│   │   └── Dockerfile
│   ├── user-service/
│   │   └── [similar structure]
│   └── course-service/
│       └── [similar structure]
└── README.md
```
