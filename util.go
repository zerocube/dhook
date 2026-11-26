package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"time"
)

// httpClient has a timeout, so dhook does not wait forever for a slow server.
var httpClient = &http.Client{Timeout: 30 * time.Second}

// sendJSONPayload posts object as JSON to url. On any failure it logs the
// error and exits with status 1. Errors are unwrapped so that the webhook URL,
// which contains the webhook token, is not logged.
func sendJSONPayload(url *string, object interface{}) {
	objectJSON, err := json.Marshal(object)
	if err != nil {
		log.Fatalln("Unable to convert object to JSON:", err)
	}

	resp, err := httpClient.Post(*url, "application/json", bytes.NewReader(objectJSON))
	if err != nil {
		log.Fatalln("Unable to send request:", errors.Unwrap(err))
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println("Unable to read response body:", err)
		}
		log.Fatalln("Error: Request received a", resp.StatusCode, "status code from the server:", string(body))
	}
}

func allStringsAreEmpty(strArr []*string) bool {
	for _, str := range strArr {
		if *str != "" {
			return false
		}
	}
	return true
}
