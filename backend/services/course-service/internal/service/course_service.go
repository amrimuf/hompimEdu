package service

import (
	"context"
	"errors"
	"time"

	"github.com/amrimuf/hompimEdu/services/course-service/api/gen/coursepb"
)

// CourseServiceServer implements the CourseService gRPC interface
type CourseServiceServer struct {
    coursepb.UnimplementedCourseServiceServer
    courses map[int32]*coursepb.Course // In-memory store for courses
    nextID  int32                      // For auto-incrementing IDs
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