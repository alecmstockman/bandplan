package handlers

import (
	"bandplan/src/database"
	"bandplan/src/models"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

func (h Handler) HandlerSetlistsPage(w http.ResponseWriter, r *http.Request) {
	log.Print("- HandlerSetlistsPage")

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		log.Println("   Unable to get AuthContext: ", err)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	user := auth.User
	band := auth.CurrentBand

	// setlists, err := database.SetlistsTableGetSetlistsByBandID(band.BandID)
	// if err != nil {
	// 	log.Println("   Unable to get settlists from db by band ID: ", nil)
	// 	http.Error(w, "Unable to get setlists", http.StatusInternalServerError)
	// }

	setlistSummaries, err := database.SetlistsTableGetSetlistSummariessByBandID(band.BandID)
	if err != nil {
		log.Println("   Unable to get setlist summaries band band id from db: ", err)
		http.Error(w, "Unable to get setlist summaries", http.StatusInternalServerError)
	}

	for set := range setlistSummaries {
		fmt.Println("Summary: ", set)
	}

	data := models.SetlistData{
		User:     user,
		Band:     band,
		Setlists: setlistSummaries,
	}

	err = h.Tmpl.ExecuteTemplate(w, "setlists.html", data)
	if err != nil {
		log.Println("   Unable to render setlists: ", err)
		http.Error(w, "Unable to load setlists", http.StatusInternalServerError)
		return
	}
}

func (h Handler) HandlerSetlist(w http.ResponseWriter, r *http.Request) {
	log.Print("- HandlerSetlist")

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		log.Println("   Unable to get AuthContext: ", err)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	user := auth.User
	band := auth.CurrentBand

	data := models.MenuPageData{
		User: user,
		Band: band,
	}

	err = h.Tmpl.ExecuteTemplate(w, "setlist.html", data)
	if err != nil {
		log.Println("   unable to get /setlist: ", err)
	}
	return
}

func (h Handler) HandlerSetlistsAddPage(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerSetlistsAddPage")

	log.Printf(
		"- HandlerSetlistsAddPage method=%s path=%s hx-request=%q",
		r.Method,
		r.URL.Path,
		r.Header.Get("HX-Request"),
	)

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

	data := models.SongDownloadData{
		User: user,
		Band: band,
	}

	err = h.Tmpl.ExecuteTemplate(w, "setlist-add.html", data)
	if err != nil {
		log.Println("   Err getting setlist-add page: ", err)
		http.Redirect(w, r, "/setlists", http.StatusSeeOther)
	}

	return
}

func (h Handler) HandlerSetlistsTempArt(w http.ResponseWriter, r *http.Request) {
	fmt.Println("--------------------------------------")
	log.Println("- HandlerSetlistsTempArt")

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		log.Println("   Unable to get AuthContext: ", err)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	band := auth.CurrentBand

	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Println("   Error parsing from while creating a new setlist: ", err)
		http.Error(w, "Unable to parse from", http.StatusBadRequest)
		return
	}

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
		previewURL, err = h.HelperSaveTempImage(r.Context(), file, imageID, band.Slug, "setlist")
		fmt.Println("previewURL: ", previewURL)
		if err != nil {
			http.Error(w, "Could not save artwork versions", http.StatusInternalServerError)
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

func (h Handler) HandlerSetlistsCreate(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerSetlistCreate")

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		log.Println("   Unable to get AuthContext: ", err)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	user := auth.User
	band := auth.CurrentBand

	title := strings.TrimSpace(r.FormValue("setlist-title"))
	slug := HelperMakeSlug(title)
	explicit := false
	notes := strings.TrimSpace(r.FormValue("notes"))

	tempArtID := strings.TrimSpace(
		r.FormValue("temporary-artwork-id"),
	)

	fmt.Println("\ntempArtID: ", tempArtID)

	artworkPath := ""

	browserPath, err := h.HelperCreatePermSetlistImage(r.Context(), tempArtID, band.Slug, slug)
	if err != nil {
		log.Println("   Unable to permanently save setlist art: ", err)
	} else {
		artworkPath = browserPath
	}

	setlist := models.Setlist{
		BandID:      band.BandID,
		Name:        title,
		Slug:        slug,
		Explicit:    explicit,
		Notes:       notes,
		ArtworkID:   tempArtID,
		ArtworkPath: artworkPath,
		CreatedBy:   user.UserID,
		UpdatedBy:   user.UserID,
	}

	err = database.SetlistsTableCreateSetlist(setlist)
	if err != nil {
		log.Println("   Unable to create setlist in database: ", err)
		http.Error(w, "/setlists", http.StatusInternalServerError)
		return
	}

	log.Println("   Created setlist: ", setlist.Name)

	w.Header().Set("HX-Redirect", "/setlists")
	w.WriteHeader(http.StatusSeeOther)
	return
}

func (h Handler) HandlerSetlistEditPage(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerSetlistEdit")

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

	err = h.Tmpl.ExecuteTemplate(w, "setlist-edit.html", data)
	if err != nil {
		log.Println("   Unable to render setlist edit page: ", err)
		http.Error(w, "Unable to load setlist edit page", http.StatusInternalServerError)
		return
	}
	return
}

func (h Handler) HandlerSetlistUpdate(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerTransitionUpdate")
	auth, err := HelperGetAuthContext(r)
	if err != nil {
		log.Println("   Unable to get AuthContext: ", err)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	user := auth.User
	band := auth.CurrentBand

	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		http.Error(w, "File too large", http.StatusBadRequest)
		return
	}

	title := strings.TrimSpace(r.FormValue("setlist-title"))
	slug := HelperMakeSlug(title)
	explicit := false
	notes := strings.TrimSpace(r.FormValue("notes"))

	artworkPath := ""
	imageID := ""

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
		artworkPath, err = h.HelperSaveSetlistImageVersions(r.Context(), file, imageID, band.Slug, slug)
		if err != nil {
			http.Error(w, "Could not save artwork versions", http.StatusInternalServerError)
			return
		}
	}

	setlist := models.Setlist{
		BandID:      band.BandID,
		Name:        title,
		Slug:        slug,
		Explicit:    explicit,
		Notes:       notes,
		ArtworkID:   imageID,
		ArtworkPath: artworkPath,
		CreatedBy:   user.UserID,
		UpdatedBy:   user.UserID,
	}

	err = database.SetlistsTableUpdateSetlist(setlist)
	if err != nil {
		log.Println("   Unable to update setlist in database: ", err)
		http.Error(w, "/setlists", http.StatusInternalServerError)
		return
	}

	log.Println("   Updated setlist: ", setlist.Name)

	w.Header().Set("HX-Redirect", "/setlists")
	w.WriteHeader(http.StatusSeeOther)
	return
}

func (h Handler) HandlerSetlistsDelete(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerSetlistsDelete")

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		log.Println("   Unable to get AuthContext: ", err)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	band := auth.CurrentBand

	setlistID := r.FormValue("setlist-id")
	imageID := r.FormValue("setlist-image-id")
	slug := r.FormValue("setlist-slug")

	err = database.SetlistsTableDeleteSetlist(setlistID)
	if err != nil {
		log.Printf("\n   Unable to delete setlist %s from setlists table: %v", setlistID, err)
	}

	err = h.HelperDeleteSetlistImageVersions(r.Context(), imageID, band.Slug, slug)
	if err != nil {
		log.Println("   Unable to delete setlist artwork from R2: ", err)
	}

	http.Redirect(w, r, "/setlists", http.StatusSeeOther)
	return
	// w.Header().Set("HX-Redirect", "/setlists")
	// w.WriteHeader(http.StatusSeeOther)
	// return
}

func (h Handler) HandlerSetlistAddSong(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerSetlistAddSong")

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		log.Println("   Unable to get AuthContext: ", err)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		w.Header().Set("HX-Redirect", "/")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	user := auth.User

	setlistID := r.FormValue("setlist-id")
	songID := r.FormValue("song-id")

	if setlistID == "" || songID == "" {
		log.Print("   Invalid setlistID or SongID: ")
		return
	}

	itemType := models.SetlistItemSong

	song, err := database.SetlistItemsTableSaveItem(itemType, songID, user.UserID, setlistID)
	if err != nil {
		log.Println("   Unable to save song to setlist: ", err)
		http.Error(w, "Unable to save song to setlist", http.StatusInternalServerError)
		return
	}
	log.Printf("   Saved %s to setlist_items\n", song.Song.Title)

	return
}

func (h Handler) HandlerSetlistPage(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerSetlistPage")

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

	fmt.Println("\n=====================")
	fmt.Printf("%+v\n", setlist)
	// fmt.Println("Setlist:")
	// fmt.Println("SetlistID: ", setlist.SetlistID)
	// fmt.Println("Name; ", setlist.Name)
	// for _, song := range setlist.Songs {
	// 	fmt.Println(song.ItemType)
	// 	fmt.Println("Item: ", song.Transition)
	// 	fmt.Printf("%+v\n", song.Transition)
	// 	fmt.Println(song.Transition.Title)
	// 	fmt.Println("")
	// }
	fmt.Println("=====================\n")

	data := models.SetlistPage{
		User:    user,
		Band:    band,
		Setlist: setlist,
	}

	err = h.Tmpl.ExecuteTemplate(w, "setlist.html", data)
	if err != nil {
		log.Println("   Unable to render setlist: ", err)
		http.Error(w, "Unable to load setlist", http.StatusInternalServerError)
		return
	}

	return
}

func (h Handler) HandlerSetlistDeleteSong(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerSetlistDeleteSong")

	_, err := HelperGetAuthContext(r)
	if err != nil {
		log.Println("   Unable to get AuthContext: ", err)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		w.Header().Set("HX-Redirect", "/")
		http.Redirect(w, r, "/", http.StatusSeeOther)
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

	err = database.SetlistItemsTableDeleteSong(songID, position, setlistID)
	if err != nil {
		log.Println("   Unable to delete song to setlist: ", err)
		http.Error(w, "Unable to delete song to setlist", http.StatusInternalServerError)
		return
	}

	redirectURL := "/setlist?id=" + url.QueryEscape(setlistID)
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

func (h Handler) HandlerBreakPage(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerBreakPage")
	return
}

func (h Handler) HandlerBreakAdd(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerBreakAdd")
	return
}

func (h Handler) HandlerBreakCreate(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerBreakCreate")
	return
}

func (h Handler) HandlerDeleteBreak(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerDeleteBreak")
	return
}
