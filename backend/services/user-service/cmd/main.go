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
    startGRPCServer()
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
    go func() {
        if err := grpcServer.Serve(lis); err != nil {
            log.Fatalf("failed to serve gRPC server: %v", err)
        }
    }()
}

func startHTTPServer() {
    grpcConn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        log.Fatalf("failed to dial gRPC server: %v", err)
    }

    userServiceClient := api.NewUserServiceClient(grpcConn)
    api.RegisterRoutes(userServiceClient)

    log.Fatal(http.ListenAndServe(":8080", nil))
}