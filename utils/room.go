package utils

import "strings"

func GetFriendFromRoom(roomId string, username string) string {

	roomId = strings.ToLower(roomId)
	username = strings.ToLower(username)

	parts := strings.Split(roomId, "_")

	if len(parts) != 2 {
		return ""
	}

	user1 := parts[0]
	user2 := parts[1]

	if user1 == username {
		return user2
	}

	if user2 == username {
		return user1
	}

	return ""
}
