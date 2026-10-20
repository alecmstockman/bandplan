package handlers

import (
	requestlog "bandplan/src/logging"
	"log"
	"log/slog"
	"net/http"
)

func (h Handler) HandlerHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")

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

	data, err := h.Services.HomePage(r.Context(), user, band)
	if err != nil {
		slog.Error(
			"unable to get homepage data",
			"request_id", requestlog.GetRequestID(r.Context()),
			"path", r.URL.Path,
			"error", err,
		)
		http.Error(w, "Unable to load page", http.StatusInternalServerError)
		return
	}

	err = h.Tmpl.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		log.Println("   template err:", err)
		return
	}
}
