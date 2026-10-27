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
	"mime/multipart"
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
	data := models.SongDownloadData{
		CSRFToken: csrf.Token(r),
	}

	err := h.Tmpl.ExecuteTemplate(w, "song-create.html", data)
	if err != nil {
		slog.Error(
			"unable to load song-create.html",
			"request_id", requestlog.GetRequestID(r.Context()),
			"error", err,
		)
		http.Error(w, "unable to load create songs page", http.StatusInternalServerError)
		return
	}
}

func (h Handler) HandlerSongTempArt(w http.ResponseWriter, r *http.Request) {
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

	r.Body = http.MaxBytesReader(w, r.Body, 12<<20)
	if err = r.ParseMultipartForm(10 << 20); err != nil {
		var sizeErr *http.MaxBytesError
		if errors.As(err, &sizeErr) {
			http.Error(w, "Upload too large", http.StatusRequestEntityTooLarge)
			return
		}
		http.Error(w, "Invalid multipart form", http.StatusBadRequest)
		return
	}
	defer r.MultipartForm.RemoveAll()

	file, fileHeader, err := r.FormFile("artwork-path")
	if err != nil {
		http.Error(w, "Unable to read artwork", http.StatusBadRequest)
		return
	}
	defer file.Close()
	if fileHeader.Size > 10<<20 {
		http.Error(w, "Upload too large", http.StatusRequestEntityTooLarge)
		return
	}

	imageID := uuid.NewString()
	previewURL, err := h.Services.ServiceSaveTempImage(
		r.Context(),
		file,
		imageID,
		auth.CurrentBand.Slug,
		"song",
	)
	if err != nil {
		slog.Error(
			"unable to save temporary song artwork",
			"request_id", requestlog.GetRequestID(r.Context()),
			"error", err,
		)
		writeImageUploadError(w, err)
		return
	}

	previousImageID := strings.TrimSpace(r.FormValue("temporary-artwork-id"))
	if previousImageID != "" {
		if _, parseErr := uuid.Parse(previousImageID); parseErr == nil {
			if deleteErr := h.Services.ServiceDeleteTempSongImageVersions(
				r.Context(),
				previousImageID,
				auth.CurrentBand.Slug,
			); deleteErr != nil {
				slog.Error(
					"unable to delete replaced temporary song artwork",
					"request_id", requestlog.GetRequestID(r.Context()),
					"image_id", previousImageID,
					"error", deleteErr,
				)
			}
		}
	}

	data := models.ArtworkPreviewData{
		ArtworkID:  imageID,
		PreviewURL: previewURL,
		IsEdit: r.FormValue("artwork-edit") == "true" ||
			strings.TrimSpace(r.FormValue("song-id")) != "",
	}
	if err = h.Tmpl.ExecuteTemplate(w, "song_artwork_preview", data); err != nil {
		if deleteErr := h.Services.ServiceDeleteTempSongImageVersions(
			r.Context(),
			imageID,
			auth.CurrentBand.Slug,
		); deleteErr != nil {
			slog.Error(
				"unable to clean up temporary song artwork after render failure",
				"request_id", requestlog.GetRequestID(r.Context()),
				"image_id", imageID,
				"error", deleteErr,
			)
		}
		http.Error(w, "Unable to render artwork preview", http.StatusInternalServerError)
		return
	}
}

func (h Handler) HandlerSongTempArtDelete(w http.ResponseWriter, r *http.Request) {
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

	isEdit := r.FormValue("artwork-edit") == "true"
	imageID := strings.TrimSpace(r.FormValue("artwork-id"))
	if imageID != "" {
		if _, err = uuid.Parse(imageID); err != nil {
			http.Error(w, "Invalid artwork ID", http.StatusBadRequest)
			return
		}

		if err = h.Services.ServiceDeleteTempSongImageVersions(
			r.Context(),
			imageID,
			auth.CurrentBand.Slug,
		); err != nil {
			slog.Error(
				"unable to delete temporary song artwork",
				"request_id", requestlog.GetRequestID(r.Context()),
				"image_id", imageID,
				"error", err,
			)
			http.Error(w, "Unable to remove artwork", http.StatusInternalServerError)
			return
		}
	} else if !isEdit {
		http.Error(w, "Artwork ID is required", http.StatusBadRequest)
		return
	}

	data := models.ArtworkPreviewData{IsEdit: isEdit}
	if err = h.Tmpl.ExecuteTemplate(w, "song_artwork_reset", data); err != nil {
		http.Error(w, "Unable to render artwork upload", http.StatusInternalServerError)
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
	r2Artwork := false

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
	temporaryArtworkID := strings.TrimSpace(r.FormValue("temporary-artwork-id"))

	if temporaryArtworkID != "" {
		if _, err = uuid.Parse(temporaryArtworkID); err != nil {
			http.Error(w, "Invalid temporary artwork ID", http.StatusBadRequest)
			return
		}

		imageID = temporaryArtworkID
		artworkPath, err = h.Services.ServiceCreatePermSongImage(r.Context(), imageID, band.Slug)
		if err != nil {
			slog.Error(
				"unable to promote temporary song artwork",
				"request_id", requestlog.GetRequestID(r.Context()),
				"image_id", imageID,
				"error", err,
			)
			http.Error(w, "Unable to save artwork", http.StatusInternalServerError)
			return
		}
		r2Artwork = true
	} else {
		file, fileHeader, fileErr := r.FormFile("artwork-path")
		if errors.Is(fileErr, http.ErrMissingFile) {
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
		} else if fileErr != nil {
			slog.Error(
				"unable to read uploaded artwork",
				"request_id", requestlog.GetRequestID(r.Context()),
				"path", r.URL.Path,
				"error", fileErr,
			)
			http.Error(w, "Unable to read artwork", http.StatusBadRequest)
			return
		} else {
			defer file.Close()
			if fileHeader.Size > 10<<20 {
				http.Error(w, "Upload too large", http.StatusRequestEntityTooLarge)
				return
			}
			imageID = uuid.NewString()
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
			r2Artwork = true
		}
	}

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
		if r2Artwork {
			if deleteErr := h.Services.ServiceDeleteArtworkImageVersions(r.Context(), imageID, band.Slug); deleteErr != nil {
				slog.Error(
					"unable to remove artwork after song save failure",
					"request_id", requestlog.GetRequestID(r.Context()),
					"image_id", imageID,
					"error", deleteErr,
				)
			}
		}
		http.Error(w, "Unable to create song", http.StatusInternalServerError)
		return
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

	songID := strings.TrimSpace(r.FormValue("song-id"))
	if songID == "" {
		http.Error(w, "Invalid song data", http.StatusBadRequest)
		return
	}

	song, err := database.SongsTableGetSongBySongID(songID, band.BandID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.NotFound(w, r)
			return
		}
		slog.Error(
			"unable to load song",
			"request_id", requestlog.GetRequestID(r.Context()),
			"song_id", songID,
			"error", err,
		)
		http.Error(w, "Unable to load song", http.StatusInternalServerError)
		return
	}
	oldArtworkID := song.ArtworkID
	oldArtworkPath := song.ArtworkPath

	songTitle := strings.TrimSpace(r.FormValue("song-title"))
	if songTitle == "" {
		http.Error(w, "Song title is required", http.StatusBadRequest)
		return
	}

	artistName := strings.TrimSpace(r.FormValue("artist-name"))
	if artistName == "" {
		http.Error(w, "Artist name is required", http.StatusBadRequest)
		return
	}
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
		minutes = 0
	}

	seconds, err := strconv.Atoi(r.FormValue("seconds"))
	if err != nil {
		seconds = 0
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

	temporaryArtworkID := strings.TrimSpace(r.FormValue("temporary-artwork-id"))
	removeArtwork := r.FormValue("remove-artwork") == "true"
	if temporaryArtworkID != "" && removeArtwork {
		http.Error(w, "Conflicting artwork changes", http.StatusBadRequest)
		return
	}

	var directFile multipart.File
	var directFileHeader *multipart.FileHeader
	if temporaryArtworkID == "" && !removeArtwork {
		directFile, directFileHeader, err = r.FormFile("artwork-path")
		if err != nil && !errors.Is(err, http.ErrMissingFile) {
			http.Error(w, "Unable to read artwork", http.StatusBadRequest)
			return
		}
		if err == nil {
			defer directFile.Close()
			if directFileHeader.Size > 10<<20 {
				http.Error(w, "Upload too large", http.StatusRequestEntityTooLarge)
				return
			}
		}
	}

	artworkChanged := false
	newArtworkCreated := false
	if temporaryArtworkID != "" {
		if _, err := uuid.Parse(temporaryArtworkID); err != nil {
			http.Error(w, "Invalid temporary artwork ID", http.StatusBadRequest)
			return
		}

		artworkPath, err := h.Services.ServiceCreatePermSongImage(r.Context(), temporaryArtworkID, band.Slug)
		if err != nil {
			slog.Error(
				"unable to promote temporary song artwork",
				"request_id", requestlog.GetRequestID(r.Context()),
				"image_id", temporaryArtworkID,
				"error", err,
			)
			http.Error(w, "Unable to save artwork", http.StatusInternalServerError)
			return
		}
		song.ArtworkID = temporaryArtworkID
		song.ArtworkPath = artworkPath
		artworkChanged = true
		newArtworkCreated = true
	} else if removeArtwork {
		song.ArtworkID = ""
		song.ArtworkPath = ""
		artworkChanged = true
	} else if directFile != nil {
		newImageID := uuid.NewString()
		artworkPath, err := h.Services.ServiceSaveSongImageVersions(r.Context(), directFile, newImageID, band.Slug)
		if err != nil {
			slog.Error(
				"unable to save song artwork",
				"request_id", requestlog.GetRequestID(r.Context()),
				"image_id", newImageID,
				"error", err,
			)
			writeImageUploadError(w, err)
			return
		}
		song.ArtworkID = newImageID
		song.ArtworkPath = artworkPath
		artworkChanged = true
		newArtworkCreated = true
	}

	song.Title = songTitle
	song.AlbumTitle = albumTitle
	song.ArtistName = artistName
	song.ReleaseDate = releaseDate
	song.Genre = genre
	song.RecordingBPM = recordingBPM
	song.LiveBPM = liveBPM
	song.TimeSignature = timeSignature
	song.OriginalKey = originalKey
	song.LiveKey = liveKey
	song.Tuning = tuning
	song.Capo = capo
	song.LengthSeconds = songLength
	song.Status = status
	song.Explicit = isExplicit
	song.IsCover = isCover
	song.SpotifyLink = spotifyLink
	song.AppleMusicLink = appleMusicLink
	song.YouTubeLink = youtubeLink
	song.AmazonMusicLink = amazonMusicLink
	song.PandoraLink = pandoraLink
	song.DeezerLink = deezerLink
	song.TidalLink = tidalLink
	song.OtherLink = otherLink
	song.Lyrics = lyrics
	song.Description = description
	song.Notes = notes
	song.UpdatedBy = user.UserID

	err = database.SongsTableUpdateSong(song, band.BandID)
	if err != nil {
		if newArtworkCreated {
			if deleteErr := h.Services.ServiceDeleteArtworkImageVersions(r.Context(), song.ArtworkID, band.Slug); deleteErr != nil {
				slog.Error(
					"unable to remove replacement artwork after song update failure",
					"request_id", requestlog.GetRequestID(r.Context()),
					"image_id", song.ArtworkID,
					"error", deleteErr,
				)
			}
		}
		if errors.Is(err, sql.ErrNoRows) {
			http.NotFound(w, r)
			return
		}
		slog.Error(
			"unable to update song",
			"request_id", requestlog.GetRequestID(r.Context()),
			"song_id", songID,
			"error", err,
		)
		http.Error(w, "Unable to update song", http.StatusInternalServerError)
		return
	}

	if artworkChanged && oldArtworkID != "" {
		if err := h.Services.ServiceDeleteSongArtwork(r.Context(), oldArtworkID, oldArtworkPath, band); err != nil {
			slog.Error(
				"unable to delete replaced song artwork",
				"request_id", requestlog.GetRequestID(r.Context()),
				"song_id", songID,
				"image_id", oldArtworkID,
				"error", err,
			)
		}
	}

	http.Redirect(w, r, "/song?id="+url.QueryEscape(songID), http.StatusSeeOther)
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
