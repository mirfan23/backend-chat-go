package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"backend-chat-go/config"
	"backend-chat-go/models"

	"go.mongodb.org/mongo-driver/bson"
)

type UserResponse struct {
	Username string `json:"username"`
	IsOnline bool   `json:"isOnline"`
}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {

	currentUser := r.Context().Value("username")
	if currentUser == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	ctx := context.Background()

	cursor, err := config.UserCollection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var users []UserResponse

	for cursor.Next(ctx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			continue
		}

		// skip diri sendiri
		if user.Username == currentUser.(string) {
			continue
		}

		users = append(users, UserResponse{
			Username: user.Username,
			IsOnline: IsUserOnline(user.Username),
		})
	}

	response := map[string]interface{}{
		"status":  200,
		"message": "Success",
		"data":    users,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
