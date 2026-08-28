package services

import (
	"bandplan/src/clients"
	"bandplan/src/models"
	"bytes"
	"encoding/json"
	"image"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/chai2010/webp"
	"github.com/disintegration/imaging"
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

func (s Service) ServiceSaveArtworkImageFromITunes(artworkURL string, imageID string) (string, error) {
	log.Println("- ServiceSaveArtworkImageFromITunes")

	file, err := clients.ClientITunesGetArtwork(artworkURL)
	if err != nil {
		log.Println("   Error getting artwork from iTunes: ", err)
		return "", nil
	}

	reader := bytes.NewReader(file)

	img, _, err := image.Decode(reader)
	if err != nil {
		log.Println("   Unable to decode file: ", err)
		return "", err
	}

	img = imaging.Fill(img, 1024, 1024, imaging.Center, imaging.Lanczos)

	sizes := map[string]int{
		"small":  128,
		"medium": 512,
		"large":  1024,
	}

	uploadDir := "./static/uploads/song-images/" + imageID
	browserPath := "/static/uploads/song-images/" + imageID

	err = os.MkdirAll(uploadDir, 0755)
	if err != nil {
		return "", err
	}

	for name, size := range sizes {
		resized := imaging.Resize(img, size, size, imaging.Lanczos)

		outputPath := filepath.Join(uploadDir, name+".webp")

		out, err := os.Create(outputPath)
		if err != nil {
			return "", err
		}

		err = webp.Encode(out, resized, &webp.Options{Quality: 85})
		closeErr := out.Close()

		if err != nil {
			log.Println("   Error encoding webp: ", err)
			return "", err
		}
		if closeErr != nil {
			return "", closeErr
		}
	}

	return browserPath, nil
}
