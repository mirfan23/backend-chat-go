package handlers

import (
	"context"
	"net/http"
	"strings"
	"time"

	"backend-chat-go/config"
	"backend-chat-go/utils"

	"go.mongodb.org/mongo-driver/bson"
)

type ChatList struct {
	RoomID          string    `json:"roomId"`
	Friend          string    `json:"friend"`
	LastMessage     string    `json:"lastMessage"`
	LastSender      string    `json:"lastSender"`
	LastMessageTime time.Time `json:"lastMessageTime"`
	IsRead          bool      `json:"isRead"`
	UnreadCount     int64     `json:"unreadCount"`
}

func GetChatList(w http.ResponseWriter, r *http.Request) {

	user := r.Context().Value("username").(string)

	pipeline := []bson.M{
		{"$sort": bson.M{"createdAt": -1}},
		{"$group": bson.M{
			"_id":           "$roomId",
			"lastMessage":   bson.M{"$first": "$text"},
			"lastSender":    bson.M{"$first": "$sender"},
			"lastCreatedAt": bson.M{"$first": "$createdAt"},
			"isRead":        bson.M{"$first": "$isRead"},
		}},
	}

	cursor, _ := config.MessageCollection.Aggregate(context.Background(), pipeline)

	var results []bson.M
	cursor.All(context.Background(), &results)

	var chats []ChatList

	for _, item := range results {

		roomID := item["_id"].(string)

		if !strings.Contains(roomID, user) {
			continue
		}

		friend := utils.ExtractFriend(roomID, user)

		count, _ := config.MessageCollection.CountDocuments(
			context.Background(),
			bson.M{
				"roomId":   roomID,
				"receiver": user,
				"isRead":   false,
			},
		)

		chats = append(chats, ChatList{
			RoomID:          roomID,
			Friend:          friend,
			LastMessage:     item["lastMessage"].(string),
			LastSender:      item["lastSender"].(string),
			LastMessageTime: item["lastCreatedAt"].(time.Time),
			IsRead:          item["isRead"].(bool),
			UnreadCount:     count,
		})
	}

	writeJSON(w, 200, "Success", chats)
}
