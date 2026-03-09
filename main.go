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
	http.HandleFunc("/users", middleware.JWTMiddleware(handlers.GetAllUsers))
	http.HandleFunc("/markRead", middleware.JWTMiddleware(handlers.MarkMessagesRead))
	http.HandleFunc("/profile", middleware.JWTMiddleware(handlers.GetProfile))
	http.HandleFunc(
		"/chatList",
		middleware.RecoveryMiddleware(
			middleware.JWTMiddleware(
				handlers.GetChatList,
			),
		),
	)
	http.HandleFunc("/userPublicKey", middleware.JWTMiddleware(handlers.GetUserPublicKey))

	log.Println("🚀 Server running on :3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
