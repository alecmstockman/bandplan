package services

import (
	"bandplan/src/helpers"
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"

	"github.com/chai2010/webp"
	"github.com/disintegration/imaging"
)

const (
	maxImageWidth  = 6000
	maxImageHeight = 6000
	maxImagePixels = 36_000_000
)

var (
	ErrInvalidImage            = errors.New("invalid image")
	ErrUnsupportedImageType    = errors.New("unsupported image type")
	ErrImageDimensionsExceeded = errors.New("image dimensions exceeded")
)

var allowedImageTypes = map[string]string{
	"image/jpeg": "jpeg",
	"image/png":  "png",
	"image/webp": "webp",
}

func decodeValidatedImage(file multipart.File) (image.Image, error) {
	log.Print("- decodeValidatedImage")

	header := make([]byte, 512)
	n, err := io.ReadFull(file, header)

	if errors.Is(err, io.EOF) {
		return nil, fmt.Errorf("decodeValidatedImage - read image header: %w", ErrInvalidImage)
	}
	if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) {
		return nil, fmt.Errorf("decodeValidatedImage - read image header: %w", err)
	}

	contentType := http.DetectContentType(header[:n])
	expectedFormat, ok := allowedImageTypes[contentType]
	if !ok {
		return nil, fmt.Errorf("decodeValidatedImage - detect content type: %w: %s", ErrUnsupportedImageType, contentType)
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return nil, fmt.Errorf("decodeValidatedImage - rewind before validation: %w", err)
	}

	config, format, err := image.DecodeConfig(file)
	if err != nil {
		return nil, fmt.Errorf("decodeValidatedImage - decode config: %w: %w", ErrInvalidImage, err)
	}

	pixels := int64(config.Width) * int64(config.Height)
	if config.Width <= 0 ||
		config.Height <= 0 ||
		config.Width > maxImageWidth ||
		config.Height > maxImageHeight ||
		pixels > maxImagePixels {
		return nil, fmt.Errorf(
			"decodeValidatedImage - validate dimensions: %w: %dx%d (%d pixels)",
			ErrImageDimensionsExceeded,
			config.Width,
			config.Height,
			pixels,
		)
	}

	if format != expectedFormat {
		return nil, fmt.Errorf("decodeValidatedImage - validate format: %w: detected %s but decoded as %s", ErrUnsupportedImageType, contentType, format)
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return nil, fmt.Errorf("decodeValidatedImage - rewind before decoding: %w", err)
	}

	img, err := helpers.NormalizeImageOrientation(file)
	if err != nil {
		return nil, fmt.Errorf("decodeValidatedImage - normalize orientation: %w: %w", ErrInvalidImage, err)
	}

	return img, nil
}

func (s Service) ServiceSaveTempImage(ctx context.Context, file multipart.File, imageID, bandSlug, format string) (string, error) {
	log.Println("- ServiceSaveTempImage")

	img, err := decodeValidatedImage(file)
	if err != nil {
		return "", fmt.Errorf("ServiceSaveTempImage - decode image: %w", err)
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

	case "event":
		sizes = map[string][2]int{
			"small":  {128, 128},
			"medium": {512, 512},
			"large":  {1024, 1024},
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
			return "", fmt.Errorf("ServiceSaveTempImage - encoding webp: %w", err)
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
			return "", fmt.Errorf("ServiceSaveTempImage - R2 upload: %w", err)
		}
		log.Printf("   Uploaded %s artwork: key=%q url=%q", name, objectKey, publicURL)

		if name == "medium" {
			previewURL = publicURL
		}
	}

	if previewURL == "" {
		return "", errors.New("ServiceSaveTempImage - temporary image preview URL was not created")
	}

	log.Println("   Preview URL:", previewURL)

	return previewURL, nil
}

func (s Service) ServiceDeleteTempImage(ctx context.Context, imageID, bandSlug string) error {
	log.Println("- ServiceDeleteTempImage")

	if imageID == "" {
		return errors.New("ServiceDeleteTempImage - image ID is empty")
	}

	key := fmt.Sprintf(
		"temp-images/%s/%s.webp",
		bandSlug,
		imageID,
	)

	err := s.Storage.Delete(ctx, key)
	if err != nil {
		return fmt.Errorf("ServiceDeleteTempImage - R2 delete temporary image: %w", err)
	}

	return nil
}

func (s Service) ServiceCreatePermChatImage(ctx context.Context, imageID, bandSlug, setlistSlug string) (string, error) {
	log.Println("- ServiceCreatePermChatImage")

	if imageID == "" {
		return "", errors.New("ServiceCreatePermChatImage - image ID is empty")
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
			return "", fmt.Errorf("ServiceCreatePermChatImage - R2 copy %s image: %w", name, err)
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
		return "", fmt.Errorf("ServiceCreatePermChatImage - build browser path: %w", err)
	}

	return browserPath, nil
}

func (s Service) ServiceCreatePermSetlistImage(ctx context.Context, imageID, bandSlug, setlistSlug string) (string, error) {
	log.Println("- ServiceCreatePermSetlistImage")

	if imageID == "" {
		return "", errors.New("ServiceCreatePermSetlistImage - image ID is empty")
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
			return "", fmt.Errorf("ServiceCreatePermSetlistImage - R2 copy %s image: %w", name, err)
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
		return "", fmt.Errorf("ServiceCreatePermSetlistImage - build browser path: %w", err)
	}

	return browserPath, nil
}

func (s Service) ServiceCreatePermEventImage(ctx context.Context, imageID, bandSlug string) (string, error) {
	if imageID == "" {
		return "", errors.New("ServiceCreatePermEventImage - image ID is empty")
	}

	sizes := []string{"small", "medium", "large"}
	for _, size := range sizes {
		sourceKey := fmt.Sprintf("temp-images/%s/%s/%s.webp", bandSlug, imageID, size)
		destinationKey := fmt.Sprintf("event-images/%s/%s/%s.webp", bandSlug, imageID, size)

		if err := s.Storage.Copy(ctx, sourceKey, destinationKey); err != nil {
			return "", fmt.Errorf("ServiceCreatePermEventImage - R2 copy %s image: %w", size, err)
		}
	}

	for _, size := range sizes {
		key := fmt.Sprintf("temp-images/%s/%s/%s.webp", bandSlug, imageID, size)
		if err := s.Storage.Delete(ctx, key); err != nil {
			log.Printf("   Unable to delete %s temporary event image: %v\n", size, err)
		}
	}

	browserPath, err := url.JoinPath(
		s.Storage.PublicURL,
		"event-images",
		bandSlug,
		imageID,
	)
	if err != nil {
		return "", fmt.Errorf("ServiceCreatePermEventImage - build browser path: %w", err)
	}

	return browserPath, nil
}

func (s Service) ServiceDeleteArtworkImageVersions(ctx context.Context, imageID, bandSlug string) error {
	log.Println("- ServiceDeleteArtworkImageVersions")

	if imageID == "" {
		return errors.New("ServiceDeleteArtworkImageVersions - image ID is empty")
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
			return fmt.Errorf("ServiceDeleteArtworkImageVersions - R2 delete %s image: %w", size, err)
		}
	}

	key := fmt.Sprintf(
		"song-images/%s/%s",
		bandSlug,
		imageID,
	)

	err := s.Storage.Delete(ctx, key)
	if err != nil {
		return fmt.Errorf("ServiceDeleteArtworkImageVersions - R2 delete image directory: %w", err)
	}

	return nil
}

func (s Service) ServiceDeleteProfileImageVersions(ctx context.Context, imageID, userSlug string) error {
	log.Println("- ServicesDeleteProfileImageVersions")

	if imageID == "" {
		return errors.New("ServiceDeleteProfileImageVersions - image ID is empty")
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
			return fmt.Errorf("ServiceDeleteProfileImageVersions - R2 delete %s image: %w", size, err)
		}
	}

	key := fmt.Sprintf(
		"profile-images/%s/%s",
		userSlug,
		imageID,
	)

	err := s.Storage.Delete(ctx, key)
	if err != nil {
		return fmt.Errorf("ServiceDeleteProfileImageVersions - R2 delete image directory: %w", err)
	}

	return nil
}

func (s Service) ServiceDeleteEventImageVersions(ctx context.Context, imageID, bandSlug string) error {
	log.Println("- ServicesDeleteProfileImageVersions")

	if imageID == "" {
		return errors.New("ServiceDeleteEventImageVersions - image ID is empty")
	}

	sizes := []string{
		"small",
		"medium",
		"large",
	}

	for _, size := range sizes {
		key := fmt.Sprintf(
			"event-images/%s/%s/%s.webp",
			bandSlug,
			imageID,
			size,
		)

		err := s.Storage.Delete(ctx, key)
		if err != nil {
			return fmt.Errorf("ServiceDeleteEventImageVersions - R2 delete %s image: %w", size, err)
		}
	}

	key := fmt.Sprintf(
		"event-images/%s/%s",
		bandSlug,
		imageID,
	)

	err := s.Storage.Delete(ctx, key)
	if err != nil {
		return fmt.Errorf("ServiceDeleteEventImageVersions - R2 delete image directory: %w", err)
	}

	return nil
}

func (s Service) ServiceSaveSongImageVersions(ctx context.Context, file multipart.File, imageID, bandSlug string) (string, error) {
	log.Println("- ServiceSaveSongImageVersions")

	img, err := decodeValidatedImage(file)
	if err != nil {
		return "", fmt.Errorf("ServiceSaveSongImageVersions - decode image: %w", err)
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
			return "", fmt.Errorf("ServiceSaveSongImageVersions - encoding webp: %w", err)
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
			return "", fmt.Errorf("ServiceSaveSongImageVersions - R2 save song artwork: %w", err)
		}
	}

	browserPath, err := url.JoinPath(
		s.Storage.PublicURL,
		"song-images",
		bandSlug,
		imageID,
	)
	if err != nil {
		return "", fmt.Errorf("ServiceSaveSongImageVersions - build browser path: %w", err)
	}

	return browserPath, nil
}

func (s Service) ServiceSaveProfileImageVersions(ctx context.Context, file multipart.File, imageID, userSlug string) (string, error) {
	log.Println("- ServiceSaveProfileImageVersions")

	img, err := decodeValidatedImage(file)
	if err != nil {
		return "", fmt.Errorf("ServiceSaveProfileImageVersions - decode image: %w", err)
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
			return "", fmt.Errorf("ServiceSaveProfileImageVersions - encoding webp: %w", err)
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
			return "", fmt.Errorf("ServiceSaveProfileImageVersions - R2 save profile image: %w", err)
		}
	}

	browserPath, err := url.JoinPath(
		s.Storage.PublicURL,
		"profile-images",
		userSlug,
		imageID,
	)
	if err != nil {
		return "", fmt.Errorf("ServiceSaveProfileImageVersions - build browser path: %w", err)
	}

	return browserPath, nil
}

func (s Service) ServiceDeleteSetlistImageVersions(ctx context.Context, imageID, bandSlug, setlistSlug string) error {
	log.Println("- ServiceDeleteSetlistImageVersions")

	if imageID == "" {
		return errors.New("ServiceDeleteSetlistImageVersions - image ID is empty")
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
			return fmt.Errorf("ServiceDeleteSetlistImageVersions - R2 delete %s image: %w", size, err)
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
		return fmt.Errorf("ServiceDeleteSetlistImageVersions - R2 delete image directory: %w", err)
	}

	return nil
}

func (s Service) ServiceSaveSetlistImageVersions(ctx context.Context, file multipart.File, imageID, bandSlug string, setlistSlug string) (string, error) {
	log.Println("- ServiceSaveSetlistImageVersions")

	img, err := decodeValidatedImage(file)
	if err != nil {
		return "", fmt.Errorf("ServiceSaveSetlistImageVersions - decode image: %w", err)
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
			return "", fmt.Errorf("ServiceSaveSetlistImageVersions - encoding webp: %w", err)
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
			return "", fmt.Errorf("ServiceSaveSetlistImageVersions - R2 save setlist image: %w", err)
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
		return "", fmt.Errorf("ServiceSaveSetlistImageVersions - build browser path: %w", err)
	}

	return browserPath, nil
}

func (s Service) ServiceSaveEventImageVersions(ctx context.Context, file multipart.File, imageID, bandSlug string) (string, error) {
	log.Println("- ServiceSaveEventImageVersions")

	img, err := decodeValidatedImage(file)
	if err != nil {
		return "", fmt.Errorf("ServiceSaveEventImageVersions - decode image: %w", err)
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
			return "", fmt.Errorf("ServiceSaveEventImageVersions - encoding webp: %w", err)
		}

		objectKey := fmt.Sprintf(
			"event-images/%s/%s/%s.webp",
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
			return "", fmt.Errorf("ServiceSaveEventImageVersions - R2 save event image: %w", err)
		}
	}

	browserPath, err := url.JoinPath(
		s.Storage.PublicURL,
		"event-images",
		bandSlug,
		imageID,
	)
	if err != nil {
		return "", fmt.Errorf("ServiceSaveEventImageVersions - build browser path: %w", err)
	}

	return browserPath, nil
}
