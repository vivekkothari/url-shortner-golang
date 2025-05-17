package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"url-shortner/server"
)

// TestShortenAndRedirect tests the URL shortening and redirection.
func TestShortenAndRedirect(t *testing.T) {
	s, err := server.NewServer(":0")
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	ts := httptest.NewServer(*s.Router)
	defer ts.Close()

	originalURL := "https://example.com"
	reqBody := `{"long_url":"` + originalURL + `"}`

	resp, err := http.Post(ts.URL+"/shortUrlCreate", "application/json", bytes.NewReader([]byte(reqBody)))
	if err != nil {
		t.Fatalf("Failed to call shorten endpoint: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Expected status 201, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response: %v", err)
	}

	shortURL := extractShortURLFromBody(t, body)
	if shortURL == "" {
		t.Fatalf("Short URL not found in response: %s", string(body))
	}

	// Replace localhost with test server URL for redirection test
	testShortURL := strings.Replace(shortURL, "http://localhost:8080", ts.URL, 1)
	redirectResp, err := http.Get(testShortURL)
	if err != nil {
		t.Fatalf("Failed to call short URL: %v", err)
	}
	defer redirectResp.Body.Close()

	if redirectResp.Request.URL.String() != originalURL {
		t.Errorf("Expected redirect to %s, got %s", originalURL, redirectResp.Request.URL.String())
	}
}

// extractShortURLFromBody extracts the short URL from the response body.
func extractShortURLFromBody(t *testing.T, body []byte) string {
	t.Helper()
	var result map[string]string
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	return result["short_url"]
}
