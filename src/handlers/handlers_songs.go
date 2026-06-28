package handlers

import (
	"bandplan/src/database"
	"bandplan/src/models"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
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

func (h *Handler) HandlerSongsSearch(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerSongsSearch")
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

	query := r.FormValue("q")

	log.Println("Search Query: ", query)

	songs, err := database.SongsTableSearchByBandID(band.BandID, query)
	if err != nil {
		http.Error(w, "Could not search songs", http.StatusInternalServerError)
		return
	}
	data := models.MenuPageData{
		User:  user,
		Band:  band,
		Songs: songs,
	}
	err = h.Tmpl.ExecuteTemplate(w, "songs-list.html", data)
	if err != nil {
		log.Println("   Err getting songs-list from search: ", err)
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
	log.Print("- HandlerSongsAdd")

	songTitle := strings.TrimSpace(r.FormValue("song-title"))
	if songTitle == "" {
		log.Print("   songTitle entry was only spaces")
		http.Redirect(w, r, "/songs/add", http.StatusSeeOther)
		return
	}

	albumTitle := strings.TrimSpace(r.FormValue("album-name"))
	genre := strings.TrimSpace(r.FormValue("genre"))
	musicalKey := strings.TrimSpace(r.FormValue("musical-key"))
	tuning := strings.TrimSpace(r.FormValue("tuning"))

	recordingBPM, err := strconv.Atoi(r.FormValue("recording-bpm"))
	if err != nil {
		recordingBPM = 0
	}

	liveBPM, err := strconv.Atoi(r.FormValue("live-bpm"))
	if err != nil {
		liveBPM = 0
	}

	minutes, err := strconv.Atoi(r.FormValue("minutes"))

	fmt.Println("****", minutes)
	if err != nil {
		log.Println("   Invalid entry for minutes: ", err)
	}

	seconds, err := strconv.Atoi(r.FormValue("seconds"))
	if err != nil {
		log.Println("   Invalid entry for seconds: ", err)
	}
	songLength := minutes*60 + seconds

	releaseDateString := r.FormValue("release-date")
	var releaseDate time.Time

	fmt.Println("   *** Release date: ", releaseDateString)
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

	_, err = database.SongsTableCreateSong(song)
	if err != nil {
		http.Redirect(w, r, "/songs", http.StatusSeeOther)
	}

	http.Redirect(w, r, "/songs", http.StatusSeeOther)
}

func (h Handler) HandlerSongPage(w http.ResponseWriter, r *http.Request) {
	log.Print("- HandlerSongPage")

	songID := r.URL.Query().Get("id")

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

	song, err := database.SongsTableGetSongBySongID(songID)
	if err != nil {
		http.Error(w, "Could not get song", http.StatusInternalServerError)
		return
	}

	data := models.SongPageData{
		User: user,
		Band: band,
		Song: song,
	}

	err = h.Tmpl.ExecuteTemplate(w, "song.html", data)
	if err != nil {
		log.Println("Unable to execute song.html:", err)
		return
	}
}

func (h Handler) HandlerSongEditPage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("================================================")
	log.Print("- HandlerSongEditPage")

	songID := r.URL.Query().Get("id")

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

	song, err := database.SongsTableGetSongBySongID(songID)
	if err != nil {
		http.Error(w, "Could not get song", http.StatusInternalServerError)
		return
	}

	data := models.SongPageData{
		User: user,
		Band: band,
		Song: song,
	}

	err = h.Tmpl.ExecuteTemplate(w, "song-edit.html", data)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
}

func (h Handler) HandlerSongUpdate(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\n -----------------------------------------------------")
	log.Print("- HandlerSongUpdate")
	fmt.Println(r.FormValue("song-id"))

	songTitle := strings.TrimSpace(r.FormValue("song-title"))
	if songTitle == "" {
		log.Print("   songTitle entry was only spaces")
		http.Redirect(w, r, "/songs/add", http.StatusSeeOther)
		return
	}

	albumTitle := strings.TrimSpace(r.FormValue("album-name"))
	genre := strings.TrimSpace(r.FormValue("genre"))
	musicalKey := strings.TrimSpace(r.FormValue("musical-key"))
	tuning := strings.TrimSpace(r.FormValue("tuning"))

	recordingBPM, err := strconv.Atoi(r.FormValue("recording-bpm"))
	if err != nil {
		recordingBPM = 0
	}

	liveBPM, err := strconv.Atoi(r.FormValue("live-bpm"))
	if err != nil {
		liveBPM = 0
	}

	minutes, err := strconv.Atoi(r.FormValue("minutes"))

	fmt.Println("****", minutes)
	if err != nil {
		log.Println("   Invalid entry for minutes: ", err)
	}

	seconds, err := strconv.Atoi(r.FormValue("seconds"))
	if err != nil {
		log.Println("   Invalid entry for seconds: ", err)
	}
	songLength := minutes*60 + seconds

	releaseDateString := r.FormValue("release-date")
	var releaseDate time.Time

	fmt.Println("   *** Release date: ", releaseDateString)
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

	err = database.SongsTableUpdateSong(song)
	if err != nil {
		http.Redirect(w, r, "/songs", http.StatusSeeOther)
	}

	data := models.SongPageData{
		User: user,
		Band: band,
		Song: song,
	}

	fmt.Println("$$$$$$$$$$$$$$$$$$$$$")

	err = h.Tmpl.ExecuteTemplate(w, "song.html", data)
	if err != nil {
		log.Println("Unable to execute song.html:", err)
		return
	}
}
