package test

import (
	"context"
	"devices-api/internal/config"
	"devices-api/internal/database"
	"devices-api/internal/models"
	"devices-api/internal/repository"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func setupTestRepository(t *testing.T) repository.DeviceRepository {
	cfg := config.Load()
	db, err := database.NewPostgresConnection(cfg.Database)
	if err != nil {
		t.Fatalf("failed to connect to test db: %v", err)
	}
	err = database.RunMigrations(db)
	if err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}
	return repository.NewPostgresDeviceRepository(db)
}

func TestPostgresRepository_DeviceCRUD(t *testing.T) {
	repo := setupTestRepository(t)
	ctx := context.Background()

	device := &models.Device{
		ID:           "test-id",
		Name:         "Test Device",
		Brand:        "Test Brand",
		State:        models.StateAvailable,
		CreationTime: time.Now(),
	}

	// Create
	err := repo.Create(ctx, device)
	assert.NoError(t, err)

	// Get
	got, err := repo.GetByID(ctx, device.ID)
	assert.NoError(t, err)
	assert.NotNil(t, got)
	assert.Equal(t, device.ID, got.ID)
	assert.Equal(t, device.Name, got.Name)
	assert.Equal(t, device.Brand, got.Brand)
	assert.Equal(t, device.State, got.State)

	// Update
	device.Name = "Updated Name"
	device.Brand = "Updated Brand"
	err = repo.Update(ctx, device)
	assert.NoError(t, err)

	updated, err := repo.GetByID(ctx, device.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Name", updated.Name)
	assert.Equal(t, "Updated Brand", updated.Brand)

	// Delete
	err = repo.Delete(ctx, device.ID)
	assert.NoError(t, err)

	deleted, err := repo.GetByID(ctx, device.ID)
	assert.Error(t, err)
	assert.Nil(t, deleted)
}
