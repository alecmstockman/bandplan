package handlers

import (
	"bandplan/src/models"
	"bandplan/src/realtime"
	"bandplan/src/services"
	"bandplan/src/storage"
	"database/sql"
	"html/template"
)

type Handler struct {
	DB             *sql.DB
	Tmpl           *template.Template
	Storage        *storage.R2Storage
	Hub            *realtime.Hub
	SetlistService *services.SetlistService
}

type contextKey string

const AuthContextKey contextKey = "auth-context"

type AuthContext struct {
	User        models.User
	CurrentBand models.Band
}
