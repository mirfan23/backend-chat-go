package utils

import (
	"backend-chat-go/models"
	"encoding/json"
	"net/http"
	"strings"
)

func ExtractFriend(roomID string, currentUser string) string {
	parts := strings.Split(roomID, "_")
	if len(parts) != 2 {
		return ""
	}
	if parts[0] == currentUser {
		return parts[1]
	}
	return parts[0]
}

func WriteJSON(w http.ResponseWriter, status int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(models.ApiResponse{
		StatusCode: status,
		Message:    message,
		Data:       data,
	})
}
