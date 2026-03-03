package handlers

import (
	"backend-chat-go/config"
	"backend-chat-go/models"
	"context"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
)

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	cursor, err := config.UserCollection.Find(
		context.Background(), bson.M{})

	if err != nil {
		writeJSON(w, 500, "Failed to fetch users", nil)
		return
	}

	var users []models.User
	if err = cursor.All(context.Background(), &users); err != nil {
		writeJSON(w, 500, "Decode Error", nil)
		return
	}

	var usernames []string
	for _, user := range users {
		usernames = append(usernames, user.Username)
	}

	writeJSON(w, 200, "Success", usernames)
}
