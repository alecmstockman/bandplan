package handlers

import (
	"bandplan/src/models"
	"log"
	"net/http"
)

func (h Handler) HandlerPromotion(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerPromotions")

	user, band, err := HelperGetAuthenticatedUserAndBand(r)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	data := models.MenuPageData{
		User: user,
		Band: band,
	}

	err = h.Tmpl.ExecuteTemplate(w, "promotion.html", data)
}

func (h Handler) HandlerGoals(w http.ResponseWriter, r *http.Request) {
	log.Print("- HandlerSetlists")

	user, band, err := HelperGetAuthenticatedUserAndBand(r)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	data := models.MenuPageData{
		User: user,
		Band: band,
	}

	err = h.Tmpl.ExecuteTemplate(w, "goals.html", data)
}

func (h Handler) HandlerCalendar(w http.ResponseWriter, r *http.Request) {
	log.Print("- HandlerCalendar")

	user, band, err := HelperGetAuthenticatedUserAndBand(r)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	data := models.MenuPageData{
		User: user,
		Band: band,
	}

	err = h.Tmpl.ExecuteTemplate(w, "calendar.html", data)
}

func (h Handler) HandlerEvents(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerEvents")

	user, band, err := HelperGetAuthenticatedUserAndBand(r)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	data := models.MenuPageData{
		User: user,
		Band: band,
	}

	err = h.Tmpl.ExecuteTemplate(w, "events.html", data)
}

func (h Handler) HandlerFiles(w http.ResponseWriter, r *http.Request) {
	log.Print("- HandlerFile")

	user, band, err := HelperGetAuthenticatedUserAndBand(r)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	data := models.MenuPageData{
		User: user,
		Band: band,
	}

	err = h.Tmpl.ExecuteTemplate(w, "files.html", data)
}
