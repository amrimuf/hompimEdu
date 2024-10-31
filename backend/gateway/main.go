package main

import (
	"log"
	"net/http"
	"os"

	"github.com/amrimuf/hompimEdu/gateway/middleware"
	"github.com/amrimuf/hompimEdu/gateway/routes"

	"github.com/gorilla/mux"
)

func main() {
    r := mux.NewRouter()

    // Register public routes without authentication (e.g., registration, login)
    authRouter := r.PathPrefix("/auth").Subrouter()
    routes.RegisterAuthRoutes(authRouter) // These routes won't require JWT auth

    // Register protected routes that require authentication
    protectedRouter := r.PathPrefix("/").Subrouter()
    protectedRouter.Use(middleware.JWTAuthMiddleware) // Apply JWT middleware only here
    routes.RegisterUserRoutes(protectedRouter) // Protected routes

    port := os.Getenv("API_GATEWAY_PORT")
    log.Printf("API Gateway running on port %s", port)
    if err := http.ListenAndServe(":"+port, r); err != nil {
        log.Fatalf("API Gateway failed to start: %v", err)
    }
}
