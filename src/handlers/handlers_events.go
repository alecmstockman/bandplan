package handlers

import (
	"bandplan/src/models"
	"fmt"
	"log"
	"net/http"
)

func (h Handler) HandlerEvents(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerEvents")

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		log.Println("   Unable to get AuthContext: ", err)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	user := auth.User
	band := auth.CurrentBand

	data := models.MenuPageData{
		User: user,
		Band: band,
	}

	err = h.Tmpl.ExecuteTemplate(w, "events.html", data)
}

func (h Handler) HandlerEventCreate(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerEventCreate")

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		log.Println("   Unable to get AuthContext: ", err)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	user := auth.User
	band := auth.CurrentBand

	data := models.MenuPageData{
		User: user,
		Band: band,
	}

	err = h.Tmpl.ExecuteTemplate(w, "event_create.html", data)
}

func (h Handler) HandlerEventSave(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\n------------------------------------")
	log.Println("- HandlerEventSave")

	_, err := HelperGetAuthContext(r)
	if err != nil {
		log.Println("   Unable to get AuthContext: ", err)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	name := r.FormValue("event-name")
	date := r.FormValue("event-date")
	startTime := r.FormValue("event-start-time")
	endTime := r.FormValue("event-end-time")
	locationName := r.FormValue("event-location-name")
	location := r.FormValue("event-location")
	linkOneName := r.FormValue("link-one-name")
	linkOne := r.FormValue("link-one")
	linkeTwoName := r.FormValue("link-two-name")
	linkTwo := r.FormValue("link-two")

	fmt.Println("event name: ", name)
	fmt.Println("event date: ", date)
	fmt.Println("event start: ", startTime)
	fmt.Println("event end: ", endTime)
	fmt.Println("location name: ", locationName)
	fmt.Println("location: ", location)
	fmt.Println("link one name: ", linkOneName)
	fmt.Println("link one: ", linkOne)
	fmt.Println("link two name: ", linkeTwoName)
	fmt.Println("link two: ", linkTwo)

}
