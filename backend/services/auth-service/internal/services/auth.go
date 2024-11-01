package services

import (
	"context"
	"database/sql"
	"errors"
	"sync"

	pb "github.com/amrimuf/hompimEdu/services/auth-service/api/gen/authpb"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
    pb.UnimplementedAuthServiceServer
    mu    sync.Mutex
    db    *sql.DB // Database connection
}

func NewAuthService(db *sql.DB) *AuthService {
    return &AuthService{
        db: db,
    }
}

func (s *AuthService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
    s.mu.Lock()
    defer s.mu.Unlock()

    // Check if the user already exists
    var exists bool
    err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username=$1)", req.Username).Scan(&exists)
    if err != nil {
        return nil, err
    }

    if exists {
        return nil, errors.New("user already exists")
    }

    // Hash the password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        return nil, err
    }

    // Insert the new user into the database
    _, err = s.db.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", req.Username, hashedPassword)
    if err != nil {
        return nil, err
    }

    return &pb.RegisterResponse{UserId: req.Username}, nil
}

func (s *AuthService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
    s.mu.Lock()
    defer s.mu.Unlock()

    var (
        hashedPassword string
        userID        int64
    )
    err := s.db.QueryRow("SELECT id, password FROM users WHERE username=$1", req.Username).Scan(&userID, &hashedPassword)
    if err != nil {
        return nil, errors.New("invalid username or password")
    }

    // Check the password against the hashed password
    err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password))
    if err != nil {
        return nil, errors.New("invalid username or password")
    }

    // Generate JWT token
    token, err := GenerateJWT(req.Username, userID)
    if err != nil {
        return nil, err
    }

    return &pb.LoginResponse{Token: token}, nil
}
