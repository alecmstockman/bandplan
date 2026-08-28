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

	transition, err := database.TransitionsTableGetTransitionByID(transitionID)
	if err != nil {
		http.Error(w, "Could not get song", http.StatusInternalServerError)
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

	link := r.FormValue("spotify-link")
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
	}

	newTransition, err := database.TransitionsTableCreateTransition(transition)
	if err != nil {
		log.Println("   Unable to save transition to db: ", err)
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

func (h Handler) HandlerTransitionEdit(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerTransitionEdit")

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		log.Println("   Unable to get AuthContext: ", err)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	user := auth.User
	band := auth.CurrentBand
	transitionID := r.URL.Query().Get("id")

	transition, err := database.TransitionsTableGetTransitionByID(transitionID)
	if err != nil {
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
		log.Println("   Err getting transitions-add page: ", err)
		http.Redirect(w, r, "/setlists", http.StatusSeeOther)
	}
}
