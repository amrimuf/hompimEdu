package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/amrimuf/hompimEdu/services/course-service/api/gen/coursepb"
	"github.com/amrimuf/hompimEdu/services/course-service/api/gen/userpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CourseServiceServer struct {
	coursepb.UnimplementedCourseServiceServer
	db *sql.DB
	userClient userpb.UserServiceClient
}

func NewCourseServiceServer(db *sql.DB, userClient userpb.UserServiceClient) *CourseServiceServer {
	return &CourseServiceServer{
		db: db,
		userClient: userClient,
	}
}

// GetCourse retrieves a course by ID
func (s *CourseServiceServer) GetCourse(ctx context.Context, req *coursepb.GetCourseRequest) (*coursepb.GetCourseResponse, error) {
	var course coursepb.Course
	err := s.db.QueryRowContext(ctx, `
		SELECT id, title, description, duration, enrollment_type, mentor_id, created_at, updated_at 
		FROM courses WHERE id = $1`, req.Id).Scan(
		&course.Id, &course.Title, &course.Description, &course.Duration,
		&course.EnrollmentType, &course.MentorId, &course.CreatedAt, &course.UpdatedAt)
	
	if err == sql.ErrNoRows {
		return nil, errors.New("course not found")
	}
	if err != nil {
		return nil, err
	}
	
	return &coursepb.GetCourseResponse{Course: &course}, nil
}

// CreateCourse creates a new course
func (s *CourseServiceServer) CreateCourse(ctx context.Context, req *coursepb.CreateCourseRequest) (*coursepb.CreateCourseResponse, error) {
	log.Printf("Attempting to verify mentor with ID: %d", req.MentorId)

	// First verify if mentor exists
	_, err := s.userClient.GetUser(ctx, &userpb.UserRequest{
		UserId: req.MentorId,
	})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, fmt.Errorf("mentor not found")
		}
		log.Printf("Error verifying mentor: %v", err)
		return nil, fmt.Errorf("failed to verify mentor: %v", err)
	}

	log.Printf("Mentor verified successfully")

	var course coursepb.Course
	now := time.Now().Format(time.RFC3339)
	
	err = s.db.QueryRowContext(ctx, `
		INSERT INTO courses (title, description, duration, enrollment_type, mentor_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, title, description, duration, enrollment_type, mentor_id, created_at, updated_at`,
		req.Title, req.Description, req.Duration, req.EnrollmentType, req.MentorId, now, now,
	).Scan(&course.Id, &course.Title, &course.Description, &course.Duration,
		&course.EnrollmentType, &course.MentorId, &course.CreatedAt, &course.UpdatedAt)
	
	if err != nil {
		return nil, fmt.Errorf("failed to create course: %v", err)
	}
	
	return &coursepb.CreateCourseResponse{Course: &course}, nil
}

// ListCourses lists all available courses
func (s *CourseServiceServer) ListCourses(ctx context.Context, req *coursepb.ListCoursesRequest) (*coursepb.ListCoursesResponse, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, title, description, duration, enrollment_type, mentor_id, created_at, updated_at 
		FROM courses`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courses []*coursepb.Course
	for rows.Next() {
		var course coursepb.Course
		err := rows.Scan(
			&course.Id, &course.Title, &course.Description, &course.Duration,
			&course.EnrollmentType, &course.MentorId, &course.CreatedAt, &course.UpdatedAt)
		if err != nil {
			return nil, err
		}
		courses = append(courses, &course)
	}
	
	if err = rows.Err(); err != nil {
		return nil, err
	}
	
	return &coursepb.ListCoursesResponse{Courses: courses}, nil
}

// UpdateCourse updates an existing course
func (s *CourseServiceServer) UpdateCourse(ctx context.Context, req *coursepb.UpdateCourseRequest) (*coursepb.UpdateCourseResponse, error) {
	var course coursepb.Course
	now := time.Now().Format(time.RFC3339)
	
	err := s.db.QueryRowContext(ctx, `
		UPDATE courses 
		SET title = $1, description = $2, duration = $3, enrollment_type = $4, 
			mentor_id = $5, updated_at = $6
		WHERE id = $7
		RETURNING id, title, description, duration, enrollment_type, mentor_id, created_at, updated_at`,
		req.Title, req.Description, req.Duration, req.EnrollmentType, 
		req.MentorId, now, req.Id,
	).Scan(&course.Id, &course.Title, &course.Description, &course.Duration,
		&course.EnrollmentType, &course.MentorId, &course.CreatedAt, &course.UpdatedAt)
	
	if err == sql.ErrNoRows {
		return nil, errors.New("course not found")
	}
	if err != nil {
		return nil, err
	}
	
	return &coursepb.UpdateCourseResponse{Course: &course}, nil
}
