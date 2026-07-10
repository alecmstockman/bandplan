package handlers

import (
	"bandplan/src/database"
	"bandplan/src/models"
	"fmt"
	"log"
	"net/http"

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

	fmt.Println("\n DATA: ", data)

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

	file, _, err := r.FormFile("profile-image")
	if err != nil {
		log.Println("   Unable to upload profile picture: ", err)
		return
	}
	defer file.Close()

	// uploadDir := "./static/uploads/profile-images" + user.UserSlug

	// err = os.MkdirAll(uploadDir, 0755)
	// if err != nil {
	// 	http.Error(w, "Could not create upload directory", http.StatusInternalServerError)
	// 	return
	// }

	// fileType := filepath.Ext(header.Filename)
	imageID := uuid.New().String()
	// fileName := imageID + fileType
	// filePath := filepath.Join(uploadDir, fileName)

	// dst, err := os.Create(filePath)
	// if err != nil {
	// 	log.Println("   Unable to create destination file: ", err)
	// 	http.Error(w, "Could not create destination file", http.StatusInternalServerError)
	// 	return
	// }

	// _, err = io.Copy(dst, file)
	// if err != nil {
	// 	log.Println("   Unable to save uploaded file: ", err)
	// 	http.Error(w, "Could not save uploaded file", http.StatusInternalServerError)
	// 	return
	// }
	// defer dst.Close()

	browserPath, err := HelperSaveProfileImageVersions(file, imageID, user.UserSlug)
	if err != nil {
		log.Println("   Unable to save image versions: ", err)
		http.Error(w, "Unable to save image versions: ", http.StatusInternalServerError)
	}

	err = database.UsersTableUpdateProfileImage(user.UserID, imageID, browserPath)
	if err != nil {
		log.Println("   Could not save image path to db: ", err)
		http.Error(w, "Cound not save image path to db", http.StatusInternalServerError)
	}

	log.Println("   Saved file to users table and:", browserPath)
	http.Redirect(w, r, "/profile-pic", http.StatusSeeOther)
}
