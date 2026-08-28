package helpers

import (
	"strings"
	"testing"
)

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

func TestLargeImagePath(t *testing.T) {
	got := LargeImagePath("users/test-user")

	expected := "users/test-user/large.webp"

	if got != expected {
		t.Errorf("LargeImagePath() = %q; want %q", got, expected)
	}
}

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

			expectedLength := len(tt.expectedPrefix) + 6

			if len(got) != expectedLength {
				t.Errorf(
					"MakeSlug(%q) length = %d; want %d",
					tt.input,
					len(got),
					expectedLength,
				)
			}
		})
	}
}
