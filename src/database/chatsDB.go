package database

import (
	"bandplan/src/models"
	"fmt"

	"github.com/google/uuid"
)

func ChatsTableCreatePrimaryBandChat(bandID string, name string, slug string, userID string) (string, error) {

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
		return "", err
	}

	return chatID, nil
}

func ChatsTableGetPrimaryChatIDByBandID(bandID string) (string, error) {

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
		return "", err
	}

	return chatID, nil
}

func ChatsTableGetPrimaryChatPreviewByBandID(bandID string) (models.ChatPreview, error) {

	query := `
		SELECT
			c.chat_id,
			c.name,
			c.is_primary,
			COALESCE(C.image_id, ''),
			COALESCE(c.image_path, ''),

			COALESCE(lm.message_id, '') AS latest_message_id,
			COALESCE(lm.user_id, '') AS latest_sender_id,
			COALESCE(u.display_name, '') AS latest_sender_name,
			COALESCE(lm.body, '') AS latest_message,
			COALESCE(lm.created_at, c.created_at) AS latest_message_time,

			c.created_at,
			c.updated_at

		FROM chats c

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

		WHERE c.band_id = $1
			AND c.is_primary = TRUE

		LIMIT 1
	`

	var chat models.ChatPreview

	err := DB.QueryRow(query, bandID).Scan(
		&chat.ChatID,
		&chat.Name,
		&chat.IsPrimary,
		&chat.ImageID,
		&chat.ImagePath,

		&chat.LatestMessageID,
		&chat.LatestSenderID,
		&chat.LatestSenderName,
		&chat.LatestMessage,
		&chat.LatestMessageTime,

		&chat.CreatedAt,
		&chat.UpdatedAt,
	)

	if err != nil {
		return models.ChatPreview{}, err
	}

	return chat, nil
}

func ChatsTableGetChatByChatID(chatID string) (models.Chat, error) {

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
		return models.Chat{}, err
	}

	return chat, nil
}

func ChatMembersTableAddMember(chatID string, userID string) error {

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
		return err
	}

	return nil
}

func ChatMembersTableRemoveMember(chatID string, userID string) (bool, error) {

	result, err := DB.Exec(`
		DELETE FROM chat_members
		WHERE chat_id = $1 AND user_id = $2
	`, chatID, userID)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected == 1, nil
}

func ChatMembersTableGetMembersByChatID(chatID string) ([]models.User, error) {

	query := `
		SELECT
			u.user_id,
			u.display_name,
			COALESCE(u.profile_image_path, '')
		FROM chat_members cm
		JOIN users u
			ON u.user_id = cm.user_id
		WHERE cm.chat_id = $1
		ORDER BY LOWER(u.display_name), u.user_id
	`

	rows, err := DB.Query(query, chatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	members := make([]models.User, 0)
	for rows.Next() {
		var member models.User
		if err := rows.Scan(
			&member.UserID,
			&member.DisplayName,
			&member.ProfileImagePath,
		); err != nil {
			return nil, err
		}
		members = append(members, member)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return members, nil
}

func ChatMembersTableGetChatIDsByUserID(userID string) (map[string]bool, error) {

	query := `
		SELECT chat_id
		FROM chat_members
		WHERE user_id = $1
	`

	rows, err := DB.Query(query, userID)
	if err != nil {
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
			return nil, err
		}
		chatIDs[chatID] = true
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return chatIDs, nil
}

func ChatsTableGetChatPreviewsByUserID(userID string) ([]models.ChatPreview, error) {

	query := `
		SELECT
			c.chat_id,
			m.user_id,
			c.name,
			c.is_primary,
			COALESCE(c.image_id, ''),
			COALESCE(c.image_path, ''),

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
			&chat.ImageID,
			&chat.ImagePath,

			&chat.LatestMessageID,
			&chat.LatestSenderID,
			&chat.LatestSenderName,
			&chat.LatestMessage,
			&chat.LatestMessageTime,

			&chat.CreatedAt,
			&chat.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		chats = append(chats, chat)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return chats, nil
}

func ChatsTableCreateChat(chat models.Chat, memberIDs []string) (string, error) {

	chatID := uuid.New().String()

	tx, err := DB.Begin()
	if err != nil {
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
			return "", err
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return "", err
		}
		if rowsAffected != 1 {
			return "", fmt.Errorf("user %q is not a member of band %q", memberID, chat.BandID)
		}
	}

	if err := tx.Commit(); err != nil {
		return "", err
	}

	return chatID, nil
}

func ChatsTableDeleteChatByChatID(chatID string) (bool, error) {

	query := `
		DELETE FROM chats
		WHERE chat_id = $1
	`

	result, err := DB.Exec(query, chatID)

	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected == 1, nil
}

func ChatsTableUpdateChat(chat models.Chat) (bool, error) {

	query := `
		UPDATE chats
		SET
			image_id = $1, 
			image_path = $2
		WHERE
			chat_id = $3
	`

	result, err := DB.Exec(
		query,
		chat.ImageID,
		chat.ImagePath,
		chat.ChatID,
	)

	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected == 1, nil
}

func ChatMembersTableUserIsMember(chatID string, userID string) (bool, error) {

	query := `
		SELECT EXISTS (
			SELECT 1
			FROM chat_members
			WHERE chat_id = $1
			  AND user_id = $2
		)
	`

	var exists bool

	err := DB.QueryRow(
		query,
		chatID,
		userID,
	).Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil
}

func ChatsTableGetThreeRecentChats(userID string) ([]models.Message, error) {

	query := `
		SELECT *
		FROM (
			SELECT DISTINCT ON (m.chat_id)
				m.id,
				m.message_id,
				m.band_id,
				m.user_id,
				u.profile_image_path,
				u.display_name AS user_name,
				m.chat_id,
				c.name AS chat_name,
				m.body,
				m.is_pinned,
				m.pinned_at,
				m.pinned_by,
				m.created_at,
				m.edited_at
			FROM messages m
			INNER JOIN chat_members cm
				ON cm.chat_id = m.chat_id
			INNER JOIN users u
				ON u.user_id = m.user_id
			INNER JOIN chats c
				ON c.chat_id = m.chat_id
			WHERE cm.user_id = $1
			ORDER BY m.chat_id, m.created_at DESC
		) AS recent_chats
		ORDER BY created_at DESC
		LIMIT 3
	`

	rows, err := DB.Query(query, userID)
	if err != nil {
		return []models.Message{}, err
	}

	defer rows.Close()

	var messages []models.Message

	for rows.Next() {
		var message models.Message

		err := rows.Scan(
			&message.ID,
			&message.MessageID,
			&message.BandID,
			&message.UserID,
			&message.ProfileImagePath,
			&message.UserName,
			&message.ChatID,
			&message.ChatName,
			&message.Body,
			&message.IsPinned,
			&message.PinnedAt,
			&message.PinnedBy,
			&message.CreatedAt,
			&message.EditedAt,
		)
		if err != nil {
			return []models.Message{}, err
		}
		messages = append(messages, message)
	}

	return messages, nil
}
