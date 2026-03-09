package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"backend-chat-go/config"
	"backend-chat-go/models"
	"backend-chat-go/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetChatList(w http.ResponseWriter, r *http.Request) {

	username := r.Context().Value("username")
	if username == nil {
		utils.WriteJSON(w, 401, "Unauthorized", nil)
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

	var chatList []models.ChatListResponse

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

		var lastMsg models.Message
		err := config.MessageCollection.FindOne(
			ctx,
			bson.M{"roomId": roomId},
			opt,
		).Decode(&lastMsg)

		if err != nil {
			continue
		}

		unreadCount, _ := config.MessageCollection.CountDocuments(
			ctx,
			bson.M{
				"roomId":   roomId,
				"receiver": username.(string),
				"isRead":   false,
			},
		)

		chatList = append(chatList, models.ChatListResponse{
			RoomId:          roomId,
			Friend:          friend,
			LastMessage:     lastMsg.Preview,
			LastMessageTime: lastMsg.CreatedAt,
			LastSender:      lastMsg.Sender,
			IsOnline:        IsUserOnline(friend),
			UnreadCount:     unreadCount,
		})

	}

	response := models.ApiResponse{
		StatusCode: 200,
		Message:    "success",
		Data:       chatList,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
