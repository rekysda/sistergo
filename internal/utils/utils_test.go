package utils

import (
	"testing"
	"time"
)

func TestHashPassword(t *testing.T) {
	password := "testpassword"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}
	if hash == "" {
		t.Error("HashPassword returned empty string")
	}
	if hash == password {
		t.Error("HashPassword returned unhashed password")
	}
}

func TestCheckPassword(t *testing.T) {
	password := "testpassword"
	hash, _ := HashPassword(password)

	if !CheckPassword(hash, password) {
		t.Error("CheckPassword failed to validate correct password")
	}

	if CheckPassword(hash, "wrongpassword") {
		t.Error("CheckPassword validated incorrect password")
	}
}

func TestGenerateJWT(t *testing.T) {
	secret := "testsecret"
	userID := uint(1)
	expiry := time.Hour

	token, err := GenerateJWT(secret, userID, expiry)
	if err != nil {
		t.Fatalf("GenerateJWT failed: %v", err)
	}
	if token == "" {
		t.Error("GenerateJWT returned empty token")
	}
}

func TestParseJWT(t *testing.T) {
	secret := "testsecret"
	userID := uint(42)
	expiry := time.Hour

	token, _ := GenerateJWT(secret, userID, expiry)

	claims, err := ParseJWT(secret, token)
	if err != nil {
		t.Fatalf("ParseJWT failed: %v", err)
	}

	// Check user_id claim
	if uid, ok := claims["user_id"].(float64); !ok || uint(uid) != userID {
		t.Errorf("Expected user_id %d, got %v", userID, claims["user_id"])
	}
}

func TestParseJWT_InvalidToken(t *testing.T) {
	secret := "testsecret"
	invalidToken := "invalid.token.string"

	_, err := ParseJWT(secret, invalidToken)
	if err == nil {
		t.Error("ParseJWT should fail with invalid token")
	}
}

func TestParseJWT_WrongSecret(t *testing.T) {
	secret := "testsecret"
	wrongSecret := "wrongsecret"
	userID := uint(1)
	expiry := time.Hour

	token, _ := GenerateJWT(secret, userID, expiry)

	_, err := ParseJWT(wrongSecret, token)
	if err == nil {
		t.Error("ParseJWT should fail with wrong secret")
	}
}
