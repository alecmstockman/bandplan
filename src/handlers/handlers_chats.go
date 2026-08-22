package handlers

import (
	"bandplan/src/database"
	"bandplan/src/models"
	"fmt"
	"log"
	"net/http"
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
	log.Println("HandlerChatsPage")

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		log.Println("   Unable to get AuthContext: ", err)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	user := auth.User
	band := auth.CurrentBand

	userChats, err := database.ChatsTableGetChatsByUserID(user.UserID)
	if err != nil {
		log.Println("   Unable to get user chats from database: ", err)
		http.Error(w, "Unable to get user chats from database", http.StatusInternalServerError)
		return
	}

	data := models.ChatsDataPage{
		User:  user,
		Band:  band,
		Chats: userChats,
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

	messages, err := database.MessagesTableGetAllMessagesByChatID(chatID)
	if err != nil {
		log.Println("    HandlerHome: messages err: ", err)
		http.Error(w, "Unable to get messages", http.StatusInternalServerError)
		return
	}
	data := models.HomePageData{
		User:     user,
		Band:     band,
		ChatID:   chatID,
		Messages: messages,
	}

	err = h.Tmpl.ExecuteTemplate(w, "chat.html", data)
	if err != nil {
		log.Println("   template err:", err)
		return
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
}

func (h Handler) HandlerChatAddPage(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerChatCreatePage")

	_, err := HelperGetAuthContext(r)
	if err != nil {
		log.Println("   Unable to get AuthContext: ", err)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	// user := auth.User
	// band := auth.CurrentBand

	err = h.Tmpl.ExecuteTemplate(w, "chat_create.html", nil)
	if err != nil {
		log.Println("   Unable to go to create chat page: ", err)
		return
	}
}

func (h Handler) HandlerChatCreate(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerChatCreate")

	_, err := HelperGetAuthContext(r)
	if err != nil {
		log.Println("   Unable to get AuthContext: ", err)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	// user := auth.User
	// band := auth.CurrentBand

	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Println("   File too large: ", err)
		http.Error(w, "File too large", http.StatusBadRequest)
		return
	}

	var chat models.Chat

	err = database.ChatsTableCreateChat(chat)
	if err != nil {
		http.Redirect(w, r, "/songs", http.StatusSeeOther)
	}

	http.Redirect(w, r, "/songs", http.StatusSeeOther)

}
