package handlers

import (
	"fmt"
	"html/template"
	"log"
	"strings"
	"time"
)

var funcMap = template.FuncMap{
	"add": func(a, b int) int {
		return a + b
	},
	"formatDuration": func(totalSeconds int) string {
		minutes := totalSeconds / 60
		seconds := totalSeconds % 60

		return fmt.Sprintf("%d:%02d", minutes, seconds)
	},
	"formatMinutes": func(totalSeconds int) int {
		minutes := totalSeconds / 60
		return minutes
	},
	"formatSeconds": func(totalSeconds int) int {
		seconds := totalSeconds % 60
		return seconds
	},
	"formatReleaseDate": func(releaseDate time.Time) string {
		if releaseDate.IsZero() {
			return ""
		}
		return releaseDate.Format("2006-01-02")
	},
	"capitalizeBandName": func(bandName string) string {
  if bandName == "" {
    return ""
  }
		splitName := strings.Split(bandName, " ")
		res := ""

		for _, word := range splitName {
			w := strings.ToUpper(word[:1]) + word[1:] + " "
			res += w
		}
		return strings.TrimSpace(res)
	},
	"boolToYesNo": func(value bool) string {
		if value == true {
			return "Yes"
		} else {
			return "No"
		}
	},
	"smallImagePath": func(dir string) string {
		return dir + "/small.webp"
	},
	"mediumImagePath": func(dir string) string {
		return dir + "/medium.webp"
	},
	"largeImagePath": func(dir string) string {
		return dir + "/large.webp"
	},
	"trimSpace": func(text string) string {
		return strings.TrimSpace(text)
	},
	"formatTime": func(date time.Time) string {
		return date.Format("3:04 PM")
	},
	"currentYear": func() int {
		return time.Now().Year()
	},
	"isAdminToRole": func(isAdmin bool) string {
		if isAdmin == true {
			return "Admin"
		} else {
			return "Member"
		}
	},
	"fetchLargeArtwork": func(url string) string {
		return strings.ReplaceAll(url, "100x100", "1000x1000")
	},
	"convertBlankField": func(field string) string {
		if field == "" {
			return "-"
		}
		return field
	},
}

func HelperSmallImagePath(dir string) string {
	return dir + "/small.webp"
}

func HelperMediumImagePath(dir string) string {
	return dir + "/medium.webp"
}

func HelperLargeImagePath(dir string) string {
	return dir + "/large.webp"
}

func HelperParseTemplates() *template.Template {
	tmpl := template.New("").Funcs(funcMap)

	template.Must(tmpl.ParseGlob("templates/*.html"))

	_, err := tmpl.ParseGlob("templates/partials/*.html")
	if err != nil {
		panic(err)
	}
	template.Must(tmpl.ParseGlob("templates/auth/*.html"))
	template.Must(tmpl.ParseGlob("templates/songs/*.html"))
	// _, err = template.Must(tmpl.ParseGlob("templates/setlists/*.html"))

	_, err = tmpl.ParseGlob("templates/setlists/*.html")
	if err != nil {
		panic(err)
	}
	_, err = tmpl.ParseGlob("templates/profile/*.html")
	if err != nil {
		panic(err)
	}

	return tmpl
}

func HelperMessagesFormatTime(date time.Time) string {
	log.Println("- HelperMessagesFormatTime")
	return date.Format("3:04 PM")
}
