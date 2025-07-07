package utils

import (
	"encoding/json"
	"net/http"
)

// WriteJSONError returns a consistent JSON error response
func WriteJSONError(w http.ResponseWriter, status int, code, message string) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{
		"error":   code,
		"message": message,
	})
}
