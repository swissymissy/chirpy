package auth

import (
	"net/http"
	"strings"
	"errors"

)

// extract token sent from user
func GetBearerToken(headers http.Header) (string, error) {

	header := headers.Get("Authorization")
	if header == "" {
		return "", errors.New("Invalid header")
	}

	// ensure that the value starts with "Bearer "
	if !strings.HasPrefix( header, "Bearer ") {
		return "", errors.New("Invalid header")
	}

	// trim the prefix "Bearer" and the space to extract the token
	token := strings.TrimSpace(strings.TrimPrefix(header, "Bearer "))
	if token == "" {
		return "", errors.New("Invalid token")
	}

	return token, nil
}