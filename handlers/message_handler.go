// package handlers

// import (
// 	"context"
// 	"net/http"
// 	"strings"

// 	"backend-chat-go/config"
// 	"backend-chat-go/models"
// 	"backend-chat-go/utils"

// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/mongo/options"
// )

// func GetMessages(w http.ResponseWriter, r *http.Request) {

// 	userId := r.Context().Value("userId")
// 	if userId == nil {
// 		utils.WriteJSON(w, 401, "Unauthorized", nil)
// 		return
// 	}

// 	roomID := r.URL.Query().Get("roomId")
// 	if roomID == "" {
// 		utils.WriteJSON(w, 400, "roomId required", nil)
// 		return
// 	}

// 	// 🔒 pastikan user bagian dari room
// 	parts := strings.Split(roomID, "_")
// 	if len(parts) != 2 {
// 		utils.WriteJSON(w, 400, "Invalid roomId", nil)
// 		return
// 	}

// 	if parts[0] != userId.(string) && parts[1] != userId.(string) {
// 		utils.WriteJSON(w, 403, "Forbidden", nil)
// 		return
// 	}

// 	cursor, err := config.MessageCollection.Find(
// 		context.Background(),
// 		bson.M{"roomId": roomID},
// 		options.Find().SetSort(bson.M{"createdAt": 1}),
// 	)

// 	if err != nil {
// 		utils.WriteJSON(w, 500, "Failed fetch messages", nil)
// 		return
// 	}

// 	defer cursor.Close(context.Background())

// 	var messages []models.Message

// 	if err := cursor.All(context.Background(), &messages); err != nil {
// 		utils.WriteJSON(w, 500, "Failed decode messages", nil)
// 		return
// 	}

// 	if messages == nil {
// 		messages = []models.Message{}
// 	}

// 	utils.WriteJSON(w, 200, "Success", messages)
// }

package handlers

import (
	"context"
	"net/http"
	"strings"

	"backend-chat-go/config"
	"backend-chat-go/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetMessages(c *gin.Context) {

	userId := c.GetString("userId")

	if userId == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
		})
		return
	}

	roomID := c.Query("roomId")

	if roomID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "roomId required",
		})
		return
	}

	// pastikan user bagian dari room
	parts := strings.Split(roomID, "_")

	if len(parts) != 2 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid roomId",
		})
		return
	}

	if parts[0] != userId && parts[1] != userId {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "Forbidden",
		})
		return
	}

	cursor, err := config.MessageCollection.Find(
		context.Background(),
		bson.M{"roomId": roomID},
		options.Find().SetSort(bson.M{"createdAt": 1}),
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed fetch messages",
		})
		return
	}

	defer cursor.Close(context.Background())

	var messages []models.Message

	if err := cursor.All(context.Background(), &messages); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed decode messages",
		})
		return
	}

	if messages == nil {
		messages = []models.Message{}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
		"data":    messages,
	})
}
