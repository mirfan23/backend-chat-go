package models

import "time"

type Message struct {
	Type         string    `json:"type,omitempty"`
	RoomID       string    `bson:"roomId" json:"roomId"`
	Sender       string    `bson:"sender" json:"sender"`
	Receiver     string    `bson:"receiver" json:"receiver"`
	CipherText   string    `bson:"cipherText" json:"cipherText"`
	EncryptedKey string    `bson:"encryptedKey" json:"encryptedKey"`
	IV           string    `bson:"iv" json:"iv"`
	Preview      string    `bson:"preview" json:"preview"`
	IsRead       bool      `bson:"isRead" json:"isRead"`
	IsTyping     bool      `json:"isTyping"`
	CreatedAt    time.Time `bson:"createdAt" json:"createdAt"`
}
