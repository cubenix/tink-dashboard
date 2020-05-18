package client

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/gauravgahlot/tink-wizard/src/pkg/types"
	log "github.com/sirupsen/logrus"
	"github.com/tinkerbell/tink/protos/template"
	"github.com/tinkerbell/tink/protos/workflow"
	"gopkg.in/yaml.v2"
)

// ListWorkflows returns a list of workflows
func ListWorkflows(ctx context.Context) ([]types.Workflow, error) {
	res, err := workflowClient.ListWorkflows(ctx, &workflow.Empty{})
	if err != nil {
		return nil, err
	}

	// updated template names
	updateTemplateNames(ctx)

	workflows := []types.Workflow{}
	var wf *workflow.Workflow
	err = nil
	for wf, err = res.Recv(); err == nil && wf.Template != ""; wf, err = res.Recv() {
		w, _ := GetWorkflow(ctx, wf.GetId(), false)
		w.CreatedAt = time.Unix(wf.GetCreatedAt().Seconds, 0).Local().Format(time.UnixDate)
		workflows = append(workflows, w)
	}
	if err != nil && err != io.EOF {
		log.Fatal(err)
	}
	return workflows, nil
}

// GetWorkflow returns details for the requested workflow ID
func GetWorkflow(ctx context.Context, id string, fillDetails bool) (types.Workflow, error) {
	w, err := workflowClient.GetWorkflow(ctx, &workflow.GetRequest{Id: id})
	if err != nil {
		return types.Workflow{}, err
	}
	if w.Data == "" {
		return types.Workflow{}, fmt.Errorf("No data found for workflow ID: %v", id)
	}
	wf := types.Workflow{
		ID:       id,
		RawData:  w.GetData(),
		Template: templateNames[w.GetTemplate()],
		State:    w.GetState().String(),
	}
	details := *parseWorkflowYAML(w.GetData())
	wf.Name = details.Name
	wf.Timeout = strconv.Itoa(details.GlobalTimeout)
	if fillDetails {
		wf.Details = details
	}
	return wf, nil
}

func parseWorkflowYAML(data string) *types.WorkflowDetails {
	var wf = types.WorkflowDetails{}
	yaml.UnmarshalStrict([]byte(data), &wf)
	return &wf
}

func updateTemplateNames(ctx context.Context) {
	ch := make(chan *template.WorkflowTemplate)
	go receiveTemplates(ctx, ch)
	for len(ch) > 0 {
		<-ch
	}
}
