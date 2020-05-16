package client

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/gauravgahlot/tink-wizard/src/pkg/types"
	log "github.com/sirupsen/logrus"
	"github.com/tinkerbell/tink/protos/template"
)

// ListTemplates returns a list of workflow templates
func ListTemplates(ctx context.Context) ([]types.Template, error) {
	res, err := templateClient.ListTemplates(ctx, &template.Empty{})
	if err != nil {
		return nil, err
	}

	templates := []types.Template{}
	var tmp *template.WorkflowTemplate
	err = nil
	for tmp, err = res.Recv(); err == nil && tmp.Name != ""; tmp, err = res.Recv() {
		templates = append(templates, types.Template{
			ID:          tmp.GetId(),
			Name:        tmp.GetName(),
			LastUpdated: time.Unix(tmp.UpdatedAt.Seconds, 0).Local().Format(time.UnixDate),
		})
	}
	if err != nil && err != io.EOF {
		log.Fatal(err)
	}
	return templates, nil
}

// GetTemplate returns details for the requested template ID
func GetTemplate(ctx context.Context, id string) (types.Template, error) {
	t, err := templateClient.GetTemplate(ctx, &template.GetRequest{Id: id})
	if err != nil {
		return types.Template{}, err
	}

	if t.Data == nil {
		return types.Template{}, fmt.Errorf("No data found of the template ID: %v", id)
	}

	return types.Template{
		ID:   t.GetId(),
		Name: t.GetName(),
		Data: string(t.GetData()),
	}, nil
}
