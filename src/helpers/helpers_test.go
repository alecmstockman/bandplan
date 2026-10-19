package helpers

import (
	"bandplan/src/models"
	"encoding/hex"
	"strings"
	"testing"
	"time"
)

func TestMakeSlug(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedPrefix string
	}{
		{
			name:           "simple title",
			input:          "My Setlist",
			expectedPrefix: "my-setlist-",
		},
		{
			name:           "leading and trailing spaces",
			input:          "   My Setlist   ",
			expectedPrefix: "my-setlist-",
		},
		{
			name:           "multiple special characters",
			input:          "My -- Cool!!! Setlist",
			expectedPrefix: "my-cool-setlist-",
		},
		{
			name:           "mixed case",
			input:          "Friday NIGHT Show",
			expectedPrefix: "friday-night-show-",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MakeSlug(tt.input)

			if !strings.HasPrefix(got, tt.expectedPrefix) {
				t.Errorf(
					"MakeSlug(%q) = %q; expected prefix %q",
					tt.input,
					got,
					tt.expectedPrefix,
				)
			}

			expectedLength := len(tt.expectedPrefix) + 12

			if len(got) != expectedLength {
				t.Errorf(
					"MakeSlug(%q) length = %d; want %d",
					tt.input,
					len(got),
					expectedLength,
				)
			}

			suffix := strings.TrimPrefix(got, tt.expectedPrefix)

			if len(suffix) != 12 {
				t.Fatalf("suffix length = %d; want 12", len(suffix))
			}

			if _, err := hex.DecodeString(suffix[:8]); err != nil {
				t.Errorf("suffix %q is not hexadecimal: %v", suffix, err)
			}
			if _, err := hex.DecodeString(suffix[9:] + "1"); err != nil {
				t.Errorf("suffix %q is not hexadecimal: %v", suffix, err)
			}
			if suffix[8] != '-' {
				t.Errorf("suffix does not contain '-': %q", suffix)
			}
		})
	}
}

func TestNormalizeEmail(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "lowercase email unchanged",
			input:    "test@example.com",
			expected: "test@example.com",
		},
		{
			name:     "uppercase converted to lowercase",
			input:    "TEST@EXAMPLE.COM",
			expected: "test@example.com",
		},
		{
			name:     "leading and trailing spaces removed",
			input:    "   test@example.com   ",
			expected: "test@example.com",
		},
		{
			name:     "mixed case and spaces",
			input:    "   Test.User@Example.COM   ",
			expected: "test.user@example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NormalizeEmail(tt.input)

			if got != tt.expected {
				t.Errorf(
					"NormalizeEmail(%q) = %q; want %q",
					tt.input,
					got,
					tt.expected,
				)
			}
		})
	}
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "regular email",
			input:    "bobross@gmail.com",
			expected: true,
		},
		{
			name:     "capital letters",
			input:    "StevenWilson@gmail.com",
			expected: true,
		},
		{
			name:     "plus addressing",
			input:    "bob+notifications@example.com",
			expected: true,
		},
		{
			name:     "subdomain",
			input:    "bob@mail.example.com",
			expected: true,
		},
		{
			name:     "period in local part",
			input:    "bob.ross@example.com",
			expected: true,
		},
		{
			name:     "hyphenated domain",
			input:    "bob@example-music.com",
			expected: true,
		},
		{
			name:     "empty",
			input:    "",
			expected: false,
		},
		{
			name:     "missing local part",
			input:    "@example.com",
			expected: false,
		},
		{
			name:     "missing at sign",
			input:    "bob.example.com",
			expected: false,
		},
		{
			name:     "multiple at signs",
			input:    "bob@@example.com",
			expected: false,
		},
		{
			name:     "multiple at signs",
			input:    "bob@example.com@",
			expected: false,
		},
		{
			name:     "missing domain",
			input:    "bob@",
			expected: false,
		},
		{
			name:     "domain without period",
			input:    "bob@localhost",
			expected: false,
		},
		{
			name:     "space in local part",
			input:    "bob ross@example.com",
			expected: false,
		},
		{
			name:     "space in domain",
			input:    "bob@example com",
			expected: false,
		},
		{
			name:     "display name",
			input:    "Bob Ross <bob@example.com>",
			expected: false,
		},
		{
			name:     "leading whitespace",
			input:    " bob@example.com",
			expected: false,
		},
		{
			name:     "trailing whitespace",
			input:    "bob@example.com ",
			expected: false,
		},
		{
			name:     "longer than 254 characters",
			input:    strings.Repeat("a", 249) + "@a.com",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := ValidateEmail(tt.input)
			if got != tt.expected {
				t.Errorf(
					"ValidateEmail(%s) = %t; want %t",
					tt.input,
					got,
					tt.expected,
				)
			}
		})
	}
}

func TestSmallImagePath(t *testing.T) {
	got := SmallImagePath("users/test-user")

	expected := "users/test-user/small.webp"

	if got != expected {
		t.Errorf("SmallImagePath() = %q; want %q", got, expected)
	}
}

func TestMediumImagePath(t *testing.T) {
	got := MediumImagePath("users/test-user")

	expected := "users/test-user/medium.webp"

	if got != expected {
		t.Errorf("MediumImagePath() = %q; want %q", got, expected)
	}
}

func TestEventIsPast(t *testing.T) {

	now := time.Now().UTC()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	yesterday := today.AddDate(0, 0, -1)
	twoDaysAgo := today.AddDate(0, 0, -2)
	tomorrow := today.AddDate(0, 0, 1)

	tests := []struct {
		name     string
		input    models.Event
		expected bool
	}{
		{
			name: "event yesterday",
			input: models.Event{
				EventDate: yesterday,
				Timezone:  "UTC",
			},
			expected: true,
		},
		{
			name: "event today",
			input: models.Event{
				EventDate: today,
				Timezone:  "UTC",
			},
			expected: false,
		},
		{
			name: "event tomorrow",
			input: models.Event{
				EventDate: tomorrow,
				Timezone:  "UTC",
			},
			expected: false,
		},
		{
			name: "time earlier today is ignored",
			input: models.Event{
				EventDate: today.Add(time.Hour),
				Timezone:  "UTC",
			},
			expected: false,
		},
		{
			name: "multi-day event ended yesterday",
			input: models.Event{
				EventDate: twoDaysAgo,
				EndDate:   &yesterday,
				Timezone:  "UTC",
			},
			expected: true,
		},
		{
			name: "multi-day event ends today",
			input: models.Event{
				EventDate: yesterday,
				EndDate:   &today,
				Timezone:  "UTC",
			},
			expected: false,
		},
		{
			name: "multi-day event ends tomorrow",
			input: models.Event{
				EventDate: yesterday,
				EndDate:   &tomorrow,
				Timezone:  "UTC",
			},
			expected: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EventIsPast(tt.input)

			if got != tt.expected {
				t.Errorf(
					"EventIsPast(%v) = %t; want %t",
					tt.input,
					got,
					tt.expected,
				)
			}
		})
	}
}

func TestLargeImagePath(t *testing.T) {
	got := LargeImagePath("users/test-user")

	expected := "users/test-user/large.webp"

	if got != expected {
		t.Errorf("LargeImagePath() = %q; want %q", got, expected)
	}
}

func TestGenerateSessionToken(t *testing.T) {
	token, err := GenerateSessionToken()
	if err != nil {
		t.Fatalf("GenerateSessionToken() error = %v", err)
	}

	decoded, err := hex.DecodeString(token)
	if err != nil {
		t.Fatalf("GenerateSessionToken() returned invalid hex: %v", err)
	}
	if len(decoded) != 32 {
		t.Fatalf("GenerateSessionToken() decoded length = %d; want 32", len(decoded))
	}
}

func TestValidateTokenLength(t *testing.T) {
	validToken := strings.Repeat("a1", 32)

	tests := []struct {
		name  string
		token string
		valid bool
	}{
		{name: "valid 32-byte hex token", token: validToken, valid: true},
		{name: "too short", token: validToken[:62], valid: false},
		{name: "too long", token: validToken + "aa", valid: false},
		{name: "non-hex", token: strings.Repeat("z", 64), valid: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateTokenLength(tt.token); got != tt.valid {
				t.Fatalf("ValidateTokenLength() = %v; want %v", got, tt.valid)
			}
		})
	}
}

func TestNormalizeAccessCode(t *testing.T) {
	got := NormalizeAccessCode("  abcd-ef12-3456  ")
	if got != "ABCD-EF12-3456" {
		t.Fatalf("NormalizeAccessCode() = %q; want %q", got, "ABCD-EF12-3456")
	}
}
