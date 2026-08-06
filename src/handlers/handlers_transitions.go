package handlers

import (
	"bandplan/src/models"
	"fmt"
	"log"
	"net/http"
)

func (h Handler) HandlerTransitionPage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("---------------------------")
	log.Println("- HandlerTransitionPage")
	return
}

func (h Handler) HandlerTransitionAddPage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("---------------------------")
	log.Println("- HandlerTransitionAdd")

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

	err = h.Tmpl.ExecuteTemplate(w, "transitions-add.html", data)
	if err != nil {
		log.Println("   Err getting transitions-add page: ", err)
		http.Redirect(w, r, "/setlists", http.StatusSeeOther)
	}

	return
}

func (h Handler) HandlerTransitionsCreate(w http.ResponseWriter, r *http.Request) {
	fmt.Println("---------------------------")
	log.Println("- HandlerTransitionCreate")

	return
}

func (h Handler) HandlerDeleteTransition(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerDeleteTransition")
	return
}
