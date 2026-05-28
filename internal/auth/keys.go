package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {

	// Get auth header, in lowercase
	authHeader := headers.Get("Authorization")

	if authHeader == "" {
		return "", errors.New("no Authorization header")
	}

	// Make sure proper prefix is given
	if !strings.HasPrefix(strings.ToLower(authHeader), "apikey ") {
		return "", errors.New("invalid Authorization header")
	}

	// Remove prefix
	runes := []rune(authHeader)

	return string(runes[7:]), nil

}
