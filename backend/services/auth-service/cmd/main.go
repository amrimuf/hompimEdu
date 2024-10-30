// main.go
package main

import (
	"log"
	"net"
	"net/http"

	pb "github.com/amrimuf/hompimEdu/services/auth-service/api/gen/authpb"
	api "github.com/amrimuf/hompimEdu/services/auth-service/api/handlers"
	"github.com/amrimuf/hompimEdu/services/auth-service/internal"
	services "github.com/amrimuf/hompimEdu/services/auth-service/internal/services"
	"google.golang.org/grpc"
)

func main() {
    // Initialize the database connection
	db, err := internal.InitDB()
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}
	defer db.Close()

    lis, err := net.Listen("tcp", ":50053")
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }

    grpcServer := grpc.NewServer()
    authService := services.NewAuthService(db) // Pass the db connection to your service
    pb.RegisterAuthServiceServer(grpcServer, authService)

    // Start the gRPC server in a separate goroutine
    go func() {
        log.Println("gRPC server listening on port 50053...")
        if err := grpcServer.Serve(lis); err != nil {
            log.Fatalf("failed to serve: %v", err)
        }
    }()

    // HTTP server
    conn, err := grpc.Dial("localhost:50053", grpc.WithInsecure())
    if err != nil {
        log.Fatalf("failed to connect to gRPC server: %v", err)
    }
    authClient := api.NewAuthServiceClient(conn)

    http.HandleFunc("/register", authClient.Register)
    http.HandleFunc("/login", authClient.Login)

    log.Fatal(http.ListenAndServe(":8084", nil))
}
