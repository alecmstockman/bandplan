package main

import (
	"log"

	"bandplan/src/models"
	services "bandplan/src/services"
)

func main() {
	log.Println("Testing iTunes API")

	band := models.Band{
		Name: "Aiming Arrows",
	}

	searchResponse, err := services.ServicesSearchITunesByArtist(band)
	if err != nil {
		log.Fatal(err)
	}

	var song models.Song

	for _, r := range searchResponse.Results {
		if r.ArtistName == band.Name {
			log.Printf("%#v\n", r)
			log.Println("")
		}
	}
	log.Println(song)
}
