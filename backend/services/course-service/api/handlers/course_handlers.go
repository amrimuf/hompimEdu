package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"github.com/gorilla/mux"
	"context"
	pb "github.com/amrimuf/hompimEdu/services/course-service/api/gen/coursepb"
)

type CourseHandler struct {
	grpcClient pb.CourseServiceClient
}

func NewCourseHandler(client pb.CourseServiceClient) *CourseHandler {
	return &CourseHandler{grpcClient: client}
}

func (h *CourseHandler) CreateCourse(w http.ResponseWriter, r *http.Request) {
	var req pb.CreateCourseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	resp, err := h.grpcClient.CreateCourse(context.Background(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp.Course)
}

func (h *CourseHandler) GetCourse(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	courseID, err := strconv.Atoi(vars["courseId"])
	if err != nil {
		http.Error(w, "Invalid course ID", http.StatusBadRequest)
		return
	}

	resp, err := h.grpcClient.GetCourse(context.Background(), &pb.GetCourseRequest{Id: int32(courseID)})
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp.Course)
}

func (h *CourseHandler) UpdateCourse(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	courseID, err := strconv.Atoi(vars["courseId"])
	if err != nil {
		http.Error(w, "Invalid course ID", http.StatusBadRequest)
		return
	}

	var req pb.UpdateCourseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	req.Id = int32(courseID)

	resp, err := h.grpcClient.UpdateCourse(context.Background(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp.Course)
} 