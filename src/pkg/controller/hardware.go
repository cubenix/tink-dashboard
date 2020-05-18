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

type hardware struct {
	templates map[string]*template.Template
}

func (h hardware) registerRoutes() {
	http.HandleFunc("/hardware/create", h.createHardware)
	http.HandleFunc("/hardware/upload", h.uploadHardware)
	http.HandleFunc("/hardware/list", h.listHardware)
	http.HandleFunc("/hardware", h.getHardware)
	http.HandleFunc("/hardware/update", h.updateHardware)
}

func (h hardware) createHardware(w http.ResponseWriter, r *http.Request) {
	data := types.Base{Title: "Hardwares"}
	err := h.templates[create].Execute(w, data)
	pkg.CheckError(err, errTemplateExecute)
}

func (h hardware) uploadHardware(w http.ResponseWriter, r *http.Request) {
	data := types.Base{Title: "Hardwares"}
	err := h.templates[upload].Execute(w, data)
	pkg.CheckError(err, errTemplateExecute)
}

func (h hardware) listHardware(w http.ResponseWriter, r *http.Request) {
	hardwares, err := client.ListHardwares(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	tmplList := types.HardwareList{Hardwares: hardwares}
	tmplList.Title = "Hardwares"
	err = h.templates[list].Execute(w, tmplList)
	pkg.CheckError(err, errTemplateExecute)
}

func (h hardware) getHardware(w http.ResponseWriter, r *http.Request) {
	// Enable CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	var req types.Get
	decErr := json.NewDecoder(r.Body).Decode(&req)
	if decErr != nil {
		log.Errorf("bad request: ", decErr)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	tmp, err := client.GetHardware(context.Background(), req.ID)
	if err != nil {
		log.Error(err)
		http.Error(w, fmt.Sprintf("No hardware found for ID: %v", req.ID), http.StatusNotFound)
		return
	}
	io.WriteString(w, tmp.Data)
}

func (h hardware) updateHardware(w http.ResponseWriter, r *http.Request) {
	var req types.UpdateTemplate
	decErr := json.NewDecoder(r.Body).Decode(&req)
	if decErr != nil {
		log.Errorf("bad request: ", decErr)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	err := client.UpdateHardware(context.Background(), req.Data)
	if err != nil {
		log.Error(err)
		http.Error(w, fmt.Sprintf("Failed to update hardware data. Error: %v", err), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusAccepted)
	io.WriteString(w, "Hardware data updated successfully")
}
