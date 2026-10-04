package handlers

import (
	requestlog "bandplan/src/logging"
	"bandplan/src/models"
	"log/slog"
	"net/http"

	"github.com/gorilla/csrf"
)

func (h Handler) HandlerGoalsPage(w http.ResponseWriter, r *http.Request) {

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

	err = h.Tmpl.ExecuteTemplate(w, "goals.html", data)
}

func (h Handler) HandlerGoalCreatePage(w http.ResponseWriter, r *http.Request) {

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

	data := models.GoalsListCreatePageData{
		CSRFToken: csrf.Token(r),
		User:      auth.User,
		Band:      auth.CurrentBand,
		GoalsList: models.GoalsList{},
	}

	err = h.Tmpl.ExecuteTemplate(w, "goals_create.html", data)
	if err != nil {
		slog.Error(
			"unable to execute goals_create.html",
			"request_id", requestlog.GetRequestID(r.Context()),
			"path", r.URL.Path,
			"error", err,
		)
		http.Error(w, "Unable to load goals list creation page", http.StatusInternalServerError)
		return
	}
}

func (h Handler) HandlerGoalAdd(w http.ResponseWriter, r *http.Request) {

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
