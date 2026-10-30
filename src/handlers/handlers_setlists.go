package handlers

import (
	"bandplan/src/database"
	"bandplan/src/helpers"
	requestlog "bandplan/src/logging"
	"bandplan/src/models"
	"bandplan/src/services"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/csrf"
)

func (h Handler) HandlerSetlistsPage(w http.ResponseWriter, r *http.Request) {

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

	setlistSummaries, err := database.SetlistsTableGetSetlistSummariesByBandIDAndUserID(band.BandID, user.UserID)
	if err != nil {
		log.Println("   Unable to get setlist summaries band band id from db: ", err)
		http.Error(w, "Unable to get setlist summaries", http.StatusInternalServerError)
		return
	}

	notifications, err := database.NotifTableGetNotificationsByUserID(r.Context(), user.UserID, band.BandID)
	if err != nil {
		slog.Error(
			"unable to load notifications",
			"request_id", requestlog.GetRequestID(r.Context()),
			"error", err,
		)
		http.Error(w, "Unable to load notifications", http.StatusInternalServerError)
		return
	}

	data := models.SetlistData{
		CSRFToken:     csrf.Token(r),
		User:          user,
		Band:          band,
		Setlists:      setlistSummaries,
		Notifications: notifications,
	}

	err = h.Tmpl.ExecuteTemplate(w, "setlists.html", data)
	if err != nil {
		log.Println("   Unable to render setlists: ", err)
		http.Error(w, "Unable to load setlists", http.StatusInternalServerError)
		return
	}
}

func (h Handler) HandlerSetlistsAddPage(w http.ResponseWriter, r *http.Request) {

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

	notifications, err := database.NotifTableGetNotificationsByUserID(r.Context(), user.UserID, band.BandID)
	if err != nil {
		slog.Error(
			"unable to load notifications",
			"request_id", requestlog.GetRequestID(r.Context()),
			"error", err,
		)
		http.Error(w, "Unable to load notifications", http.StatusInternalServerError)
		return
	}

	data := models.SongDownloadData{
		CSRFToken:     csrf.Token(r),
		User:          user,
		Band:          band,
		Notifications: notifications,
	}

	err = h.Tmpl.ExecuteTemplate(w, "setlist-add.html", data)
	if err != nil {
		log.Println("   Err getting setlist-add page: ", err)
		http.Redirect(w, r, "/setlists", http.StatusSeeOther)
	}

	return
}

func (h Handler) HandlerSetlistsTempArt(w http.ResponseWriter, r *http.Request) {

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

	band := auth.CurrentBand

	r.Body = http.MaxBytesReader(w, r.Body, 12<<20)

	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Println("   Error parsing from while creating a new setlist: ", err)
		http.Error(w, "Unable to parse from", http.StatusBadRequest)
		return
	}
	defer r.MultipartForm.RemoveAll()

	imageID := ""
	previewURL := ""

	file, _, err := r.FormFile("artwork-path")
	if err != nil {
		if err != http.ErrMissingFile {
			log.Println("   Unable to read artwork file:", err)
			http.Error(w, "Unable to read artwork", http.StatusBadRequest)
			return
		}

		log.Println("   No setlist artwork uploaded")
	} else {
		defer file.Close()

		imageID = uuid.New().String()

		previewURL, err = h.Services.ServiceSaveTempImage(r.Context(), file, imageID, band.Slug, "setlist")
		if err != nil {
			slog.Error(
				"unable to save temp image",
				"request_id", requestlog.GetRequestID(r.Context()),
				"path", r.URL.Path,
				"error", err,
			)
			writeImageUploadError(w, err)
			return
		}
	}

	data := models.ArtworkPreviewData{
		ArtworkID:  imageID,
		PreviewURL: previewURL,
	}

	err = h.Tmpl.ExecuteTemplate(w, "setlist_artwork_preview", data)
	if err != nil {
		http.Error(w, "Unable to render preview", http.StatusInternalServerError)
		return
	}
}

func (h Handler) HandlerSetlistsTempArtDelete(w http.ResponseWriter, r *http.Request) {

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

	data := models.SongDownloadData{
		User: user,
		Band: band,
	}

	err = h.Tmpl.ExecuteTemplate(w, "setlist_artwork_reset", data)
	if err != nil {
		log.Println("   Unable to exececute : ", err)
		http.Error(w, "Unable to render preview", http.StatusInternalServerError)
		return
	}
}

func (h Handler) HandlerSetlistsCreate(w http.ResponseWriter, r *http.Request) {

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

	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)
	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		var sizeErr *http.MaxBytesError
		if errors.As(err, &sizeErr) {
			http.Error(w, "Upload too large", http.StatusRequestEntityTooLarge)
			return
		}
		http.Error(w, "Invalid multipart form", http.StatusBadRequest)
		return
	}
	defer r.MultipartForm.RemoveAll()

	input := services.CreateSetlistInput{
		Title:     r.FormValue("setlist-title"),
		Notes:     r.FormValue("notes"),
		TempArtID: r.FormValue("temporary-artwork-id"),
	}

	setlist, err := h.Services.SetlistCreate(r.Context(), auth.User, auth.CurrentBand, input)
	if err != nil {
		log.Println("   Unable to create setlist: ", err)
		http.Error(w, "Unable to create setlist", http.StatusInternalServerError)
		return
	}

	notif := models.Notification{
		BandID:    &band.BandID,
		UserID:    user.UserID,
		Type:      models.NotificationSetlistCreated,
		ItemID:    setlist.SetlistID,
		Body:      models.NotifBodySetlistCreated,
		CreatedBy: user.UserID,
	}

	err = h.Services.CreateNotification(r.Context(), notif)
	if err != nil {
		slog.Error(
			"unable to create notification",
			"request_id", requestlog.GetRequestID(r.Context()),
			"type", notif.Type,
			"error", err,
		)
	}

	log.Println("   Created setlist: ", input.Title)

	w.Header().Set("HX-Redirect", "/setlists")
	w.WriteHeader(http.StatusSeeOther)
	return
}

func (h Handler) HandlerSetlistEditPage(w http.ResponseWriter, r *http.Request) {

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
	setlistID := r.URL.Query().Get("id")
	if setlistID == "" {
		http.Error(w, "Setlist ID is required", http.StatusBadRequest)
		return
	}

	setlist, err := database.SetlistsTableGetSetlistByIDAndUserID(setlistID, user.UserID, band.BandID)
	if err != nil {
		log.Println("   Unable to get setlist: ", err)
		http.Error(w, "Unable to load setlist", http.StatusNotFound)
		return
	}
	if setlist.BandID != band.BandID {
		http.Error(w, "Setlist does not belong to the current band", http.StatusForbidden)
		return
	}

	notifications, err := database.NotifTableGetNotificationsByUserID(r.Context(), user.UserID, band.BandID)
	if err != nil {
		slog.Error(
			"unable to load notifications",
			"request_id", requestlog.GetRequestID(r.Context()),
			"error", err,
		)
		http.Error(w, "Unable to load notifications", http.StatusInternalServerError)
		return
	}

	data := models.SetlistPage{
		CSRFToken:     csrf.Token(r),
		User:          user,
		Band:          band,
		Setlist:       setlist,
		Notifications: notifications,
	}

	err = h.Tmpl.ExecuteTemplate(w, "setlist-edit.html", data)
	if err != nil {
		log.Println("   Unable to render setlist edit page: ", err)
		http.Error(w, "Unable to load setlist edit page", http.StatusInternalServerError)
		return
	}
	return
}

func (h Handler) HandlerSetlistUpdate(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

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

	r.Body = http.MaxBytesReader(w, r.Body, 12<<20)

	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}
	defer r.MultipartForm.RemoveAll()

	setlistID := r.FormValue("setlist-id")
	if setlistID == "" {
		http.Error(w, "Setlist ID is required", http.StatusBadRequest)
		return
	}

	existingSetlist, err := database.SetlistsTableGetSetlistByIDAndUserID(setlistID, user.UserID, band.BandID)
	if err != nil {
		log.Println("   Unable to get setlist: ", err)
		http.Error(w, "Setlist not found", http.StatusNotFound)
		return
	}
	if existingSetlist.BandID != band.BandID {
		http.Error(w, "Setlist does not belong to the current band", http.StatusForbidden)
		return
	}

	title := strings.TrimSpace(r.FormValue("setlist-title"))
	if title == "" {
		http.Error(w, "Setlist title is required", http.StatusBadRequest)
		return
	}
	slug := helpers.MakeSlug(title)
	notes := strings.TrimSpace(r.FormValue("notes"))

	artworkPath := existingSetlist.ArtworkPath
	imageID := existingSetlist.ArtworkID
	temporaryArtworkID := r.FormValue("temporary-artwork-id")

	if temporaryArtworkID != "" {
		artworkPath, err = h.Services.ServiceCreatePermSetlistImage(r.Context(), temporaryArtworkID, band.Slug, slug)
		if err != nil {
			slog.Error(
				"unable to save permanent setlist images",
				"request_id", requestlog.GetRequestID(r.Context()),
				"error", err,
			)
			http.Error(w, "Could not save artwork versions", http.StatusInternalServerError)
			return
		}
		imageID = temporaryArtworkID
	} else if r.FormValue("remove-artwork") == "true" {
		imageID = ""
		artworkPath = ""
	}

	setlist := models.Setlist{
		SetlistID:   setlistID,
		BandID:      band.BandID,
		Name:        title,
		Slug:        slug,
		Explicit:    existingSetlist.Explicit,
		Notes:       notes,
		ArtworkID:   imageID,
		ArtworkPath: artworkPath,
		UpdatedBy:   auth.User.UserID,
	}

	updated, err := database.SetlistsTableUpdateSetlist(setlist, user.UserID)
	if err != nil {
		log.Println("   Unable to update setlist in database: ", err)
		http.Error(w, "Unable to update setlist", http.StatusInternalServerError)
		return
	}
	if !updated {
		http.Error(w, "Setlist not found", http.StatusNotFound)
		return
	}

	log.Println("   Updated setlist: ", setlist.Name)

	w.Header().Set("HX-Redirect", "/setlist?id="+setlistID)
	w.WriteHeader(http.StatusOK)
	return
}

func (h Handler) HandlerSetlistsDelete(w http.ResponseWriter, r *http.Request) {

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

	setlistID := r.FormValue("setlist-id")

	setlist, err := database.SetlistsTableGetSetlistByIDAndUserID(setlistID, user.UserID, band.BandID)
	if err != nil {
		http.Error(w, "Setlist not found", http.StatusNotFound)
		return
	}

	err = database.SetlistsTableDeleteSetlist(setlistID, band.BandID, auth.User.UserID)
	if err != nil {
		log.Printf("\n   Unable to delete setlist %s from setlists table: %v", setlistID, err)
		http.Error(w, "Setlist not found", http.StatusNotFound)
		return
	}

	err = h.Services.ServiceDeleteSetlistImageVersions(r.Context(), setlist.ArtworkID, band.Slug, setlist.Slug)
	if err != nil {
		log.Println("   Unable to delete setlist artwork from R2: ", err)
	}

	http.Redirect(w, r, "/setlists", http.StatusSeeOther)
	return
}

func (h Handler) HandlerSetlistAddSong(w http.ResponseWriter, r *http.Request) {

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

	setlistID := r.FormValue("setlist-id")
	songID := r.FormValue("song-id")

	if setlistID == "" || songID == "" {
		log.Print("   Invalid setlistID or SongID: ")
		return
	}

	itemType := models.SetlistItemSong

	item, err := database.SetlistItemsTableSaveItem(itemType, songID, user.UserID, setlistID, band.BandID)
	if err != nil {
		log.Println("   Unable to save song to setlist: ", err)
		http.Error(w, "Unable to save song to setlist", http.StatusInternalServerError)
		return
	}

	log.Printf("   Saved item ID: %s of type %v to setlist id: %s\n", item.ItemID, item.ItemType, setlistID)

	return
}

func (h Handler) HandlerSetlistPage(w http.ResponseWriter, r *http.Request) {

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
	setlistID := r.URL.Query().Get("id")

	setlist, err := database.SetlistsTableGetSetlistByIDAndUserID(setlistID, user.UserID, band.BandID)
	if err != nil {
		log.Println("   Unable to get setlist: ", err)
		return
	}

	members, err := database.BandMembersGetMemberNameAndID(user.UserID, band.BandID)
	if err != nil {
		log.Println("   Unable to get band members: ", err)
		http.Error(w, "Unable to get band members", http.StatusInternalServerError)
		return
	}

	data := models.SetlistPage{
		CSRFToken: csrf.Token(r),
		User:      user,
		Band:      band,
		Members:   members,
		Setlist:   setlist,
	}

	err = h.Tmpl.ExecuteTemplate(w, "setlist.html", data)
	if err != nil {
		log.Println("   Unable to render setlist: ", err)
		http.Error(w, "Unable to load setlist", http.StatusInternalServerError)
		return
	}
	return
}

func (h Handler) HandlerSetlistInfoPage(w http.ResponseWriter, r *http.Request) {
	auth, err := HelperGetAuthContext(r)
	if err != nil {
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	setlist, err := database.SetlistsTableGetSetlistByIDAndUserID(
		r.URL.Query().Get("id"),
		auth.User.UserID,
		auth.CurrentBand.BandID,
	)
	if err != nil {
		http.Error(w, "Unable to load setlist", http.StatusInternalServerError)
		return
	}

	data := models.SetlistPage{
		CSRFToken: csrf.Token(r),
		User:      auth.User,
		Band:      auth.CurrentBand,
		Setlist:   setlist,
	}

	if err := h.Tmpl.ExecuteTemplate(w, "setlist-info.html", data); err != nil {
		http.Error(w, "Unable to load setlist info", http.StatusInternalServerError)
	}
}

func (h Handler) HandlerSetlistReorder(w http.ResponseWriter, r *http.Request) {

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
	setlistID := r.URL.Query().Get("id")

	setlist, err := database.SetlistsTableGetSetlistByIDAndUserID(setlistID, user.UserID, band.BandID)
	if err != nil {
		log.Println("   Unable to get setlist: ", err)
		return
	}

	data := models.SetlistPage{
		User:    user,
		Band:    band,
		Setlist: setlist,
	}

	err = h.Tmpl.ExecuteTemplate(w, "setlist-reorder", data)
	if err != nil {
		log.Println("   Unable to render setlist-reorder.html: ", err)
		http.Error(w, "Unable to load setlist-reorder.html", http.StatusInternalServerError)
		return
	}
}

func (h Handler) HandlerSetlistReorderSave(w http.ResponseWriter, r *http.Request) {

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

	setlistID := r.FormValue("setlistID")
	orderJSON := r.FormValue("order")

	var order []models.ReorderItem

	err = json.Unmarshal([]byte(orderJSON), &order)
	if err != nil {
		log.Println("   Unable to decode reorder response body: ", err)
		http.Error(w, "Invalid reorder data", http.StatusBadRequest)
		return
	}

	err = database.SetlistItemsUpdateOrder(setlistID, order, auth.CurrentBand.BandID, auth.User.UserID)
	if err != nil {
		log.Println("   Unable to save new order to db: ", err)
		http.Error(w, "Unable to save new order to database", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/setlist?id="+setlistID, http.StatusSeeOther)
	return
}

func (h Handler) HandlerSetlistDeleteSong(w http.ResponseWriter, r *http.Request) {

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

	setlistID := r.FormValue("setlist-id")
	songID := r.FormValue("song-id")
	position, err := strconv.Atoi(r.FormValue("position"))
	if err != nil {
		log.Println("   Unable to get song position: ", err)
		return
	}

	if setlistID == "" || songID == "" {
		log.Print("   Invalid setlistID or SongID: ")
		return
	}

	err = database.SetlistItemsTableDeleteSong(
		songID,
		position,
		setlistID,
		auth.CurrentBand.BandID,
		auth.User.UserID,
	)
	if err != nil {
		log.Println("   Unable to delete song to setlist: ", err)
		http.Error(w, "Unable to delete song to setlist", http.StatusInternalServerError)
		return
	}

	redirectURL := "/setlist?id=" + setlistID
	w.Header().Set("HX-Redirect", redirectURL)
	w.WriteHeader(http.StatusOK)
}

func (h Handler) HandlerSetlistDeleteTransition(w http.ResponseWriter, r *http.Request) {

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		slog.Error(
			"unable to load auth context",
			"request_id", requestlog.GetRequestID(r.Context()),
			"error", err,
		)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	setlistID := r.FormValue("setlist-id")
	transitionID := r.FormValue("transition-id")

	position, err := strconv.Atoi(r.FormValue("position"))
	if err != nil {
		log.Println("   Unable to get transition position: ", err)
		return
	}

	if setlistID == "" || transitionID == "" {
		log.Print("   Invalid setlistID or transitionID: ")
		return
	}

	err = database.SetlistItemsTableDeleteTransition(
		transitionID,
		position,
		setlistID,
		auth.CurrentBand.BandID,
		auth.User.UserID,
	)
	if err != nil {
		log.Println("   Unable to delete transition from setlist: ", err)
		http.Error(w, "Unable to delete transition from setlist", http.StatusInternalServerError)
		return
	}

	redirectURL := "/setlist?id=" + setlistID
	w.Header().Set("HX-Redirect", redirectURL)
	w.WriteHeader(http.StatusOK)
	return
}

func (h Handler) HandlerSetlistDeleteBreak(w http.ResponseWriter, r *http.Request) {
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

	setlistID := strings.TrimSpace(r.FormValue("setlist-id"))
	setlistItemID := strings.TrimSpace(r.FormValue("setlist-item-id"))
	if setlistID == "" || setlistItemID == "" {
		http.Error(w, "Invalid setlist break", http.StatusBadRequest)
		return
	}

	err = database.SetlistItemsTableDeleteBreak(
		setlistItemID,
		setlistID,
		auth.CurrentBand.BandID,
		auth.User.UserID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.NotFound(w, r)
			return
		}
		slog.Error(
			"unable to delete break from setlist",
			"request_id", requestlog.GetRequestID(r.Context()),
			"setlist_id", setlistID,
			"setlist_item_id", setlistItemID,
			"error", err,
		)
		http.Error(w, "Unable to delete break from setlist", http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Redirect", "/setlist?id="+url.QueryEscape(setlistID))
	w.WriteHeader(http.StatusOK)
}

func (h Handler) HandlerSetlistsSearch(w http.ResponseWriter, r *http.Request) {

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

	query := r.FormValue("q")

	_, err = database.SetlistsTableSearchSetlistByBandIDAndUserID(band.BandID, user.UserID, query)
	if err != nil {
		log.Println("   Error searching songs by Band ID: ", band.BandID)
		http.Error(w, "Could not search songs", http.StatusInternalServerError)
		return
	}
}

func (h Handler) HandlerSetlistDuplicate(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerSetlistDuplicate")

}

func (h Handler) HandlerSetlistShare(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerSetlistShare")

}

func (h Handler) HandlerPrintSetist(w http.ResponseWriter, r *http.Request) {

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

	setlistID := r.PathValue("id")

	options := models.SetlistPDFOptions{
		Numbered:    r.FormValue("numbered") == "true",
		Songs:       r.FormValue("songs") == "true",
		Transitions: r.FormValue("transitions") == "true",
		Breaks:      r.FormValue("breaks") == "true",
		Length:      r.FormValue("length") == "true",
		Tempo:       r.FormValue("tempo") == "true",
		Key:         r.FormValue("key") == "true",
		Tuning:      r.FormValue("tuning") == "true",
		Capo:        r.FormValue("capo") == "true",
	}

	fmt.Printf("%+v\n", options)

	setlist, err := database.SetlistsTableGetSetlistByIDAndUserID(setlistID, user.UserID, band.BandID)
	if err != nil {
		slog.Error(
			"unable to load selist",
			"request_id", requestlog.GetRequestID(r.Context()),
			"path", r.URL.Path,
			"error", err,
		)
		http.Error(w, "Unable to load setlist", http.StatusInternalServerError)
		return
	}

	fmt.Println("user: ", user.UserID)
	fmt.Println("band: ", band.BandID)
	fmt.Println("setlistID: ", setlistID)

	pdfBytes, err := h.Services.GenerateSetlistPDF(setlist, options)
	if err != nil {
		slog.Error(
			"unable to generate pdf",
			"request_id", requestlog.GetRequestID(r.Context()),
			"path", r.URL.Path,
			"error", err,
		)
		http.Error(w, "Could not generate PDF", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set(
		"Content-Disposition",
		`attachment; filename="setlist.pdf"`,
	)
	w.Write(pdfBytes)

}
