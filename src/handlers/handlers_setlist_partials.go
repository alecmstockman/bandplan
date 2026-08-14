package handlers

import (
	"bandplan/src/database"
	"bandplan/src/models"
	"fmt"
	"log"
	"net/http"
	"strconv"
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

func (h Handler) HandlerSetlistEditInfoCard(w http.ResponseWriter, r *http.Request) {
	fmt.Println("-----------------------------")
	log.Println("- HandlerSetlistEditInfoCard")

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		log.Println("   Unable to get Auth Context: ", err)
		w.Header().Set("HX-Redirect", "/")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	user := auth.User
	band := auth.CurrentBand
	song_id := r.URL.Query().Get("id")

	fmt.Println("song: ", song_id)

	setlistID := r.FormValue("setlist-id")
	itemType := r.FormValue("item-type")
	songID := r.FormValue("song-id")
	itemTitle := r.FormValue("item-title")
	length := r.FormValue("item-length")
	fmt.Println("length: ", length)
	itemLength, err := strconv.Atoi(r.FormValue("item-length"))

	if err != nil {
		log.Println("   Unable to convert length to int: ", err)
		http.Error(w, "Unable to get item length", http.StatusInternalServerError)
		return
	}

	fmt.Println("setlistIDValue: ", setlistID)
	fmt.Println("itemType: ", itemType)
	fmt.Println("songID: ", songID)

	var newItemType models.SetlistItemType

	switch itemType {
	case "song":
		newItemType = models.SetlistItemSong
	case "transition":
		newItemType = models.SetlistItemTransition
	case "break":
		newItemType = models.SetlistItemBreak
	}

	setlist, err := database.SetlistItemsGetItem(setlistID, newItemType, songID)
	if err != nil {
		log.Println("   Unable to get setlist item from database: ", err)
		http.Error(w, "Unable to get setlist item from database", http.StatusInternalServerError)
		return
	}

	data := models.SetlistItemPage{
		User:   user,
		Band:   band,
		Item:   setlist,
		Title:  itemTitle,
		Length: itemLength,
	}

	err = h.Tmpl.ExecuteTemplate(w, "setlist_item_edit", data)
	if err != nil {
		log.Println("   Unable to get setlist_item_edit template: ", err)
		http.Error(w, "Unable to get setlist_item_editm template", http.StatusInternalServerError)
	}
	fmt.Println("end of handler")
}

func (h Handler) HandlerSetlistSaveInfoCard(w http.ResponseWriter, r *http.Request) {
	fmt.Println("-------------------------------")
	log.Println("- HandlerSetlistSaveInfoCard")

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		log.Println("   Unable to get Auth Context: ", err)
		w.Header().Set("HX-Redirect", "/")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	user := auth.User
	band := auth.CurrentBand

	setlistID := r.FormValue("setlist-id")
	itemType := r.FormValue("item-type")
	songID := r.FormValue("song-id")
	itemTitle := r.FormValue("item-title")
	length := r.FormValue("item-length")
	fmt.Println("length: ", length)

	itemLength, err := strconv.Atoi(r.FormValue("item-length"))
	if err != nil {
		log.Println("   Unable to convert item lenght to int: ", err)
		http.Error(w, "Unable to convert item lenght to int", http.StatusInternalServerError)
		return
	}

	pauseAfter, err := strconv.Atoi(r.FormValue("pause-after"))
	fmt.Println("pause-after: ", pauseAfter)
	if err != nil {
		log.Println("   Unable to convert pause after to int: ", err)
		http.Error(w, "Unable to convert pause after to int", http.StatusInternalServerError)
		return
	}

	fmt.Println("setlistIDValue: ", setlistID)
	fmt.Println("itemType: ", itemType)
	fmt.Println("songID: ", songID)
	fmt.Println("Pause after: ", pauseAfter)

	var newItemType models.SetlistItemType

	switch itemType {
	case "song":
		newItemType = models.SetlistItemSong
	case "transition":
		newItemType = models.SetlistItemTransition
	case "break":
		newItemType = models.SetlistItemBreak
	}

	err = database.SetlistItemsUpdateItem(setlistID, newItemType, songID, pauseAfter)
	if err != nil {
		log.Println("   Unable to get setlist item from database: ", err)
		http.Error(w, "Unable to get setlist item from database", http.StatusInternalServerError)
		return
	}

	setlist, err := database.SetlistItemsGetItem(setlistID, newItemType, songID)
	if err != nil {
		log.Println("   Unable to get setlist item from database: ", err)
		http.Error(w, "Unable to get setlist item from database", http.StatusInternalServerError)
		return
	}

	data := models.SetlistItemPage{
		User:   user,
		Band:   band,
		Item:   setlist,
		Title:  itemTitle,
		Length: itemLength,
	}

	err = h.Tmpl.ExecuteTemplate(w, "setlist_item_popup", data)
	if err != nil {
		log.Println("   Unable to get setlist_item_edit template: ", err)
		http.Error(w, "Unable to get setlist_item_editm template", http.StatusInternalServerError)
	}
	fmt.Println("end of handler")
}

func (h Handler) HandlerSetlistPopupInfoCard(w http.ResponseWriter, r *http.Request) {
	fmt.Println("-----------------------------")
	log.Println("- HandlerSetlistPopupInfoCard")

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		log.Println("   Unable to get Auth Context: ", err)
		w.Header().Set("HX-Redirect", "/")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	user := auth.User
	band := auth.CurrentBand
	song_id := r.URL.Query().Get("id")

	fmt.Println("song: ", song_id)

	setlistID := r.FormValue("setlist-id")
	itemType := r.FormValue("item-type")
	songID := r.FormValue("song-id")
	itemTitle := r.FormValue("item-title")
	length := r.FormValue("item-length")
	fmt.Println("length: ", length)
	itemLength, err := strconv.Atoi(r.FormValue("item-length"))

	if err != nil {
		log.Println("   Unable to convert length to int: ", err)
		http.Error(w, "Unable to get item length", http.StatusInternalServerError)
		return
	}

	fmt.Println("setlistIDValue: ", setlistID)
	fmt.Println("itemType: ", itemType)
	fmt.Println("songID: ", songID)

	var newItemType models.SetlistItemType

	switch itemType {
	case "song":
		newItemType = models.SetlistItemSong
	case "transition":
		newItemType = models.SetlistItemTransition
	case "break":
		newItemType = models.SetlistItemBreak
	}

	setlist, err := database.SetlistItemsGetItem(setlistID, newItemType, songID)
	if err != nil {
		log.Println("   Unable to get setlist item from database: ", err)
		http.Error(w, "Unable to get setlist item from database", http.StatusInternalServerError)
		return
	}

	data := models.SetlistItemPage{
		User:   user,
		Band:   band,
		Item:   setlist,
		Title:  itemTitle,
		Length: itemLength,
	}

	err = h.Tmpl.ExecuteTemplate(w, "setlist_item_popup", data)
	if err != nil {
		log.Println("   Unable to get setlist_item_edit template: ", err)
		http.Error(w, "Unable to get setlist_item_editm template", http.StatusInternalServerError)
	}
	fmt.Println("end of handler")

}
