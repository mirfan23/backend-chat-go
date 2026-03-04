package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"backend-chat-go/config"
	"backend-chat-go/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ChatListResponse struct {
	RoomId          string    `json:"roomId"`
	Friend          string    `json:"friend"`
	LastMessage     string    `json:"lastMessage"`
	LastMessageTime time.Time `json:"lastMessageTime"`
	LastSender      string    `json:"lastSender"`
	IsOnline        bool      `json:"isOnline"`
}

func GetChatList(w http.ResponseWriter, r *http.Request) {

	username := r.Context().Value("username")
	if username == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	ctx := context.Background()

	filter := bson.M{
		"roomId": bson.M{
			"$regex": username.(string),
		},
	}

	roomIds, err := config.MessageCollection.Distinct(ctx, "roomId", filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var chatList []ChatListResponse

	for _, rId := range roomIds {
		roomId := rId.(string)

		parts := strings.Split(roomId, "_")
		var friend string

		if parts[0] == username.(string) {
			friend = parts[1]
		} else {
			friend = parts[0]
		}

		opt := options.FindOne().SetSort(bson.D{{Key: "createdAt", Value: -1}})

		var lastMsg models.Chat
		err := config.MessageCollection.FindOne(
			ctx,
			bson.M{"roomId": roomId},
			opt,
		).Decode(&lastMsg)

		if err != nil {
			continue
		}

		chatList = append(chatList, ChatListResponse{
			RoomId:          roomId,
			Friend:          friend,
			LastMessage:     lastMsg.Text,
			LastMessageTime: lastMsg.CreatedAt.Time(),
			LastSender:      lastMsg.Sender,
			IsOnline:        IsUserOnline(friend),
		})
	}

	response := map[string]interface{}{
		"status":  200,
		"message": "Success",
		"data":    chatList,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
