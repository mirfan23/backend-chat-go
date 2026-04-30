package main

import (
	"log"

	"backend-chat-go/config"
	"backend-chat-go/handlers"
	"backend-chat-go/middleware"

	"github.com/gin-gonic/gin"
)

func main() {

	config.InitMongo()

	r := gin.Default()

	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)

	r.GET("/ws", handlers.WsHandler)

	auth := r.Group("/")
	auth.Use(middleware.JWTMiddleware())
	{
		auth.GET("/profile", handlers.GetProfile)
		auth.POST("/updateProfile", handlers.UpdateProfile)
		auth.GET("/messages", handlers.GetMessages)
		auth.GET("/users", handlers.GetAllUsers)
		auth.GET("/userPublicKey", handlers.GetUserPublicKey)
		auth.GET("/markRead", handlers.MarkMessagesRead)
	}

	// http.HandleFunc("/ws", handlers.WsHandler)

	log.Println("🚀 Server running on :3000")
	log.Fatal(r.Run("0.0.0.0:3000"))
}
