package handlers

import (
	"bandplan/src/database"
	"bandplan/src/models"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
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
	albumTitle := r.FormValue("album-name")
	genre := r.FormValue("genre")

	musicalKey := r.FormValue("musical-key")
	tuning := r.FormValue("tuning")

	recordingBPM, err := strconv.Atoi(r.FormValue("recording-bpm"))
	if err != nil {
		fmt.Println("   Invalid entry for recordingBPM: ", err)
	}

	liveBPM, err := strconv.Atoi(r.FormValue("live-bpm"))
	if err != nil {
		log.Println("   Invalid entry for liveBPM: ", err)
	}

	lengthMinutes := r.FormValue("minutes")
	minutes, err := strconv.Atoi(r.FormValue("minutes"))
	if err != nil {
		log.Println("   Invalid entry for minutes: ", err)
	}

	lengthSeconds := r.FormValue("seconds")
	seconds, err := strconv.Atoi(r.FormValue("minutes"))
	if err != nil {
		log.Println("   Invalid entry for seconds: ", err)
	}
	songLength := minutes*60 + seconds

	fmt.Println(lengthMinutes, ":", lengthSeconds)

	releaseDateString := r.FormValue("release-date")
	var releaseDate time.Time

	if releaseDateString != "" {
		parsedDate, err := time.Parse("2006-01-02", releaseDateString)
		if err != nil {
			http.Error(w, "Invalid release date", http.StatusBadRequest)
			return
		}
		releaseDate = parsedDate
	}

	lyrics := r.FormValue("lyrics")

	spotifyLink := r.FormValue("spotify-link")
	appleMusicLink := r.FormValue("apple-music-link")
	youtubeLink := r.FormValue("youtube-link")
	otherLink := r.FormValue("other-link")

	notes := r.FormValue("notes")

	fmt.Println("Song:")
	fmt.Println(songTitle, albumTitle, genre, musicalKey, tuning, recordingBPM, liveBPM, songLength, releaseDate, lyrics, spotifyLink, appleMusicLink, youtubeLink, otherLink, notes)

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

	song := models.Song{
		Title:      songTitle,
		AlbumTitle: albumTitle,
		BandID:     band.BandID,
		Genre:      genre,

		MusicalKey:    musicalKey,
		Tuning:        tuning,
		RecordingBPM:  recordingBPM,
		LiveBPM:       liveBPM,
		LengthSeconds: songLength,

		ReleaseDate:    releaseDate,
		Lyrics:         lyrics,
		SpotifyLink:    spotifyLink,
		AppleMusicLink: appleMusicLink,
		YouTubeLink:    youtubeLink,
		OtherLink:      otherLink,

		Notes: notes,
	}

	Song, err := database.SongsTableCreateSong(song)
	if err != nil {
		http.Redirect(w, r, "/songs", http.StatusSeeOther)
	}
	fmt.Println("\nSong: ")
	fmt.Println(Song)

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
