package handlers

import (
	"bandplan/src/database"
	"bandplan/src/models"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
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
	log.Println("- HandlerProfilePicAdd")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user, err := HelperGetAuthenticatedUser(r)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "File too large", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("profile-image")
	if err != nil {
		fmt.Println("   Unable to upload profile picture: ", err)
		return
	}
	defer file.Close()

	uploadDir := "./static/uploads/profile-images"

	err = os.MkdirAll(uploadDir, 0755)
	if err != nil {
		http.Error(w, "Could not create upload directory", http.StatusInternalServerError)
		return
	}

	fileType := filepath.Ext(header.Filename)
	imageID := uuid.New().String()
	fileName := imageID + fileType
	filePath := filepath.Join(uploadDir, fileName)

	dst, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Could not create destination file", http.StatusInternalServerError)
		return
	}

	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "Could not save uploaded file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	browserPath := "/static/uploads/profile-images/" + fileName

	err = database.UsersTableAddProfileImagePath(user.UserID, browserPath)
	if err != nil {
		http.Error(w, "Cound not save image path to db", http.StatusInternalServerError)
	}

	fmt.Println("Saved file to users table and:", filePath)
	http.Redirect(w, r, "/profile-pic", http.StatusSeeOther)
}
