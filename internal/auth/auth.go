package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization key don't exists")
	}
	s := strings.Split(authHeader, " ")
	if len(s) < 2 || s[0] != "Bearer" {
		return "", errors.New("malformed authorization header")
	}
	return s[1], nil
}

func GetAPIKey(headers http.Header) (string, error) {
	apiHeader := headers.Get("Authorization")
	if apiHeader == "" {
		return "", errors.New("authorization key don't exists")
	}
	s := strings.Split(apiHeader, " ")
	if len(s) < 2 || s[0] != "ApiKey" {
		return "", errors.New("malformed authorization header")
	}
	return s[1], nil
}
