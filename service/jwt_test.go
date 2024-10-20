package service

import (
	"log"
	"testing"
)

func TestJWTService(t *testing.T) {
	secret := "test" // Your test secret
	jwtService := NewJWTService(secret)

	// Test generating a token
	userID := 1
	token, err := jwtService.GenerateToken(userID)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
		return // Exit early to prevent using an invalid token
	}

	// Test validating the generated token
	validatedUserID, err := jwtService.ValidateToken(token)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
		return // Exit early if validation fails
	}
	log.Println(validatedUserID)
	if validatedUserID != userID {
		t.Errorf("expected userID %d, got %d", userID, validatedUserID)
	}
}
