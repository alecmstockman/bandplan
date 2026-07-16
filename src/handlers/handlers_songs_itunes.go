package handlers

import (
	"bandplan/src/database"
	"bandplan/src/models"
	services "bandplan/src/services"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

func (h Handler) HandlerSongsITunesQueryPage(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerSongsITunesSearch")

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		log.Println("   Unable to get AuthContext: ", err)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	user := auth.User
	band := auth.CurrentBand

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
	log.Println("- HandlerSongsITunesQuery")

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		log.Println("   Unable to get AuthContext: ", err)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	band := auth.CurrentBand

	artistQuery := strings.TrimSpace(r.FormValue("itunes-query-artist-name"))
	songQuery := strings.TrimSpace(r.FormValue("itunes-query-song-title"))

	searchResponse, err := services.ServicesSearchITunesByArtistAndSong(artistQuery, songQuery)
	if err != nil {
		log.Println("   Unable to get reponse from iTunes API: ", err)
		http.Redirect(w, r, "/songs", http.StatusSeeOther)
	}

	if len(searchResponse.Results) <= 0 {
		log.Println("   No results returned")
		http.Redirect(w, r, "/songs", http.StatusSeeOther)
		return
	}

	// fmt.Println("\n RESULT COUNT: ", len(searchResponse.Results))
	// for _, s := range searchResponse.Results {
	// 	fmt.Println("\nTrackID:", s.TrackID)
	// 	fmt.Println("ArtistName:", s.ArtistName)
	// 	fmt.Println("TrackName:", s.TrackName)
	// 	fmt.Println("CollectionName:", s.CollectionName)
	// 	fmt.Println("ArtworkURL100:", s.ArtworkURL100)
	// 	fmt.Println("TrackViewURL:", s.TrackViewURL)
	// 	fmt.Println("PreviewURL:", s.PreviewURL)
	// 	fmt.Println("TrackTimeMillis:", s.TrackTimeMillis)
	// 	fmt.Println("PrimaryGenreName:", s.PrimaryGenreName)
	// 	fmt.Println("ReleaseDate:", s.ReleaseDate)
	// 	fmt.Println("\n")
	// }

	// res := searchResponse.Results[0]

	if len(searchResponse.Results) == 1 {
		fmt.Println("ONE RESULT: -----------------")
		res := searchResponse.Results[0]

		length := res.TrackTimeMillis / 1000

		releaseDate, err := time.Parse(time.RFC3339, "2024-02-01T12:00:00Z")
		if err != nil {
			log.Println("   unable to parse release date:", err)
		}

		var cover bool
		if res.ArtistName == band.Name {
			cover = false
		} else {
			cover = true
		}

		imageID := uuid.New().String()

		artworkPath, err := HelperSaveArtworkImageFromITunes(res.ArtworkURL100, imageID)
		if err != nil {
			log.Println("   Unable to save artwork image from iTunes: ", err)
			imageID = ""
			artworkPath = ""
		}

		newSong := models.Song{
			BandID:         band.BandID,
			Title:          res.TrackName,
			ArtistName:     res.ArtistName,
			AlbumTitle:     res.CollectionName,
			ArtworkID:      imageID,
			ArtworkPath:    artworkPath,
			ReleaseDate:    releaseDate,
			Genre:          res.PrimaryGenreName,
			LengthSeconds:  length,
			IsCover:        cover,
			AppleMusicLink: res.TrackViewURL,
		}

		err = h.Tmpl.ExecuteTemplate(w, "songs-add-itunes.html", newSong)
		if err != nil {
			log.Println("   Unable to execute songs-add-itunes.html:", err)
		}

		w.Header().Set("HX-Redirect", "/songs")
		w.WriteHeader(http.StatusOK)
		return

	} else if len(searchResponse.Results) > 1 {
		fmt.Println("MULTIPLE RESULTS: -----------------")
		var resList models.ITunesSearchResponse

		for _, s := range searchResponse.Results {

			newSong := models.ITunesSong{
				TrackID:         s.TrackID,
				TrackName:       s.TrackName,
				ArtistName:      s.ArtistName,
				CollectionName:  s.CollectionName,
				ArtworkURL100:   s.ArtworkURL100,
				TrackViewURL:    s.TrackViewURL,
				PreviewURL:      s.PreviewURL,
				TrackTimeMillis: s.TrackTimeMillis,
			}

			resList.Results = append(resList.Results, newSong)
		}

		err = h.Tmpl.ExecuteTemplate(w, "songs-itunes-results.html", resList)
		if err != nil {
			log.Println("   Unable to execute songs-itunes-results.html:", err)
			return
		}

		w.Header().Set("HX-Redirect", "/songs")
		w.WriteHeader(http.StatusOK)
		return

	} else {
		fmt.Println("NO RESULTS: -----------------")
		w.Header().Set("HX-Redirect", "/songs")
		w.WriteHeader(http.StatusOK)
		return
	}

}

func (h Handler) HandlerSongsITunesResults(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerSongsITunesResults")

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		log.Println("   Unable to get AuthContext: ", err)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	user := auth.User
	band := auth.CurrentBand

	data := models.SongDownloadData{
		User: user,
		Band: band,
	}

	err = h.Tmpl.ExecuteTemplate(w, "songs-itunes-results.html", data)
	if err != nil {
		log.Println("   Err getting songs-download page: ", err)
		http.Redirect(w, r, "/songs", http.StatusSeeOther)
	}
}

func (h Handler) HandlerSongsITunesResultsAddSong(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerSongsItunesResultsAddSong")

	itunesID := r.FormValue("track-id")

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		log.Println("   Unable to get AuthContext: ", err)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	user := auth.User
	band := auth.CurrentBand

	res, err := services.ServicesSearchITunesByITunesID(itunesID)
	if err != nil {
		log.Println("   Unable to get song from iTunes by id:", err)
		http.Redirect(w, r, "/songs", http.StatusSeeOther)
		return
	}

	song := res.Results[0]

	length := song.TrackTimeMillis / 1000
	releaseDate, err := time.Parse(time.RFC3339, "2024-02-01T12:00:00Z")
	if err != nil {
		log.Println("   unable to parse release date:", err)
	}

	var cover bool
	if song.ArtistName == band.Name {
		cover = false
	} else {
		cover = true
	}

	imageID := uuid.New().String()

	artworkPath, err := HelperSaveArtworkImageFromITunes(song.ArtworkURL100, imageID)
	if err != nil {
		log.Println("   Unable to save artwork image from iTunes: ", err)
		imageID = ""
		artworkPath = ""
	}

	newSong := models.Song{
		BandID: band.BandID,

		Title:      song.TrackName,
		TitleSlug:  "",
		ArtistName: song.ArtistName,
		AlbumTitle: song.CollectionName,
		AlbumSlug:  "",

		ArtworkID:     imageID,
		ArtworkPath:   artworkPath,
		ReleaseDate:   releaseDate,
		Genre:         song.PrimaryGenreName,
		LengthSeconds: length,

		IsCover: cover,

		CreatedBy: user.UserID,
		UpdatedBy: user.UserID,
	}

	_, err = database.SongsTableCreateSong(newSong)
	if err != nil {
		log.Println("   Unable to to create song in database: ", err)
		return
	}

	return
}
