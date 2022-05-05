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

	"github.com/gauravgahlot/tinker/internal/types"
)

type Connection interface {
	// ConnectionOK checks api server connection status.
	// ConnectionOK() bool

	Hardware
	Template
	Workflow
}

type Hardware interface {
	// CreateHardware creates a new workflow hardware configuration
	// returns hardware configuration identifier
	CreateHardware(ctx context.Context, data string) (string, error)

	// ListHardwares returns a list of workflow hardwares
	ListHardwares(ctx context.Context) ([]types.Hardware, error)

	// GetHardware returns details for the requested hardware ID
	GetHardware(ctx context.Context, id string) (types.Hardware, error)

	// UpdateHardware updates the given workflow hardware configuration
	UpdateHardware(ctx context.Context, id string, data string) error
}

type Template interface {
	// CreateTemplate creates a new workflow template
	CreateTemplate(ctx context.Context, name, data string) (string, error)

	// ListTemplates returns a list of workflow templates
	ListTemplates(ctx context.Context) ([]types.Template, error)

	// GetTemplate returns details for the requested template ID
	GetTemplate(ctx context.Context, id string) (types.Template, error)

	// UpdateTemplate updates the give template
	UpdateTemplate(ctx context.Context, id string, data string) error
}

type Workflow interface {
	// CreateWorkflow creates a new workflow with
	// selected template and hardware devices
	CreateWorkflow(ctx context.Context, templateID string, hardware string) (string, error)

	// ListWorkflows returns a list of workflows
	ListWorkflows(ctx context.Context) ([]types.Workflow, error)

	// GetWorkflow returns details for the requested workflow ID
	GetWorkflow(ctx context.Context, id string) (types.Workflow, error)
}
