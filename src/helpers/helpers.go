package helpers

import (
	"image"
	"io"
	"log"
	"mime/multipart"
	"strings"
	"unicode"

	"github.com/disintegration/imaging"
	"github.com/google/uuid"
	"github.com/rwcarlsen/goexif/exif"
)

func MakeSlug(text string) string {
	log.Println("- Helper - MakeSlug")
	name := strings.ToLower(strings.TrimSpace(text))

	var b strings.Builder
	lastDash := false

	for _, r := range name {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			b.WriteRune(r)
		} else if !lastDash {
			b.WriteRune('-')
			lastDash = true
		}
	}
	id := "-" + uuid.NewString()[:6]

	return strings.Trim(b.String(), "-") + id
}

func NormalizeImageOrientation(file multipart.File) (image.Image, error) {
	log.Println("- Helper - NormalizeImageOrientation")

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
