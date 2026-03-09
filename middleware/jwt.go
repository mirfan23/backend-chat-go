package middleware

import (
	"context"
	"net/http"
	"strings"

	"backend-chat-go/config"
	"backend-chat-go/utils"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func unauthorized(w http.ResponseWriter, message string) {
	utils.WriteJSON(w, 401, message, nil)
}

func JWTMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			unauthorized(w, "Missing Authorization header")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			unauthorized(w, "Invalid Authorization format")
			return
		}

		tokenString := parts[1]

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(
			tokenString,
			claims,
			func(token *jwt.Token) (interface{}, error) {
				return config.JwtSecret, nil
			},
		)

		if err != nil || !token.Valid {
			unauthorized(w, "Token expired or invalid")
			return
		}

		ctx := context.WithValue(r.Context(), "username", claims.Username)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
