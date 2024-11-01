package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(bearerToken[1], claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Safely extract user_id from claims
		userID, ok := claims["user_id"]
		if !ok {
			http.Error(w, "Invalid token: missing user_id", http.StatusUnauthorized)
			return
		}

		// Convert user_id to float64 first, then to int64
		userIDFloat, ok := userID.(float64)
		if !ok {
			http.Error(w, "Invalid token: user_id format invalid", http.StatusUnauthorized)
			return
		}

		// Add user info to request headers
		r.Header.Set("X-User-ID", fmt.Sprintf("%d", int64(userIDFloat)))
		
		// Forward the request
		next.ServeHTTP(w, r)
	})
}
