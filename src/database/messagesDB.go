package database

import (
	"bandplan/src/models"
	"log"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

func MessagesTableCreateMessage(bandID string, userID string, userName string, chatID string, body string) (models.Message, error) {
	log.Println("- MessagesTableCreateMessage")

	messageID := uuid.New().String()

	var message models.Message

	query := `
	INSERT INTO messages(
		message_id,
		band_id,
		user_id,
		chat_id,
		body
	)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, message_id, band_id, user_id, chat_id, body, created_at
	`
	err := DB.QueryRow(
		query,
		messageID,
		bandID,
		userID,
		chatID,
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
		log.Println("   Unable to save message to db: ", err)
		return models.Message{}, err
	}
	message.UserName = userName

	return message, nil
}

func MessagesTableGetAllMessages() ([]models.Message, error) {
	log.Println("- MessagesTableGetAllMessages")
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
		log.Println("   Unable to get all messages: ", err)
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
	log.Println("- MessagesTableGetLatestMessages")
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
	log.Println("- MessagesTableDeleteAll")
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
	log.Println("- MessagesTableGetAllMessagesByBandID")

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
		log.Println("   unable to get messages by band id: ", err)
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

func MessagesTableGetAllMessagesByChatID(chatID string) ([]models.Message, error) {
	log.Println("- MessagesTableGetAllMessagesByChatID")

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
		WHERE messages.chat_id = $1
		ORDER BY messages.created_at ASC
	`

	rows, err := DB.Query(query, chatID)
	if err != nil {
		log.Println("   unable to get messages by band id: ", err)
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
