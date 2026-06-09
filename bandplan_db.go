package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func connectDB() *sql.DB {
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

func createTables(db *sql.DB) {
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

func messagesTableInsertMessage(message string) error {
	query := `
	INSERT INTO messages(body)
	VALUES ($1)
	`
	_, err := db.Exec(query, message)
	return err
}

func messagesTableGetAllMessages() ([]string, error) {
	query := `
	SELECT body
	FROM messages
	ORDER BY id
	`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}

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

func messagesTableDeleteAll() error {
	query := `
	DELETE FROM messages
	`
	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}
