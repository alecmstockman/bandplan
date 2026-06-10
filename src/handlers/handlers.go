package handlers

import (
	"bandplan/src/database"
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
)

type Handler struct {
	DB   *sql.DB
	Tmpl *template.Template
}

func (h Handler) HomeHandler(w http.ResponseWriter, r *http.Request) {
	messages, err := database.MessagesTableGetAllMessages()
	if err != nil {
		fmt.Println("Unable to get messages from db")
	}

	h.Tmpl.Execute(w, messages)
}

func (h Handler) SendHandler(w http.ResponseWriter, r *http.Request) {
	message := r.FormValue("message")

	err := database.MessagesTableInsertMessage(message)
	if err != nil {
		http.Error(w, "Could not save message", http.StatusInternalServerError)
		return
	}

	html := fmt.Sprintf(
		"<li>%s</li>",
		message,
	)

	w.Write([]byte(html))
}

func (h Handler) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	err := database.MessagesTableDeleteAll()
	if err != nil {
		http.Error(w, "Could not delete messages", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(""))
	return
}

func (h Handler) MessagesHandler(w http.ResponseWriter, r *http.Request) {
	messages, err := database.MessagesTableGetAllMessages()

	if err != nil {
		http.Error(w, "Could not fetch latest messages", http.StatusInternalServerError)
		return
	}

	for _, msg := range messages {
		html := fmt.Sprintf("<li>%s</li>", msg)
		w.Write([]byte(html))
	}
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	band := r.FormValue("band")
	email := r.FormValue("email")
	password := r.FormValue("password")

	fmt.Println(name, band, email, password)

	w.Write([]byte("User registered!"))
}
