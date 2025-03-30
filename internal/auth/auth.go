package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("authorization key don't exists")
	}
	s := strings.Split(authHeader, " ")
	return s[len(s)-1], nil

}
