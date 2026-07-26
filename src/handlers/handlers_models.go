package handlers

import (
	"bandplan/src/models"
	"bandplan/src/realtime"
	"database/sql"
	"html/template"
)

type Handler struct {
	DB   *sql.DB
	Tmpl *template.Template
	Hub  *realtime.Hub
}

type contextKey string

const AuthContextKey contextKey = "auth-context"

type AuthContext struct {
	User        models.User
	CurrentBand models.Band
}
