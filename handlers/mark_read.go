// package handlers

// import (
// 	"context"
// 	"net/http"

// 	"backend-chat-go/config"
// 	"backend-chat-go/utils"

// 	"go.mongodb.org/mongo-driver/bson"
// )

// func MarkMessagesRead(w http.ResponseWriter, r *http.Request) {

// 	if r.Method != http.MethodGet {
// 		utils.WriteJSON(w, 405, "Method Not Allowed", nil)
// 		return
// 	}

// 	userId := r.Context().Value("userId")
// 	roomId := r.URL.Query().Get("roomId")

// 	ctx := context.Background()

// 	filter := bson.M{
// 		"roomId":   roomId,
// 		"receiver": userId,
// 		"isRead":   false,
// 	}

// 	update := bson.M{
// 		"$set": bson.M{
// 			"isRead": true,
// 		},
// 	}

// 	_, err := config.MessageCollection.UpdateMany(ctx, filter, update)

// 	if err != nil {
// 		utils.WriteJSON(w, 500, "Failed update", nil)
// 		return
// 	}

// 	utils.WriteJSON(w, 200, "Updated", nil)
// }

package handlers

import (
	"context"
	"net/http"

	"backend-chat-go/config"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func MarkMessagesRead(c *gin.Context) {

	userId := c.GetString("userId")

	if userId == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
		})
		return
	}

	roomId := c.Query("roomId")

	if roomId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "roomId required",
		})
		return
	}

	ctx := context.Background()

	filter := bson.M{
		"roomId":   roomId,
		"receiver": userId,
		"isRead":   false,
	}

	update := bson.M{
		"$set": bson.M{
			"isRead": true,
		},
	}

	_, err := config.MessageCollection.UpdateMany(
		ctx,
		filter,
		update,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed update",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Updated",
	})
}
