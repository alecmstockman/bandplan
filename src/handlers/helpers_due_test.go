package handlers

import "testing"

func TestParseOptionalDueFields(t *testing.T) {
	tests := []struct {
		name       string
		date       string
		time       string
		timezone   string
		wantDate   string
		wantTime   string
		wantZone   string
		wantErr    bool
		wantAbsent bool
	}{
		{name: "absent", wantAbsent: true},
		{
			name:     "date and time",
			date:     "2026-10-05",
			time:     "14:30",
			timezone: "America/Chicago",
			wantDate: "2026-10-05",
			wantTime: "14:30",
			wantZone: "America/Chicago",
		},
		{name: "invalid date", date: "not-a-date", timezone: "UTC", wantErr: true},
		{name: "invalid time", time: "25:00", timezone: "UTC", wantErr: true},
		{name: "invalid timezone", date: "2026-10-05", timezone: "invalid", wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dueDate, dueTime, dueTimezone, err := parseOptionalDueFields(test.date, test.time, test.timezone)
			if test.wantErr {
				if err == nil {
					t.Fatal("expected an error")
				}
				return
			}
			if err != nil {
				t.Fatalf("parse due fields: %v", err)
			}
			if test.wantAbsent {
				if dueDate != nil || dueTime != nil || dueTimezone != nil {
					t.Fatal("expected all due fields to be absent")
				}
				return
			}
			if dueDate == nil || dueDate.Format("2006-01-02") != test.wantDate {
				t.Errorf("due date = %v; want %q", dueDate, test.wantDate)
			}
			if dueTime == nil || dueTime.Format("15:04") != test.wantTime {
				t.Errorf("due time = %v; want %q", dueTime, test.wantTime)
			}
			if dueTimezone == nil || *dueTimezone != test.wantZone {
				t.Errorf("due timezone = %v; want %q", dueTimezone, test.wantZone)
			}
		})
	}
}
