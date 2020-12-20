package server

import (
	"text/template" // there's html/template, but I don't know the difference so I'll use this
)

const templatePath = `./templates/*`

var Tmpl *template.Template

func initTemplates() error {
	t, err := template.ParseGlob(templatePath)
	if err != nil {
		return err
	}
	Tmpl = t
	return nil
}
