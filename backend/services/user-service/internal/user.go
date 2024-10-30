package internal

import (
	"context"
	"database/sql"
	"errors"
	"log"

	pb "github.com/amrimuf/hompimEdu/services/user-service/api/gen/userpb"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Server struct {
	pb.UnimplementedUserServiceServer
	db *sql.DB // Add DB connection for actual user data
}

func (s *Server) AuthenticateJWT(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errors.New("missing metadata")
	}

	// Log the incoming metadata for debugging
	log.Printf("Incoming metadata: %v", md)

	if len(md["authorization"]) == 0 {
		return "", errors.New("missing token")
	}

	// Expect the token to be in the format "Bearer <token>"
	token := md["authorization"][0]
	if len(token) < 7 || token[:7] != "Bearer " {
		return "", errors.New("invalid token format")
	}

	// Strip the "Bearer " prefix to get the token
	token = token[7:]

	username, err := ValidateJWT(token)
	if err != nil {
		return "", err
	}
	return username, nil
}



func (s *Server) GetUser(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {
	// Here you would normally fetch user details from a database
	return &pb.UserResponse{Id: req.UserId, Name: "John Doe", Email: "john@example.com"}, nil
}

func (s *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	// Hash the password before saving it
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Here you would normally save the user details to the database, including the hashed password
	_, err = s.db.Exec("INSERT INTO users (name, email, password) VALUES ($1, $2, $3)", req.Name, req.Email, hashedPassword)
	if err != nil {
		return nil, err
	}

	// Assume the ID is generated after insertion (you should fetch it accordingly if needed)
	return &pb.CreateUserResponse{Id: 1, Name: req.Name, Email: req.Email}, nil
}


func (s *Server) ListUsers(ctx context.Context, req *emptypb.Empty) (*pb.ListUsersResponse, error) {
    // Log incoming context and metadata
    // md, _ := metadata.FromIncomingContext(ctx)
    // log.Printf("Incoming metadata: %v", md)

    // // Authenticate JWT
    // _, err := s.AuthenticateJWT(ctx)
    // if err != nil {
    //     return nil, err
    // }

    // Example user data - replace this with actual DB query
    users := []*pb.User{
        {Id: "1", Name: "Alice"},
        {Id: "2", Name: "Bob"},
    }

    return &pb.ListUsersResponse{Users: users}, nil
}
