package main

import (
	"bandplan/src/database"
	"bandplan/src/handlers"
	"bandplan/src/middleware"
	"log"
	"net/http"
)

var messages []string

func handleAuth(pattern string, handler http.HandlerFunc) {
	http.Handle(pattern, middleware.RequireAuth(handler))
}

func main() {
	log.Println("MAIN")

	database.DB = database.ConnectDB()
	defer database.DB.Close()

	if err := database.DB.Ping(); err != nil {
		log.Fatal("Unable to connect to database:", err)
	}

	log.Println("   Database connection successful")

	tmpl := handlers.HelperParseTemplates()

	h := handlers.Handler{
		DB:   database.DB,
		Tmpl: tmpl,
	}

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	http.HandleFunc("/", h.HandlerHome)
	http.HandleFunc("/access", h.HandlerAccessCodePage)
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

	handleAuth("/songs", h.HandlerSongsPage)
	handleAuth("/songs/add", h.HandlerSongsAddPage)
	handleAuth("/songs/create", h.HandlerSongsAdd)
	handleAuth("/songs/search", h.HandlerSongsSearch)

	handleAuth("/songs/itunes/query", h.HandlerSongsITunesQueryPage)
	handleAuth("/songs/itunes/download", h.HandlerSongsITunesQuery)
	handleAuth("/songs/itunes/results", h.HandlerSongsITunesResults)
	handleAuth("/songs/itunes/create", h.HandlerSongsITunesResultsAddSong)

	handleAuth("/song", h.HandlerSongPage)
	handleAuth("/song/edit", h.HandlerSongEditPage)
	handleAuth("/song/update", h.HandlerSongUpdate)
	handleAuth("/song/delete", h.HandlerSongDelete)

	handleAuth("/setlists", h.HandlerSetlists)
	handleAuth("/setlists/add", h.HandlerSetlistsAddPage)
	handleAuth("/setlists/create", h.HandlerSetlistCreate)
	handleAuth("/setlist", h.HandlerSetlist)

	http.HandleFunc("/promotion", h.HandlerPromotion)
	http.HandleFunc("/goals", h.HandlerGoals)
	http.HandleFunc("/calendar", h.HandlerCalendar)
	http.HandleFunc("/events", h.HandlerEvents)
	http.HandleFunc("/files", h.HandlerFiles)

	http.HandleFunc("/profile-pic", h.HandlerProfilePicPage)
	http.HandleFunc("/profile-pic/add", h.HandlerProfilePicAdd)
	// http.HandleFunc("/admin", h.HandlerAdmin)
	http.Handle("/admin", middleware.RequireAuth(http.HandlerFunc(h.HandlerAdmin)))
	http.Handle("/admin/access-code", middleware.RequireAuth(http.HandlerFunc(h.HandlerCreateAccessCode)))
	http.HandleFunc("/settings", h.HandlerSettingsPage)

	log.Println("   Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
