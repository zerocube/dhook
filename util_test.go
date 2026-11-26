package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAllStringsAreEmpty(t *testing.T) {
	empty, set := "", "x"
	tests := []struct {
		name string
		in   []*string
		want bool
	}{
		{"no strings", nil, true},
		{"all empty", []*string{&empty, &empty}, true},
		{"one set", []*string{&empty, &set}, false},
	}
	for _, tt := range tests {
		if got := allStringsAreEmpty(tt.in); got != tt.want {
			t.Errorf("%s: allStringsAreEmpty() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestSendJSONPayload(t *testing.T) {
	var method, contentType string
	var body map[string]string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method, contentType = r.Method, r.Header.Get("Content-Type")
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("request body is not valid JSON: %v", err)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	url := server.URL
	sendJSONPayload(&url, map[string]string{"content": "hello"})

	if method != http.MethodPost {
		t.Errorf("method = %q, want %q", method, http.MethodPost)
	}
	if contentType != "application/json" {
		t.Errorf("Content-Type = %q, want %q", contentType, "application/json")
	}
	if body["content"] != "hello" {
		t.Errorf("content = %q, want %q", body["content"], "hello")
	}
}
