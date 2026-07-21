package handlers

import (
	"bandplan/src/database"
	"bandplan/src/models"
	"log"
	"net/http"
	"strings"
)

func (h Handler) HandlerSetlists(w http.ResponseWriter, r *http.Request) {
	log.Print("- HandlerSetlists")

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		log.Println("   Unable to get AuthContext: ", err)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	user := auth.User
	band := auth.CurrentBand

	setlists, err := database.SetlistsTableGetSetlistsByBandID(band.BandID)
	if err != nil {
		log.Println("   Unable to get settlists from db by band ID: ", nil)
		http.Error(w, "Unable to get setlists", http.StatusInternalServerError)
	}

	data := models.SetlistData{
		User:     user,
		Band:     band,
		Setlists: setlists,
	}

	err = h.Tmpl.ExecuteTemplate(w, "setlists.html", data)
	if err != nil {
		log.Println("   Unable to render setlists: ", err)
		http.Error(w, "Unable to load setlists", http.StatusInternalServerError)
		return
	}
}

func (h Handler) HandlerSetlist(w http.ResponseWriter, r *http.Request) {
	log.Print("- HandlerSetlist")

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

	err = h.Tmpl.ExecuteTemplate(w, "setlist.html", data)
	if err != nil {
		log.Println("   unablet to get /setlist: ", err)
	}
	return
}

func (h Handler) HandlerSetlistsAddPage(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerSetlistsAddPage")

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		log.Println("   Unable to get AuthContext: ", err)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		w.Header().Set("HX-Redirect", "/")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	user := auth.User
	band := auth.CurrentBand

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

func (h Handler) HandlerSetlistCreate(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerSetlistCreate")

	user, band, err := HelperGetAuthenticatedUserAndBand(r)
	if err != nil {
		log.Println("   Unable to authenticate user:", err)
		w.Header().Set("HX-Redirect", "/")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	title := strings.TrimSpace(r.FormValue("setlist-title"))
	artworkPath := strings.TrimSpace(r.FormValue("artwork-path"))

	setlist := models.Setlist{
		BandID:        band.BandID,
		Name:          title,
		ArtworkPath:   artworkPath,
		LastUpdatedBy: user.UserID,
		CreatedBy:     user.UserID,
	}

	log.Println("   Created setlist: ", setlist)

	http.Redirect(w, r, "/setlists", http.StatusSeeOther)
	return
}
