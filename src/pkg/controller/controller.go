package controller

import (
	"html/template"
	"net/http"
)

const publicDir = "src/app/public"

var (
	homeController     home
	templateController tmpl
	hardwareController hardware
)

// Startup registers all the HTTP request handlers
func Startup(templates map[string]*template.Template) {
	homeController.template = templates["home.html"]
	homeController.registerRoutes()

	templateController.template = templates["template.html"]
	templateController.registerRoutes()

	hardwareController.template = templates["hardware.html"]
	hardwareController.registerRoutes()

	http.Handle("/img/", http.FileServer(http.Dir(publicDir)))
	http.Handle("/plugin/", http.FileServer(http.Dir(publicDir)))
}
