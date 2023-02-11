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
	db.InitDB()

	r := mux.NewRouter()
	r.HandleFunc("/shorten", shorten.Shorten).Methods("POST")
	r.HandleFunc("/", home.Home)
	r.HandleFunc("/{shortLink}", redirect.Redirect)

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Panic(err)
	}
}
