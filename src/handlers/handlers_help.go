package handlers

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
)

type helpSection struct {
	Title string
	Body  string
}

type helpTopic struct {
	Slug     string
	Title    string
	Feature  string
	Sections []helpSection
}

var helpTopics = []helpTopic{
	{
		Slug:    "bands",
		Title:   "Bands",
		Feature: "/bands",
		Sections: []helpSection{
			{"Join or switch bands", "Open Bands in the right sidebar. Enter a code in the access code field to join a band. If you already belong to multiple bands, select a band to switch to it."},
			{"Invite a bandmate", "If you are a band admin, open Members in the right sidebar and select Generate Access Code. Send the code to your bandmate so they can join."},
		},
	},
	{
		Slug:    "chats",
		Title:   "Chats",
		Feature: "/chats",
		Sections: []helpSection{
			{"Message your band", "Open Chats and select your band's main chat. Type a message in the field at the bottom and send it."},
			{"Reply to a message", "Press and hold a message to open its menu, then select Reply. Your reply will show the original message for context."},
		},
	},
	{
		Slug:    "songs",
		Title:   "Songs",
		Feature: "/songs",
		Sections: []helpSection{
			{"Add or import a song", "In Songs, choose Add Song to enter it manually, or choose Import to find a track on iTunes. Review the details before saving."},
			{"Find your live arrangement", "Search for a song and open it to see its details. Use the live key and tempo fields to keep track of how your band plays it. Add lyrics and notes so bandmates can refer to them at rehearsal."},
		},
	},
	{
		Slug:    "setlists",
		Title:   "Setlists",
		Feature: "/setlists",
		Sections: []helpSection{
			{"Build a setlist", "Create a setlist, then open it and select Add to include songs, transitions, or breaks. Select Order and drag items by their handles, then save the new order."},
			{"Save a PDF", "Open the setlist and select Share. Choose the PDF option, select the details you want to include, then export the file."},
		},
	},
	{
		Slug:    "events",
		Title:   "Events",
		Feature: "/events",
		Sections: []helpSection{
			{"Schedule a gig or rehearsal", "Open Events and create an event. Enter its date, time, and location. In the Invite section, choose the full band or individual members before saving."},
			{"Respond to an invitation", "Open the event from your notifications, Events, or Calendar. If you were invited, use the response options on the event page to share your availability."},
		},
	},
	{
		Slug:    "calendar",
		Title:   "Calendar",
		Feature: "/calendar",
		Sections: []helpSection{
			{"Find an upcoming event", "Open Calendar in the left sidebar under Schedule. Tap a date in the month view to see that day's schedule. Select an event from the day view to open its details."},
			{"Change the view", "Use the calendar controls to switch between month, week, day, and list views."},
		},
	},
	{
		Slug:    "to-do",
		Title:   "To Do",
		Feature: "/todos",
		Sections: []helpSection{
			{"Create and assign a task", "Open To Do and choose New List or open an existing list. Select Add Item, enter the task, and choose an assignee and due date as needed."},
			{"Track completed work", "Use the checkbox next to a task to mark it complete. To reopen a task, tap its checked circle in the Completed view."},
		},
	},
	{
		Slug:    "notifications",
		Title:   "Notifications",
		Feature: "/notifications",
		Sections: []helpSection{
			{"See band updates", "Select the bell in the app header to see notifications about updates such as task assignments and event invitations. Open a notification to go to the related item."},
			{"Enable push notifications", "If you see Enable push on this device on the Notifications page, select it and allow notifications when prompted. On iPhone, open BandPlan from your Home Screen first."},
		},
	},
}

func (h Handler) HandlerHelpPage(w http.ResponseWriter, r *http.Request) {
	data := struct {
		BackURL  string
		HelpURL  string
		ReturnTo string
		Topics   []helpTopic
		Topic    *helpTopic
	}{
		BackURL: "/",
		HelpURL: "/help",
		Topics:  helpTopics,
	}

	returnTo, err := url.ParseRequestURI(r.URL.Query().Get("return"))
	if err == nil && returnTo.Host == "" && returnTo.Scheme == "" &&
		strings.HasPrefix(returnTo.Path, "/") && !strings.HasPrefix(returnTo.Path, "//") &&
		!strings.Contains(returnTo.Path, "\\") {
		data.ReturnTo = returnTo.RequestURI()
		data.BackURL = data.ReturnTo
		data.HelpURL += "?" + url.Values{"return": {data.ReturnTo}}.Encode()
	}

	if slug := r.PathValue("feature"); slug != "" {
		for i := range helpTopics {
			if helpTopics[i].Slug == slug {
				data.Topic = &helpTopics[i]
				break
			}
		}
		if data.Topic == nil {
			http.NotFound(w, r)
			return
		}
	}

	var page bytes.Buffer
	if err := h.Tmpl.ExecuteTemplate(&page, "help.html", data); err != nil {
		slog.Error("render help page", "error", err)
		http.Error(w, "Unable to load help", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(page.Bytes())
}
