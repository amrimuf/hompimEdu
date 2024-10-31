package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/amrimuf/hompimEdu/services/user-service/api"
	"github.com/amrimuf/hompimEdu/services/user-service/api/gen/userpb"
	"github.com/amrimuf/hompimEdu/services/user-service/internal"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

func main() {
    // Get ports from environment
    grpcPort := os.Getenv("GRPC_PORT")
    if grpcPort == "" {
        grpcPort = "50051" // Default port
    }

    httpPort := os.Getenv("HTTP_PORT")
    if httpPort == "" {
        httpPort = "8083" // Default port
    }

    db, err := internal.InitDB()
    if err != nil {
        log.Fatalf("failed to connect to the database: %v", err)
    }
    defer db.Close()

    // Start both servers and wait for them to finish
    go startGRPCServer(db, grpcPort)
    startHTTPServer(grpcPort, httpPort)
}

func startGRPCServer(db *sql.DB, port string) {
    lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }

    server := internal.NewServer(db)

    grpcServer := grpc.NewServer()
    userpb.RegisterUserServiceServer(grpcServer, server)
    reflection.Register(grpcServer)

    log.Printf("gRPC server listening on port %s", port)
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("failed to serve gRPC server: %v", err)
    }
}

func startHTTPServer(grpcPort, httpPort string) {
    // Use environment-based service discovery
    grpcHost := os.Getenv("GRPC_HOST")
    if grpcHost == "" {
        grpcHost = "localhost" // Default for local development
    }

    grpcConn, err := grpc.Dial(
        fmt.Sprintf("%s:%s", grpcHost, grpcPort),
        grpc.WithTransportCredentials(insecure.NewCredentials()),
    )
    if err != nil {
        log.Fatalf("failed to dial gRPC server: %v", err)
    }
    defer grpcConn.Close()

    userServiceClient := api.NewUserServiceClient(grpcConn)

    api.RegisterRoutes(userServiceClient)

    log.Printf("HTTP server listening on port %s", httpPort)
    log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", httpPort), nil))
}
