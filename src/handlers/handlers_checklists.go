package handlers

import (
	requestlog "bandplan/src/logging"
	"bandplan/src/models"
	"log"
	"log/slog"
	"net/http"

	"github.com/gorilla/csrf"
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

func (h Handler) HandlerChecklistCreatePage(w http.ResponseWriter, r *http.Request) {

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

	data := models.ChecklistCreatePageData{
		CSRFToken: csrf.Token(r),
		User:      auth.User,
		Band:      auth.CurrentBand,
		Checklist: models.Checklist{},
	}

	err = h.Tmpl.ExecuteTemplate(w, "checklist_create.html", data)
	if err != nil {
		slog.Error(
			"unable to execute checklist_create.html",
			"request_id", requestlog.GetRequestID(r.Context()),
			"path", r.URL.Path,
			"error", err,
		)
		http.Error(w, "Unable to load checklist creation page", http.StatusInternalServerError)
		return
	}
}

func (h Handler) HandlerChecklistAdd(w http.ResponseWriter, r *http.Request) {

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
