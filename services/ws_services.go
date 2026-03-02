package services

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"

	"backend-chat-go/config"
	"backend-chat-go/models"
)

var Clients = make(map[*websocket.Conn]string)
var Rooms = make(map[string]map[*websocket.Conn]bool)
var Mutex sync.RWMutex

// =======================
// CLIENT MANAGEMENT
// =======================

func AddClient(ws *websocket.Conn, username string) {
	Mutex.Lock()
	Clients[ws] = username
	Mutex.Unlock()

	SendUserList()
}

func RemoveClient(ws *websocket.Conn) {
	Mutex.Lock()

	delete(Clients, ws)

	for roomID := range Rooms {
		delete(Rooms[roomID], ws)

		if len(Rooms[roomID]) == 0 {
			delete(Rooms, roomID)
		}
	}

	Mutex.Unlock()

	SendUserList()
}

// =======================
// ROOM MANAGEMENT
// =======================

func JoinRoom(ws *websocket.Conn, roomID string) {
	Mutex.Lock()
	defer Mutex.Unlock()

	if Rooms[roomID] == nil {
		Rooms[roomID] = make(map[*websocket.Conn]bool)
	}

	Rooms[roomID][ws] = true
}

// =======================
// MESSAGE
// =======================

func HandleMessage(ws *websocket.Conn, msg models.Message) {

	username := Clients[ws]
	msg.Sender = username

	switch msg.Type {

	case "joinRoom":
		JoinRoom(ws, msg.RoomID)

	case "sendMessage":
		msg.CreatedAt = time.Now()
		SaveMessage(msg)
		BroadcastToRoom(msg)

	case "typing":
		BroadcastTyping(ws, msg)
	}
}

func SaveMessage(msg models.Message) {

	message := bson.M{
		"roomId":    msg.RoomID,
		"sender":    msg.Sender,
		"text":      msg.Text,
		"createdAt": msg.CreatedAt,
	}

	_, err := config.MessageCollection.InsertOne(context.Background(), message)
	if err != nil {
		log.Println("❌ Failed save message:", err)
	}
}

func BroadcastToRoom(msg models.Message) {

	Mutex.RLock()
	room := Rooms[msg.RoomID]

	var clientsInRoom []*websocket.Conn
	for client := range room {
		clientsInRoom = append(clientsInRoom, client)
	}
	Mutex.RUnlock()

	if room == nil {
		return
	}

	msg.Type = "newMessage"

	for _, client := range clientsInRoom {
		err := client.WriteJSON(msg)
		if err != nil {
			log.Println("Broadcast error:", err)
		}
	}
}

// =======================
// USER LIST
// =======================

func SendUserList() {

	Mutex.RLock()

	var userList []string
	var clientList []*websocket.Conn

	for client, username := range Clients {
		userList = append(userList, username)
		clientList = append(clientList, client)
	}

	Mutex.RUnlock()

	data, _ := json.Marshal(map[string]interface{}{
		"type":  "userList",
		"users": userList,
	})

	for _, client := range clientList {
		err := client.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			log.Println("UserList send error:", err)
		}
	}
}

// =======================
// TYPING
// =======================

func BroadcastTyping(ws *websocket.Conn, msg models.Message) {

	Mutex.RLock()
	room := Rooms[msg.RoomID]

	var clientsInRoom []*websocket.Conn
	for client := range room {
		if client != ws {
			clientsInRoom = append(clientsInRoom, client)
		}
	}
	Mutex.RUnlock()

	for _, client := range clientsInRoom {
		err := client.WriteJSON(msg)
		if err != nil {
			log.Println("Typing broadcast error:", err)
		}
	}
}
