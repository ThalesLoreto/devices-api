package service

import (
	"context"
	"devices-api/internal/models"
	"devices-api/internal/repository"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// DeviceService defines the interface for device business logic operations
type DeviceService interface {
	CreateDevice(ctx context.Context, req CreateDeviceRequest) (*models.Device, error)
	GetDevice(ctx context.Context, id string) (*models.Device, error)
	GetAllDevices(ctx context.Context) ([]*models.Device, error)
	GetDevicesByBrand(ctx context.Context, brand string) ([]*models.Device, error)
	GetDevicesByState(ctx context.Context, state models.DeviceState) ([]*models.Device, error)
	UpdateDevice(ctx context.Context, id string, req UpdateDeviceRequest) (*models.Device, error)
	DeleteDevice(ctx context.Context, id string) error
}

// CreateDeviceRequest represents the request to create a new device
type CreateDeviceRequest struct {
	Name  string             `json:"name" validate:"required"`
	Brand string             `json:"brand" validate:"required"`
	State models.DeviceState `json:"state" validate:"required"`
}

// UpdateDeviceRequest represents the request to update a device
type UpdateDeviceRequest struct {
	Name  *string             `json:"name,omitempty"`
	Brand *string             `json:"brand,omitempty"`
	State *models.DeviceState `json:"state,omitempty"`
}

// DeviceServiceImpl implements DeviceService
type DeviceServiceImpl struct {
	deviceRepo repository.DeviceRepository
}

// NewDeviceService creates a new device service
func NewDeviceService(deviceRepo repository.DeviceRepository) DeviceService {
	return &DeviceServiceImpl{
		deviceRepo: deviceRepo,
	}
}

// CreateDevice creates a new device
func (s *DeviceServiceImpl) CreateDevice(ctx context.Context, req CreateDeviceRequest) (*models.Device, error) {
	// Validate input
	if err := s.validateCreateRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Generate unique ID
	deviceID := uuid.New().String()

	// Create device entity
	device, err := models.NewDevice(deviceID, req.Name, req.Brand, req.State)
	if err != nil {
		return nil, fmt.Errorf("failed to create device entity: %w", err)
	}

	// Save to repository
	if err := s.deviceRepo.Create(ctx, device); err != nil {
		return nil, fmt.Errorf("failed to save device: %w", err)
	}
	return device, nil
}

// GetDevice retrieves a device by ID
func (s *DeviceServiceImpl) GetDevice(ctx context.Context, id string) (*models.Device, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("device ID cannot be empty")
	}

	device, err := s.deviceRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get device: %w", err)
	}
	return device, nil
}

// GetAllDevices retrieves all devices
func (s *DeviceServiceImpl) GetAllDevices(ctx context.Context) ([]*models.Device, error) {
	devices, err := s.deviceRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all devices: %w", err)
	}
	return devices, nil
}

// GetDevicesByBrand retrieves devices by brand
func (s *DeviceServiceImpl) GetDevicesByBrand(ctx context.Context, brand string) ([]*models.Device, error) {
	if strings.TrimSpace(brand) == "" {
		return nil, fmt.Errorf("brand cannot be empty")
	}

	devices, err := s.deviceRepo.GetByBrand(ctx, brand)
	if err != nil {
		return nil, fmt.Errorf("failed to get devices by brand: %w", err)
	}
	return devices, nil
}

// GetDevicesByState retrieves devices by state
func (s *DeviceServiceImpl) GetDevicesByState(ctx context.Context, state models.DeviceState) ([]*models.Device, error) {
	if !state.IsValid() {
		return nil, fmt.Errorf("invalid device state: %s", state)
	}

	devices, err := s.deviceRepo.GetByState(ctx, state)
	if err != nil {
		return nil, fmt.Errorf("failed to get devices by state: %w", err)
	}
	return devices, nil
}

// UpdateDevice updates an existing device
func (s *DeviceServiceImpl) UpdateDevice(ctx context.Context, id string, req UpdateDeviceRequest) (*models.Device, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("device ID cannot be empty")
	}

	// Get existing device
	device, err := s.deviceRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get device: %w", err)
	}

	// Apply updates
	if err := s.applyUpdates(device, req); err != nil {
		return nil, fmt.Errorf("failed to apply updates: %w", err)
	}

	// Save updated device
	if err := s.deviceRepo.Update(ctx, device); err != nil {
		return nil, fmt.Errorf("failed to update device: %w", err)
	}
	return device, nil
}

// DeleteDevice deletes a device
func (s *DeviceServiceImpl) DeleteDevice(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		return fmt.Errorf("device ID cannot be empty")
	}

	// Get device to check if it can be deleted
	device, err := s.deviceRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get device: %w", err)
	}

	// Check business rules
	if !device.CanDelete() {
		return fmt.Errorf("cannot delete device in use")
	}

	// Delete device
	if err := s.deviceRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete device: %w", err)
	}
	return nil
}

// validateCreateRequest validates the create device request
func (s *DeviceServiceImpl) validateCreateRequest(req CreateDeviceRequest) error {
	if strings.TrimSpace(req.Name) == "" {
		return fmt.Errorf("device name is required")
	}
	if strings.TrimSpace(req.Brand) == "" {
		return fmt.Errorf("device brand is required")
	}
	if !req.State.IsValid() {
		return fmt.Errorf("invalid device state: %s", req.State)
	}
	return nil
}

// applyUpdates applies the update request to the device
func (s *DeviceServiceImpl) applyUpdates(device *models.Device, req UpdateDeviceRequest) error {
	// Update state if provided
	if req.State != nil {
		if err := device.UpdateState(*req.State); err != nil {
			return err
		}
	}

	// Update name and brand if provided
	if req.Name != nil || req.Brand != nil {
		newName := device.Name
		newBrand := device.Brand

		if req.Name != nil {
			newName = strings.TrimSpace(*req.Name)
		}
		if req.Brand != nil {
			newBrand = strings.TrimSpace(*req.Brand)
		}

		if err := device.UpdateNameAndBrand(newName, newBrand); err != nil {
			return err
		}
	}
	return nil
}
