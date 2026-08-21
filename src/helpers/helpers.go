package helpers

import (
	"strings"
	"unicode"

	"github.com/google/uuid"
)

func MakeSlug(text string) string {
	name := strings.ToLower(strings.TrimSpace(text))

	var b strings.Builder
	lastDash := false

	for _, r := range name {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			b.WriteRune(r)
		} else if !lastDash {
			b.WriteRune('-')
			lastDash = true
		}
	}
	id := "-" + uuid.NewString()[:6]

	return strings.Trim(b.String(), "-") + id
}
