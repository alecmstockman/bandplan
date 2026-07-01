package handlers

import (
	"fmt"
	"html/template"
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
		fmt.Println("\n\n---------------------")
		fmt.Println(bandName)

		splitName := strings.Split(bandName, " ")
		fmt.Println(splitName)

		res := ""

		for _, word := range splitName {
			w := strings.ToUpper(word[:1]) + word[1:] + " "
			res += w
		}

		return strings.TrimSpace(res)
	},
}

func HelperParseTemplates() *template.Template {
	tmpl := template.New("").Funcs(funcMap)

	template.Must(tmpl.ParseGlob("templates/*.html"))
	// err := template.Must(tmpl.ParseGlob("templates/partials/*.html"))
	// if err != nil {
	// 	fmt.Println("\n\n\n", err)
	// }

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
