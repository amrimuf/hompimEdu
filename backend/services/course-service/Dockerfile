FROM golang:1.21.0 AS builder

WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Copy the entire project structure
COPY . .

# Download dependencies
RUN go mod download

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o course-service ./cmd

# Step 2: Use a minimal image to run the binary
FROM gcr.io/distroless/static-debian10

WORKDIR /app

# Copy the statically linked binary from the builder stage
COPY --from=builder /app/course-service .

# Command to run the executable
CMD ["/app/course-service"]