package main

import (
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
    // Start both servers and wait for them to finish
    go startGRPCServer()
    startHTTPServer()
}

func startGRPCServer() {
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }

    grpcServer := grpc.NewServer()
    userpb.RegisterUserServiceServer(grpcServer, &internal.Server{})
    reflection.Register(grpcServer)

    log.Printf("gRPC server listening on %v", lis.Addr())
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("failed to serve gRPC server: %v", err)
    }
}

func startHTTPServer() {
    // Correctly use grpc.Dial to create a client connection
    grpcConn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        log.Fatalf("failed to dial gRPC server: %v", err)
    }
    defer grpcConn.Close() // Ensure the connection is closed when done

    userServiceClient := api.NewUserServiceClient(grpcConn)

    // Ensure RegisterRoutes attaches the routes to the default HTTP mux
    api.RegisterRoutes(userServiceClient)

    log.Println("HTTP server listening on :8083")
    log.Fatal(http.ListenAndServe(":8083", nil))
}
