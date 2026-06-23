package handlers

import (
	"bandplan/src/database"
	"bandplan/src/models"
	"fmt"
	"log"
	"net/http"
	"strconv"
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
		log.Println("   err gettings songs.html: ", err)
		return
	}
}

func (h Handler) HandlerSongsAddPage(w http.ResponseWriter, r *http.Request) {
	log.Print("- HandlerSongsAddPage")

	err := h.Tmpl.ExecuteTemplate(w, "songs-add.html", nil)
	if err != nil {
		log.Println("   Unable to go to add songs page: ", err)
		return
	}

	songTitle := r.FormValue("song-name")
	albumName := r.FormValue("album-name")
	genre := r.FormValue("genre")

	musicalKey := r.FormValue("musical-key")
	tuning := r.FormValue("tuning")
	recordingBPM := r.FormValue("recording-bpm")
	liveBPM := r.FormValue("live-bpm")
	lengthSeconds := r.FormValue("length-seconds")

	releaseDate := r.FormValue("release-date")
	lyrics := r.FormValue("lyrics")

	spotifyLink := r.FormValue("spotify-link")
	appleMusicLink := r.FormValue("apple-music-link")
	youtubeLink := r.FormValue("youtube-link")
	otherLink := r.FormValue("other-link")

	notes := r.FormValue("notes")

	fmt.Println(songTitle, albumName, genre, musicalKey, tuning, recordingBPM, liveBPM, lengthSeconds, releaseDate, lyrics, spotifyLink, appleMusicLink, youtubeLink, otherLink, notes)

}

func (h Handler) HandlerSongsAdd(w http.ResponseWriter, r *http.Request) {
	fmt.Println("----------------------------------------")
	log.Print("- HandlerSongsAdd")

	if r.Method == http.MethodPost {
		fmt.Println("POST")
	} else {
		fmt.Println("NOT POST")
	}

	songTitle := r.FormValue("song-title")
	albumName := r.FormValue("album-name")
	genre := r.FormValue("genre")

	musicalKey := r.FormValue("musical-key")
	tuning := r.FormValue("tuning")
	recordingBPMString := r.FormValue("recording-bpm")

	recordingBPM, err := strconv.Atoi(recordingBPMString)
	if err != nil {
		fmt.Println("invalid recordingBPM: ", err)
	}

	liveBPM := r.FormValue("live-bpm")
	lengthSeconds := r.FormValue("length-seconds")

	releaseDate := r.FormValue("release-date")
	lyrics := r.FormValue("lyrics")

	spotifyLink := r.FormValue("spotify-link")
	appleMusicLink := r.FormValue("apple-music-link")
	youtubeLink := r.FormValue("youtube-link")
	otherLink := r.FormValue("other-link")

	notes := r.FormValue("notes")

	fmt.Println("HANDLER SONGS ADD")
	fmt.Println(songTitle, albumName, genre, musicalKey, tuning, recordingBPM, liveBPM, lengthSeconds, releaseDate, lyrics, spotifyLink, appleMusicLink, youtubeLink, otherLink, notes)

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
		log.Println("   err gettings songs.html: ", err)
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
