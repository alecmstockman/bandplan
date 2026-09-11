package handlers

import (
	"bandplan/src/helpers"
	"bandplan/src/models"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

func (h Handler) HandlerEvents(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerEvents")

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		log.Println("   Unable to get AuthContext: ", err)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	user := auth.User
	band := auth.CurrentBand

	data := models.MenuPageData{
		User: user,
		Band: band,
	}

	err = h.Tmpl.ExecuteTemplate(w, "events.html", data)
	if err != nil {
		log.Println("   Unable to open events page: ", err)
		return
	}
}

func (h Handler) HandlerEventCreate(w http.ResponseWriter, r *http.Request) {
	log.Println("- HandlerEventCreate")

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		log.Println("   Unable to get AuthContext: ", err)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	user := auth.User
	band := auth.CurrentBand

	data := models.MenuPageData{
		User: user,
		Band: band,
	}

	err = h.Tmpl.ExecuteTemplate(w, "event_create.html", data)
}

func (h Handler) HandlerEventSave(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\n------------------------------------")
	log.Println("- HandlerEventSave")

	auth, err := HelperGetAuthContext(r)
	if err != nil {
		log.Println("   Unable to get AuthContext: ", err)
		http.Error(w, "Unable to load authenticated user", http.StatusInternalServerError)
		return
	}

	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Println("   File too large: ", err)
		http.Error(w, "File too large", http.StatusBadRequest)
		return
	}

	imagePath := ""
	imageID := ""

	file, _, err := r.FormFile("image-path")
	if err != nil {
		log.Println("   Error with provided image path: ", err)
		imageID = r.FormValue("existing-image-path")
		imagePath = r.FormValue("existing-image-path")
	} else {
		defer file.Close()

		imageID = uuid.NewString()

	}

	name := strings.TrimSpace(r.FormValue("event-name"))
	date := strings.TrimSpace(r.FormValue("event-date"))
	eventType := strings.TrimSpace(r.FormValue("event-type"))
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

	newEvent := models.Event{
		EventID: eventID,
		BandID:  auth.CurrentBand.BandID,

		Name:      name,
		Slug:      helpers.MakeSlug(name),
		ImageID:   imageID,
		ImagePath: imagePath,

		EventDate:  eventDate,
		EventType:  validatedEventType,
		Recurrence: validatedRecurrence,

		Address:   eventAddress,
		StartTime: parsedStartTime,
		EndTime:   parsedEndTime,
		Timezone:  timeZone,

		SetLocation:        setLocation,
		LoadInTime:         parsedLoadInTime,
		LoadInInstructions: loadInInstructions,
		SetTime:            parsedSetTime,
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

	log.Printf("   New event: %+v\n", newEvent)

	fmt.Println("event name:    ", name)
	fmt.Println("event date:    ", date)
	fmt.Println("event type:    ", eventType)
	fmt.Println("venue name:    ", venueName)
	fmt.Println("event address: ", eventAddress)

	fmt.Println("event start:   ", startTime)
	fmt.Println("event end:     ", endTime)
	fmt.Println("time zone:     ", timeZone)
	fmt.Println("event repeats: ", recurrence)

	fmt.Println("ticket link:   ", ticketLink)
	fmt.Println("presale price: ", presaleTicketPrice)
	fmt.Println("ticket price:  ", ticketPrice)

	fmt.Println("set location:  ", setLocation)
	fmt.Println("set time:      ", setTime)
	fmt.Println("set length:    ", setLength)
	fmt.Println("load in inst.: ", loadInInstructions)

	fmt.Println("link one name: ", linkOneName)
	fmt.Println("link one:      ", linkOne)
	fmt.Println("link two name: ", linkTwoName)
	fmt.Println("link two:      ", linkTwo)

	fmt.Println("event notes:   ", eventNotes)

	err = h.Tmpl.ExecuteTemplate(w, "events.html", nil)
	if err != nil {
		log.Println("   Unable to open events page: ", err)
		return
	}

}
