package device

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Device represents a network device that can be woken up.
type Device struct {
	Name        string `yaml:"name"`
	Mac         string `yaml:"mac"`
	Broadcast   string `yaml:"broadcast"`
	Description string `yaml:"description,omitempty"`
}

// Config represents the structure of the devices configuration file.
type Config struct {
	Devices []Device `yaml:"devices"`
}

// Manager handles device configuration and lookup.
type Manager struct {
	devices map[string]Device
}

// NewManager creates a new device manager from a configuration file.
func NewManager(configPath string) (*Manager, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	manager := &Manager{
		devices: make(map[string]Device),
	}

	for _, device := range config.Devices {
		if device.Name == "" {
			return nil, fmt.Errorf("device name cannot be empty")
		}
		if device.Mac == "" {
			return nil, fmt.Errorf("device MAC address cannot be empty for device: %s", device.Name)
		}
		if device.Broadcast == "" {
			return nil, fmt.Errorf("device broadcast address cannot be empty for device: %s", device.Name)
		}
		manager.devices[device.Name] = device
	}

	return manager, nil
}

// GetDevice retrieves a device by its name.
func (m *Manager) GetDevice(name string) (Device, error) {
	device, exists := m.devices[name]
	if !exists {
		return Device{}, fmt.Errorf("device not found: %s", name)
	}
	return device, nil
}

// ListDevices returns all registered devices.
func (m *Manager) ListDevices() []Device {
	devices := make([]Device, 0, len(m.devices))
	for _, device := range m.devices {
		devices = append(devices, device)
	}
	return devices
}

// HasDevice checks if a device exists by name.
func (m *Manager) HasDevice(name string) bool {
	_, exists := m.devices[name]
	return exists
}
