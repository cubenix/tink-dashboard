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

// HardwareList is the view model for hardware list page
type HardwareList struct {
	Base
	Hardwares []Hardware
}

// Hardware represents a worker hardware
type Hardware struct {
	ID     string
	Data   string
	Fields map[string]string
}

// UpdateHardware represents an update request for a hardware
type UpdateHardware struct {
	ID   string
	Name string
	Data string
}
