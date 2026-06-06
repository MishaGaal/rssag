package auth

import (
	"errors"
	"net/http"
	"strings"
)

const (
	AuthorizationHeader = "Authorization"
	ApiKeyPrefix        = "ApiKey"
)

// Authorization : ApiKey: {api_key}
func GetApiKey(headers http.Header) (string, error) {
	auth := headers.Get(AuthorizationHeader)
	if auth == "" {
		return "", errors.New("Authorization header not found")
	}

	authParts := strings.Split(auth, " ")
	if len(authParts) != 2 {
		return "", errors.New("Authorization header has invalid format")
	}

	if authParts[0] != ApiKeyPrefix {
		return "", errors.New("Authorization header has invalid prefix")
	}
	return authParts[1], nil
}
