package main

import (
	"database/sql"
	"log"
	"net"
	"net/http"

	"github.com/amrimuf/hompimEdu/services/user-service/api"
	"github.com/amrimuf/hompimEdu/services/user-service/api/gen/userpb"
	"github.com/amrimuf/hompimEdu/services/user-service/internal"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

func main() {
    db, err := internal.InitDB()
    if err != nil {
        log.Fatalf("failed to connect to the database: %v", err)
    }
    defer db.Close()

    // Start both servers and wait for them to finish
    go startGRPCServer(db)
    startHTTPServer()
}

func startGRPCServer(db *sql.DB) {
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }

    server := internal.NewServer(db) // Use NewServer to initialize

    grpcServer := grpc.NewServer()
    userpb.RegisterUserServiceServer(grpcServer, server)
    reflection.Register(grpcServer)

    log.Printf("gRPC server listening on %v", lis.Addr())
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("failed to serve gRPC server: %v", err)
    }
}

func startHTTPServer() {
    grpcConn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        log.Fatalf("failed to dial gRPC server: %v", err)
    }
    defer grpcConn.Close()

    userServiceClient := api.NewUserServiceClient(grpcConn)

    api.RegisterRoutes(userServiceClient)

    log.Println("HTTP server listening on :8083")
    log.Fatal(http.ListenAndServe(":8083", nil))
}
