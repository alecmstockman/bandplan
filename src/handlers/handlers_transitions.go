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

func (h Handler) HandlerTransitionPage(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerTransitionPage")

	transitionID := r.URL.Query().Get("id")
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

	transition, err := database.TransitionsTableGetTransitionByID(transitionID, band.BandID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "Could not get transition", http.StatusInternalServerError)
		return
	}

	data := models.TransitionPageData{
		BackURL:    backURL,
		User:       user,
		Band:       band,
		Transition: transition,
	}

	err = h.Tmpl.ExecuteTemplate(w, "transition.html", data)
	if err != nil {
		log.Println("Unable to execute transition.html:", err)
		return
	}

	return
}

func (h Handler) HandlerTransitionCreatePage(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerTransitionCreatePage")

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		log.Println("   Unable to get AuthContext: ", err)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
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

	err = h.Tmpl.ExecuteTemplate(w, "transition-create.html", data)
	if err != nil {
		log.Println("   Err getting transitions-add page: ", err)
		http.Redirect(w, r, "/setlists", http.StatusSeeOther)
	}
}

func (h Handler) HandlerTransitionSave(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerTransitionSave")

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		log.Println("   Unable to get AuthContext: ", err)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	user := auth.User
	band := auth.CurrentBand
	setlistID := r.FormValue("setlist_id")

	title := strings.TrimSpace(r.FormValue("transition-title"))
	if title == "" {
		log.Println("   songTitle entry was only spaces")
		http.Redirect(w, r, "/transition/create", http.StatusSeeOther)
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
	songLength := minutes*60 + seconds

	bpm, err := strconv.Atoi(r.FormValue("bpm"))
	if err != nil {
		bpm = 0
	}

	timeSignature := strings.TrimSpace(r.FormValue("time-signature"))
	key := strings.TrimSpace(r.FormValue("key"))
	tuning := strings.TrimSpace(r.FormValue("tuning"))
	capo := strings.TrimSpace(r.FormValue("capo"))

	explicitEntry := r.FormValue("explicit")

	var explicit bool

	if explicitEntry == "on" {
		explicit = true
	} else {
		explicit = false
	}

	link := strings.TrimSpace(r.FormValue("link-one"))
	lyrics := r.FormValue("lyrics")
	notes := r.FormValue("notes")

	transition := models.Transition{
		BandID: band.BandID,
		Title:  title,
		Slug:   titleSlug,

		LengthSeconds: songLength,
		BPM:           bpm,
		TimeSignature: timeSignature,
		Key:           key,

		Tuning:   tuning,
		Capo:     capo,
		Explicit: explicit,

		LinkOne: link,
		Lyrics:  lyrics,
		Notes:   notes,

		CreatedBy: user.UserID,
		UpdatedBy: user.UserID,
	}

	newTransition, err := database.TransitionsTableCreateTransition(transition)
	if err != nil {
		log.Println("   Unable to save transition to db: ", err)
		http.Error(w, "Unable to save transition", http.StatusInternalServerError)
		return
	}

	itemType := models.SetlistItemTransition

	_, err = database.SetlistItemsTableSaveItem(itemType, newTransition.TransitionID, user.UserID, setlistID)
	if err != nil {
		log.Println("   Unable to save transition to setlist: ", err)
		http.Error(w, "Unable to save transtion to setlist", http.StatusInternalServerError)
		return
	}

	redirectURL := "/setlist?id=" + setlistID
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
	return
}

func (h Handler) HandlerDeleteTransition(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerDeleteTransition")

	_, err := HelperGetAuthContext(r)
	if err != nil {
		log.Println("   Unable to get AuthContext: ", err)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	setlistID := r.FormValue("setlist-id")
	transitionID := r.FormValue("transition-id")

	position, err := strconv.Atoi(r.FormValue("position"))
	if err != nil {
		log.Println("   Unable to get transition position: ", err)
		return
	}

	if setlistID == "" || transitionID == "" {
		log.Print("   Invalid setlistID or SongID: ")
		return
	}

	err = database.SetlistItemsTableDeleteTransition(transitionID, position, setlistID)
	if err != nil {
		log.Printf("   Unable to delete transition id %v to setlist: %v\n", transitionID, err)
		http.Error(w, "Unable to delete transition to setlist", http.StatusInternalServerError)
		return
	}

	redirectURL := "/setlist?id=" + url.QueryEscape(setlistID)
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
	return
}

func (h Handler) HandlerTransitionEditPage(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerTransitionEditPage")

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

	user := auth.User
	band := auth.CurrentBand
	transitionID := strings.TrimSpace(r.URL.Query().Get("id"))
	if transitionID == "" {
		http.Error(w, "Missing transition ID", http.StatusBadRequest)
		return
	}

	transition, err := database.TransitionsTableGetTransitionByID(transitionID, band.BandID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.NotFound(w, r)
			return
		}
		log.Println("   Unable to get transition: ", err)
		http.Error(w, "Could not get transition", http.StatusInternalServerError)
		return
	}

	data := models.TransitionPageData{
		BackURL:    "/transition?id=" + transitionID,
		User:       user,
		Band:       band,
		Transition: transition,
	}

	err = h.Tmpl.ExecuteTemplate(w, "transition-edit.html", data)
	if err != nil {
		log.Println("   Unable to execute transition-edit.html: ", err)
		http.Error(w, "Unable to load transition editor", http.StatusInternalServerError)
	}
}

func (h Handler) HandlerTransitionUpdate(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerTransitionUpdate")

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

	transitionID := strings.TrimSpace(r.FormValue("transition-id"))
	if transitionID == "" {
		http.Error(w, "Missing transition ID", http.StatusBadRequest)
		return
	}

	existingTransition, err := database.TransitionsTableGetTransitionByID(transitionID, auth.CurrentBand.BandID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.NotFound(w, r)
			return
		}
		log.Println("   Unable to get transition: ", err)
		http.Error(w, "Unable to get transition", http.StatusInternalServerError)
		return
	}

	title := strings.TrimSpace(r.FormValue("transition-title"))
	if title == "" {
		redirectURL := "/transition/edit?id=" + url.QueryEscape(transitionID)
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
		return
	}
	if len(title) > 300 {
		http.Error(w, "Transition title is too long", http.StatusBadRequest)
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

	bpm := 0
	bpmEntry := strings.TrimSpace(r.FormValue("bpm"))
	if bpmEntry != "" {
		bpm, err = strconv.Atoi(bpmEntry)
		if err != nil || bpm < 0 || bpm > 500 {
			http.Error(w, "Invalid BPM", http.StatusBadRequest)
			return
		}
	}

	slug := existingTransition.Slug
	if title != existingTransition.Title {
		slug = helpers.MakeSlug(title)
	}

	transition := models.Transition{
		TransitionID: transitionID,
		BandID:       auth.CurrentBand.BandID,
		Title:        title,
		Slug:         slug,

		LengthSeconds: minutes*60 + seconds,
		BPM:           bpm,
		TimeSignature: strings.TrimSpace(r.FormValue("time-signature")),
		Key:           strings.TrimSpace(r.FormValue("key")),
		Tuning:        strings.TrimSpace(r.FormValue("tuning")),
		Capo:          strings.TrimSpace(r.FormValue("capo")),
		Explicit:      r.FormValue("explicit") == "on",

		LinkOne:   strings.TrimSpace(r.FormValue("link-one")),
		Lyrics:    r.FormValue("lyrics"),
		Notes:     r.FormValue("notes"),
		UpdatedBy: auth.User.UserID,
	}

	updated, err := database.TransitionsTableUpdateTransition(transition)
	if err != nil {
		log.Println("   Unable to update transition: ", err)
		http.Error(w, "Unable to update transition", http.StatusInternalServerError)
		return
	}
	if !updated {
		http.NotFound(w, r)
		return
	}

	redirectURL := "/transition?id=" + url.QueryEscape(transitionID)
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}
