package controller

import (
	"html/template"
)

var homeController home

// Startup registers all the HTTP request handlers
func Startup(templates map[string]*template.Template) {
	homeController.homeTemplate = templates["home.html"]
	homeController.registerRoutes()

	// http.Handle("/js/", http.FileServer(http.Dir("app/public")))
	// http.Handle("/img/", http.FileServer(http.Dir("app/public")))
	// http.Handle("/vendor/", http.FileServer(http.Dir("app/public")))
	// http.Handle("/css/", http.FileServer(http.Dir("app/public")))
}
