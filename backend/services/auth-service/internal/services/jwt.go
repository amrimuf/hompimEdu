package services

import (
	"errors"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	Username string `json:"username"`
	UserID   int64  `json:"user_id"`
	jwt.StandardClaims
}

func GenerateJWT(username string, userID int64) (string, error) {
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		return "", errors.New("JWT secret key not set")
	}

	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		Username: username,
		UserID:   userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secretKey))
}

func ValidateJWT(tokenString string) (*Claims, error) {
	secretKey := os.Getenv("JWT_SECRET")
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
