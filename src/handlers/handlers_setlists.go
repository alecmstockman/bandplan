package handlers

import (
	"bandplan/src/database"
	"bandplan/src/models"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func (h Handler) HandlerSetlistsPage(w http.ResponseWriter, r *http.Request) {
	log.Print("- HandlerSetlistsPage")

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
		log.Println("   unable to get /setlist: ", err)
	}
	return
}

func (h Handler) HandlerSetlistsAddPage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("----------------------------------")
	log.Println("- HandlerSetlistsAddPage")

	log.Printf(
		"- HandlerSetlistsAddPage method=%s path=%s hx-request=%q",
		r.Method,
		r.URL.Path,
		r.Header.Get("HX-Request"),
	)

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
	fmt.Println("----------------------------------")
	log.Println("- HandlerSetlistCreate")

	user, band, err := HelperGetAuthenticatedUserAndBand(r)
	if err != nil {
		log.Println("   Unable to authenticate user:", err)
		w.Header().Set("HX-Redirect", "/")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	title := strings.TrimSpace(r.FormValue("setlist-title"))
	notes := strings.TrimSpace(r.FormValue("notes"))

	file, header, err := r.FormFile("artwork-path")
	if err != nil && err != http.ErrMissingFile {
		log.Println("Unable to read artwork file:", err)
		http.Error(w, "Unable to read artwork", http.StatusBadRequest)
		return
	}

	if err == nil {
		defer file.Close()

		fmt.Println("Artwork filename:", header.Filename)
		fmt.Println("Artwork size:", header.Size)

		// Save the file and assign the resulting server path.
	}
	fmt.Println("   Title: ", title)
	// fmt.Println("   artworkPath: ", artworkPath)

	setlist := models.Setlist{
		BandID: band.BandID,
		Name:   title,
		Notes:  notes,
		// ArtworkPath: artworkPath,
		CreatedBy: user.UserID,
		UpdatedBy: user.UserID,
	}

	err = database.SetlistsTableCreateSetlist(setlist)
	if err != nil {
		log.Println("   Unable to create setlist in database: ", err)
		http.Error(w, "/setlists", http.StatusInternalServerError)
		return
	}

	log.Println("   Created setlist: ", setlist.Name)

	w.Header().Set("HX-Redirect", "/setlists")
	w.WriteHeader(http.StatusSeeOther)
	return
}

func (h Handler) HandlerSetlistsDelete(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerSetlistsDelete")

	_, err := HelperGetAuthContext(r)
	if err != nil {
		log.Println("   Unable to get AuthContext: ", err)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	setlistID := r.FormValue("setlist-id")

	err = database.SetlistsTableDeleteSetlist(setlistID)
	if err != nil {
		log.Printf("\n   Unable to delete setlist %s from setlists table: ", setlistID, err)
	}

	log.Println("   Deleted setlist: ", setlistID)

	w.Header().Set("HX-Redirect", "/setlists")
	w.WriteHeader(http.StatusSeeOther)
	return
}

func (h Handler) HandlerSetlistAddSong(w http.ResponseWriter, r *http.Request) {
	fmt.Println("---------------------------------")
	log.Println("- HandlerSetlistAddSong")

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		log.Println("   Unable to get AuthContext: ", err)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		w.Header().Set("HX-Redirect", "/")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	user := auth.User

	setlistID := r.FormValue("setlist-id")
	songID := r.FormValue("song-id")

	if setlistID == "" || songID == "" {
		log.Print("   Invalid setlistID or SongID: ")
		return
	}

	song, err := database.SetlistItemsTableSaveSong(songID, user.UserID, setlistID)
	if err != nil {
		log.Println("   Unable to save song to setlist: ", err)
		http.Error(w, "Unable to save song to setlist", http.StatusInternalServerError)
		return
	}
	log.Printf("   Saved %s to setlist_songs\n", song.Song.Title)

	return
}

func (h Handler) HandlerSetlistDeleteSong(w http.ResponseWriter, r *http.Request) {
	fmt.Println("---------------------------------")
	log.Println("- HandlerSetlistDeleteSong")

	_, err := HelperGetAuthContext(r)
	if err != nil {
		log.Println("   Unable to get AuthContext: ", err)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		w.Header().Set("HX-Redirect", "/")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	setlistID := r.FormValue("setlist-id")
	songID := r.FormValue("song-id")

	fmt.Println("\nsetlistID: ", setlistID)
	fmt.Println("\nsongID; ", songID)

	if setlistID == "" || songID == "" {
		log.Print("   Invalid setlistID or SongID: ")
		return
	}

	err = database.SetlistItemsTableDeleteSong(songID, setlistID)
	if err != nil {
		log.Println("   Unable to delete song to setlist: ", err)
		http.Error(w, "Unable to delete song to setlist", http.StatusInternalServerError)
		return
	}

	redirectURL := "/setlist?id=" + url.QueryEscape(setlistID)

	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

func (h Handler) HandlerSetlistPage(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerSetlistPage")

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
	setlistID := r.URL.Query().Get("id")

	setlist, err := database.SetlistsTableGetSetlistByID(setlistID)
	if err != nil {
		log.Println("   Unable to get setlist: ", err)
		return
	}

	data := models.SetlistPage{
		User:    user,
		Band:    band,
		Setlist: setlist,
	}

	err = h.Tmpl.ExecuteTemplate(w, "setlist.html", data)
	if err != nil {
		log.Println("   Unable to render setlist: ", err)
		http.Error(w, "Unable to load setlist", http.StatusInternalServerError)
		return
	}

	return
}
