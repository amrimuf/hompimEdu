package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/amrimuf/hompimEdu/services/course-service/api"
	"github.com/amrimuf/hompimEdu/services/course-service/api/gen/coursepb"
	"github.com/amrimuf/hompimEdu/services/course-service/internal/service"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"github.com/amrimuf/hompimEdu/services/course-service/api/gen/userpb"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	// Get configuration from environment variables
	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "50052" // Default port
	}

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8082" // Default port
	}

	// Initialize database connection
	db, err := service.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize user service client with retry
	userServiceHost := os.Getenv("USER_SERVICE_HOST")
	userServicePort := os.Getenv("USER_SERVICE_GRPC_PORT")
	if userServiceHost == "" || userServicePort == "" {
		log.Fatal("USER_SERVICE_HOST and USER_SERVICE_GRPC_PORT must be set")
	}

	var userConn *grpc.ClientConn
	var userClient userpb.UserServiceClient
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		log.Printf("Attempting to connect to user service (attempt %d/%d)", i+1, maxRetries)
		var err error
		userConn, err = grpc.Dial(
			fmt.Sprintf("%s:%s", userServiceHost, userServicePort),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			log.Printf("Failed to connect: %v", err)
			time.Sleep(time.Second * 5)
			continue
		}

		userClient = userpb.NewUserServiceClient(userConn)
		
		// Test the connection
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		_, err = userClient.GetUser(ctx, &userpb.UserRequest{UserId: 1})
		cancel()

		if err != nil && status.Code(err) != codes.NotFound {
			// If error is not "NotFound", it's a real error
			log.Printf("Failed to call user service: %v", err)
			userConn.Close()
			time.Sleep(time.Second * 5)
			continue
		}

		log.Printf("Successfully connected to user service")
		break
	}

	if userClient == nil {
		log.Fatal("Failed to connect to user service after maximum retries")
	}

	// Don't close the connection until the program exits
	defer userConn.Close()

	// Start both servers
	go func() { startGRPCServer(db, userClient, grpcPort) }()
	startHTTPServer(grpcPort, httpPort)
}

func startGRPCServer(db *sql.DB, userClient userpb.UserServiceClient, port string) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	courseService := service.NewCourseServiceServer(db, userClient)
	coursepb.RegisterCourseServiceServer(grpcServer, courseService)
	reflection.Register(grpcServer)

	log.Printf("gRPC server listening on port %s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func startHTTPServer(grpcPort, httpPort string) {
	// Use environment-based service discovery
	grpcHost := os.Getenv("GRPC_HOST")
	if grpcHost == "" {
		grpcHost = "localhost" // Default for local development
	}

	conn, err := grpc.Dial(
		fmt.Sprintf("%s:%s", grpcHost, grpcPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()

	router := mux.NewRouter()
	grpcClient := coursepb.NewCourseServiceClient(conn)
	api.RegisterRoutes(router, grpcClient)

	log.Printf("HTTP server listening on port %s", httpPort)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", httpPort), router); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
