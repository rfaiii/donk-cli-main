// Package node defines the local node/device registry used by the UI and
// agent tools. This is intentionally dependency-light so it can be shared
// without import cycles.
package node

import (
	"fmt"
	"sync"
)

// DeviceStatus represents the connectivity state of a node/device.
type DeviceStatus string

const (
	DeviceStatusOffline DeviceStatus = "offline"
	DeviceStatusOnline  DeviceStatus = "online"
	DeviceStatusError   DeviceStatus = "error"
)

// Device represents a discoverable node/device in the local network or via IP.
type Device struct {
	ID             string
	Name           string
	Nickname       string
	ConnectionType string
	Address        string
	Status         DeviceStatus
}

// DisplayName returns the best human-readable name for the device.
func (d Device) DisplayName() string {
	if d.Nickname != "" {
		return d.Nickname
	}
	if d.Name != "" {
		return d.Name
	}
	return d.ID
}

// ConnectionLabel returns a short connection descriptor.
func (d Device) ConnectionLabel() string {
	if d.ConnectionType == "" {
		return "local"
	}
	return d.ConnectionType
}

// Manager stores the current set of known devices.
type Manager struct {
	mu       sync.Mutex
	devices  map[string]Device
	onChange func()
}

// NewManager creates a new device manager.
func NewManager() *Manager {
	return &Manager{
		devices: make(map[string]Device),
	}
}

// OnChange registers a callback invoked after mutations.
func (m *Manager) OnChange(fn func()) {
	m.mu.Lock()
	m.onChange = fn
	m.mu.Unlock()
}

func (m *Manager) changed() {
	if fn := func() func() {
		m.mu.Lock()
		defer m.mu.Unlock()
		return m.onChange
	}(); fn != nil {
		fn()
	}
}

// Devices returns a snapshot of known devices.
func (m *Manager) Devices() []Device {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]Device, 0, len(m.devices))
	for _, d := range m.devices {
		out = append(out, d)
	}
	return out
}

// Upsert adds or updates a device.
func (m *Manager) Upsert(device Device) {
	if device.ID == "" {
		device.ID = fmt.Sprintf("%s-%s", device.Name, device.Address)
	}
	m.mu.Lock()
	m.devices[device.ID] = device
	m.mu.Unlock()
	m.changed()
}

// SetStatus updates only the status of a known device.
func (m *Manager) SetStatus(id string, status DeviceStatus) {
	m.mu.Lock()
	if d, ok := m.devices[id]; ok {
		d.Status = status
		m.devices[id] = d
	}
	m.mu.Unlock()
	m.changed()
}

// SetNickname updates the nickname for a known device.
func (m *Manager) SetNickname(id, nickname string) {
	m.mu.Lock()
	if d, ok := m.devices[id]; ok {
		d.Nickname = nickname
		m.devices[id] = d
	}
	m.mu.Unlock()
	m.changed()
}

// Remove deletes a device.
func (m *Manager) Remove(id string) {
	m.mu.Lock()
	delete(m.devices, id)
	m.mu.Unlock()
	m.changed()
}

// EnsureDefault adds the default local device if none exist.
func (m *Manager) EnsureDefault() {
	m.mu.Lock()
	defer m.mu.Unlock()
	if len(m.devices) > 0 {
		return
	}
	m.devices["local"] = Device{
		ID:             "local",
		Name:           "local",
		Nickname:       "Local Device",
		ConnectionType: "local",
		Address:        "127.0.0.1",
		Status:         DeviceStatusOnline,
	}
	if m.onChange != nil {
		m.onChange()
	}
}

// defaultManager is the process-wide device registry used by UI panels.
var defaultManager = NewManager()

// Devices returns a snapshot of the default device registry.
func Devices() []Device {
	return defaultManager.Devices()
}

// UpsertDevice adds or updates a device in the default registry.
func UpsertDevice(device Device) {
	defaultManager.Upsert(device)
}

// SetDeviceStatus updates the status of a device in the default registry.
func SetDeviceStatus(id string, status DeviceStatus) {
	defaultManager.SetStatus(id, status)
}

// SetDeviceNickname updates the nickname of a device in the default registry.
func SetDeviceNickname(id, nickname string) {
	defaultManager.SetNickname(id, nickname)
}

// RemoveDevice deletes a device from the default registry.
func RemoveDevice(id string) {
	defaultManager.Remove(id)
}

// EnsureDefaultDevice adds the default local device if the registry is empty.
func EnsureDefaultDevice() {
	defaultManager.EnsureDefault()
}
