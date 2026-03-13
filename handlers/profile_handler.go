package handlers

import (
	"context"
	"net/http"

	"backend-chat-go/config"
	"backend-chat-go/models"
	"backend-chat-go/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetProfile(w http.ResponseWriter, r *http.Request) {

	userId, ok := r.Context().Value("userId").(string)
	if !ok {
		utils.WriteJSON(w, 401, "Unauthorized", nil)
		return
	}

	objId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		utils.WriteJSON(w, 400, "Invalid userId", nil)
		return
	}

	var user models.User

	err = config.UserCollection.FindOne(
		context.Background(),
		bson.M{"_id": objId},
	).Decode(&user)

	if err != nil {
		utils.WriteJSON(w, 404, "User not found", nil)
		return
	}

	utils.WriteJSON(w, 200, "success", user)
}
