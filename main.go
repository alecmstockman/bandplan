package main

import (
	"html/template"
	"log"
	"net/http"
)

var tmpl = template.Must(template.ParseFiles("index.html"))

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/messages", messagesHandler)

	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl.Execute(w, nil)
}

func messagesHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`<li>New message from Go server</li>`))
}
