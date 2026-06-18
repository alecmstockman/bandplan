package database

import (
	"bandplan/src/models"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDB() *sql.DB {
	connStr := "postgres://alecstockman:yourpassword@localhost:5432/bandplan?sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func CreateMessagesTable(db *sql.DB) {
	log.Println("- CreateMessagesTable")
	query := `
	CREATE TABLE IF NOT EXISTS messages (
		id SERIAL PRIMARY KEY,
		message_id TEXT,
		band_id TEXT REFERENCES bands(band_id),
		user_id TEXT REFERENCES users(user_id),
		body TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)
	`
	_, err := db.Exec(query)
	if err != nil {
		fmt.Println("Unable to create or load messages table")
		log.Fatal(err)
	}
}

func MessagesTableCreateMessage(bandID string, userID string, userName string, body string) (models.Message, error) {
	log.Println("- MessagesTableCreateMessage")
	fmt.Println("++ BAND ID: ", bandID)

	messageID := uuid.New().String()

	var message models.Message

	query := `
	INSERT INTO messages(
		message_id,
		band_id,
		user_id,
		body
	)
	VALUES ($1, $2, $3, $4)
	RETURNING id, message_id, band_id, user_id, body, created_at
	`
	err := DB.QueryRow(
		query,
		messageID,
		bandID,
		userID,
		body,
	).Scan(
		&message.ID,
		&message.MessageID,
		&message.BandID,
		&message.UserID,
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
		messages.user_id,
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
		messages.user_id,
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
			&message.UserID,
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
	SELECT * 
	FROM messages
	WHERE messages.band_id = $1
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
