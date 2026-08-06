package handlers

import (
	"bandplan/src/database"
	"bandplan/src/models"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func (h Handler) HandlerProfilePage(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerProfilePicPage")

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

	err = h.Tmpl.ExecuteTemplate(w, "profile.html", data)
	if err != nil {
		log.Println("   Err getting profile pic page: ", err)
		return
	}
}

func (h Handler) HandlerProfilePicAdd(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerProfilePicAdd")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		log.Println("   Unable to get AuthContext: ", err)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	user := auth.User

	oldImageID := user.ProfileImageID

	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Println("   Error parsing multipart form: ", err)
		http.Error(w, "File too large", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("profile-image")
	if err != nil {
		log.Println("   Unable to upload profile picture: ", err)
		return
	}
	defer file.Close()

	imageID := uuid.New().String()

	browserPath, err := h.HelperSaveProfileImageVersions(r.Context(), file, imageID, user.Slug)
	if err != nil {
		log.Println("   Unable to save image versions: ", err)
		http.Error(w, "Unable to save image versions: ", http.StatusInternalServerError)
	}

	fmt.Println("\nPROFILE IMAGE BROWSER PATH: ")

	err = database.UsersTableUpdateProfileImage(user.UserID, imageID, browserPath)
	if err != nil {
		log.Println("   Could not save image path to db: ", err)
		http.Error(w, "Cound not save image path to db", http.StatusInternalServerError)
	}

	fmt.Println("Browser path: ", browserPath)
	fmt.Println("new user profile image: ", imageID)
	fmt.Println("old user profile image: ", oldImageID)

	err = h.HelperDeleteProfileImageVersions(r.Context(), oldImageID, user.Slug)
	if err != nil {
		log.Println("   Could not delete old image path from cloud: ", err)
		// http.Error(w, "Could not delete old image path from cloud", http.StatusInternalServerError)
	}

	log.Println("   Saved file to users table and:", browserPath)
	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

func (h Handler) HandlerSettingsPage(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerSettingsPage")

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

	err = h.Tmpl.ExecuteTemplate(w, "settings.html", data)
	if err != nil {
		log.Println("   Err getting profile pic page: ", err)
		return
	}
}

func (h Handler) HandlerAdmin(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerAdmin")

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		log.Println("   Unable to get AuthContext: ", err)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	user := auth.User
	band := auth.CurrentBand

	users, err := database.BandMembersGetMembersByBandID(band.BandID)
	if err != nil {
		log.Println("   Unable to get band members from database: ", err)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	data := models.AdminPageData{
		User:  user,
		Band:  band,
		Users: users,
	}

	err = h.Tmpl.ExecuteTemplate(w, "admin.html", data)
	if err != nil {
		log.Println("   err getting admin page.html: ", err)
		return
	}
	return
}
