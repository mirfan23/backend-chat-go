package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"backend-chat-go/config"
	"backend-chat-go/models"
	"backend-chat-go/services"
	"backend-chat-go/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserResponse struct {
	UserId   string `json:"userId"`
	Username string `json:"username"`
	IsOnline bool   `json:"isOnline"`
}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {

	currentUser := r.Context().Value("userId")
	if currentUser == nil {
		utils.WriteJSON(w, 401, "Unauthorized", nil)
		return
	}

	ctx := context.Background()

	cursor, err := config.UserCollection.Find(ctx, bson.M{})
	if err != nil {
		utils.WriteJSON(w, 500, "Internal server error", nil)
		return
	}
	defer cursor.Close(ctx)

	var users []UserResponse

	for cursor.Next(ctx) {

		var user models.User

		if err := cursor.Decode(&user); err != nil {
			continue
		}

		userId := user.ID.Hex()

		// skip diri sendiri
		if userId == currentUser.(string) {
			continue
		}

		users = append(users, UserResponse{
			UserId:   userId,
			Username: user.Username,
			IsOnline: services.IsUserOnline(userId),
		})
	}

	response := models.ApiResponse{
		StatusCode: 200,
		Message:    "success",
		Data:       users,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetUserPublicKey(w http.ResponseWriter, r *http.Request) {

	userId := r.Context().Value("userId")
	if userId == nil {
		utils.WriteJSON(w, 401, "Unauthorized", nil)
		return
	}

	objId, err := primitive.ObjectIDFromHex(userId.(string))
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

	utils.WriteJSON(w, 200, "success", map[string]string{
		"publicKey": user.PublicKey,
	})
}
