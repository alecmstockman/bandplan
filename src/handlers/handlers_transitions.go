package handlers

import (
	"bandplan/src/models"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func (h Handler) HandlerTransitionPage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("---------------------------")
	log.Println("- HandlerTransitionPage")
	return
}

func (h Handler) HandlerTransitionCreatePage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("---------------------------")
	log.Println("- HandlerTransitionCreatePage")

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

	err = h.Tmpl.ExecuteTemplate(w, "transition-create.html", data)
	if err != nil {
		log.Println("   Err getting transitions-add page: ", err)
		http.Redirect(w, r, "/setlists", http.StatusSeeOther)
	}

	return
}

func (h Handler) HandlerTransitionSave(w http.ResponseWriter, r *http.Request) {
	fmt.Println("---------------------------")
	log.Println("- HandlerTransitionSave")

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		log.Println("   Unable to get AuthContext: ", err)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	band := auth.CurrentBand

	title := strings.TrimSpace(r.FormValue("transition-title"))
	if title == "" {
		log.Println("   songTitle entry was only spaces")
		http.Redirect(w, r, "/songs/add", http.StatusSeeOther)
		return
	}

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

	fmt.Println(transition)

	return
}

func (h Handler) HandlerDeleteTransition(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerDeleteTransition")
	return
}
