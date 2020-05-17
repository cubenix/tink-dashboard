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
