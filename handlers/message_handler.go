package handlers

import (
	"context"
	"net/http"

	"backend-chat-go/config"
	"backend-chat-go/models"
	"backend-chat-go/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetMessages(w http.ResponseWriter, r *http.Request) {

	roomID := r.URL.Query().Get("roomId")

	cursor, err := config.MessageCollection.Find(
		context.Background(),
		bson.M{"roomId": roomID},
		options.Find().SetSort(bson.M{"createdAt": 1}),
	)

	if err != nil {
		utils.WriteJSON(w, 500, "Failed fetch messages", nil)
		return
	}

	var messages []models.Message
	cursor.All(context.Background(), &messages)

	utils.WriteJSON(w, 200, "Success", messages)
}
