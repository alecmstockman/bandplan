package handlers

import (
	"bandplan/src/database"
	"bandplan/src/models"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
)

func (h Handler) HandlerChatPage(w http.ResponseWriter, r *http.Request) {
	log.Printf("- HandlerChatPage")

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		log.Println("   Unable to get AuthContext: ", err)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	user := auth.User
	band := auth.CurrentBand

	chatID := r.URL.Query().Get("id")
	if chatID == "" {
		http.Error(w, "Chat ID is required", http.StatusBadRequest)
		return
	}

	messages, err := database.MessagesTableGetAllMessagesByChatID(chatID)
	if err != nil {
		log.Println("    HandlerHome: messages err: ", err)
		http.Error(w, "Unable to get messages", http.StatusInternalServerError)
		return
	}

	chat, err := database.ChatsTableGetChatByChatID(chatID)
	if err != nil {
		log.Println("   Unable to get chat: ", err)
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Chat not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Unable to get chat", http.StatusInternalServerError)
		return
	}

	pageData := models.ChatPageData{
		User:     user,
		Band:     band,
		Chat:     chat,
		Messages: messages,
	}

	templateName := "chat.html"
	if r.Header.Get("HX-Request") == "true" {
		templateName = "chat_main_content"
	}

	err = h.Tmpl.ExecuteTemplate(w, templateName, pageData)
	if err != nil {
		log.Println("   template err:", err)
		return
	}
}

func (h Handler) HandlerChatMessageReaction(w http.ResponseWriter, r *http.Request) {
	fmt.Println("--------------------------")
	log.Println("- HandlerChatMessageReaction")

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		log.Println("   Unable to get AuthContext: ", err)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	reaction := r.FormValue("reaction")
	messageID := r.FormValue("message-id")
	chatID := r.FormValue("chat-id")

	fmt.Println("reaction: ", reaction)
	fmt.Println("message-id: ", messageID)
	fmt.Println("chat-id: ", chatID)

	err = database.MessagesReactionTableAddReaction(messageID, auth.User.UserID, reaction)
	if err != nil {
		log.Println("Unable to add reaction to message_reactions table: ", err)
		http.Error(w, "Unable to save reaction", http.StatusInternalServerError)
		return
	}

}

func (h Handler) HandlerChatMessageReply(w http.ResponseWriter, r *http.Request) {
	fmt.Println("--------------------------")
	log.Println("- HandlerChatMessageReply")
}

func (h Handler) HandlerChatMessageCopy(w http.ResponseWriter, r *http.Request) {
	fmt.Println("--------------------------")
	log.Println("- HandlerChatMessageCopy")
}

func (h Handler) HandlerChatMessagePinAdd(w http.ResponseWriter, r *http.Request) {
	fmt.Println("--------------------------")
	log.Println("- HandlerChatMessagePinAdd")
}

func (h Handler) HandlerChatMessagePinRemove(w http.ResponseWriter, r *http.Request) {
	fmt.Println("--------------------------")
	log.Println("- HandlerChatMessagePinRemove")
}

func (h Handler) HandlerChatMessageDelete(w http.ResponseWriter, r *http.Request) {
	fmt.Println("--------------------------")
	log.Println("- HandlerChatMessageDelete")
}
