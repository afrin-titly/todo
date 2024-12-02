package utility

import (
	"errors"
	"net/http"
	"strings"
)

func ExtractTokenFromHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("Authorization header is required")
	}
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", errors.New("authorization header must start with 'Bearer'")
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == "" {
		return "", errors.New("token is required")
	}

	return token, nil
}
