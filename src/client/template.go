package client

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/gauravgahlot/tink-dashboard/src/pkg/types"
	log "github.com/sirupsen/logrus"
	"github.com/tinkerbell/tink/protos/template"
)

// CreateNewTemplate creates a new workflow template
func CreateNewTemplate(ctx context.Context, name, data string) (string, error) {
	res, err := templateClient.CreateTemplate(ctx, &template.WorkflowTemplate{
		Name: name,
		Data: data,
	})
	if err != nil {
		return "", err
	}

	return res.Id, nil
}

// ListTemplates returns a list of workflow templates
func ListTemplates(ctx context.Context) ([]types.Template, error) {
	return listTemplatesFromServer(ctx)
}

// GetTemplate returns details for the requested template ID
func GetTemplate(ctx context.Context, id string) (types.Template, error) {
	return getTemplateFromServer(ctx, id)
}

// UpdateTemplate updates the give template
func UpdateTemplate(ctx context.Context, id string, data string) error {
	_, err := templateClient.UpdateTemplate(ctx, &template.WorkflowTemplate{
		Id:   id,
		Data: data,
	})
	if err != nil {
		return err
	}

	return nil
}

func listTemplatesFromServer(ctx context.Context) ([]types.Template, error) {
	res, err := templateClient.ListTemplates(ctx, &template.ListRequest{
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
		data, err := getTemplateData(ctx, tmp.GetId())
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
		log.Fatal(err)
	}

	return templates, nil
}

func getTemplateData(ctx context.Context, id string) (string, error) {
	t, err := templateClient.GetTemplate(ctx, &template.GetRequest{GetBy: &template.GetRequest_Id{Id: id}})
	if err != nil {
		return "", err
	}

	if t.Data == "" {
		return "", fmt.Errorf("no data found for template ID: %v", id)
	}

	return t.GetData(), nil
}

func getTemplateFromServer(ctx context.Context, id string) (types.Template, error) {
	t, err := templateClient.GetTemplate(ctx, &template.GetRequest{GetBy: &template.GetRequest_Id{Id: id}})
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
