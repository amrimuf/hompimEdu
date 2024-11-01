package api

import (
	"github.com/gorilla/mux"
	"github.com/amrimuf/hompimEdu/services/course-service/api/handlers"
	pb "github.com/amrimuf/hompimEdu/services/course-service/api/gen/coursepb"
)

func RegisterRoutes(router *mux.Router, grpcClient pb.CourseServiceClient) {
	handler := handlers.NewCourseHandler(grpcClient)

	router.HandleFunc("/api/courses", handler.CreateCourse).Methods("POST")
	router.HandleFunc("/api/courses/{courseId}", handler.GetCourse).Methods("GET")
	router.HandleFunc("/api/courses/{courseId}", handler.UpdateCourse).Methods("PUT")
} 