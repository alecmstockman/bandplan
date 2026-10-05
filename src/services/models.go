package services

import (
	"bandplan/src/storage"
	"database/sql"
)

type Service struct {
	DB                *sql.DB
	Storage           *storage.R2Storage
	LoginEmailLimiter LoginLimiter
}

type CreateSetlistInput struct {
	Title     string
	Notes     string
	TempArtID string
}
