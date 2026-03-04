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
	http.HandleFunc("/markRead", middleware.JWTMiddleware(handlers.MarkAsRead))
	http.HandleFunc(
		"/chatList",
		middleware.RecoveryMiddleware(
			middleware.JWTMiddleware(
				handlers.GetChatList,
			),
		),
	)

	log.Println("🚀 Server running on :3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
