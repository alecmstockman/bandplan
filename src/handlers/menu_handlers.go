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

func (h Handler) HandlerSongsPage(w http.ResponseWriter, r *http.Request) {
	log.Print("- HandlerSongsPage")

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

	songs, err := database.SongsTableGetAllSongsByBandID(band.BandID)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	data := models.MenuPageData{
		User:  user,
		Band:  band,
		Songs: songs,
	}

	err = h.Tmpl.ExecuteTemplate(w, "songs.html", data)
	if err != nil {
		log.Println("   err gettings songs.html: ", err)
		return
	}
}

func (h Handler) HandlerSongs(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerSongs")

	user, err := HelperGetAuthenticatedUser(r)
	if err != nil {
		fmt.Println("HandlerSend: Unable to get authenticated user: ", err)
		w.Header().Set("HX-Redirect", "/login")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	band, err := database.BandsTableGetBandByUserID(user.UserID)
	if err != nil {
		log.Println("   HandlerSend: Unable to get band by user id: ", err)
		return
	}

	songs, err := database.SongsTableGetAllSongsByBandID(band.BandID)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	data := models.MenuPageData{
		User:  user,
		Band:  band,
		Songs: songs,
	}
	fmt.Println(data)

	for _, song := range songs {
		html := fmt.Sprintf(`
			<li class="songs>
				<div class="song">%v</div>
			</li>
		`, song.Title,
		)
		w.Write([]byte(html))
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
	fmt.Println("============================================")
	log.Print("- HandlerSongsAdd")

	songTitle := r.FormValue("song-title")
	albumTitle := r.FormValue("album-name")
	genre := r.FormValue("genre")

	musicalKey := r.FormValue("musical-key")
	tuning := r.FormValue("tuning")

	recordingBPM, err := strconv.Atoi(r.FormValue("recording-bpm"))
	if err != nil {
		recordingBPM = 0
	}

	liveBPM, err := strconv.Atoi(r.FormValue("live-bpm"))
	if err != nil {
		liveBPM = 0
	}

	minutes, err := strconv.Atoi(r.FormValue("minutes"))
	if err != nil {
		log.Println("   Invalid entry for minutes: ", err)
	}

	seconds, err := strconv.Atoi(r.FormValue("minutes"))
	if err != nil {
		log.Println("   Invalid entry for seconds: ", err)
	}
	songLength := minutes*60 + seconds

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

	fmt.Println("\n\nSong: ", song)

	_, err = database.SongsTableCreateSong(song)
	if err != nil {
		http.Redirect(w, r, "/songs", http.StatusSeeOther)
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
