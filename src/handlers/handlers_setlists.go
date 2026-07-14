package handlers

import (
	"bandplan/src/models"
	"fmt"
	"log"
	"net/http"
)

func (h Handler) HandlerSetlists(w http.ResponseWriter, r *http.Request) {
	log.Print("- HandlerSetlists")

	user, band, err := HelperGetAuthenticatedUserAndBand(r)
	if err != nil {
		log.Println("   Unable to authenticate user: ", err)
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

func (h Handler) HandlerSetlist(w http.ResponseWriter, r *http.Request) {
	log.Print("- HandlerSetlist")

	user, band, err := HelperGetAuthenticatedUserAndBand(r)
	if err != nil {
		log.Println("   Unable to authenticate user: ", err)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	data := models.MenuPageData{
		User: user,
		Band: band,
	}

	err = h.Tmpl.ExecuteTemplate(w, "setlist.html", data)
	if err != nil {
		log.Println("   unablet to get /setlist: ", err)
	}
	return
}

func (h Handler) HandlerSetlistsAddPage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("--------------------------------------")
	log.Println("- HandlerSetlistsAddPage")

	user, band, err := HelperGetAuthenticatedUserAndBand(r)
	if err != nil {
		log.Println("   Unable to authenticate user: ", err)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	data := models.SongDownloadData{
		User: user,
		Band: band,
	}

	err = h.Tmpl.ExecuteTemplate(w, "setlist-add.html", data)
	if err != nil {
		log.Println("   Err getting setlist-add page: ", err)
		http.Redirect(w, r, "/setlists", http.StatusSeeOther)
	}

	return
}
