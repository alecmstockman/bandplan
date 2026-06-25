package handlers

import "html/template"

var funcMap = template.FuncMap{
	"add": func(a, b int) int {
		return a + b
	},
}

func HelperParseTemplates() *template.Template {
	tmpl := template.New("").Funcs(funcMap)

	template.Must(tmpl.ParseGlob("templates/*.html"))
	template.Must(tmpl.ParseGlob("templates/songs/*.html"))

	return tmpl
}
