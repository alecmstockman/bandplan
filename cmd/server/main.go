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
	http.Handle("/songs/add", middleware.RequireAuth(http.HandlerFunc(h.HandlerSongsAddPage)))
	http.Handle("/songs/create", middleware.RequireAuth(http.HandlerFunc(h.HandlerSongsAdd)))
	http.Handle("/songs/search", middleware.RequireAuth(http.HandlerFunc(h.HandlerSongsSearch)))

	http.Handle("/songs/itunes/query", middleware.RequireAuth(http.HandlerFunc(h.HandlerSongsITunesQueryPage)))
	http.Handle("/songs/itunes/download", middleware.RequireAuth(http.HandlerFunc(h.HandlerSongsITunesQuery)))
	http.Handle("/songs/itunes/results", middleware.RequireAuth(http.HandlerFunc(h.HandlerSongsITunesResults)))
	http.Handle("/songs/itunes/create", middleware.RequireAuth(http.HandlerFunc(h.HandlerSongsITunesResultsAddSong)))

	http.Handle("/song", middleware.RequireAuth(http.HandlerFunc(h.HandlerSongPage)))
	http.Handle("/song/edit", middleware.RequireAuth(http.HandlerFunc(h.HandlerSongEditPage)))
	http.Handle("/song/update", middleware.RequireAuth(http.HandlerFunc(h.HandlerSongUpdate)))
	http.Handle("/song/delete", middleware.RequireAuth(http.HandlerFunc(h.HandlerSongDelete)))

	http.Handle("/setlists", middleware.RequireAuth(http.HandlerFunc(h.HandlerSetlists)))
	http.Handle("/setlists/add", middleware.RequireAuth(http.HandlerFunc(h.HandlerSetlistsAddPage)))
	http.Handle("/setlists/create", middleware.RequireAuth(http.HandlerFunc(h.HandlerSetlistCreate)))
	http.Handle("/setlist", middleware.RequireAuth(http.HandlerFunc(h.HandlerSetlist)))

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
