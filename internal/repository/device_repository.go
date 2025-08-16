package repository

import (
	"context"
	"devices-api/internal/models"
)

// DeviceRepository defines the interface for device data access operations
type DeviceRepository interface {
	Create(ctx context.Context, device *models.Device) error
	GetByID(ctx context.Context, id string) (*models.Device, error)
	GetByBrand(ctx context.Context, brand string) ([]*models.Device, error)
	GetByState(ctx context.Context, state models.DeviceState) ([]*models.Device, error)
	GetAll(ctx context.Context) ([]*models.Device, error)
	Update(ctx context.Context, device *models.Device) error
	Delete(ctx context.Context, id string) error
	Exists(ctx context.Context, id string) (bool, error)
}
