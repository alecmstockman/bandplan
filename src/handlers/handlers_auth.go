package handlers

import (
	"bandplan/src/database"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (h Handler) HandlerRegisterPage(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerRegisterPage")

	user, err := HelperGetAuthenticatedUser(r)
	if err == nil {
		log.Println("   Already logged in: ", user.Name)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	h.Tmpl.ExecuteTemplate(w, "register.html", nil)
	return
}

func (h Handler) HandlerRegister(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerRegister")

	if r.Method == http.MethodPost {
		name := r.FormValue("name")
		slug := HelperMakeUserSlug(name)
		displayName := r.FormValue("display-name")
		bandNameEntry := r.FormValue("band")
		email := r.FormValue("email")
		password := r.FormValue("password")

		bandName := ProcessBandNameEntry(bandNameEntry)

		user, err := database.UsersTableCreateUser(name, displayName, slug, email, password)
		if err != nil {
			log.Println("   register err: ", err)
			http.Error(w, "Could not create user", http.StatusInternalServerError)
			return
		}

		band, err := database.BandsTableGetBandByName(bandName)
		if err != nil {
			band, err = database.BandsTableCreateBand(bandName, user.UserID)
			if err != nil {
				log.Println("   register err: ", err)
				http.Error(w, "Could not create band", http.StatusInternalServerError)
				return
			}
		}

		err = database.BandMembersCreateMember(band.BandID, user.UserID)
		if err != nil {
			http.Error(w, "Could not create band member", http.StatusInternalServerError)
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

	email := r.FormValue("email")
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

	session, err := database.SessionsTableCreateSession(user.UserID, token)
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
