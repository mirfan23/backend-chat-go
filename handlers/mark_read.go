package handlers

import (
	"context"
	"net/http"

	"backend-chat-go/config"
	"backend-chat-go/utils"

	"go.mongodb.org/mongo-driver/bson"
)

func MarkMessagesRead(w http.ResponseWriter, r *http.Request) {

	username := r.Context().Value("username")
	roomId := r.URL.Query().Get("roomId")

	ctx := context.Background()

	_, err := config.MessageCollection.UpdateMany(
		ctx,
		bson.M{
			"roomId":   roomId,
			"receiver": username.(string),
			"isRead":   false,
		},
		bson.M{
			"$set": bson.M{
				"isRead": true,
			},
		},
	)

	if err != nil {
		utils.WriteJSON(w, 500, "Failed update", nil)
		return
	}

	utils.WriteJSON(w, 200, "Updated", nil)
}
