package controller

import (
	"html/template"
	"net/http"
)

const publicDir = "src/app/public"

var homeController home

// Startup registers all the HTTP request handlers
func Startup(templates map[string]*template.Template) {
	homeController.homeTemplate = templates["home.html"]
	homeController.registerRoutes()

	http.Handle("/img/", http.FileServer(http.Dir(publicDir)))
}
