package handlers

import (
	"bandplan/src/clients"
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"log"
	"mime/multipart"
	"net/url"
	"os"
	"path/filepath"

	"github.com/chai2010/webp"
	"github.com/disintegration/imaging"
)

func (h Handler) HelperDeleteArtworkImageVersions(ctx context.Context, imageID string, bandSlug string) error {
	fmt.Println("")
	log.Println("- HelperDeleteArtworkImageVersions")

	if imageID == "" {
		log.Println("   No imageID provided")
		return errors.New("imageID is empty")
	}

	sizes := []string{
		"small",
		"medium",
		"large",
	}

	for _, size := range sizes {

		key := fmt.Sprintf(
			"song-images/%s/%s/%s.webp",
			bandSlug,
			imageID,
			size,
		)

		err := h.Storage.Delete(ctx, key)

		if err != nil {
			log.Println("   Unable to delete Song Artwork from R2: ", err)
			return errors.New("Unable to delete Song Artwork from R2")
		}

	}

	key := fmt.Sprintf(
		"song-images/%s/%s",
		bandSlug,
		imageID,
	)

	err := h.Storage.Delete(ctx, key)
	if err != nil {
		log.Println("   Unable to delete Song Artwork Folder from R2: ", err)
		return errors.New("Unable to delete Song Artwork Folder from R2")
	}

	return nil
}

func (h Handler) HelperSaveArtworkImageVersions(ctx context.Context, file multipart.File, imageID string, bandSlug string) (string, error) {
	fmt.Println("--------------------------------------------")
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

	// uploadDir := "./static/uploads/song-images/" + imageID
	// browserPath := "/static/uploads/song-images/" + imageID

	// err = os.MkdirAll(uploadDir, 0755)
	// if err != nil {
	// 	return "", err
	// }

	for name, size := range sizes {
		resized := imaging.Resize(img, size, size, imaging.Lanczos)

		// outputPath := filepath.Join(uploadDir, name+".webp")

		// out, err := os.Create(outputPath)
		// if err != nil {
		// 	return "", err
		// }

		fmt.Println("name: ", name, " size: ", size)

		var buffer bytes.Buffer

		err = webp.Encode(&buffer, resized, &webp.Options{Quality: 85})
		if err != nil {
			log.Println("   Error encoding webp: ", err)
			return "", err
		}

		objectKey := fmt.Sprintf(
			"song-images/%s/%s/%s.webp",
			bandSlug,
			imageID,
			name,
		)

		fmt.Println("objectKey: ", objectKey)

		_, err := h.Storage.Upload(
			ctx,
			objectKey,
			&buffer,
			"image/webp",
		)
		if err != nil {
			log.Println("   Unable to upload image to R2: ", err)
			return "", err
		}
	}

	browserPath, err := url.JoinPath(
		h.Storage.PublicURL,
		"song-images",
		bandSlug,
		imageID,
	)

	fmt.Println("browserPath: ", browserPath)

	return browserPath, nil
}

func HelperSaveProfileImageVersions(file multipart.File, imageID string, slug string) (string, error) {
	log.Println("- HelperSaveProfileImageVersions")

	img, _, err := image.Decode(file)
	if err != nil {
		log.Println("   Unable to decode file: ", err)
		return "", err
	}

	img = imaging.Fill(img, 512, 512, imaging.Center, imaging.Lanczos)

	sizes := map[string]int{
		"small":  64,
		"medium": 256,
		"large":  512,
	}

	uploadDir := "./static/uploads/profile-images/" + slug + "/" + imageID
	browserPath := "/static/uploads/profile-images/" + slug + "/" + imageID

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
