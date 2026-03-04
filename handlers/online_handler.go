package handlers

import "sync"

var (
	onlineUsers = make(map[string]bool)
	onlineMutex sync.RWMutex
)

func SetUserOnline(username string) {
	onlineMutex.Lock()
	defer onlineMutex.Unlock()
	onlineUsers[username] = true
}

func SetUserOffline(username string) {
	onlineMutex.Lock()
	defer onlineMutex.Unlock()
	delete(onlineUsers, username)
}

func IsUserOnline(username string) bool {
	onlineMutex.RLock()
	defer onlineMutex.RUnlock()
	return onlineUsers[username]
}
