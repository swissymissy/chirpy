package auth

import (
	"net/http"
	"strings"
	"errors"
)

// function used for extracting API key in header
func GetAPIKey(headers http.Header) (string, error) {
	header := headers.Get("Authorization")			// get "authorization" from header list
	if header == "" {
		return "", errors.New("Invalid header")
	}

	// ensure that the value starts with "ApiKey "
	if !strings.HasPrefix(header, "ApiKey ") {
		return "", errors.New("Invalid header")
	}

	// strip out the prefix and the space to get the api key
	apiKey := strings.TrimSpace(strings.TrimPrefix(header, "ApiKey "))
	if apiKey == "" {
		return "", errors.New("Invalid api key")
	}
	return apiKey, nil
}