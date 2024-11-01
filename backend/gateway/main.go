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

    // Public routes (no auth required)
    authRouter := r.PathPrefix("/auth").Subrouter()
    routes.RegisterAuthRoutes(authRouter)

    // Protected routes
    protectedRouter := r.PathPrefix("/api").Subrouter()
    protectedRouter.Use(middleware.JWTAuthMiddleware)
    routes.RegisterUserRoutes(protectedRouter)
    routes.RegisterCourseRoutes(protectedRouter)

    port := os.Getenv("API_GATEWAY_PORT")
    if port == "" {
        port = "8085"
    }

    log.Printf("API Gateway running on port %s", port)
    if err := http.ListenAndServe(":"+port, r); err != nil {
        log.Fatalf("API Gateway failed to start: %v", err)
    }
}
