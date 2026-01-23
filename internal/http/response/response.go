package response

import (
	"encoding/json"
	"net/http"
)

type messageResponse struct {
	Message string `json:"message"`
}

func RespondJSON(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(messageResponse{
		Message: message,
	})
}

func Unauthorized(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": "authorization required",
	})
}
