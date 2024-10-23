package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	pb "github.com/amrimuf/hompimEdu/services/user-service/api/gen/userpb"
	"google.golang.org/grpc"
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
    
    // Convert userID from string to int32
    id, err := strconv.Atoi(userID)
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

    // Create a gRPC request
    req := &pb.UserRequest{UserId: int32(id)} // Use int32
    ctx := context.Background()
    
    // Call the gRPC method
    res, err := usc.client.GetUser(ctx, req)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }
    
    // Return the user data as JSON
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(res)
}

func (usc *UserServiceClient) CreateUser(w http.ResponseWriter, r *http.Request) {
    var user pb.CreateUserRequest
    
    // Decode the request body into the user struct
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    ctx := context.Background()
    
    // Call the gRPC method to create a user
    res, err := usc.client.CreateUser(ctx, &user)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    // Return the created user data as JSON
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(res)
}

func RegisterRoutes(client *UserServiceClient) {
    http.HandleFunc("/users", client.GetUser)
    http.HandleFunc("/users/create", client.CreateUser)
}