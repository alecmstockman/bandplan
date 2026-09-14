package handlers

import (
	requestlog "bandplan/src/logging"
	"bandplan/src/models"
	"log"
	"log/slog"
	"net/http"
)

func (h Handler) HandlerPromotion(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerPromotions")

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		slog.Error(
			"request started",
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
	log.Print("- HandlerSetlists")

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		slog.Error(
			"request started",
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
	log.Print("- HandlerCalendar")

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		slog.Error(
			"request started",
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
	log.Print("- HandlerFile")

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		slog.Error(
			"request started",
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
