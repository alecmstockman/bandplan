package handlers

import (
	"bandplan/src/database"
	"bandplan/src/helpers"
	requestlog "bandplan/src/logging"
	"bandplan/src/models"
	"bandplan/src/services"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/google/uuid"
)

const (
	currentTermsVersion   = "2026-07-14"
	currentPrivacyVersion = "2026-07-14"
)

func (h Handler) HandlerRegisterAccessCodePage(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerRegisterAccessCodePage")
	w.Header().Set("Cache-Control", "no-store")

	r.Body = http.MaxBytesReader(w, r.Body, 16*1024)

	if err := r.ParseForm(); err != nil {
		var sizeErr *http.MaxBytesError
		if errors.As(err, &sizeErr) {
			http.Error(w, "Request too large", http.StatusRequestEntityTooLarge)
		} else {
			http.Error(w, "Invalid form", http.StatusBadRequest)
		}
		return
	}

	registrationToken := strings.TrimSpace(r.FormValue("registration-id"))
	accessCode := helpers.NormalizeAccessCode(r.FormValue("access-code"))

	data, err := h.Services.RegistrationLoadAccessData(r.Context(), registrationToken, accessCode)
	if err != nil {
		slog.Error(
			"unable to load auth context",
			"request_id", requestlog.GetRequestID(r.Context()),
			"error", err,
		)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	err = h.Tmpl.ExecuteTemplate(w, "register-page1-access-code.html", data)
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
	w.Header().Set("Cache-Control", "no-store")

	r.Body = http.MaxBytesReader(w, r.Body, 16*1024)

	if err := r.ParseForm(); err != nil {
		var sizeErr *http.MaxBytesError
		if errors.As(err, &sizeErr) {
			http.Error(w, "Request too large", http.StatusRequestEntityTooLarge)
		} else {
			http.Error(w, "Invalid form", http.StatusBadRequest)
		}
		return
	}

	accessCode := helpers.NormalizeAccessCode(r.FormValue("access-code"))
	if r.FormValue("skip-access-code") == "true" {
		accessCode = ""
	}
	registrationToken := strings.TrimSpace(r.FormValue("registration-id"))

	bandID := ""
	accessCodeHash := ""

	if accessCode != "" {
		accessCodeHash = helpers.HashRegistrationCode(accessCode)

		existingBandID, err := database.AccessCodesTableValidateCodeReturnBandID(r.Context(), accessCodeHash)
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
		AccessCodeHash: accessCodeHash,
		BandID:         bandID,
	}

	if registrationToken != "" {
		valid := helpers.ValidateTokenLength(registrationToken)
		if valid != true {
			http.Error(w, "Unable to load registration", http.StatusNotFound)
			return
		}

		registrationTokenHash := helpers.HashRegistrationCode(registrationToken)

		registeredUser, err := database.UsersRegTableGetRegistrationToken(registrationTokenHash)
		if err != nil {
			log.Println("Unable to get registered user", err)
			http.Error(w, "Unable to load registration", http.StatusNotFound)
			return
		}
		newUser = registeredUser
	}

	newUser.RegistrationToken = registrationToken

	data := models.RegistrationPages{
		User:              newUser,
		AccessCode:        accessCode,
		RegistrationToken: registrationToken,
	}

	err := h.Tmpl.ExecuteTemplate(w, "register-page2-name.html", data)
	if err != nil {
		slog.Error(
			"unable to validate access code",
			"request_id", requestlog.GetRequestID(r.Context()),
			"path", r.URL.Path,
			"error", err,
		)
		http.Error(w, "Unable to load page", http.StatusInternalServerError)
		return
	}
}

func (h Handler) HandlerRegisterUserInfoSubmit(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerRegisterUserInfoSubmit")
	w.Header().Set("Cache-Control", "no-store")

	r.Body = http.MaxBytesReader(w, r.Body, 16*1024)

	if err := r.ParseForm(); err != nil {
		var sizeErr *http.MaxBytesError
		if errors.As(err, &sizeErr) {
			http.Error(w, "Request too large", http.StatusRequestEntityTooLarge)
		} else {
			http.Error(w, "Invalid form", http.StatusBadRequest)
		}
		return
	}

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

	valid := helpers.ValidateNameEntryLength(firstName)
	if valid != true {
		http.Error(w, "Invalid entry for first name", http.StatusBadRequest)
		return
	}

	lastName := strings.TrimSpace(r.FormValue("last-name"))

	valid = helpers.ValidateNameEntryMaxLength(lastName)
	if valid != true {
		http.Error(w, "Invalid entry for last name", http.StatusBadRequest)
		return
	}

	displayName := strings.TrimSpace(r.FormValue("display-name"))

	valid = helpers.ValidateNameEntryMaxLength(displayName)
	if valid != true {
		http.Error(w, "Invalid entry for display name", http.StatusBadRequest)
		return
	}
	if displayName == "" {
		displayName = firstName
	}

	timezone := r.FormValue("timezone")
	accessCode := helpers.NormalizeAccessCode(r.FormValue("access-code"))
	registrationToken := strings.TrimSpace(r.FormValue("registration-id"))

	bandID := ""
	accessCodeHash := ""

	if accessCode != "" {
		accessCodeHash = helpers.HashRegistrationCode(accessCode)

		existingBandID, err := database.AccessCodesTableValidateCodeReturnBandID(r.Context(), accessCodeHash)
		if err != nil {
			log.Println("Unable to get bandID by access code: ", err)
			http.Error(w, "Invalid access code", http.StatusBadRequest)
			return
		}
		bandID = existingBandID
	}

	newUser := models.UserRegistration{
		FirstName:         firstName,
		LastName:          lastName,
		DisplayName:       displayName,
		Email:             normalizedEmail,
		Timezone:          timezone,
		BandID:            bandID,
		RegistrationToken: registrationToken,
		AccessCodeHash:    accessCodeHash,
	}

	newUser, err := h.Services.RegistrationSaveUserProfile(newUser)
	if err != nil {
		slog.Error(
			"unable to create user registration profile",
			"request_id", requestlog.GetRequestID(r.Context()),
			"path", r.URL.Path,
			"error", err,
		)
		http.Error(w, "unable to create user registration profile", http.StatusInternalServerError)
		return
	}

	var band models.Band

	if newUser.BandID != "" {
		band, err = database.BandsTableGetBandByBandID(newUser.BandID)
		if err != nil {
			http.Error(w, "Band not found", http.StatusNotFound)
			return
		}
	}

	data := models.RegistrationPages{
		User:              newUser,
		Band:              band,
		RegistrationToken: newUser.RegistrationToken,
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

func (h Handler) HandlerRegisterBandPage(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerRegisterBandPage")
	w.Header().Set("Cache-Control", "no-store")

	r.Body = http.MaxBytesReader(w, r.Body, 16*1024)

	if err := r.ParseForm(); err != nil {
		var sizeErr *http.MaxBytesError
		if errors.As(err, &sizeErr) {
			http.Error(w, "Request too large", http.StatusRequestEntityTooLarge)
		} else {
			http.Error(w, "Invalid form", http.StatusBadRequest)
		}
		return
	}

	registrationToken := strings.TrimSpace(r.FormValue("registration-id"))
	if !helpers.ValidateTokenLength(registrationToken) {
		http.Error(w, "Invalid registration token", http.StatusBadRequest)
		return
	}

	registrationTokenHash := helpers.HashRegistrationCode(registrationToken)

	user, err := database.UsersRegTableGetRegistrationToken(registrationTokenHash)
	if err != nil {
		slog.Error(
			"unable to load user registration",
			"request_id", requestlog.GetRequestID(r.Context()),
			"path", r.URL.Path,
			"error", err,
		)
		http.Error(w, "Unable to load registration", http.StatusBadRequest)
		return
	}

	band := models.Band{
		Name: strings.TrimSpace(r.FormValue("band-name")),
	}
	if band.Name == "" && user.BandID != "" {
		band, err = database.BandsTableGetBandByBandID(user.BandID)
		if err != nil {
			slog.Error(
				"unable to load registration band",
				"request_id", requestlog.GetRequestID(r.Context()),
				"path", r.URL.Path,
				"error", err,
			)
			http.Error(w, "Unable to load band", http.StatusBadRequest)
			return
		}
	}

	data := models.RegistrationPages{
		User:              user,
		Band:              band,
		RegistrationToken: registrationToken,
	}

	if err := h.Tmpl.ExecuteTemplate(w, "register-page3-band.html", data); err != nil {
		slog.Error(
			"unable to load register-page3-band.html",
			"request_id", requestlog.GetRequestID(r.Context()),
			"path", r.URL.Path,
			"error", err,
		)
		http.Error(w, "Unable to load page", http.StatusInternalServerError)
	}
}

func (h Handler) HandlerRegisterBandPageSubmit(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerRegisterBandPageSubmit")
	w.Header().Set("Cache-Control", "no-store")

	r.Body = http.MaxBytesReader(w, r.Body, 16*1024)

	if err := r.ParseForm(); err != nil {
		var sizeErr *http.MaxBytesError
		if errors.As(err, &sizeErr) {
			http.Error(w, "Request too large", http.StatusRequestEntityTooLarge)
		} else {
			http.Error(w, "Invalid form", http.StatusBadRequest)
		}
		return
	}

	registrationToken := strings.TrimSpace(r.FormValue("registration-id"))
	bandName := strings.TrimSpace(r.FormValue("band-name"))

	data, err := h.Services.RegistrationBandPageSubmit(r.Context(), registrationToken, bandName)
	if err != nil {
		slog.Error(
			"unable to validate access code",
			"request_id", requestlog.GetRequestID(r.Context()),
			"path", r.URL.Path,
			"error", err,
		)
		http.Error(w, "Invalid entry", http.StatusInternalServerError)
		return
	}

	err = h.Tmpl.ExecuteTemplate(w, "register-page4-password.html", data)
	if err != nil {
		slog.Error(
			"unable to load register-page4-password.html",
			"request_id", requestlog.GetRequestID(r.Context()),
			"path", r.URL.Path,
			"error", err,
		)
		http.Error(w, "Unable to load page", http.StatusInternalServerError)
		return
	}
}

func (h Handler) HandlerRegisterPassword(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerRegisterPassword")
	w.Header().Set("Cache-Control", "no-store")

	r.Body = http.MaxBytesReader(w, r.Body, 16*1024)

	if err := r.ParseForm(); err != nil {
		var sizeErr *http.MaxBytesError
		if errors.As(err, &sizeErr) {
			http.Error(w, "Request too large", http.StatusRequestEntityTooLarge)
		} else {
			http.Error(w, "Invalid form", http.StatusBadRequest)
		}
		return
	}

	accepted := r.FormValue("legal-agreement") == "accepted"

	if accepted != true {
		http.Error(w, "Terms of Service and Privacy Policy not accepted", http.StatusBadRequest)
		return
	}

	password := r.FormValue("password")
	passwordConfirmation := r.FormValue("password-confirmation")

	if password != passwordConfirmation {
		http.Error(w, "passwords do not match", http.StatusBadRequest)
		return
	}

	valid := helpers.PasswordValidateLength(password)
	if valid != true {
		http.Error(w, "Invalid password length", http.StatusBadRequest)
		return
	}

	registrationToken := strings.TrimSpace(r.FormValue("registration-id"))
	if !helpers.ValidateTokenLength(registrationToken) {
		http.Error(w, "invalid registration token", http.StatusBadRequest)
		return
	}

	registrationTokenHash := helpers.HashRegistrationCode(registrationToken)

	valid, err := database.UserRegTableValidateRegistrationToken(registrationTokenHash)
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

	user, err := h.Services.RegistrationSavePassword(registrationToken, password)
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

	fullName := user.FirstName + " " + user.LastName

	newUser := models.User{
		UserID:         uuid.NewString(),
		Name:           fullName,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		DisplayName:    user.DisplayName,
		Email:          user.Email,
		Slug:           helpers.MakeSlug(user.FirstName),
		PasswordHash:   user.PasswordHash,
		IsAdmin:        false,
		LegalAccepted:  accepted,
		TermsVersion:   currentTermsVersion,
		PrivacyVersion: currentPrivacyVersion,
		TimeZone:       user.Timezone,
	}

	// band := models.Band{}

	var bandID string

	if user.AccessCodeHash != "" {
		validatedBandID, err := database.AccessCodesTableValidateCodeReturnBandID(r.Context(), user.AccessCodeHash)
		if err != nil {
			slog.Error(
				"unable to validate registration code",
				"request_id", requestlog.GetRequestID(r.Context()),
				"path", r.URL.Path,
				"error", err,
			)
			http.Error(w, "Invalid registration code", http.StatusBadRequest)
			return
		}
		bandID = validatedBandID
	}

	if bandID != "" {
		_, err = database.BandsTableGetBandByBandID(bandID)
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

		_, err = database.RegisterNewBandUser(r.Context(), newUser, bandID, chatID, user.RegistrationTokenHash, user.AccessCodeHash)
		if err != nil {
			slog.Error(
				"unable to save user",
				"request_id", requestlog.GetRequestID(r.Context()),
				"path", r.URL.Path,
				"error", err,
			)
			http.Error(w, "Unable to save user", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return

	} else {

		// data := models.RegistrationPages{
		// 	User: user,
		// 	Band: band,
		// }

		newBand := models.Band{
			BandID: uuid.NewString(),
			Name:   bandName,
			Slug:   helpers.MakeSlug(bandName),
		}

		newUser.IsAdmin = true

		err = database.RegisterInitialUserBandAndChat(r.Context(), newUser, newBand, user.RegistrationTokenHash)
		if err != nil {
			slog.Error(
				"unable to register new user",
				"request_id", requestlog.GetRequestID(r.Context()),
				"path", r.URL.Path,
				"error", err,
			)
			http.Error(w, "Unable to register user", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
}

func (h Handler) HandlerLoginPage(w http.ResponseWriter, r *http.Request) {

	_, err := HelperGetAuthenticatedUser(r)
	if err == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	h.Tmpl.ExecuteTemplate(w, "login.html", nil)
	return
}

func (h Handler) HandlerLogin(w http.ResponseWriter, r *http.Request) {

	r.Body = http.MaxBytesReader(w, r.Body, 16*1024)

	if err := r.ParseForm(); err != nil {
		var sizeErr *http.MaxBytesError
		if errors.As(err, &sizeErr) {
			http.Error(w, "Request too large", http.StatusRequestEntityTooLarge)
		} else {
			http.Error(w, "Invalid form", http.StatusBadRequest)
		}
		return
	}

	email := strings.TrimSpace(r.FormValue("email"))
	password := r.FormValue("password")

	session, err := h.Services.LoginValidation(r.Context(), email, password)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			w.Write([]byte("* Invalid email or password * "))
			return
		}
		slog.Error(
			"user unable to log in",
			"request_id", requestlog.GetRequestID(r.Context()),
			"path", r.URL.Path,
			"error", err,
		)
		http.Error(w, "unable to log in, please try again", http.StatusInternalServerError)
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

	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusOK)
	return
}

func (h Handler) HandlerLogout(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\n\n- HandlerLogout")

	cookie, err := r.Cookie("session_token")
	if err != nil {
		fmt.Println("Unable to get cookie")
	}

	token := cookie.Value

	err = database.SessionsTableDeleteSessionByToken(token)
	if err != nil {
		slog.Error(
			"unable to delete session token by user id or token",
			"request_id", requestlog.GetRequestID(r.Context()),
			"path", r.URL.Path,
			"error", err,
		)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		MaxAge:   -1,
	})

	w.Header().Set("HX-Redirect", "/login")
	w.WriteHeader(http.StatusOK)
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
	band := auth.CurrentBand

	html, err := h.Services.RegistrationCreateAccessCode(r.Context(), user, band)
	if err != nil {
		slog.Error(
			"unable to load register-page1-access-code.html",
			"request_id", requestlog.GetRequestID(r.Context()),
			"path", r.URL.Path,
			"error", err,
		)
		http.Error(w, "Unable to generate access code", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(html))
	return
}
