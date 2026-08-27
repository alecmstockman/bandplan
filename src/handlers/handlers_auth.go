package handlers

import (
	"bandplan/src/database"
	"bandplan/src/models"
	"fmt"
	"log"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func (h Handler) HandlerRegisterPage(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerRegisterPage")

	user, err := HelperGetAuthenticatedUser(r)
	if err == nil {
		log.Println("   User already logged in: ", user.Name)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	code := r.FormValue("access-code")

	band := models.Band{}

	bandID, err := database.AccessCodesTableValidateCode(code)
	if bandID == "" {
		log.Println("   Unable to validate code: ", err)
	} else {
		band, err = database.BandsTableGetBandByBandID(bandID)
		if err != nil {
			log.Println("   Unable to get band by band ID: ", err)
		}
	}

	h.Tmpl.ExecuteTemplate(w, "register.html", band)
	return
}

func (h Handler) HandlerRegister(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerRegister")

	if r.Method == http.MethodPost {
		name := strings.TrimSpace(r.FormValue("name"))
		slug := HelperMakeSlug(name)
		displayName := strings.TrimSpace(r.FormValue("display-name"))
		bandNameEntry := strings.TrimSpace(r.FormValue("band"))
		email := HelperNormalizeEmail(r.FormValue("email"))
		password := r.FormValue("password")
		isAdmin := true

		bandName := HelperProcessBandNameEntry(bandNameEntry)
		_, err := database.BandsTableGetBandByName(bandName)
		if err == nil {
			isAdmin = false
		}

		user, err := database.UsersTableCreateUser(name, displayName, slug, email, password, isAdmin)
		if err != nil {
			log.Println("   register err: ", err)
			http.Error(w, "Could not create user", http.StatusInternalServerError)
			return
		}

		band, err := database.BandsTableGetBandByName(bandName)
		if err != nil {
			bandSlug := HelperMakeSlug(bandName)
			band, err = database.BandsTableCreateBand(bandName, user.UserID, bandSlug)
			if err != nil {
				log.Println("   register err: ", err)
				http.Error(w, "Could not create band", http.StatusInternalServerError)
				return
			}

			chatName := fmt.Sprintf("%s (Band Chat)", band.Name)
			chatSlug := HelperMakeSlug(chatName)
			_, err := database.ChatsTableCreatePrimaryBandChat(band.BandID, chatName, chatSlug, user.UserID)

			log.Println("   Created primary band chat id: ", err)

			if err != nil {
				log.Printf("   Unable to create primary band chat for band: %v, bandID: %v, error: %v", band.Name, band.BandID, err)
				http.Error(w, "Unable to create primary band chat, please try again", http.StatusInternalServerError)
				return
			}
			// err = database.ChatMembersTableAddMember(chatID, user.UserID)
			// if err != nil {
			// 	log.Println("   Unable to add member to the chat_members table: ", err)
			// 	http.Error(w, "Unable to add member to the chat_members table", http.StatusInternalServerError)
			// 	return
			// }
		}

		chatID, err := database.ChatsTableGetPrimaryChatIDByBandID(band.BandID)
		if err != nil {
			log.Println("   Unable to get primary chatID by bandID: ", err)
			http.Error(w, "Unable to get primary chatID by bandID", http.StatusInternalServerError)
			return
		}

		err = database.ChatMembersTableAddMember(chatID, user.UserID)
		if err != nil {
			log.Println("   Unable to add new user to chat_members table: ", err)
			http.Error(w, "Unable to add new user to chat_members table", http.StatusInternalServerError)
			return
		}

		err = database.BandMembersCreateMember(band.BandID, user.UserID)
		if err != nil {
			log.Printf("   Unable to add user: %v to band members table: %v\n", user.UserID, err)
			http.Error(w, "Could not create band member", http.StatusInternalServerError)
			return
		}

		w.Header().Set("HX-Redirect", "/login")
		w.WriteHeader(http.StatusOK)
		return
	}
	return
}

func (h Handler) HandlerLoginPage(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerLoginPage")

	user, err := HelperGetAuthenticatedUser(r)
	if err == nil {
		log.Println("   Already logged in: ", user.Name)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	h.Tmpl.ExecuteTemplate(w, "login.html", nil)
	return
}

func (h Handler) HandlerLogin(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerLogin")

	email := HelperNormalizeEmail(r.FormValue("email"))
	password := r.FormValue("password")

	user, err := database.UsersTableGetUserByEmail(email)
	if err != nil {
		log.Println("   HandlerLogin: Unable to get user: ", err)
		w.Write([]byte("Invalid email or password"))
		return
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	)
	if err != nil {
		log.Println("   Invalid email or password")
		w.Write([]byte("* Invalid email or password * "))
		return
	}

	token, err := HelperGenerateSessionToken()
	if err != nil {
		log.Println("   Unabe to generate token: ", err)
		return
	}

	params := models.CreateSessionParams{
		UserID: user.UserID,
		Token:  token,
	}

	session, err := database.SessionsTableCreateSession(params)
	if err != nil {
		log.Println("   Unable to create session: ", err)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    session.Token,
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	cookie, err := r.Cookie("session_token")
	if err != nil {
		log.Print("\n  - err with session_token: ", err)
	}
	log.Println("\n  - cookie: ", cookie)

	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusOK)
	return
}

func (h Handler) HandlerUserAgreementPage(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerUserAgreementPage")

	h.Tmpl.ExecuteTemplate(w, "user-agreement.html", nil)
	return
}

func (h Handler) HandlerUserAgreement(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerUserAgreement")

	h.Tmpl.ExecuteTemplate(w, "login.html", nil)
	return
}

func (h Handler) HandlerTermsPage(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerTermPage")

	err := h.Tmpl.ExecuteTemplate(w, "terms.html", nil)
	if err != nil {
		log.Println("Unable to render terms page:", err)
		http.Error(
			w,
			"Unable to load page",
			http.StatusInternalServerError,
		)
	}
	return
}

func (h Handler) HandlerPrivacyPage(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerPrivacyPage")

	h.Tmpl.ExecuteTemplate(w, "privacy.html", nil)
	return
}

func (h Handler) HandlerAccessCodePage(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerAccessCode")

	h.Tmpl.ExecuteTemplate(w, "access.html", nil)
	return
}

func (h Handler) HandlerCreateAccessCode(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerCreateAccessCode")

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		log.Println("   Unable to get AuthContext: ", err)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	user := auth.User
	band := auth.CurrentBand

	code, err := database.AccessCodesTablesCreateCode(band.BandID, user.UserID)
	if err != nil {
		http.Error(w, "Unable to generate access code", http.StatusInternalServerError)
		return
	}

	html := fmt.Sprintf(`
		<div class="admin-access-code-box">
            <span class="admin-access-code">%v</span>
			<button
              type="button"
              class="admin-display-field-copy"
              onclick="copyToClipboard('%s', this)">
              <span class="copy-icon">
                  <svg 
					xmlns="http://www.w3.org/2000/svg" 
					width="20" height="20" 
					viewBox="0 0 24 24" fill="none" 
					stroke="currentColor" stroke-width="2" 
					stroke-linecap="round" stroke-linejoin="round" 
					class="lucide lucide-copy-icon lucide-copy">
					<rect width="14" height="14" x="8" y="8" rx="2" ry="2"/>
					<path d="M4 16c-1.1 0-2-.9-2-2V4c0-1.1.9-2 2-2h10c1.1 0 2 .9 2 2"/>
				</svg>
              </span>

              <span class="check-icon">
                  <svg 
					xmlns="http://www.w3.org/2000/svg" 
					width="20" height="20" 
					viewBox="0 0 24 24" fill="none" 
					stroke="currentColor" stroke-width="2" 
					stroke-linecap="round" stroke-linejoin="round" 
					class="lucide lucide-check-icon lucide-check">
					<path d="M20 6 9 17l-5-5"/>
				</svg>
              </span>
            </button>
		</div>
	`, code, code)
	w.Write([]byte(html))
	return
}
