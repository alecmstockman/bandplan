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
	if err != nil {
		log.Println("   Unable to open events page: ", err)
		return
	}
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
	eventType := r.FormValue("event-type")
	venueName := r.FormValue("venue-name")
	eventAddress := r.FormValue("event-location")

	startTime := r.FormValue("event-start-time")
	endTime := r.FormValue("event-end-time")
	timeZone := r.FormValue("event-timezone")
	recurrence := (r.FormValue("event-recurrence"))

	ticketLink := r.FormValue("event-ticket-link")
	presaleTicketPrice := r.FormValue("event-presale-ticket-price")
	ticketPrice := r.FormValue("event-ticket-price")

	setLocation := r.FormValue("event-set-location")
	setTime := r.FormValue("event-set-time")
	setLength := r.FormValue("event-set-length")
	loadinInstructions := r.FormValue("load-in-instructions")

	linkOneName := r.FormValue("link-one-name")
	linkOne := r.FormValue("link-one")
	linkeTwoName := r.FormValue("link-two-name")
	linkTwo := r.FormValue("link-two")

	eventNotes := r.FormValue("event-notes")

	fmt.Println("event name:    ", name)
	fmt.Println("event date:    ", date)
	fmt.Println("event type:    ", eventType)
	fmt.Println("venue name:    ", venueName)
	fmt.Println("event address: ", eventAddress)

	fmt.Println("event start:   ", startTime)
	fmt.Println("event end:     ", endTime)
	fmt.Println("time zone:     ", timeZone)
	fmt.Println("event repeats: ", recurrence)

	fmt.Println("ticket link:   ", ticketLink)
	fmt.Println("presale price: ", presaleTicketPrice)
	fmt.Println("ticket price:  ", ticketPrice)

	fmt.Println("set location:  ", setLocation)
	fmt.Println("set time:      ", setTime)
	fmt.Println("set length:    ", setLength)
	fmt.Println("load in inst.: ", loadinInstructions)

	fmt.Println("link one name: ", linkOneName)
	fmt.Println("link one:      ", linkOne)
	fmt.Println("link two name: ", linkeTwoName)
	fmt.Println("link two:      ", linkTwo)

	fmt.Println("event notes:   ", eventNotes)

}
