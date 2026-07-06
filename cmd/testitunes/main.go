package main

import (
	"fmt"
	"log"

	"bandplan/src/models"
	"bandplan/src/services"
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

	fmt.Println("\nSearch Response: ")
	fmt.Println(searchResponse.ResultCount)
	fmt.Println("\nSearch Results")

	var song models.Song

	for _, r := range searchResponse.Results {
		if r.ArtistName == band.Name {
			fmt.Printf("%#v\n", r)
			fmt.Println("")

		}

	}
	fmt.Println(song)
}
