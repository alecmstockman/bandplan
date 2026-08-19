package main

import (
	"bandplan/src/database"
	"bandplan/src/handlers"
	"bandplan/src/middleware"
	"bandplan/src/realtime"
	"bandplan/src/storage"
	"context"
	"log"
	"net/http"
	"os"
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

	log.Println("Database connection successful")

	tmpl := handlers.HelperParseTemplates()

	r2Storage, err := storage.NewR2Storage(context.Background())
	if err != nil {
		log.Fatal("Unable to initialize R2 storage: ", err)
	}

	hub := realtime.NewHub()
	go hub.Run()

	h := handlers.Handler{
		DB:      database.DB,
		Tmpl:    tmpl,
		Storage: r2Storage,
		Hub:     hub,
	}

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.WriteHeader(http.StatusOK)

		if _, err := w.Write([]byte("ok")); err != nil {
			log.Println("Unable to write health-check response:", err)
		}
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

	http.HandleFunc("/delete", h.HandlerDelete)
	http.HandleFunc("/messages", h.HandlerMessages)

	handleAuth("/ws/chat", h.HandlerChatWebSocket)

	handleAuth("/chats", h.HandlerChatsPage)
	handleAuth("/chat", h.HandlerChatPage)

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
	handleAuth("/song/lyrics", h.HandlerSongLyrics)

	handleAuth("/setlists", h.HandlerSetlistsPage)
	handleAuth("/setlists/add", h.HandlerSetlistsAddPage)
	handleAuth("/setlists/temp-art/add", h.HandlerSetlistsTempArt)
	handleAuth("/setlists/temp-art/delete", h.HandlerSetlistsTempArtDelete)
	handleAuth("/setlists/create", h.HandlerSetlistsCreate)
	handleAuth("/setlists/delete", h.HandlerSetlistsDelete)

	handleAuth("/setlist", h.HandlerSetlistPage)
	handleAuth("/setlist/add", h.HandlerSetlistAddSong)
	handleAuth("/setlist/edit", h.HandlerSetlistEditPage)
	handleAuth("/setlist/update", h.HandlerSetlistUpdate)
	handleAuth("/setlist/delete", h.HandlerSetlistDeleteSong)
	handleAuth("/setlist/transition/delete", h.HandlerSetlistDeleteTransition)

	handleAuth("/setlist/songs", h.HandlerSetlistSongs)
	handleAuth("/setlist/transitions", h.HandlerSetlistTransitions)
	handleAuth("/setlist/breaks", h.HandlerSetlistBreaks)
	handleAuth("/setlist/items", h.HandlerSetlistItems)

	handleAuth("/setlist/delete/song", h.HandlerSetlistItems)
	handleAuth("/setlist/delete/transition", h.HandlerSetlistDeleteTransition)
	handleAuth("/setlist/delete/break", h.HandlerSetlistItems)

	handleAuth("/setlist/reorder", h.HandlerSetlistReorder)
	handleAuth("/setlist/reorder/save", h.HandlerSetlistReorderSave)

	handleAuth("/setlist/time/songs", h.HandlerSetlistUpdateTimeButtonSongs)
	handleAuth("/setlist/time/transitions", h.HandlerSetlistUpdateTimeButtonTransitions)
	handleAuth("/setlist/time/breaks", h.HandlerSetlistUpdateTimeButtonBreaks)
	handleAuth("/setlist/time/items", h.HandlerSetlistUpdateTimeButtonItems)

	handleAuth("/setlist/notes", h.HandlerSetlistOpenNotesPage)
	handleAuth("/setlist/notes/edit", h.HandlerSetlistEditNotesPage)
	handleAuth("/setlist/notes/save", h.HandlerSetlistSaveNotesPage)

	handleAuth("/setlist/item/edit", h.HandlerSetlistEditInfoCard)
	handleAuth("/setlist/item/save", h.HandlerSetlistSaveInfoCard)
	handleAuth("/setlist/item/cancel", h.HandlerSetlistPopupInfoCard)

	handleAuth("/transition", h.HandlerTransitionPage)
	handleAuth("/transition/create", h.HandlerTransitionCreatePage)
	handleAuth("/transition/save", h.HandlerTransitionSave)
	handleAuth("/transition/edit", h.HandlerTransitionEdit)
	handleAuth("/transition/delete", h.HandlerDeleteTransition)

	handleAuth("/break", h.HandlerBreakPage)
	handleAuth("/break/create", h.HandlerBreakCreatePage)
	handleAuth("/break/save", h.HandlerBreakSave)
	handleAuth("/break/edit", h.HandlerBreakSave)
	handleAuth("/break/delete", h.HandlerDeleteBreak)

	handleAuth("/promotion", h.HandlerPromotion)
	handleAuth("/goals", h.HandlerGoals)
	handleAuth("/calendar", h.HandlerCalendar)
	handleAuth("/events", h.HandlerEvents)
	handleAuth("/files", h.HandlerFiles)

	handleAuth("/profile", h.HandlerProfilePage)
	handleAuth("/profile/add", h.HandlerProfilePicAdd)

	handleAuth("/admin", h.HandlerAdmin)
	handleAuth("/admin/access-code", h.HandlerCreateAccessCode)

	http.HandleFunc("/settings", h.HandlerSettingsPage)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	address := "0.0.0.0:" + port

	log.Printf("BandPlan listening on %s", address)

	if err := http.ListenAndServe(address, nil); err != nil {
		log.Fatal("Server failed: ", err)
	}
}
