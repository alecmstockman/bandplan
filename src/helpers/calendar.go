package helpers

import (
	"bandplan/src/models"
	"net/url"
)

func GoogleCalendarURL(event models.Event) string {
	const layout = "20060102T150405Z"

	params := url.Values{}
	params.Set("action", "TEMPLATE")
	params.Set("text", event.Name)

	params.Set("details", event.Notes)
	// params.Set("location", location)

	if event.StartTime != nil && event.EndTime != nil {
		params.Set("dates",
			event.StartTime.UTC().Format(layout)+"/"+
				event.EndTime.UTC().Format(layout),
		)
	}

	return "https://calendar.google.com/calendar/r/eventedit?" + params.Encode()
}
