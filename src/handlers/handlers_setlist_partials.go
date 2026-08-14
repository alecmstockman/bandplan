package handlers

import (
	"bandplan/src/database"
	"bandplan/src/models"
	"fmt"
	"log"
	"net/http"
)

func (h Handler) HandlerSetlistUpdateCountButtonSongs(w http.ResponseWriter, r *http.Request) {
	fmt.Println("----------------------------------")
	log.Println("- HandlerSetlistUpdateCountButtonSong")

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		log.Println("   Unable to get AuthContext: ", err)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		w.Header().Set("HX-Redirect", "/")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	user := auth.User
	band := auth.CurrentBand
	setlistID := r.URL.Query().Get("id")
	fmt.Println("Setlist Page, SetlistID: ", setlistID)

	setlist, err := database.SetlistsTableGetSetlistByID(setlistID)
	if err != nil {
		log.Println("   Unable to get setlist: ", err)
		return
	}

	data := models.SetlistPage{
		User:    user,
		Band:    band,
		Setlist: setlist,
	}

	err = h.Tmpl.ExecuteTemplate(w, "setlist_song_count", data)
	if err != nil {
		log.Println("   Unable to get setlist_count template: ", err)
		http.Error(w, "Unable to get setlist_count template", http.StatusInternalServerError)
	}
}

func (h Handler) HandlerSetlistUpdateCountButtonTransitions(w http.ResponseWriter, r *http.Request) {
	fmt.Println("----------------------------------")
	log.Println("- HandlerSetlistUpdateCountButtonTransitions")

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		log.Println("   Unable to get AuthContext: ", err)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		w.Header().Set("HX-Redirect", "/")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	user := auth.User
	band := auth.CurrentBand
	setlistID := r.URL.Query().Get("id")
	fmt.Println("Setlist Page, SetlistID: ", setlistID)

	setlist, err := database.SetlistsTableGetSetlistByID(setlistID)
	if err != nil {
		log.Println("   Unable to get setlist: ", err)
		return
	}

	data := models.SetlistPage{
		User:    user,
		Band:    band,
		Setlist: setlist,
	}

	err = h.Tmpl.ExecuteTemplate(w, "setlist_transition_count", data)
	if err != nil {
		log.Println("   Unable to get setlist_count template: ", err)
		http.Error(w, "Unable to get setlist_count template", http.StatusInternalServerError)
	}
}

func (h Handler) HandlerSetlistUpdateCountButtonBreaks(w http.ResponseWriter, r *http.Request) {
	fmt.Println("----------------------------------")
	log.Println("- HandlerSetlistUpdateCountButtonTransitions")

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		log.Println("   Unable to get AuthContext: ", err)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		w.Header().Set("HX-Redirect", "/")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	user := auth.User
	band := auth.CurrentBand
	setlistID := r.URL.Query().Get("id")
	fmt.Println("Setlist Page, SetlistID: ", setlistID)

	setlist, err := database.SetlistsTableGetSetlistByID(setlistID)
	if err != nil {
		log.Println("   Unable to get setlist: ", err)
		return
	}

	data := models.SetlistPage{
		User:    user,
		Band:    band,
		Setlist: setlist,
	}

	err = h.Tmpl.ExecuteTemplate(w, "setlist_break_count", data)
	if err != nil {
		log.Println("   Unable to get setlist_break template: ", err)
		http.Error(w, "Unable to get setlist_break template", http.StatusInternalServerError)
	}
}

func (h Handler) HandlerSetlistUpdateCountButtonItems(w http.ResponseWriter, r *http.Request) {
	fmt.Println("----------------------------------")
	log.Println("- HandlerSetlistUpdateCountButtonItems")

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		log.Println("   Unable to get AuthContext: ", err)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		w.Header().Set("HX-Redirect", "/")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	user := auth.User
	band := auth.CurrentBand
	setlistID := r.URL.Query().Get("id")
	fmt.Println("Setlist Page, SetlistID: ", setlistID)

	setlist, err := database.SetlistsTableGetSetlistByID(setlistID)
	if err != nil {
		log.Println("   Unable to get setlist: ", err)
		return
	}

	data := models.SetlistPage{
		User:    user,
		Band:    band,
		Setlist: setlist,
	}

	err = h.Tmpl.ExecuteTemplate(w, "setlist_item_count", data)
	if err != nil {
		log.Println("   Unable to get setlist_item template: ", err)
		http.Error(w, "Unable to get setlist_item template", http.StatusInternalServerError)
	}
}

func (h Handler) HandlerSetlistUpdateTimeButtonSongs(w http.ResponseWriter, r *http.Request) {
	fmt.Println("----------------------------------")
	log.Println("- HandlerSetlistUpdateTimeButtonSongs")

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		log.Println("   Unable to get AuthContext: ", err)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		w.Header().Set("HX-Redirect", "/")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	user := auth.User
	band := auth.CurrentBand
	setlistID := r.URL.Query().Get("id")
	fmt.Println("Setlist Page, SetlistID: ", setlistID)

	setlist, err := database.SetlistsTableGetSetlistByID(setlistID)
	if err != nil {
		log.Println("   Unable to get setlist: ", err)
		return
	}

	data := models.SetlistPage{
		User:    user,
		Band:    band,
		Setlist: setlist,
	}

	err = h.Tmpl.ExecuteTemplate(w, "setlist_songs_time", data)
	if err != nil {
		log.Println("   Unable to get setlist_time_songs template: ", err)
		http.Error(w, "Unable to get setlist_time_songs template", http.StatusInternalServerError)
	}
}

func (h Handler) HandlerSetlistUpdateTimeButtonTransitions(w http.ResponseWriter, r *http.Request) {
	fmt.Println("----------------------------------")
	log.Println("- HandlerSetlistUpdateTimeButtonTransitions")

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		log.Println("   Unable to get AuthContext: ", err)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		w.Header().Set("HX-Redirect", "/")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	user := auth.User
	band := auth.CurrentBand
	setlistID := r.URL.Query().Get("id")
	fmt.Println("Setlist Page, SetlistID: ", setlistID)

	setlist, err := database.SetlistsTableGetSetlistByID(setlistID)
	if err != nil {
		log.Println("   Unable to get setlist: ", err)
		return
	}

	data := models.SetlistPage{
		User:    user,
		Band:    band,
		Setlist: setlist,
	}

	err = h.Tmpl.ExecuteTemplate(w, "setlist_transitions_time", data)
	if err != nil {
		log.Println("   Unable to get setlist_item template: ", err)
		http.Error(w, "Unable to get setlist_item template", http.StatusInternalServerError)
	}
}

func (h Handler) HandlerSetlistUpdateTimeButtonBreaks(w http.ResponseWriter, r *http.Request) {
	fmt.Println("----------------------------------")
	log.Println("- HandlerSetlistUpdateTimeButtonBreaks")

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		log.Println("   Unable to get AuthContext: ", err)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		w.Header().Set("HX-Redirect", "/")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	user := auth.User
	band := auth.CurrentBand
	setlistID := r.URL.Query().Get("id")
	fmt.Println("Setlist Page, SetlistID: ", setlistID)

	setlist, err := database.SetlistsTableGetSetlistByID(setlistID)
	if err != nil {
		log.Println("   Unable to get setlist: ", err)
		return
	}

	data := models.SetlistPage{
		User:    user,
		Band:    band,
		Setlist: setlist,
	}

	err = h.Tmpl.ExecuteTemplate(w, "setlist_breaks_time", data)
	if err != nil {
		log.Println("   Unable to get setlist_item template: ", err)
		http.Error(w, "Unable to get setlist_item template", http.StatusInternalServerError)
	}
}

func (h Handler) HandlerSetlistUpdateTimeButtonItems(w http.ResponseWriter, r *http.Request) {
	fmt.Println("----------------------------------")
	log.Println("- HandlerSetlistUpdateTimeButtonItems")

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		log.Println("   Unable to get AuthContext: ", err)
		w.Header().Set("HX-Redirect", "/")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	user := auth.User
	band := auth.CurrentBand
	setlistID := r.URL.Query().Get("id")
	fmt.Println("Setlist Page, SetlistID: ", setlistID)

	setlist, err := database.SetlistsTableGetSetlistByID(setlistID)
	if err != nil {
		log.Println("   Unable to get setlist: ", err)
		return
	}

	data := models.SetlistPage{
		User:    user,
		Band:    band,
		Setlist: setlist,
	}

	err = h.Tmpl.ExecuteTemplate(w, "setlist_items_time", data)
	if err != nil {
		log.Println("   Unable to get setlist_item template: ", err)
		http.Error(w, "Unable to get setlist_item template", http.StatusInternalServerError)
	}
}

func (h Handler) HandlerSetlistOpenNotesPage(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerSetlistOpenNotesPage")

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		log.Println("   Unable to get AuthContext: ", err)
		w.Header().Set("HX-Redirect", "/")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	user := auth.User
	band := auth.CurrentBand
	setlistID := r.URL.Query().Get("id")

	setlist, err := database.SetlistsTableGetSetlistByID(setlistID)
	if err != nil {
		log.Println("   Unable to get setlist: ", err)
		return
	}

	data := models.SetlistPage{
		User:    user,
		Band:    band,
		Setlist: setlist,
		BackURL: "/setlist?id=" + setlistID,
	}

	err = h.Tmpl.ExecuteTemplate(w, "setlist_notes", data)
	if err != nil {
		log.Println("   Unable to get setlist_item template: ", err)
		http.Error(w, "Unable to get setlist_item template", http.StatusInternalServerError)
	}
}

func (h Handler) HandlerSetlistEditNotesPage(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerSetlistEditNotesPage")

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		log.Println("   Unable to get AuthContext: ", err)
		w.Header().Set("HX-Redirect", "/")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	user := auth.User
	band := auth.CurrentBand
	setlistID := r.URL.Query().Get("id")

	setlist, err := database.SetlistsTableGetSetlistByID(setlistID)
	if err != nil {
		log.Println("   Unable to get setlist: ", err)
		return
	}

	data := models.SetlistPage{
		User:    user,
		Band:    band,
		Setlist: setlist,
		BackURL: "/setlist?id=" + setlistID,
	}

	err = h.Tmpl.ExecuteTemplate(w, "setlist_notes_edit", data)
	if err != nil {
		log.Println("   Unable to get setlist_item template: ", err)
		http.Error(w, "Unable to get setlist_item template", http.StatusInternalServerError)
	}
}

func (h Handler) HandlerSetlistSaveNotesPage(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerSetlistSaveNotesPage")

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		log.Println("   Unable to get AuthContext: ", err)
		w.Header().Set("HX-Redirect", "/")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	user := auth.User
	band := auth.CurrentBand
	setlistID := r.URL.Query().Get("id")

	setlist, err := database.SetlistsTableGetSetlistByID(setlistID)
	if err != nil {
		log.Println("   Unable to get setlist: ", err)
		return
	}

	data := models.SetlistPage{
		User:    user,
		Band:    band,
		Setlist: setlist,
		BackURL: "/setlist?id=" + setlistID,
	}

	err = h.Tmpl.ExecuteTemplate(w, "setlist.html", data)
	if err != nil {
		log.Println("   Unable to render setlist: ", err)
		http.Error(w, "Unable to load setlist", http.StatusInternalServerError)
		return
	}
	return
}
