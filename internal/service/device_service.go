package service

import (
	"context"
	"devices-api/internal/models"
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
