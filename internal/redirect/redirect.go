package redirect

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/url-shortening-service/pkg/db"
)

// invalid writes a JSON response to the client with a status code of Bad Request, indicating that the provided short link is invalid
func invalid(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := make(map[string]string)
	w.WriteHeader(http.StatusBadRequest)
	response["error"] = "Bad Request"
	response["message"] = "Invalid link"
	json.NewEncoder(w).Encode(response)
	w.WriteHeader(http.StatusOK)
}

// fetchLongLink fetches the corresponding long link for a given short link from the database
func fetchLongLink(link string) (string, error) {
	var short sql.NullString
	err := db.Db.QueryRow("SELECT longLink FROM links WHERE shortLink=? AND created>=now() - INTERVAL 1 DAY", link).Scan(&short)
	if err != nil && err != sql.ErrNoRows {
		return "", err
	}
	return short.String, nil
}

// Redirect redirects the client to the corresponding long link for a given short link. If the short link is not found or an error occurs, it sends a JSON response with an error message
func Redirect(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := make(map[string]string)
	errKey := "error"
	msgKey := "message"

	vars := mux.Vars(r)
	shortLink, ok := vars["shortLink"]
	if !ok {
		invalid(w, r)
		return
	}
	longLink, err := fetchLongLink(shortLink)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		response[errKey] = err.Error()
		response[msgKey] = "Error in fetching long link"
		json.NewEncoder(w).Encode(response)
		w.WriteHeader(http.StatusOK)
	}
	if len(longLink) == 0 {
		invalid(w, r)
		return
	}
	http.Redirect(w, r, longLink, http.StatusTemporaryRedirect)
}
