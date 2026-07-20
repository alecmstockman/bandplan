package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDB() *sql.DB {
	databaseURL := os.Getenv("DATABASE_URL")

	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatal("Unable to open database: ", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("Unable to connect ot database: ", err)
	}

	log.Println("Connected to PostgreSQL")

	return db
}
