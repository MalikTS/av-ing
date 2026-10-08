package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"av-ing/internal/config"
)

func TestAuthMiddleware(t *testing.T) {
	secret := "12345678901234567890123456789012"
	InitAuthMiddleware(&config.Config{JWTSecret: secret})
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  "user-id",
		"is_admin": true,
		"exp":      time.Now().Add(time.Hour).Unix(),
	})
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("SignedString() error = %v", err)
	}

	handler := AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := GetUserIDFromContext(r); got != "user-id" {
			t.Fatalf("user ID = %q, want user-id", got)
		}
		if !IsAdminFromContext(r) {
			t.Fatal("admin claim was not added to context")
		}
		w.WriteHeader(http.StatusNoContent)
	}))

	request := httptest.NewRequest(http.MethodGet, "/me", nil)
	request.Header.Set("Authorization", "Bearer "+signed)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusNoContent)
	}
}

func TestAuthMiddlewareRejectsMalformedHeader(t *testing.T) {
	InitAuthMiddleware(&config.Config{JWTSecret: "12345678901234567890123456789012"})
	handler := AuthMiddleware(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		t.Fatal("next handler must not be called")
	}))

	request := httptest.NewRequest(http.MethodGet, "/me", nil)
	request.Header.Set("Authorization", "not-a-bearer-token")
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusUnauthorized)
	}
}
