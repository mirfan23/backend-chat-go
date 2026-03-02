package main

import (
	"log"
	"net/http"

	"backend-chat-go/config"
	"backend-chat-go/handlers"
	"backend-chat-go/middleware"
)

func main() {

	config.InitMongo()

	http.HandleFunc("/register", handlers.Register)
	http.HandleFunc("/login", handlers.Login)
	http.HandleFunc("/messages", middleware.JWTMiddleware(handlers.GetMessages))
	http.HandleFunc("/ws", handlers.WsHandler)

	log.Println("🚀 Server running on :3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}

// package main

// import (
// 	"context"
// 	"encoding/json"

// 	// "go/token"
// 	"log"
// 	// "maps"
// 	"net/http"
// 	"sync"
// 	"time"

// 	"github.com/golang-jwt/jwt/v5"
// 	"github.com/gorilla/websocket"
// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/mongo"
// 	"go.mongodb.org/mongo-driver/mongo/options"
// 	"golang.org/x/crypto/bcrypt"
// )

// var jwtSecret = []byte("secret")

// var mongoClient *mongo.Client
// var userCollection *mongo.Collection
// var messageCollection *mongo.Collection

// var upgrader = websocket.Upgrader{
// 	CheckOrigin: func(r *http.Request) bool {
// 		return true
// 	},
// }

// var clients = make(map[*websocket.Conn]string)
// var rooms = make(map[string]map[*websocket.Conn]bool)
// var mutex = &sync.Mutex{}

// type Message struct {
// 	Type     string `json:"type"`
// 	RoomID   string `json:"roomId,omitempty"`
// 	Sender   string `json:"sender,omitempty"`
// 	Text     string `json:"text,omitempty"`
// 	IsTyping bool   `json:"isTyping,omitempty"`
// }

// type ApiResponse struct {
// 	Success bool        `json:"success"`
// 	Message string      `json:"message,omitempty"`
// 	Data    interface{} `json:"data,omitempty"`
// }

// func main() {
// 	initMongo()

// 	http.HandleFunc("/register", registerHandler)
// 	http.HandleFunc("/login", loginHandler)
// 	http.HandleFunc("/messages", jwtMiddleware(getMessagesHandler))
// 	http.HandleFunc("/ws", wsHandler)

// 	log.Println("🚀 Server running on :3000")
// 	log.Fatal(http.ListenAndServe(":3000", nil))
// }

// func initMongo() {
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
// 	if err != nil {
// 		panic(err)
// 	}

// 	mongoClient = client
// 	userCollection = client.Database("chat").Collection("users")
// 	messageCollection = client.Database("chat").Collection("messages")

// 	log.Println("✅ MongoDB connected")
// }

// func writeJSON(w http.ResponseWriter, status int, response ApiResponse) {
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(status)
// 	json.NewEncoder(w).Encode(response)
// }

// func registerHandler(w http.ResponseWriter, r *http.Request) {

// 	var data map[string]string
// 	json.NewDecoder(r.Body).Decode(&data)

// 	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), bcrypt.DefaultCost)

// 	user := bson.M{
// 		"username":  data["username"],
// 		"password":  string(hashedPassword),
// 		"createdAt": time.Now(),
// 	}

// 	_, err := userCollection.InsertOne(context.Background(), user)
// 	if err != nil {
// 		writeJSON(w, 400, ApiResponse{
// 			Success: false,
// 			Message: "User already exists",
// 		})
// 		return
// 	}

// 	writeJSON(w, 200, ApiResponse{
// 		Success: true,
// 		Message: "Register success",
// 	})
// }

// func loginHandler(w http.ResponseWriter, r *http.Request) {

// 	var data map[string]string
// 	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
// 		writeJSON(w, 400, ApiResponse{
// 			Success: false,
// 			Message: "Invalid request body",
// 		})
// 		return
// 	}

// 	var user bson.M
// 	err := userCollection.FindOne(context.Background(),
// 		bson.M{"username": data["username"]}).Decode(&user)

// 	if err != nil {
// 		writeJSON(w, 400, ApiResponse{
// 			Success: false,
// 			Message: "User not found",
// 		})
// 		return
// 	}

// 	err = bcrypt.CompareHashAndPassword(
// 		[]byte(user["password"].(string)),
// 		[]byte(data["password"]),
// 	)

// 	if err != nil {
// 		writeJSON(w, 400, ApiResponse{
// 			Success: false,
// 			Message: "Wrong password",
// 		})
// 		return
// 	}

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
// 		"username": data["username"],
// 		"exp":      time.Now().Add(24 * time.Hour).Unix(),
// 	})

// 	tokenString, _ := token.SignedString(jwtSecret)

// 	writeJSON(w, 200, ApiResponse{
// 		Success: true,
// 		Data: map[string]string{
// 			"token": tokenString,
// 		},
// 	})
// }

// func jwtMiddleware(next http.HandlerFunc) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		tokenString := r.Header.Get("Authorization")
// 		if tokenString == "" {
// 			http.Error(w, "Unauthorized", 401)
// 			return
// 		}

// 		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 			return jwtSecret, nil
// 		})

// 		if err != nil || !token.Valid {
// 			http.Error(w, "Invalid token", 401)
// 			return
// 		}

// 		next(w, r)
// 	}
// }

// func getMessagesHandler(w http.ResponseWriter, r *http.Request) {

// 	roomID := r.URL.Query().Get("roomId")

// 	cursor, err := messageCollection.Find(
// 		context.Background(),
// 		bson.M{"roomId": roomID},
// 		options.Find().SetSort(bson.M{"createdAt": 1}),
// 	)

// 	if err != nil {
// 		writeJSON(w, 500, ApiResponse{
// 			Success: false,
// 			Message: "Failed to fetch messages",
// 		})
// 		return
// 	}

// 	var messages []bson.M
// 	cursor.All(context.Background(), &messages)

// 	writeJSON(w, 200, ApiResponse{
// 		Success: true,
// 		Data:    messages,
// 	})
// }
