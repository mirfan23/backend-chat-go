package services

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"

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

	username := Clients[ws]
	msg.Sender = username

	switch msg.Type {

	case "joinRoom":
		JoinRoom(ws, msg.RoomID)

	case "sendMessage":
		msg.CreatedAt = time.Now()
		msg.Receiver = utils.ExtractFriend(msg.RoomID, username)
		msg.IsRead = false

		SaveMessage(msg)
		BroadcastToRoom(msg)
		BroadcastToFriendList(msg)
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

func BroadcastToRoom(msg models.Message) {

	Mutex.RLock()
	room := Rooms[msg.RoomID]
	Mutex.RUnlock()

	if room == nil {
		return
	}

	msg.Type = "newMessage"

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
			// tambahkan type: newMessage
			sendMsg := msg
			sendMsg.Type = "newMessage"
			client.WriteJSON(sendMsg)
		}
	}
}
