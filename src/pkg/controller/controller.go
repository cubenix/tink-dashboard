package controller

import (
	"html/template"
	"net/http"
)

const (
	publicDir  = "src/app/public"
	index      = "index"
	about      = "about"
	create     = "create"
	upload     = "upload"
	details    = "details"
	list       = "list"
	ext        = ".html"
	hwPrefix   = "hw-"
	tmplPrefix = "tmpl-"
	wfPrefix   = "wf-"
)

var view = struct {
	index, about, create, upload, details, list string
}{
	index:   index + ext,
	about:   about + ext,
	create:  create + ext,
	upload:  upload + ext,
	details: details + ext,
	list:    list + ext,
}

var (
	homeController     home
	templateController tmpl
	hardwareController hardware
	workflowController workflow
)

// Startup registers all the HTTP request handlers
func Startup(templates map[string]*template.Template) {
	homeController.templates = map[string]*template.Template{
		index: templates[view.index],
		about: templates[view.about],
	}
	homeController.registerRoutes()

	templateController.templates = map[string]*template.Template{
		create: templates[tmplPrefix+view.create],
		upload: templates[tmplPrefix+view.upload],
		list:   templates[tmplPrefix+view.list],
	}
	templateController.registerRoutes()

	hardwareController.templates = map[string]*template.Template{
		create: templates[hwPrefix+view.create],
		upload: templates[hwPrefix+view.upload],
		list:   templates[hwPrefix+view.list],
	}
	hardwareController.registerRoutes()

	workflowController.templates = map[string]*template.Template{
		create:  templates[wfPrefix+view.create],
		list:    templates[wfPrefix+view.list],
		details: templates[wfPrefix+view.details],
	}
	workflowController.registerRoutes()

	http.Handle("/css/", http.FileServer(http.Dir(publicDir)))
	http.Handle("/js/", http.FileServer(http.Dir(publicDir)))
	http.Handle("/img/", http.FileServer(http.Dir(publicDir)))
	http.Handle("/plugin/", http.FileServer(http.Dir(publicDir)))
}
