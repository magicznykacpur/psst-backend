package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetBearerToken(header http.Header) (string, error) {
	authorization := header.Get("Authorization")
	if authorization == "" {
		return "", fmt.Errorf("authorization header not provided")
	}

	bearerParts := strings.Split(authorization, " ")
	if len(bearerParts) != 2 {
		return "", fmt.Errorf("authorization header malformed")
	}

	if bearerParts[0] != "Bearer" {
		return "", fmt.Errorf("bearer token malformed")
	}

	return bearerParts[1], nil
}
