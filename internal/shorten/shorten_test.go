package shorten

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

func TestShorten(t *testing.T) {
	const lonLink = "https://www.google.com"
	tests := []struct {
		name       string
		url        string
		statusCode int
		response   string
	}{
		{
			name:       "Valid url",
			url:        lonLink,
			statusCode: http.StatusOK,
			response:   `{"message":"Success","shortUrl":"shortLink"}`,
		},
		{
			name:       "Invalid url",
			url:        "google.com",
			statusCode: http.StatusBadGateway,
			response:   `{"error":"Invalid link","message":"Error checking long link"}`,
		},
	}
	// Create a new mux router
	router := mux.NewRouter()

	// Register the Redirect function as a handler for the "/{shortLink}" route
	router.HandleFunc("/{shortLink}", Shorten)

	// Start an HTTP server using the mux router as the handler
	server := httptest.NewServer(router)
	defer server.Close()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			requestBody := `{"url":"` + test.url + `"}`
			resp, err := http.Post("http://localhost:8080/shorten", "application/json", strings.NewReader(requestBody))

			if err != nil {
				t.Fatalf("Unable to make request: %v", err)
			}

			// Check the response status code
			if resp.StatusCode != test.statusCode {
				t.Errorf("Expected status code %d, got %d", http.StatusTemporaryRedirect, resp.StatusCode)
			}

		})
	}

}
