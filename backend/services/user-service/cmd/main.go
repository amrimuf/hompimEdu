package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/amrimuf/hompimEdu/services/user-service/api"
	pb "github.com/amrimuf/hompimEdu/services/user-service/api/gen/userpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Define your server that implements the gRPC service
type server struct {
    pb.UnimplementedUserServiceServer
}

func (s *server) ListUsers(ctx context.Context, req *emptypb.Empty) (*pb.ListUsersResponse, error) {
    // Example user data
    users := []*pb.User{
        {Id: "1", Name: "Alice"},
        {Id: "2", Name: "Bob"},
    }

    return &pb.ListUsersResponse{Users: users}, nil
}

func main() {
    // gRPC server
    go func() {
        lis, err := net.Listen("tcp", ":50051")
        if err != nil {
            log.Fatalf("failed to listen: %v", err)
        }
        s := grpc.NewServer()
        pb.RegisterUserServiceServer(s, &server{}) // Register your server implementation
        reflection.Register(s)
        log.Printf("gRPC server listening on %v", lis.Addr())
        if err := s.Serve(lis); err != nil {
            log.Fatalf("failed to serve: %v", err)
        }
    }()

    // HTTP server
    api.RegisterRoutes()
    log.Fatal(http.ListenAndServe(":8080", nil))
}
