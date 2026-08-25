package database

import (
	"bandplan/src/models"
	"fmt"
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

func ChatsTableGetPrimaryChatByBandID(bandID string) (string, error) {
	log.Println("- ChatsTableGetChatIDByBandID")

	query := `
		SELECT chat_id
		FROM chats
		WHERE chats.band_id = $1
			AND chats.is_primary = TRUE;
	`

	var chatID string

	err := DB.QueryRow(query, bandID).Scan(
		&chatID,
	)
	if err != nil {
		log.Println("   Unable to get primary chat id by band id: ", err)
		return "", err
	}

	return chatID, nil
}

func ChatsTableGetChatByChatID(chatID string) (models.Chat, error) {
	log.Println("- ChatsTableGetChatByChatID")

	query := `
		SELECT
			id,
			chat_id,
			band_id,
			name,
			slug,
			is_primary,
			COALESCE(image_id, ''),
			COALESCE(image_path, ''),
			created_at,
			created_by,
			updated_at,
			updated_by
		FROM chats
		WHERE chat_id = $1
	`

	var chat models.Chat

	err := DB.QueryRow(query, chatID).Scan(
		&chat.ID,
		&chat.ChatID,
		&chat.BandID,
		&chat.Name,
		&chat.Slug,
		&chat.IsPrimary,
		&chat.ImageID,
		&chat.ImagePath,
		&chat.CreatedAt,
		&chat.CreatedBy,
		&chat.UpdatedAt,
		&chat.UpdatedBy,
	)
	if err != nil {
		log.Println("   Unable to get chat from database: ", err)
		return models.Chat{}, err
	}

	return chat, nil
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

func ChatMembersTableRemoveMember(chatID string, userID string) (bool, error) {
	log.Println("- ChatMembersTableRemoveMember")

	result, err := DB.Exec(`
		DELETE FROM chat_members
		WHERE chat_id = $1 AND user_id = $2
	`, chatID, userID)
	if err != nil {
		log.Println("   Unable to remove user from chat members table: ", err)
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("   Unable to confirm chat member removal: ", err)
		return false, err
	}

	return rowsAffected == 1, nil
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

	fmt.Println("UserID: ", userID)

	query := `
		SELECT
			c.chat_id,
			m.user_id,
			c.name,
			c.is_primary,

			COALESCE(lm.message_id, '') AS latest_message_id,
			COALESCE(lm.user_id, '') AS latest_sender_id,
			COALESCE(u.display_name, '') AS latest_sender_name,
			COALESCE(lm.body, '') AS latest_message,
			COALESCE(lm.created_at, c.created_at) AS latest_message_time,

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

func ChatsTableCreateChat(chat models.Chat, memberIDs []string) (string, error) {
	log.Println("- ChatsTableCreateChat")

	chatID := uuid.New().String()

	tx, err := DB.Begin()
	if err != nil {
		log.Println("   Unable to start chat creation transaction: ", err)
		return "", err
	}
	defer tx.Rollback()

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

	_, err = tx.Exec(
		query,
		chatID,
		chat.BandID,
		chat.Name,
		chat.Slug,
		chat.IsPrimary,
		chat.CreatedBy,
		chat.UpdatedBy,
	)
	if err != nil {
		log.Println("   Unable to create chat: ", err)
		return "", err
	}

	uniqueMemberIDs := make(map[string]struct{}, len(memberIDs)+1)
	uniqueMemberIDs[chat.CreatedBy] = struct{}{}
	for _, memberID := range memberIDs {
		uniqueMemberIDs[memberID] = struct{}{}
	}

	for memberID := range uniqueMemberIDs {
		result, err := tx.Exec(`
			INSERT INTO chat_members (chat_id, user_id)
			SELECT $1, user_id
			FROM band_members
			WHERE band_id = $2 AND user_id = $3
		`, chatID, chat.BandID, memberID)
		if err != nil {
			log.Println("   Unable to add member to chat: ", err)
			return "", err
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			log.Println("   Unable to confirm chat member insertion: ", err)
			return "", err
		}
		if rowsAffected != 1 {
			return "", fmt.Errorf("user %q is not a member of band %q", memberID, chat.BandID)
		}
	}

	if err := tx.Commit(); err != nil {
		log.Println("   Unable to commit chat creation: ", err)
		return "", err
	}

	return chatID, nil
}

// func ChatsTableGetChatCreatorByChatID()

func ChatsTableDeleteChatByChatID(chatID string) (bool, error) {
	log.Println("- ChatsTableDeleteChatByChatID")

	query := `
		DELETE FROM chats
		WHERE chat_id = $1
	`

	result, err := DB.Exec(query, chatID)

	if err != nil {
		log.Printf("   Unable to delete chat, id# %v due to: %v\n", chatID, err)
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("   Unable to confirm chat deletion, id# %v due to: %v\n", chatID, err)
		return false, err
	}

	return rowsAffected == 1, nil
}
