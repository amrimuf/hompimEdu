package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	pb "github.com/amrimuf/hompimEdu/services/user-service/api/gen/userpb"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserServiceClient struct {
	client pb.UserServiceClient
}

func NewUserServiceClient(conn *grpc.ClientConn) *UserServiceClient {
	return &UserServiceClient{
		client: pb.NewUserServiceClient(conn),
	}
}

// GetUser retrieves a user by ID
func (usc *UserServiceClient) GetUser(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("id") // Expecting ?id=123
	id, err := strconv.Atoi(userID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	req := &pb.UserRequest{UserId: int64(id)}
	ctx := context.Background()
	res, err := usc.client.GetUser(ctx, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

// CreateUser creates a new user
func (usc *UserServiceClient) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user pb.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	res, err := usc.client.CreateUser(ctx, &user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

// GetUsers retrieves a list of users
func (usc *UserServiceClient) GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	res, err := usc.client.GetUsers(ctx, &emptypb.Empty{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

// RegisterRoutes registers the HTTP routes
func RegisterRoutes(client *UserServiceClient) {
	http.HandleFunc("/users", client.GetUsers)           // Route for getting all users
	http.HandleFunc("/users/create", client.CreateUser)   // Route for creating a new user
	http.HandleFunc("/user", client.GetUser)              // Route for getting a specific user by ID
}
