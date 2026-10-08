package handlers

import (
	"bandplan/src/models"
	"bytes"
	"html/template"
	"os"
	"strings"
	"testing"
	"time"
)

func TestFormatTimeValue(t *testing.T) {
	formatTimeValue := funcMap["formatTimeValue"].(func(*time.Time, string) string)
	value := time.Date(2026, time.September, 17, 15, 30, 0, 0, time.UTC)

	tests := []struct {
		name     string
		value    *time.Time
		timeZone string
		want     string
	}{
		{name: "local time", value: &value, timeZone: "America/Chicago", want: "10:30"},
		{name: "invalid timezone fallback", value: &value, timeZone: "invalid", want: "15:30"},
		{name: "nil time", value: nil, timeZone: "UTC", want: ""},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := formatTimeValue(test.value, test.timeZone); got != test.want {
				t.Errorf("formatTimeValue() = %q; want %q", got, test.want)
			}
		})
	}
}

func TestFormatEventAges(t *testing.T) {
	formatEventAges := funcMap["formatEventAges"].(func(models.EventAges) string)

	tests := []struct {
		ages models.EventAges
		want string
	}{
		{ages: models.EventNA, want: "N/A"},
		{ages: models.EventAllAges, want: "All Ages"},
		{ages: models.Event18Plus, want: "18+"},
		{ages: models.Event21Plus, want: "21+"},
		{ages: models.EventAges("unknown"), want: "N/A"},
	}

	for _, test := range tests {
		if got := formatEventAges(test.ages); got != test.want {
			t.Errorf("formatEventAges(%q) = %q; want %q", test.ages, got, test.want)
		}
	}
}

func TestEventTemplatesParse(t *testing.T) {
	_, err := template.New("").Funcs(funcMap).ParseFiles(
		"../../templates/events/event_create.html",
		"../../templates/events/event-edit.html",
		"../../templates/events/event.html",
	)
	if err != nil {
		t.Fatalf("parse event templates: %v", err)
	}
}

func TestAllTemplatesParse(t *testing.T) {
	tmpl := template.New("").Funcs(funcMap)
	patterns := []string{
		"../../templates/*.html",
		"../../templates/partials/*.html",
		"../../templates/auth/*.html",
		"../../templates/home/*.html",
		"../../templates/songs/*.html",
		"../../templates/chats/*.html",
		"../../templates/transitions/*.html",
		"../../templates/breaks/*.html",
		"../../templates/events/*.html",
		"../../templates/bands/*.html",
		"../../templates/setlists/*.html",
		"../../templates/todo/*.html",
		"../../templates/profile/*.html",
	}

	for _, pattern := range patterns {
		if _, err := tmpl.ParseGlob(pattern); err != nil {
			t.Fatalf("parse templates matching %s: %v", pattern, err)
		}
	}
}

func TestNativeFallbackFormsIncludeCSRFToken(t *testing.T) {
	tests := []struct {
		path   string
		action string
	}{
		{path: "../../templates/chats/chat_create.html", action: `action="/chats/create"`},
		{path: "../../templates/checklist_create.html", action: `action="/checklist/add"`},
		{path: "../../templates/goals_create.html", action: `action="/goal/add"`},
		{path: "../../templates/setlists/setlist-add.html", action: `action="/setlists/create"`},
		{path: "../../templates/setlists/setlist-edit.html", action: `action="/setlist/update"`},
		{path: "../../templates/todo/todo_create.html", action: `action="/todo/add"`},
		{path: "../../templates/todo/todo.html", action: `action="/todo/item/add"`},
		{path: "../../templates/profile/settings.html", action: `action="/logout"`},
	}

	for _, test := range tests {
		t.Run(test.path, func(t *testing.T) {
			contents, err := os.ReadFile(test.path)
			if err != nil {
				t.Fatalf("read template: %v", err)
			}

			templateContents := string(contents)
			if !strings.Contains(templateContents, test.action) {
				t.Errorf("template does not contain native fallback %s", test.action)
			}
			if !strings.Contains(templateContents, `name="gorilla.csrf.Token"`) {
				t.Error("template does not contain a CSRF form field")
			}
		})
	}
}

func TestTodoItemPopupMarkup(t *testing.T) {
	contents, err := os.ReadFile("../../templates/todo/todo.html")
	if err != nil {
		t.Fatalf("read todo template: %v", err)
	}

	templateContents := string(contents)
	wants := []string{
		`id="todo-item-box-popup"`,
		`id="todo-item-popup-edit"`,
		`class="todo-list-item-card"`,
		`class="todo-item-popup-trigger"`,
		`href="/todo/item/delete?item-id={{ .ItemID }}"`,
		`aria-haspopup="dialog"`,
		`data-item-name="{{ .Name }}"`,
		`todoItemPopup.showModal()`,
	}
	for _, want := range wants {
		if !strings.Contains(templateContents, want) {
			t.Errorf("todo template does not contain %s", want)
		}
	}

	if strings.Count(templateContents, `id="todo-item-box-popup"`) != 1 {
		t.Error("todo template should contain one shared item popup")
	}
}

func TestTodoTemplateRendersItemData(t *testing.T) {
	tmpl := template.New("").Funcs(funcMap)
	if _, err := tmpl.ParseGlob("../../templates/partials/*.html"); err != nil {
		t.Fatalf("parse partial templates: %v", err)
	}
	if _, err := tmpl.ParseFiles("../../templates/todo/todo.html"); err != nil {
		t.Fatalf("parse todo template: %v", err)
	}

	assigneeID := "user-1"
	timezone := "America/Chicago"
	dueDate := time.Date(2026, time.October, 8, 0, 0, 0, 0, time.UTC)
	dueTime := time.Date(1, time.January, 1, 14, 30, 0, 0, time.UTC)
	data := models.ToDoListPage{
		User: models.User{TimeZone: timezone},
		List: models.ToDoList{
			ToDoListID: "list-1",
			Items: []models.ToDoItem{
				{
					ItemID:      "item-1",
					Name:        "Book rehearsal",
					Body:        "Confirm the room",
					AssignedTo:  &assigneeID,
					DueDate:     &dueDate,
					DueTime:     &dueTime,
					DueTimezone: &timezone,
				},
			},
		},
		Members: []models.User{{UserID: assigneeID, DisplayName: "Band Member"}},
	}

	var output bytes.Buffer
	if err := tmpl.ExecuteTemplate(&output, "todo.html", data); err != nil {
		t.Fatalf("render todo template: %v", err)
	}

	rendered := output.String()
	for _, want := range []string{
		`href="/todo/item/delete?item-id=item-1"`,
		`data-item-name="Book rehearsal"`,
		`data-assigned-to="user-1"`,
		`data-due-date="2026-10-08"`,
		`data-due-time="14:30"`,
	} {
		if !strings.Contains(rendered, want) {
			t.Errorf("rendered todo template does not contain %s", want)
		}
	}
}

func TestBandsNavigationAndActions(t *testing.T) {
	sidebar, err := os.ReadFile("../../templates/partials/right_sidebar.html")
	if err != nil {
		t.Fatalf("read right sidebar template: %v", err)
	}
	if !strings.Contains(string(sidebar), `href="/bands"`) {
		t.Error("right sidebar does not link to /bands")
	}

	page, err := os.ReadFile("../../templates/profile/bands.html")
	if err != nil {
		t.Fatalf("read bands template: %v", err)
	}
	pageContents := string(page)
	if strings.Contains(pageContents, "<form") || strings.Contains(pageContents, "hx-") {
		t.Error("bands page actions should not submit until their handlers are implemented")
	}
	if strings.Count(pageContents, `type="button"`) != 2 {
		t.Error("bands page should render inert Join Band and Create New Band buttons")
	}
	if !strings.Contains(pageContents, `href="/band?band-id={{ .BandID }}"`) {
		t.Error("band rows should link to their band detail page")
	}

	bandPage, err := os.ReadFile("../../templates/bands/band.html")
	if err != nil {
		t.Fatalf("read band template: %v", err)
	}
	bandPageContents := string(bandPage)
	if !strings.Contains(bandPageContents, `href="/bands"`) {
		t.Error("band page does not link back to /bands")
	}
	if strings.Contains(bandPageContents, "<form") || strings.Contains(bandPageContents, "hx-") || strings.Contains(bandPageContents, "<button") {
		t.Error("band page controls should remain inert until their handlers are implemented")
	}
}
