package routes

import (
	"net/http"

	"github.com/amrimuf/hompimEdu/gateway/services"
	"github.com/gorilla/mux"
)

func RegisterUserRoutes(r *mux.Router) {
	r.HandleFunc("/users", listUsersHandler).Methods("GET")
}

func listUsersHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := services.UserClient("GET", "/users", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	copyResponse(w, resp)
}
