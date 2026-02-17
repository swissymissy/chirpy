package auth

import (
	"fmt"


	"github.com/alexedwards/argon2id"
)

func HashPassword(password string) (string, error) {

	// hash the given password
	hashpwd, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", fmt.Errorf("Error hashing password: %w", err)
	}
	return hashpwd, nil
}

func CheckPasswordHash(password, hash string) (bool, error) {

	// compare entered passwd from user with the hashed passwd
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, fmt.Errorf("Error comparing password and hashed string: %w", err)
	}
	return match, nil
}