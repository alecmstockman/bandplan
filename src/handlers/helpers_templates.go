package handlers

import (
	"fmt"
	"html/template"
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

	return tmpl
}
