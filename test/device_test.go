package test

import (
	"devices-api/internal/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewDevice(t *testing.T) {
	tests := []struct {
		name        string
		id          string
		deviceName  string
		brand       string
		state       models.DeviceState
		expectError bool
	}{
		{
			name:        "Valid device creation",
			id:          "device-1",
			deviceName:  "iPhone 15",
			brand:       "Apple",
			state:       models.StateAvailable,
			expectError: false,
		},
		{
			name:        "Empty ID",
			id:          "",
			deviceName:  "iPhone 15",
			brand:       "Apple",
			state:       models.StateAvailable,
			expectError: true,
		},
		{
			name:        "Empty name",
			id:          "device-1",
			deviceName:  "",
			brand:       "Apple",
			state:       models.StateAvailable,
			expectError: true,
		},
		{
			name:        "Empty brand",
			id:          "device-1",
			deviceName:  "iPhone 15",
			brand:       "",
			state:       models.StateAvailable,
			expectError: true,
		},
		{
			name:        "Invalid state",
			id:          "device-1",
			deviceName:  "iPhone 15",
			brand:       "Apple",
			state:       models.DeviceState("invalid"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			device, err := models.NewDevice(tt.id, tt.deviceName, tt.brand, tt.state)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, device)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, device)
				if device != nil {
					assert.Equal(t, tt.id, device.ID)
					assert.Equal(t, tt.deviceName, device.Name)
					assert.Equal(t, tt.brand, device.Brand)
					assert.Equal(t, tt.state, device.State)
					assert.NotNil(t, device.CreationTime)
				}
			}
		})
	}
}

func TestDeviceState_IsValid(t *testing.T) {
	tests := []struct {
		state models.DeviceState
		valid bool
	}{
		{models.StateAvailable, true},
		{models.StateInUse, true},
		{models.StateInactive, true},
		{models.DeviceState("invalid"), false},
		{models.DeviceState(""), false},
	}

	for _, tt := range tests {
		t.Run(string(tt.state), func(t *testing.T) {
			assert.Equal(t, tt.valid, tt.state.IsValid())
		})
	}
}

func TestDevice_CanUpdateNameAndBrand(t *testing.T) {
	tests := []struct {
		state     models.DeviceState
		canUpdate bool
	}{
		{models.StateAvailable, true},
		{models.StateInactive, true},
		{models.StateInUse, false},
	}

	for _, tt := range tests {
		t.Run(string(tt.state), func(t *testing.T) {
			device := &models.Device{
				ID:           "test-id",
				Name:         "Test Device",
				Brand:        "Test Brand",
				State:        tt.state,
				CreationTime: time.Now(),
			}
			assert.Equal(t, tt.canUpdate, device.CanUpdateNameAndBrand())
		})
	}
}

func TestDevice_CanDelete(t *testing.T) {
	tests := []struct {
		state     models.DeviceState
		canDelete bool
	}{
		{models.StateAvailable, true},
		{models.StateInactive, true},
		{models.StateInUse, false},
	}

	for _, tt := range tests {
		t.Run(string(tt.state), func(t *testing.T) {
			device := &models.Device{
				ID:           "test-id",
				Name:         "Test Device",
				Brand:        "Test Brand",
				State:        tt.state,
				CreationTime: time.Now(),
			}
			assert.Equal(t, tt.canDelete, device.CanDelete())
		})
	}
}

func TestDevice_UpdateState(t *testing.T) {
	device := &models.Device{
		ID:           "test-id",
		Name:         "Test Device",
		Brand:        "Test Brand",
		State:        models.StateAvailable,
		CreationTime: time.Now(),
	}

	// Test valid state update
	err := device.UpdateState(models.StateInUse)
	assert.NoError(t, err)
	assert.Equal(t, models.StateInUse, device.State)

	// Test invalid state update
	err = device.UpdateState(models.DeviceState("invalid"))
	assert.Error(t, err)
	assert.Equal(t, models.StateInUse, device.State)
}

func TestDevice_UpdateNameAndBrand(t *testing.T) {
	tests := []struct {
		name         string
		initialState models.DeviceState
		newName      string
		newBrand     string
		expectError  bool
	}{
		{
			name:         "Update available device",
			initialState: models.StateAvailable,
			newName:      "New Name",
			newBrand:     "New Brand",
			expectError:  false,
		},
		{
			name:         "Update inactive device",
			initialState: models.StateInactive,
			newName:      "New Name",
			newBrand:     "New Brand",
			expectError:  false,
		},
		{
			name:         "Update in-use device",
			initialState: models.StateInUse,
			newName:      "New Name",
			newBrand:     "New Brand",
			expectError:  true,
		},
		{
			name:         "Empty name",
			initialState: models.StateAvailable,
			newName:      "",
			newBrand:     "New Brand",
			expectError:  true,
		},
		{
			name:         "Empty brand",
			initialState: models.StateAvailable,
			newName:      "New Name",
			newBrand:     "",
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			device := &models.Device{
				ID:           "test-id",
				Name:         "Original Name",
				Brand:        "Original Brand",
				State:        tt.initialState,
				CreationTime: time.Now(),
			}

			err := device.UpdateNameAndBrand(tt.newName, tt.newBrand)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.newName, device.Name)
				assert.Equal(t, tt.newBrand, device.Brand)
			}
		})
	}
}
