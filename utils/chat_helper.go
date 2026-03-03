package utils

import "strings"

// extartFriend mengambil username lawan chat dari roomId
// format roomId : username1_username2

func ExtractFriend(roomID string, currentUser string) string {
	parts := strings.Split(roomID, "_")
	if len(parts) != 2 {
		return ""
	}
	if parts[0] == currentUser {
		return parts[1]
	}
	return parts[0]
}
