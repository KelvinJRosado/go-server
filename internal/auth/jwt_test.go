package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestJWTValidation(t *testing.T) {

	userId := uuid.New()
	tokenSecret := "foo"
	expiresIn := time.Minute * 60

	token, err := MakeJWT(userId, tokenSecret, expiresIn)
	if err != nil {
		t.Errorf("Got an unexpected error making JWT: %s", err)
	}

	if token == "" {
		t.Errorf("Expected a token, got an empty string")
	}

	// Validate happy path
	parsedId, err := ValidateJWT(token, tokenSecret)
	if err != nil {
		t.Errorf("Got an unexpected error validating JWT: %s", err)
	}

	if parsedId != userId {
		t.Errorf("Expected parsed ID to match user ID, got %s instead of %s", parsedId, userId)
	}

	// Validate wrong token secret is rejected
	_, err = ValidateJWT(token, "wrongSecret")
	if err == nil {
		t.Errorf("Expected an error validating JWT due to wrong secret, got nil")
	}
}

func TestExpiredToken(t *testing.T) {

	userId := uuid.New()
	tokenSecret := "foo"
	expiresIn := time.Microsecond

	token, err := MakeJWT(userId, tokenSecret, expiresIn)
	if err != nil {
		t.Errorf("Got an unexpected error making JWT: %s", err)
	}

	if token == "" {
		t.Errorf("Expected a token, got an empty string")
	}

	time.Sleep(time.Microsecond * 2)

	_, err = ValidateJWT(token, tokenSecret)
	if err == nil {
		t.Errorf("Expected an error validating JWT, got nil")
	}
}
