package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/amrimuf/hompimEdu/services/course-service/api"
	"github.com/amrimuf/hompimEdu/services/course-service/api/gen/coursepb"
	"github.com/amrimuf/hompimEdu/services/course-service/internal/service"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
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
	// Start both servers
	go func() { startGRPCServer(db, grpcPort) }()
	startHTTPServer(grpcPort, httpPort)
}

func startGRPCServer(db *sql.DB, port string) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	courseService := service.NewCourseServiceServer(db)
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
