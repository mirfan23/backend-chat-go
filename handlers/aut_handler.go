// package handlers

// import (
// 	"context"
// 	"encoding/json"
// 	"net/http"
// 	"time"

// 	"backend-chat-go/config"
// 	"backend-chat-go/models"
// 	"backend-chat-go/utils"

// 	"github.com/golang-jwt/jwt/v5"
// 	"go.mongodb.org/mongo-driver/bson"
// 	"golang.org/x/crypto/bcrypt"
// )

// func Register(w http.ResponseWriter, r *http.Request) {

// 	if r.Method != http.MethodPost {
// 		utils.WriteJSON(w, 405, "Method Not Allowed", nil)
// 		return
// 	}

// 	var data map[string]string
// 	json.NewDecoder(r.Body).Decode(&data)

// 	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), bcrypt.DefaultCost)

// 	user := models.User{
// 		Username:  data["username"],
// 		Password:  string(hashedPassword),
// 		PublicKey: data["publicKey"],
// 		CreatedAt: time.Now(),
// 	}

// 	_, err := config.UserCollection.InsertOne(context.Background(), user)
// 	if err != nil {
// 		utils.WriteJSON(w, 400, "User already exists", nil)
// 		return
// 	}

// 	utils.WriteJSON(w, 200, "Register success", nil)
// }

// func Login(w http.ResponseWriter, r *http.Request) {

// 	if r.Method != http.MethodPost {
// 		utils.WriteJSON(w, 405, "Method Not Allowed", nil)
// 		return
// 	}

// 	var data map[string]string
// 	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
// 		utils.WriteJSON(w, 400, "Invalid request body", nil)
// 		return
// 	}

// 	if data["username"] == "" || data["password"] == "" {
// 		utils.WriteJSON(w, 400, "Username and password required", nil)
// 		return
// 	}

// 	var user models.User
// 	err := config.UserCollection.FindOne(
// 		context.Background(),
// 		bson.M{"username": data["username"]},
// 	).Decode(&user)

// 	if err != nil {
// 		utils.WriteJSON(w, 400, "User not found", nil)
// 		return
// 	}

// 	err = bcrypt.CompareHashAndPassword(
// 		[]byte(user.Password),
// 		[]byte(data["password"]),
// 	)

// 	if err != nil {
// 		utils.WriteJSON(w, 400, "Wrong password", nil)
// 		return
// 	}

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
// 		"userId":   user.ID.Hex(),
// 		"username": user.Username,
// 		"exp":      time.Now().Add(24 * time.Hour).Unix(),
// 	})

// 	tokenString, err := token.SignedString(config.JwtSecret)
// 	if err != nil {
// 		utils.WriteJSON(w, 500, "Failed to generate token", nil)
// 		return
// 	}

// 	utils.WriteJSON(w, 200, "Login success", map[string]string{
// 		"token":    tokenString,
// 		"userId":   user.ID.Hex(),
// 		"username": user.Username,
// 	})
// }

package handlers

import (
	"context"
	"net/http"
	"time"

	"backend-chat-go/config"
	"backend-chat-go/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	PublicKey string `json:"publicKey"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Register(c *gin.Context) {

	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
		})
		return
	}

	if req.Username == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Username and password required",
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed hash password",
		})
		return
	}

	user := models.User{
		Username:  req.Username,
		Password:  string(hashedPassword),
		PublicKey: req.PublicKey,
		CreatedAt: time.Now(),
	}

	_, err = config.UserCollection.InsertOne(
		context.Background(),
		user,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "User already exists",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Register success",
	})
}

func Login(c *gin.Context) {

	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
		})
		return
	}

	if req.Username == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Username and password required",
		})
		return
	}

	var user models.User

	err := config.UserCollection.FindOne(
		context.Background(),
		bson.M{
			"username": req.Username,
		},
	).Decode(&user)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "User not found",
		})
		return
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(req.Password),
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Wrong password",
		})
		return
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"userId":   user.ID.Hex(),
			"username": user.Username,
			"exp":      time.Now().Add(24 * time.Hour).Unix(),
		},
	)

	tokenString, err := token.SignedString(config.JwtSecret)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed generate token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login success",
		"data": gin.H{
			"token":    tokenString,
			"userId":   user.ID.Hex(),
			"username": user.Username,
		},
	})
}
