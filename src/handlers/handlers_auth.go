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

	"github.com/alexedwards/argon2id"
	"github.com/google/uuid"
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

	bandID, err := database.AccessCodesTableValidateCodeReturnBandID(code)
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

func (h Handler) HandlerRegisterAccessCodePage(w http.ResponseWriter, r *http.Request) {
	log.Println("\n\n- HandlerRegisterAccessCodePage")

	accessCode := strings.TrimSpace(r.FormValue("access-code"))

	err := h.Tmpl.ExecuteTemplate(w, "register-page1-access-code.html", accessCode)
	if err != nil {
		log.Println("Unable to execute register-page1-access-code.html")
		return
	}
}

func (h Handler) HandlerRegisterUserInfoPage(w http.ResponseWriter, r *http.Request) {
	log.Println("\n\n- HandlerRegisterUserInfoPage")

	accessCode := r.FormValue("access-code")
	registrationID := r.FormValue("registration-id")

	bandID := ""

	if accessCode != "" {
		existingBandID, err := database.AccessCodesTableValidateCodeReturnBandID(accessCode)
		if err != nil {
			fmt.Println("Unable to get bandID by access code")
			http.Error(w, "Invalid access code", http.StatusBadRequest)
			return
		}
		bandID = existingBandID
	}

	newUser := models.UserRegistration{
		AccessCodeHash: accessCode,
	}

	if registrationID != "" {
		registeredUser, err := database.UsersRegTableGetUserByID(registrationID)
		if err != nil {
			log.Println("Unable to get registered user", err)
			http.Error(w, "Unable to load registration", http.StatusInternalServerError)
			return
		}
		newUser = registeredUser
	}

	newUser.BandID = bandID

	err := h.Tmpl.ExecuteTemplate(w, "register-page2-name.html", newUser)
	if err != nil {
		log.Println("Unable to execute register-page1-access-code.html", err)
		return
	}
}

func (h Handler) HandlerRegisterUserInfoSubmit(w http.ResponseWriter, r *http.Request) {
	log.Println("\n\n- HandlerRegisterUserInfoSubmit")

	email := strings.TrimSpace(r.FormValue("email"))
	emailConfirmation := strings.TrimSpace(r.FormValue("email-confirmation"))

	fmt.Println("email: ", email)
	fmt.Println("cmail: ", emailConfirmation)

	if email != emailConfirmation {
		log.Println("email addresses do not match")
		http.Error(w, "Email addresses do not match", http.StatusBadRequest)
		return
	}

	normalizedEmail := helpers.NormalizeEmail(email)
	validatedEmail := helpers.ValidateEmail(normalizedEmail)

	if !validatedEmail {
		log.Println("invalid email provided")
		http.Error(w, "Invalid Email addresses", http.StatusBadRequest)
		return
	}

	firstName := strings.TrimSpace(r.FormValue("first-name"))
	lastName := strings.TrimSpace(r.FormValue("last-name"))
	displayName := strings.TrimSpace(r.FormValue("display-name"))

	timezone := r.FormValue("timezone")
	accessCode := r.FormValue("access-code")
	registrationID := r.FormValue("registration-id")

	// fmt.Println("first name:     ", firstName)
	// fmt.Println("last name:      ", lastName)
	// fmt.Println("display name:   ", displayName)
	// fmt.Println("timezone:       ", timezone)
	// fmt.Println("access code:    ", accessCode)
	// fmt.Println("registrationID: ", registrationID)

	bandID := ""

	if accessCode != "" {
		existingBandID, err := database.AccessCodesTableValidateCodeReturnBandID(accessCode)
		fmt.Println("existingBandID: ", existingBandID)
		if err != nil {
			fmt.Println("Unable to get bandID by access code")
			http.Error(w, "Invalid access code", http.StatusBadRequest)
			return
		}
		bandID = existingBandID
		fmt.Println("bandID: ", bandID)
	}

	fmt.Println("bandID: ", bandID)

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
		log.Println("   unable to create user registration profile", err)
		http.Error(w, "unable to create user registration profile", http.StatusInternalServerError)
		return
	}

	// fmt.Printf("%+v\n", newUser)
	// fmt.Println("test")

	var band models.Band

	if newUser.BandID != "" {
		band, err = database.BandsTableGetBandByBandID(newUser.BandID)
		// fmt.Println("BandID: ", band.ID)
		// fmt.Println("BandName: ", band.Name)
		if err != nil {
			http.Error(w, "Band not found", http.StatusBadRequest)
			return
		}
	}

	// fmt.Println("BandID: ", band.ID)
	// fmt.Println("BandName: ", band.Name)

	data := models.RegistrationPages{
		User: newUser,
		Band: band,
	}

	// http.Redirect(w, r, "/register/3", http.StatusSeeOther)

	err = h.Tmpl.ExecuteTemplate(w, "register-page3-band.html", data)
	if err != nil {
		log.Println("Unable to execute register-page3-band.html", err)
		return
	}
}

// func (h Handler) HandlerRegisterBandPage(w http.ResponseWriter, r *http.Request) {
// 	log.Println("\n\n- HandlerRegisterBandPage")

// 	registrationID := r.FormValue("registration-id")
// 	fmt.Println("registration-id: ", registrationID)

// 	bandName := r.FormValue("band-name")

// 	newUser, err := database.UsersRegTableGetUserByID(registrationID)
// 	if err != nil {
// 		log.Println("   unable to get user registration profile", err)
// 		http.Error(w, "unable to get user registration profile", http.StatusInternalServerError)
// 		return
// 	}

// 	fmt.Println("newUser Access Code: ", newUser.AccessCodeHash)

// 	data := models.RegistrationPages{
// 		User:     newUser,
// 		BandName: bandName,
// 	}

// 	err = h.Tmpl.ExecuteTemplate(w, "register-page3-band.html", data)
// 	if err != nil {
// 		log.Println("Unable to execute register-page3-band.html", err)
// 		return
// 	}
// }

func (h Handler) HandlerRegisterBandPageSubmit(w http.ResponseWriter, r *http.Request) {
	log.Println("\n\n- HandlerRegisterBandPageSubmit")

	registrationID := r.FormValue("registration-id")
	bandName := r.FormValue("band-name")

	fmt.Println("registratinoID: ", registrationID)
	fmt.Println("bandName: ", bandName)

	data := models.RegistrationPages{
		RegistrationID: registrationID,
		BandName:       bandName,
	}

	err := h.Tmpl.ExecuteTemplate(w, "register-page4-password.html", data)
	if err != nil {
		log.Println("Unable to execute register-page3-password.html", err)
		return
	}
}

func (h Handler) HandlerRegisterPageFour(w http.ResponseWriter, r *http.Request) {
	log.Println("\n\n- HandlerRegisterPageFour")

	password := r.FormValue("password")
	passwordConfirmation := r.FormValue("password-confirmation")

	if password != passwordConfirmation {
		log.Println("passwords do not match")
		http.Error(w, "passwords do not match", http.StatusBadRequest)
		return
	}

	registrationID := r.FormValue("registration-id")
	accessCode := r.FormValue("access-code-hash")

	// fmt.Println("password: ", password)
	// fmt.Println("registration-id: ", registrationID)
	// fmt.Println("access-code-hash: ", accessCode)

	bandName := r.FormValue("band-name")
	if bandName == "" {
		log.Println("no band name provided")
		http.Error(w, "no band name provided", http.StatusBadRequest)
		return
	}
	// fmt.Println("band-name")

	user, err := h.Services.RegistrationSavePassword(registrationID, password)
	if err != nil {
		log.Println("   Unable to save password", err)
		http.Error(w, "Unable to save password", http.StatusInternalServerError)
		return
	}

	// fmt.Println("password saved")

	if registrationID != "" {
		user.UserRegistrationID = registrationID
	}

	band := models.Band{}

	bandID, err := database.AccessCodesTableValidateCodeReturnBandID(accessCode)
	if bandID == "" {
		log.Println("   Unable to validate code: ", err)
	} else {
		band, err = database.BandsTableGetBandByBandID(bandID)
		if err != nil {
			log.Println("   Unable to get band by band ID: ", err)
		}
	}

	data := models.RegistrationPages{
		User: user,
		Band: band,
	}

	newUser := models.User{
		UserID:       uuid.NewString(),
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		DisplayName:  user.DisplayName,
		Email:        user.Email,
		Slug:         helpers.MakeSlug(user.LastName),
		PasswordHash: user.PasswordHash,
		IsAdmin:      true,
		TimeZone:     user.Timezone,
	}

	newBand := models.Band{
		BandID: uuid.NewString(),
		Name:   bandName,
		Slug:   helpers.MakeSlug(bandName),
	}

	// fmt.Println("\n\nbandID: ", newBand.BandID)
	// fmt.Printf("\n\nUSER: %+v", newUser)
	// fmt.Printf("\n\nBAND: %+v", newBand)

	err = database.RegisterNewUserBandAndChat(newUser, newBand)
	if err != nil {
		fmt.Println("unable to register new user: ", err)
		http.Error(w, "Unable to register user", http.StatusInternalServerError)
		return
	}

	fmt.Println("End of registration: ")

	err = h.Tmpl.ExecuteTemplate(w, "login.html", data)
	if err != nil {
		log.Println("Unable to execute register-page4-band.html", err)
		return
	}
}

// func (h Handler) HandlerRegisterPageFourSubmit(w http.ResponseWriter, r *http.Request) {
// 	log.Println("\n\n- HandlerRegisterPageFourSubmit")

// 	password := r.FormValue("password")
// 	passwordConfirmation := r.FormValue("password-confirmation")

// 	if password != passwordConfirmation {
// 		log.Println("passwords do not match")
// 		http.Error(w, "passwords do not match", http.StatusBadRequest)
// 		return
// 	}

// 	passwordHash, err := helpers.HashPassword(password)
// 	if err != nil {
// 		log.Println("error hashing password: ", passwordHash)
// 		http.Error(w, "Error hashing password", http.StatusInternalServerError)
// 		return
// 	}

// 	registrationID := r.FormValue("registration-id")
// 	// fmt.Println("registration-id: ", registrationID)
// 	if registrationID == "" {
// 		http.Error(w, "Unable to get registration id", http.StatusNotFound)
// 		return
// 	}

// 	accessCode := r.FormValue("access-code-hash")
// 	// fmt.Println("access-code-hash: ", accessCode)

// 	bandName := r.FormValue("band-name")
// 	if bandName == "" {
// 		http.Error(w, "Band name is required", http.StatusBadRequest)
// 		return
// 	}

// 	user, err := h.Services.RegistrationSavePassword(registrationID, password)
// 	if err != nil {
// 		http.Error(w, "Unable to save password", http.StatusInternalServerError)
// 		return
// 	}

// 	if registrationID != "" {
// 		user.UserRegistrationID = registrationID
// 	}

// 	band := models.Band{}
// 	bandSlug := helpers.MakeSlug(bandName)
// 	fmt.Println("bandslug: ", bandSlug)

// 	bandID, err := database.AccessCodesTableValidateCodeReturnBandID(accessCode)
// 	if bandID == "" {
// 		log.Println("   Unable to validate code: ", err)
// 	} else {
// 		existingBand, err := database.BandsTableGetBandByBandID(bandID)
// 		if err == nil {
// 			band = existingBand
// 		}
// 		if err != nil {
// 			log.Println("   Unable to get band by band ID: ", err)
// 		}
// 	}

// 	data := models.RegistrationPages{
// 		User: user,
// 		Band: band,
// 	}

// 	fmt.Printf("User: %+v\n", data.User)
// 	fmt.Printf("Band: %+v\n", data.Band)

// 	fullName := fmt.Sprintf("%s %s", data.User.FirstName, data.User.LastName)
// 	isAdmin := true
// 	slug := helpers.MakeSlug(fullName)

// 	newUser, err := database.UsersTableCreateUser(fullName, data.User.DisplayName, slug, data.User.Email, password, passwordConfirmation, isAdmin)
// 	if err != nil {
// 		log.Println("   register err: ", err)
// 		http.Error(w, "Could not create user", http.StatusInternalServerError)
// 		return
// 	}

// 	fmt.Printf("created new user:\n %+v\n\n", newUser)

// 	bandName = HelperProcessBandNameEntry(data.Band.Name)
// 	_, err = database.BandsTableGetBandByName(bandName)
// 	if err == nil {
// 		isAdmin = false
// 	}

// 	fmt.Println("name: ", bandName)
// 	band, err = database.BandsTableGetBandByName(bandName)

// 	fmt.Printf("get band: %+v\n", band)
// 	if band.BandID == "" || band.Name == "" || err != nil {
// 		log.Println("No band exists")
// 	}
// 	if err != nil {
// 		fmt.Println("- err: ", err)
// 		bandSlug := helpers.MakeSlug(bandName)
// 		band, err = database.BandsTableCreateBand(bandName, newUser.UserID, bandSlug)
// 		if err != nil {
// 			log.Println("   register err: ", err)
// 			http.Error(w, "Could not create band", http.StatusInternalServerError)
// 			return
// 		}

// 		chatName := fmt.Sprintf("%s (Band Chat)", band.Name)
// 		chatSlug := helpers.MakeSlug(chatName)
// 		_, err := database.ChatsTableCreatePrimaryBandChat(band.BandID, chatName, chatSlug, newUser.UserID)

// 		log.Println("   Created primary band chat id: ", err)

// 		if err != nil {
// 			log.Printf("   Unable to create primary band chat for band: %v, bandID: %v, error: %v", band.Name, band.BandID, err)
// 			http.Error(w, "Unable to create primary band chat, please try again", http.StatusInternalServerError)
// 			return
// 		}
// 		// err = database.ChatMembersTableAddMember(chatID, user.UserID)
// 		// if err != nil {
// 		// 	log.Println("   Unable to add member to the chat_members table: ", err)
// 		// 	http.Error(w, "Unable to add member to the chat_members table", http.StatusInternalServerError)
// 		// 	return
// 		// }
// 	}

// 	chatID, err := database.ChatsTableGetPrimaryChatIDByBandID(band.BandID)
// 	if err != nil {
// 		log.Println("   Unable to get primary chatID by bandID: ", err)
// 		http.Error(w, "Unable to get primary chatID by bandID", http.StatusInternalServerError)
// 		return
// 	}

// 	err = database.ChatMembersTableAddMember(chatID, newUser.UserID)
// 	if err != nil {
// 		log.Println("   Unable to add new user to chat_members table: ", err)
// 		http.Error(w, "Unable to add new user to chat_members table", http.StatusInternalServerError)
// 		return
// 	}

// 	err = database.BandMembersCreateMember(band.BandID, newUser.UserID)
// 	if err != nil {
// 		log.Printf("   Unable to add user: %v to band members table: %v\n", newUser.UserID, err)
// 		http.Error(w, "Could not create band member", http.StatusInternalServerError)
// 		return
// 	}

// 	fmt.Println("New User: ", newUser)

// 	// err = h.Tmpl.ExecuteTemplate(w, "login.html", data)
// 	// if err != nil {
// 	// 	log.Println("Unable to execute register-page4-band.html", err)
// 	// 	return
// 	// }

// 	http.Redirect(w, r, "/login", http.StatusSeeOther)
// }

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

	match, err := argon2id.ComparePasswordAndHash(
		password,
		user.PasswordHash,
	)
	if err != nil {
		log.Println("   Invalid email or password")
		w.Write([]byte("* Invalid email or password * "))
		return
	}

	if !match {
		log.Println("   Invalid email or password")
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
