package authorization

import (
	"net/http"
	"strings"

	"github.com/Anuolu-2020/sync-meili/pkg/env"
)

func GetAuthorizationToken(r *http.Request) (string, bool) {
	authToken := r.Header.Get("Authorization")

	token := strings.Split(authToken, " ")
	if token[0] != "Bearer" {
		return "", false
	}

	return token[1], true
}

func VerifyAuthToken(token string) bool {
	return token == env.GetEnv("WEBHOOK_AUTH_TOKEN")
}
