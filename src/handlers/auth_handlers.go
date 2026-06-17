package handlers

import (
	"bandplan/src/database"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (h Handler) HandlerRegisterPage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\nHandlerRegisterPage")

	user, err := HelperGetAuthenticatedUser(r)
	if err == nil {
		fmt.Println("Already logged in: ", user.Name)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	h.Tmpl.ExecuteTemplate(w, "register.html", nil)
	return
}

func (h Handler) HandlerRegister(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("\nHandlerRegister: %s %s\n", r.Method, r.URL.Path)

	if r.Method == http.MethodPost {
		name := r.FormValue("name")
		band := r.FormValue("band")
		email := r.FormValue("email")
		password := r.FormValue("password")

		_, err := database.UsersTableCreateUser(name, band, email, password)
		if err != nil {
			fmt.Println("register err: ", err)
			http.Error(w, "Could not create user", http.StatusInternalServerError)
			return
		}

		fmt.Println(name, band, email, password)

		w.Header().Set("HX-Redirect", "/login")
		w.WriteHeader(http.StatusOK)
		return
	}
	return
}

func (h Handler) HandlerLoginPage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\nHandlerLoginPage")

	user, err := HelperGetAuthenticatedUser(r)
	if err == nil {
		fmt.Println("Already logged in: ", user.Name)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	h.Tmpl.ExecuteTemplate(w, "login.html", nil)
	return
}

func (h Handler) HandlerLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Println("HANDLER LOGIN")

	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		password := r.FormValue("password")
		user, err := database.UsersTableGetUserByEmail(email)

		if err != nil {
			fmt.Println("HandlerLogin: Unable to get user: ", err)
			return
		}

		err = bcrypt.CompareHashAndPassword(
			[]byte(user.PasswordHash),
			[]byte(password),
		)

		if err != nil {
			fmt.Println("Invalid email or password")
			w.Write([]byte("Invalid email or password"))
			return
		}

		token, err := HelperGenerateSessionToken()

		if err != nil {
			fmt.Println("Unabe to generate token: ", err)
			return
		}

		session, err := database.SessionsTableCreateSession(user.UserID, token)

		if err != nil {
			fmt.Println("Unable to create session: ", err)
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    session.Token,
			Path:     "/",
			MaxAge:   3600,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		})

		w.Header().Set("HX-Redirect", "/")
		w.WriteHeader(http.StatusOK)
		return
	}
	return
}
