package models

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestUserPasswordHashIsNotSerialized(t *testing.T) {
	data, err := json.Marshal(User{ID: "user-id", PasswordHash: "secret-hash"})
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}
	if strings.Contains(string(data), "secret-hash") || strings.Contains(string(data), "PasswordHash") {
		t.Fatalf("password hash leaked into JSON: %s", data)
	}
}
