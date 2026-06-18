package handlers

import (
	"bandplan/src/database"
	"bandplan/src/models"
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type Handler struct {
	DB   *sql.DB
	Tmpl *template.Template
}

type HomePageData struct {
	User     models.User
	Messages []models.Message
}

func (h Handler) HandlerHome(w http.ResponseWriter, r *http.Request) {
	log.Printf("- HandlerHome: %s %s\n", r.Method, r.URL.Path)

	user, err := HelperGetAuthenticatedUser(r)
	if err != nil {
		fmt.Println("Not authenticated: ", err)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	messages, err := database.MessagesTableGetAllMessages()
	if err != nil {
		fmt.Println("Unable to get messages from db: ", err)
		http.Error(w, "Unable to get messages", http.StatusInternalServerError)
		return
	}
	data := HomePageData{
		User:     user,
		Messages: messages,
	}

	err = h.Tmpl.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		fmt.Println("template err:", err)
	}
}

func (h Handler) HandlerChatPage(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerChatPage")
	http.ServeFile(w, r, "templates/index.html")
	return
}

func (h Handler) HandlerSend(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerSend")
	input := r.FormValue("message")

	if input == "" {
		return
	}

	user, err := HelperGetAuthenticatedUser(r)
	if err != nil {
		fmt.Println("HandlerSend: Unable to get authenticated user: ", err)
		return
	}

	message := models.Message{
		MessageID: uuid.New().String(),
		UserID:    user.UserID,
		UserName:  user.Name,
		Body:      input,
	}

	err = database.MessagesTableInsertMessage(message)
	if err != nil {
		fmt.Println("- err: ", err)
		http.Error(w, "Could not save message", http.StatusInternalServerError)
		return
	}

	if user.Name == message.UserName {
		html := fmt.Sprintf(`
			<li class="message-own">
				<div class="message-header">%s</div>
				<div class="message-body">%s</div>
			</li>
			`, message.UserName, message.Body,
		)
		w.Write([]byte(html))
	} else {
		html := fmt.Sprintf(`
			<li class="message-other">
				<div class="message-header">%s</div>
				<div class="message-body">%s</div>
			</li>
			`, message.UserName, message.Body,
		)
		w.Write([]byte(html))
	}
}

func (h Handler) HandlerDelete(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerDelete")
	err := database.MessagesTableDeleteAll()
	if err != nil {
		http.Error(w, "Could not delete messages", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(""))
	return
}

func (h Handler) HandlerMessages(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerMessages")
	messages, err := database.MessagesTableGetAllMessages()
	if err != nil {
		fmt.Printf("- Handler messages: MessagesTableGetAll err: ", err)
		http.Error(w, "Could not fetch latest messages", http.StatusInternalServerError)
		return
	}

	user, err := HelperGetAuthenticatedUser(r)
	if err != nil {
		fmt.Println(" - Handler messages: HelperGetAuthenticateduser err: ", err)
		return
	}

	for _, message := range messages {
		if user.Name == message.UserName {
			html := fmt.Sprintf(`
				<li class="message-own">
					<div class="message-header">%s</div>
					<div class="message-body">%s</div>
				</li>
				`, message.UserName, message.Body,
			)
			w.Write([]byte(html))
		} else {
			html := fmt.Sprintf(`
				<li class="message-other">
					<div class="message-header">%s</div>
					<div class="message-body">%s</div>
				</li>
				`, message.UserName, message.Body,
			)
			w.Write([]byte(html))
		}
	}
}

func (h Handler) HandlerLogout(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerLogout")
	http.SetCookie(w, &http.Cookie{
		Name:   "session_token",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	w.Header().Set("HX-Redirect", "/login")
	w.WriteHeader(http.StatusOK)
	return
}
