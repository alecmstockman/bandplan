package handlers

import (
	"bandplan/src/database"
	"bandplan/src/helpers"
	requestlog "bandplan/src/logging"
	"bandplan/src/models"
	"database/sql"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/csrf"
)

func (h Handler) HandlerProfilePage(w http.ResponseWriter, r *http.Request) {

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

	data := models.MenuPageData{
		CSRFToken: csrf.Token(r),
		User:      user,
		Band:      band,
	}

	err = h.Tmpl.ExecuteTemplate(w, "profile.html", data)
	if err != nil {
		log.Println("   Err getting profile pic page: ", err)
		return
	}
}

func (h Handler) HandlerBandsPage(w http.ResponseWriter, r *http.Request) {
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

	bands, err := database.BandsTableGetBandsByUserID(auth.User.UserID)
	if err != nil {
		slog.Error(
			"unable to load user bands",
			"request_id", requestlog.GetRequestID(r.Context()),
			"user_id", auth.User.UserID,
			"error", err,
		)
		http.Error(w, "Unable to load bands", http.StatusInternalServerError)
		return
	}

	data := models.BandsPageData{
		CSRFToken: csrf.Token(r),
		User:      auth.User,
		Band:      auth.CurrentBand,
		Bands:     bands,
	}

	if err = h.Tmpl.ExecuteTemplate(w, "bands.html", data); err != nil {
		slog.Error(
			"unable to render bands page",
			"request_id", requestlog.GetRequestID(r.Context()),
			"error", err,
		)
	}
}

func (h Handler) HandlerBandPage(w http.ResponseWriter, r *http.Request) {
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

	bandID := strings.TrimSpace(r.URL.Query().Get("band-id"))
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
			"unable to load band",
			"request_id", requestlog.GetRequestID(r.Context()),
			"user_id", auth.User.UserID,
			"band_id", bandID,
			"error", err,
		)
		http.Error(w, "Unable to load band", http.StatusInternalServerError)
		return
	}

	data := models.BandPageData{
		User:         auth.User,
		Band:         auth.CurrentBand,
		SelectedBand: band,
	}

	if err = h.Tmpl.ExecuteTemplate(w, "band.html", data); err != nil {
		slog.Error(
			"unable to render band page",
			"request_id", requestlog.GetRequestID(r.Context()),
			"band_id", bandID,
			"error", err,
		)
	}
}

func (h Handler) HandlerBandSwitch(w http.ResponseWriter, r *http.Request) {
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

	bandID := strings.TrimSpace(r.FormValue("band-id"))
	if bandID == "" {
		http.Error(w, "Band is required", http.StatusBadRequest)
		return
	}

	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, "Unable to load session", http.StatusUnauthorized)
		return
	}

	tokenHash := helpers.HashSessionToken(cookie.Value)
	updated, err := database.SessionsTableSetCurrentBand(tokenHash, auth.User.UserID, bandID)
	if err != nil {
		slog.Error(
			"unable to switch current band",
			"request_id", requestlog.GetRequestID(r.Context()),
			"user_id", auth.User.UserID,
			"band_id", bandID,
			"error", err,
		)
		http.Error(w, "Unable to switch band", http.StatusInternalServerError)
		return
	}
	if !updated {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, "/bands", http.StatusSeeOther)
}

func (h Handler) HandlerProfilePicAdd(w http.ResponseWriter, r *http.Request) {

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

	user := auth.User

	oldImageID := user.ProfileImageID

	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Println("   Error parsing multipart form: ", err)
		http.Error(w, "File too large", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("profile-image")
	if err != nil {
		log.Println("   Unable to upload profile picture: ", err)
		return
	}
	defer file.Close()

	imageID := uuid.New().String()

	browserPath, err := h.Services.ServiceSaveProfileImageVersions(r.Context(), file, imageID, user.Slug)
	if err != nil {
		log.Println("   Unable to save image versions: ", err)
		http.Error(w, "Unable to save image versions: ", http.StatusInternalServerError)
		return
	}

	err = database.UsersTableUpdateProfileImage(user.UserID, imageID, browserPath)
	if err != nil {
		log.Println("   Could not save image path to db: ", err)
		http.Error(w, "Cound not save image path to db", http.StatusInternalServerError)
		return
	}

	err = h.Services.ServiceDeleteProfileImageVersions(r.Context(), oldImageID, user.Slug)
	if err != nil {
		log.Println("   Could not delete old image path from cloud: ", err)
	}

	log.Println("   Saved file to users table and:", browserPath)
	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

func (h Handler) HandlerSettingsPage(w http.ResponseWriter, r *http.Request) {

	user, band, err := HelperGetAuthenticatedUserAndBand(r)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	data := models.SettingsPageData{
		CSRFToken:   csrf.Token(r),
		User:        user,
		CurrentBand: band,
	}

	err = h.Tmpl.ExecuteTemplate(w, "settings.html", data)
	if err != nil {
		log.Println("   Err getting profile pic page: ", err)
		return
	}
}

func (h Handler) HandlerAdmin(w http.ResponseWriter, r *http.Request) {

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

	if user.IsAdmin == false {
		http.Error(w, "Unable to get admin page", http.StatusForbidden)
		return
	}

	bandMembers, err := database.BandMembersGetMembersByBandID(band.BandID)
	if err != nil {
		log.Println("   Unable to get band members from database: ", err)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	data := models.AdminPageData{
		User:  user,
		Band:  band,
		Users: bandMembers,
	}

	err = h.Tmpl.ExecuteTemplate(w, "admin.html", data)
	if err != nil {
		log.Println("   err getting admin page.html: ", err)
		return
	}
	return
}

func (h Handler) HandlerBandJoin(w http.ResponseWriter, r *http.Request) {
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

	accessCode := helpers.NormalizeAccessCode(r.FormValue("band-access-code"))
	if accessCode == "" {
		http.Error(w, "Access code is required", http.StatusBadRequest)
		return
	}
	accessCodeHash := helpers.HashRegistrationCode(accessCode)

	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, "Unable to load session", http.StatusUnauthorized)
		return
	}

	tokenHash := helpers.HashSessionToken(cookie.Value)
	_, err = database.BandMembersJoinWithAccessCode(
		r.Context(),
		auth.User.UserID,
		tokenHash,
		accessCodeHash,
	)
	if errors.Is(err, database.ErrAccessCodeExpired) {
		http.Error(w, "Access code is invalid or expired", http.StatusGone)
		return
	}
	if errors.Is(err, database.ErrBandAlreadyJoined) {
		http.Error(w, "You are already a member of this band", http.StatusConflict)
		return
	}
	if err != nil {
		slog.Error(
			"unable to add user to band",
			"request_id", requestlog.GetRequestID(r.Context()),
			"path", r.URL.Path,
			"error", err,
		)
		http.Error(w, "Unable to add user to band", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/chats", http.StatusSeeOther)
}
