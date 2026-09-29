package helpers

import (
	"unicode/utf8"
)

func ValidateNameEntryLength(name string) bool {
	if name == "" || utf8.RuneCountInString(name) > 100 {
		return false
	}
	return true
}

func ValidateNameEntryMaxLength(name string) bool {
	if utf8.RuneCountInString(name) > 100 {
		return false
	}
	return true
}
