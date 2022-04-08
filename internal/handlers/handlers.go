package handlers

import (
	"log"
	"net/http"

	"github.com/CloudyKit/jet/v6"
)

var views = jet.NewSet(
	jet.NewOSFileSystemLoader("./html"),
	jet.InDevelopmentMode(),
)

// var upgradeConnection =
//TODO: finish here https://www.udemy.com/course/working-with-websockets-in-go/learn/lecture/24201988#overview 0:44

func Home(w http.ResponseWriter, r *http.Request) {
	err := renderPage(w, "home.jet", nil)
	if err != nil {
		log.Panicln(err)
	}
}

func renderPage(w http.ResponseWriter, tmpl string, data jet.VarMap) error {
	view, err := views.GetTemplate(tmpl)
	if err != nil {
		log.Panicln(err)
		return err
	}
	err = view.Execute(w, data, nil)
	if err != nil {
		log.Panicln(err)
		return err
	}
	return nil
}
