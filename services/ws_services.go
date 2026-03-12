package services

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"backend-chat-go/config"
	"backend-chat-go/models"
	"backend-chat-go/utils"
)

var Clients = make(map[*websocket.Conn]string)
var Rooms = make(map[string]map[*websocket.Conn]bool)
var OnlineUsers = make(map[string]int)
var Mutex sync.RWMutex

// ================= CLIENT =================

func AddClient(ws *websocket.Conn, username string) {
	Mutex.Lock()
	Clients[ws] = username
	Mutex.Unlock()
}

func RemoveClient(ws *websocket.Conn) {
	Mutex.Lock()
	delete(Clients, ws)

	for roomID := range Rooms {
		delete(Rooms[roomID], ws)
	}
	Mutex.Unlock()
}

// ================= ROOM =================

func JoinRoom(ws *websocket.Conn, roomID string) {
	Mutex.Lock()
	if Rooms[roomID] == nil {
		Rooms[roomID] = make(map[*websocket.Conn]bool)
	}
	Rooms[roomID][ws] = true
	Mutex.Unlock()
}

// ================= HANDLE MESSAGE =================

func HandleMessage(ws *websocket.Conn, msg models.Message) {

	log.Println("📩 WS message:", msg.Type)

	username := Clients[ws]
	msg.Sender = username

	switch msg.Type {

	case "joinRoom":
		JoinRoom(ws, msg.RoomID)

	case "sendMessage":

		msg.Type = "newMessage"
		msg.CreatedAt = time.Now()
		msg.Receiver = utils.ExtractFriend(msg.RoomID, username)
		msg.IsRead = false

		SaveMessage(msg)
		BroadcastMessage(msg)
		BroadcastToFriendList(msg)

	case "chatList":
		log.Println("🔥 chatList requested by", username)
		SendChatList(ws, username)

	case "typing":
		log.Println("⌨️ typing event:", msg.Sender)
		BroadcastTyping(msg, ws)
	}

}

// ================= BROADCAST CHAT LIST =================

func BroadcastChatList(username string) {

	for client, user := range Clients {
		if user == username {
			SendChatList(client, username)
		}
	}
}

// ================= SAVE =================

func SaveMessage(msg models.Message) {

	_, err := config.MessageCollection.InsertOne(
		context.Background(),
		msg,
	)

	if err != nil {
		log.Println("❌ Save error:", err)
	}
}

// ================= BROADCAST =================

func BroadcastMessage(msg models.Message) {

	Mutex.RLock()
	room := Rooms[msg.RoomID]
	Mutex.RUnlock()

	if room == nil {
		return
	}

	for client := range room {
		client.WriteJSON(msg)
	}
}

// ================= BROADCAST USER STATUS =================
func BroadcastUserStatus(username string, isOnline bool) {

	Mutex.RLock()
	defer Mutex.RUnlock()

	message := map[string]interface{}{
		"type":     "user_status",
		"username": username,
		"isOnline": isOnline,
	}

	for client := range Clients {
		client.WriteJSON(message)
	}
}

// ================= SEND ONLINE USERS =================
func SendOnlineUsers(ws *websocket.Conn) {

	Mutex.RLock()
	defer Mutex.RUnlock()

	var users []string
	for user := range OnlineUsers {
		users = append(users, user)
	}

	ws.WriteJSON(map[string]interface{}{
		"type":  "online_users",
		"users": users,
	})
}

// ================= BROADCAST TO FRIEND LIST =================
func BroadcastToFriendList(msg models.Message) {
	Mutex.RLock()
	defer Mutex.RUnlock()

	for client, username := range Clients {
		if username == msg.Receiver {
			sendMsg := msg
			sendMsg.Type = "newMessage"
			client.WriteJSON(sendMsg)
		}
	}
	BroadcastChatList(msg.Sender)
	BroadcastChatList(msg.Receiver)
}

// ================= SEND CHAT LIST =================
func SendChatList(ws *websocket.Conn, username string) {
	ctx := context.Background()

	filter := bson.M{
		"roomId": bson.M{
			"$regex": username,
		},
	}

	roomIds, err := config.MessageCollection.Distinct(ctx, "roomId", filter)
	if err != nil {
		return
	}

	var chatList []models.ChatListResponse

	for _, rId := range roomIds {
		roomId := rId.(string)

		opts := options.FindOne().
			SetSort(bson.D{{Key: "createdAt", Value: -1}})

		var lastMsg models.Message

		err := config.MessageCollection.FindOne(ctx, bson.M{"roomId": roomId}, opts).Decode(&lastMsg)

		if err != nil {
			continue
		}

		friend := utils.GetFriendFromRoom(roomId, username)

		unreadFilter := bson.M{
			"roomId":   roomId,
			"receiver": username,
			"isRead":   false,
		}

		unreadCount, _ := config.MessageCollection.CountDocuments(ctx, unreadFilter)

		chatList = append(chatList, models.ChatListResponse{
			RoomId:          roomId,
			Friend:          friend,
			LastMessage:     lastMsg.Preview,
			LastMessageTime: lastMsg.CreatedAt,
			LastSender:      lastMsg.Sender,
			IsOnline:        IsUserOnline(friend),
			UnreadCount:     unreadCount,
		})
		log.Println("roomId:", roomId)
		log.Println("username:", username)
		log.Println("friend:", friend)
	}

	log.Println("📤 sending chat list:", len(chatList))

	ws.WriteJSON(map[string]interface{}{
		"type": "chatList",
		"data": chatList,
	})
}

// ================= TYPING =================
func BroadcastTyping(msg models.Message, senderConn *websocket.Conn) {

	Mutex.RLock()
	room := Rooms[msg.RoomID]
	Mutex.RUnlock()

	if room == nil {
		return
	}

	typingPayload := map[string]interface{}{
		"type":     "typing",
		"roomId":   msg.RoomID,
		"sender":   msg.Sender,
		"isTyping": msg.IsTyping,
	}

	for client := range room {

		// jangan kirim ke pengirim
		if client == senderConn {
			continue
		}

		client.WriteJSON(typingPayload)
	}
}
