package handlers

import (
	requestlog "bandplan/src/logging"
	"bandplan/src/models"
	"log/slog"
	"net/http"
)

func (h Handler) HandlerNotificationSettingsPage(w http.ResponseWriter, r *http.Request) {
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

	data := models.SettingsPageData{
		User:        auth.User,
		CurrentBand: auth.CurrentBand,
	}

	if err := h.Tmpl.ExecuteTemplate(w, "notification-settings.html", data); err != nil {
		slog.Error(
			"unable to load notification settings",
			"request_id", requestlog.GetRequestID(r.Context()),
			"path", r.URL.Path,
			"error", err,
		)
		http.Error(w, "Unable to load notification settings", http.StatusInternalServerError)
	}
}
