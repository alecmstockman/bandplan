package main

import (
	"bandplan/src/database"
	"bandplan/src/handlers"
	"bandplan/src/middleware"
	"bandplan/src/realtime"
	"bandplan/src/services"
	"bandplan/src/storage"
	"context"
	"log"
	"net/http"
	"os"
)

var messages []string

func handleAuth(mux *http.ServeMux, pattern string, handler http.HandlerFunc) {
	mux.Handle(
		pattern,
		middleware.MiddlewareRecover(
			middleware.RequestID(
				middleware.RequireAuth(
					middleware.RequestLogging(
						handler,
					),
				),
			),
		),
	)
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
		Services: &services.Service{
			DB:      database.DB,
			Storage: r2Storage,
		},
	}

	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	mux.HandleFunc("GET /health", h.HandlerHealth)

	mux.HandleFunc("GET /{$}", h.HandlerHome)
	mux.HandleFunc("/access", h.HandlerAccessCodePage)
	mux.HandleFunc("/register", h.HandlerRegisterPage)

	mux.HandleFunc("GET /register/1", h.HandlerRegisterPageOne)
	mux.HandleFunc("POST /register/1", h.HandlerRegisterPageOne)

	mux.HandleFunc("GET /register/2", h.HandlerRegisterPageTwo)
	mux.HandleFunc("POST /register/2", h.HandlerRegisterPageTwo)
	mux.HandleFunc("POST /register/2-submit", h.HandlerRegisterPageTwoSubmit)

	mux.HandleFunc("POST /register/3", h.HandlerRegisterPageThreeSubmit)

	mux.HandleFunc("GET /register/4", h.HandlerRegisterPageFour)
	mux.HandleFunc("POST /register/4", h.HandlerRegisterPageFourSubmit)

	mux.HandleFunc("POST /register/5", h.HandlerRegisterPageFiveCreate)

	mux.HandleFunc("POST /register/create", h.HandlerRegister)
	mux.HandleFunc("POST /register/user-agreement", h.HandlerUserAgreementPage)
	mux.HandleFunc("POST /register/user-agreed", h.HandlerUserAgreement)
	mux.HandleFunc("GET /terms", h.HandlerTermsPage)
	mux.HandleFunc("GET /privacy", h.HandlerPrivacyPage)

	mux.HandleFunc("GET /login", h.HandlerLoginPage)
	mux.HandleFunc("POST /login/enter", h.HandlerLogin)
	handleAuth(mux, "POST /logout", h.HandlerLogout)

	handleAuth(mux, "POST /delete", h.HandlerDelete)
	handleAuth(mux, "GET /messages", h.HandlerMessages)

	handleAuth(mux, "GET /ws/chat", h.HandlerChatWebSocket)

	handleAuth(mux, "GET /chats", h.HandlerChatsPage)
	handleAuth(mux, "GET /chats/add", h.HandlerChatAddPage)
	handleAuth(mux, "POST /chats/members/select", h.HandlerChatSelectMember)
	handleAuth(mux, "POST /chats/members/remove", h.HandlerChatRemoveMember)
	handleAuth(mux, "POST /chats/create", h.HandlerChatCreate)

	handleAuth(mux, "GET /chat", h.HandlerChatPage)
	handleAuth(mux, "GET /chat/settings", h.HandlerChatSettings)
	handleAuth(mux, "GET /chat/settings/members", h.HandlerChatSettingsMembers)
	handleAuth(mux, "GET /chat/settings/pinned-chats", h.HandlerPinnedChats)
	handleAuth(mux, "POST /chat/leave", h.HandlerChatLeave)
	handleAuth(mux, "POST /chat/delete", h.HandlerChatDelete)

	handleAuth(mux, "GET /chat/image", h.HandlerChatAddImagePage)
	handleAuth(mux, "POST /chat/image/save", h.HandlerChatImageSave)
	handleAuth(mux, "POST /chat/temp-art/add", h.HandlerChatTempArt)
	handleAuth(mux, "DELETE /chat/temp-art/delete", h.HandlerChatTempArtDelete)

	handleAuth(mux, "POST /chat/message/reaction", h.HandlerChatMessageReaction)
	handleAuth(mux, "GET /chat/message/reply", h.HandlerChatMessageReply)
	handleAuth(mux, "GET /chat/message/edit", h.HandlerChatMessageEdit)
	handleAuth(mux, "POST /chat/message/pin", h.HandlerChatMessagePinAdd)
	handleAuth(mux, "POST /chat/message/pin/remove", h.HandlerChatMessagePinRemove)
	handleAuth(mux, "POST /chat/message/delete", h.HandlerChatMessageDelete)

	handleAuth(mux, "GET /songs", h.HandlerSongsPage)
	handleAuth(mux, "GET /songs/add", h.HandlerSongsAddPage)
	handleAuth(mux, "POST /songs/create", h.HandlerSongsAdd)
	handleAuth(mux, "GET /songs/search", h.HandlerSongsSearch)
	handleAuth(mux, "GET /songs/itunes/query", h.HandlerSongsITunesQueryPage)
	handleAuth(mux, "GET /songs/itunes/download", h.HandlerSongsITunesQuery)
	handleAuth(mux, "GET /songs/itunes/results", h.HandlerSongsITunesResults)
	handleAuth(mux, "POST /songs/itunes/create", h.HandlerSongsITunesResultsAddSong)
	handleAuth(mux, "GET /song", h.HandlerSongPage)
	handleAuth(mux, "GET /song/edit", h.HandlerSongEditPage)
	handleAuth(mux, "POST /song/update", h.HandlerSongUpdate)
	handleAuth(mux, "POST /song/delete", h.HandlerSongDelete)
	handleAuth(mux, "GET /song/lyrics", h.HandlerSongLyrics)

	handleAuth(mux, "GET /setlists", h.HandlerSetlistsPage)
	handleAuth(mux, "GET /setlists/add", h.HandlerSetlistsAddPage)
	handleAuth(mux, "POST /setlists/temp-art/add", h.HandlerSetlistsTempArt)
	handleAuth(mux, "DELETE /setlists/temp-art/delete", h.HandlerSetlistsTempArtDelete)
	handleAuth(mux, "POST /setlists/create", h.HandlerSetlistsCreate)
	handleAuth(mux, "POST /setlists/delete", h.HandlerSetlistsDelete)
	handleAuth(mux, "GET /setlist", h.HandlerSetlistPage)
	handleAuth(mux, "POST /setlist/add", h.HandlerSetlistAddSong)
	handleAuth(mux, "GET /setlist/edit", h.HandlerSetlistEditPage)
	handleAuth(mux, "POST /setlist/update", h.HandlerSetlistUpdate)
	handleAuth(mux, "POST /setlist/delete", h.HandlerSetlistDeleteSong)
	handleAuth(mux, "POST /setlist/transition/delete", h.HandlerSetlistDeleteTransition)

	handleAuth(mux, "GET /setlist/pdf/print", h.HandlerSetlistPDFPrint)
	handleAuth(mux, "POST /setlist/pdf/save", h.HandlerSetlistPDFSave)
	handleAuth(mux, "POST /setlist/share", h.HandlerSetlistPDFSave)
	handleAuth(mux, "GET /setlist/songs", h.HandlerSetlistSongs)
	handleAuth(mux, "GET /setlist/transitions", h.HandlerSetlistTransitions)
	handleAuth(mux, "GET /setlist/breaks", h.HandlerSetlistBreaks)
	handleAuth(mux, "GET /setlist/items", h.HandlerSetlistItems)
	handleAuth(mux, "POST /setlist/delete/song", h.HandlerSetlistDeleteSong)
	handleAuth(mux, "POST /setlist/delete/transition", h.HandlerSetlistDeleteTransition)
	handleAuth(mux, "POST /setlist/delete/break", h.HandlerSetlistItems)
	handleAuth(mux, "POST /setlist/dublicate", h.HandlerSetlistDuplicate)

	handleAuth(mux, "GET /setlist/reorder", h.HandlerSetlistReorder)
	handleAuth(mux, "POST /setlist/reorder/save", h.HandlerSetlistReorderSave)

	handleAuth(mux, "GET /setlist/time/songs", h.HandlerSetlistUpdateTimeButtonSongs)
	handleAuth(mux, "GET /setlist/time/transitions", h.HandlerSetlistUpdateTimeButtonTransitions)
	handleAuth(mux, "GET /setlist/time/breaks", h.HandlerSetlistUpdateTimeButtonBreaks)
	handleAuth(mux, "GET /setlist/time/items", h.HandlerSetlistUpdateTimeButtonItems)

	handleAuth(mux, "GET /setlist/notes", h.HandlerSetlistOpenNotesPage)
	handleAuth(mux, "GET /setlist/notes/edit", h.HandlerSetlistEditNotesPage)
	handleAuth(mux, "POST /setlist/notes/save", h.HandlerSetlistSaveNotesPage)

	handleAuth(mux, "GET /setlist/item/edit", h.HandlerSetlistEditInfoCard)
	handleAuth(mux, "POST /setlist/item/save", h.HandlerSetlistSaveInfoCard)
	handleAuth(mux, "GET /setlist/item/cancel", h.HandlerSetlistPopupInfoCard)

	handleAuth(mux, "GET /transition", h.HandlerTransitionPage)
	handleAuth(mux, "GET /transition/create", h.HandlerTransitionCreatePage)
	handleAuth(mux, "POST /transition/save", h.HandlerTransitionSave)
	handleAuth(mux, "GET /transition/edit", h.HandlerTransitionEditPage)
	handleAuth(mux, "POST /transition/update", h.HandlerTransitionUpdate)
	handleAuth(mux, "POST /transition/delete", h.HandlerDeleteTransition)

	handleAuth(mux, "GET /break", h.HandlerBreakPage)
	handleAuth(mux, "GET /break/create", h.HandlerBreakCreatePage)
	handleAuth(mux, "POST /break/save", h.HandlerBreakSave)
	handleAuth(mux, "GET /break/edit", h.HandlerBreakEditPage)
	handleAuth(mux, "POST /break/update", h.HandlerBreakUpdate)
	handleAuth(mux, "POST /break/delete", h.HandlerDeleteBreak)

	handleAuth(mux, "GET /events", h.HandlerEventsPage)
	handleAuth(mux, "GET /events/create", h.HandlerEventCreate)
	handleAuth(mux, "POST /events/save", h.HandlerEventSave)

	handleAuth(mux, "GET /event", h.HandlerEventPage)
	handleAuth(mux, "GET /event/edit", h.HandlerEventEdit)
	handleAuth(mux, "POST /event/update", h.HandlerEventUpdate)
	handleAuth(mux, "POST /event/delete", h.HandlerEventDelete)
	handleAuth(mux, "POST /event/temp-art/add", h.HandlerEventTempArt)
	handleAuth(mux, "DELETE /event/temp-art/delete", h.HandlerEventTempArtDelete)

	handleAuth(mux, "GET /todo", h.HandlerToDo)
	handleAuth(mux, "GET /todo/add", h.HandlerToDoAddPage)

	handleAuth(mux, "GET /promotion", h.HandlerPromotion)
	handleAuth(mux, "GET /goals", h.HandlerGoals)
	handleAuth(mux, "GET /calendar", h.HandlerCalendar)
	handleAuth(mux, "GET /files", h.HandlerFiles)

	handleAuth(mux, "GET /profile", h.HandlerProfilePage)
	handleAuth(mux, "POST /profile/add", h.HandlerProfilePicAdd)

	handleAuth(mux, "GET /admin", h.HandlerAdmin)
	handleAuth(mux, "POST /admin/access-code", h.HandlerCreateAccessCode)

	mux.HandleFunc("/settings", h.HandlerSettingsPage)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	address := "0.0.0.0:" + port

	log.Printf("BandPlan listening on %s", address)

	if err := http.ListenAndServe(address, mux); err != nil {
		log.Fatal("Server failed: ", err)
	}
}
