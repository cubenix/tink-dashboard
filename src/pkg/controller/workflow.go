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

type workflow struct {
	templates map[string]*template.Template
}

func (wf workflow) registerRoutes() {
	http.HandleFunc("/workflow/list", wf.listWorkflow)
	http.HandleFunc("/workflow", wf.getWorkflow)
}

func (wf workflow) getWorkflow(w http.ResponseWriter, r *http.Request) {
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
	workflow, err := client.GetWorkflow(context.Background(), req.ID, false)
	if err != nil {
		log.Error(err)
		http.Error(w, fmt.Sprintf("No workflow data found for ID: %v", req.ID), http.StatusNotFound)
		return
	}
	io.WriteString(w, workflow.RawData)
}

func (wf workflow) listWorkflow(w http.ResponseWriter, r *http.Request) {
	workflows, err := client.ListWorkflows(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	wfList := types.WorkflowList{Workflows: workflows}
	wfList.Title = "Workflows"
	err = wf.templates[list].Execute(w, wfList)
	pkg.CheckError(err, errTemplateExecute)
}
