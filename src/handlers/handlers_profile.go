package handlers

import (
	"bandplan/src/clients"
	"bandplan/src/database"
	"bandplan/src/helpers"
	requestlog "bandplan/src/logging"
	"bandplan/src/models"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/csrf"
	"github.com/lib/pq"
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

	// fmt.Printf("\nuser: \n%+v", user)

	notifications, err := database.NotifTableGetNotificationsByUserID(r.Context(), user.UserID, band.BandID)
	if err != nil {
		slog.Error(
			"unable to load notifications",
			"request_id", requestlog.GetRequestID(r.Context()),
			"error", err,
		)
		http.Error(w, "Unable to load notifications", http.StatusInternalServerError)
		return
	}

	data := models.MenuPageData{
		CSRFToken:     csrf.Token(r),
		User:          user,
		Band:          band,
		Notifications: notifications,
	}

	err = h.Tmpl.ExecuteTemplate(w, "profile.html", data)
	if err != nil {
		slog.Error(
			"unable to loading profile.html",
			"request_id", requestlog.GetRequestID(r.Context()),
			"path", r.URL.Path,
			"error", err,
		)
		return
	}
}

func (h Handler) HandlerProfileEditPage(w http.ResponseWriter, r *http.Request) {

	fmt.Println("- HandlerProfileEditPage")

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

	notifications, err := database.NotifTableGetNotificationsByUserID(r.Context(), user.UserID, band.BandID)
	if err != nil {
		slog.Error(
			"unable to load notifications",
			"request_id", requestlog.GetRequestID(r.Context()),
			"error", err,
		)
		http.Error(w, "Unable to load notifications", http.StatusInternalServerError)
		return
	}

	data := models.MenuPageData{
		CSRFToken:     csrf.Token(r),
		User:          user,
		Band:          band,
		Notifications: notifications,
	}

	err = h.Tmpl.ExecuteTemplate(w, "profile-edit.html", data)
	if err != nil {
		log.Println("   Err getting profile pic page: ", err)
		return
	}
}

func (h Handler) HandlerProfileUpdate(w http.ResponseWriter, r *http.Request) {

	fmt.Println("- HandlerProfileUpdate")

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

	firstName := strings.TrimSpace(r.FormValue("first-name"))
	if !helpers.ValidateNameEntryLength(firstName) {
		http.Error(w, "Invalid entry for first name", http.StatusBadRequest)
		return
	}

	lastName := strings.TrimSpace(r.FormValue("last-name"))
	if !helpers.ValidateNameEntryMaxLength(lastName) {
		http.Error(w, "Invalid entry for last name", http.StatusBadRequest)
		return
	}

	displayName := strings.TrimSpace(r.FormValue("display-name"))
	if !helpers.ValidateNameEntryLength(displayName) {
		http.Error(w, "Invalid entry for display name", http.StatusBadRequest)
		return
	}

	email := helpers.NormalizeEmail(r.FormValue("email"))
	if !helpers.ValidateEmail(email) {
		http.Error(w, "Invalid email entry", http.StatusBadRequest)
		return
	}

	timeZone := strings.TrimSpace(r.FormValue("timezone"))
	if _, err := time.LoadLocation(timeZone); err != nil {
		http.Error(w, "Invalid timezone", http.StatusBadRequest)
		return
	}

	user := auth.User
	user.Name = strings.TrimSpace(firstName + " " + lastName)
	user.FirstName = firstName
	user.LastName = lastName
	user.DisplayName = displayName
	user.Email = email
	user.TimeZone = timeZone

	emailChanged, err := database.UsersTableUpdateProfile(r.Context(), user)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			http.Error(w, "Email address is already in use", http.StatusConflict)
			return
		}

		slog.Error(
			"unable to update profile",
			"request_id", requestlog.GetRequestID(r.Context()),
			"user_id", user.UserID,
			"error", err,
		)
		http.Error(w, "Unable to update profile", http.StatusInternalServerError)
		return
	}

	if !emailChanged {
		http.Redirect(w, r, "/profile", http.StatusSeeOther)
		return
	}

	destination := "/profile/email-verification-sent"
	if err := h.Services.RegistrationSendAccountVerificationEmail(r.Context(), user.UserID, user.Email); err != nil {
		slog.Error(
			"unable to send updated email verification",
			"request_id", requestlog.GetRequestID(r.Context()),
			"user_id", user.UserID,
			"error", err,
		)
		destination += "?delivery=failed"
	}

	http.Redirect(w, r, destination, http.StatusSeeOther)
}

func (h Handler) HandlerProfileEmailVerificationSent(w http.ResponseWriter, r *http.Request) {
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

	w.Header().Set("Cache-Control", "no-store")
	data := models.ProfileVerificationPage{
		CSRFToken: csrf.Token(r),
		Email:     auth.User.Email,
		Sent:      r.URL.Query().Get("delivery") != "failed",
	}

	if err := h.Tmpl.ExecuteTemplate(w, "profile-email-verification.html", data); err != nil {
		slog.Error(
			"unable to load profile email verification page",
			"request_id", requestlog.GetRequestID(r.Context()),
			"error", err,
		)
		http.Error(w, "Unable to load page", http.StatusInternalServerError)
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

	notifications, err := database.NotifTableGetNotificationsByUserID(r.Context(), auth.User.UserID, auth.CurrentBand.BandID)
	if err != nil {
		slog.Error(
			"unable to load notifications",
			"request_id", requestlog.GetRequestID(r.Context()),
			"error", err,
		)
		http.Error(w, "Unable to load notifications", http.StatusInternalServerError)
		return
	}

	data := models.BandsPageData{
		CSRFToken:     csrf.Token(r),
		User:          auth.User,
		Band:          auth.CurrentBand,
		Bands:         bands,
		Notifications: notifications,
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

	notifications, err := database.NotifTableGetNotificationsByUserID(r.Context(), auth.User.UserID, auth.CurrentBand.BandID)
	if err != nil {
		slog.Error(
			"unable to load notifications",
			"request_id", requestlog.GetRequestID(r.Context()),
			"error", err,
		)
		http.Error(w, "Unable to load notifications", http.StatusInternalServerError)
		return
	}

	data := models.BandPageData{
		User:          auth.User,
		Band:          auth.CurrentBand,
		SelectedBand:  band,
		Notifications: notifications,
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

	file, _, err := r.FormFile("profile-image")
	if err != nil {
		log.Println("   Unable to upload profile picture: ", err)
		return
	}
	defer file.Close()

	imageID := uuid.New().String()

	browserPath, err := h.Services.ServiceSaveProfileImageVersions(r.Context(), file, imageID, user.Slug)
	if err != nil {
		slog.Error(
			"unable to save profile images",
			"request_id", requestlog.GetRequestID(r.Context()),
			"path", r.URL.Path,
			"error", err,
		)
		writeImageUploadError(w, err)
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

	user, band, err := helpers.HelperGetAuthenticatedUserAndBand(r)
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

func (h Handler) HandlerProfileFeedback(w http.ResponseWriter, r *http.Request) {
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

	feedback := strings.TrimSpace(r.FormValue("feedback"))
	if feedback == "" {
		http.Error(w, "Feedback is required", http.StatusBadRequest)
		return
	}
	if len([]rune(feedback)) > 4000 {
		http.Error(w, "Feedback is too long", http.StatusBadRequest)
		return
	}

	if err := clients.SendFeedbackEmail(auth.User, auth.CurrentBand.BandID, feedback); err != nil {
		slog.Error(
			"unable to send user feedback",
			"request_id", requestlog.GetRequestID(r.Context()),
			"user_id", auth.User.UserID,
			"band_id", auth.CurrentBand.BandID,
			"error", err,
		)
		http.Error(w, "Unable to send feedback", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/settings?toast=feedback-sent", http.StatusSeeOther)
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
		log.Println("User is not an admin")
		http.Error(w, "Unable to get admin page", http.StatusForbidden)
		return
	}

	bandMembers, err := database.BandMembersGetMembersByBandID(band.BandID)
	if err != nil {
		log.Println("   Unable to get band members from database: ", err)
		slog.Error(
			"Unable to get band members from database",
			"request_id", requestlog.GetRequestID(r.Context()),
			"path", r.URL.Path,
			"error", err,
		)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	notifications, err := database.NotifTableGetNotificationsByUserID(r.Context(), user.UserID, band.BandID)
	if err != nil {
		slog.Error(
			"unable to load notifications",
			"request_id", requestlog.GetRequestID(r.Context()),
			"error", err,
		)
		http.Error(w, "Unable to load notifications", http.StatusInternalServerError)
		return
	}

	data := models.AdminPageData{
		User:          user,
		Band:          band,
		Users:         bandMembers,
		Notifications: notifications,
	}

	err = h.Tmpl.ExecuteTemplate(w, "admin.html", data)
	if err != nil {
		slog.Error(
			"err getting admin page.html",
			"request_id", requestlog.GetRequestID(r.Context()),
			"path", r.URL.Path,
			"error", err,
		)
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
