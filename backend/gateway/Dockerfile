# Use a minimal base image for production
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files for dependency installation
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the API Gateway binary
RUN go build -o gateway .

# Start a new, lightweight image for running the app
FROM alpine:latest

# Copy the compiled binary from the builder
COPY --from=builder /app/gateway /gateway

# Set environment variables for the services
ENV USER_SERVICE_URL="http://user-service:8083"
ENV AUTH_SERVICE_URL="http://auth-service:8084"
ENV API_GATEWAY_PORT=8085

# Expose the API Gateway port
EXPOSE 8085

# Command to run the API Gateway
CMD ["/gateway"]
