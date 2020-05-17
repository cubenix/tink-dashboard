package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/gauravgahlot/tink-wizard/src/client"
	"github.com/gauravgahlot/tink-wizard/src/pkg"
	"github.com/gauravgahlot/tink-wizard/src/pkg/types"
)

type tmpl struct {
	template *template.Template
}

func (t tmpl) registerRoutes() {
	http.HandleFunc("/template/list", t.listTemplates)
	http.HandleFunc("/template", t.getTemplate)
	http.HandleFunc("/template/update", t.updateTemplate)
}

func (t tmpl) listTemplates(w http.ResponseWriter, r *http.Request) {
	templates, err := client.ListTemplates(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	tmplList := types.TemplateList{Templates: templates}
	tmplList.Title = "Templates"
	err = t.template.Execute(w, tmplList)
	pkg.CheckError(err, errTemplateExecute)
}

func (t tmpl) getTemplate(w http.ResponseWriter, r *http.Request) {
	// Enable CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	var req types.Get
	decErr := json.NewDecoder(r.Body).Decode(&req)
	if decErr != nil {
		log.Errorf("bad request: ", decErr)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	tmp, err := client.GetTemplate(context.Background(), req.ID)
	if err != nil {
		log.Error(err)
		http.Error(w, fmt.Sprintf("No template found for ID: %v", req.ID), http.StatusNotFound)
		return
	}
	io.WriteString(w, tmp.Data)
}

func (t tmpl) updateTemplate(w http.ResponseWriter, r *http.Request) {
	var req types.UpdateTemplate
	decErr := json.NewDecoder(r.Body).Decode(&req)
	if decErr != nil {
		log.Errorf("bad request: ", decErr)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	err := client.UpdateTemplate(context.Background(), req.ID, req.Data)
	if err != nil {
		log.Error(err)
		http.Error(w, fmt.Sprintf("No template found for ID: %v", req.ID), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusAccepted)
	io.WriteString(w, "Template updated successfully")
}
