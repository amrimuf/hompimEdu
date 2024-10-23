package main

import (
	"log"
	"net"

	"github.com/amrimuf/hompimEdu/services/course-service/api/gen/coursepb"
	"github.com/amrimuf/hompimEdu/services/course-service/internal/service"
	"google.golang.org/grpc"
)


func main() {
    // Initialize user service client
    userServiceClient, err := service.NewUserServiceClient("user-service:50051")
    if err != nil {
        log.Fatalf("Failed to connect to user service: %v", err)
    }

    // Call ListUsers
    userServiceClient.CallListUsers()

    lis, err := net.Listen("tcp", ":50052")
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }

    grpcServer := grpc.NewServer()
    courseService := service.NewCourseServiceServer()

    coursepb.RegisterCourseServiceServer(grpcServer, courseService)

    log.Println("Course service listening on port 50052...")
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}
