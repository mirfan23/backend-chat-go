package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Message struct {
	Type   string `json:"type"`
	RoomID string `json:"roomId,omitempty"`
	Sender string `json:"sender,omitempty"`
	Text   string `json:"text,omitempty"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients = make(map[*websocket.Conn]string)

// 🔥 rooms pakai map supaya tidak duplicate
var rooms = make(map[string]map[*websocket.Conn]bool)

var mutex = &sync.Mutex{}

func main() {
	http.HandleFunc("/ws", handleConnections)

	log.Println("🚀 Server running on :3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	log.Println("🔥 New client connected")

	defer func() {
		removeClient(ws)
		ws.Close()
		log.Println("❌ Client disconnected")
	}()

	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			break
		}

		log.Println("📩 Received:", msg)

		switch msg.Type {

		case "login":
			handleLogin(ws, msg)

		case "joinRoom":
			handleJoinRoom(ws, msg)

		case "sendMessage":
			handleSendMessage(ws, msg)
		}
	}
}

func handleLogin(ws *websocket.Conn, msg Message) {
	mutex.Lock()
	clients[ws] = msg.Sender
	mutex.Unlock()

	sendUserList()
}

func handleJoinRoom(ws *websocket.Conn, msg Message) {
	mutex.Lock()

	if rooms[msg.RoomID] == nil {
		rooms[msg.RoomID] = make(map[*websocket.Conn]bool)
	}

	rooms[msg.RoomID][ws] = true

	log.Println("👥 Room:", msg.RoomID, "Members:", len(rooms[msg.RoomID]))

	mutex.Unlock()
}

func handleSendMessage(ws *websocket.Conn, msg Message) {
	mutex.Lock()

	// 🔥 pastikan room ada
	if rooms[msg.RoomID] == nil {
		rooms[msg.RoomID] = make(map[*websocket.Conn]bool)
	}

	// 🔥 pastikan sender ada di room
	rooms[msg.RoomID][ws] = true

	room := rooms[msg.RoomID]

	mutex.Unlock()

	msg.Type = "newMessage"

	data, _ := json.Marshal(msg)

	log.Println("📢 Broadcasting to room:", msg.RoomID, "Members:", len(room))

	for client := range room {
		err := client.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			log.Println("Write error:", err)
		}
	}
}

func removeClient(ws *websocket.Conn) {
	mutex.Lock()

	// hapus dari clients
	delete(clients, ws)

	// hapus dari semua rooms
	for roomID := range rooms {
		delete(rooms[roomID], ws)

		if len(rooms[roomID]) == 0 {
			delete(rooms, roomID)
		}
	}

	mutex.Unlock()

	sendUserList()
}

func sendUserList() {
	mutex.Lock()

	var userList []string
	var clientList []*websocket.Conn

	for client, username := range clients {
		userList = append(userList, username)
		clientList = append(clientList, client)
	}

	mutex.Unlock()

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
