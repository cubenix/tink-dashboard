package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/gauravgahlot/tink-wizard/src/pkg/redis"
	"github.com/gauravgahlot/tink-wizard/src/pkg/types"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"github.com/tinkerbell/tink/protos/hardware"
)

// CreateNewHardware creates a new workflow template
func CreateNewHardware(ctx context.Context, data string) (string, error) {
	_, err := hardwareClient.Push(ctx, &hardware.PushRequest{Data: data})
	if err != nil {
		return "", err
	}
	h := fillHardwareFromJSON(data)
	cache.Set(redis.CacheKeys.Hardwares, h.ID, h)
	return h.ID, nil
}

// ListHardwares returns a list of workflow hardwares
func ListHardwares(ctx context.Context) ([]types.Hardware, error) {
	hws, err := cache.GetAll(redis.CacheKeys.Hardwares)
	if err != nil || hws == nil || len(hws) == 0 {
		return listHardwaresFromServer(ctx)
	}

	hardwares := []types.Hardware{}
	for id, hw := range hws {
		var h types.Hardware
		if err := json.Unmarshal([]byte(hw), &h); err != nil {
			log.Error(err)
			cache.Delete(redis.CacheKeys.Hardwares, id)
			continue
		}
		hardwares = append(hardwares, h)
	}
	return hardwares, nil
}

// GetHardware returns details for the requested hardware ID
func GetHardware(ctx context.Context, id string) (types.Hardware, error) {
	result, err := cache.Get(redis.CacheKeys.Hardwares, id)
	if err != nil || result == "" {
		return getHardwareFromServer(ctx, id)
	}
	var hw types.Hardware
	json.Unmarshal([]byte(result), &hw)
	return hw, nil
}

// UpdateHardware updates the give hardware
func UpdateHardware(ctx context.Context, id string, data string) error {
	_, err := hardwareClient.Push(ctx, &hardware.PushRequest{Data: data})
	if err != nil {
		return err
	}
	result, _ := cache.Get(redis.CacheKeys.Hardwares, id)
	var hw types.Hardware
	json.Unmarshal([]byte(result), &hw)
	hw.Data = data
	if err := cache.Set(redis.CacheKeys.Hardwares, id, data); err != nil {
		log.Error(err)
	}
	return nil
}

func listHardwaresFromServer(ctx context.Context) ([]types.Hardware, error) {
	res, err := hardwareClient.All(ctx, &hardware.Empty{})
	if err != nil {
		return nil, err
	}
	hardwares := []types.Hardware{}
	var hw *hardware.Hardware
	for hw, err = res.Recv(); err == nil && hw.JSON != ""; hw, err = res.Recv() {
		if err != nil {
			log.Error("error receiving hardware data")
			continue
		}
		h := fillHardwareFromJSON(hw.JSON)
		cache.Set(redis.CacheKeys.Hardwares, h.ID, h)
		hardwares = append(hardwares, h)
	}
	if err != nil && err != io.EOF {
		log.Fatal(err)
	}
	return hardwares, nil
}

func getHardwareFromServer(ctx context.Context, id string) (types.Hardware, error) {
	h, err := hardwareClient.ByID(ctx, &hardware.GetRequest{ID: id})
	if err != nil {
		return types.Hardware{}, err
	}
	if h.JSON == "" {
		return types.Hardware{}, fmt.Errorf("no data found for hardware ID: %v", id)
	}
	hw := fillHardwareFromJSON(h.JSON)
	cache.Set(redis.CacheKeys.Hardwares, id, hw)
	return hw, nil
}

func fillHardwareFromJSON(json string) types.Hardware {
	data := gjson.Parse(json)
	return types.Hardware{
		ID:   data.Get("id").String(),
		Data: json,

		// setting hardcoded fields for now
		// TODO: get fields from settings page
		Fields: map[string]string{
			"Architecture":   data.Get("arch").String(),
			"Allow Workflow": data.Get("allow_workflow").String(),
			"MAC":            data.Get("network_ports.0.data.mac").String(),
			"Requested IP":   data.Get("ip_addresses.0.address").String(),
		},
	}
}
