package models

import "time"

type ChatListResponse struct {
	RoomId          string    `json:"roomId"`
	Friend          string    `json:"friend"`
	LastMessage     string    `json:"lastMessage"`
	LastMessageTime time.Time `json:"lastMessageTime"`
	LastSender      string    `json:"lastSender"`
	IsOnline        bool      `json:"isOnline"`
	IsRead          bool      `json:"isRead"`
	UnreadCount     int64     `json:"unreadCount"`
	FriendName      string    `json:"friendName"`
}
