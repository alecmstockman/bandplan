package handlers

import (
	"bandplan/src/database"
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	DB   *sql.DB
	Tmpl *template.Template
}

func (h Handler) HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("HomeHandler: %s %s\n", r.Method, r.URL.Path)

	cookie, err := r.Cookie("session_token")
	if err != nil {
		fmt.Println("cookie err: ", err)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	fmt.Println("cookie: ", cookie)

	user, err := database.SessionsTableGetUserByToken(cookie.Value)
	if err != nil {
		fmt.Println("Unable to get User with token: ", err)
		return
	}

	fmt.Println("user: ", user)

	messages, err := database.MessagesTableGetAllMessages()
	if err != nil {
		fmt.Println("Unable to get messages from db: ", err)
		return
	}

	h.Tmpl.ExecuteTemplate(w, "index.html", messages)
}

func (h Handler) RegisterPageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\nRegisterPageHandler")

	http.ServeFile(w, r, "templates/register.html")
}

func (h Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("\nRegisterHandler: %s %s\n", r.Method, r.URL.Path)
	if r.Method == http.MethodGet {
		http.ServeFile(w, r, "templates/register.html")
		return
	}

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

func (h Handler) LoginPageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\nLoginPageHandler")
	http.ServeFile(w, r, "templates/login.html")
	return
}

func (h Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\nLogin Handler")
	if r.Method == http.MethodGet {
		fmt.Println("GET")
		http.ServeFile(w, r, "templates/login.html")
		return
	}

	if r.Method == http.MethodPost {
		fmt.Println("Post")
		email := r.FormValue("email")
		password := r.FormValue("password")

		fmt.Printf("\nemail: %s", email)
		fmt.Printf("\npassword: %s\n", password)

		user, err := database.UsersTableGetUserByEmail(email)

		if err != nil {
			fmt.Println("LoginHandler: Unable to get user: ", err)
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

		token, err := GenerateSessionToken()

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

func (h Handler) ChatPageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ChatPageHandler")
	http.ServeFile(w, r, "templates/index.html")
	return
}

func (h Handler) SendHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("SendHandler")
	message := r.FormValue("message")

	err := database.MessagesTableInsertMessage(message)
	if err != nil {
		http.Error(w, "Could not save message", http.StatusInternalServerError)
		return
	}

	html := fmt.Sprintf(
		"<li>%s</li>",
		message,
	)

	w.Write([]byte(html))
}

func (h Handler) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	err := database.MessagesTableDeleteAll()
	if err != nil {
		http.Error(w, "Could not delete messages", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(""))
	return
}

func (h Handler) MessagesHandler(w http.ResponseWriter, r *http.Request) {
	messages, err := database.MessagesTableGetAllMessages()

	if err != nil {
		http.Error(w, "Could not fetch latest messages", http.StatusInternalServerError)
		return
	}

	for _, msg := range messages {
		html := fmt.Sprintf("<li>%s</li>", msg)
		w.Write([]byte(html))
	}
}

func (h Handler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	w.Header().Set("HX-Redirect", "/login")
	w.WriteHeader(http.StatusOK)
	return
}
