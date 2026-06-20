package handlers

import (
	"bandplan/src/database"
	"bandplan/src/models"
	"fmt"
	"log"
	"net/http"
)

func (h Handler) HandlerSongs(w http.ResponseWriter, r *http.Request) {
	log.Print("- HandlerSongs")

	user, err := HelperGetAuthenticatedUser(r)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	band, err := database.BandsTableGetBandByUserID(user.UserID)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	data := models.MenuPageData{
		User: user,
		Band: band,
	}

	err = h.Tmpl.ExecuteTemplate(w, "songs.html", data)
	if err != nil {
		fmt.Println("\n\n----------------------------------------\nerr gettings songs.html: ", err)
		return
	}
}

func (h Handler) HandlerSongsAdd(w http.ResponseWriter, r *http.Request) {
	fmt.Println("----------------------------------------")
	log.Print("- HandlerSongsAdd")

	err := h.Tmpl.ExecuteTemplate(w, "songs-add.html", nil)
	if err != nil {
		log.Println("   Unable to go to songs add: ", err)
		return
	}
}

func (h Handler) HandlerSetlists(w http.ResponseWriter, r *http.Request) {
	log.Print("- HandlerSetlists")

	user, err := HelperGetAuthenticatedUser(r)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	band, err := database.BandsTableGetBandByUserID(user.UserID)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	data := models.MenuPageData{
		User: user,
		Band: band,
	}

	err = h.Tmpl.ExecuteTemplate(w, "setlists.html", data)
	return
}

func (h Handler) HandlerGoals(w http.ResponseWriter, r *http.Request) {
	log.Print("- HandlerSetlists")

	user, err := HelperGetAuthenticatedUser(r)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	band, err := database.BandsTableGetBandByUserID(user.UserID)
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

	user, err := HelperGetAuthenticatedUser(r)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	band, err := database.BandsTableGetBandByUserID(user.UserID)
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

func (h Handler) HandlerFiles(w http.ResponseWriter, r *http.Request) {
	fmt.Println("-----------------------------------------------")
	log.Print("- HandlerFile")

	user, err := HelperGetAuthenticatedUser(r)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	band, err := database.BandsTableGetBandByUserID(user.UserID)
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
