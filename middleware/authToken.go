package middleware

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
)

var jwtSecret = []byte("advanced")

//func AuthToken(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		// Extract the token from the Authorization header.
//		authHeader := r.Header.Get("Authorization")
//		if authHeader == "" {
//			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
//			return
//		}
//
//		// Check if the header is in the format "Bearer <token>".
//		tokenParts := strings.Split(authHeader, " ")
//		if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
//			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
//			return
//		}
//
//		tokenString := tokenParts[1]
//
//		// Parse and validate the token.
//		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
//			// Validate the signing method.
//			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
//				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
//			}
//			return jwtSecret, nil
//		})
//
//		if err != nil {
//			http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
//			return
//		}
//
//		// Check if the token is valid.
//		if !token.Valid {
//			http.Error(w, "Invalid token", http.StatusUnauthorized)
//			return
//		}
//
//		// If the token is valid, pass the request to the next handler.
//		next.ServeHTTP(w, r)
//	})
//}

type UserClaims struct {
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.RegisteredClaims
}

// AuthToken проверяет JWT и добавляет данные пользователя в контекст
func AuthToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims := &UserClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// Здесь должен быть твой секретный ключ
			return []byte("advanced"), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Добавляем данные о пользователе в контекст запроса
		ctx := context.WithValue(r.Context(), "userRole", claims.Role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
