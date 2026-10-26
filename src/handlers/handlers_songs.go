package handlers

import (
	"bandplan/src/database"
	requestlog "bandplan/src/logging"
	"bandplan/src/models"
	"bandplan/src/services"
	"database/sql"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/csrf"
)

func (h Handler) HandlerSongsPage(w http.ResponseWriter, r *http.Request) {
	auth, err := HelperGetAuthContext(r)
	if err != nil {
		slog.Error(
			"unable to load auth context",
			"request_id", requestlog.GetRequestID(r.Context()),
			"error", err,
		)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	user := auth.User
	band := auth.CurrentBand

	data, err := h.Services.SongsPage(r.Context(), auth.User, auth.CurrentBand)
	if err != nil {
		slog.Error(
			"failed to load songs",
			"request_id", requestlog.GetRequestID(r.Context()),
			"path", r.URL.Path,
			"error", err,
		)
		http.Error(w, "Could not load song page", http.StatusInternalServerError)
		return
	}

	notifications, err := database.NotifTableGetNotificationsByUserID(r.Context(), user.UserID, band.BandID)
	if err != nil {
		slog.Error(
			"unable to get notifications by user id",
			"request_id", requestlog.GetRequestID(r.Context()),
			"path", r.URL.Path,
			"error", err,
		)
	}

	data.Notifications = notifications

	err = h.Tmpl.ExecuteTemplate(w, "songs.html", data)
	if err != nil {
		slog.Error(
			"failed to load songs.html",
			"request_id", requestlog.GetRequestID(r.Context()),
			"error", err,
		)
		http.Error(w, "Could not load song page", http.StatusInternalServerError)
		return
	}
}

func (h Handler) HandlerSongCreatePage(w http.ResponseWriter, r *http.Request) {
	err := h.Tmpl.ExecuteTemplate(w, "songs-add.html", nil)
	if err != nil {
		slog.Error(
			"unable to load songs-add.html",
			"request_id", requestlog.GetRequestID(r.Context()),
			"error", err,
		)
		http.Error(w, "unable to load create songs page", http.StatusInternalServerError)
		return
	}
}

func (h Handler) HandlerSongSave(w http.ResponseWriter, r *http.Request) {

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		slog.Error(
			"unable to load auth context",
			"request_id", requestlog.GetRequestID(r.Context()),
			"error", err,
		)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	user := auth.User
	band := auth.CurrentBand

	r.Body = http.MaxBytesReader(w, r.Body, 12<<20)

	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		var sizeErr *http.MaxBytesError
		if errors.As(err, &sizeErr) {
			http.Error(w, "Upload too large", http.StatusRequestEntityTooLarge)
			return
		}

		http.Error(w, "Invalid multipart form", http.StatusBadRequest)
		return
	}
	defer r.MultipartForm.RemoveAll()

	artworkPath := ""
	imageID := ""

	songTitle := strings.TrimSpace(r.FormValue("song-title"))
	if songTitle == "" {
		slog.Error(
			"no song title provided",
			"request_id", requestlog.GetRequestID(r.Context()),
		)
		http.Redirect(w, r, "/song/create", http.StatusSeeOther)
		return
	}

	artistName := strings.TrimSpace(r.FormValue("artist-name"))
	if artistName == "" {
		http.Error(w, "Invalid artist name entry", http.StatusBadRequest)
		return
	}

	albumTitle := strings.TrimSpace(r.FormValue("album-name"))
	trackIDValue := strings.TrimSpace(r.FormValue("itunes-track-id"))
	var itunesTrackID *int64
	if trackIDValue != "" {
		trackID, parseErr := strconv.ParseInt(trackIDValue, 10, 64)
		if parseErr != nil || trackID <= 0 {
			http.Error(w, "Invalid iTunes track ID", http.StatusBadRequest)
			return
		}
		itunesTrackID = &trackID
	}

	file, _, err := r.FormFile("artwork-path")
	if errors.Is(err, http.ErrMissingFile) {
		if itunesTrackID != nil {
			searchResult, lookupErr := services.ServicesSearchITunesByITunesID(trackIDValue)
			if lookupErr != nil {
				slog.Error(
					"unable to load iTunes artwork source",
					"request_id", requestlog.GetRequestID(r.Context()),
					"track_id", *itunesTrackID,
					"error", lookupErr,
				)
				http.Error(w, "Unable to validate artwork source", http.StatusBadGateway)
				return
			}
			if len(searchResult.Results) != 1 || searchResult.Results[0].TrackID != *itunesTrackID {
				http.Error(w, "Invalid artwork source", http.StatusBadRequest)
				return
			}

			imageID = uuid.NewString()
			artworkPath, err = h.Services.ServiceSaveArtworkImageFromITunes(
				searchResult.Results[0].ArtworkURL100,
				imageID,
			)
			if err != nil {
				slog.Error(
					"unable to save iTunes artwork",
					"request_id", requestlog.GetRequestID(r.Context()),
					"track_id", *itunesTrackID,
					"path", r.URL.Path,
					"error", err,
				)
				http.Error(w, "Unable to save artwork", http.StatusBadGateway)
				return
			}
		}

	} else if err != nil {
		slog.Error(
			"unable to read uploaded artwork",
			"request_id", requestlog.GetRequestID(r.Context()),
			"path", r.URL.Path,
			"error", err,
		)
		http.Error(w, "Unable to read artwork", http.StatusBadRequest)
		return

	} else {
		defer file.Close()

		imageID = uuid.New().String()

		artworkPath, err = h.Services.ServiceSaveSongImageVersions(r.Context(), file, imageID, band.Slug)
		if err != nil {
			slog.Error(
				"unable to save artwork image versions",
				"request_id", requestlog.GetRequestID(r.Context()),
				"path", r.URL.Path,
				"error", err,
			)
			writeImageUploadError(w, err)
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
	notes := r.FormValue("notes")

	song := models.Song{

		SongID:        uuid.New().String(),
		BandID:        band.BandID,
		ITunesTrackID: itunesTrackID,

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
		slog.Error(
			"unable to save song",
			"request_id", requestlog.GetRequestID(r.Context()),
			"path", r.URL.Path,
			"error", err,
		)
		http.Error(w, "Unable to create song", http.StatusInternalServerError)
	}

	notif := models.Notification{
		BandID:    &band.BandID,
		UserID:    user.UserID,
		Type:      models.NotificationSongCreated,
		ItemID:    song.SongID,
		Body:      models.NotifBodySongCreated,
		CreatedBy: user.UserID,
	}

	err = h.Services.CreateNotification(r.Context(), notif)
	if err != nil {
		slog.Error(
			"unable to create notification",
			"request_id", requestlog.GetRequestID(r.Context()),
			"type", notif.Type,
			"error", err,
		)
	}

	http.Redirect(w, r, "/songs", http.StatusSeeOther)
}

func (h Handler) HandlerSongPage(w http.ResponseWriter, r *http.Request) {

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		slog.Error(
			"unable to load auth context",
			"request_id", requestlog.GetRequestID(r.Context()),
			"error", err,
		)
		http.Error(w, "unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	songID := r.URL.Query().Get("id")
	setlistID := r.URL.Query().Get("setlist-id")

	backURL := "/songs"
	if setlistID != "" {
		backURL = "/setlist?id=" + url.QueryEscape(setlistID)
	}
	user := auth.User
	band := auth.CurrentBand

	data, err := h.Services.SongPage(r.Context(), user, band, songID)
	if err != nil {
		slog.Error(
			"error with services.SongPage",
			"request_id", requestlog.GetRequestID(r.Context()),
			"path", r.URL.Path,
			"error", err,
		)
	}

	data.BackURL = backURL

	err = h.Tmpl.ExecuteTemplate(w, "song.html", data)
	if err != nil {
		log.Println("Unable to execute song.html:", err)
		return
	}
}

func (h Handler) HandlerSongEditPage(w http.ResponseWriter, r *http.Request) {

	songID := r.URL.Query().Get("id")

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		slog.Error(
			"unable to load auth context",
			"request_id", requestlog.GetRequestID(r.Context()),
			"error", err,
		)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	user := auth.User
	band := auth.CurrentBand

	data, err := h.Services.SongData(r.Context(), user, band, songID)
	if err != nil {
		slog.Error(
			"unable to load song data",
			"request_id", requestlog.GetRequestID(r.Context()),
			"path", r.URL.Path,
			"error", err,
		)
		http.Error(w, "Unable to load song", http.StatusInternalServerError)
		return
	}

	data.CSRFToken = csrf.Token(r)

	err = h.Tmpl.ExecuteTemplate(w, "song-edit.html", data)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
}

func (h Handler) HandlerSongUpdate(w http.ResponseWriter, r *http.Request) {

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		slog.Error(
			"unable to load auth context",
			"request_id", requestlog.GetRequestID(r.Context()),
			"error", err,
		)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	user := auth.User
	band := auth.CurrentBand

	r.Body = http.MaxBytesReader(w, r.Body, 12<<20)

	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		var sizeErr *http.MaxBytesError
		if errors.As(err, &sizeErr) {
			http.Error(w, "Upload too large", http.StatusRequestEntityTooLarge)
			return
		}

		http.Error(w, "Invalid multipart form", http.StatusBadRequest)
		return
	}
	defer r.MultipartForm.RemoveAll()

	songID := r.FormValue("song-id")
	if songID == "" {
		slog.Error(
			"no song-id provided",
			"request_id", requestlog.GetRequestID(r.Context()),
		)
		http.Error(w, "Invalid song data", http.StatusBadRequest)
		return
	}

	imageID, artworkPath, err := database.SongsTableGetImageIDAndPathBySongID(songID, band.BandID, user.UserID)
	if err != nil {
		slog.Error(
			"unable to load image id",
			"request_id", requestlog.GetRequestID(r.Context()),
			"path", r.URL.Path,
			"error", err,
		)
	}

	file, _, err := r.FormFile("artwork-path")
	if err != nil {
		log.Println("   Error with provided artwork-path: ", err)
	} else {
		defer file.Close()

		newImageID := uuid.New().String()

		artworkPath, err = h.Services.ServiceSaveSongImageVersions(r.Context(), file, newImageID, band.Slug)
		if err != nil {
			slog.Error(
				"unable to save song images",
				"request_id", requestlog.GetRequestID(r.Context()),
				"path", r.URL.Path,
				"error", err,
			)
			writeImageUploadError(w, err)
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

		err = database.SongsTableUpdateSong(song, band.BandID)
		if err != nil {
			log.Println("   Unable to udpate song in database: ", err)
			http.Redirect(w, r, "/songs", http.StatusSeeOther)
		}
	} else {
		err = database.SongsTableUpdateSongWithoutArt(song, band.BandID)
		if err != nil {
			log.Println("   Error updating songs table without art: ", err)
			http.Redirect(w, r, "/songs", http.StatusSeeOther)
		}
	}

	redirectURL := "/song?id=" + url.QueryEscape(songID)

	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

func (h Handler) HandlerSongLyrics(w http.ResponseWriter, r *http.Request) {

	songID := r.FormValue("song-id")

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		slog.Error(
			"unable to load auth context",
			"request_id", requestlog.GetRequestID(r.Context()),
			"error", err,
		)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	user := auth.User
	band := auth.CurrentBand

	data, err := h.Services.SongData(r.Context(), user, band, songID)
	if err != nil {
		slog.Error(
			"unable to load auth context",
			"request_id", requestlog.GetRequestID(r.Context()),
			"path", r.URL.Path,
			"error", err,
		)
		http.Error(w, "Unable to load song", http.StatusSeeOther)
	}

	data.CSRFToken = csrf.Token(r)

	err = h.Tmpl.ExecuteTemplate(w, "lyrics.html", data.Song)
	if err != nil {
		slog.Error(
			"unable to load lyrics.html",
			"request_id", requestlog.GetRequestID(r.Context()),
			"path", r.URL.Path,
			"error", err,
		)
		http.Error(w, "Unable to load lyrics page", http.StatusInternalServerError)
	}
	return
}

func (h Handler) HandlerSongDelete(w http.ResponseWriter, r *http.Request) {

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		slog.Error(
			"unable to load auth context",
			"request_id", requestlog.GetRequestID(r.Context()),
			"error", err,
		)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	user := auth.User
	band := auth.CurrentBand

	songID := strings.TrimSpace(r.FormValue("song-id"))

	if songID == "" {
		http.Error(w, "Missing song ID", http.StatusBadRequest)
		return
	}

	err = h.Services.SongDelete(r.Context(), user, band, songID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.NotFound(w, r)
			return
		}
		slog.Error(
			"unable to delete song",
			"request_id", requestlog.GetRequestID(r.Context()),
			"path", r.URL.Path,
			"error", err,
		)
		http.Error(w, "Unable to delete song", http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Redirect", "/songs?toast=song-deleted")
	w.WriteHeader(http.StatusOK)
	return

	// http.Redirect(w, r, "/songs?toast=song-deleted", http.StatusSeeOther)
}

func (h *Handler) HandlerSongsSearch(w http.ResponseWriter, r *http.Request) {

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		slog.Error(
			"unable to load auth context",
			"request_id", requestlog.GetRequestID(r.Context()),
			"error", err,
		)
		http.Error(w, "unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	user := auth.User
	band := auth.CurrentBand

	query := r.FormValue("q")

	songs, err := database.SongsTableSearchByBandID(band.BandID, query)
	if err != nil {
		slog.Error(
			"unable to search songs by band ID",
			"request_id", requestlog.GetRequestID(r.Context()),
			"error", err,
		)
		http.Error(w, "unable to get songs", http.StatusInternalServerError)
		return
	}

	data := models.MenuPageData{
		User:  user,
		Band:  band,
		Songs: songs,
	}
	err = h.Tmpl.ExecuteTemplate(w, "songs-list.html", data)
	if err != nil {
		slog.Error(
			"unable to execute songs-list.html template",
			"request_id", requestlog.GetRequestID(r.Context()),
			"template", "songs-list.html",
			"error", err,
		)
		http.Error(w, "Unable to load songs-list templage", http.StatusInternalServerError)
		return
	}
}
