package handlers

import (
	"bandplan/src/database"
	"bandplan/src/models"
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func (h Handler) HandlerChatSettings(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerChatSettings")

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		log.Println("   Unable to get AuthContext: ", err)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	chatID := r.URL.Query().Get("id")
	if chatID == "" {
		http.Error(w, "Chat ID is required", http.StatusBadRequest)
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
		User: auth.User,
		Band: auth.CurrentBand,
		Chat: chat,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.Tmpl.ExecuteTemplate(w, "chat_settings.html", pageData); err != nil {
		log.Println("   Unable to render chat settings: ", err)
		http.Error(w, "Unable to load chat settings", http.StatusInternalServerError)
	}
}

func (h Handler) HandlerChatAddImagePage(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerChatAddImage")

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		log.Println("   Unable to get AuthContext: ", err)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	chatID := r.URL.Query().Get("id")

	chat, err := database.ChatsTableGetChatByChatID(chatID)
	if err != nil {
		log.Println("   Unable to get chat by chat id: ", err)
		http.Error(w, "Unable to get chat info", http.StatusInternalServerError)
		return
	}

	data := models.ChatPageData{
		User: auth.User,
		Band: auth.CurrentBand,
		Chat: chat,
	}

	err = h.Tmpl.ExecuteTemplate(w, "chat_image.html", data)
	if err != nil {
		log.Println("   Unable to get setlist_item template: ", err)
		http.Error(w, "Unable to get setlist_item template", http.StatusInternalServerError)
	}
}

func (h Handler) HandlerChatImageSave(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerChatImageSave")

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		log.Println("   Unable to get AuthContext: ", err)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	band := auth.CurrentBand

	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	chatID := r.FormValue("chat-id")
	if chatID == "" {
		http.Error(w, "Chat ID is required", http.StatusBadRequest)
		return
	}

	existingChat, err := database.ChatsTableGetChatByChatID(chatID)
	if err != nil {
		log.Println("   Unable to get chat by chat id: ", err)
		http.Error(w, "Unable to get chat info", http.StatusInternalServerError)
		return
	}

	imagePath := existingChat.ImagePath
	imageID := existingChat.ImageID
	temporaryArtworkID := r.FormValue("temporary-artwork-id")

	if temporaryArtworkID != "" {
		imagePath, err = h.Services.ServiceCreatePermChatImage(
			r.Context(),
			temporaryArtworkID,
			band.Slug,
			existingChat.Slug,
		)
		if err != nil {
			log.Println("   Unable to save temporary chat image versions: ", err)
			http.Error(w, "Could not save tempoerary iamage versions", http.StatusInternalServerError)
			return
		}
		imageID = temporaryArtworkID
	} else if r.FormValue("remove-artwork") == "true" {
		imageID = ""
		imagePath = ""
	}

	chat := models.Chat{
		ChatID:    chatID,
		ImageID:   imageID,
		ImagePath: imagePath,
	}

	updated, err := database.ChatsTableUpdateChat(chat)
	if err != nil {
		log.Println("   Unable to update chat in database: ", err)
		http.Error(w, "Unable to update chat in database", http.StatusInternalServerError)
		return
	}

	if !updated {
		http.Error(w, "Chat not found", http.StatusNotFound)
		return
	}

	log.Println("   Updated chat: ", chat.Name, chat.ChatID)

	w.Header().Set("HX-Redirect", "/chat?id="+chatID)
	w.WriteHeader(http.StatusOK)
	return
}

func (h Handler) HandlerChatTempArt(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerChatTempArt")

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		log.Println("   Unable to get AuthContext: ", err)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	band := auth.CurrentBand

	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Println("   Error parsing from while creating a new chat: ", err)
		http.Error(w, "Unable to parse from", http.StatusBadRequest)
		return
	}

	imageID := ""
	previewURL := ""

	file, _, err := r.FormFile("artwork-path")
	if err != nil {
		if err != http.ErrMissingFile {
			log.Println("   Unable to read artwork file:", err)
			http.Error(w, "Unable to read artwork", http.StatusBadRequest)
			return
		}

		log.Println("   No setlist artwork uploaded")
	} else {
		defer file.Close()

		imageID = uuid.New().String()

		previewURL, err = h.Services.ServiceSaveTempImage(r.Context(), file, imageID, band.Slug, "chat")
		if err != nil {
			http.Error(w, "Could not save artwork versions", http.StatusInternalServerError)
			return
		}
	}

	data := models.ArtworkPreviewData{
		ArtworkID:  imageID,
		PreviewURL: previewURL,
	}

	err = h.Tmpl.ExecuteTemplate(w, "chat_artwork_preview", data)
	if err != nil {
		http.Error(w, "Unable to render preview", http.StatusInternalServerError)
		return
	}
}

func (h Handler) HandlerChatTempArtDelete(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerChatTempArtDelete")

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		log.Println("   Unable to get AuthContext: ", err)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	user := auth.User
	band := auth.CurrentBand

	data := models.SongDownloadData{
		User: user,
		Band: band,
	}

	err = h.Tmpl.ExecuteTemplate(w, "chat_artwork_reset", data)
	if err != nil {
		log.Println("   Unable to exececute : ", err)
		http.Error(w, "Unable to render preview", http.StatusInternalServerError)
		return
	}
}
