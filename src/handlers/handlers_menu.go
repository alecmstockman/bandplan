package handlers

import (
	requestlog "bandplan/src/logging"
	"bandplan/src/models"
	"log/slog"
	"net/http"
)

func (h Handler) HandlerPromotion(w http.ResponseWriter, r *http.Request) {

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		slog.Error(
			"unable to load authenticated user",
			"request_id", requestlog.GetRequestID(r.Context()),
			"user_id", auth.User.UserID,
			"band_id", auth.CurrentBand.BandID,
			"method", r.Method,
			"path", r.URL.Path,
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

	err = h.Tmpl.ExecuteTemplate(w, "promotion.html", data)
}

func (h Handler) HandlerGoals(w http.ResponseWriter, r *http.Request) {

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		slog.Error(
			"unable to load authenticated user",
			"request_id", requestlog.GetRequestID(r.Context()),
			"user_id", auth.User.UserID,
			"band_id", auth.CurrentBand.BandID,
			"method", r.Method,
			"path", r.URL.Path,
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

func (h Handler) HandlerCalendar(w http.ResponseWriter, r *http.Request) {

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		slog.Error(
			"unable to load authenticated user",
			"request_id", requestlog.GetRequestID(r.Context()),
			"user_id", auth.User.UserID,
			"band_id", auth.CurrentBand.BandID,
			"method", r.Method,
			"path", r.URL.Path,
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

	err = h.Tmpl.ExecuteTemplate(w, "calendar.html", data)
}

func (h Handler) HandlerFiles(w http.ResponseWriter, r *http.Request) {

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		slog.Error(
			"unable to load authenticated user",
			"request_id", requestlog.GetRequestID(r.Context()),
			"user_id", auth.User.UserID,
			"band_id", auth.CurrentBand.BandID,
			"method", r.Method,
			"path", r.URL.Path,
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

	err = h.Tmpl.ExecuteTemplate(w, "files.html", data)
}
