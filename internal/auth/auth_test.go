package auth

import (
	"testing"
)

func TestHashingWorks(t *testing.T) {

	password := "foo"

	hash, err := HashPassword(password)
	if err != nil {
		t.Errorf("Got an unexpected error hashing password: %s", err)
	}

	matches, err := CheckPasswordHash(password, hash)
	if err != nil {
		t.Errorf("Got an unexpected error checking password hash: %s", err)
	}
	if !matches {
		t.Errorf("Password does not match hash")
	}

}

func TestHashingWrongPassword(t *testing.T) {

	password := "foo"

	hash, err := HashPassword(password)
	if err != nil {
		t.Errorf("Got an unexpected error hashing password: %s", err)
	}

	wrongPassword := "bar"

	matches, err := CheckPasswordHash(wrongPassword, hash)
	if err != nil {
		t.Errorf("Got an unexpected error checking password hash: %s", err)
	}
	if matches {
		t.Errorf("Wrong password should not match hash")
	}

}

func TestHashingSalt(t *testing.T) {

	password := "foo"

	hash1, err := HashPassword(password)
	if err != nil {
		t.Errorf("Got an unexpected error hashing password: %s", err)
	}

	hash2, err := HashPassword(password)
	if err != nil {
		t.Errorf("Got an unexpected error hashing password: %s", err)
	}

	if hash1 == hash2 {
		t.Errorf("Hashed passwords should not match")
	}

}
