package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	authString := headers.Get("Authorization")
	authSlice := strings.Fields(authString)

	if len(authSlice) != 2 {
		return "", errors.New("there is no valid authorization header in the header")
	}

	if authSlice[0] != "ApiKey" || authSlice[1] == "" {
		return "", errors.New("invalid api key format")
	}

	return authSlice[1], nil
}
