package home

import (
	"encoding/json"
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := make(map[string]string)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
	w.WriteHeader(http.StatusOK)
}
