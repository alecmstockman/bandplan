package services

import (
	"bandplan/src/helpers"
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"log"
	"mime/multipart"
	"net/url"

	"github.com/chai2010/webp"
	"github.com/disintegration/imaging"
)

func (s Service) ServiceSaveTempImage(ctx context.Context, file multipart.File, imageID string, bandSlug string, format string) (string, error) {
	log.Println("- ServiceSaveTempImage")

	img, err := helpers.NormalizeImageOrientation(file)
	if err != nil {
		log.Printf("   Error normalizing image file id: %v due to : %v\n", imageID, err)
		return "", err
	}

	var sizes map[string][2]int

	switch format {
	case "setlist":
		sizes = map[string][2]int{
			"small":  {256, 192},
			"medium": {512, 384},
		}

	case "profile":
		sizes = map[string][2]int{
			"small":  {64, 64},
			"medium": {256, 256},
			"large":  {512, 512},
		}

	case "chat":
		sizes = map[string][2]int{
			"small":  {64, 64},
			"medium": {256, 256},
			"large":  {512, 512},
		}

	default:
		sizes = map[string][2]int{
			"small":  {128, 128},
			"medium": {512, 512},
			"large":  {1024, 1024},
		}
	}

	previewURL := ""

	for name, size := range sizes {

		resized := imaging.Fill(
			img,
			size[0],
			size[1],
			imaging.Center,
			imaging.Lanczos,
		)

		var buffer bytes.Buffer

		err = webp.Encode(&buffer, resized, &webp.Options{Quality: 85})
		if err != nil {
			log.Println("   Error encoding webp: ", err)
			return "", err
		}

		objectKey := fmt.Sprintf(
			"temp-images/%s/%s/%s.webp",
			bandSlug,
			imageID,
			name,
		)

		publicURL, err := s.Storage.Upload(
			ctx,
			objectKey,
			&buffer,
			"image/webp",
		)

		if err != nil {
			log.Println("   Unable to upload song image to R2: ", err)
			return "", err
		}
		log.Printf("   Uploaded %s artwork: key=%q url=%q", name, objectKey, publicURL)

		if name == "medium" {
			previewURL = publicURL
		}
	}

	if previewURL == "" {
		return "", errors.New("temporary image preview URL was not created")
	}

	log.Println("   Preview URL:", previewURL)

	return previewURL, nil
}

func (s Service) ServiceDeleteTempImage(ctx context.Context, imageID string, bandSlug string) error {
	log.Println("- ServiceDeleteTempImage")

	if imageID == "" {
		log.Println("   No imageID provided")
		return errors.New("imageID empty")
	}

	key := fmt.Sprintf(
		"temp-images/%s/%s.webp",
		bandSlug,
		imageID,
	)

	err := s.Storage.Delete(ctx, key)
	if err != nil {
		log.Println("   Unable to delete temporary image from R2:", err)
		return errors.New("Unable to delete temporary image from R2")
	}

	return nil
}

func (s Service) ServiceCreatePermChatImage(ctx context.Context, imageID string, bandSlug string, setlistSlug string) (string, error) {
	log.Println("- ServiceCreatePermChatImage")

	if imageID == "" {
		log.Println("   No imageID provided")
		return "", errors.New("imageID empty")
	}

	sizes := map[string][2]int{
		"small":  {64, 64},
		"medium": {256, 256},
		"large":  {512, 512},
	}

	for name, _ := range sizes {
		sourceKey := fmt.Sprintf(
			"temp-images/%s/%s/%s.webp",
			bandSlug,
			imageID,
			name,
		)

		destinationKey := fmt.Sprintf(
			"chat-images/%s/%s/%s/%s.webp",
			bandSlug,
			setlistSlug,
			imageID,
			name,
		)

		err := s.Storage.Copy(
			ctx,
			sourceKey,
			destinationKey,
		)
		if err != nil {
			log.Printf("   Unable to copy %s chat image to R2 permanent storage: %v\n", name, err)
			return "", err
		}

		err = s.ServiceDeleteTempImage(ctx, imageID, bandSlug)
		if err != nil {
			log.Printf("   Unable to delete %s temporary chat image: %v\n", name, err)
		}

	}

	browserPath, err := url.JoinPath(
		s.Storage.PublicURL,
		"chat-images",
		bandSlug,
		setlistSlug,
		imageID,
	)
	if err != nil {
		log.Println("   Unable to create browser path for chat images: ", err)
		return "", err
	}

	return browserPath, nil
}

func (s Service) ServiceCreatePermSetlistImage(ctx context.Context, imageID string, bandSlug string, setlistSlug string) (string, error) {
	log.Println("- ServiceCreatePermSetlistImage")

	if imageID == "" {
		log.Println("   No imageID provided")
		return "", errors.New("imageID empty")
	}

	sizes := map[string][2]int{
		"small":  {256, 192},
		"medium": {512, 384},
	}

	for name, _ := range sizes {

		sourceKey := fmt.Sprintf(
			"temp-images/%s/%s/%s.webp",
			bandSlug,
			imageID,
			name,
		)

		destinationKey := fmt.Sprintf(
			"setlist-images/%s/%s/%s/%s.webp",
			bandSlug,
			setlistSlug,
			imageID,
			name,
		)

		err := s.Storage.Copy(
			ctx,
			sourceKey,
			destinationKey,
		)
		if err != nil {
			log.Printf("   Unable to copy %s song image to R2 permanent storage: %v\n", name, err)
			return "", err
		}

		err = s.ServiceDeleteTempImage(ctx, imageID, bandSlug)
		if err != nil {
			log.Printf("   Unable to delete %s temporary setlist art: %v\n", name, err)
		}

	}

	browserPath, err := url.JoinPath(
		s.Storage.PublicURL,
		"setlist-images",
		bandSlug,
		setlistSlug,
		imageID,
	)
	if err != nil {
		log.Println("   Unable to create browser path for setlist images: ", err)
		return "", err
	}

	return browserPath, nil
}

func (s Service) ServiceDeleteArtworkImageVersions(ctx context.Context, imageID string, bandSlug string) error {
	log.Println("- ServiceDeleteArtworkImageVersions")

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

		err := s.Storage.Delete(ctx, key)

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

	err := s.Storage.Delete(ctx, key)
	if err != nil {
		log.Println("   Unable to delete Song Artwork directory from R2: ", err)
		return errors.New("Unable to delete Song Artwork directory from R2")
	}

	return nil
}

func (s Service) ServiceDeleteProfileImageVersions(ctx context.Context, imageID string, userSlug string) error {
	log.Println("- ServicesDeleteProfileImageVersions")

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

		err := s.Storage.Delete(ctx, key)
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

	err := s.Storage.Delete(ctx, key)
	if err != nil {
		log.Println("   Unable to delete profile image directory from R2:", err)
		return errors.New("Unable to delete profile image directory from R2")
	}

	return nil
}

func (s Service) ServiceSaveArtworkImageVersions(ctx context.Context, file multipart.File, imageID string, bandSlug string) (string, error) {
	log.Println("- ServiceSaveArtworkImageVersions")

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

		_, err := s.Storage.Upload(
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
		s.Storage.PublicURL,
		"song-images",
		bandSlug,
		imageID,
	)

	return browserPath, nil
}

func (s Service) ServiceSaveProfileImageVersions(ctx context.Context, file multipart.File, imageID string, userSlug string) (string, error) {
	log.Println("- ServiceSaveProfileImageVersions")

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

		_, err := s.Storage.Upload(
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
		s.Storage.PublicURL,
		"profile-images",
		userSlug,
		imageID,
	)

	return browserPath, nil
}

func (s Service) ServiceDeleteSetlistImageVersions(ctx context.Context, imageID string, bandSlug string, setlistSlug string) error {
	log.Println("- ServiceDeleteSetlistImageVersions")

	if imageID == "" {
		log.Println("   No imageID provided")
		return errors.New("imageID empty")
	}

	sizes := []string{
		"small",
		"medium",
	}

	for _, size := range sizes {
		key := fmt.Sprintf(
			"setlist-images/%s/%s/%s/%s.webp",
			bandSlug,
			setlistSlug,
			imageID,
			size,
		)

		err := s.Storage.Delete(ctx, key)
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

	err := s.Storage.Delete(ctx, key)
	if err != nil {
		log.Println("   Unable to delete profile image directory from R2:", err)
		return errors.New("Unable to delete profile image directory from R2")
	}

	return nil
}

func (s Service) ServiceSaveSetlistImageVersions(ctx context.Context, file multipart.File, imageID string, bandSlug string, setlistSlug string) (string, error) {
	log.Println("- ServiceSaveSetlistImageVersions")

	img, err := helpers.NormalizeImageOrientation(file)
	if err != nil {
		log.Printf("   Error normalizing image file id: %v due to : %v\n", imageID, err)
		return "", err
	}

	img = imaging.Fill(img, 1024, 768, imaging.Center, imaging.Lanczos)

	sizes := map[string][2]int{
		"small":  {256, 192},
		"medium": {512, 384},
	}

	for name, size := range sizes {
		resized := imaging.Resize(img, size[0], size[1], imaging.Lanczos)

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

		_, err := s.Storage.Upload(
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
		s.Storage.PublicURL,
		"setlist-images",
		bandSlug,
		setlistSlug,
		imageID,
	)
	return browserPath, nil
}
