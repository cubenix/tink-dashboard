// Copyright 2022 Tinker codeowners.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package client

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/gauravgahlot/tinker/internal/types"
	"github.com/rs/zerolog/log"
	"github.com/tinkerbell/tink/protos/workflow"
	"gopkg.in/yaml.v2"
)

// ListWorkflows returns a list of workflows
func (c Client) ListWorkflows(ctx context.Context) ([]types.Workflow, error) {
	res, err := c.workflow.ListWorkflows(ctx, &workflow.Empty{})
	if err != nil {
		return nil, err
	}

	workflows := []types.Workflow{}
	var wf *workflow.Workflow
	for wf, err = res.Recv(); err == nil && wf.Template != ""; wf, err = res.Recv() {
		w, err := c.getWorkflow(ctx, wf.GetId())
		if err == nil {
			w.CreatedAt = time.Unix(wf.GetCreatedAt().Seconds, 0).Local().Format(time.UnixDate)
			workflows = append(workflows, w)
		}
	}
	if err != nil && err != io.EOF {
		log.Error().Err(err)
	}

	return workflows, nil
}

// GetWorkflow returns details for the requested workflow ID
func (c Client) GetWorkflow(ctx context.Context, id string) (types.Workflow, error) {
	wf, err := c.getWorkflow(ctx, id)
	if err != nil {
		return types.Workflow{}, err
	}

	return wf, nil
}

// ParseWorkflowTemplate parses a template into workflow details
func (c Client) ParseWorkflowTemplate(data string) *types.WorkflowDetails {
	return parseWorkflowYAML(data)
}

// CreateWorkflow creates a new workflow with
// selected template and hardware devices
func (c Client) CreateWorkflow(ctx context.Context, templateID string, hardware string) (string, error) {
	res, err := c.workflow.CreateWorkflow(ctx, &workflow.CreateRequest{
		Template: templateID,
		Hardware: hardware,
	})

	if err != nil {
		return "", err
	}
	return res.Id, nil
}

func (c Client) getWorkflow(ctx context.Context, id string) (types.Workflow, error) {
	w, err := c.workflow.GetWorkflow(ctx, &workflow.GetRequest{Id: id})
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
