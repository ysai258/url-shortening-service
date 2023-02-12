// main is the entry point of the program
package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/url-shortening-service/internal/home"
	"github.com/url-shortening-service/internal/redirect"
	"github.com/url-shortening-service/internal/shorten"
	"github.com/url-shortening-service/pkg/db"
)

func main() {
	// Initialize the database
	db.InitDB()

	// Create a new router
	r := mux.NewRouter()

	// Define routes and associated handlers
	r.HandleFunc("/shorten", shorten.Shorten).Methods("POST")
	r.HandleFunc("/", home.Home)
	r.HandleFunc("/{shortLink}", redirect.Redirect)

	// Start the HTTP server
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Panic(err)
	}
}
