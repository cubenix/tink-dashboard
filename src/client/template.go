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

var templateNames = make(map[string]string)

// ListTemplates returns a list of workflow templates
func ListTemplates(ctx context.Context) ([]types.Template, error) {
	ch := make(chan *template.WorkflowTemplate)
	go receiveTemplates(ctx, ch)

	templates := []types.Template{}
	for tmp := range ch {
		templates = append(templates, types.Template{
			ID:          tmp.GetId(),
			Name:        tmp.GetName(),
			LastUpdated: time.Unix(tmp.UpdatedAt.Seconds, 0).Local().Format(time.UnixDate),
		})
	}
	return templates, nil
}

func receiveTemplates(ctx context.Context, ch chan *template.WorkflowTemplate) {
	defer close(ch)
	res, err := templateClient.ListTemplates(ctx, &template.Empty{})
	if err != nil {
		log.Error(err)
		return
	}
	var tmp *template.WorkflowTemplate
	err = nil
	for tmp, err = res.Recv(); err == nil && tmp.Name != ""; tmp, err = res.Recv() {
		ch <- tmp
		// update templateNames
		templateNames[tmp.Id] = tmp.Name
	}
	if err != nil && err != io.EOF {
		log.Fatal(err)
	}
}

// GetTemplate returns details for the requested template ID
func GetTemplate(ctx context.Context, id string) (types.Template, error) {
	t, err := templateClient.GetTemplate(ctx, &template.GetRequest{Id: id})
	if err != nil {
		return types.Template{}, err
	}

	if t.Data == nil {
		return types.Template{}, fmt.Errorf("no data found for template ID: %v", id)
	}

	return types.Template{
		ID:   t.GetId(),
		Name: t.GetName(),
		Data: string(t.GetData()),
	}, nil
}

// UpdateTemplate updates the give template
func UpdateTemplate(ctx context.Context, id string, data string) error {
	_, err := templateClient.UpdateTemplate(ctx, &template.WorkflowTemplate{
		Id:   id,
		Data: []byte(data),
	})
	return err
}
