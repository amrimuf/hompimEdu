package service

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/amrimuf/hompimEdu/services/course-service/api/gen/coursepb"
	"github.com/amrimuf/hompimEdu/services/course-service/api/gen/userpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

// CourseServiceServer implements the CourseService gRPC interface
type CourseServiceServer struct {
    coursepb.UnimplementedCourseServiceServer
    courses map[int32]*coursepb.Course // In-memory store for courses
    nextID  int32                      // For auto-incrementing IDs
}

type UserServiceClient struct {
    client userpb.UserServiceClient
}

// NewCourseServiceServer creates a new CourseServiceServer
func NewCourseServiceServer() *CourseServiceServer {
    return &CourseServiceServer{
        courses: make(map[int32]*coursepb.Course),
        nextID:  1,
    }
}

// GetCourse retrieves a course by ID
func (s *CourseServiceServer) GetCourse(ctx context.Context, req *coursepb.GetCourseRequest) (*coursepb.GetCourseResponse, error) {
    course, exists := s.courses[req.Id]
    if !exists {
        return nil, errors.New("course not found")
    }
    return &coursepb.GetCourseResponse{Course: course}, nil
}

// CreateCourse creates a new course
func (s *CourseServiceServer) CreateCourse(ctx context.Context, req *coursepb.CreateCourseRequest) (*coursepb.CreateCourseResponse, error) {
    course := &coursepb.Course{
        Id:             s.nextID,
        Title:          req.Title,
        Description:    req.Description,
        Duration:       req.Duration,
        EnrollmentType: req.EnrollmentType,
        MentorId:       req.MentorId,
        CreatedAt:      time.Now().Format(time.RFC3339),
        UpdatedAt:      time.Now().Format(time.RFC3339),
    }
    s.courses[s.nextID] = course
    s.nextID++
    return &coursepb.CreateCourseResponse{Course: course}, nil
}

// ListCourses lists all available courses
func (s *CourseServiceServer) ListCourses(ctx context.Context, req *coursepb.ListCoursesRequest) (*coursepb.ListCoursesResponse, error) {
    var courses []*coursepb.Course
    for _, course := range s.courses {
        courses = append(courses, course)
    }
    return &coursepb.ListCoursesResponse{Courses: courses}, nil
}

// NewUserServiceClient initializes a new UserServiceClient
func NewUserServiceClient(address string) (*UserServiceClient, error) {
    conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        return nil, err
    }

    return &UserServiceClient{
        client: userpb.NewUserServiceClient(conn),
    }, nil
}

// CallGetUsers makes a call to the GetUsers method of the UserService
func (usc *UserServiceClient) CallGetUsers() {
    maxRetries := 5
    for i := 0; i < maxRetries; i++ {
        resp, err := usc.client.GetUsers(context.Background(), &emptypb.Empty{})
        if err != nil {
            log.Printf("Could not list users: %v. Retrying in 5 seconds...", err)
            time.Sleep(5 * time.Second)
            continue
        }

        for _, user := range resp.Users {
            log.Printf("User: %d, Username: %s", user.Id, user.Username)
        }
        return
    }
    log.Fatalf("Failed to connect to user-service after %d attempts", maxRetries)
}
