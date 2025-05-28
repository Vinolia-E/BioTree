package util

import (
	"encoding/json"
	"net/http"
)

/*
RespondError writes a JSON response with status "error" and a message.
It sets the HTTP status code to 500 (Internal Server Error).
Example response:
{
  "status": "error",
  "message": "Something went wrong"
}

RespondSuccess writes a JSON response with status "ok" indicating success.
Example response:
{
  "status": "ok"
}
*/


func RespondError(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "error",
		"message": message,
	})
}

func RespondSuccess(w http.ResponseWriter) {
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}
