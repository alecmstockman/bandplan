package handlers

import (
	"bandplan/src/database"
	"bandplan/src/helpers"
	"bandplan/src/models"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func (h Handler) HandlerBreakPage(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerBreakPage")

	breakID := r.URL.Query().Get("id")
	setlistID := r.URL.Query().Get("setlist-id")

	backURL := "/songs"
	if setlistID != "" {
		backURL = "/setlist?id=" + url.QueryEscape(setlistID)
	}

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		log.Println("   Unable to get AuthContext: ", err)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	user := auth.User
	band := auth.CurrentBand

	breakItem, err := database.BreaksTableGetBreakByID(breakID)
	if err != nil {
		http.Error(w, "Could not get song", http.StatusInternalServerError)
		return
	}

	data := models.BreakPageData{
		BackURL: backURL,
		User:    user,
		Band:    band,
		Break:   breakItem,
	}

	err = h.Tmpl.ExecuteTemplate(w, "break.html", data)
	if err != nil {
		log.Println("Unable to execute break.html:", err)
		return
	}
}

func (h Handler) HandlerBreakCreatePage(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerBreakCreate")

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

	data := models.TransitionCreateData{
		User:      user,
		Band:      band,
		SetlistID: setlistID,
	}

	err = h.Tmpl.ExecuteTemplate(w, "break-create.html", data)
	if err != nil {
		log.Println("   Err getting transitions-add page: ", err)
		http.Redirect(w, r, "/setlists", http.StatusSeeOther)
	}

}

func (h Handler) HandlerBreakSave(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerBreakSave")

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		log.Println("   Unable to get AuthContext: ", err)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	user := auth.User
	band := auth.CurrentBand
	setlistID := r.FormValue("setlist_id")

	title := strings.TrimSpace(r.FormValue("break-title"))
	if title == "" {
		log.Println("   breakTitle entry was only spaces")
		http.Redirect(w, r, "/break/create", http.StatusSeeOther)
		return
	}
	titleSlug := helpers.MakeSlug(title)

	minutes, err := strconv.Atoi(r.FormValue("minutes"))
	if err != nil {
		log.Println("   Invalid entry for minutes: ", err)
	}

	seconds, err := strconv.Atoi(r.FormValue("seconds"))
	if err != nil {
		log.Println("   Invalid entry for seconds: ", err)
	}
	breakLength := minutes*60 + seconds

	linkOne := strings.TrimSpace(r.FormValue("spotify-link"))
	linkTwo := strings.TrimSpace(r.FormValue("spotify-link"))
	notes := r.FormValue("notes")

	breakItem := models.Break{
		BandID:        band.BandID,
		Title:         title,
		Slug:          titleSlug,
		Notes:         notes,
		LengthSeconds: breakLength,
		LinkOne:       linkOne,
		LinkTwo:       linkTwo,
		CreatedBy:     user.UserID,
		UpdatedBy:     user.UserID,
	}

	newBreak, err := database.BreaksTableCreateBreak(breakItem)
	if err != nil {
		log.Println("   Unable to save break to database: ", err)
		http.Error(w, "Unable to save break to database", http.StatusInternalServerError)
		return
	}

	itemType := models.SetlistItemBreak

	_, err = database.SetlistItemsTableSaveItem(itemType, newBreak.BreakID, user.UserID, setlistID)
	if err != nil {
		log.Println("   Unable to save break to setlist: ", err)
		http.Error(w, "Unable to save break to setlist", http.StatusInternalServerError)
		return
	}

	redirectURL := "/setlist?id=" + setlistID
	log.Println("redirect URL: ", redirectURL)
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

func (h Handler) HandlerBreakEdit(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerBreakEdit")

	return
}

func (h Handler) HandlerDeleteBreak(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerDeleteBreak")

	_, err := HelperGetAuthContext(r)
	if err != nil {
		log.Println("   Unable to get AuthContext: ", err)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	setlistID := r.FormValue("setlist-id")
	breakID := r.FormValue("break-id")

	if setlistID == "" || breakID == "" {
		log.Print("   Invalid setlistID or SongID: ")
		return
	}

	err = database.BreaksTableDeleteBreak(breakID)
	if err != nil {
		log.Println("   Unable to delete break from breaks table: ", err)
		http.Error(w, "Unable to delete break from breaks table", http.StatusInternalServerError)
		return
	}

	redirectURL := "/setlist?id=" + url.QueryEscape(setlistID)
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
	return
}
