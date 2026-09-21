package helpers

import (
	"image"
	"io"
	"mime/multipart"
	"net/mail"
	"strings"
	"unicode"

	"github.com/disintegration/imaging"
	"github.com/google/uuid"
	"github.com/rwcarlsen/goexif/exif"
)

func MakeSlug(text string) string {
	name := strings.ToLower(strings.TrimSpace(text))

	var b strings.Builder
	lastDash := false

	for _, r := range name {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			b.WriteRune(r)
			lastDash = false
		} else if !lastDash {
			b.WriteRune('-')
			lastDash = true
		}
	}
	id := "-" + uuid.NewString()[:6]

	return strings.Trim(b.String(), "-") + id
}

func NormalizeImageOrientation(file multipart.File) (image.Image, error) {

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
		return nil, err
	}

	img, err := imaging.Decode(file)
	if err != nil {
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

func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func ValidateEmail(email string) bool {

	if email == "" || len(email) > 254 || strings.Count(email, "@") != 1 {
		return false
	}

	address, err := mail.ParseAddress(email)
	if err != nil {
		return false
	}

	if address.Address != email {
		return false
	}

	parts := strings.SplitN(email, "@", 2)
	return strings.Contains(parts[1], ".")
}

func SmallImagePath(dir string) string {
	return dir + "/small.webp"
}

func MediumImagePath(dir string) string {
	return dir + "/medium.webp"
}

func LargeImagePath(dir string) string {
	return dir + "/large.webp"
}
