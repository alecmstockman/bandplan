package handlers

import (
	"bandplan/src/clients"
	"bytes"
	"errors"
	"fmt"
	"image"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/chai2010/webp"
	"github.com/disintegration/imaging"
)

func HelperDeleteArtworkImageVersions(imageID string) error {
	log.Println("- HelperDeleteArtworkImageVersions")

	if imageID == "" {
		log.Println("   No imageID provided")
		return errors.New("imageID is empty")
	}

	uploadDir := "./static/uploads/song-images/"
	dirPath := filepath.Join(uploadDir, imageID)

	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		log.Println("   artwork directory does not exist:", dirPath)
		return nil
	}

	fmt.Println("\n dirpath: ", dirPath)
	err := os.RemoveAll(dirPath)
	if err != nil {
		log.Println("   Unable to delete artwork image versions: ", err)
		return err
	}
	return nil
}

func HelperSaveArtworkImageVersions(file multipart.File, imageID string) (string, error) {
	log.Println("- HelperSaveArtworkImageVersions")

	img, _, err := image.Decode(file)
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

func HelperSaveProfileImageVersions(file multipart.File, imageID string) (string, error) {
	log.Println("- HelperSaveProfileImageVersions")

	img, _, err := image.Decode(file)
	if err != nil {
		log.Println("   Unable to decode file: ", err)
		return "", err
	}

	img = imaging.Fill(img, 1024, 1024, imaging.Center, imaging.Lanczos)

	sizes := map[string]int{
		"small":  64,
		"medium": 256,
		"large":  512,
	}

	uploadDir := "./static/uploads/profile-images/" + imageID
	browserPath := "/static/uploads/profile-images/" + imageID

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

func HelperSaveArtworkImageFromITunes(artworkURL string, imageID string) (string, error) {
	log.Println("- HelperSaveArtworkImageVersions")

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
