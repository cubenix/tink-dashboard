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
	"time"

	"github.com/gauravgahlot/tinker/internal/types"
	"github.com/rs/zerolog/log"
	"github.com/tinkerbell/tink/protos/template"
)

// CreateTemplate creates a new workflow template
func (c Client) CreateTemplate(ctx context.Context, name, data string) (string, error) {
	res, err := c.template.CreateTemplate(ctx, &template.WorkflowTemplate{
		Name: name,
		Data: data,
	})
	if err != nil {
		return "", err
	}

	return res.Id, nil
}

// ListTemplates returns a list of workflow templates
func (c Client) ListTemplates(ctx context.Context) ([]types.Template, error) {
	res, err := c.template.ListTemplates(ctx, &template.ListRequest{
		FilterBy: &template.ListRequest_Name{
			Name: "*",
		},
	})
	if err != nil {
		return nil, err
	}

	templates := []types.Template{}
	var tmp *template.WorkflowTemplate
	for tmp, err = res.Recv(); err == nil && tmp.Name != ""; tmp, err = res.Recv() {
		data, err := c.getTemplateData(ctx, tmp.GetId())
		if err == nil && data != "" {
			t := types.Template{
				ID:          tmp.GetId(),
				Name:        tmp.GetName(),
				Data:        data,
				LastUpdated: time.Unix(tmp.UpdatedAt.Seconds, 0).Local().Format(time.UnixDate),
			}

			templates = append(templates, t)
		}
	}

	if err != nil && err != io.EOF {
		log.Error().Err(err)
	}

	return templates, nil
}

// GetTemplate returns details for the requested template ID
func (c Client) GetTemplate(ctx context.Context, id string) (types.Template, error) {
	t, err := c.template.GetTemplate(ctx, &template.GetRequest{GetBy: &template.GetRequest_Id{Id: id}})
	if err != nil {
		return types.Template{}, err
	}

	if t.Data == "" {
		return types.Template{}, fmt.Errorf("no data found for template ID: %v", id)
	}

	tmpl := types.Template{
		ID:   t.GetId(),
		Name: t.GetName(),
		Data: t.GetData(),
	}

	return tmpl, nil
}

// UpdateTemplate updates the give template
func (c Client) UpdateTemplate(ctx context.Context, id string, data string) error {
	_, err := c.template.UpdateTemplate(ctx, &template.WorkflowTemplate{
		Id:   id,
		Data: data,
	})
	if err != nil {
		return err
	}

	return nil
}

func (c Client) getTemplateData(ctx context.Context, id string) (string, error) {
	t, err := c.template.GetTemplate(ctx, &template.GetRequest{GetBy: &template.GetRequest_Id{Id: id}})
	if err != nil {
		return "", err
	}

	if t.Data == "" {
		return "", fmt.Errorf("no data found for template ID: %v", id)
	}

	return t.GetData(), nil
}
