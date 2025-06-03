package util

import (
	"encoding/json"
	"net/http"
)

// RespondError sends an error response to the client
func RespondError(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)

	response := map[string]interface{}{
		"status":  "error",
		"message": message,
	}

	json.NewEncoder(w).Encode(response)
}

// RespondSuccess sends a success response to the client
func RespondSuccess(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")

	response := map[string]interface{}{
		"status": "ok",
		"data":   data,
	}

	json.NewEncoder(w).Encode(response)
}
