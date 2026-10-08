package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"

	"av-ing/internal/config"
)

type contextKey string

const UserContextKey contextKey = "userID"
const AdminContextKey contextKey = "IsAdmin"

var cfg *config.Config

func InitAuthMiddleware(c *config.Config) {
	cfg = c
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authParts := strings.Fields(r.Header.Get("Authorization"))
		if len(authParts) == 0 {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}
		if len(authParts) != 2 || !strings.EqualFold(authParts[0], "Bearer") {
			http.Error(w, "Invalid authorization header", http.StatusUnauthorized)
			return
		}

		tokenString := authParts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(cfg.JWTSecret), nil
		}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		userID, ok := claims["user_id"].(string)
		if !ok {
			http.Error(w, "Invalid user_id in token", http.StatusUnauthorized)
			return
		}

		isAdmin := false
		if val, ok := claims["is_admin"].(bool); ok {
			isAdmin = val
		}

		ctx := context.WithValue(r.Context(), UserContextKey, userID)
		ctx = context.WithValue(ctx, AdminContextKey, isAdmin)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserIDFromContext(r *http.Request) string {
	userID, ok := r.Context().Value(UserContextKey).(string)
	if !ok {
		return ""
	}
	return userID
}

func IsAdminFromContext(r *http.Request) bool {
	isAdmin, ok := r.Context().Value(AdminContextKey).(bool)
	if !ok {
		return false
	}
	return isAdmin
}
