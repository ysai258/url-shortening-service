// Package shorten is the main package of a URL shortening service.
package shorten

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/url-shortening-service/internal/constants"
	"github.com/url-shortening-service/pkg/db"
)

// ShortBody represents the body of a request to shorten a URL.
type ShortBody struct {
	Url string `json:"url"`
}

// generateShortLink creates a random short link.
func generateShortLink() string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, constants.SHORT_LINK_LENGTH)
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	n := len(letterBytes)
	for i := range b {
		b[i] = letterBytes[rand.Int()%n]
	}
	return string(b)
}

// isLinkValid checks if a short link already exists in the database.
func isLinkValid(link string) (bool, error) {
	var count int
	err := db.Db.QueryRow("SELECT COUNT(*) FROM links WHERE shortLink=?", link).Scan(&count)
	if err != nil {
		return false, err
	}
	return count == 0, nil
}

// getShortLink returns a unique short link.
func getShortLink() (string, error) {
	link := generateShortLink()
	checkCode, err := isLinkValid(link)
	for !checkCode {
		link = generateShortLink()
		checkCode, err = isLinkValid(link)
		if err != nil {
			return "", nil
		}
	}
	return link, err
}

// validLink checks if a URL is valid.
func validLink(link string) bool {
	r, err := regexp.Compile("^(http|https)://")
	if err != nil {
		return false
	}
	link = strings.TrimSpace(link)
	return r.MatchString(link)
}

// Shorten is the main function to shorten a URL.
func Shorten(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	payload := make(map[string]string)
	errKey := "error"
	msgKey := "message"
	respKey := "shortUrl"
	var body ShortBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		payload[msgKey] = "Error decoding request body"
		payload[errKey] = err.Error()
		json.NewEncoder(w).Encode(payload)
		return
	}
	longLink := body.Url
	if !validLink(longLink) {
		w.WriteHeader(http.StatusBadGateway)
		payload[msgKey] = "Error checking long link"
		payload[errKey] = "Invalid link"
		json.NewEncoder(w).Encode(payload)
		return

	}

	shortLink, err := getShortLink()
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		payload[msgKey] = "Error generating short link"
		payload[errKey] = err.Error()
		json.NewEncoder(w).Encode(payload)
		return
	}
	res, err := db.Db.Exec(`INSERT INTO links (shortLink, longLink) SELECT ?,? WHERE (SELECT count(*) FROM links WHERE shortLink = ?)=0`, shortLink, longLink, shortLink)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		payload[msgKey] = "Error storing short link"
		payload[errKey] = err.Error()
		json.NewEncoder(w).Encode(payload)
		return
	}
	rowsAff, err := res.RowsAffected()
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		payload[msgKey] = "Error fetching rows effected"
		payload[errKey] = err.Error()
		json.NewEncoder(w).Encode(payload)
		return
	}
	for rowsAff == 0 {
		shortLink, err = getShortLink()
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			payload[msgKey] = "Error generating short link 2"
			payload[errKey] = err.Error()
			json.NewEncoder(w).Encode(payload)
			return
		}
		res, err = db.Db.Exec(`INSERT INTO links (shortLink, longLink) SELECT ?,? WHERE (SELECT count(*) FROM links WHERE shortLink = ?)=0`, shortLink, longLink, shortLink)
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			payload[msgKey] = "Error storing short link 2"
			payload[errKey] = err.Error()
			json.NewEncoder(w).Encode(payload)
			return
		}
		rowsAff, err = res.RowsAffected()
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			payload[msgKey] = "Error fetching rows effected 2"
			payload[errKey] = err.Error()
			json.NewEncoder(w).Encode(payload)
			return
		}
	}

	w.WriteHeader(http.StatusOK)

	payload[respKey] = fmt.Sprintf("http://localhost:8080/%s", shortLink)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(payload)
}
