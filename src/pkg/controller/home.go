package controller

import (
	"html/template"
	"net/http"

	"github.com/gauravgahlot/frawn/src/pkg"
	"github.com/gauravgahlot/frawn/src/pkg/types"
)

const (
	errTemplateExecute = "failed to execute template"
)

type home struct {
	homeTemplate *template.Template
}

func (h home) registerRoutes() {
	http.HandleFunc("/", h.handleHome)
	http.HandleFunc("/home", h.handleHome)
}

func (h home) handleHome(w http.ResponseWriter, r *http.Request) {
	data := types.Home{Message: "Hello World!"}
	data.Title = "Home"
	err := h.homeTemplate.Execute(w, data)
	pkg.CheckError(err, errTemplateExecute)
}
