package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/gauravgahlot/tink-wizard/src/client"
	"github.com/gauravgahlot/tink-wizard/src/pkg"
	"github.com/gauravgahlot/tink-wizard/src/pkg/types"
)

type workflow struct {
	templates map[string]*template.Template
}

func (wf workflow) registerRoutes() {
	http.HandleFunc("/workflow/create", wf.createWorkflow)
	http.HandleFunc("/workflow/new", wf.createNewWorkflow)
	http.HandleFunc("/workflow/list", wf.listWorkflow)
	http.HandleFunc("/workflow/define", wf.defineWorkflowTemplateAndDevices)
	http.HandleFunc("/workflow", wf.getWorkflow)
}

func (wf workflow) createWorkflow(w http.ResponseWriter, r *http.Request) {
	templates, err := client.ListTemplates(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	tmplList := types.TemplateList{Templates: templates}
	tmplList.Title = "Workflows"
	err = wf.templates[create].Execute(w, tmplList)
	pkg.CheckError(err, errTemplateExecute)
}

func (wf workflow) createNewWorkflow(w http.ResponseWriter, r *http.Request) {
	var req types.NewWorkflow
	decErr := json.NewDecoder(r.Body).Decode(&req)
	if decErr != nil {
		log.Errorf("bad request: %v", decErr)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	hw, err := json.Marshal(req.Devices)
	if err != nil {
		log.Error(err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	id, err := client.CreateNewWorkflow(context.Background(), req.TemplateID, string(hw))
	if err != nil {
		log.Error(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	io.WriteString(w, id)
}

func (wf workflow) defineWorkflowTemplateAndDevices(w http.ResponseWriter, r *http.Request) {
	// Enable CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	var req types.Get
	decErr := json.NewDecoder(r.Body).Decode(&req)
	if decErr != nil {
		log.Errorf("bad request: %v", decErr)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	tmp, err := client.GetTemplate(context.Background(), req.ID)
	if err != nil {
		log.Error(err)
		http.Error(w, fmt.Sprintf("No template found for ID: %v", req.ID), http.StatusNotFound)
		return
	}
	details := client.ParseWorkflowTemplate(tmp.Data)
	devices := make([]string, len(details.Tasks))
	for i, t := range details.Tasks {
		device := strings.TrimLeft(t.WorkerAddr, "{{.")
		device = strings.TrimRight(device, "}}")
		devices[i] = device
	}
	res, _ := json.Marshal(types.WorkflowDefinition{Data: tmp.Data, Devices: devices})
	io.WriteString(w, string(res))
}

func (wf workflow) getWorkflow(w http.ResponseWriter, r *http.Request) {
	// Enable CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	var req types.Get
	decErr := json.NewDecoder(r.Body).Decode(&req)
	if decErr != nil {
		log.Errorf("bad request: %v", decErr)
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
