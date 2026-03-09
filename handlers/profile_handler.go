package handlers

import (
	"context"
	"net/http"

	"backend-chat-go/config"
	"backend-chat-go/models"
	"backend-chat-go/utils"

	"go.mongodb.org/mongo-driver/bson"
)

func GetProfile(w http.ResponseWriter, r *http.Request) {
	// username := r.Context().Value("username").(string)
	username, ok := r.Context().Value("username").(string)
	if !ok {
		utils.WriteJSON(w, 401, "Unauthorized", nil)
		return
	}

	var user models.User
	err := config.UserCollection.FindOne(
		context.Background(),
		bson.M{"username": username},
	).Decode(&user)

	if err != nil {
		utils.WriteJSON(w, 400, "User not found", nil)
		return
	}

	utils.WriteJSON(w, 200, "success", user)

	// json.NewEncoder(w).Encode(response)
}
