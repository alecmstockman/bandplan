package database

import (
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
	fmt.Println("CreateMessagesTable")
	query := `
	CREATE TABLE IF NOT EXISTS messages (
		id SERIAL PRIMARY KEY,
		body TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}

func MessagesTableInsertMessage(message string) error {
	query := `
	INSERT INTO messages(body)
	VALUES ($1)
	`
	_, err := DB.Exec(query, message)
	return err
}

func MessagesTableGetAllMessages() ([]string, error) {
	query := `
	SELECT body
	FROM messages
	ORDER BY id
	`
	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []string

	for rows.Next() {
		var body string

		err := rows.Scan(&body)
		if err != nil {
			return nil, err
		}

		messages = append(messages, body)
	}

	return messages, nil
}

func MessagesTableGetLatestMessages() ([]string, error) {
	fmt.Println("messagesTableGetAllMessages")
	t := time.Now().Add(-2 * time.Second)

	query := `
	SELECT *
	FROM messages
	WHERE created_at > $1
	`
	rows, err := DB.Query(query, t)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []string

	for rows.Next() {
		var id int
		var body string
		var createdAt time.Time

		err := rows.Scan(&id, &body, &createdAt)
		if err != nil {
			return nil, err
		}

		messages = append(messages, body)
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
