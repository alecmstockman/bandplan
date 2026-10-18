package helpers

import (
	"bandplan/src/models"
	"errors"
	"image"
	"io"
	"mime/multipart"
	"net/mail"
	"strings"
	"time"
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
	id := "-" + uuid.NewString()[:12]

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

func PastEvents(events []models.Event) []models.Event {
	now := time.Now()

	today := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		0,
		0,
		0,
		0,
		now.Location(),
	)

	var pastEvents []models.Event

	for _, event := range events {
		eventDate := time.Date(
			event.EventDate.Year(),
			event.EventDate.Month(),
			event.EventDate.Day(),
			0,
			0,
			0,
			0,
			today.Location(),
		)

		if eventDate.Before(today) {
			pastEvents = append(pastEvents, event)
		}
	}
	return pastEvents
}

func UpcomingEvents(events []models.Event) []models.Event {
	now := time.Now()

	today := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		0,
		0,
		0,
		0,
		now.Location(),
	)

	var upcomingEvents []models.Event

	for _, event := range events {
		eventDate := time.Date(
			event.EventDate.Year(),
			event.EventDate.Month(),
			event.EventDate.Day(),
			0,
			0,
			0,
			0,
			today.Location(),
		)

		if !eventDate.Before(today) {
			upcomingEvents = append(upcomingEvents, event)
		}
	}
	return upcomingEvents
}

func FormatOptionalTime(value *time.Time, timezone string) string {
	if value == nil {
		return ""
	}

	location, err := time.LoadLocation(timezone)
	if err != nil {
		location = time.UTC
	}

	return value.In(location).Format("3:04 PM")
}

func ParseOptionalDueFields(dateValue, timeValue, timezoneValue string) (*time.Time, *time.Time, *string, error) {
	if dateValue == "" && timeValue == "" {
		return nil, nil, nil, nil
	}

	location, err := time.LoadLocation(timezoneValue)
	if err != nil {
		return nil, nil, nil, errors.New("Invalid due timezone")
	}

	var dueDate *time.Time
	if dateValue != "" {
		parsedDate, err := time.ParseInLocation("2006-01-02", dateValue, location)
		if err != nil {
			return nil, nil, nil, errors.New("Invalid due date")
		}
		dueDate = &parsedDate
	}

	var dueTime *time.Time
	if timeValue != "" {
		parsedTime, err := time.ParseInLocation("15:04", timeValue, location)
		if err != nil {
			return nil, nil, nil, errors.New("Invalid due time")
		}
		dueTime = &parsedTime
	}

	return dueDate, dueTime, &timezoneValue, nil
}

func ValidatePriceEntry(price string) bool {
	decimalCount := 0

	for _, r := range price {
		if r == '.' && len(price) == 1 {
			return false
		}
		if r == '.' {
			decimalCount++

			if decimalCount > 1 {
				return false
			}
			continue
		}

		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}
