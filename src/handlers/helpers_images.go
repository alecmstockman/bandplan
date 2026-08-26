package handlers

import (
	"bandplan/src/clients"
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"io"
	"log"
	"mime/multipart"
	"net/url"
	"os"
	"path/filepath"

	"github.com/chai2010/webp"
	"github.com/disintegration/imaging"
	"github.com/rwcarlsen/goexif/exif"
)

func (h Handler) HelperSaveArtworkImageVersions(ctx context.Context, file multipart.File, imageID string, bandSlug string) (string, error) {
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
		"song-images",
		bandSlug,
		imageID,
	)

	return browserPath, nil
}

func (h Handler) HelperDeleteArtworkImageVersions(ctx context.Context, imageID string, bandSlug string) error {
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
		"profile-images",
		userSlug,
		imageID,
	)

	return browserPath, nil
}

func (h Handler) HelperSaveTempImage(ctx context.Context, file multipart.File, imageID string, bandSlug string, format string) (string, error) {
	log.Println("- HelperSaveTempImage")

	img, err := HelperNormalizeImageOrientation(file)
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

		fmt.Println("name: ", name, " size: ", size)

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

		fmt.Println("Object Key: ", objectKey)

		publicURL, err := h.Storage.Upload(
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

func HelperNormalizeImageOrientation(file multipart.File) (image.Image, error) {
	log.Println("- HelperNormalizeImageOrientation")

	orientation := 1

	x, err := exif.Decode(file)
	if err == nil {
		tag, err := x.Get(exif.Orientation)
		if err == nil {
			value, err := tag.Int(0)
			if err == nil {
				orientation = value
			}
		}
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		log.Println("   Error returning file seek to start: ", err)
		return nil, err
	}

	img, err := imaging.Decode(file)
	if err != nil {
		log.Println("   Unable to decode image file: ", err)
		return nil, err
	}

	switch orientation {
	case 2:
		img = imaging.FlipH(img)
	case 3:
		img = imaging.Rotate180(img)
	case 4:
		img = imaging.FlipV(img)
	case 5:
		img = imaging.Rotate270(imaging.FlipH(img))
	case 6:
		img = imaging.Rotate270(img)
	case 7:
		img = imaging.Rotate90(imaging.FlipH(img))
	case 8:
		img = imaging.Rotate90(img)
	}

	return img, nil
}

func (h Handler) HelperSaveSetlistImageVersions(ctx context.Context, file multipart.File, imageID string, bandSlug string, setlistSlug string) (string, error) {
	log.Println("- HelperSaveSetlistImageVersions")

	// img, _, err := image.Decode(file)
	// if err != nil {
	// 	log.Println("   Unable to decode file: ", err)
	// 	return "", err
	// }

	img, err := HelperNormalizeImageOrientation(file)
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
		"setlist-images",
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

func (h Handler) HelperCreatePermSetlistImage(ctx context.Context, imageID string, bandSlug string, setlistSlug string) (string, error) {
	log.Println("- HelperCreatePermSetlistImage")

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

		fmt.Println("sourcekey: ", sourceKey)

		destinationKey := fmt.Sprintf(
			"setlist-images/%s/%s/%s/%s.webp",
			bandSlug,
			setlistSlug,
			imageID,
			name,
		)

		err := h.Storage.Copy(
			ctx,
			sourceKey,
			destinationKey,
		)
		if err != nil {
			log.Printf("   Unable to copy %s song image to R2 permanent storage: %v\n", name, err)
			return "", err
		}

		err = h.HelperDeleteTempImage(ctx, imageID, bandSlug)
		if err != nil {
			log.Printf("   Unable to delete %s temporary setlist art: %v\n", name, err)
		}

	}

	browserPath, err := url.JoinPath(
		h.Storage.PublicURL,
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

func (h Handler) HelperDeleteTempImage(ctx context.Context, imageID string, bandSlug string) error {
	log.Println("- HelperDeleteTempImage")

	if imageID == "" {
		log.Println("   No imageID provided")
		return errors.New("imageID empty")
	}

	key := fmt.Sprintf(
		"temp-images/%s/%s.webp",
		bandSlug,
		imageID,
	)

	err := h.Storage.Delete(ctx, key)
	if err != nil {
		log.Println("   Unable to delete temporary image from R2:", err)
		return errors.New("Unable to delete temporary image from R2")
	}

	return nil
}
