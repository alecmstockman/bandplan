package handlers

import (
	"bandplan/src/database"
	requestlog "bandplan/src/logging"
	"bandplan/src/models"
	"database/sql"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/uuid"
)

func (h Handler) HandlerBandImageAdd(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

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
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		var sizeErr *http.MaxBytesError
		if errors.As(err, &sizeErr) {
			http.Error(w, "Upload too large", http.StatusRequestEntityTooLarge)
			return
		}

		http.Error(w, "Invalid multipart form", http.StatusBadRequest)
		return
	}
	defer r.MultipartForm.RemoveAll()

	bandID := strings.TrimSpace(r.FormValue("band-id"))
	if bandID == "" {
		http.Error(w, "Band is required", http.StatusBadRequest)
		return
	}

	band, err := database.BandsTableGetBandByBandIDAndUserID(bandID, auth.User.UserID)
	if errors.Is(err, sql.ErrNoRows) {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		slog.Error(
			"unable to load band for profile image upload",
			"request_id", requestlog.GetRequestID(r.Context()),
			"user_id", auth.User.UserID,
			"band_id", bandID,
			"error", err,
		)
		http.Error(w, "Unable to load band", http.StatusInternalServerError)
		return
	}

	isAdmin, err := database.BandMembersGetAdminStatus(auth.User.UserID, band.BandID)
	if err != nil {
		slog.Error(
			"unable to load band admin status",
			"request_id", requestlog.GetRequestID(r.Context()),
			"user_id", auth.User.UserID,
			"band_id", band.BandID,
			"error", err,
		)
		http.Error(w, "Unable to verify band access", http.StatusInternalServerError)
		return
	}
	if !isAdmin {
		http.Error(w, "Band admin access required", http.StatusForbidden)
		return
	}

	file, _, err := r.FormFile("band-image")
	if errors.Is(err, http.ErrMissingFile) {
		http.Error(w, "Band image is required", http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, "Unable to read band image", http.StatusBadRequest)
		return
	}
	defer file.Close()

	imageID := uuid.NewString()
	browserPath, err := h.Services.ServiceSaveBandImageVersions(r.Context(), file, imageID, band.Slug)
	if err != nil {
		slog.Error(
			"unable to save band profile images",
			"request_id", requestlog.GetRequestID(r.Context()),
			"band_id", band.BandID,
			"error", err,
		)
		writeImageUploadError(w, err)
		return
	}

	if err := database.BandsTableUpdateProfileImage(band.BandID, imageID, browserPath); err != nil {
		slog.Error(
			"unable to update band profile image",
			"request_id", requestlog.GetRequestID(r.Context()),
			"band_id", band.BandID,
			"error", err,
		)
		http.Error(w, "Unable to update band image", http.StatusInternalServerError)
		return
	}

	if band.ProfileImageID != "" {
		if err := h.Services.ServiceDeleteBandImageVersions(r.Context(), band.ProfileImageID, band.Slug); err != nil {
			slog.Error(
				"unable to delete old band profile images",
				"request_id", requestlog.GetRequestID(r.Context()),
				"band_id", band.BandID,
				"image_id", band.ProfileImageID,
				"error", err,
			)
		}
	}

	http.Redirect(w, r, "/band?band-id="+url.QueryEscape(band.BandID), http.StatusSeeOther)
}

func (h Handler) HandlerBandTempArt(w http.ResponseWriter, r *http.Request) {

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

	band := auth.CurrentBand

	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Println("   Error parsing from while creating a new setlist: ", err)
		http.Error(w, "Unable to parse from", http.StatusBadRequest)
		return
	}

	imageID := ""
	previewURL := ""

	file, _, err := r.FormFile("artwork-path")
	if err != nil {
		if err != http.ErrMissingFile {
			log.Println("   Unable to read artwork file:", err)
			http.Error(w, "Unable to read artwork", http.StatusBadRequest)
			return
		}

		log.Println("   No setlist artwork uploaded")
	} else {
		defer file.Close()

		imageID = uuid.New().String()

		previewURL, err = h.Services.ServiceSaveTempImage(r.Context(), file, imageID, band.Slug, "setlist")
		if err != nil {
			slog.Error(
				"unable to save temp image",
				"request_id", requestlog.GetRequestID(r.Context()),
				"path", r.URL.Path,
				"error", err,
			)
			writeImageUploadError(w, err)
			return
		}
	}

	data := models.ArtworkPreviewData{
		ArtworkID:  imageID,
		PreviewURL: previewURL,
	}

	err = h.Tmpl.ExecuteTemplate(w, "setlist_artwork_preview", data)
	if err != nil {
		http.Error(w, "Unable to render preview", http.StatusInternalServerError)
		return
	}

}

func (h Handler) HandlerBandTempArtDelete(w http.ResponseWriter, r *http.Request) {

	_, err := HelperGetAuthContext(r)
	if err != nil {
		slog.Error(
			"unable to load auth context",
			"request_id", requestlog.GetRequestID(r.Context()),
			"error", err,
		)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

}
