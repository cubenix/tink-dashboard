package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/tinkerbell/portal/src/pkg/redis"
	"github.com/tinkerbell/portal/src/pkg/types"
	"github.com/tinkerbell/tink/protos/template"
)

// CreateNewTemplate creates a new workflow template
func CreateNewTemplate(ctx context.Context, name, data string) (string, error) {
	res, err := templateClient.CreateTemplate(ctx, &template.WorkflowTemplate{
		Name: name,
		Data: []byte(data),
	})
	if err != nil {
		return "", err
	}
	cache.Set(redis.CacheKeys.Templates, res.Id, types.Template{
		ID:          res.Id,
		Name:        name,
		Data:        data,
		LastUpdated: time.Now().Local().Format(time.UnixDate),
	})
	return res.Id, nil
}

// ListTemplates returns a list of workflow templates
func ListTemplates(ctx context.Context) ([]types.Template, error) {
	tmpls, err := cache.GetAll(redis.CacheKeys.Templates)
	if err != nil || tmpls == nil || len(tmpls) == 0 {
		return listTemplatesFromServer(ctx)
	}
	templates := []types.Template{}
	for id, tmpl := range tmpls {
		var t types.Template
		if err := json.Unmarshal([]byte(tmpl), &t); err != nil {
			log.Error(err)
			cache.Delete(redis.CacheKeys.Templates, id)
			continue
		}
		templates = append(templates, t)
	}
	return templates, nil
}

// GetTemplate returns details for the requested template ID
func GetTemplate(ctx context.Context, id string) (types.Template, error) {
	result, err := cache.Get(redis.CacheKeys.Templates, id)
	if err != nil || result == "" {
		return getTemplateFromServer(ctx, id)
	}
	var tmpl types.Template
	json.Unmarshal([]byte(result), &tmpl)
	return tmpl, nil
}

// UpdateTemplate updates the give template
func UpdateTemplate(ctx context.Context, id string, data string) error {
	_, err := templateClient.UpdateTemplate(ctx, &template.WorkflowTemplate{
		Id:   id,
		Data: []byte(data),
	})
	if err != nil {
		return err
	}
	result, _ := cache.Get(redis.CacheKeys.Templates, id)
	var tmpl types.Template
	json.Unmarshal([]byte(result), &tmpl)
	tmpl.Data = data
	if err := cache.Set(redis.CacheKeys.Templates, id, tmpl); err != nil {
		log.Error(err)
	}
	return nil
}

func listTemplatesFromServer(ctx context.Context) ([]types.Template, error) {
	res, err := templateClient.ListTemplates(ctx, &template.Empty{})
	if err != nil {
		return nil, err
	}
	templates := []types.Template{}
	var tmp *template.WorkflowTemplate
	for tmp, err = res.Recv(); err == nil && tmp.Name != ""; tmp, err = res.Recv() {
		data, err := getTemplateData(ctx, tmp.GetId())
		if err == nil && data != nil {
			t := types.Template{
				ID:          tmp.GetId(),
				Name:        tmp.GetName(),
				Data:        string(data),
				LastUpdated: time.Unix(tmp.UpdatedAt.Seconds, 0).Local().Format(time.UnixDate),
			}
			if err := cache.Set(redis.CacheKeys.Templates, t.ID, t); err != nil {
				log.Error(err)
			}
			if err := cache.Set(redis.CacheKeys.TemplateNames, t.ID, t.Name); err != nil {
				log.Error(err)
			}
			templates = append(templates, t)
		}
	}
	if err != nil && err != io.EOF {
		log.Fatal(err)
	}
	return templates, nil
}

func getTemplateData(ctx context.Context, id string) ([]byte, error) {
	t, err := templateClient.GetTemplate(ctx, &template.GetRequest{Id: id})
	if err != nil {
		return nil, err
	}
	if t.Data == nil {
		return nil, fmt.Errorf("no data found for template ID: %v", id)
	}
	return t.GetData(), nil
}

func getTemplateFromServer(ctx context.Context, id string) (types.Template, error) {
	t, err := templateClient.GetTemplate(ctx, &template.GetRequest{Id: id})
	if err != nil {
		return types.Template{}, err
	}
	if t.Data == nil {
		return types.Template{}, fmt.Errorf("no data found for template ID: %v", id)
	}
	tmpl := types.Template{
		ID:          t.GetId(),
		Name:        t.GetName(),
		Data:        string(t.GetData()),
		LastUpdated: time.Unix(t.UpdatedAt.Seconds, 0).Local().Format(time.UnixDate),
	}
	cache.Set(redis.CacheKeys.Templates, id, tmpl)
	return tmpl, nil
}
