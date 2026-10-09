package services

import (
	"bandplan/src/database"
	"bandplan/src/models"
	"context"
	"errors"
	"fmt"
)

func (s Service) SongsPage(ctx context.Context, user models.User, band models.Band) (models.MenuPageData, error) {

	songs, err := database.SongsTableGetAllSongsByBandID(band.BandID)
	if err != nil {
		return models.MenuPageData{}, errors.New("Unable to get songs from DB")
	}

	setlists, err := database.SetlistsTableGetSetlistsByBandIDAndUserID(band.BandID, user.UserID)
	if err != nil {
		return models.MenuPageData{}, fmt.Errorf("Unable to get setlists: ", err)
	}

	data := models.MenuPageData{
		User:     user,
		Band:     band,
		Songs:    songs,
		Setlists: setlists,
	}
	return data, nil
}
