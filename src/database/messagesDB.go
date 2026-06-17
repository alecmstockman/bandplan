package database

import (
	"bandplan/src/models"
	"database/sql"
	"fmt"
	"log"
	"time"

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
	query := `
	CREATE TABLE IF NOT EXISTS messages (
		id SERIAL PRIMARY KEY,
		message_id TEXT,
		user_id TEXT REFERENCES users(user_id),
		band_id TEXT REFERENCES bands(band_id),
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

func MessagesTableInsertMessage(message models.Message) error {
	fmt.Println("\nINSERT")
	fmt.Println(message)
	query := `
	INSERT INTO messages(
		message_id,
		user_id,
		body
	)
	VALUES ($1, $2, $3)
	`
	_, err := DB.Exec(
		query,
		message.MessageID,
		message.UserID,
		message.Body,
	)
	return err
}

func MessagesTableGetAllMessages() ([]models.Message, error) {
	fmt.Println("MessagesTableGetAllMessages")
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

	return messages, nil
}

func MessagesTableGetLatestMessages() ([]models.Message, error) {
	fmt.Println("messagesTableGetAllMessages")
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
	query := `
	TRUNCATE messages RESTART IDENTITY
	`
	_, err := DB.Exec(query)
	if err != nil {
		return err
	}

	return nil
}
