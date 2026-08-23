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

func ChatsTableGetChatPreviewsByUserID(userID string) ([]models.ChatPreview, error) {
	log.Println("- ChatsTableGetChatPreviewsByUserID")

	query := `
		SELECT
			c.chat_id,
			m.user_id,
			c.name,
			c.is_primary,

			lm.message_id,
			lm.user_id AS latest_sender_id,
			u.display_name AS latest_sender_name,
			lm.body AS latest_message,
			lm.created_at AS latest_message_time,

			c.created_at,
			c.updated_at

		FROM chat_members m

		JOIN chats c
			ON c.chat_id = m.chat_id

		LEFT JOIN LATERAL (
			SELECT
				msg.message_id,
				msg.user_id,
				msg.body,
				msg.created_at
			FROM messages msg
			WHERE msg.chat_id = c.chat_id
			ORDER BY msg.created_at DESC
			LIMIT 1
		) lm ON true

		LEFT JOIN users u
			ON u.user_id = lm.user_id

		WHERE m.user_id = $1

		ORDER BY
			COALESCE(lm.created_at, c.updated_at) DESC
	`

	rows, err := DB.Query(query, userID)
	if err != nil {
		log.Println("   Unable to get user chat previews: ", err)
		return nil, err
	}
	defer rows.Close()

	var chats []models.ChatPreview

	for rows.Next() {
		var chat models.ChatPreview

		err := rows.Scan(
			&chat.ChatID,
			&chat.UserID,
			&chat.Name,
			&chat.IsPrimary,

			&chat.LatestMessageID,
			&chat.LatestSenderID,
			&chat.LatestSenderName,
			&chat.LatestMessage,
			&chat.LatestMessageTime,

			&chat.CreatedAt,
			&chat.UpdatedAt,
		)
		if err != nil {
			log.Println("   Unable to scan chat preview: ", err)
			return nil, err
		}

		chats = append(chats, chat)
	}

	if err := rows.Err(); err != nil {
		log.Println("   Error iterating user chat previews: ", err)
		return nil, err
	}

	return chats, nil
}

func ChatsTableCreateChat(chat models.Chat) error {
	log.Println("- ChatsTableCreateChat")

	return nil
}
