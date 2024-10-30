package api

import (
	"context"
	"encoding/json"
	"net/http"

	pb "github.com/amrimuf/hompimEdu/services/auth-service/api/gen/authpb"
	"google.golang.org/grpc"
)

type AuthServiceClient struct {
	client pb.AuthServiceClient
}

func NewAuthServiceClient(conn *grpc.ClientConn) *AuthServiceClient {
	return &AuthServiceClient{
		client: pb.NewAuthServiceClient(conn),
	}
}

func (asc *AuthServiceClient) Register(w http.ResponseWriter, r *http.Request) {
	var req pb.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	res, err := asc.client.Register(ctx, &req)
	if err != nil {
		http.Error(w, "Registration failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (asc *AuthServiceClient) Login(w http.ResponseWriter, r *http.Request) {
	var req pb.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	res, err := asc.client.Login(ctx, &req)
	if err != nil {
		http.Error(w, "Login failed: "+err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
