package handlers

import (
	"bandplan/src/database"
	"bandplan/src/helpers"
	requestlog "bandplan/src/logging"
	"bandplan/src/models"
	"fmt"
	"log"
	"log/slog"
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

func (h Handler) HandlerRegisterPageOne(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerRegisterPageOne")

	err := h.Tmpl.ExecuteTemplate(w, "register-page1-access-code.html", nil)
	if err != nil {
		log.Println("Unable to execute register-page1-access-code.html")
		return
	}
}

func (h Handler) HandlerRegisterPageTwo(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerRegisterPageTwo")

	accessCode := r.FormValue("access-code")

	fmt.Println("accessCode: ", accessCode)

	err := h.Tmpl.ExecuteTemplate(w, "register-page2-name.html", accessCode)
	if err != nil {
		log.Println("Unable to execute register-page1-access-code.html", err)
		return
	}
}

func (h Handler) HandlerRegisterPageThree(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerRegisterPageThree")

	firstName := strings.TrimSpace(r.FormValue("first-name"))
	lastName := strings.TrimSpace(r.FormValue("last-name"))
	displayName := strings.TrimSpace(r.FormValue("display-name"))
	timezone := r.FormValue("timezone")
	accessCode := r.FormValue("access-code")

	fmt.Println("first name:   ", firstName)
	fmt.Println("last name:    ", lastName)
	fmt.Println("display name: ", displayName)
	fmt.Println("timezone:     ", timezone)
	fmt.Println("access code:  ", accessCode)

	newUser, err := h.Services.RegistrationCreateUserProfile(firstName, lastName, displayName, timezone, accessCode)
	if err != nil {
		log.Println("   unable to create user registration profile", err)
		http.Error(w, "unable to create user registration profile", http.StatusInternalServerError)
		return
	}

	err = h.Tmpl.ExecuteTemplate(w, "register-page3-password.html", newUser)
	if err != nil {
		log.Println("Unable to execute register-page3-password.html", err)
		return
	}
}

func (h Handler) HandlerRegisterPageFour(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerRegisterPageThree")

	password := r.FormValue("password")
	passwordConfirmation := r.FormValue("password-confirmation")

	if password != passwordConfirmation {
		log.Println("passwords do not match")
		http.Error(w, "passwords do not match", http.StatusBadRequest)
		return
	}

	registrationID := r.FormValue("registration-id")
	accessCodeHash := r.FormValue("access-code-hash")

	fmt.Println("registration-id: ", registrationID)
	fmt.Println("access-code-hash: ", accessCodeHash)

	_, err := h.Services.RegistrationSavePassword(registrationID, password)
	if err != nil {
		log.Println("   Unable to save password", err)
		http.Error(w, "Unable to save password", http.StatusInternalServerError)
		return
	}

	fmt.Println("password:     ", password)
	fmt.Println("password Con: ", passwordConfirmation)

	err = h.Tmpl.ExecuteTemplate(w, "register-page4-band.html", registrationID)
	if err != nil {
		log.Println("Unable to execute register-page4-band.html", err)
		return
	}
}

func (h Handler) HandlerRegisterPageFive(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerRegisterPageThree")

	err := h.Tmpl.ExecuteTemplate(w, "register-page5-register.html", nil)
	if err != nil {
		log.Println("Unable to execute register-page5-register.html", err)
		return
	}
}

func (h Handler) HandlerRegister(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerRegister")

	if r.Method == http.MethodPost {
		name := strings.TrimSpace(r.FormValue("name"))
		slug := helpers.MakeSlug(name)
		displayName := strings.TrimSpace(r.FormValue("display-name"))
		bandNameEntry := strings.TrimSpace(r.FormValue("band"))
		email := helpers.NormalizeEmail(r.FormValue("email"))
		emailConfirmation := helpers.NormalizeEmail(r.FormValue("email-confirmation"))

		if email != emailConfirmation {
			log.Println("email does not match confirmation")
			http.Error(w, "emails must match", http.StatusBadRequest)
			return
		}

		validEmail := helpers.ValidateEmail(email)
		validEmailConfirmation := helpers.ValidateEmail(emailConfirmation)

		if validEmail != true || validEmailConfirmation != true {
			log.Println("invalid email address provided")
			http.Error(w, "Invalid email address provided", http.StatusBadRequest)
			return
		}

		password := r.FormValue("password")
		passwordConfirmation := r.FormValue("password-confirmation")

		if password != passwordConfirmation {
			log.Println("password does not match confirmation")
			http.Error(w, "passwords must match", http.StatusBadRequest)
			return
		}
		if len(password) < 8 || len(passwordConfirmation) < 8 {
			log.Println(" len of password too short: ", len(password), len(passwordConfirmation))
			http.Error(w, "passwords must be at least 8 characters long", http.StatusBadRequest)
			return
		}
		isAdmin := true

		bandName := HelperProcessBandNameEntry(bandNameEntry)
		_, err := database.BandsTableGetBandByName(bandName)
		if err == nil {
			isAdmin = false
		}

		user, err := database.UsersTableCreateUser(name, displayName, slug, email, password, passwordConfirmation, isAdmin)
		if err != nil {
			log.Println("   register err: ", err)
			http.Error(w, "Could not create user", http.StatusInternalServerError)
			return
		}

		band, err := database.BandsTableGetBandByName(bandName)
		if err != nil {
			bandSlug := helpers.MakeSlug(bandName)
			band, err = database.BandsTableCreateBand(bandName, user.UserID, bandSlug)
			if err != nil {
				log.Println("   register err: ", err)
				http.Error(w, "Could not create band", http.StatusInternalServerError)
				return
			}

			chatName := fmt.Sprintf("%s (Band Chat)", band.Name)
			chatSlug := helpers.MakeSlug(chatName)
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

	email := helpers.NormalizeEmail(r.FormValue("email"))
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

	token, err := helpers.GenerateSessionToken()
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

	_, err = r.Cookie("session_token")
	if err != nil {
		slog.Error(
			"unable to get session token",
			"request_id", requestlog.GetRequestID(r.Context()),
			"error", err,
		)
	}

	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusOK)
	return
}

func (h Handler) HandlerUserAgreementPage(w http.ResponseWriter, r *http.Request) {

	h.Tmpl.ExecuteTemplate(w, "user-agreement.html", nil)
	return
}

func (h Handler) HandlerUserAgreement(w http.ResponseWriter, r *http.Request) {

	h.Tmpl.ExecuteTemplate(w, "login.html", nil)
	return
}

func (h Handler) HandlerTermsPage(w http.ResponseWriter, r *http.Request) {

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

	err := h.Tmpl.ExecuteTemplate(w, "privacy.html", nil)
	if err != nil {
		slog.Error(
			"unable to load privacy.html",
			"request_id", requestlog.GetRequestID(r.Context()),
			"error", err,
		)
		http.Error(w, "Unable to load privacy page", http.StatusInternalServerError)
		return
	}
}

func (h Handler) HandlerAccessCodePage(w http.ResponseWriter, r *http.Request) {

	err := h.Tmpl.ExecuteTemplate(w, "access.html", nil)
	if err != nil {
		slog.Error(
			"unable to load access.html",
			"request_id", requestlog.GetRequestID(r.Context()),
			"error", err,
		)
		http.Error(w, "Unable to load access page", http.StatusInternalServerError)
		return
	}
}

func (h Handler) HandlerCreateAccessCode(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerCreateAccessCode")

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
