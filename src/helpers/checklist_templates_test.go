package helpers

import (
	"bandplan/src/models"
	"io"
	"os"
	"testing"
)

func TestChecklistTemplatesRender(t *testing.T) {
	workingDirectory, err := os.Getwd()
	if err != nil {
		t.Fatalf("get working directory: %v", err)
	}
	if err := os.Chdir("../.."); err != nil {
		t.Fatalf("change to project root: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(workingDirectory); err != nil {
			t.Errorf("restore working directory: %v", err)
		}
	})

	tmpl := HelperParseTemplates()
	user := models.User{UserID: "user-id", TimeZone: "UTC"}
	band := models.Band{BandID: "band-id"}
	item := models.ChecklistItem{
		ItemID:      "item-id",
		ChecklistID: "checklist-id",
		Name:        "Pack cables",
	}

	pageData := models.ChecklistPageData{
		User: user,
		Band: band,
		Checklist: models.Checklist{
			ChecklistID: "checklist-id",
			Name:        "Show checklist",
			Items:       []models.ChecklistItem{item},
		},
	}
	if err := tmpl.ExecuteTemplate(io.Discard, "checklist.html", pageData); err != nil {
		t.Fatalf("render checklist page: %v", err)
	}

	listsData := models.ChecklistsPageData{
		User: user,
		Band: band,
		Checklists: []models.Checklist{{
			ChecklistID:    "checklist-id",
			Name:           "Show checklist",
			ItemCount:      2,
			CompletedCount: 1,
		}},
	}
	if err := tmpl.ExecuteTemplate(io.Discard, "checklists.html", listsData); err != nil {
		t.Fatalf("render checklists page: %v", err)
	}

	itemData := models.ChecklistItemPage{
		User: user,
		Band: band,
		Item: item,
	}
	if err := tmpl.ExecuteTemplate(io.Discard, "checklist_item_edit", itemData); err != nil {
		t.Fatalf("render checklist item edit: %v", err)
	}
	if err := tmpl.ExecuteTemplate(io.Discard, "checklist_item_update", itemData); err != nil {
		t.Fatalf("render checklist item update: %v", err)
	}
}
