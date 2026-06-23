package main

import (
	"bandplan/src/database"
	"bandplan/src/handlers"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

var messages []string

func main() {
	database.DB = database.ConnectDB()
	defer database.DB.Close()

	tmpl := template.Must(template.ParseGlob("templates/*.html"))

	h := handlers.Handler{
		DB:   database.DB,
		Tmpl: tmpl,
	}

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	fmt.Println("Connecting to database...")

	database.CreateBandsTable(database.DB)
	database.CreateUsersTable(database.DB)
	database.CreateBandMembersTable(database.DB)
	database.CreateMessagesTable(database.DB)
	database.CreateSesssionsTable(database.DB)
	database.CreateSongsTable(database.DB)

	fmt.Println("Database connection succesful")

	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	http.HandleFunc("/", h.HandlerHome)
	http.HandleFunc("/register", h.HandlerRegisterPage)
	http.HandleFunc("/register/create", h.HandlerRegister)

	http.HandleFunc("/login", h.HandlerLoginPage)
	http.HandleFunc("/login/enter", h.HandlerLogin)
	http.HandleFunc("/logout", h.HandlerLogout)

	http.HandleFunc("/send", h.HandlerSend)
	http.HandleFunc("/delete", h.HandlerDelete)
	http.HandleFunc("/messages", h.HandlerMessages)

	http.HandleFunc("/songs", h.HandlerSongsPage)
	http.HandleFunc("/songlist", h.HandlerSongs)
	http.HandleFunc("/songs/add", h.HandlerSongsAddPage)
	http.HandleFunc("/songs/create", h.HandlerSongsAdd)

	http.HandleFunc("/setlists", h.HandlerSetlists)
	http.HandleFunc("/goals", h.HandlerGoals)
	http.HandleFunc("/calendar", h.HandlerCalendar)
	http.HandleFunc("/files", h.HandlerFiles)

	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
