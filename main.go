package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

var db *sql.DB
var tmpl = template.Must(template.ParseFiles("index.html"))
var messages []string

func main() {
	db = connectDB()
	defer db.Close()

	fmt.Println("Connected to database")

	createTables(db)

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/send", sendHandler)
	http.HandleFunc("/delete", deleteHandler)

	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	messages, err := messagesTableGetAllMessages()
	if err != nil {
		fmt.Println("Unable to get messages from db")
	}

	tmpl.Execute(w, messages)
}

func sendHandler(w http.ResponseWriter, r *http.Request) {
	message := r.FormValue("message")

	err := messagesTableInsertMessage(message)
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

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	err := messagesTableDeleteAll()
	if err != nil {
		http.Error(w, "Could not delete messages", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(""))
	return
}
