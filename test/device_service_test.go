package test

import (
	"context"
	"devices-api/internal/models"
	"devices-api/internal/service"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// MockDeviceRepository is a mock implementation of DeviceRepository for testing
type MockDeviceRepository struct {
	devices map[string]*models.Device
	nextID  int
}

func NewMockDeviceRepository() *MockDeviceRepository {
	return &MockDeviceRepository{
		devices: make(map[string]*models.Device),
		nextID:  1,
	}
}

func (m *MockDeviceRepository) Create(ctx context.Context, device *models.Device) error {
	if _, exists := m.devices[device.ID]; exists {
		return errors.New("device already exists")
	}
	m.devices[device.ID] = device
	return nil
}

func (m *MockDeviceRepository) GetByID(ctx context.Context, id string) (*models.Device, error) {
	device, exists := m.devices[id]
	if !exists {
		return nil, errors.New("device not found")
	}
	return device, nil
}

func (m *MockDeviceRepository) GetAll(ctx context.Context) ([]*models.Device, error) {
	var devices []*models.Device
	for _, device := range m.devices {
		devices = append(devices, device)
	}
	return devices, nil
}

func (m *MockDeviceRepository) GetByBrand(ctx context.Context, brand string) ([]*models.Device, error) {
	var devices []*models.Device
	for _, device := range m.devices {
		if device.Brand == brand {
			devices = append(devices, device)
		}
	}
	return devices, nil
}

func (m *MockDeviceRepository) GetByState(ctx context.Context, state models.DeviceState) ([]*models.Device, error) {
	var devices []*models.Device
	for _, device := range m.devices {
		if device.State == state {
			devices = append(devices, device)
		}
	}
	return devices, nil
}

func (m *MockDeviceRepository) Update(ctx context.Context, device *models.Device) error {
	if _, exists := m.devices[device.ID]; !exists {
		return errors.New("device not found")
	}
	m.devices[device.ID] = device
	return nil
}

func (m *MockDeviceRepository) Delete(ctx context.Context, id string) error {
	if _, exists := m.devices[id]; !exists {
		return errors.New("device not found")
	}
	delete(m.devices, id)
	return nil
}

func (m *MockDeviceRepository) Exists(ctx context.Context, id string) (bool, error) {
	_, exists := m.devices[id]
	return exists, nil
}

func TestDeviceService_CreateDevice(t *testing.T) {
	mockRepo := NewMockDeviceRepository()
	deviceService := service.NewDeviceService(mockRepo)
	ctx := context.Background()

	tests := []struct {
		name        string
		request     service.CreateDeviceRequest
		expectError bool
	}{
		{
			name: "Valid device creation",
			request: service.CreateDeviceRequest{
				Name:  "iPhone 15",
				Brand: "Apple",
				State: models.StateAvailable,
			},
			expectError: false,
		},
		{
			name: "Empty name",
			request: service.CreateDeviceRequest{
				Name:  "",
				Brand: "Apple",
				State: models.StateAvailable,
			},
			expectError: true,
		},
		{
			name: "Empty brand",
			request: service.CreateDeviceRequest{
				Name:  "iPhone 15",
				Brand: "",
				State: models.StateAvailable,
			},
			expectError: true,
		},
		{
			name: "Invalid state",
			request: service.CreateDeviceRequest{
				Name:  "iPhone 15",
				Brand: "Apple",
				State: models.DeviceState("invalid"),
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			device, err := deviceService.CreateDevice(ctx, tt.request)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, device)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, device)
				if device != nil {
					assert.Equal(t, tt.request.Name, device.Name)
					assert.Equal(t, tt.request.Brand, device.Brand)
					assert.Equal(t, tt.request.State, device.State)
					assert.NotEmpty(t, device.ID)
				}
			}
		})
	}
}

func TestDeviceService_GetDevice(t *testing.T) {
	mockRepo := NewMockDeviceRepository()
	deviceService := service.NewDeviceService(mockRepo)
	ctx := context.Background()

	// Create a test device
	testDevice := &models.Device{
		ID:           "test-id",
		Name:         "Test Device",
		Brand:        "Test Brand",
		State:        models.StateAvailable,
		CreationTime: time.Now(),
	}
	mockRepo.devices["test-id"] = testDevice

	tests := []struct {
		name        string
		deviceID    string
		expectError bool
	}{
		{
			name:        "Get existing device",
			deviceID:    "test-id",
			expectError: false,
		},
		{
			name:        "Get non-existing device",
			deviceID:    "non-existing",
			expectError: true,
		},
		{
			name:        "Empty device ID",
			deviceID:    "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			device, err := deviceService.GetDevice(ctx, tt.deviceID)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, device)
				if device != nil {
					assert.Equal(t, tt.deviceID, device.ID)
				}
			}
		})
	}
}

func TestDeviceService_UpdateDevice(t *testing.T) {
	mockRepo := NewMockDeviceRepository()
	deviceService := service.NewDeviceService(mockRepo)
	ctx := context.Background()

	// Create test devices
	availableDevice := &models.Device{
		ID:           "available-device",
		Name:         "Available Device",
		Brand:        "Test Brand",
		State:        models.StateAvailable,
		CreationTime: time.Now(),
	}
	inUseDevice := &models.Device{
		ID:           "in-use-device",
		Name:         "In Use Device",
		Brand:        "Test Brand",
		State:        models.StateInUse,
		CreationTime: time.Now(),
	}
	mockRepo.devices["available-device"] = availableDevice
	mockRepo.devices["in-use-device"] = inUseDevice

	tests := []struct {
		name        string
		deviceID    string
		request     service.UpdateDeviceRequest
		expectError bool
	}{
		{
			name:     "Update available device name and brand",
			deviceID: "available-device",
			request: service.UpdateDeviceRequest{
				Name:  stringPtr("New Name"),
				Brand: stringPtr("New Brand"),
			},
			expectError: false,
		},
		{
			name:     "Update device state",
			deviceID: "available-device",
			request: service.UpdateDeviceRequest{
				State: statePtr(models.StateInactive),
			},
			expectError: false,
		},
		{
			name:     "Update in-use device name (should fail)",
			deviceID: "in-use-device",
			request: service.UpdateDeviceRequest{
				Name: stringPtr("New Name"),
			},
			expectError: true,
		},
		{
			name:     "Update non-existing device",
			deviceID: "non-existing",
			request: service.UpdateDeviceRequest{
				Name: stringPtr("New Name"),
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			device, err := deviceService.UpdateDevice(ctx, tt.deviceID, tt.request)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, device)
			}
		})
	}
}

func TestDeviceService_DeleteDevice(t *testing.T) {
	mockRepo := NewMockDeviceRepository()
	deviceService := service.NewDeviceService(mockRepo)
	ctx := context.Background()

	// Create test devices
	availableDevice := &models.Device{
		ID:           "available-device",
		Name:         "Available Device",
		Brand:        "Test Brand",
		State:        models.StateAvailable,
		CreationTime: time.Now(),
	}
	inUseDevice := &models.Device{
		ID:           "in-use-device",
		Name:         "In Use Device",
		Brand:        "Test Brand",
		State:        models.StateInUse,
		CreationTime: time.Now(),
	}
	mockRepo.devices["available-device"] = availableDevice
	mockRepo.devices["in-use-device"] = inUseDevice

	tests := []struct {
		name        string
		deviceID    string
		expectError bool
	}{
		{
			name:        "Delete available device",
			deviceID:    "available-device",
			expectError: false,
		},
		{
			name:        "Delete in-use device (should fail)",
			deviceID:    "in-use-device",
			expectError: true,
		},
		{
			name:        "Delete non-existing device",
			deviceID:    "non-existing",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := deviceService.DeleteDevice(ctx, tt.deviceID)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func statePtr(s models.DeviceState) *models.DeviceState {
	return &s
}
