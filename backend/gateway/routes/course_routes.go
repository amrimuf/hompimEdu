package routes

import (
	"net/http"

	"github.com/amrimuf/hompimEdu/gateway/services"
	"github.com/gorilla/mux"
)

func RegisterCourseRoutes(r *mux.Router) {
	r.HandleFunc("/courses", createCourseHandler).Methods("POST")
	r.HandleFunc("/courses/{courseId}", getCourseHandler).Methods("GET")
	r.HandleFunc("/courses/{courseId}", updateCourseHandler).Methods("PUT")
}

func createCourseHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := services.CourseClient("POST", "/api/courses", r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	copyResponse(w, resp)
}

func getCourseHandler(w http.ResponseWriter, r *http.Request) {
	courseID := mux.Vars(r)["courseId"]
	resp, err := services.CourseClient("GET", "/api/courses/"+courseID, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	copyResponse(w, resp)
}

func updateCourseHandler(w http.ResponseWriter, r *http.Request) {
	courseID := mux.Vars(r)["courseId"]
	resp, err := services.CourseClient("PUT", "/api/courses/"+courseID, r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	copyResponse(w, resp)
} 