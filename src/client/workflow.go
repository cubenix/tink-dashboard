package client

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/gauravgahlot/tink-dashboard/src/pkg/types"
	log "github.com/sirupsen/logrus"
	"github.com/tinkerbell/tink/protos/workflow"
	"gopkg.in/yaml.v2"
)

// ListWorkflows returns a list of workflows
func ListWorkflows(ctx context.Context) ([]types.Workflow, error) {
	return listWorkflowsFromServer(ctx)
}

// GetWorkflow returns details for the requested workflow ID
func GetWorkflow(ctx context.Context, id string, fillDetails bool) (types.Workflow, error) {
	wf, err := getWorkflow(ctx, id)
	if err != nil {
		return types.Workflow{}, err
	}

	if fillDetails {
		wf.Details = *parseWorkflowYAML(wf.RawData)
	}

	return wf, nil
}

// ParseWorkflowTemplate parses a template into workflow details
func ParseWorkflowTemplate(data string) *types.WorkflowDetails {
	return parseWorkflowYAML(data)
}

// CreateNewWorkflow creates a new workflow with
// selected template and hardware devices
func CreateNewWorkflow(ctx context.Context, templateID string, hardware string) (string, error) {
	res, err := workflowClient.CreateWorkflow(ctx, &workflow.CreateRequest{
		Template: templateID,
		Hardware: hardware,
	})
	if err != nil {
		return "", err
	}

	return res.Id, nil
}

func listWorkflowsFromServer(ctx context.Context) ([]types.Workflow, error) {
	res, err := workflowClient.ListWorkflows(ctx, &workflow.Empty{})
	if err != nil {
		return nil, err
	}

	workflows := []types.Workflow{}
	var wf *workflow.Workflow
	for wf, err = res.Recv(); err == nil && wf.Template != ""; wf, err = res.Recv() {
		w, err := getWorkflow(ctx, wf.GetId())
		if err == nil {
			w.CreatedAt = time.Unix(wf.GetCreatedAt().Seconds, 0).Local().Format(time.UnixDate)
			workflows = append(workflows, w)
		}
	}

	if err != nil && err != io.EOF {
		log.Fatal(err)
	}

	return workflows, nil
}

func getWorkflow(ctx context.Context, id string) (types.Workflow, error) {
	w, err := workflowClient.GetWorkflow(ctx, &workflow.GetRequest{Id: id})
	if err != nil {
		return types.Workflow{}, err
	}

	if w.Data == "" {
		return types.Workflow{}, fmt.Errorf("no data found for workflow ID: %v", id)
	}

	wf := types.Workflow{
		ID:      id,
		RawData: w.GetData(),
		State:   w.GetState().String(),
	}

	setNameAndTimeout(&wf)
	return wf, nil
}

func setNameAndTimeout(wf *types.Workflow) {
	details := *parseWorkflowYAML(wf.RawData)
	wf.Name = details.Name
	wf.Timeout = strconv.Itoa(details.GlobalTimeout)
}

func parseWorkflowYAML(data string) *types.WorkflowDetails {
	var wf types.WorkflowDetails
	yaml.UnmarshalStrict([]byte(data), &wf)

	return &wf
}
