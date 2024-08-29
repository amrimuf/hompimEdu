package main

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/amrimuf/hompimEdu/services/course-service/api/gen/coursepb"
	pb "github.com/amrimuf/hompimEdu/services/course-service/api/gen/userpb"
	"github.com/amrimuf/hompimEdu/services/course-service/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

func callUserService() {
    maxRetries := 5
    for i := 0; i < maxRetries; i++ {
        conn, err := grpc.NewClient("user-service:50051", grpc.WithInsecure())
        if err != nil {
            log.Printf("Failed to connect: %v. Retrying in 5 seconds...", err)
            time.Sleep(5 * time.Second)
            continue
        }
        defer conn.Close()

        client := pb.NewUserServiceClient(conn)

        // Call ListUsers method
        resp, err := client.ListUsers(context.Background(), &emptypb.Empty{})
        if err != nil {
            log.Printf("Could not list users: %v. Retrying in 5 seconds...", err)
            time.Sleep(5 * time.Second)
            continue
        }

        for _, user := range resp.Users {
            log.Printf("User: %s, Name: %s", user.Id, user.Name)
        }
        return
    }
    log.Fatalf("Failed to connect to user-service after %d attempts", maxRetries)
}


func main() {
    callUserService()

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
