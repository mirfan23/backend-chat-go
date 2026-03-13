package handlers

import "sync"

var (
	onlineUsers = make(map[string]bool) // key = userId
	onlineMutex sync.RWMutex
)

func SetUserOnline(userId string) {
	onlineMutex.Lock()
	defer onlineMutex.Unlock()
	onlineUsers[userId] = true
}

func SetUserOffline(userId string) {
	onlineMutex.Lock()
	defer onlineMutex.Unlock()
	delete(onlineUsers, userId)
}

func IsUserOnline(userId string) bool {
	onlineMutex.RLock()
	defer onlineMutex.RUnlock()
	return onlineUsers[userId]
}
