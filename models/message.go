package models

import "time"

type Message struct {
	Type      string    `json:"type,omitempty"`
	RoomID    string    `bson:"roomId" json:"roomId"`
	Sender    string    `bson:"sender" json:"sender"`
	Receiver  string    `bson:"receiver" json:"receiver"`
	Text      string    `bson:"text" json:"text"`
	IsRead    bool      `bson:"isRead" json:"isRead"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
}
