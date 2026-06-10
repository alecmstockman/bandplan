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

	tmpl := template.Must(template.ParseFiles("templates/register.html"))

	h := handlers.Handler{
		DB:   database.DB,
		Tmpl: tmpl,
	}

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	fmt.Println("Connected to database")

	database.CreateMessagesTable(database.DB)
	database.CreateUsersTable(database.DB)

	http.HandleFunc("/", h.RegisterPageHandler)
	http.HandleFunc("/register", h.RegisterHandler)
	http.HandleFunc("/login", h.LoginPageHandler)

	http.HandleFunc("/send", h.SendHandler)
	http.HandleFunc("/delete", h.DeleteHandler)
	http.HandleFunc("/messages", h.MessagesHandler)

	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
