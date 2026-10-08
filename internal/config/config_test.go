package config

import "testing"

func TestValidate(t *testing.T) {
	valid := &Config{
		DBHost: "localhost", DBPort: "5432", DBUser: "user",
		DBPassword: "password", DBName: "av_ing",
		JWTSecret: "12345678901234567890123456789012",
	}
	if err := valid.Validate(); err != nil {
		t.Fatalf("Validate() unexpected error = %v", err)
	}

	invalid := *valid
	invalid.JWTSecret = "short"
	if err := invalid.Validate(); err == nil {
		t.Fatal("Validate() expected an error for a short JWT secret")
	}
}
