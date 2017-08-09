package views

import (
	"html/template"
)

func GetHtml(sw int) (*template.Template, error) {
	//return template.ParseFiles("webUI.html")
	switch sw {
	case 1:
		return template.ParseFiles("view/upload.gtpl")
	case 2:
		return template.ParseFiles("view/webUI.html")
	default:
		return template.ParseFiles("view/upload.gtpl")
	}
}
