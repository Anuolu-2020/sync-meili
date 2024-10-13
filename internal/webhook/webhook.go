package webhook

import (
	"bufio"
	"encoding/json"
	"log"
	"net/http"

	"github.com/Anuolu-2020/sync-meili/internal/authorization"
)

func WebhookHandler(w http.ResponseWriter, r *http.Request) {
	// if no token is provided, log and return
	token, ok := authorization.GetAuthorizationToken(r)
	if !ok {
		log.Printf("Authorization token not provided for request from IP: %s", r.RemoteAddr)
		return
	}

	// Check if token is authentic
	if isVerified := authorization.VerifyAuthToken(token); !isVerified {
		log.Printf("Authorization token verification failed for request from IP: %s", r.RemoteAddr)
		return
	}

	scanner := bufio.NewScanner(r.Body)
	defer r.Body.Close()

	for scanner.Scan() {
		line := scanner.Text()

		var webHookBody TaskDetails
		err := json.Unmarshal([]byte(line), &webHookBody)
		if err != nil {
			log.Printf("Error decoding JSON: %v\n", err)
			continue
		}

		// Logic to store successfully task
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading webhook request body: %v\n", err)
	}
}
