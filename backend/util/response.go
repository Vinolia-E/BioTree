package util

import (
	"encoding/json"
	"net/http"
)

func RespondError(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "error",
		"message": message,
	})
}

func RespondSuccess(w http.ResponseWriter, data any) {
	json.NewEncoder(w).Encode(map[string]any{
		"status": "ok",
		"units":  data,
	})
}
