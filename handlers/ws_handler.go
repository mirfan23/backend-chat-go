// package handlers

// import (
// 	"log"
// 	"net/http"

// 	"backend-chat-go/config"
// 	"backend-chat-go/models"
// 	"backend-chat-go/services"
// 	"backend-chat-go/utils"

// 	"github.com/golang-jwt/jwt/v5"
// 	"github.com/gorilla/websocket"
// )

// var upgrader = websocket.Upgrader{
// 	CheckOrigin: func(r *http.Request) bool {
// 		return true
// 	},
// }

// func WsHandler(w http.ResponseWriter, r *http.Request) {

// 	tokenString := r.URL.Query().Get("token")
// 	if tokenString == "" {
// 		utils.WriteJSON(w, 401, "Unauthorized", nil)
// 		return
// 	}

// 	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 		return config.JwtSecret, nil
// 	})

// 	if err != nil || !token.Valid {
// 		utils.WriteJSON(w, 401, "Invalid Token", nil)
// 		return
// 	}

// 	// safer claims parsing
// 	claims, ok := token.Claims.(jwt.MapClaims)
// 	if !ok {
// 		utils.WriteJSON(w, 401, "Invalid token claims", nil)
// 		return
// 	}

// 	userId, ok := claims["userId"].(string)
// 	if !ok {
// 		utils.WriteJSON(w, 401, "Invalid userId in token", nil)
// 		return
// 	}

// 	username, ok := claims["username"].(string)
// 	if !ok {
// 		utils.WriteJSON(w, 401, "Invalid username in token", nil)
// 		return
// 	}

// 	ws, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		log.Println("Upgrade error:", err)
// 		return
// 	}

// 	log.Println("User connected:", username, "| ID:", userId)

// 	// register client using userId
// 	services.AddClient(ws, userId)

// 	services.Mutex.Lock()
// 	services.OnlineUsers[userId]++
// 	services.Mutex.Unlock()

// 	// broadcast user online
// 	services.BroadcastUserStatus(userId, true)

// 	// send current online users
// 	services.SendOnlineUsers(ws)

// 	defer func() {

// 		services.RemoveClient(ws)

// 		services.Mutex.Lock()
// 		services.OnlineUsers[userId]--

// 		if services.OnlineUsers[userId] <= 0 {
// 			delete(services.OnlineUsers, userId)
// 			services.Mutex.Unlock()

// 			// broadcast offline
// 			services.BroadcastUserStatus(userId, false)
// 		} else {
// 			services.Mutex.Unlock()
// 		}

// 		log.Println("User disconnected:", username, "| ID:", userId)

// 		ws.Close()
// 	}()

// 	for {
// 		var msg models.Message

// 		err := ws.ReadJSON(&msg)
// 		if err != nil {
// 			log.Println("Read error:", err)
// 			break
// 		}

// 		services.HandleMessage(ws, msg)
// 	}
// }

package handlers

import (
	"log"
	"net/http"

	"backend-chat-go/config"
	"backend-chat-go/models"
	"backend-chat-go/services"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WsHandler(c *gin.Context) {

	// ambil token dari query
	tokenString := c.Query("token")

	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
		})
		return
	}

	// parse jwt
	token, err := jwt.Parse(
		tokenString,
		func(token *jwt.Token) (interface{}, error) {
			return config.JwtSecret, nil
		},
	)

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid Token",
		})
		return
	}

	// claims parsing
	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid token claims",
		})
		return
	}

	userId, ok := claims["userId"].(string)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid userId in token",
		})
		return
	}

	username, ok := claims["username"].(string)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid username in token",
		})
		return
	}

	// websocket upgrade
	ws, err := upgrader.Upgrade(
		c.Writer,
		c.Request,
		nil,
	)

	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	log.Println("User connected:", username, "| ID:", userId)

	// register client
	services.AddClient(ws, userId)

	// online counter
	services.Mutex.Lock()
	services.OnlineUsers[userId]++
	services.Mutex.Unlock()

	// broadcast online
	services.BroadcastUserStatus(userId, true)

	// kirim online users
	services.SendOnlineUsers(ws)

	defer func() {

		services.RemoveClient(ws)

		services.Mutex.Lock()

		services.OnlineUsers[userId]--

		if services.OnlineUsers[userId] <= 0 {

			delete(services.OnlineUsers, userId)

			services.Mutex.Unlock()

			// broadcast offline
			services.BroadcastUserStatus(userId, false)

		} else {
			services.Mutex.Unlock()
		}

		log.Println("User disconnected:", username, "| ID:", userId)

		ws.Close()
	}()

	for {

		var msg models.Message

		err := ws.ReadJSON(&msg)

		if err != nil {
			log.Println("Read error:", err)
			break
		}

		services.HandleMessage(ws, msg)
	}
}
