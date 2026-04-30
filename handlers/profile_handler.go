package handlers

import (
	"context"
	"net/http"

	"backend-chat-go/config"
	"backend-chat-go/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UpdateProfileRequest struct {
	Username string `json:"username"`
	Photo    string `json:"photo"`
}

func GetProfile(c *gin.Context) {

	userId := c.GetString("userId")

	objId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid userId",
		})
		return
	}

	var user models.User

	err = config.UserCollection.FindOne(
		context.Background(),
		bson.M{"_id": objId},
	).Decode(&user)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    user,
	})
}

func UpdateProfile(c *gin.Context) {

	userId := c.GetString("userId")

	objId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid userId",
		})
		return
	}

	var req UpdateProfileRequest

	// otomatis decode json body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
		})
		return
	}

	updateData := bson.M{}

	if req.Username != "" {
		updateData["username"] = req.Username
	}

	if req.Photo != "" {
		updateData["photo"] = req.Photo
	}

	_, err = config.UserCollection.UpdateOne(
		context.Background(),
		bson.M{"_id": objId},
		bson.M{
			"$set": updateData,
		},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed update profile",
		})
		return
	}

	var updatedUser models.User

	err = config.UserCollection.FindOne(
		context.Background(),
		bson.M{"_id": objId},
	).Decode(&updatedUser)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed get updated user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    updatedUser,
	})
}

// func GetProfile(w http.ResponseWriter, r *http.Request) {

// 	userId, ok := r.Context().Value("userId").(string)
// 	if !ok {
// 		utils.WriteJSON(w, 401, "Unauthorized", nil)
// 		return
// 	}

// 	objId, err := primitive.ObjectIDFromHex(userId)
// 	if err != nil {
// 		utils.WriteJSON(w, 400, "Invalid userId", nil)
// 		return
// 	}

// 	var user models.User

// 	err = config.UserCollection.FindOne(
// 		context.Background(),
// 		bson.M{"_id": objId},
// 	).Decode(&user)

// 	if err != nil {
// 		utils.WriteJSON(w, 404, "User not found", nil)
// 		return
// 	}

// 	utils.WriteJSON(w, 200, "success", user)
// }

// func UpdateProfile(w http.ResponseWriter, r *http.Request) {
// 	userId, ok := r.Context().Value("userId").(string)
// 	if !ok {
// 		utils.WriteJSON(w, 401, "Unauthorized", nil)
// 		return
// 	}

// 	objId, err := primitive.ObjectIDFromHex(userId)
// 	if err != nil {
// 		utils.WriteJSON(w, 400, "Invalid userId", nil)
// 		return
// 	}

// 	var req UpdateProfileRequest

// 	err = json.NewDecoder(r.Body).Decode(&req)
// 	if err != nil {
// 		utils.WriteJSON(w, 400, "Invalid request body", nil)
// 		return
// 	}

// 	updateData := bson.M{}

// 	if req.Username != "" {
// 		updateData["username"] = req.Username
// 	}

// 	if req.Photo != "" {
// 		updateData["photo"] = req.Photo
// 	}

// 	_, err = config.UserCollection.UpdateOne(
// 		context.Background(),
// 		bson.M{"_id": objId},
// 		bson.M{
// 			"$set": updateData,
// 		},
// 	)

// 	if err != nil {
// 		utils.WriteJSON(w, 400, "Failed update profile", nil)
// 		return
// 	}

// 	var updatedUser models.User

// 	err = config.UserCollection.FindOne(
// 		context.Background(),
// 		bson.M{"_id": objId},
// 	).Decode(&updatedUser)

// 	if err != nil {
// 		utils.WriteJSON(w, 400, "Failed get updated user", nil)
// 		return
// 	}

// 	utils.WriteJSON(w, 200, "success", nil)
// }
