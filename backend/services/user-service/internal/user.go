package internal

import (
	"context"
	"database/sql"
	"log"

	pb "github.com/amrimuf/hompimEdu/services/user-service/api/gen/userpb"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Server represents the gRPC server and holds the database connection.
type Server struct {
	pb.UnimplementedUserServiceServer
	db *sql.DB // Database connection for user data
}

// NewServer creates a new Server with a database connection.
func NewServer(db *sql.DB) *Server {
	return &Server{db: db}
}

// GetUser fetches a user by their ID from the database.
func (s *Server) GetUser(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {
    var user pb.UserResponse
    err := s.db.QueryRowContext(ctx, `
        SELECT id, username, email 
        FROM users 
        WHERE id = $1`, req.UserId).Scan(&user.Id, &user.Username, &user.Email)
    
    if err == sql.ErrNoRows {
        return nil, status.Errorf(codes.NotFound, "user not found")
    }
    if err != nil {
        return nil, status.Errorf(codes.Internal, "failed to fetch user: %v", err)
    }
    
    return &user, nil
}

// CreateUser creates a new user and saves it to the database.
func (s *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	// Hash the password before saving it
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Insert the user details into the database, including the hashed password
	var userID int64
	err = s.db.QueryRow("INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id", req.Username, req.Email, hashedPassword).Scan(&userID)
	if err != nil {
		log.Printf("Error inserting user: %v", err)
		return nil, err // Handle other database errors
	}

	return &pb.CreateUserResponse{Id: userID, Username: req.Username, Email: req.Email}, nil
}

// GetUsers retrieves a list of all users from the database.
func (s *Server) GetUsers(ctx context.Context, req *emptypb.Empty) (*pb.GetUsersResponse, error) {
	// Query the database for all users
	rows, err := s.db.Query("SELECT id, username FROM users")
	if err != nil {
		log.Printf("Error fetching users: %v", err)
		return nil, err // Handle database errors
	}
	defer rows.Close()

	var users []*pb.User
	for rows.Next() {
		var user pb.User
		if err := rows.Scan(&user.Id, &user.Username); err != nil {
			log.Printf("Error scanning user: %v", err)
			return nil, err // Handle scanning errors
		}
		users = append(users, &user)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		log.Printf("Error during rows iteration: %v", err)
		return nil, err // Handle iteration errors
	}

	return &pb.GetUsersResponse{Users: users}, nil // Return updated response type
}
