package handlers

import (
	"bandplan/src/database"
	"bandplan/src/helpers"
	requestlog "bandplan/src/logging"
	"bandplan/src/models"
	services "bandplan/src/services"
	"bytes"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

type songITunesQueryData struct {
	ArtistName string
	SongTitle  string
	Message    string
}

func (h Handler) HandlerSongsITunesQueryPage(w http.ResponseWriter, r *http.Request) {
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

	h.renderSongsITunesFragment(w, r, "song-download.html", songITunesQueryData{
		ArtistName: auth.CurrentBand.Name,
	})
}

func (h Handler) HandlerSongsITunesQuery(w http.ResponseWriter, r *http.Request) {
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

	query := songITunesQueryData{
		ArtistName: strings.TrimSpace(r.FormValue("itunes-query-artist-name")),
		SongTitle:  strings.TrimSpace(r.FormValue("itunes-query-song-title")),
	}

	search, err := services.ServicesSearchITunesByArtistAndSong(query.ArtistName, query.SongTitle)
	if err != nil {
		slog.Error(
			"unable to search iTunes",
			"request_id", requestlog.GetRequestID(r.Context()),
			"error", err,
		)
		query.Message = "iTunes is unavailable right now. Try again in a moment."
		h.renderSongsITunesFragment(w, r, "song-download.html", query)
		return
	}

	switch len(search.Results) {
	case 0:
		query.Message = "No songs found. Check the artist name or try a different song title."
		h.renderSongsITunesFragment(w, r, "song-download.html", query)
	case 1:
		res := search.Results[0]
		band := auth.CurrentBand
		releaseDate, err := time.Parse(time.RFC3339, "2024-02-01T12:00:00Z")
		if err != nil {
			log.Println("   unable to parse release date:", err)
		}

		h.renderSongsITunesFragment(w, r, "songs-add-itunes.html", models.SongImportFormData{
			Song: models.Song{
				BandID:         band.BandID,
				Title:          res.TrackName,
				ArtistName:     res.ArtistName,
				AlbumTitle:     res.CollectionName,
				ReleaseDate:    releaseDate,
				Genre:          res.PrimaryGenreName,
				LengthSeconds:  res.TrackTimeMillis / 1000,
				IsCover:        res.ArtistName != band.Name,
				AppleMusicLink: res.TrackViewURL,
			},
			ITunesTrackID: res.TrackID,
		})
	default:
		h.renderSongsITunesFragment(w, r, "songs-itunes-results.html", search)
	}
}

func (h Handler) renderSongsITunesFragment(w http.ResponseWriter, r *http.Request, name string, data any) {
	var fragment bytes.Buffer
	if err := h.Tmpl.ExecuteTemplate(&fragment, name, data); err != nil {
		slog.Error(
			"unable to render iTunes search",
			"request_id", requestlog.GetRequestID(r.Context()),
			"template", name,
			"error", err,
		)
		http.Error(w, "Unable to display iTunes search", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fragment.WriteTo(w)
}

func (h Handler) HandlerSongsITunesResults(w http.ResponseWriter, r *http.Request) {

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

	data := models.SongDownloadData{
		User: user,
		Band: band,
	}

	err = h.Tmpl.ExecuteTemplate(w, "songs-itunes-results.html", data)
	if err != nil {
		log.Println("   Err getting songs-download page: ", err)
		http.Redirect(w, r, "/songs", http.StatusSeeOther)
	}
}

func (h Handler) HandlerSongsITunesResultsAddSong(w http.ResponseWriter, r *http.Request) {

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

	r.Body = http.MaxBytesReader(w, r.Body, 16*1024)
	if err := r.ParseForm(); err != nil {
		var sizeErr *http.MaxBytesError
		if errors.As(err, &sizeErr) {
			http.Error(w, "Request too large", http.StatusRequestEntityTooLarge)
		} else {
			http.Error(w, "Invalid form", http.StatusBadRequest)
		}
		return
	}

	itunesID := r.FormValue("track-id")

	user := auth.User
	band := auth.CurrentBand

	res, err := services.ServicesSearchITunesByITunesID(itunesID)
	if err != nil {
		log.Println("   Unable to get song from iTunes by id:", err)
		http.Redirect(w, r, "/songs", http.StatusSeeOther)
		return
	}

	song := res.Results[0]

	length := song.TrackTimeMillis / 1000
	releaseDate, err := time.Parse(time.RFC3339, "2024-02-01T12:00:00Z")
	if err != nil {
		log.Println("   unable to parse release date:", err)
	}

	var cover bool
	if song.ArtistName == band.Name {
		cover = false
	} else {
		cover = true
	}

	imageID := uuid.New().String()
	// artworkPath, err := h.Services.ServiceSaveArtworkImageFromITunes(song.ArtworkURL100, imageID)
	// if err != nil {
	// 	log.Println("   Unable to save artwork image from iTunes: ", err)
	// 	imageID = ""
	// 	artworkPath = ""
	// }

	itunesTrackID := song.TrackID

	newSong := models.Song{
		SongID:        uuid.New().String(),
		BandID:        band.BandID,
		ITunesTrackID: &itunesTrackID,

		Title:      song.TrackName,
		TitleSlug:  helpers.MakeSlug(song.TrackName),
		ArtistName: song.ArtistName,
		AlbumTitle: song.CollectionName,
		AlbumSlug:  "",

		ArtworkID:     imageID,
		ArtworkPath:   song.ArtworkURL100,
		ReleaseDate:   releaseDate,
		Genre:         song.PrimaryGenreName,
		LengthSeconds: length,

		IsCover: cover,

		AppleMusicLink: song.TrackViewURL,
		CreatedBy:      user.UserID,
		UpdatedBy:      user.UserID,
	}

	_, err = database.SongsTableCreateSong(newSong)
	if err != nil {
		log.Println("   Unable to to create song in database: ", err)
		return
	}

	notif := models.Notification{
		BandID:    &band.BandID,
		UserID:    user.UserID,
		Type:      models.NotificationSongImported,
		ItemID:    newSong.SongID,
		Body:      models.NotifBodySongImported,
		CreatedBy: user.UserID,
	}

	_, err = h.Services.CreateBandNotification(r.Context(), notif)
	if err != nil {
		slog.Error(
			"unable to create notification",
			"request_id", requestlog.GetRequestID(r.Context()),
			"type", notif.Type,
			"error", err,
		)
	}

	return
}
