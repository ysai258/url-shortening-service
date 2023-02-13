package redirect

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestRedirect(t *testing.T) {
	// Create a new mux router
	router := mux.NewRouter()

	// Register the Redirect function as a handler for the "/{shortLink}" route
	router.HandleFunc("/{shortLink}", Redirect)

	// GIVE SHORT LINK WHICH IS IN DATABASE
	shortLink := "cwlh5LV"
	url := fmt.Sprintf("http://localhost:8080/%s", shortLink)

	// Start an HTTP server using the mux router as the handler
	server := httptest.NewServer(router)
	defer server.Close()

	t.Run("Valid short link", func(t *testing.T) {
		// Make a request to the server with a short link
		res, err := http.Get(url)
		if err != nil {
			t.Fatalf("Unable to make request: %v", err)
		}

		// Check the response status code
		if res.StatusCode != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, res.StatusCode)
		}
	})

	t.Run("Invalid short link", func(t *testing.T) {
		// Make a request to the server with an invalid short link
		res, err := http.Get("http://localhost:8080/invalid")
		if err != nil {
			t.Fatalf("Unable to make request: %v", err)
		}

		// Check the response status code
		if res.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, res.StatusCode)
		}
	})
}
