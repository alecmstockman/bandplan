package handlers

import (
	requestlog "bandplan/src/logging"
	"bandplan/src/models"
	"log"
	"log/slog"
	"net/http"
)

func (h Handler) HandlerChecklistsPage(w http.ResponseWriter, r *http.Request) {
	log.Print("- HandlerChecklistsPage")

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
		User: user,
		Band: band,
	}

	err = h.Tmpl.ExecuteTemplate(w, "checklists.html", data)
}
