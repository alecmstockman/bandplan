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
	tmpl := template.Must(template.ParseFiles("templates/index.html"))

	h := handlers.Handler{
		DB:   database.DB,
		Tmpl: tmpl,
	}

	database.DB = database.ConnectDB()
	defer database.DB.Close()

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	fmt.Println("Connected to database")

	database.CreateMessagesTable(database.DB)

	http.HandleFunc("/", h.HomeHandler)
	http.HandleFunc("/send", h.SendHandler)
	http.HandleFunc("/delete", h.DeleteHandler)
	http.HandleFunc("/messages", h.MessagesHandler)

	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
