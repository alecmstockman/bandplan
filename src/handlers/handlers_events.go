package handlers

import (
	"bandplan/src/database"
	"bandplan/src/helpers"
	requestlog "bandplan/src/logging"
	"bandplan/src/models"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

func (h Handler) HandlerEventsPage(w http.ResponseWriter, r *http.Request) {

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		slog.Error(
			"unable to load auth context",
			"request_id", requestlog.GetRequestID(r.Context()),
			"error", err,
		)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	user := auth.User
	band := auth.CurrentBand

	events, err := database.EventsTableGetAllEventsByBandIDAndUserID(band.BandID, user.UserID)
	if err != nil {
		log.Println("   Unable to get events: ", err)
		slog.Error(
			"unable to get all events",
			"request_id", requestlog.GetRequestID(r.Context()),
			"path", r.URL.Path,
			"error", err,
		)
		http.Error(w, "Unable to get events: ", http.StatusInternalServerError)
		return
	}

	UpcomingEvents(events)

	data := models.EventsPageData{
		User:        user,
		Band:        band,
		Events:      events,
		UpcomingLen: len(UpcomingEvents(events)),
		PastLen:     len(PastEvents(events)),
	}

	err = h.Tmpl.ExecuteTemplate(w, "events.html", data)
	if err != nil {
		slog.Error(
			"unable to execute events.html template",
			"request_id", requestlog.GetRequestID(r.Context()),
			"path", r.URL.Path,
			"error", err,
		)
		return
	}
}

func (h Handler) HandlerEventCreate(w http.ResponseWriter, r *http.Request) {

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		slog.Error(
			"unable to load auth context",
			"request_id", requestlog.GetRequestID(r.Context()),
			"error", err,
		)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	user := auth.User
	band := auth.CurrentBand

	setlists, err := database.SetlistsTableGetSetlistNamesAndIDs(band.BandID, user.UserID)
	if err != nil {
		slog.Error(
			"unable to get setlist names and ids",
			"request_id", requestlog.GetRequestID(r.Context()),
			"path", r.URL.Path,
			"error", err,
		)
		http.Error(w, "Unable to load setlists", http.StatusInternalServerError)
		return
	}

	data := models.EventCreatePageData{
		User:     user,
		Band:     band,
		Setlists: setlists,
	}

	err = h.Tmpl.ExecuteTemplate(w, "event_create.html", data)
	if err != nil {
		slog.Error(
			"unable to execute event_create.html",
			"request_id", requestlog.GetRequestID(r.Context()),
			"path", r.URL.Path,
			"error", err,
		)
		http.Error(w, "Unable to load event creation page", http.StatusInternalServerError)
		return
	}
}

func (h Handler) HandlerEventSave(w http.ResponseWriter, r *http.Request) {

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		slog.Error(
			"unable to load auth context",
			"request_id", requestlog.GetRequestID(r.Context()),
			"path", r.URL.Path,
			"error", err,
		)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	band := auth.CurrentBand

	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		slog.Error(
			"unable to upload file, too large",
			"request_id", requestlog.GetRequestID(r.Context()),
			"path", r.URL.Path,
			"error", err,
		)
		http.Error(w, "File too large", http.StatusBadRequest)
		return
	}

	setlistID := strings.TrimSpace(r.FormValue("event-setlist"))
	if setlistID != "" {
		validSetlist, err := database.SetlistsTableSetlistBelongsToBandAndUser(
			setlistID,
			band.BandID,
			auth.User.UserID,
		)
		if err != nil {
			slog.Error(
				"unable to validate event setlist",
				"request_id", requestlog.GetRequestID(r.Context()),
				"error", err,
			)
			http.Error(w, "Unable to validate setlist", http.StatusInternalServerError)
			return
		}
		if !validSetlist {
			http.Error(w, "Invalid setlist", http.StatusBadRequest)
			return
		}
	}

	imagePath := ""
	imageID := ""
	temporaryImageID := r.FormValue("temporary-artwork-id")

	if temporaryImageID != "" {
		imageID = temporaryImageID
		imagePath, err = h.Services.ServiceCreatePermEventImage(
			r.Context(),
			imageID,
			band.Slug,
		)
		if err != nil {
			slog.Error(
				"unable to save image versions",
				"request_id", requestlog.GetRequestID(r.Context()),
				"path", r.URL.Path,
				"error", err,
			)
			http.Error(w, "could not save image versions", http.StatusInternalServerError)
			return
		}
	} else {
		file, _, fileErr := r.FormFile("image-path")
		if fileErr == nil {
			defer file.Close()

			imageID = uuid.NewString()
			imagePath, err = h.Services.ServiceSaveEventImageVersions(r.Context(), file, imageID, band.Slug)
			if err != nil {
				http.Error(w, "could not save image versions", http.StatusInternalServerError)
				return
			}
		} else if !errors.Is(fileErr, http.ErrMissingFile) {
			slog.Error(
				"unable to read event image",
				"request_id", requestlog.GetRequestID(r.Context()),
				"path", r.URL.Path,
				"error", err,
			)
			http.Error(w, "could not read event image", http.StatusBadRequest)
			return
		}
	}

	name := strings.TrimSpace(r.FormValue("event-name"))
	date := strings.TrimSpace(r.FormValue("event-date"))
	eventType := strings.TrimSpace(r.FormValue("event-type"))
	eventAges := strings.TrimSpace(r.FormValue("event-ages"))
	venueName := strings.TrimSpace(r.FormValue("venue-name"))
	eventAddress := strings.TrimSpace(r.FormValue("event-location"))

	startTime := strings.TrimSpace(r.FormValue("event-start-time"))
	endTime := strings.TrimSpace(r.FormValue("event-end-time"))
	timeZone := strings.TrimSpace(r.FormValue("event-timezone"))
	recurrence := strings.TrimSpace((r.FormValue("event-recurrence")))

	ticketLink := strings.TrimSpace(r.FormValue("event-ticket-link"))
	presaleTicketPrice := strings.TrimSpace(r.FormValue("event-presale-ticket-price"))
	ticketPrice := strings.TrimSpace(r.FormValue("event-ticket-price"))

	setLocation := strings.TrimSpace(r.FormValue("event-set-location"))
	loadInTime := strings.TrimSpace(r.FormValue("event-load-in-time"))
	soundCheckTime := strings.TrimSpace(r.FormValue("event-sound-check-time"))
	setTime := strings.TrimSpace(r.FormValue("event-set-time"))
	setLength := strings.TrimSpace(r.FormValue("event-set-length"))
	loadInInstructions := strings.TrimSpace(r.FormValue("load-in-instructions"))

	linkOneName := strings.TrimSpace(r.FormValue("link-one-name"))
	linkOne := strings.TrimSpace(r.FormValue("link-one"))
	linkTwoName := strings.TrimSpace(r.FormValue("link-two-name"))
	linkTwo := strings.TrimSpace(r.FormValue("link-two"))

	eventNotes := strings.TrimSpace(r.FormValue("event-notes"))

	if name == "" {
		http.Error(w, "Event name is required", http.StatusBadRequest)
		return
	}
	if date == "" {
		http.Error(w, "Event date is required", http.StatusBadRequest)
		return
	}

	validatedEventType := models.EventType(eventType)
	switch validatedEventType {
	case models.EventTypeNone,
		models.EventTypeShow,
		models.EventTypeGig,
		models.EventTypeFestival,
		models.EventTypeRehearsal,
		models.EventTypePractice,
		models.EventTypeWriting,
		models.EventTypeRecording,
		models.EventTypeMeeting,
		models.EventTypePhotos,
		models.EventTypePress,
		models.EventtypeOther:
	default:
		http.Error(w, "Invalid event type", http.StatusBadRequest)
		return
	}

	validatedEventAges := models.EventAges(eventAges)
	switch validatedEventAges {
	case models.EventNA,
		models.EventAllAges,
		models.Event18Plus,
		models.Event21Plus:
	default:
		http.Error(w, "Invalid event ages", http.StatusBadRequest)
		return
	}

	validatedRecurrence := models.EventRecurrence(recurrence)
	switch validatedRecurrence {
	case models.EventRecurrenceNone,
		models.EventRecurrenceDaily,
		models.EventRecurrenceWeekly,
		models.EventRecurrenceBiweekly,
		models.EventRecurrenceMonthly,
		models.EventRecurrenceYearly:
	default:
		http.Error(w, "Invalid event recurrence", http.StatusBadRequest)
		return
	}

	location, err := time.LoadLocation(timeZone)
	if err != nil {
		http.Error(w, "Invalid event timezone", http.StatusBadRequest)
		return
	}
	eventDate, err := time.ParseInLocation("2006-01-02", date, location)
	if err != nil {
		http.Error(w, "Invalid event date", http.StatusBadRequest)
		return
	}

	parseEventTime := func(value string) (time.Time, error) {
		if value == "" {
			return time.Time{}, nil
		}
		return time.ParseInLocation("2006-01-02 15:04", date+" "+value, location)
	}

	parsedStartTime, err := parseEventTime(startTime)
	if err != nil {
		http.Error(w, "Invalid event start time", http.StatusBadRequest)
		return
	}
	parsedEndTime, err := parseEventTime(endTime)
	if err != nil {
		http.Error(w, "Invalid event end time", http.StatusBadRequest)
		return
	}
	if !parsedStartTime.IsZero() && !parsedEndTime.IsZero() && parsedEndTime.Before(parsedStartTime) {
		parsedEndTime = parsedEndTime.AddDate(0, 0, 1)
	}

	parsedLoadInTime, err := parseEventTime(loadInTime)
	if err != nil {
		http.Error(w, "Invalid load-in time", http.StatusBadRequest)
		return
	}
	parsedSoundCheckTime, err := parseEventTime(soundCheckTime)
	if err != nil {
		http.Error(w, "Invalid sound check time", http.StatusBadRequest)
		return
	}
	parsedSetTime, err := parseEventTime(setTime)
	if err != nil {
		http.Error(w, "Invalid set time", http.StatusBadRequest)
		return
	}

	setLengthSeconds := 0
	if setLength != "" {
		setLengthMinutes, err := strconv.Atoi(setLength)
		if err != nil || setLengthMinutes < 0 {
			http.Error(w, "Invalid set length", http.StatusBadRequest)
			return
		}
		setLengthSeconds = setLengthMinutes * 60
	}

	eventID := uuid.NewString()

	fmt.Println("\nimageID: ", imageID)
	fmt.Println("imagePath: ", imagePath)

	newEvent := models.Event{
		EventID: eventID,
		BandID:  auth.CurrentBand.BandID,

		Name:      name,
		Slug:      helpers.MakeSlug(name),
		ImageID:   imageID,
		ImagePath: imagePath,

		EventDate:  eventDate,
		EventType:  validatedEventType,
		Ages:       validatedEventAges,
		Recurrence: validatedRecurrence,

		Address:   eventAddress,
		StartTime: &parsedStartTime,
		EndTime:   &parsedEndTime,
		Timezone:  timeZone,

		SetlistID:          setlistID,
		SetLocation:        setLocation,
		LoadInTime:         &parsedLoadInTime,
		LoadInInstructions: loadInInstructions,
		SoundCheckTime:     &parsedSoundCheckTime,
		SetTime:            &parsedSetTime,
		SetLengthSeconds:   setLengthSeconds,

		VenueName: venueName,

		PresaleTicketPrice: presaleTicketPrice,
		TicketPrice:        ticketPrice,
		TicketLink:         ticketLink,

		Notes: eventNotes,

		LinkOneName: linkOneName,
		LinkOne:     linkOne,
		LinkTwoName: linkTwoName,
		LinkTwo:     linkTwo,

		CreatedBy: auth.User.UserID,
		UpdatedBy: auth.User.UserID,
	}

	_, err = database.EventsTableCreateEvent(newEvent)
	if err != nil {
		log.Println("   Unable to save event to database: ", err)
		http.Error(w, "Unable to save event to database", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/events", http.StatusSeeOther)
}

func (h Handler) HandlerEventDelete(w http.ResponseWriter, r *http.Request) {

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		slog.Error(
			"unable to load auth context",
			"request_id", requestlog.GetRequestID(r.Context()),
			"error", err,
		)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	user := auth.User

	eventID := r.URL.Query().Get("id")

	err = database.EventsTableDeleteEvent(eventID, user.UserID)
	if err != nil {
		slog.Error(
			"unable to delete event",
			"request_id", requestlog.GetRequestID(r.Context()),
			"event-id", eventID,
			"path", r.URL.Path,
			"method", r.Method,
			"error", err,
		)
	}

	http.Redirect(w, r, "/events", http.StatusSeeOther)
}

func (h Handler) HandlerEventPage(w http.ResponseWriter, r *http.Request) {

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		slog.Error(
			"unable to load auth context",
			"request_id", requestlog.GetRequestID(r.Context()),
			"error", err,
		)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	user := auth.User
	band := auth.CurrentBand

	eventID := r.URL.Query().Get("event-id")

	event, err := database.EventsTableGetEventByEventIDAndBandID(eventID, band.BandID)
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "Event not found", http.StatusNotFound)
		return
	}
	if err != nil {
		slog.Error(
			"unable to get event",
			"request_id", requestlog.GetRequestID(r.Context()),
			"event_id", eventID,
			"band_id", band.BandID,
			"error", err,
		)
		http.Error(w, "2 Unable to get event", http.StatusInternalServerError)
		return
	}

	data := models.EventPageData{
		User:  user,
		Band:  band,
		Event: event,
		Time:  time.Now(),
	}

	err = h.Tmpl.ExecuteTemplate(w, "event.html", data)
	if err != nil {
		log.Println("   Unable to open events page: ", err)
		return
	}
}

func (h Handler) HandlerEventEdit(w http.ResponseWriter, r *http.Request) {

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		slog.Error(
			"unable to load auth context",
			"request_id", requestlog.GetRequestID(r.Context()),
			"error", err,
		)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	user := auth.User
	band := auth.CurrentBand

	eventID := r.URL.Query().Get("event-id")

	event, err := database.EventsTableGetEventByEventIDAndBandID(eventID, band.BandID)
	if err != nil {
		log.Println("   Unable to get event: ", err)
		http.Error(w, "3 Unable to get event", http.StatusInternalServerError)
		return
	}

	setlists, err := database.SetlistsTableGetSetlistNamesAndIDs(band.BandID, user.UserID)
	if err != nil {
		slog.Error(
			"unable to get setlist names and ids",
			"request_id", requestlog.GetRequestID(r.Context()),
			"error", err,
		)
		http.Error(w, "Unable to load setlists", http.StatusInternalServerError)
		return
	}

	data := models.EventPageData{
		User:     user,
		Band:     band,
		Event:    event,
		Setlists: setlists,
	}

	err = h.Tmpl.ExecuteTemplate(w, "event-edit.html", data)
	if err != nil {
		log.Println("   Unable to open events page: ", err)
		return
	}
}

func (h Handler) HandlerEventUpdate(w http.ResponseWriter, r *http.Request) {

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		slog.Error(
			"unable to load auth context",
			"request_id", requestlog.GetRequestID(r.Context()),
			"error", err,
		)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	eventID := r.FormValue("event-id")

	if eventID == "" {
		slog.Error(
			"No event ID on event update form",
			"request_id", requestlog.GetRequestID(r.Context()),
		)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	band := auth.CurrentBand

	event, err := database.EventsTableGetEventByEventIDAndBandID(eventID, band.BandID)
	if err != nil {
		slog.Error(
			"unable to get event from database",
			"request_id", requestlog.GetRequestID(r.Context()),
			"error", err,
		)
		http.Error(w, "Unable to update event", http.StatusInternalServerError)
		return
	}

	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "File too large", http.StatusBadRequest)
		return
	}

	setlistID := strings.TrimSpace(r.FormValue("event-setlist"))
	if setlistID != "" {
		validSetlist, err := database.SetlistsTableSetlistBelongsToBandAndUser(
			setlistID,
			band.BandID,
			auth.User.UserID,
		)
		if err != nil {
			slog.Error(
				"unable to validate event setlist",
				"request_id", requestlog.GetRequestID(r.Context()),
				"error", err,
			)
			http.Error(w, "Unable to validate setlist", http.StatusInternalServerError)
			return
		}
		if !validSetlist {
			http.Error(w, "Invalid setlist", http.StatusBadRequest)
			return
		}
	}
	event.SetlistID = setlistID

	oldImageID := event.ImageID
	imageChanged := false
	temporaryImageID := r.FormValue("temporary-artwork-id")
	removeImage := r.FormValue("remove-artwork") == "true"

	name := strings.TrimSpace(r.FormValue("event-name"))
	if name == "" {
		http.Error(w, "Event name is required", http.StatusBadRequest)
		return
	}

	date := strings.TrimSpace(r.FormValue("event-date"))

	if date == "" {
		http.Error(w, "Event date is required", http.StatusBadRequest)
		return
	}

	timeZone := strings.TrimSpace(r.FormValue("event-timezone"))

	location, err := time.LoadLocation(timeZone)
	if err != nil {
		http.Error(w, "Invalid event timezone", http.StatusBadRequest)
		return
	}
	event.Timezone = timeZone

	eventDate, err := time.ParseInLocation("2006-01-02", date, location)
	if err != nil {
		http.Error(w, "Invalid event date", http.StatusBadRequest)
		return
	}
	event.EventDate = eventDate

	if name != event.Name {
		event.Name = name
		event.Slug = helpers.MakeSlug(name)
	}

	event.Location = strings.TrimSpace(r.FormValue("location"))
	event.Address = strings.TrimSpace(r.FormValue("event-location"))
	event.VenueName = strings.TrimSpace(r.FormValue("venue-name"))
	event.AddressOne = strings.TrimSpace(r.FormValue("address-one"))
	event.AddressTwo = strings.TrimSpace(r.FormValue("address-two"))
	event.City = strings.TrimSpace(r.FormValue("city"))
	event.State = strings.TrimSpace(r.FormValue("state"))
	event.ZipCode = strings.TrimSpace(r.FormValue("zip-code"))

	event.SetLocation = strings.TrimSpace(r.FormValue("event-set-location"))
	event.LoadInInstructions = strings.TrimSpace(r.FormValue("load-in-instructions"))

	presaleTicketPrice := strings.TrimSpace(r.FormValue("event-presale-ticket-price"))
	ticketPrice := strings.TrimSpace(r.FormValue("event-ticket-price"))
	event.TicketLink = strings.TrimSpace(r.FormValue("event-ticket-link"))

	if ValidatePriceEntry(presaleTicketPrice) {
		event.PresaleTicketPrice = presaleTicketPrice
	} else {
		slog.Error(
			"unable to validate presales ticket price entry",
			"request_id", requestlog.GetRequestID(r.Context()),
		)
		http.Error(w, "invalid presale ticket price entry", http.StatusBadRequest)
		return
	}

	if ValidatePriceEntry(ticketPrice) {
		event.TicketPrice = ticketPrice
	} else {
		slog.Error(
			"unable to validate ticket price entry",
			"request_id", requestlog.GetRequestID(r.Context()),
		)
		http.Error(w, "invalid ticket price entry", http.StatusBadRequest)
		return
	}

	event.Notes = strings.TrimSpace(r.FormValue("event-notes"))
	event.LinkOneName = strings.TrimSpace(r.FormValue("link-one-name"))
	event.LinkOne = strings.TrimSpace(r.FormValue("link-one"))
	event.LinkTwoName = strings.TrimSpace(r.FormValue("link-two-name"))
	event.LinkTwo = strings.TrimSpace(r.FormValue("link-two"))
	event.UpdatedBy = auth.User.UserID

	eventType := strings.TrimSpace(r.FormValue("event-type"))

	validatedEventType := models.EventType(eventType)
	switch validatedEventType {
	case models.EventTypeNone,
		models.EventTypeShow,
		models.EventTypeGig,
		models.EventTypeFestival,
		models.EventTypeRehearsal,
		models.EventTypePractice,
		models.EventTypeWriting,
		models.EventTypeRecording,
		models.EventTypeMeeting,
		models.EventTypePhotos,
		models.EventTypePress,
		models.EventtypeOther:
	default:
		http.Error(w, "Invalid event type", http.StatusBadRequest)
		return
	}
	event.EventType = validatedEventType

	eventAges := models.EventAges(strings.TrimSpace(r.FormValue("event-ages")))
	switch eventAges {
	case models.EventNA,
		models.EventAllAges,
		models.Event18Plus,
		models.Event21Plus:
	default:
		http.Error(w, "Invalid event ages", http.StatusBadRequest)
		return
	}
	event.Ages = eventAges

	parseEventTime := func(value string) (*time.Time, error) {
		if value == "" {
			return nil, nil
		}
		parsed, err := time.ParseInLocation("2006-01-02 15:04", date+" "+value, location)
		if err != nil {
			return nil, err
		}

		return &parsed, nil
	}

	startTime := strings.TrimSpace(r.FormValue("event-start-time"))
	parsedStartTime, err := parseEventTime(startTime)
	if err != nil {
		http.Error(w, "Invalid event start time", http.StatusBadRequest)
		return
	}

	endTime := strings.TrimSpace(r.FormValue("event-end-time"))
	parsedEndTime, err := parseEventTime(endTime)
	if err != nil {
		http.Error(w, "Invalid event end time", http.StatusBadRequest)
		return
	}

	if parsedStartTime != nil &&
		parsedEndTime != nil &&
		parsedEndTime.Before(*parsedStartTime) {
		adjustedEndTime := parsedEndTime.AddDate(0, 0, 1)
		parsedEndTime = &adjustedEndTime
	}

	event.StartTime = parsedStartTime
	event.EndTime = parsedEndTime

	recurrence := strings.TrimSpace((r.FormValue("event-recurrence")))
	validatedRecurrence := models.EventRecurrence(recurrence)
	switch validatedRecurrence {
	case models.EventRecurrenceNone,
		models.EventRecurrenceDaily,
		models.EventRecurrenceWeekly,
		models.EventRecurrenceBiweekly,
		models.EventRecurrenceMonthly,
		models.EventRecurrenceYearly:
	default:
		http.Error(w, "Invalid event recurrence", http.StatusBadRequest)
		return
	}
	event.Recurrence = validatedRecurrence

	loadInTime := strings.TrimSpace(r.FormValue("event-load-in-time"))
	parsedLoadInTime, err := parseEventTime(loadInTime)
	if err != nil {
		http.Error(w, "Invalid load-in time", http.StatusBadRequest)
		return
	}
	event.LoadInTime = parsedLoadInTime

	soundCheckTime := strings.TrimSpace(r.FormValue("event-sound-check-time"))
	parsedSoundCheckTime, err := parseEventTime(soundCheckTime)
	if err != nil {
		http.Error(w, "Invalid sound check time", http.StatusBadRequest)
		return
	}
	event.SoundCheckTime = parsedSoundCheckTime

	setTime := strings.TrimSpace(r.FormValue("event-set-time"))
	parsedSetTime, err := parseEventTime(setTime)
	if err != nil {
		http.Error(w, "Invalid set time", http.StatusBadRequest)
		return
	}
	event.SetTime = parsedSetTime

	setLength := strings.TrimSpace(r.FormValue("event-set-length"))
	setLengthSeconds := 0

	if setLength != "" {
		setLengthMinutes, err := strconv.Atoi(setLength)
		if err != nil || setLengthMinutes < 0 {
			http.Error(w, "Invalid set length", http.StatusBadRequest)
			return
		}
		setLengthSeconds = setLengthMinutes * 60
	}
	event.SetLengthSeconds = setLengthSeconds

	newImageID := ""

	if temporaryImageID != "" {
		imagePath, err := h.Services.ServiceCreatePermEventImage(r.Context(), temporaryImageID, band.Slug)
		if err != nil {
			slog.Error(
				"unable to save permanent event images",
				"request_id", requestlog.GetRequestID(r.Context()),
				"error", err,
			)
			http.Error(w, "Could not save event image versions", http.StatusInternalServerError)
			return
		}

		newImageID = temporaryImageID
		event.ImageID = temporaryImageID
		event.ImagePath = imagePath
		imageChanged = true

	} else if removeImage {

		event.ImageID = ""
		event.ImagePath = ""
		imageChanged = true
	}

	_, err = database.EventsTableUpdateEvent(event)
	if err != nil {
		if newImageID != "" {
			cleanupErr := h.Services.ServiceDeleteEventImageVersions(r.Context(), newImageID, band.Slug)
			if cleanupErr != nil {
				slog.Error(
					"unable to clean up new event image after database failure",
					"request_id", requestlog.GetRequestID(r.Context()),
					"image_id", newImageID,
					"error", cleanupErr,
				)
			}
		}

		slog.Error(
			"unable to save event to database",
			"request_id", requestlog.GetRequestID(r.Context()),
			"image_id", newImageID,
			"error", err,
		)
		http.Error(w, "Unable to save event to database", http.StatusInternalServerError)
		return
	}

	if imageChanged && oldImageID != "" {
		err = h.Services.ServiceDeleteEventImageVersions(r.Context(), oldImageID, band.Slug)
		if err != nil {
			slog.Error(
				"unable to delete old event image",
				"request_id", requestlog.GetRequestID(r.Context()),
				"image_id", oldImageID,
				"error", err,
			)
		}
	}

	eventURL := fmt.Sprintf("/event?event-id=%v", eventID)
	http.Redirect(w, r, eventURL, http.StatusSeeOther)
}

func (h Handler) HandlerEventTempArt(w http.ResponseWriter, r *http.Request) {

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		slog.Error(
			"unable to load auth context",
			"request_id", requestlog.GetRequestID(r.Context()),
			"error", err,
		)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	band := auth.CurrentBand

	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Println("   Error parsing from while creating a new event: ", err)
		http.Error(w, "Unable to parse from", http.StatusBadRequest)
		return
	}

	imageID := ""
	previewURL := ""

	file, _, err := r.FormFile("image-path")
	if err != nil {
		if err != http.ErrMissingFile {
			log.Println("   Unable to read image file:", err)
			http.Error(w, "Unable to read image", http.StatusBadRequest)
			return
		}

		log.Println("   No event image uploaded")
	} else {
		defer file.Close()

		imageID = uuid.New().String()

		previewURL, err = h.Services.ServiceSaveTempImage(r.Context(), file, imageID, band.Slug, "event")
		if err != nil {
			http.Error(w, "Could not save image versions", http.StatusInternalServerError)
			return
		}
	}

	data := models.ArtworkPreviewData{
		ArtworkID:  imageID,
		PreviewURL: previewURL,
	}

	err = h.Tmpl.ExecuteTemplate(w, "event_image_preview", data)
	if err != nil {
		http.Error(w, "Unable to render preview", http.StatusInternalServerError)
		return
	}
}

func (h Handler) HandlerEventTempArtDelete(w http.ResponseWriter, r *http.Request) {

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		slog.Error(
			"unable to load auth context",
			"request_id", requestlog.GetRequestID(r.Context()),
			"error", err,
		)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	user := auth.User
	band := auth.CurrentBand

	temporaryImageID := r.FormValue("artwork-id")
	if temporaryImageID != "" {
		err = h.Services.ServiceDeleteTempImage(r.Context(), temporaryImageID, band.Slug)
		if err != nil {
			slog.Error(
				"unable to remove temporary event images",
				"request_id", requestlog.GetRequestID(r.Context()),
				"image_id", temporaryImageID,
				"error", err,
			)
			http.Error(w, "Unable to remove image", http.StatusInternalServerError)
			return
		}
	}

	data := models.SongDownloadData{
		User: user,
		Band: band,
	}

	err = h.Tmpl.ExecuteTemplate(w, "event_image_reset", data)
	if err != nil {
		log.Println("   Unable to exececute : ", err)
		http.Error(w, "Unable to render preview", http.StatusInternalServerError)
		return
	}
}
