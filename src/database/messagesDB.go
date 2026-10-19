package database

import (
	"bandplan/src/models"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

func MessagesTableCreateMessage(bandID, userID, userName, chatID, messageID, body, timeZone string) (models.Message, error) {

	var message models.Message

	query := `
	INSERT INTO messages(
		band_id,
		user_id,
		chat_id,
		message_id,
		body
	)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, message_id, band_id, user_id, chat_id, body, created_at
	`
	err := DB.QueryRow(
		query,
		bandID,
		userID,
		chatID,
		messageID,
		body,
	).Scan(
		&message.ID,
		&message.MessageID,
		&message.BandID,
		&message.UserID,
		&message.ChatID,
		&message.Body,
		&message.CreatedAt,
	)
	if err != nil {
		return models.Message{}, err
	}
	message.UserName = userName

	return message, nil
}

func MessagesTableGetAllMessages() ([]models.Message, error) {

	query := `
	SELECT 
		messages.id,
		messages.message_id,
		messages.band_id,
		messages.user_id,
		COALESCE(users.profile_image_path, ''),
		users.name,
		messages.body,
		messages.created_at
	FROM messages
	JOIN users ON messages.user_id = users.user_id
	ORDER BY messages.created_at ASC
	`
	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message
	var message models.Message

	for rows.Next() {
		err := rows.Scan(
			&message.ID,
			&message.MessageID,
			&message.BandID,
			&message.UserID,
			&message.ProfileImagePath,
			&message.UserName,
			&message.Body,
			&message.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	return messages, nil
}

func MessagesTableGetLatestMessages() ([]models.Message, error) {

	t := time.Now().Add(-2 * time.Second)

	query := `
	SELECT 
		messages.id,
		messages.message_id,
		messages.band_id,
		messages.user_id,
		COALESCE(users.profile_image_path, ''),
		users.name,
		messages.body,
		messages.created_at
	FROM messages
	JOIN users ON messages.user_id = users.user_id
	WHERE created_at > $1
	`
	rows, err := DB.Query(query, t)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message
	var message models.Message

	for rows.Next() {
		err := rows.Scan(
			&message.ID,
			&message.MessageID,
			&message.BandID,
			&message.UserID,
			&message.ProfileImagePath,
			&message.UserName,
			&message.Body,
			&message.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		messages = append(messages, message)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

func MessagesTableDeleteAll() error {

	query := `
	TRUNCATE messages RESTART IDENTITY
	`
	_, err := DB.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func MessagesTableGetAllMessagesByBandID(bandID string) ([]models.Message, error) {

	query := `
	SELECT
		messages.id,
		messages.message_id,
		messages.band_id,
		messages.user_id,
		COALESCE(users.profile_image_path, ''),
		users.name,
		messages.body,
		messages.created_at
	FROM messages
	JOIN users ON messages.user_id = users.user_id
	WHERE messages.band_id = $1
	ORDER BY messages.created_at ASC
	`

	rows, err := DB.Query(query, bandID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message
	var message models.Message

	for rows.Next() {
		err := rows.Scan(
			&message.ID,
			&message.MessageID,
			&message.BandID,
			&message.UserID,
			&message.ProfileImagePath,
			&message.UserName,
			&message.Body,
			&message.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	return messages, nil
}

func MessagesTableGetAllMessagesByChatID(chatID, bandID, userID string) ([]models.Message, error) {

	query := `
		SELECT
			m.id,
			m.message_id,
			m.band_id,
			m.user_id,
			COALESCE(users.profile_image_path, ''),
			users.name,
			m.body,
			m.created_at,
			COALESCE(r.reactions, '[]'::json) AS reactions
		FROM messages m
		JOIN users ON m.user_id = users.user_id
		JOIN chats c
			ON c.chat_id = m.chat_id
			AND c.band_id = m.band_id
		JOIN chat_members cm
			ON cm.chat_id = c.chat_id
			AND cm.user_id = $3

		LEFT JOIN LATERAL (
			SELECT json_agg(
				json_build_object(
					'reaction_id', mr.reaction_id,
					'user_id', mr.user_id,
					'reaction', mr.reaction,
					'created_at', mr.created_at
				)
				ORDER BY mr.created_at
			) AS reactions
			FROM message_reactions mr
			WHERE mr.message_id = m.message_id
		) r ON true

		WHERE m.chat_id = $1
			AND m.band_id = $2
			AND c.band_id = $2
		ORDER BY m.created_at ASC
		LIMIT 200
	`

	rows, err := DB.Query(query, chatID, bandID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message
	var message models.Message
	var reactionsJSON []byte

	for rows.Next() {
		err := rows.Scan(
			&message.ID,
			&message.MessageID,
			&message.BandID,
			&message.UserID,
			&message.ProfileImagePath,
			&message.UserName,
			&message.Body,
			&message.CreatedAt,
			&reactionsJSON,
		)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(reactionsJSON, &message.Reactions)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	return messages, nil
}

func MessageReactionsTableAddReaction(messageID, userID, bandID, reaction string) (string, error) {

	reactionID := uuid.New().String()

	query := `
		WITH authorized_message AS (
			SELECT m.message_id
			FROM messages m
			JOIN chats c
				ON c.chat_id = m.chat_id
				AND c.band_id = m.band_id
			JOIN chat_members cm
				ON cm.chat_id = c.chat_id
				AND cm.user_id = $2
			WHERE m.message_id = $1
				AND m.band_id = $3
				AND c.band_id = $3
		),
		deleted AS (
			DELETE FROM message_reactions AS mr
			WHERE mr.message_id = $1
				AND mr.user_id = $2
				AND mr.reaction = $4
				AND EXISTS (SELECT 1 FROM authorized_message)
			RETURNING 1
		),
		updated AS (
			UPDATE message_reactions AS mr
			SET reaction = $4
			WHERE mr.message_id = $1
				AND mr.user_id = $2
				AND mr.reaction <> $4
				AND EXISTS (SELECT 1 FROM authorized_message)
			RETURNING 1
		),
		inserted AS (
			INSERT INTO message_reactions (
				reaction_id,
				message_id,
				user_id,
				reaction
			)
			SELECT $5, am.message_id, $2, $4
			FROM authorized_message am
			WHERE NOT EXISTS (
				SELECT 1 FROM deleted
			)
			AND NOT EXISTS (
				SELECT 1 FROM updated
			)
			RETURNING 1
		)
		SELECT 'removed' AS action FROM deleted
		UNION ALL
		SELECT 'updated' AS action FROM updated
		UNION ALL
		SELECT 'added' AS action FROM inserted
	`

	var action string
	err := DB.QueryRow(
		query,
		messageID,
		userID,
		bandID,
		reaction,
		reactionID,
	).Scan(&action)

	if err != nil {
		return "", err
	}

	return action, nil
}

func MessageReactionsTableGetReactionsByMessageID(messageID, userID, bandID string) ([]models.MessageReaction, error) {

	query := `
		SELECT 
			m.id, 
			m.reaction_id,
			m.message_id,
			u.display_name,
			m.user_id,
			m.reaction,
			m.created_at
		FROM
			message_reactions m
		JOIN messages msg
			ON msg.message_id = m.message_id
		JOIN chats c
			ON c.chat_id = msg.chat_id
			AND c.band_id = msg.band_id
		JOIN chat_members cm
			ON cm.chat_id = c.chat_id
			AND cm.user_id = $2
		LEFT JOIN users u
			ON m.user_id = u.user_id
		WHERE 
			m.message_id = $1
			AND msg.band_id = $3
			AND c.band_id = $3
	`

	rows, err := DB.Query(query, messageID, userID, bandID)
	if err != nil {
		return []models.MessageReaction{}, err
	}

	defer rows.Close()

	var reactions []models.MessageReaction

	for rows.Next() {
		var reaction models.MessageReaction

		err := rows.Scan(
			&reaction.ID,
			&reaction.ReactionID,
			&reaction.MessageID,
			&reaction.UserName,
			&reaction.UserID,
			&reaction.Reaction,
			&reaction.CreatedAt,
		)
		if err != nil {
			return []models.MessageReaction{}, err
		}

		reactions = append(reactions, reaction)
	}

	if err := rows.Err(); err != nil {
		return []models.MessageReaction{}, err
	}

	return reactions, nil
}

func MessagesTablePinMessage(messageID, chatID, bandID, userID string) (bool, error) {

	query := `
		UPDATE messages AS m
		SET
			is_pinned = TRUE,
			pinned_at = CURRENT_TIMESTAMP,
			pinned_by = $4
		FROM chats c
		JOIN chat_members cm
			ON cm.chat_id = c.chat_id
			AND cm.user_id = $4
		WHERE m.message_id = $1
			AND m.chat_id = $2
			AND m.band_id = $3
			AND c.chat_id = m.chat_id
			AND c.band_id = $3
			AND m.band_id = c.band_id
	`

	result, err := DB.Exec(query, messageID, chatID, bandID, userID)
	if err != nil {
		return false, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	return rowsAffected == 1, nil
}

func MessagesTableUnPinMessage(messageID, chatID, bandID, userID string) (bool, error) {

	query := `
		UPDATE messages AS m
		SET
			is_pinned = FALSE,
			pinned_at = NULL,
			pinned_by = NULL
		FROM chats c
		JOIN chat_members cm
			ON cm.chat_id = c.chat_id
			AND cm.user_id = $4
		WHERE m.message_id = $1
			AND m.chat_id = $2
			AND m.band_id = $3
			AND c.chat_id = m.chat_id
			AND c.band_id = $3
			AND m.band_id = c.band_id
	`
	result, err := DB.Exec(query, messageID, chatID, bandID, userID)
	if err != nil {
		return false, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	return rowsAffected == 1, nil
}

func MessagesTableGetPinnedMessagesByChatID(chatID, bandID, userID string) ([]models.Message, error) {

	query := `
        SELECT
            m.id,
            m.message_id,
            m.band_id,
            m.user_id,
            COALESCE(u.profile_image_path, ''),
            u.display_name,
            m.body,
            m.is_pinned,
            m.pinned_at,
            COALESCE(m.pinned_by, ''),
            m.created_at,
            m.edited_at
		FROM messages m
		JOIN chats c
			ON c.chat_id = m.chat_id
			AND c.band_id = m.band_id
		JOIN chat_members cm
			ON cm.chat_id = c.chat_id
			AND cm.user_id = $3
		JOIN users u ON m.user_id = u.user_id
		WHERE m.chat_id = $1
			AND m.band_id = $2
			AND c.band_id = $2
			AND m.is_pinned = TRUE
		ORDER BY m.pinned_at DESC
	`

	rows, err := DB.Query(query, chatID, bandID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pinnedMessages []models.Message

	for rows.Next() {
		var message models.Message

		err := rows.Scan(
			&message.ID,
			&message.MessageID,
			&message.BandID,
			&message.UserID,
			&message.ProfileImagePath,
			&message.UserName,
			&message.Body,
			&message.IsPinned,
			&message.PinnedAt,
			&message.PinnedBy,
			&message.CreatedAt,
			&message.EditedAt,
		)
		if err != nil {
			return nil, err
		}

		pinnedMessages = append(pinnedMessages, message)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return pinnedMessages, nil
}
