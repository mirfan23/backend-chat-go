package services

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"backend-chat-go/config"
	"backend-chat-go/models"
	"backend-chat-go/utils"
)

var Clients = make(map[*websocket.Conn]string) // conn -> userId
var Rooms = make(map[string]map[*websocket.Conn]bool)
var OnlineUsers = make(map[string]int) // userId -> connection count
var Mutex sync.RWMutex

// ================= CLIENT =================

func AddClient(ws *websocket.Conn, userId string) {
	Mutex.Lock()
	Clients[ws] = userId
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

	userId := Clients[ws]
	msg.Sender = userId

	switch msg.Type {

	case "joinRoom":
		JoinRoom(ws, msg.RoomID)

	case "sendMessage":

		msg.Type = "newMessage"
		msg.CreatedAt = time.Now()
		msg.Receiver = utils.ExtractFriend(msg.RoomID, userId)
		msg.IsRead = false

		SaveMessage(msg)
		BroadcastMessage(msg)
		BroadcastToFriendList(msg)

	case "chatList":
		log.Println("🔥 chatList requested by", userId)
		SendChatList(ws, userId)

	case "typing":
		log.Println("⌨️ typing event:", msg.Sender)
		BroadcastTyping(msg, ws)

	case "readMessage":
		HandleReadMessage(msg)

	}
}

// ================= HANDLE READ MESSAGE =================
func HandleReadMessage(msg models.Message) {

	ctx := context.Background()

	friend := utils.ExtractFriend(msg.RoomID, msg.Sender)

	_, err := config.MessageCollection.UpdateMany(
		ctx,
		bson.M{
			"roomId":   msg.RoomID,
			"sender":   friend,
			"receiver": msg.Sender,
			"isRead":   false,
		},
		bson.M{
			"$set": bson.M{"isRead": true},
		},
	)

	if err != nil {
		log.Println("❌ Read error:", err)
		return
	}

	payload := map[string]interface{}{
		"type":   "messagesRead",
		"roomId": msg.RoomID,
		"reader": msg.Sender,
	}

	BroadcastRoomPayLoad(msg.RoomID, payload)
	BroadcastToUser(msg.Sender, payload)
	BroadcastToUser(friend, payload)
}

func BroadcastToUser(userId string, payload map[string]interface{}) {

	Mutex.RLock()
	defer Mutex.RUnlock()

	for client, uid := range Clients {
		if uid == userId {
			client.WriteJSON(payload)
		}
	}
}

func BroadcastRoomPayLoad(roomId string, payload map[string]interface{}) {
	Mutex.RLock()
	room := Rooms[roomId]
	Mutex.RUnlock()

	if room == nil {
		return
	}

	for client := range room {
		client.WriteJSON(payload)
	}
}

// ================= BROADCAST CHAT LIST =================
func BroadcastChatList(userId string) {

	for client, uid := range Clients {
		if uid == userId {
			SendChatList(client, userId)
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

// ================= BROADCAST MESSAGE =================

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

func BroadcastUserStatus(userId string, isOnline bool) {

	Mutex.RLock()
	defer Mutex.RUnlock()

	message := map[string]interface{}{
		"type":     "user_status",
		"userId":   userId,
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

	for client, userId := range Clients {

		if userId == msg.Receiver {

			sendMsg := msg
			sendMsg.Type = "newMessage"

			client.WriteJSON(sendMsg)
		}
	}

	BroadcastChatList(msg.Sender)
	BroadcastChatList(msg.Receiver)
}

// ================= SEND CHAT LIST =================

func SendChatList(ws *websocket.Conn, userId string) {

	ctx := context.Background()

	// pipeline untuk aggregate last message + friend info + unread count
	pipeline := []bson.M{
		{"$match": bson.M{
			"$or": []bson.M{
				{"sender": userId},
				{"receiver": userId},
			},
		}},
		{"$sort": bson.M{"createdAt": -1}}, // urutkan newest first
		{"$group": bson.M{
			"_id":             "$roomId",
			"lastMessage":     bson.M{"$first": "$preview"},
			"lastMessageTime": bson.M{"$first": "$createdAt"},
			"lastSender":      bson.M{"$first": "$sender"},
			"lastIsRead":      bson.M{"$first": "$isRead"},
			"friendId": bson.M{"$first": bson.M{
				"$cond": []interface{}{
					bson.M{"$eq": []interface{}{"$sender", userId}},
					"$receiver",
					"$sender",
				},
			}},
		}},
	}

	cursor, err := config.MessageCollection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Println("❌ chat list aggregation error:", err)
		return
	}
	defer cursor.Close(ctx)

	var chatList []models.ChatListResponse

	for cursor.Next(ctx) {
		var doc struct {
			RoomID          string             `bson:"_id"`
			LastMessage     string             `bson:"lastMessage"`
			LastMessageTime primitive.DateTime `bson:"lastMessageTime"`
			LastSender      string             `bson:"lastSender"`
			LastIsRead      bool               `bson:"lastIsRead"`
			FriendID        string             `bson:"friendId"`
		}

		if err := cursor.Decode(&doc); err != nil {
			log.Println("❌ decode chat doc error:", err)
			continue
		}

		// ambil username friend
		friendName := doc.FriendID
		friendObjectID, err := primitive.ObjectIDFromHex(doc.FriendID)
		if err == nil {
			var user models.User
			if err := config.UserCollection.FindOne(ctx, bson.M{"_id": friendObjectID}).Decode(&user); err == nil {
				friendName = user.Username
			}
		}

		// hitung unread messages
		unreadCount, _ := config.MessageCollection.CountDocuments(ctx, bson.M{
			"roomId":   doc.RoomID,
			"receiver": userId,
			"isRead":   false,
		})

		isRead := true

		if doc.LastSender == userId {
			isRead = doc.LastIsRead
		}

		chatList = append(chatList, models.ChatListResponse{
			RoomId:          doc.RoomID,
			Friend:          doc.FriendID,
			FriendName:      friendName,
			LastMessage:     doc.LastMessage,
			LastMessageTime: doc.LastMessageTime.Time(),
			LastSender:      doc.LastSender,
			IsOnline:        IsUserOnline(doc.FriendID),
			UnreadCount:     unreadCount,
			IsRead:          isRead,
		})
	}

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
