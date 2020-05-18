package client

import (
	"context"
	"fmt"
	"io"

	"github.com/gauravgahlot/tink-wizard/src/pkg/types"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"github.com/tinkerbell/tink/protos/hardware"
)

// ListHardwares returns a list of workflow hardwares
func ListHardwares(ctx context.Context) ([]types.Hardware, error) {
	res, err := hardwareClient.All(ctx, &hardware.Empty{})
	if err != nil {
		return nil, err
	}

	hardwares := []types.Hardware{}
	var hw *hardware.Hardware
	err = nil
	for hw, err = res.Recv(); err == nil && hw.JSON != ""; hw, err = res.Recv() {
		if err != nil {
			log.Error("Invalid hardware data")
			continue
		}
		// set custom fields here read from config.json
		// hardcoding a few for now
		hardwares = append(hardwares, getHardware(hw.JSON, false))
	}
	if err != nil && err != io.EOF {
		log.Fatal(err)
	}
	return hardwares, nil
}

// GetHardware returns details for the requested hardware ID
func GetHardware(ctx context.Context, id string) (types.Hardware, error) {
	t, err := hardwareClient.ByID(ctx, &hardware.GetRequest{ID: id})
	if err != nil {
		return types.Hardware{}, err
	}
	if t.JSON == "" {
		return types.Hardware{}, fmt.Errorf("No data found of the hardware ID: %v", id)
	}
	return getHardware(t.JSON, true), nil
}

// UpdateHardware updates the give hardware
func UpdateHardware(ctx context.Context, data string) error {
	_, err := hardwareClient.Push(ctx, &hardware.PushRequest{Data: data})
	return err
}

func getHardware(json string, setData bool) types.Hardware {
	data := gjson.Parse(json)
	hw := types.Hardware{
		ID: data.Get("id").String(),
		Fields: map[string]string{
			"Architecture":   data.Get("arch").String(),
			"Allow Workflow": data.Get("allow_workflow").String(),
			"MAC":            data.Get("network_ports.0.data.mac").String(),
			"Requested IP":   data.Get("ip_addresses.0.address").String(),
		},
	}
	if setData {
		hw.Data = json
	}
	return hw
}
