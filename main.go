package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

var tmpl = template.Must(template.ParseFiles("index.html"))

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/send", secondHandler)

	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl.Execute(w, nil)
}

func secondHandler(w http.ResponseWriter, r *http.Request) {
	message := r.FormValue("message")

	fmt.Println("received: ", message)

	html := fmt.Sprintf(
		"<li>%s</li>",
		message,
	)

	w.Write([]byte(html))
}
