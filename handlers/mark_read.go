package handlers

import (
	"context"
	"net/http"

	"backend-chat-go/config"

	"go.mongodb.org/mongo-driver/bson"
)

// PUT /mark-read?roomId=xxx
func MarkAsRead(w http.ResponseWriter, r *http.Request) {

	roomID := r.URL.Query().Get("roomId")
	if roomID == "" {
		writeJSON(w, 400, "roomId required", nil)
		return
	}

	// username dari JWT middleware
	user := r.Context().Value("username").(string)

	filter := bson.M{
		"roomId":   roomID,
		"receiver": user,
		"isRead":   false,
	}

	update := bson.M{
		"$set": bson.M{
			"isRead": true,
		},
	}

	_, err := config.MessageCollection.UpdateMany(
		context.Background(),
		filter,
		update,
	)

	if err != nil {
		writeJSON(w, 500, "Failed update", nil)
		return
	}

	writeJSON(w, 200, "Updated", nil)
}
