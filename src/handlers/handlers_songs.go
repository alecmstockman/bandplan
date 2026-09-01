package handlers

import (
	"bandplan/src/database"
	"bandplan/src/models"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

func (h Handler) HandlerSongsPage(w http.ResponseWriter, r *http.Request) {
	log.Print("- HandlerSongsPage")

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		log.Println("   Unable to get AuthContext: ", err)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	user := auth.User
	band := auth.CurrentBand

	songs, err := database.SongsTableGetAllSongsByBandID(band.BandID)
	if err != nil {
		log.Println("   Unable to get all songs by band id: ", err)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	setlists, err := database.SetlistsTableGetSetlistsByBandID(band.BandID)
	if err != nil {
		http.Error(w, "Could not get setlists by bandID", http.StatusInternalServerError)
		return
	}

	data := models.MenuPageData{
		User:     user,
		Band:     band,
		Songs:    songs,
		Setlists: setlists,
	}

	err = h.Tmpl.ExecuteTemplate(w, "songs.html", data)
	if err != nil {
		log.Println("   err getting songs.html: ", err)
		return
	}
}

func (h *Handler) HandlerSongsSearch(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerSongsSearch")

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		log.Println("   Unable to get AuthContext: ", err)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	user := auth.User
	band := auth.CurrentBand

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

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		log.Println("   Unable to get AuthContext: ", err)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	user := auth.User
	band := auth.CurrentBand

	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Println("   File too large: ", err)
		http.Error(w, "File too large", http.StatusBadRequest)
		return
	}

	artworkPath := ""
	imageID := ""

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

	file, _, err := r.FormFile("artwork-path")
	if err != nil {
		log.Println("   Error with provided artwork-path: ", err)
		imageID = r.FormValue("existing-artwork-id")

		artworkPath = r.FormValue("existing-artwork-path")

	} else {
		defer file.Close()

		imageID = uuid.New().String()

		artworkPath, err = h.Services.ServiceSaveArtworkImageVersions(r.Context(), file, imageID, band.Slug)
		if err != nil {
			http.Error(w, "Could not save artwork versions", http.StatusInternalServerError)
			return
		}
	}

	genre := strings.TrimSpace(r.FormValue("genre"))
	originalKey := strings.TrimSpace(r.FormValue("original-key"))
	liveKey := strings.TrimSpace(r.FormValue("live-key"))
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

	timeSignature := strings.TrimSpace(r.FormValue("time-signature"))

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
	tidalLink := r.FormValue("tidal-link")
	otherLink := r.FormValue("other-link")

	lyrics := r.FormValue("lyrics")
	description := r.FormValue("description")
	notes := r.FormValue("notes")

	song := models.Song{

		BandID: band.BandID,

		Title:     songTitle,
		TitleSlug: "",

		AlbumTitle: albumTitle,
		AlbumSlug:  "",

		ArtistName:  artistName,
		ArtworkID:   imageID,
		ArtworkPath: artworkPath,
		ReleaseDate: releaseDate,
		Genre:       genre,

		RecordingBPM:  recordingBPM,
		LiveBPM:       liveBPM,
		TimeSignature: timeSignature,
		OriginalKey:   originalKey,
		LiveKey:       liveKey,
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
		TidalLink:       tidalLink,
		OtherLink:       otherLink,

		Lyrics:      lyrics,
		Description: description,
		Notes:       notes,
		CreatedBy:   user.UserID,
		UpdatedBy:   user.UserID,
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

	song, err := database.SongsTableGetSongBySongID(songID)
	if err != nil {
		http.Error(w, "Could not get song", http.StatusInternalServerError)
		return
	}
	setlists, err := database.SetlistsTableGetSetlistsByBandID(band.BandID)
	if err != nil {
		http.Error(w, "Could not get setlists by bandID", http.StatusInternalServerError)
		return
	}

	data := models.SongPageData{
		BackURL:  backURL,
		User:     user,
		Band:     band,
		Song:     song,
		Setlists: setlists,
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

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		log.Println("   Unable to get AuthContext: ", err)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	user := auth.User
	band := auth.CurrentBand

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

func (h Handler) HandlerSongLyrics(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerSongLyrics")

	songID := r.FormValue("song-id")

	song, err := database.SongsTableGetSongBySongID(songID)
	if err != nil {
		log.Printf("   Unable to get song %s from database: %v", songID, err)
		return
	}

	err = h.Tmpl.ExecuteTemplate(w, "lyrics.html", song)
	if err != nil {
		log.Println("   Unable to execute lyrics.html")
		http.Error(w, "Unable to load lyrics page", http.StatusInternalServerError)
	}

	return
}

func (h Handler) HandlerSongUpdate(w http.ResponseWriter, r *http.Request) {
	log.Print("- HandlerSongUpdate")

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		log.Println("   Unable to get AuthContext: ", err)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	user := auth.User
	band := auth.CurrentBand

	err = r.ParseMultipartForm(10 << 20)
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

		artworkPath, err = h.Services.ServiceSaveArtworkImageVersions(r.Context(), file, newImageID, band.Slug)
		if err != nil {
			http.Error(w, "Could not save artwork versions", http.StatusInternalServerError)
			return
		}

		err = h.Services.ServiceDeleteArtworkImageVersions(r.Context(), imageID, band.Slug)
		if err != nil {
			log.Println("   Unable to delete artwork image versions: ", err)
		}
		imageID = newImageID
	}

	songTitle := strings.TrimSpace(r.FormValue("song-title"))
	if songTitle == "" {
		log.Print("   songTitle entry was only spaces")
		redirectURL := "/song/edit?id=" + url.QueryEscape(songID)
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
		return
	}

	artistName := strings.TrimSpace(r.FormValue("artist-name"))
	albumTitle := strings.TrimSpace(r.FormValue("album-name"))
	genre := strings.TrimSpace(r.FormValue("genre"))
	originalKey := strings.TrimSpace(r.FormValue("original-key"))
	liveKey := strings.TrimSpace(r.FormValue("live-key"))
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

	timeSignature := strings.TrimSpace(r.FormValue("time-signature"))

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

	status := strings.TrimSpace(r.FormValue("status"))

	isExplicit := r.FormValue("explicit") == "on"
	isCover := r.FormValue("is-cover") == "on"

	spotifyLink := strings.TrimSpace(r.FormValue("spotify-link"))
	appleMusicLink := strings.TrimSpace(r.FormValue("apple-music-link"))
	youtubeLink := strings.TrimSpace(r.FormValue("youtube-link"))
	amazonMusicLink := strings.TrimSpace(r.FormValue("amazon-music-link"))
	pandoraLink := strings.TrimSpace(r.FormValue("pandora-link"))
	deezerLink := strings.TrimSpace(r.FormValue("deezer-link"))
	tidalLink := strings.TrimSpace(r.FormValue("tidal-link"))
	otherLink := strings.TrimSpace(r.FormValue("other-link"))

	lyrics := strings.TrimSpace(r.FormValue("lyrics"))
	description := strings.TrimSpace(r.FormValue("description"))
	notes := strings.TrimSpace(r.FormValue("notes"))

	song := models.Song{
		SongID:     songID,
		Title:      songTitle,
		AlbumTitle: albumTitle,
		BandID:     band.BandID,
		ArtistName: artistName,

		ReleaseDate: releaseDate,
		Genre:       genre,

		RecordingBPM:  recordingBPM,
		LiveBPM:       liveBPM,
		TimeSignature: timeSignature,
		OriginalKey:   originalKey,
		LiveKey:       liveKey,
		Tuning:        tuning,
		Capo:          capo,
		LengthSeconds: songLength,

		Status:   status,
		Explicit: isExplicit,
		IsCover:  isCover,

		SpotifyLink:     spotifyLink,
		AppleMusicLink:  appleMusicLink,
		YouTubeLink:     youtubeLink,
		AmazonMusicLink: amazonMusicLink,
		PandoraLink:     pandoraLink,
		DeezerLink:      deezerLink,
		TidalLink:       tidalLink,
		OtherLink:       otherLink,

		Lyrics:      lyrics,
		Description: description,
		Notes:       notes,
		UpdatedBy:   user.UserID,
	}

	if artworkPath != "" {
		song.ArtworkID = imageID
		song.ArtworkPath = artworkPath

		err = database.SongsTableUpdateSong(song)
		if err != nil {
			log.Println("   Unable to udpate song in database: ", err)
			http.Redirect(w, r, "/songs", http.StatusSeeOther)
		}
	} else {
		err = database.SongsTableUpdateSongWithoutArt(song)
		if err != nil {
			log.Println("   Error updating songs table without art: ", err)
			http.Redirect(w, r, "/songs", http.StatusSeeOther)
		}
	}

	redirectURL := "/song?id=" + url.QueryEscape(songID)

	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

func (h Handler) HandlerSongDelete(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerSongDelete")

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		log.Println("   Unable to get AuthContext: ", err)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	band := auth.CurrentBand

	songID := r.FormValue("song-id")
	imageID := r.FormValue("artwork-id")

	log.Printf("Song ID: %q", songID)
	log.Printf("Artwork ID: %q", imageID)

	if songID == "" {
		http.Error(w, "Missing song ID", http.StatusBadRequest)
		return
	}

	err = h.Services.ServiceDeleteArtworkImageVersions(r.Context(), imageID, band.Slug)
	if err != nil {
		log.Println("   Unable to delete artwork image versions: ", err)
	}

	err = database.SongsTableDeleteSongByID(songID)
	if err != nil {
		log.Println("  Unable to delete song: ", err)
		http.Redirect(w, r, "/songs", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/songs", http.StatusSeeOther)

}
