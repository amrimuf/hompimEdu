package api

import (
	"net/http"
)

func RegisterRoutes() {
    http.HandleFunc("/users", usersHandler)
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodGet {
        w.Write([]byte("List of users"))
    } else {
        http.Error(w, "Unsupported HTTP method", http.StatusMethodNotAllowed)
    }
}
