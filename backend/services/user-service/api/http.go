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

func (usc *UserServiceClient) GetUser(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("id")
	id, err := strconv.Atoi(userID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	req := &pb.UserRequest{UserId: int32(id)}
	ctx := context.Background()
	res, err := usc.client.GetUser(ctx, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

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

// New method to list users
func (usc *UserServiceClient) ListUsers(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	res, err := usc.client.ListUsers(ctx, &emptypb.Empty{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func RegisterRoutes(client *UserServiceClient) {
	http.HandleFunc("/users", client.GetUser)
	http.HandleFunc("/users/create", client.CreateUser)
	http.HandleFunc("/users/list", client.ListUsers) // Register ListUsers route
}
