package webhook

import (
	"bufio"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Anuolu-2020/sync-meili/internal/authorization"
)

func WebhookHandler(w http.ResponseWriter, r *http.Request) {
	if _, ok := os.LookupEnv("WEBHOOK_AUTH_TOKEN"); ok {
		// if no token is provided, log and return
		token, ok := authorization.GetAuthorizationToken(r)
		if !ok {
			log.Printf(
				"Authorization token not provided for webhook request from IP: %s",
				r.RemoteAddr,
			)
			return
		}

		// Check if token is authentic
		if isVerified := authorization.VerifyAuthToken(token); !isVerified {
			log.Printf(
				"Authorization token verification failed for webhook request from IP: %s",
				r.RemoteAddr,
			)
			return
		}
	} else {
		log.Printf("Webhook auth token not set, authorizating webhook request without token")
	}

	defer r.Body.Close()

	bodyReader := r.Body

	// Check if content is gzip compressed
	if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
		var err error
		bodyReader, err = gzip.NewReader(r.Body)
		if err != nil {
			log.Printf("Error creating gzip reader: %v\n", err)
			return
		}
		defer bodyReader.Close()
	}

	scanner := bufio.NewScanner(bodyReader)
	for scanner.Scan() {
		var webHookBody TaskDetails
		line := scanner.Text()

		// Strip any non-JSON prefix characters, if they exist
		//line = strings.TrimLeft(line, "\x1f\x1e")

		// Decode each line as a JSON object
		if err := json.Unmarshal([]byte(line), &webHookBody); err != nil {
			log.Printf("Error decoding NDJSON line: %v\n", err)
			continue // Optionally continue to the next line on error
		}

		// Logic to store successfully parsed task
		fmt.Printf("webhook: %v\n", webHookBody)
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading webhook request body: %v\n", err)
	}
}
