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

	fmt.Printf(
		"- HandlerHome method=%s path=%s",
		r.Method,
		r.URL.Path,
	)

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
		log.Println("   Unable to get authenticated user: ", err)
		return
	}
	band, err := database.BandsTableGetBandByUserID(user.UserID)
	if err != nil {
		log.Println("   HandlerSend: Unable to get band by user id: ", err)
		return
	}

	message, err := database.MessagesTableCreateMessage(band.BandID, user.UserID, user.Name, input)
	if err != nil {
		log.Println("   Could not create message in message table: ", err)
		http.Error(w, "Could not save message", http.StatusInternalServerError)
		return
	}

	if user.UserID == message.UserID {
		html := fmt.Sprintf(`
			<li class="message-own">
				<div class="message-body">%s
					<div class="message-body-footer">%v</div>
				</div>
			</li>
			`,
			message.Body,
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
				</div>


				<div class="test-message-content">
					<div class="message-header-other">%s</div>

					<div class="message-body-other">%s
						<div class="message-body-footer">%s</div>
					</div>
				</div> 

			</li>
			`,
			message.ProfileImagePath,
			message.UserName,
			message.Body,
			message.CreatedAt.Format("3:04 PM"),
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
}

func (h Handler) HandlerMessages(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerMessages")

	user, err := HelperGetAuthenticatedUser(r)
	if err != nil {
		log.Println("   HandlerSend: Unable to get authenticated user: ", err)
		w.Header().Set("HX-Redirect", "/login")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	band, err := database.BandsTableGetBandByUserID(user.UserID)
	if err != nil {
		log.Println("   HandlerSend: Unable to get band by user id: ", err)
		return
	}

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
