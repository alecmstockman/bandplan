package main

import (
	"bandplan/src/database"
	"bandplan/src/handlers"
	"bandplan/src/middleware"
	"log"
	"net/http"
)

var messages []string

func main() {
	log.Println("MAIN")

	database.DB = database.ConnectDB()
	defer database.DB.Close()

	tmpl := handlers.HelperParseTemplates()

	h := handlers.Handler{
		DB:   database.DB,
		Tmpl: tmpl,
	}

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Println("   Connecting to database...")

	database.CreateBandsTable(database.DB)
	database.CreateUsersTable(database.DB)
	database.CreateBandMembersTable(database.DB)
	database.CreateMessagesTable(database.DB)
	database.CreateSesssionsTable(database.DB)
	database.CreateSongsTable(database.DB)
	database.CreateSetlistsTable(database.DB)
	database.CreateSetlistSongsTable(database.DB)

	log.Println("   Database connection succesful")

	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	http.HandleFunc("/", h.HandlerHome)
	http.HandleFunc("/register", h.HandlerRegisterPage)
	http.HandleFunc("/register/create", h.HandlerRegister)
	http.HandleFunc("/register/user-agreement", h.HandlerUserAgreementPage)
	http.HandleFunc("/register/user-agreed", h.HandlerUserAgreement)
	http.HandleFunc("/terms", h.HandlerTermsPage)
	http.HandleFunc("/privacy", h.HandlerPrivacyPage)

	http.HandleFunc("/login", h.HandlerLoginPage)
	http.HandleFunc("/login/enter", h.HandlerLogin)
	http.HandleFunc("/logout", h.HandlerLogout)

	http.HandleFunc("/send", h.HandlerSend)
	http.HandleFunc("/delete", h.HandlerDelete)
	http.HandleFunc("/messages", h.HandlerMessages)

	// http.HandleFunc("/songs", h.HandlerSongsPage)
	http.Handle("/songs", middleware.RequireAuth(http.HandlerFunc(h.HandlerSongsPage)))
	// http.HandleFunc("/songlist", h.HandlerSongs)
	// http.HandleFunc("/songs/add", h.HandlerSongsAddPage)
	http.Handle("/songs/add", middleware.RequireAuth(http.HandlerFunc(h.HandlerSongsAddPage)))
	// http.HandleFunc("/songs/create", h.HandlerSongsAdd)
	http.Handle("/songs/create", middleware.RequireAuth(http.HandlerFunc(h.HandlerSongsAdd)))
	// http.HandleFunc("/songs/search", h.HandlerSongsSearch)
	http.Handle("/songs/search", middleware.RequireAuth(http.HandlerFunc(h.HandlerSongsSearch)))

	http.HandleFunc("/songs/itunes/query", h.HandlerSongsITunesQueryPage)
	http.HandleFunc("/songs/itunes/download", h.HandlerSongsITunesQuery)
	http.HandleFunc("/songs/itunes/results", h.HandlerSongsITunesResults)
	http.HandleFunc("/songs/itunes/create", h.HandlerSongsITunesResultsAddSong)

	// http.HandleFunc("/song", h.HandlerSongPage)
	http.Handle("/song", middleware.RequireAuth(http.HandlerFunc(h.HandlerSongPage)))
	// http.HandleFunc("/song/edit", h.HandlerSongEditPage)
	http.Handle("/song/edit", middleware.RequireAuth(http.HandlerFunc(h.HandlerSongEditPage)))
	// http.HandleFunc("/song/update", h.HandlerSongUpdate)
	http.Handle("/song/update", middleware.RequireAuth(http.HandlerFunc(h.HandlerSongUpdate)))
	// http.HandleFunc("/song/delete", h.HandlerSongDelete)
	http.Handle("/song/delete", middleware.RequireAuth(http.HandlerFunc(h.HandlerSongDelete)))

	http.Handle("/setlists", middleware.RequireAuth(http.HandlerFunc(h.HandlerSetlists)))
	http.HandleFunc("/setlists/add", h.HandlerSetlistsAddPage)
	http.HandleFunc("/setlists/create", h.HandlerSetlistsAddPage)

	http.HandleFunc("/setlist", h.HandlerSetlist)

	http.HandleFunc("/promotion", h.HandlerPromotion)
	http.HandleFunc("/goals", h.HandlerGoals)
	http.HandleFunc("/calendar", h.HandlerCalendar)
	http.HandleFunc("/events", h.HandlerEvents)
	http.HandleFunc("/files", h.HandlerFiles)

	http.HandleFunc("/profile-pic", h.HandlerProfilePicPage)
	http.HandleFunc("/profile-pic/add", h.HandlerProfilePicAdd)

	log.Println("   Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
