package handlers

import (
	"bandplan/src/database"
	"bandplan/src/models"
	"log"
	"net/http"
)

func (h Handler) HandlerProfilePicPage(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerProfilePicPage")

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

	data := models.MenuPageData{
		User: user,
		Band: band,
	}

	err = h.Tmpl.ExecuteTemplate(w, "profile-pic.html", data)
	if err != nil {
		log.Println("   Err getting profile pic page: ", err)
		return
	}

}

func (h Handler) HandlerProfilePicAdd(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerProfilePicadd")
}
