package services

import (
	"bandplan/src/models"
	"context"
	"database/sql"
)

func (s Service) SessionSetFlash(ctx context.Context, db *sql.DB, token models.Session, message string) {

}
