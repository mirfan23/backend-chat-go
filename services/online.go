package services

func IsUserOnline(username string) bool {
	Mutex.RLock()
	defer Mutex.RUnlock()

	count, ok := OnlineUsers[username]
	return ok && count > 0
}
