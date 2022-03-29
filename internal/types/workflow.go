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

package types

// WorkflowList is the view model for workflow list page
type WorkflowList struct {
	Base
	Workflows []Workflow
}

// Workflow represents a workflow
type Workflow struct {
	ID        string
	Name      string
	State     string
	CreatedAt string
	Template  string
	Timeout   string
	RawData   string
	Details   WorkflowDetails
}

// WorkflowDetails represents a workflow to be executed
type WorkflowDetails struct {
	Version       string `yaml:"version"`
	Name          string `yaml:"name"`
	ID            string `yaml:"id"`
	GlobalTimeout int    `yaml:"global_timeout"`
	Tasks         []Task `yaml:"tasks"`
}

// Task represents a task to be performed in a worflow
type Task struct {
	Name        string            `yaml:"name"`
	WorkerAddr  string            `yaml:"worker"`
	Actions     []Action          `yaml:"actions"`
	Volumes     []string          `yaml:"volumes"`
	Environment map[string]string `yaml:"environment"`
}

// Action is the basic executional unit for a workflow
type Action struct {
	Name        string            `yaml:"name"`
	Image       string            `yaml:"image"`
	Timeout     int64             `yaml:"timeout"`
	Command     []string          `yaml:"command"`
	OnTimeout   []string          `yaml:"on-timeout"`
	OnFailure   []string          `yaml:"on-failure"`
	Volumes     []string          `yaml:"volumes,omitempty"`
	Environment map[string]string `yaml:"environment,omitempty"`
}

// WorkflowDefinition defines the template data
// and device keys for that workflow
type WorkflowDefinition struct {
	Data    string   `json:"data"`
	Devices []string `json:"devices"`
}

// NewWorkflow represents a create new workflow request
type NewWorkflow struct {
	TemplateID string            `json:"templateID"`
	Devices    map[string]string `json:"devices"`
}
