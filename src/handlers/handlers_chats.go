package handlers

import (
	"bandplan/src/database"
	"bandplan/src/models"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

func (h Handler) HandlerHome(w http.ResponseWriter, r *http.Request) {
	log.Printf("- HandlerHome")

	user, err := HelperGetAuthenticatedUser(r)
	if err != nil {
		log.Println("   Not authenticated: ", err)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	band, err := database.BandsTableGetBandByUserID(user.UserID)
	if err != nil {
		log.Println("   HandlerHome: Unable to get band by user id: ", err)
		http.Error(w, "Unable to load band", http.StatusInternalServerError)
		return
	}

	messages, err := database.MessagesTableGetAllMessagesByBandID(band.BandID)
	if err != nil {
		log.Println("    HandlerHome: messages err: ", err)
		http.Error(w, "Unable to get messages", http.StatusInternalServerError)
		return
	}
	data := models.HomePageData{
		User:     user,
		Band:     band,
		Messages: messages,
	}

	err = h.Tmpl.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		log.Println("   template err:", err)
		return
	}
}

func (h Handler) HandlerChatsPage(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerChatsPage")

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		log.Println("   Unable to get AuthContext: ", err)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	user := auth.User
	band := auth.CurrentBand

	chatPreviews, err := database.ChatsTableGetChatPreviewsByUserID(user.UserID)
	if err != nil {
		log.Println("   Unable to get user chats from database: ", err)
		http.Error(w, "Unable to get user chats from database", http.StatusInternalServerError)
		return
	}

	members, err := database.UsersTableGetUsersByBand(band.BandID)
	if err != nil {
		log.Println("   Unable to get users by band id: ", err)
		http.Error(w, "Unable to load band members", http.StatusInternalServerError)
		return
	}

	chatPreview, err := database.ChatsTableGetPrimaryChatPreviewByBandID(band.BandID)
	if err != nil {
		log.Println("   Unable to get ")
		http.Error(w, "Unable to load primary chat info", http.StatusInternalServerError)
		return
	}

	data := models.ChatsPageData{
		User:        user,
		Band:        band,
		Members:     members,
		PrimaryChat: chatPreview,
		Chats:       chatPreviews,
	}

	err = h.Tmpl.ExecuteTemplate(w, "chats.html", data)
	if err != nil {
		log.Println("   Unable to render chats page: ", err)
		http.Error(w, "Unable to load chats page", http.StatusInternalServerError)
		return
	}
}

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
	// data := models.HomePageData{
	// 	User:     user,
	// 	Band:     band,
	// 	ChatID:   chatID,
	// 	Messages: messages,
	// }

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

func (h Handler) HandlerChatLeave(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerChatLeave")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		log.Println("   Unable to get AuthContext: ", err)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	chatID := r.FormValue("chat-id")
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
	if chat.BandID != auth.CurrentBand.BandID {
		http.Error(w, "Chat does not belong to the current band", http.StatusForbidden)
		return
	}

	removed, err := database.ChatMembersTableRemoveMember(chatID, auth.User.UserID)
	if err != nil {
		http.Error(w, "Unable to leave chat", http.StatusInternalServerError)
		return
	}
	if !removed {
		http.Error(w, "You are not a member of this chat", http.StatusNotFound)
		return
	}

	w.Header().Set("HX-Redirect", "/chats")
	w.WriteHeader(http.StatusOK)
}

func (h Handler) HandlerChatDelete(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerChatDelete")
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		log.Println("   Unable to get AuthContext: ", err)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	chatID := r.FormValue("chat-id")
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
	if chat.BandID != auth.CurrentBand.BandID {
		http.Error(w, "Chat does not belong to the current band", http.StatusForbidden)
		return
	}

	deleted, err := database.ChatsTableDeleteChatByChatID(chatID)
	if err != nil {
		log.Printf("   Unable to delete chat id %v due to: %v\n", chatID, err)
		http.Error(w, "Unable to delete chat", http.StatusInternalServerError)
		return
	}
	if !deleted {
		http.Error(w, "Chat not found", http.StatusNotFound)
		return
	}

	w.Header().Set("HX-Redirect", "/chats")
	w.WriteHeader(http.StatusOK)
	return
}

func (h Handler) HandlerDelete(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerDelete")

	err := database.MessagesTableDeleteAll()
	if err != nil {
		http.Error(w, "Could not delete messages", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(""))
}

func (h Handler) HandlerMessages(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerMessages")

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		log.Println("   Unable to get AuthContext: ", err)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	user := auth.User
	band := auth.CurrentBand

	// user, err := HelperGetAuthenticatedUser(r)
	// if err != nil {
	// 	log.Println("   HandlerSend: Unable to get authenticated user: ", err)
	// 	http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
	// 	return
	// }

	// band, err := database.BandsTableGetBandByUserID(user.UserID)
	// if err != nil {
	// 	log.Println("   HandlerSend: Unable to get band by user id: ", err)
	// 	return
	// }

	messages, err := database.MessagesTableGetAllMessagesByBandID(band.BandID)
	if err != nil {
		http.Error(w, "Could not fetch latest messages", http.StatusInternalServerError)
		return
	}
	if len(messages) == 0 {
		return
	}
	for _, message := range messages {
		if user.UserID == message.UserID {
			html := fmt.Sprintf(`
				<li class="message-own">
					<div class="message-body">%s
						<div class="message-body-footer">%v</div>
					</div>
				</li>
				`, message.Body,
				message.CreatedAt.Format("3:04 PM"),
			)
			w.Write([]byte(html))

		} else {
			html := fmt.Sprintf(`
				<li class="test-message-other">

					<div class="test-message-sender-pic-box">
						<img
							class="test-message-sender-pic"
							src="%s"
							alt=""
							>
					</div>


					<div class="test-message-content">
						<div class="message-header-other">%s</div>

						<div class="message-body-other">%s

							<div class="message-body-footer">%s</div>
						</div>
					</div> 

				</li>
				`,
				HelperSmallImagePath(message.ProfileImagePath),
				message.UserName,
				message.Body,
				message.CreatedAt.Format("3:04 PM"),
			)
			w.Write([]byte(html))
		}
	}
}

func (h Handler) HandlerChatAddPage(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerChatAddPage")

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		log.Println("   Unable to get AuthContext: ", err)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	members, err := database.UsersTableGetUsersByBand(auth.CurrentBand.BandID)
	if err != nil {
		log.Println("   Unable to get users by band id: ", err)
		http.Error(w, "Unable to load band members", http.StatusInternalServerError)
		return
	}
	availableMembers := make([]models.User, 0, len(members))
	for _, member := range members {
		if member.UserID != auth.User.UserID {
			availableMembers = append(availableMembers, member)
		}
	}

	data := models.ChatsPageData{
		User:    auth.User,
		Band:    auth.CurrentBand,
		Members: availableMembers,
	}

	err = h.Tmpl.ExecuteTemplate(w, "chat_create.html", data)
	if err != nil {
		log.Println("   Unable to go to create chat page: ", err)
		return
	}
}

func (h Handler) HandlerChatSelectMember(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerChatSelectMember")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		log.Println("   Unable to get AuthContext: ", err)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	memberID := r.FormValue("member-id")
	members, err := database.UsersTableGetUsersByBand(auth.CurrentBand.BandID)
	if err != nil {
		log.Println("   Unable to get users by band id: ", err)
		http.Error(w, "Unable to load band members", http.StatusInternalServerError)
		return
	}

	var selectedMember models.User
	for _, member := range members {
		if member.UserID == memberID {
			selectedMember = member
			break
		}
	}

	if selectedMember.UserID == "" {
		http.Error(w, "Band member not found", http.StatusNotFound)
		return
	}
	if selectedMember.UserID == auth.User.UserID {
		http.Error(w, "The current user cannot be selected", http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.Tmpl.ExecuteTemplate(w, "chat_create_member_added", selectedMember); err != nil {
		log.Println("   Unable to render selected chat member: ", err)
		http.Error(w, "Unable to select band member", http.StatusInternalServerError)
	}
}

func (h Handler) HandlerChatRemoveMember(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerChatRemoveMember")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		log.Println("   Unable to get AuthContext: ", err)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	memberID := r.FormValue("member-id")
	members, err := database.UsersTableGetUsersByBand(auth.CurrentBand.BandID)
	if err != nil {
		log.Println("   Unable to get users by band id: ", err)
		http.Error(w, "Unable to load band members", http.StatusInternalServerError)
		return
	}

	var removedMember models.User
	for _, member := range members {
		if member.UserID == memberID {
			removedMember = member
			break
		}
	}

	if removedMember.UserID == "" {
		http.Error(w, "Band member not found", http.StatusNotFound)
		return
	}
	if removedMember.UserID == auth.User.UserID {
		http.Error(w, "The current user cannot be added to this list", http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.Tmpl.ExecuteTemplate(w, "chat_create_member_restored", removedMember); err != nil {
		log.Println("   Unable to render restored chat member: ", err)
		http.Error(w, "Unable to restore band member", http.StatusInternalServerError)
	}
}

func (h Handler) HandlerChatCreate(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerChatCreate")

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		log.Println("   Unable to get AuthContext: ", err)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	user := auth.User
	band := auth.CurrentBand

	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Println("   File too large: ", err)
		http.Error(w, "File too large", http.StatusBadRequest)
		return
	}
	chatName := strings.TrimSpace(r.FormValue("setlist-title"))
	if chatName == "" {
		http.Error(w, "Chat name is required", http.StatusBadRequest)
		return
	}
	slug := HelperMakeSlug(chatName)

	memberIDs := r.Form["member-id"]
	log.Println("   selected member IDs: ", memberIDs)

	chat := models.Chat{
		BandID:    band.BandID,
		Name:      chatName,
		Slug:      slug,
		IsPrimary: false,
		CreatedBy: user.UserID,
		UpdatedBy: user.UserID,
	}

	_, err = database.ChatsTableCreateChat(chat, memberIDs)
	if err != nil {
		log.Println("   Unable to create chat: ", err)
		http.Error(w, "Unable to create chat", http.StatusInternalServerError)
		return
	}

	// http.Redirect(w, r, "/chats", http.StatusSeeOther)
	w.Header().Set("HX-Redirect", "/chats")
	w.WriteHeader(http.StatusOK)

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
	fmt.Println("-----------------------------")
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
		imagePath, err = h.SetlistService.HelperCreatePermSetlistImage(
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
	fmt.Println("\n--------------------------------")
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

		previewURL, err = h.HelperSaveTempImage(r.Context(), file, imageID, band.Slug, "chat")
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
	fmt.Println("\n-----------------------------------")
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
