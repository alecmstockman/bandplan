package handlers

import (
	services "bandplan/src/Services"
	"bandplan/src/database"
	"bandplan/src/models"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
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
		log.Println("   Unable to get all songs by band id: ", err)
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
		log.Println("HandlerSend: Unable to get authenticated user: ", err)
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

	songs, err := database.SongsTableSearchByBandID(band.BandID, query)
	if err != nil {
		log.Println("   Error searching songs by Band ID: ", songs)
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

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "File too large", http.StatusBadRequest)
		return
	}

	artworkPath := ""
	imageID := ""

	file, _, err := r.FormFile("artwork-path")
	if err != nil {
		log.Println("   Error with provided artwork-path: ", err)
	} else {
		defer file.Close()

		imageID := uuid.New().String()

		artworkPath, err = HelperSaveArtworkImageVersions(file, imageID)
		if err != nil {
			http.Error(w, "Could not save artwork versions", http.StatusInternalServerError)
			return
		}
	}

	songTitle := strings.TrimSpace(r.FormValue("song-title"))
	if songTitle == "" {
		log.Println("   songTitle entry was only spaces")
		http.Redirect(w, r, "/songs/add", http.StatusSeeOther)
		return
	}

	artistName := strings.TrimSpace(r.FormValue("artist-name"))
	if songTitle == "" {
		log.Println("   artistName entry was only spaces")
		http.Redirect(w, r, "/songs/add", http.StatusSeeOther)
		return
	}

	albumTitle := strings.TrimSpace(r.FormValue("album-name"))
	genre := strings.TrimSpace(r.FormValue("genre"))
	musicalKey := strings.TrimSpace(r.FormValue("musical-key"))
	tuning := strings.TrimSpace(r.FormValue("tuning"))
	capo := strings.TrimSpace(r.FormValue("capo"))

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

	seconds, err := strconv.Atoi(r.FormValue("seconds"))
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

	status := r.FormValue("status")
	explicitEntry := r.FormValue("explicit")
	isCoverEntry := r.FormValue("is-cover")

	var explicit bool
	var isCover bool

	if explicitEntry == "on" {
		explicit = true
	} else {
		explicit = false
	}

	if isCoverEntry == "on" {
		isCover = true
	} else {
		isCover = false
	}

	spotifyLink := r.FormValue("spotify-link")
	appleMusicLink := r.FormValue("apple-music-link")
	youtubeLink := r.FormValue("youtube-link")
	amazonMusicLink := r.FormValue("amazon-music-link")
	pandoraLink := r.FormValue("pandora-link")
	deezerLink := r.FormValue("deezer-link")
	otherLink := r.FormValue("other-link")

	lyrics := r.FormValue("lyrics")
	description := r.FormValue("description")
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

		BandID: band.BandID,

		Title:      songTitle,
		TitleSlug:  "",
		ArtistName: artistName,
		AlbumTitle: albumTitle,
		AlbumSlug:  "",

		ArtworkID:   imageID,
		ArtworkPath: artworkPath,
		ReleaseDate: releaseDate,
		Genre:       genre,

		RecordingBPM:  recordingBPM,
		LiveBPM:       liveBPM,
		MusicalKey:    musicalKey,
		Tuning:        tuning,
		Capo:          capo,
		LengthSeconds: songLength,

		Status:   status,
		Explicit: explicit,
		IsCover:  isCover,

		SpotifyLink:     spotifyLink,
		AppleMusicLink:  appleMusicLink,
		YouTubeLink:     youtubeLink,
		AmazonMusicLink: amazonMusicLink,
		PandoraLink:     pandoraLink,
		DeezerLink:      deezerLink,
		OtherLink:       otherLink,

		Lyrics:      lyrics,
		Description: description,
		Notes:       notes,
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
	log.Println("- HandlerSongEditPage")

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
	log.Print("- HandlerSongUpdate")

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Println("   Unable to parse multipart form: ", err)
		http.Error(w, "File too large", http.StatusBadRequest)
		return
	}

	songID := r.FormValue("song-id")

	imageID, artworkPath, err := database.SongsTableGetImageIDAndPathBySongID(songID)
	if err != nil {
		log.Println("   Unable to get image ID: ", err)
	}

	file, _, err := r.FormFile("artwork-path")
	if err != nil {
		log.Println("   Error with provided artwork-path: ", err)
	} else {
		defer file.Close()

		newImageID := uuid.New().String()

		artworkPath, err = HelperSaveArtworkImageVersions(file, newImageID)
		if err != nil {
			http.Error(w, "Could not save artwork versions", http.StatusInternalServerError)
			return
		}
		err = HelperDeleteArtworkImageVersions(imageID)
		if err != nil {
			log.Println("   Unable to delete artwork image versions: ", err)
		}
		imageID = newImageID
	}

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

	if releaseDateString != "" {
		parsedDate, err := time.Parse("2006-01-02", releaseDateString)
		if err != nil {
			http.Error(w, "Invalid release date", http.StatusBadRequest)
			return
		}
		releaseDate = parsedDate
	}

	spotifyLink := strings.TrimSpace(r.FormValue("spotify-link"))
	appleMusicLink := strings.TrimSpace(r.FormValue("apple-music-link"))
	youtubeLink := strings.TrimSpace(r.FormValue("youtube-link"))
	amazonMusicLink := strings.TrimSpace(r.FormValue("amazon-music-link"))
	pandoraLink := strings.TrimSpace(r.FormValue("pandora-link"))
	deezerLink := strings.TrimSpace(r.FormValue("deezer-link"))
	otherLink := strings.TrimSpace(r.FormValue("other-link"))

	lyrics := strings.TrimSpace(r.FormValue("lyrics"))
	description := strings.TrimSpace(r.FormValue("description"))
	notes := strings.TrimSpace(r.FormValue("notes"))

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
		SongID:     songID,
		Title:      songTitle,
		AlbumTitle: albumTitle,
		BandID:     band.BandID,
		Genre:      genre,

		RecordingBPM:  recordingBPM,
		LiveBPM:       liveBPM,
		MusicalKey:    musicalKey,
		Tuning:        tuning,
		LengthSeconds: songLength,

		ReleaseDate: releaseDate,

		SpotifyLink:     spotifyLink,
		AppleMusicLink:  appleMusicLink,
		YouTubeLink:     youtubeLink,
		AmazonMusicLink: amazonMusicLink,
		PandoraLink:     pandoraLink,
		DeezerLink:      deezerLink,
		OtherLink:       otherLink,

		Lyrics:      lyrics,
		Description: description,
		Notes:       notes,
	}

	if artworkPath != "" {
		song.ArtworkID = imageID
		song.ArtworkPath = artworkPath

		err = database.SongsTableUpdateSong(song)
		if err != nil {
			http.Redirect(w, r, "/songs", http.StatusSeeOther)
		}
	} else {
		err = database.SongsTableUpdateSongWithoutArt(song)
		if err != nil {
			http.Redirect(w, r, "/songs", http.StatusSeeOther)
		}
	}

	updatedSong, err := database.SongsTableGetSongBySongID(songID)
	if err != nil {
		log.Println("   Unable to get updated song:", err)
		http.Redirect(w, r, "/songs", http.StatusSeeOther)
		return
	}

	data := models.SongPageData{
		User: user,
		Band: band,
		Song: updatedSong,
	}

	err = h.Tmpl.ExecuteTemplate(w, "song.html", data)
	if err != nil {
		log.Println("Unable to execute song.html:", err)
		return
	}
}

func (h Handler) HandlerSongDelete(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerSongDelete")

	songID := r.FormValue("song-id")

	if songID == "" {
		http.Error(w, "Missing song ID", http.StatusBadRequest)
		return
	}

	err := database.SongsTableDeleteSongByID(songID)
	if err != nil {
		log.Println("  Unable to delete song: ", err)
		http.Redirect(w, r, "/songs", http.StatusSeeOther)
	}

	http.Redirect(w, r, "/songs", http.StatusSeeOther)

}

func (h Handler) HandlerSongsITunesQueryPage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("------------------------------------------")
	log.Println("- HandlerSongsITunesSearch")

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

	fmt.Println("   band name: ", band.Name)

	data := models.SongDownloadData{
		User: user,
		Band: band,
	}

	err = h.Tmpl.ExecuteTemplate(w, "song-download.html", data)
	if err != nil {
		log.Println("   Err getting songs-download page: ", err)
		http.Redirect(w, r, "/songs", http.StatusSeeOther)
	}
}

func (h Handler) HandlerSongsITunesQuery(w http.ResponseWriter, r *http.Request) {
	fmt.Println("------------------------------------------")
	log.Println("- HandlerSongsITunesQuery")

	user, err := HelperGetAuthenticatedUser(r)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	_, err = database.BandsTableGetBandByUserID(user.UserID)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	artistQuery := strings.TrimSpace(r.FormValue("itunes-query-artist-name"))
	songQuery := strings.TrimSpace(r.FormValue("itunes-query-song-title"))

	fmt.Println("   artistQuery: ", artistQuery)
	fmt.Println("   songQuery; ", songQuery)

	searchResponse, err := services.ServicesSearchITunesByArtistAndSong(artistQuery, songQuery)
	if err != nil {
		log.Println("   Unable to get reponse from iTunes API: ", err)
	}

	fmt.Println("\nSearch Response: ")
	fmt.Println(searchResponse.ResultCount)
	fmt.Println("\nSearch Results")
	for _, r := range searchResponse.Results {
		fmt.Printf("%#v\n", r)
		fmt.Println("")
	}

	http.Redirect(w, r, "/songs", http.StatusSeeOther)
	return
}
