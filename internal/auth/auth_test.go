package auth 

import (
	"testing"
	"time"
	"github.com/google/uuid"
)


func TestMakeJWTAndValidateJWT (t *testing.T) {
	// create test object
	secret := "test-object-secret"
	userID := uuid.New()
	expiresIn := time.Hour 			

	// test function
	token, err := MakeJWT(userID, secret, expiresIn)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// validate
	returnredID, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// asset
	if returnredID != userID {
		t.Errorf("expected %v, got %v", userID, returnredID) 
	}
}

func TestJWTEdgeCase (t *testing.T) {
	// create test object with expired duration
	secret := "test-object-secret"
	userID := uuid.New()
	expiresIn := -1 * time.Hour			// expired an hour ago already

	// test Make function
	token, err := MakeJWT( userID, secret, expiresIn)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// validate
	_, err = ValidateJWT(token, secret)
	if err == nil {
		t.Errorf("expected an error for an expired token, but got nil")
	}
}