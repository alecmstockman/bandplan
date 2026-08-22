package database

import (
	"bandplan/src/models"
	"log"

	"github.com/google/uuid"
)

func ChatsTableCreatePrimaryBandChat(bandID string, name string, slug string, userID string) (string, error) {
	log.Println("- ChatsTableCreatePrimaryBandChat")

	chatID := uuid.New().String()

	query := `
	INSERT INTO chats (
		chat_id,
		band_id,
		name,
		slug,
		is_primary,
		created_by,
		updated_by
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7
	)
	`

	_, err := DB.Exec(
		query,
		chatID,
		bandID,
		name,
		slug,
		true,
		userID,
		userID,
	)
	if err != nil {
		log.Println("   Unable to create new default band chat: ", err)
		return "", err
	}

	return chatID, nil
}

func ChatMembersTableAddMember(chatID string, userID string) error {
	log.Println("- ChatMembersTableAddMember")

	query := `
	INSERT INTO chat_members(
		chat_id,
		user_id
	) VALUES (
		$1, $2
	)
	`

	_, err := DB.Exec(query, chatID, userID)
	if err != nil {
		log.Printf("   Unable to add user to chat members table: %v\n", err)
		return err
	}

	return nil
}

func ChatMembersTableGetChatIDsByUserID(userID string) (map[string]bool, error) {
	log.Println("- ChatMembersTableGetChatIDsByUserID")

	query := `
		SELECT chat_id
		FROM chat_members
		WHERE user_id = $1
	`

	rows, err := DB.Query(query, userID)
	if err != nil {
		log.Println("   Unable to get user chatIDs from chat_members table: ", err)
		return nil, err
	}

	defer rows.Close()

	chatIDs := make(map[string]bool)

	for rows.Next() {
		var chatID string

		err = rows.Scan(
			&chatID,
		)
		if err != nil {
			log.Println("   Unable to chats by userID: ", err)
			return nil, err
		}
		chatIDs[chatID] = true
	}

	if err = rows.Err(); err != nil {
		log.Println("   Unable to iterate chatIDs: ", err)
		return nil, err
	}

	return chatIDs, nil
}

func ChatsTableGetChatsByUserID(userID string) ([]models.Chat, error) {
	log.Println("- ChatsTableGetChatsByUserID")

	query := `
		SELECT
			c.id,
			c.chat_id,
			c.band_id,
			c.name,
			c.slug,
			c.is_primary,
			c.created_at,
			c.created_by,
			c.updated_at,
			c.updated_by
		FROM chat_members m
		JOIN chats c
			ON c.chat_id = m.chat_id
		WHERE m.user_id = $1
	`

	rows, err := DB.Query(query, userID)
	if err != nil {
		log.Println("   Unable to get all user chats from chats database: ", err)
		return nil, err
	}

	defer rows.Close()

	var chats []models.Chat

	for rows.Next() {
		var chat models.Chat

		err = rows.Scan(
			&chat.ID,
			&chat.ChatID,
			&chat.BandID,
			&chat.Name,
			&chat.Slug,
			&chat.IsPrimary,
			&chat.CreatedAt,
			&chat.CreatedBy,
			&chat.UpdatedAt,
			&chat.UpdatedBy,
		)
		if err != nil {
			log.Println("   Unable to get user chats from database: ", err)
			return nil, err
		}
		chat.UserID = userID

		chats = append(chats, chat)
	}

	if err := rows.Err(); err != nil {
		log.Println("   Error iterating user chats: ", err)
		return nil, err
	}

	return chats, nil
}
