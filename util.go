package main

import (
	"encoding/json"
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
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Println("WARNING: Request received a", resp.StatusCode, " status code from the server.")
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("Unable to read response body:", err)
		}
		log.Println(string(body))
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
