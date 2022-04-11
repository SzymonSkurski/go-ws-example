package handlers

import (
	"log"
	"net/http"

	"github.com/CloudyKit/jet/v6"
)

var views = jet.NewSet(
	jet.NewOSFileSystemLoader("./html"), //will scan given directory for *.jet
	jet.InDevelopmentMode(),             //this should be a env variable
)

// var upgradeConnection =

func Home(w http.ResponseWriter, r *http.Request) {
	//here could pass data map instead of nil
	err := renderPage(w, "home.jet", nil)
	if err != nil {
		log.Panicln(err)
	}
}

// renderPage render page from jet template and write response
func renderPage(w http.ResponseWriter, tmpl string, data jet.VarMap) error {
	view, err := views.GetTemplate(tmpl) //read temple from html/<tmpl>.jet file
	if err != nil {
		log.Panicln(err)
		return err //handle as 404 error - not found
	}
	err = view.Execute(w, data, nil) //pass data (map) to template and write
	if err != nil {
		log.Panicln(err)
		return err
	}
	return nil
}
