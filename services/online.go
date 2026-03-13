package services

func IsUserOnline(userId string) bool {
	Mutex.RLock()
	defer Mutex.RUnlock()

	count, ok := OnlineUsers[userId]
	return ok && count > 0
}
