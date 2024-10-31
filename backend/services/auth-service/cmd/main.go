// main.go
package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	pb "github.com/amrimuf/hompimEdu/services/auth-service/api/gen/authpb"
	api "github.com/amrimuf/hompimEdu/services/auth-service/api/handlers"
	"github.com/amrimuf/hompimEdu/services/auth-service/internal"
	services "github.com/amrimuf/hompimEdu/services/auth-service/internal/services"
	"google.golang.org/grpc"
)

func main() {
    // Get ports from environment variables with defaults
    grpcPort := os.Getenv("GRPC_PORT")
    if grpcPort == "" {
        grpcPort = "50053" // Default port
    }

    httpPort := os.Getenv("HTTP_PORT")
    if httpPort == "" {
        httpPort = "8084" // Default port
    }

    // Initialize the database connection
    db, err := internal.InitDB()
    if err != nil {
        log.Fatalf("failed to connect to the database: %v", err)
    }
    defer db.Close()

    // Start gRPC server
    lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }

    grpcServer := grpc.NewServer()
    authService := services.NewAuthService(db)
    pb.RegisterAuthServiceServer(grpcServer, authService)

    // Start the gRPC server in a separate goroutine
    go func() {
        log.Printf("gRPC server listening on port %s...", grpcPort)
        if err := grpcServer.Serve(lis); err != nil {
            log.Fatalf("failed to serve: %v", err)
        }
    }()

    // HTTP server
    // Use environment-based service discovery
    grpcHost := os.Getenv("GRPC_HOST")
    if grpcHost == "" {
        grpcHost = "localhost" // Default for local development
    }
    
    conn, err := grpc.Dial(
        fmt.Sprintf("%s:%s", grpcHost, grpcPort),
        grpc.WithInsecure(),
    )
    if err != nil {
        log.Fatalf("failed to connect to gRPC server: %v", err)
    }
    authClient := api.NewAuthServiceClient(conn)

    http.HandleFunc("/register", authClient.Register)
    http.HandleFunc("/login", authClient.Login)

    log.Printf("HTTP server listening on port %s...", httpPort)
    log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", httpPort), nil))
}
