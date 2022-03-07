package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/gauravgahlot/tink-dashboard/src/pkg/types"
	log "github.com/sirupsen/logrus"
	"github.com/tinkerbell/tink/pkg"
	"github.com/tinkerbell/tink/protos/hardware"
)

// CreateNewHardware creates a new workflow hardware configuration
// returns hardware configuration identifier
func CreateNewHardware(ctx context.Context, data string) (string, error) {
	// use HardwareWrapper for adapted json marshal/unmarshal code
	var hw pkg.HardwareWrapper
	err := json.Unmarshal([]byte(data), &hw)
	if err != nil {
		return "", err
	}
	_, err = hardwareClient.Push(ctx, &hardware.PushRequest{Data: hw.Hardware})
	if err != nil {
		return "", err
	}

	hardware := fillHardwareFromWrapper(&hw)

	return hardware.ID, nil
}

// ListHardwares returns a list of workflow hardwares
func ListHardwares(ctx context.Context) ([]types.Hardware, error) {
	return listHardwaresFromServer(ctx)
}

// GetHardware returns details for the requested hardware ID
func GetHardware(ctx context.Context, id string) (types.Hardware, error) {
	return getHardwareFromServer(ctx, id)
}

// UpdateHardware updates the given workflow hardware configuration
func UpdateHardware(ctx context.Context, id string, data string) error {
	// use HardwareWrapper for adapted json marshal/unmarshal code
	var hw pkg.HardwareWrapper
	err := json.Unmarshal([]byte(data), &hw)
	if err != nil {
		log.Fatal(err)
	}
	_, err = hardwareClient.Push(ctx, &hardware.PushRequest{Data: hw.Hardware})
	if err != nil {
		return err
	}

	return nil
}

func listHardwaresFromServer(ctx context.Context) ([]types.Hardware, error) {
	res, err := hardwareClient.All(ctx, &hardware.Empty{})
	if err != nil {
		return nil, err
	}

	// use HardwareWrapper for adapted json marshal/unmarshal code
	var hw pkg.HardwareWrapper
	hardwares := []types.Hardware{}
	for hw.Hardware, err = res.Recv(); err == nil && hw.Hardware != nil; hw.Hardware, err = res.Recv() {
		if err != nil {
			log.Error("error receiving hardware data")
			continue
		}

		h := fillHardwareFromWrapper(&hw)
		hardwares = append(hardwares, h)
	}

	if err != nil && err != io.EOF {
		log.Fatal(err)
	}

	return hardwares, nil
}

func getHardwareFromServer(ctx context.Context, id string) (types.Hardware, error) {
	// use HardwareWrapper for adapted json marshal/unmarshal code
	var h pkg.HardwareWrapper
	var err error

	h.Hardware, err = hardwareClient.ByID(ctx, &hardware.GetRequest{Id: id})
	if err != nil {
		return types.Hardware{}, err
	}
	if h.Hardware == nil {
		return types.Hardware{}, fmt.Errorf("no data found for hardware ID: %v", id)
	}

	hw := fillHardwareFromWrapper(&h)
	return hw, nil
}

func fillHardwareFromWrapper(hw *pkg.HardwareWrapper) types.Hardware {
	data, _ := json.Marshal(hw)
	interfaces := hw.GetNetwork().GetInterfaces()
	allowWorkflow := "false"
	if interfaces[0].GetNetboot().GetAllowWorkflow() {
		allowWorkflow = "true"
	}
	return types.Hardware{
		ID:   hw.GetId(),
		Data: string(data),

		// setting hardcoded fields for now
		// TODO: get fields from settings page
		Fields: map[string]string{
			"Architecture":   interfaces[0].GetDhcp().GetArch(),
			"Allow Workflow": allowWorkflow,
			"MAC":            interfaces[0].GetDhcp().GetMac(),
			"Requested IP":   interfaces[0].GetDhcp().GetIp().GetAddress(),
		},
	}
}
