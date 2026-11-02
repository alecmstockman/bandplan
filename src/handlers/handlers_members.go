package handlers

import (
	"bandplan/src/database"
	requestlog "bandplan/src/logging"
	"bandplan/src/models"
	"log"
	"log/slog"
	"net/http"
)

func (h Handler) HandlerMembers(w http.ResponseWriter, r *http.Request) {

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

	err = h.Tmpl.ExecuteTemplate(w, "members.html", data)
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
