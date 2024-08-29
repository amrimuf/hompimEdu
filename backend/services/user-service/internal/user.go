package internal

import (
	"context"

	pb "github.com/amrimuf/hompimEdu/services/user-service/api/gen/userpb"
)

type Server struct {
    pb.UnimplementedUserServiceServer
}

func (s *Server) GetUser(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {
    return &pb.UserResponse{Id: req.UserId, Name: "John Doe", Email: "john@example.com"}, nil
}

func (s *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
    return &pb.CreateUserResponse{Id: 1, Name: req.Name, Email: req.Email}, nil
}
