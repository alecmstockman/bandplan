package services

import (
	"bandplan/src/clients"
	"bandplan/src/models"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

func ServicesSearchITunesByArtist(band models.Band) (models.ITunesSearchResponse, error) {
	log.Println("- ServicesSearchITunesByArtist")

	term := strings.ReplaceAll(band.Name, " ", "+")

	body, err := clients.ITunesSearchByArtist(term)
	if err != nil {
		log.Println("   Unable to search iTunes by artist: ", err)
		return models.ITunesSearchResponse{}, err
	}

	var searchResponse models.ITunesSearchResponse

	err = json.Unmarshal(body, &searchResponse)
	if err != nil {
		return models.ITunesSearchResponse{}, err
	}

	return searchResponse, nil
}

func ServicesSearchITunesByArtistAndSong(artistName string, songName string) (models.ITunesSearchResponse, error) {
	log.Println("- ServicesSearchITunesByArtistAndSong")

	term := artistName + " " + songName

	body, err := clients.ITunesSearchByArtistAndSong(term)
	if err != nil {
		return models.ITunesSearchResponse{}, err
	}

	var searchResponse models.ITunesSearchResponse

	err = json.Unmarshal(body, &searchResponse)
	if err != nil {
		return models.ITunesSearchResponse{}, fmt.Errorf("unmarshal iTunes artist and song search response: %w", err)
	}

	return searchResponse, nil
}

func ServicesSearchITunesByITunesID(itunesID string) (models.ITunesSearchResponse, error) {
	log.Println("- ServicesSearchITunesByITunesID")

	body, err := clients.ITunesSearchByITunesID(itunesID)
	if err != nil {
		return models.ITunesSearchResponse{}, err
	}

	var searchResponse models.ITunesSearchResponse

	err = json.Unmarshal(body, &searchResponse)
	if err != nil {
		return models.ITunesSearchResponse{}, fmt.Errorf("unmarshal iTunes ID lookup response: %w", err)
	}

	return searchResponse, nil
}

func (s Service) ServiceDeleteITunesArtworkImageVersions(imageID string) error {
	if _, err := uuid.Parse(imageID); err != nil {
		return fmt.Errorf("invalid imported artwork ID: %w", err)
	}

	uploadDir := filepath.Join("static", "uploads", "song-images", imageID)
	if err := os.RemoveAll(uploadDir); err != nil {
		return fmt.Errorf("delete imported artwork: %w", err)
	}

	return nil
}
