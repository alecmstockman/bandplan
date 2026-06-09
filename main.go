package main

import (
	"bandplan/src/database"
	"bandplan/src/handlers"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

// var db *sql.DB
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

	database.CreateTables(database.DB)

	http.HandleFunc("/", h.HomeHandler)
	http.HandleFunc("/send", h.SendHandler)
	http.HandleFunc("/delete", h.DeleteHandler)
	http.HandleFunc("/messages", h.MessagesHandler)

	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// func homeHandler(w http.ResponseWriter, r *http.Request) {
// 	messages, err := handlers.messagesTableGetAllMessages()
// 	if err != nil {
// 		fmt.Println("Unable to get messages from db")
// 	}

// 	tmpl.Execute(w, messages)
// }

// func sendHandler(w http.ResponseWriter, r *http.Request) {
// 	message := r.FormValue("message")

// 	err := handlers.messagesTableInsertMessage(message)
// 	if err != nil {
// 		http.Error(w, "Could not save message", http.StatusInternalServerError)
// 		return
// 	}

// 	html := fmt.Sprintf(
// 		"<li>%s</li>",
// 		message,
// 	)

// 	w.Write([]byte(html))
// }

// func deleteHandler(w http.ResponseWriter, r *http.Request) {
// 	err := messagesTableDeleteAll()
// 	if err != nil {
// 		http.Error(w, "Could not delete messages", http.StatusInternalServerError)
// 		return
// 	}

// 	w.Write([]byte(""))
// 	return
// }

// func messagesHandler(w http.ResponseWriter, r *http.Request) {
// 	messages, err := messagesTableGetAllMessages()

// 	if err != nil {
// 		http.Error(w, "Could not fetch latest messages", http.StatusInternalServerError)
// 		return
// 	}

// 	for _, msg := range messages {
// 		html := fmt.Sprintf("<li>%s</li>", msg)
// 		w.Write([]byte(html))
// 	}
// }
