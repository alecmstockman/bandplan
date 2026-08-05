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

	for name, size := range sizes {
		resized := imaging.Resize(img, size, size, imaging.Lanczos)

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
			log.Println("   Unable to upload song image to R2: ", err)
			return "", err
		}
	}

	browserPath, err := url.JoinPath(
		h.Storage.PublicURL,
		"/bandplan/song-images",
		bandSlug,
		imageID,
	)

	fmt.Println("browserPath: ", browserPath)

	return browserPath, nil
}

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
		log.Println("   Unable to delete Song Artwork directory from R2: ", err)
		return errors.New("Unable to delete Song Artwork directory from R2")
	}

	return nil
}

func (h Handler) HelperSaveProfileImageVersions(ctx context.Context, file multipart.File, imageID string, userSlug string) (string, error) {
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

	for name, size := range sizes {
		resized := imaging.Resize(img, size, size, imaging.Lanczos)

		var buffer bytes.Buffer

		err = webp.Encode(&buffer, resized, &webp.Options{Quality: 85})
		if err != nil {
			log.Println("   Error encoding webp: ", err)
			return "", err
		}

		objectKey := fmt.Sprintf(
			"profile-images/%s/%s/%s.webp",
			userSlug,
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
			log.Println("   Unable to upload profile image to R2: ", err)
			return "", err
		}
	}

	browserPath, err := url.JoinPath(
		h.Storage.PublicURL,
		"/bandplan/profile-images",
		userSlug,
		imageID,
	)

	fmt.Println(" func browserpath: ", browserPath)

	return browserPath, nil
}

func (h Handler) HelperSaveTempImage(ctx context.Context, file multipart.File, imageID string, bandSlug string, format string) (string, error) {
	log.Println("- HelperSaveTempImage")

	img, _, err := image.Decode(file)
	if err != nil {
		log.Println("   Unable to decode file: ", err)
		return "", err
	}

	if format == "setlist" {
		img = imaging.Fill(img, 512, 384, imaging.Center, imaging.Lanczos)
	} else {
		img = imaging.Fill(img, 512, 512, imaging.Center, imaging.Lanczos)
	}

	var buffer bytes.Buffer

	err = webp.Encode(&buffer, img, &webp.Options{Quality: 85})
	if err != nil {
		log.Println(".  Error encoding webp: ", err)
		return "", err
	}

	objectKey := fmt.Sprintf(
		"temp-images/%s/%s.webp",
		bandSlug,
		imageID,
	)

	_, err = h.Storage.Upload(
		ctx,
		objectKey,
		&buffer,
		"image/webp",
	)
	if err != nil {
		log.Println("   Unable to upload profile image to R2: ", err)
		return "", err
	}

	browserPath, err := url.JoinPath(
		h.Storage.PublicURL,
		"/bandplan/temp-images",
		bandSlug,
		imageID,
	)

	return browserPath + ".webp", nil
}

func (h Handler) HelperDeleteProfileImageVersions(ctx context.Context, imageID string, userSlug string) error {
	fmt.Println("")
	log.Println("- HelperDeleteProfileImageVersions")

	if imageID == "" {
		log.Println("   No imageID provided")
		return errors.New("imageID empty")
	}

	sizes := []string{
		"small",
		"medium",
		"large",
	}

	for _, size := range sizes {
		key := fmt.Sprintf(
			"profile-images/%s/%s/%s.webp",
			userSlug,
			imageID,
			size,
		)

		err := h.Storage.Delete(ctx, key)
		if err != nil {
			log.Println("   Unable to delete profile artwork")
			return errors.New("Unable to delete profile artwork")
		}
	}

	key := fmt.Sprintf(
		"profile-images/%s/%s",
		userSlug,
		imageID,
	)

	err := h.Storage.Delete(ctx, key)
	if err != nil {
		log.Println("   Unable to delete profile image directory from R2:", err)
		return errors.New("Unable to delete profile image directory from R2")
	}

	return nil
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

func (h Handler) HelperSaveSetlistImageVersions(ctx context.Context, file multipart.File, imageID string, bandSlug string, setlistSlug string) (string, error) {
	log.Println("- HelperSaveSetlistImageVersions")

	img, _, err := image.Decode(file)
	if err != nil {
		log.Println("   Unable to decode file: ", err)
		return "", err
	}

	img = imaging.Fill(img, 1024, 768, imaging.Center, imaging.Lanczos)

	sizes := map[string][2]int{
		"small":  {256, 192},
		"medium": {512, 384},
		"large":  {1024, 768},
	}

	for name, size := range sizes {
		resized := imaging.Resize(img, size[0], size[1], imaging.Lanczos)

		fmt.Println("name: ", name, " size: ", size)

		var buffer bytes.Buffer

		err = webp.Encode(&buffer, resized, &webp.Options{Quality: 85})
		if err != nil {
			log.Println("   Error encoding webp: ", err)
			return "", err
		}

		objectKey := fmt.Sprintf(
			"setlist-images/%s/%s/%s/%s.webp",
			bandSlug,
			setlistSlug,
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
			log.Println("   Unable to upload song image to R2: ", err)
			return "", err
		}
	}
	browserPath, err := url.JoinPath(
		h.Storage.PublicURL,
		"/bandplan/setlist-images",
		bandSlug,
		setlistSlug,
		imageID,
	)

	fmt.Println("browserPath: ", browserPath)

	return browserPath, nil
}

func (h Handler) HelperDeleteSetlistImageVersions(ctx context.Context, imageID string, bandSlug string, setlistSlug string) error {
	log.Println("- HelperDeleteSetlistImageVersions")

	if imageID == "" {
		log.Println("   No imageID provided")
		return errors.New("imageID empty")
	}

	sizes := []string{
		"small",
		"medium",
		"large",
	}

	for _, size := range sizes {
		key := fmt.Sprintf(
			"setlist-images/%s/%s/%s/%s.webp",
			bandSlug,
			setlistSlug,
			imageID,
			size,
		)

		err := h.Storage.Delete(ctx, key)
		if err != nil {
			log.Println("   Unable to delete setlist artwork")
			return errors.New("Unable to delete setlist artwork")
		}
	}

	key := fmt.Sprintf(
		"setlist-images/%s/%s/%s",
		bandSlug,
		setlistSlug,
		imageID,
	)

	err := h.Storage.Delete(ctx, key)
	if err != nil {
		log.Println("   Unable to delete profile image directory from R2:", err)
		return errors.New("Unable to delete profile image directory from R2")
	}

	return nil
}
