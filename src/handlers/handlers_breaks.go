package handlers

import (
	"bandplan/src/database"
	"bandplan/src/helpers"
	"bandplan/src/models"
	"database/sql"
	"errors"
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
		backURL = "/setlist?id=" + url.QueryEscape(setlistID) + "&from=setlist&setlist-id=" + url.QueryEscape(setlistID)
	}

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		log.Println("   Unable to get AuthContext: ", err)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	user := auth.User
	band := auth.CurrentBand

	breakItem, err := database.BreaksTableGetBreakByID(breakID, band.BandID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "Could not get break", http.StatusInternalServerError)
		return
	}

	data := models.BreakPageData{
		SetlistID: setlistID,
		BackURL:   backURL,
		User:      user,
		Band:      band,
		Break:     breakItem,
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

	linkOne := strings.TrimSpace(r.FormValue("link-one"))
	linkTwo := strings.TrimSpace(r.FormValue("link-two"))
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

func (h Handler) HandlerBreakEditPage(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerBreakEditPage")

	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		log.Println("   Unable to get AuthContext: ", err)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	breakID := strings.TrimSpace(r.URL.Query().Get("id"))
	if breakID == "" {
		http.Error(w, "Missing break ID", http.StatusBadRequest)
		return
	}
	setlistID := strings.TrimSpace(r.URL.Query().Get("setlist-id"))
	if breakID == "" {
		http.Error(w, "Missing setlist ID", http.StatusBadRequest)
		return
	}

	breakItem, err := database.BreaksTableGetBreakByID(breakID, auth.CurrentBand.BandID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.NotFound(w, r)
			return
		}
		log.Println("   Unable to get break: ", err)
		http.Error(w, "Unable to get break", http.StatusInternalServerError)
		return
	}

	data := models.BreakPageData{
		SetlistID: setlistID,
		BackURL:   "/break?id=" + url.QueryEscape(breakID) + "&from=setlist&setlist-id=" + url.QueryEscape(setlistID),
		User:      auth.User,
		Band:      auth.CurrentBand,
		Break:     breakItem,
	}

	if err := h.Tmpl.ExecuteTemplate(w, "break-edit.html", data); err != nil {
		log.Println("   Unable to execute break-edit.html: ", err)
		http.Error(w, "Unable to load break editor", http.StatusInternalServerError)
	}
}

func (h Handler) HandlerBreakUpdate(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerBreakUpdate")

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		log.Println("   Unable to get AuthContext: ", err)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	breakID := strings.TrimSpace(r.FormValue("break-id"))
	if breakID == "" {
		http.Error(w, "Missing break ID", http.StatusBadRequest)
		return
	}

	existingBreak, err := database.BreaksTableGetBreakByID(breakID, auth.CurrentBand.BandID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.NotFound(w, r)
			return
		}
		log.Println("   Unable to get break: ", err)
		http.Error(w, "Unable to get break", http.StatusInternalServerError)
		return
	}

	title := strings.TrimSpace(r.FormValue("break-title"))
	if title == "" {
		redirectURL := "/break/edit?id=" + url.QueryEscape(breakID)
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
		return
	}
	if len(title) > 300 {
		http.Error(w, "Break title is too long", http.StatusBadRequest)
		return
	}

	minutes, err := strconv.Atoi(r.FormValue("minutes"))
	if err != nil || minutes < 0 || minutes > 99 {
		http.Error(w, "Invalid minutes", http.StatusBadRequest)
		return
	}

	seconds, err := strconv.Atoi(r.FormValue("seconds"))
	if err != nil || seconds < 0 || seconds > 59 {
		http.Error(w, "Invalid seconds", http.StatusBadRequest)
		return
	}

	slug := existingBreak.Slug
	if title != existingBreak.Title {
		slug = helpers.MakeSlug(title)
	}

	breakItem := models.Break{
		BreakID:       breakID,
		BandID:        auth.CurrentBand.BandID,
		Title:         title,
		Slug:          slug,
		LengthSeconds: minutes*60 + seconds,
		LinkOne:       strings.TrimSpace(r.FormValue("link-one")),
		LinkTwo:       strings.TrimSpace(r.FormValue("link-two")),
		Notes:         r.FormValue("notes"),
		UpdatedBy:     auth.User.UserID,
	}

	updated, err := database.BreaksTableUpdateBreak(breakItem)
	if err != nil {
		log.Println("   Unable to update break: ", err)
		http.Error(w, "Unable to update break", http.StatusInternalServerError)
		return
	}
	if !updated {
		http.NotFound(w, r)
		return
	}

	redirectURL := "/break?id=" + url.QueryEscape(breakID)
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
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
