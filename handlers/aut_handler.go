package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"backend-chat-go/config"
	"backend-chat-go/models"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func Register(w http.ResponseWriter, r *http.Request) {

	var data map[string]string
	json.NewDecoder(r.Body).Decode(&data)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), bcrypt.DefaultCost)

	user := models.User{
		Username:  data["username"],
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
	}

	_, err := config.UserCollection.InsertOne(context.Background(), user)
	if err != nil {
		writeJSON(w, 400, "User already exists", nil)
		return
	}

	writeJSON(w, 200, "Register success", nil)
}

func Login(w http.ResponseWriter, r *http.Request) {

	var data map[string]string
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		writeJSON(w, 400, "Invalid request body", nil)
		return
	}

	if data["username"] == "" || data["password"] == "" {
		writeJSON(w, 400, "Username and password required", nil)
		return
	}

	var user models.User
	err := config.UserCollection.FindOne(
		context.Background(),
		bson.M{"username": data["username"]},
	).Decode(&user)

	if err != nil {
		writeJSON(w, 400, "User not found", nil)
		return
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(data["password"]),
	)

	if err != nil {
		writeJSON(w, 400, "Wrong password", nil)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(config.JwtSecret)
	if err != nil {
		writeJSON(w, 500, "Failed to generate token", nil)
		return
	}

	writeJSON(w, 200, "Login success", map[string]string{
		"token": tokenString,
	})
}

func writeJSON(w http.ResponseWriter, status int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(models.ApiResponse{
		StatusCode: status,
		Message:    message,
		Data:       data,
	})
}
