package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"backend-chat-go/config"
	"backend-chat-go/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var MessageCollection *mongo.Collection

func GetChatList(w http.ResponseWriter, r *http.Request) {

	usernameVal := r.Context().Value("username")
	if usernameVal == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	username, ok := usernameVal.(string)
	if !ok {
		http.Error(w, "Invalid user context", http.StatusUnauthorized)
		return
	}

	if config.MessageCollection == nil {
		http.Error(w, "Database not initialized", http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"$or": []bson.M{
			{"sender": username},
			{"receiver": username},
		},
	}

	cursor, err := config.MessageCollection.Find(ctx, filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var chats []models.Chat
	if err := cursor.All(ctx, &chats); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type ChatResponse struct {
		ID        string `json:"id"`
		Sender    string `json:"sender"`
		Receiver  string `json:"receiver"`
		Message   string `json:"message"`
		IsRead    bool   `json:"is_read"`
		CreatedAt string `json:"created_at"`
	}

	var response []ChatResponse

	for _, chat := range chats {

		isRead := false
		if chat.IsRead != nil {
			isRead = *chat.IsRead
		}

		response = append(response, ChatResponse{
			ID:        chat.ID.Hex(),
			Sender:    chat.Sender,
			Receiver:  chat.Receiver,
			Message:   chat.Message,
			IsRead:    isRead,
			CreatedAt: chat.CreatedAt.Time().Format(time.RFC3339),
		})
	}

	// 🔥 Bungkus dalam JSON object
	result := map[string]interface{}{
		"status":  200,
		"message": "Success",
		"data":    response,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}
