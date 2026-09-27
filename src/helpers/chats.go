package helpers

import (
	"bandplan/src/models"
	"slices"
)

func GetChatNonMembers(chatMembers, bandMembers []models.User) []models.User {
	var chatMemberIDs []string

	for _, chatMember := range chatMembers {
		chatMemberIDs = append(chatMemberIDs, chatMember.UserID)
	}

	var chatNonMembers []models.User

	for _, bandMember := range bandMembers {
		if !slices.Contains(chatMemberIDs, bandMember.UserID) {
			chatNonMembers = append(chatNonMembers, bandMember)
		}
	}

	return chatNonMembers
}
