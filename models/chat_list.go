package models

import "time"

type ChatListResponse struct {
	RoomId          string    `json:"roomId"`
	Friend          string    `json:"friend"`
	LastMessage     string    `json:"lastMessage"`
	LastMessageTime time.Time `json:"lastMessageTime"`
	LastSender      string    `json:"lastSender"`
	IsOnline        bool      `json:"isOnline"`
	UnreadCount     int64     `json:"unreadCount"`
}
