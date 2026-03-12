package handlers

import (
	"context"
	"net/http"

	"backend-chat-go/config"
	"backend-chat-go/utils"

	"go.mongodb.org/mongo-driver/bson"
)

func MarkMessagesRead(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		utils.WriteJSON(w, 405, "Method Not Allowed", nil)
		return
	}

	username := r.Context().Value("username")
	roomId := r.URL.Query().Get("roomId")

	ctx := context.Background()

	filter := bson.M{
		"roomId":   roomId,
		"receiver": username,
		"isRead":   false,
	}

	update := bson.M{
		"$set": bson.M{
			"isRead": true,
		},
	}

	_, err := config.MessageCollection.UpdateMany(ctx, filter, update)

	if err != nil {
		utils.WriteJSON(w, 500, "Failed update", nil)
		return
	}

	utils.WriteJSON(w, 200, "Updated", nil)
}
