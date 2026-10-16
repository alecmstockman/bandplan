package handlers

import (
	"bandplan/src/database"
	requestlog "bandplan/src/logging"
	"bandplan/src/models"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"strings"
)

func (h Handler) HandlerChatPage(w http.ResponseWriter, r *http.Request) {
	log.Printf("- HandlerChatPage")

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		slog.Error(
			"unable to load auth context",
			"request_id", requestlog.GetRequestID(r.Context()),
			"error", err,
		)
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

	chat, err := database.ChatsTableGetChatByChatIDForMember(chatID, band.BandID, user.UserID)
	if err != nil {
		log.Println("   Unable to get chat: ", err)
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Chat not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Unable to get chat", http.StatusInternalServerError)
		return
	}

	messages, err := database.MessagesTableGetAllMessagesByChatID(chatID, band.BandID, user.UserID)
	if err != nil {
		log.Println("    HandlerChatPage: messages err: ", err)
		http.Error(w, "Unable to get messages", http.StatusInternalServerError)
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
		templateName = "chat-main-content"
	}

	err = h.Tmpl.ExecuteTemplate(w, templateName, pageData)
	if err != nil {
		log.Println("   template err:", err)
		return
	}
}

func (h Handler) HandlerChatMessageReaction(w http.ResponseWriter, r *http.Request) {

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		slog.Error(
			"unable to load auth context",
			"request_id", requestlog.GetRequestID(r.Context()),
			"error", err,
		)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	user := auth.User
	band := auth.CurrentBand

	reaction := strings.TrimSpace(r.FormValue("reaction"))
	messageID := strings.TrimSpace(r.FormValue("message-id"))
	if messageID == "" {
		http.Error(w, "Message ID is required", http.StatusBadRequest)
		return
	}

	switch reaction {
	case "heart", "laugh", "shocked", "anger", "sad", "horns":
	default:
		http.Error(w, "Invalid reaction", http.StatusBadRequest)
		return
	}

	_, err = database.MessageReactionsTableAddReaction(messageID, user.UserID, band.BandID, reaction)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.NotFound(w, r)
			return
		}
		log.Println("Unable to add reaction to message_reactions table: ", err)
		http.Error(w, "Unable to save reaction", http.StatusInternalServerError)
		return
	}

	reactions, err := database.MessageReactionsTableGetReactionsByMessageID(messageID, user.UserID, band.BandID)
	if err != nil {
		log.Println("   Unable to get message reactions from database: ", err)
		http.Error(w, "Unable to get load message reactions", http.StatusNotFound)
		return
	}

	for _, reaction := range reactions {
		var emoji string

		switch reaction.Reaction {
		case "heart":
			emoji = "❤️"
			html := fmt.Sprintf(`
				<div class="chat-reaction">%s</div>
			`, emoji)
			w.Write([]byte(html))
		case "laugh":
			emoji = "😂"
			html := fmt.Sprintf(`
				<div class="chat-reaction">%s</div>
			`, emoji)
			w.Write([]byte(html))
		case "shocked":
			emoji = "😮"
			html := fmt.Sprintf(`
				<div class="chat-reaction">%s</div>
			`, emoji)
			w.Write([]byte(html))
		case "anger":
			emoji = "😡"
			html := fmt.Sprintf(`
				<div class="chat-reaction">%s</div>
			`, emoji)
			w.Write([]byte(html))
		case "sad":
			emoji = "😢"
			html := fmt.Sprintf(`
				<div class="chat-reaction">%s</div>
			`, emoji)
			w.Write([]byte(html))
		case "horns":
			emoji = "🤘"
			html := fmt.Sprintf(`
				<div class="chat-reaction">%s</div>
			`, emoji)
			w.Write([]byte(html))
		}
	}
}

func (h Handler) HandlerChatMessageReply(w http.ResponseWriter, r *http.Request) {
	fmt.Println("--------------------------")
	log.Println("- HandlerChatMessageReply")
}

func (h Handler) HandlerChatMessageEdit(w http.ResponseWriter, r *http.Request) {
	fmt.Println("--------------------------")
	log.Println("- HandlerChatMessageEdit")

}

func (h Handler) HandlerChatMessagePinAdd(w http.ResponseWriter, r *http.Request) {
	fmt.Println("--------------------------")
	log.Println("- HandlerChatMessagePinAdd")

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		slog.Error(
			"unable to load auth context",
			"request_id", requestlog.GetRequestID(r.Context()),
			"error", err,
		)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	user := auth.User
	band := auth.CurrentBand

	chatID := strings.TrimSpace(r.FormValue("chat-id"))
	messageID := strings.TrimSpace(r.FormValue("message-id"))
	if chatID == "" || messageID == "" {
		http.Error(w, "Chat ID and message ID are required", http.StatusBadRequest)
		return
	}

	updated, err := database.MessagesTablePinMessage(messageID, chatID, band.BandID, user.UserID)
	if err != nil {
		log.Println("   Unable to pin message: ", err)
		http.Error(w, "Unable to pin message", http.StatusInternalServerError)
		return
	}
	if !updated {
		http.NotFound(w, r)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h Handler) HandlerChatMessagePinRemove(w http.ResponseWriter, r *http.Request) {
	fmt.Println("--------------------------")
	log.Println("- HandlerChatMessagePinRemove")

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		slog.Error(
			"unable to load auth context",
			"request_id", requestlog.GetRequestID(r.Context()),
			"error", err,
		)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	chatID := strings.TrimSpace(r.FormValue("chat-id"))
	messageID := strings.TrimSpace(r.FormValue("message-id"))
	if chatID == "" || messageID == "" {
		http.Error(w, "Chat ID and message ID are required", http.StatusBadRequest)
		return
	}

	updated, err := database.MessagesTableUnPinMessage(
		messageID,
		chatID,
		auth.CurrentBand.BandID,
		auth.User.UserID,
	)
	if err != nil {
		log.Println("   Unable to unpin message: ", err)
		http.Error(w, "Unable to unpin message", http.StatusInternalServerError)
		return
	}
	if !updated {
		http.NotFound(w, r)
		return
	}

	messages, err := database.MessagesTableGetPinnedMessagesByChatID(
		chatID,
		auth.CurrentBand.BandID,
		auth.User.UserID,
	)
	if err != nil {
		log.Println("   Unable to get pinned messages: ", err)
		http.Error(w, "Unable to get pinned messages", http.StatusInternalServerError)
		return
	}

	data := models.PinnedChatsPageData{
		User:     auth.User,
		Band:     auth.CurrentBand,
		ChatID:   chatID,
		Messages: messages,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.Tmpl.ExecuteTemplate(w, "chat_pinned_messages", data); err != nil {
		log.Println("   Template error: ", err)
		http.Error(w, "Unable to load pinned messages", http.StatusInternalServerError)
	}
}

func (h Handler) HandlerChatMessageDelete(w http.ResponseWriter, r *http.Request) {
	fmt.Println("--------------------------")
	log.Println("- HandlerChatMessageDelete")
}

func (h Handler) HandlerPinnedChats(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerPinnedChats")

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		slog.Error(
			"unable to load auth context",
			"request_id", requestlog.GetRequestID(r.Context()),
			"error", err,
		)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	user := auth.User
	band := auth.CurrentBand
	chatID := strings.TrimSpace(r.FormValue("chat_id"))
	if chatID == "" {
		http.Error(w, "Chat ID is required", http.StatusBadRequest)
		return
	}

	if _, err := database.ChatsTableGetChatByChatIDForMember(chatID, band.BandID, user.UserID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.NotFound(w, r)
			return
		}
		log.Println("   Unable to authorize pinned messages: ", err)
		http.Error(w, "Unable to get pinned messages", http.StatusInternalServerError)
		return
	}

	messages, err := database.MessagesTableGetPinnedMessagesByChatID(chatID, band.BandID, user.UserID)
	if err != nil {
		log.Println("   Unable to get pinned messages: ", err)
		http.Error(w, "Unable to get pinned messages", http.StatusInternalServerError)
		return
	}

	data := models.PinnedChatsPageData{
		User:     user,
		Band:     band,
		ChatID:   chatID,
		Messages: messages,
	}

	err = h.Tmpl.ExecuteTemplate(w, "chat_pinned_messages", data)
	if err != nil {
		log.Println("   template err:", err)
		return
	}
}
