package handlers

import (
	services "bandplan/src/Services"
	"bandplan/src/clients"
	"bandplan/src/database"
	"bandplan/src/models"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

func (h Handler) HandlerSongsITunesQueryPage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("------------------------------------------")
	log.Println("- HandlerSongsITunesSearch")

	user, err := HelperGetAuthenticatedUser(r)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	band, err := database.BandsTableGetBandByUserID(user.UserID)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	fmt.Println("   band name: ", band.Name)

	data := models.SongDownloadData{
		User: user,
		Band: band,
	}

	err = h.Tmpl.ExecuteTemplate(w, "song-download.html", data)
	if err != nil {
		log.Println("   Err getting songs-download page: ", err)
		http.Redirect(w, r, "/songs", http.StatusSeeOther)
	}
}

func (h Handler) HandlerSongsITunesQuery(w http.ResponseWriter, r *http.Request) {
	fmt.Println("------------------------------------------")
	log.Println("- HandlerSongsITunesQuery")

	user, err := HelperGetAuthenticatedUser(r)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	band, err := database.BandsTableGetBandByUserID(user.UserID)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	artistQuery := strings.TrimSpace(r.FormValue("itunes-query-artist-name"))
	songQuery := strings.TrimSpace(r.FormValue("itunes-query-song-title"))

	fmt.Println("   artistQuery: ", artistQuery)
	fmt.Println("   songQuery; ", songQuery)

	searchResponse, err := services.ServicesSearchITunesByArtistAndSong(artistQuery, songQuery)
	if err != nil {
		log.Println("   Unable to get reponse from iTunes API: ", err)
		http.Redirect(w, r, "/songs", http.StatusSeeOther)
	}

	fmt.Println("\nSearch Response: ")
	fmt.Println(searchResponse.ResultCount)
	fmt.Println("\nSearch Results")
	for _, r := range searchResponse.Results {
		fmt.Printf("%#v\n", r)
		fmt.Println("")
	}
	if len(searchResponse.Results) <= 0 {
		log.Println("   No results returned")
		http.Redirect(w, r, "/songs", http.StatusSeeOther)
		return
	}
	res := searchResponse.Results[0]

	length := res.TrackTimeMillis / 1000
	releaseDate, err := time.Parse(time.RFC3339, "2024-02-01T12:00:00Z")
	if err != nil {
		log.Println("unable to parse release date:", err)
	}

	var cover bool
	if res.ArtistName == band.Name {
		cover = false
	} else {
		cover = true
	}

	newSong := models.Song{
		BandID:         band.BandID,
		Title:          res.TrackName,
		ArtistName:     res.ArtistName,
		AlbumTitle:     res.CollectionName,
		ReleaseDate:    releaseDate,
		Genre:          res.PrimaryGenreName,
		LengthSeconds:  length,
		IsCover:        cover,
		AppleMusicLink: res.TrackViewURL,
	}

	fmt.Println(newSong)
	fmt.Println(res.ArtworkURL100)
	clients.ClientITunesSearchGetArtwork(res.ArtworkURL100)

	err = h.Tmpl.ExecuteTemplate(w, "songs-add-itunes.html", newSong)
	if err != nil {
		log.Println("   Unable to execute songs-add-itunes.html:", err)
	}

	// http.Redirect(w, r, "/songs", http.StatusSeeOther)
	w.Header().Set("HX-Redirect", "/songs")
	w.WriteHeader(http.StatusOK)
	return
}
