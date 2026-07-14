package services

import (
	"bandplan/src/clients"
	"bandplan/src/models"
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

func ServicesSearchITunesByArtist(band models.Band) (models.ITunesSearchResponse, error) {
	log.Println("- ServicesSearchITunesByArtist")

	term := strings.ReplaceAll(band.Name, " ", "+")

	body, err := clients.ITunesSearchByArtist(term)
	if err != nil {
		fmt.Println("Unable to search iTunes by artist: ", err)
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
		log.Println("   Unable to search iTunes by artist: ", err)
		return models.ITunesSearchResponse{}, err
	}

	var searchResponse models.ITunesSearchResponse

	err = json.Unmarshal(body, &searchResponse)
	if err != nil {
		log.Println("   Error unmarshalling search response: ", err)
		return models.ITunesSearchResponse{}, err
	}

	return searchResponse, nil
}

func ServicesSearchITunesByITunesID(itunesID string) (models.ITunesSearchResponse, error) {
	log.Println("- ServicesSearchITunesByITunesID")

	body, err := clients.ITunesSearchByITunesID(itunesID)
	if err != nil {
		log.Println("   Unable to search iTunes by artist ID: ", err)
		return models.ITunesSearchResponse{}, err
	}

	var searchResponse models.ITunesSearchResponse

	err = json.Unmarshal(body, &searchResponse)
	if err != nil {
		log.Println("   Error unmarshalling search response: ", err)
		return models.ITunesSearchResponse{}, err
	}

	return searchResponse, nil
}
