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
	"unicode/utf8"

	"github.com/alexedwards/argon2id"
	"github.com/google/uuid"
)

func (h Handler) HandlerRegisterAccessCodePage(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerRegisterAccessCodePage")

	accessCode := strings.TrimSpace(r.FormValue("access-code"))

	if accessCode != "" {
		_, err := database.AccessCodesTableValidateCodeReturnBandID(r.Context(), accessCode)
		if err != nil {
			http.Error(w, "invalid access code", http.StatusBadRequest)
			return
		}
	}

	registrationID := strings.TrimSpace(r.FormValue("registration-id"))

	var user models.UserRegistration

	if registrationID != "" {
		registeredUser, err := database.UsersRegTableGetUserRegistrationID(registrationID)
		if err != nil {
			http.Error(w, "Invalid registration id", http.StatusBadRequest)
			return
		}
		user = registeredUser
	}

	data := models.RegistrationPages{
		AccessCode:     accessCode,
		RegistrationID: registrationID,
		User:           user,
	}

	err := h.Tmpl.ExecuteTemplate(w, "register-page1-access-code.html", data)
	if err != nil {
		slog.Error(
			"unable to load register-page1-access-code.html",
			"request_id", requestlog.GetRequestID(r.Context()),
			"error", err,
		)
		return
	}
}

func (h Handler) HandlerRegisterUserInfoPage(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerRegisterUserInfoPage")

	accessCode := strings.TrimSpace(r.FormValue("access-code"))
	registrationID := strings.TrimSpace(r.FormValue("registration-id"))

	bandID := ""

	if accessCode != "" {
		existingBandID, err := database.AccessCodesTableValidateCodeReturnBandID(r.Context(), accessCode)
		if err != nil {
			slog.Error(
				"unable to validate access code",
				"request_id", requestlog.GetRequestID(r.Context()),
				"path", r.URL.Path,
				"error", err,
			)
			http.Error(w, "Invalid access code", http.StatusBadRequest)
			return
		}
		bandID = existingBandID
	}

	newUser := models.UserRegistration{
		AccessCodeHash: accessCode,
	}

	if registrationID != "" {
		registeredUser, err := database.UsersRegTableGetUserRegistrationID(registrationID)
		if err != nil {
			log.Println("Unable to get registered user", err)
			http.Error(w, "Unable to load registration", http.StatusInternalServerError)
			return
		}
		newUser = registeredUser
	}

	newUser.BandID = bandID

	data := models.RegistrationPages{
		User:           newUser,
		RegistrationID: registrationID,
	}

	err := h.Tmpl.ExecuteTemplate(w, "register-page2-name.html", data)
	if err != nil {
		log.Println("Unable to execute register-page1-access-code.html", err)
		return
	}
}

func (h Handler) HandlerRegisterUserInfoSubmit(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerRegisterUserInfoSubmit")

	email := strings.TrimSpace(r.FormValue("email"))
	emailConfirmation := strings.TrimSpace(r.FormValue("email-confirmation"))

	if email == "" || len(email) > 254 {
		http.Error(w, "Invalid email entry", http.StatusBadRequest)
		return
	}

	if email != emailConfirmation {
		http.Error(w, "Email addresses do not match", http.StatusBadRequest)
		return
	}

	normalizedEmail := helpers.NormalizeEmail(email)
	validatedEmail := helpers.ValidateEmail(normalizedEmail)

	if !validatedEmail {
		http.Error(w, "Invalid email entry", http.StatusBadRequest)
		return
	}

	firstName := strings.TrimSpace(r.FormValue("first-name"))

	if firstName == "" || utf8.RuneCountInString(firstName) > 100 {
		http.Error(w, "Invalid entry for first name", http.StatusBadRequest)
		return
	}

	lastName := strings.TrimSpace(r.FormValue("last-name"))

	if utf8.RuneCountInString(lastName) > 100 {
		http.Error(w, "Invalid entry for last name", http.StatusBadRequest)
		return
	}

	displayName := strings.TrimSpace(r.FormValue("display-name"))

	if displayName == "" || utf8.RuneCountInString(displayName) > 100 {
		http.Error(w, "Invalid entry for display name", http.StatusBadRequest)
		return
	}

	timezone := r.FormValue("timezone")
	accessCode := strings.TrimSpace(r.FormValue("access-code"))
	registrationID := strings.TrimSpace(r.FormValue("registration-id"))

	bandID := ""

	if accessCode != "" {
		existingBandID, err := database.AccessCodesTableValidateCodeReturnBandID(r.Context(), accessCode)
		if err != nil {
			log.Println("Unable to get bandID by access code: ", err)
			http.Error(w, "Invalid access code", http.StatusBadRequest)
			return
		}
		bandID = existingBandID
	}

	newUser := models.UserRegistration{
		FirstName:          firstName,
		LastName:           lastName,
		DisplayName:        displayName,
		Email:              normalizedEmail,
		Timezone:           timezone,
		BandID:             bandID,
		UserRegistrationID: registrationID,
		AccessCodeHash:     accessCode,
	}

	newUser, err := h.Services.RegistrationSaveUserProfile(newUser)
	if err != nil {
		slog.Error(
			"unable to create user registration profile",
			"request_id", requestlog.GetRequestID(r.Context()),
			"path", r.URL.Path,
			"first_name", firstName,
			"last_name", email,
			"error", err,
		)
		http.Error(w, "unable to create user registration profile", http.StatusInternalServerError)
		return
	}

	var band models.Band

	if newUser.BandID != "" {
		band, err = database.BandsTableGetBandByBandID(newUser.BandID)
		if err != nil {
			http.Error(w, "Band not found", http.StatusBadRequest)
			return
		}
	}

	data := models.RegistrationPages{
		User: newUser,
		Band: band,
	}

	err = h.Tmpl.ExecuteTemplate(w, "register-page3-band.html", data)
	if err != nil {
		slog.Error(
			"unable to load register-page3-band.html",
			"request_id", requestlog.GetRequestID(r.Context()),
			"error", err,
		)
		return
	}
}

func (h Handler) HandlerRegisterBandPageSubmit(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerRegisterBandPageSubmit")

	registrationID := strings.TrimSpace(r.FormValue("registration-id"))
	if _, err := uuid.Parse(registrationID); err != nil {
		http.Error(w, "Invalid registration ID", http.StatusBadRequest)
		return
	}

	bandName := strings.TrimSpace(r.FormValue("band-name"))
	if bandName == "" || utf8.RuneCountInString(bandName) > 100 {
		http.Error(w, "Invalid entry for band name", http.StatusBadRequest)
		return
	}

	data := models.RegistrationPages{
		RegistrationID: registrationID,
		BandName:       bandName,
	}

	err := h.Tmpl.ExecuteTemplate(w, "register-page4-password.html", data)
	if err != nil {
		slog.Error(
			"unable to load register-page4-password.html",
			"request_id", requestlog.GetRequestID(r.Context()),
			"path", r.URL.Path,
			"error", err,
		)
		return
	}
}

func (h Handler) HandlerRegisterPassword(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerRegisterPassword")

	accepted := r.FormValue("legal-agreement") == "accepted"

	if accepted != true {
		http.Error(w, "Terms of Service and Privacy Policy not accepted", http.StatusBadRequest)
		return
	}

	password := r.FormValue("password")
	passwordConfirmation := r.FormValue("password-confirmation")

	if len(password) < 8 || len(password) > 255 {
		http.Error(w, "Invalid password", http.StatusSeeOther)
		return
	}

	if password != passwordConfirmation {
		http.Error(w, "passwords do not match", http.StatusBadRequest)
		return
	}

	registrationID := strings.TrimSpace(r.FormValue("registration-id"))
	if _, err := uuid.Parse(registrationID); err != nil {
		http.Error(w, "Invalid registration ID", http.StatusBadRequest)
		return
	}

	valid, err := database.UserRegTableValidateRegistrationID(registrationID)
	if err != nil || valid != true {
		http.Error(w, "invalid registration id", http.StatusBadRequest)
		return
	}

	bandName := strings.TrimSpace(r.FormValue("band-name"))
	if bandName == "" {
		http.Error(w, "No band name provided", http.StatusBadRequest)
		return
	}

	if utf8.RuneCountInString(bandName) > 100 {
		http.Error(w, "Band names must be 100 characters or fewer", http.StatusBadRequest)
		return
	}

	user, err := h.Services.RegistrationSavePassword(registrationID, password)
	if err != nil {
		slog.Error(
			"unable to save password",
			"request_id", requestlog.GetRequestID(r.Context()),
			"path", r.URL.Path,
			"error", err,
		)
		http.Error(w, "Unable to save password", http.StatusInternalServerError)
		return
	}

	if registrationID != "" {
		user.UserRegistrationID = registrationID
	}

	fullName := user.FirstName + " " + user.LastName

	newUser := models.User{
		UserID:       uuid.NewString(),
		Name:         fullName,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		DisplayName:  user.DisplayName,
		Email:        user.Email,
		Slug:         helpers.MakeSlug(user.FirstName),
		PasswordHash: user.PasswordHash,
		IsAdmin:      false,
		TimeZone:     user.Timezone,
	}

	band := models.Band{}

	bandID, err := database.AccessCodesTableValidateCodeReturnBandID(r.Context(), user.AccessCodeHash)
	if err != nil {
		log.Println("Error: ", err)
	}
	if bandID != "" {
		band, err = database.BandsTableGetBandByBandID(bandID)
		if err != nil {
			slog.Error(
				"unable to load band",
				"request_id", requestlog.GetRequestID(r.Context()),
				"path", r.URL.Path,
				"error", err,
			)
			http.Error(w, "Unable to load band", http.StatusSeeOther)
			return
		}

		chatID, err := database.ChatsTableGetPrimaryChatIDByBandID(bandID)
		if err != nil {
			log.Println("Unable to get primary band chat: ", err)
			http.Error(w, "Unable to get primary band chat", http.StatusSeeOther)
			return
		}

		err = database.RegisterNewBandUser(newUser, bandID, chatID, user.AccessCodeHash)
		if err != nil {
			slog.Error(
				"unable to save user",
				"request_id", requestlog.GetRequestID(r.Context()),
				"path", r.URL.Path,
				"first_name", user.FirstName,
				"last_name", user.LastName,
				"email", user.Email,
				"error", err,
			)
			http.Error(w, "Unable to save user", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return

	} else {

		data := models.RegistrationPages{
			User: user,
			Band: band,
		}

		newBand := models.Band{
			BandID: uuid.NewString(),
			Name:   bandName,
			Slug:   helpers.MakeSlug(bandName),
		}

		newUser.IsAdmin = true

		err = database.RegisterInitialUserBandAndChat(newUser, newBand, registrationID)
		if err != nil {
			slog.Error(
				"unable to register new user",
				"request_id", requestlog.GetRequestID(r.Context()),
				"path", r.URL.Path,
				"first_name", newUser.FirstName,
				"last_name", newUser.LastName,
				"email", newUser.Email,
				"error", err,
			)
			http.Error(w, "Unable to register user", http.StatusInternalServerError)
			return
		}

		err = h.Tmpl.ExecuteTemplate(w, "login.html", data)
		if err != nil {
			slog.Error(
				"unable to load login.html",
				"request_id", requestlog.GetRequestID(r.Context()),
				"path", r.URL.Path,
				"error", err,
			)
			http.Error(w, "Error getting login page", http.StatusInternalServerError)
			return
		}
	}
}

func (h Handler) HandlerLoginPage(w http.ResponseWriter, r *http.Request) {

	user, err := HelperGetAuthenticatedUser(r)
	if err == nil {
		log.Println("   Already logged in: ", user.Name, err)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	h.Tmpl.ExecuteTemplate(w, "login.html", nil)
	return
}

func (h Handler) HandlerLogin(w http.ResponseWriter, r *http.Request) {

	email := strings.TrimSpace(r.FormValue("email"))

	if email == "" || len(email) > 254 {
		http.Error(w, "Invalid email entry", http.StatusBadRequest)
		return
	}

	normalizedEmail := helpers.NormalizeEmail(email)
	validatedEmail := helpers.ValidateEmail(normalizedEmail)

	if !validatedEmail {
		http.Error(w, "Invalid email entry", http.StatusBadRequest)
		return
	}

	password := r.FormValue("password")

	if len(password) < 8 || len(password) > 255 {
		http.Error(w, "Invalid password", http.StatusSeeOther)
		return
	}

	user, err := database.UsersTableGetUserByEmail(normalizedEmail)
	if err != nil {
		log.Println("   HandlerLogin: Unable to get user: ", err)
		w.Write([]byte("Invalid email or password"))
		return
	}

	match, err := argon2id.ComparePasswordAndHash(
		password,
		user.PasswordHash,
	)
	if err != nil {
		w.Write([]byte("* Invalid email or password * "))
		return
	}

	if !match {
		w.Write([]byte("* Invalid email or password * "))
		return
	}

	token, err := helpers.GenerateSessionToken()
	if err != nil {
		slog.Error(
			"unable to generate session token",
			"request_id", requestlog.GetRequestID(r.Context()),
			"error", err,
		)

		http.Error(w, "Unable to log in", http.StatusInternalServerError)
		return
	}

	params := models.CreateSessionParams{
		UserID: user.UserID,
		Token:  token,
	}

	session, err := database.SessionsTableCreateSession(params)
	if err != nil {
		log.Println("   unable to create session: ", err)
		w.Write([]byte("* Unable to login * "))
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    session.Token,
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	// _, err = r.Cookie("session_token")
	// if err != nil {
	// 	slog.Error(
	// 		"unable to get session token",
	// 		"request_id", requestlog.GetRequestID(r.Context()),
	// 		"error", err,
	// 	)
	// }

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

	if user.IsAdmin != true {
		http.Error(w, "User must be admin to generate access code", http.StatusForbidden)
		return
	}
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
