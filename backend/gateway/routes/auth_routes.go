package routes

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/amrimuf/hompimEdu/gateway/services"

	"github.com/gorilla/mux"
)

func RegisterAuthRoutes(r *mux.Router) {
	r.HandleFunc("/login", loginHandler).Methods("POST")
	r.HandleFunc("/register", registerHandler).Methods("POST")
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var credentials map[string]string

	log.Printf("Received login request: %v", r)

	// Decode the request body and handle potential errors
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		log.Printf("Error decoding login request: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	log.Printf("Decoded credentials: %+v", credentials)

	// Call the authentication service
	resp, err := services.AuthClient("POST", "/login", credentials)
	if err != nil {
		log.Printf("Error calling auth service for login: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	log.Printf("Received response from auth service for login: %v", resp.Status)

	copyResponse(w, resp)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	var userData map[string]string
	json.NewDecoder(r.Body).Decode(&userData)

	resp, err := services.AuthClient("POST", "/register", userData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	copyResponse(w, resp)
}

func copyResponse(w http.ResponseWriter, resp *http.Response) {
	// Set the status code of the response writer
	w.WriteHeader(resp.StatusCode)

	// Copy the headers from the response to the response writer
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// Copy the response body to the response writer
	_, _ = io.Copy(w, resp.Body)
}
