package handlers

import (
	"bandplan/src/database"
	"bandplan/src/models"
	"fmt"
	"log"
	"net/http"
	"slices"
	"strconv"
)

func (h Handler) HandlerSetlistSongs(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerSetlistSongs")

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		log.Println("   Unable to get AuthContext: ", err)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
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
	}

	err = h.Tmpl.ExecuteTemplate(w, "setlist_song_count", data)
	if err != nil {
		log.Println("   Unable to get setlist_count template: ", err)
		http.Error(w, "Unable to get setlist_count template", http.StatusInternalServerError)
	}
}

func (h Handler) HandlerSetlistTransitions(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerSetlistTransitions")

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		log.Println("   Unable to get AuthContext: ", err)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
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
	}

	err = h.Tmpl.ExecuteTemplate(w, "setlist_transition_count", data)
	if err != nil {
		log.Println("   Unable to get setlist_count template: ", err)
		http.Error(w, "Unable to get setlist_count template", http.StatusInternalServerError)
	}
}

func (h Handler) HandlerSetlistBreaks(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerSetlistBreaks")

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		log.Println("   Unable to get AuthContext: ", err)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
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
	}

	err = h.Tmpl.ExecuteTemplate(w, "setlist_break_count", data)
	if err != nil {
		log.Println("   Unable to get setlist_break template: ", err)
		http.Error(w, "Unable to get setlist_break template", http.StatusInternalServerError)
	}
}

func (h Handler) HandlerSetlistItems(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerSetlistItems")

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		log.Println("   Unable to get AuthContext: ", err)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
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
	}

	err = h.Tmpl.ExecuteTemplate(w, "setlist_item_count", data)
	if err != nil {
		log.Println("   Unable to get setlist_break template: ", err)
		http.Error(w, "Unable to get setlist_break template", http.StatusInternalServerError)
	}
}

func (h Handler) HandlerSetlistUpdateCountButtonItems(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerSetlistUpdateCountButtonItems")

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		log.Println("   Unable to get AuthContext: ", err)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
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
	}

	err = h.Tmpl.ExecuteTemplate(w, "setlist_item_count", data)
	if err != nil {
		log.Println("   Unable to get setlist_item template: ", err)
		http.Error(w, "Unable to get setlist_item template", http.StatusInternalServerError)
	}
}

func (h Handler) HandlerSetlistUpdateTimeButtonSongs(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerSetlistUpdateTimeButtonSongs")

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		log.Println("   Unable to get AuthContext: ", err)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
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
	}

	err = h.Tmpl.ExecuteTemplate(w, "setlist_songs_time", data)
	if err != nil {
		log.Println("   Unable to get setlist_time_songs template: ", err)
		http.Error(w, "Unable to get setlist_time_songs template", http.StatusInternalServerError)
	}
}

func (h Handler) HandlerSetlistUpdateTimeButtonTransitions(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerSetlistUpdateTimeButtonTransitions")

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		log.Println("   Unable to get AuthContext: ", err)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
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
	}

	err = h.Tmpl.ExecuteTemplate(w, "setlist_transitions_time", data)
	if err != nil {
		log.Println("   Unable to get setlist_item template: ", err)
		http.Error(w, "Unable to get setlist_item template", http.StatusInternalServerError)
	}
}

func (h Handler) HandlerSetlistUpdateTimeButtonBreaks(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerSetlistUpdateTimeButtonBreaks")

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		log.Println("   Unable to get AuthContext: ", err)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
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
	}

	err = h.Tmpl.ExecuteTemplate(w, "setlist_breaks_time", data)
	if err != nil {
		log.Println("   Unable to get setlist_item template: ", err)
		http.Error(w, "Unable to get setlist_item template", http.StatusInternalServerError)
	}
}

func (h Handler) HandlerSetlistUpdateTimeButtonItems(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
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
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
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
	fmt.Println("-----------------------------")
	log.Println("- HandlerSetlistSaveNotesPage")

	_, err := HelperGetAuthContext(r)
	if err != nil {
		log.Println("   Unable to get AuthContext: ", err)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	setlistID := r.FormValue("setlist-id")
	notes := r.FormValue("notes")

	err = database.SetlistsTableUpdateNotes(setlistID, notes)
	if err != nil {
		log.Println("   Unable to save notes to database: ", err)
		http.Error(w, "Unable to save notes to database", http.StatusInternalServerError)
		return
	}

	redirectURL := "/setlist?id=" + setlistID
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

func (h Handler) HandlerSetlistEditInfoCard(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerSetlistEditInfoCard")

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		log.Println("   Unable to get Auth Context: ", err)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	user := auth.User
	band := auth.CurrentBand
	song_id := r.URL.Query().Get("id")

	fmt.Println("song: ", song_id)

	setlistID := r.FormValue("setlist-id")
	itemType := r.FormValue("item-type")

	itemID := r.FormValue("item-id")
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
	fmt.Println("songID: ", itemID)

	var newItemType models.SetlistItemType

	switch itemType {
	case "song":
		newItemType = models.SetlistItemSong
	case "transition":
		newItemType = models.SetlistItemTransition
	case "break":
		newItemType = models.SetlistItemBreak
	}

	setlist, err := database.SetlistItemsGetItem(setlistID, newItemType, itemID)
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
		Type:   newItemType,
	}

	fmt.Println("TEST!")
	err = h.Tmpl.ExecuteTemplate(w, "setlist_item_edit", data)
	if err != nil {
		log.Println("   Unable to get setlist_item_edit template: ", err)
		http.Error(w, "Unable to get setlist_item_editm template", http.StatusInternalServerError)
	}
}

func (h Handler) HandlerSetlistSaveInfoCard(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerSetlistSaveInfoCard")

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		log.Println("   Unable to get Auth Context: ", err)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	user := auth.User
	band := auth.CurrentBand

	setlistID := r.FormValue("setlist-id")
	itemType := r.FormValue("item-type")
	itemID := r.FormValue("item-id")
	itemTitle := r.FormValue("item-title")

	itemPosition, err := strconv.Atoi(r.FormValue("item-position"))
	if err != nil {
		log.Println("   Unable to convert item-position to integer: ", err)
		http.Error(w, "Unable to get item position", http.StatusInternalServerError)
		return
	}
	position := r.FormValue("position")

	itemLength, err := strconv.Atoi(r.FormValue("item-length"))
	if err != nil {
		log.Println("   Unable to convert item lenght to int: ", err)
		http.Error(w, "Unable to convert item lenght to int", http.StatusInternalServerError)
		return
	}

	order, err := database.SetlistItemsGetSetlistOrder(setlistID)
	if err != nil {
		log.Println("   Unable to get setlist order: ", err)
		http.Error(w, "Error updating setlit item", http.StatusInternalServerError)
		return
	}

	newPosition, err := strconv.Atoi(position)
	if err != nil {
		log.Println("   Unable to convert position entry to integer: ", err)
		http.Error(w, "Unable to convert position to integer", http.StatusBadRequest)
		return
	}

	var tempOrder []models.ReorderItem

	if newPosition > len(order) {
		newPosition = len(order)
	}

	tempItem := order[itemPosition-1]
	tempOrder = slices.Delete(order, itemPosition-1, itemPosition)
	newOrder := slices.Insert(tempOrder, newPosition-1, tempItem)

	for i, _ := range newOrder {
		newOrder[i].Position = i + 1
	}

	err = database.SetlistItemsUpdateOrder(setlistID, order)
	if err != nil {
		log.Println("   Unable to update setlist order: ", err)
		http.Error(w, "Unable to update setlist order", http.StatusInternalServerError)
	}

	fmt.Println("order: ", order)
	pauseAfter, err := strconv.Atoi(r.FormValue("pause-after"))
	if err != nil {
		log.Println("   Unable to convert pause after to int: ", err)
		http.Error(w, "Unable to convert pause after to int", http.StatusInternalServerError)
		return
	}

	var newItemType models.SetlistItemType

	switch itemType {
	case "song":
		newItemType = models.SetlistItemSong
	case "transition":
		newItemType = models.SetlistItemTransition
	case "break":
		newItemType = models.SetlistItemBreak
	}

	err = database.SetlistItemsUpdateItem(setlistID, newItemType, itemID, pauseAfter)
	if err != nil {
		log.Println("   Unable to get setlist item from database: ", err)
		http.Error(w, "Unable to get setlist item from database", http.StatusInternalServerError)
		return
	}

	setlist, err := database.SetlistItemsGetItem(setlistID, newItemType, itemID)
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
}

func (h Handler) HandlerSetlistPopupInfoCard(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerSetlistPopupInfoCard")

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		log.Println("   Unable to get Auth Context: ", err)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
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
