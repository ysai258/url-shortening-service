package shorten

import (
	"database/sql"
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

type ShortBody struct {
	Url string `json:"url"`
}

// Create random short link
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
func isLinkValid(link string) (bool, error) {
	var count int
	err := db.Db.QueryRow("SELECT COUNT(*) FROM links WHERE shortLink=?", link).Scan(&count)
	if err != nil {
		return false, err
	}
	return count == 0, nil
}
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
func checkLinkExists(link string) (string, error) {
	var short sql.NullString
	err := db.Db.QueryRow("SELECT shortLink FROM links WHERE longLink=? AND created>=now() - INTERVAL 1 DAY", link).Scan(&short)
	if err != nil && err != sql.ErrNoRows {
		return "", err
	}
	return short.String, nil
}
func validLink(link string) bool {
	r, err := regexp.Compile("^(http|https)://")
	if err != nil {
		return false
	}
	link = strings.TrimSpace(link)
	// Check if string matches the regex
	return r.MatchString(link)
}
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
	prevLink, err := checkLinkExists(longLink)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		payload[msgKey] = "Error checking long link"
		payload[errKey] = err.Error()
		json.NewEncoder(w).Encode(payload)
		return
	}
	if len(prevLink) > 0 {
		w.WriteHeader(http.StatusOK)
		payload[respKey] = fmt.Sprintf("http://localhost:8080/%s", prevLink)
		w.WriteHeader(http.StatusOK)
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
