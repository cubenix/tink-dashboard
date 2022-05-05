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
	"encoding/json"
	"fmt"
	"io"

	"github.com/gauravgahlot/tinker/internal/types"
	"github.com/rs/zerolog/log"
	"github.com/tinkerbell/tink/pkg"
	"github.com/tinkerbell/tink/protos/hardware"
)

// CreateHardware creates a new workflow hardware configuration
// returns hardware configuration identifier
func (c Client) CreateHardware(ctx context.Context, data string) (string, error) {
	// use HardwareWrapper for adapted json marshal/unmarshal code
	var hw pkg.HardwareWrapper
	err := json.Unmarshal([]byte(data), &hw)
	if err != nil {
		return "", err
	}

	_, err = c.hardware.Push(ctx, &hardware.PushRequest{Data: hw.Hardware})
	if err != nil {
		return "", err
	}

	hardware := fillHardwareFromWrapper(&hw)

	return hardware.ID, nil
}

// ListHardwares returns a list of workflow hardwares
func (c Client) ListHardwares(ctx context.Context) ([]types.Hardware, error) {
	res, err := c.hardware.All(ctx, &hardware.Empty{})
	if err != nil {
		return nil, err
	}

	// use HardwareWrapper for adapted json marshal/unmarshal code
	var hw pkg.HardwareWrapper
	hardwares := []types.Hardware{}
	for hw.Hardware, err = res.Recv(); err == nil && hw.Hardware != nil; hw.Hardware, err = res.Recv() {
		if err != nil {
			log.Error().Err(err)
			continue
		}

		h := fillHardwareFromWrapper(&hw)
		hardwares = append(hardwares, h)
	}

	if err != nil && err != io.EOF {
		log.Error().Err(err)
	}

	return hardwares, nil
}

// GetHardware returns details for the requested hardware ID
func (c Client) GetHardware(ctx context.Context, id string) (types.Hardware, error) {
	// use HardwareWrapper for adapted json marshal/unmarshal code
	var h pkg.HardwareWrapper
	var err error
	h.Hardware, err = c.hardware.ByID(ctx, &hardware.GetRequest{Id: id})
	if err != nil {
		return types.Hardware{}, err
	}

	if h.Hardware == nil {
		return types.Hardware{}, fmt.Errorf("no data found for hardware ID: %v", id)
	}

	hw := fillHardwareFromWrapper(&h)

	return hw, nil
}

// UpdateHardware updates the given workflow hardware configuration
func (c Client) UpdateHardware(ctx context.Context, id string, data string) error {
	// use HardwareWrapper for adapted json marshal/unmarshal code
	var hw pkg.HardwareWrapper
	err := json.Unmarshal([]byte(data), &hw)
	if err != nil {
		log.Error().Err(err)
	}

	_, err = c.hardware.Push(ctx, &hardware.PushRequest{Data: hw.Hardware})
	if err != nil {
		return err
	}

	return nil
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
