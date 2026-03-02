package handlers

import (
	"log"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"

	"backend-chat-go/config"
	"backend-chat-go/models"
	"backend-chat-go/services"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WsHandler(w http.ResponseWriter, r *http.Request) {

	tokenString := r.URL.Query().Get("token")
	if tokenString == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return config.JwtSecret, nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	claims := token.Claims.(jwt.MapClaims)
	username := claims["username"].(string)

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	services.AddClient(ws, username)

	defer func() {
		services.RemoveClient(ws)
		ws.Close()
	}()

	for {
		var msg models.Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			break
		}

		services.HandleMessage(ws, msg)
	}
}
