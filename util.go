package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func sendJSONPayload(url *string, object interface{}) error {

	// Serialise object to JSON
	objectJSON, err := json.Marshal(object)
	if err != nil {
		log.Println("Unable to convert object to JSON:", object)
		log.Fatalln(err)
	}

	req, err := http.NewRequest(http.MethodPost, *url, strings.NewReader(string(objectJSON)))
	if err != nil {
		log.Println("Unable to create request:", err)
		log.Fatalln(req)
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalln("Unable to send request:", errors.Unwrap(err))
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("Unable to read response body:", err)
		}
		log.Fatalln("Error: Request received a", resp.StatusCode, "status code from the server:", string(body))
	}
	return err
}

func allStringsAreEmpty(strArr []*string) bool {
	for _, str := range strArr {
		if *str != "" {
			return false
		}
	}
	return true
}
