package internal

import (
	"errors"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
    Username string `json:"username"`
    jwt.StandardClaims
}

func GenerateJWT(username string) (string, error) {
    secretKey := os.Getenv("JWT_SECRET")
    if secretKey == "" {
        return "", errors.New("JWT secret key not set")
    }

    expirationTime := time.Now().Add(24 * time.Hour) // Token expires in 24 hours

    claims := &Claims{
        Username: username,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: expirationTime.Unix(),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    return token.SignedString([]byte(secretKey))
}

func ValidateJWT(tokenString string) (string, error) {
    secretKey := os.Getenv("JWT_SECRET")
    claims := &Claims{}

    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        return []byte(secretKey), nil
    })

    if err != nil {
        return "", err
    }

    if !token.Valid {
        return "", errors.New("invalid token")
    }

    return claims.Username, nil
}
