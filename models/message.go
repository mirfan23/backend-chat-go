package models

import "time"

type Message struct {
	Type      string    `json:"type"`
	RoomID    string    `json:"roomId,omitempty"`
	Sender    string    `json:"sender,omitempty"`
	Text      string    `json:"text,omitempty"`
	IsTyping  bool      `json:"isTyping,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
}
