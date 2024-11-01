package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/amrimuf/hompimEdu/services/course-service/api/gen/coursepb"
	"github.com/amrimuf/hompimEdu/services/course-service/internal/service"
	"google.golang.org/grpc"
	"github.com/gorilla/mux"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
	"github.com/amrimuf/hompimEdu/services/course-service/api"
	"google.golang.org/grpc/reflection"
)

func main() {
    // Get configuration from environment variables
    grpcPort := os.Getenv("GRPC_PORT")
    if grpcPort == "" {
        grpcPort = "50052" // Default port
    }

    // Get user service connection details from environment
    userServiceHost := os.Getenv("USER_SERVICE_HOST")
    if userServiceHost == "" {
        userServiceHost = "user-service" // Default for docker/k8s
    }

    userServicePort := os.Getenv("USER_SERVICE_GRPC_PORT")
    if userServicePort == "" {
        userServicePort = "50051" // Default port
    }

    // Initialize user service client with environment-based configuration
    userServiceClient, err := service.NewUserServiceClient(
        fmt.Sprintf("%s:%s", userServiceHost, userServicePort),
    )
    if err != nil {
        log.Fatalf("Failed to connect to user service: %v", err)
    }
    // Call GetUsers (consider making this optional or part of health check)
    userServiceClient.CallGetUsers()
    if err != nil {
        log.Printf("Warning: Failed to call GetUsers: %v", err)
        // Consider whether this should be fatal or just a warning
    }

    // Start gRPC server
    go func() {
        lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
        if err != nil {
            log.Fatalf("failed to listen: %v", err)
        }

        grpcServer := grpc.NewServer()
        courseService := service.NewCourseServiceServer()
        coursepb.RegisterCourseServiceServer(grpcServer, courseService)
        
        reflection.Register(grpcServer)

        log.Printf("gRPC server listening on port %s", grpcPort)
        if err := grpcServer.Serve(lis); err != nil {
            log.Fatalf("failed to serve: %v", err)
        }
    }()

    // HTTP server
    router := mux.NewRouter()
    
    // Create gRPC client for internal communication
    conn, err := grpc.Dial(
        fmt.Sprintf("localhost:%s", grpcPort),
        grpc.WithTransportCredentials(insecure.NewCredentials()),
    )
    if err != nil {
        log.Fatalf("failed to dial: %v", err)
    }
    defer conn.Close()
    
    grpcClient := coursepb.NewCourseServiceClient(conn)
    api.RegisterRoutes(router, grpcClient)

    httpPort := os.Getenv("HTTP_PORT")
    if httpPort == "" {
        httpPort = "8082"
    }

    log.Printf("HTTP server listening on port %s", httpPort)
    if err := http.ListenAndServe(fmt.Sprintf(":%s", httpPort), router); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}
