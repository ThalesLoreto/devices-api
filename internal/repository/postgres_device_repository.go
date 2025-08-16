package repository

import (
	"context"
	"database/sql"
	"devices-api/internal/models"
	"fmt"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

// PostgresDeviceRepository implements DeviceRepository using PostgreSQL
type PostgresDeviceRepository struct {
	db *sql.DB
}

// NewPostgresDeviceRepository creates a new PostgreSQL device repository
func NewPostgresDeviceRepository(db *sql.DB) *PostgresDeviceRepository {
	return &PostgresDeviceRepository{
		db: db,
	}
}

// Create inserts a new device into the database
func (r *PostgresDeviceRepository) Create(ctx context.Context, device *models.Device) error {
	query := `
		INSERT INTO devices (id, name, brand, state, creation_time)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.ExecContext(ctx, query,
		device.ID,
		device.Name,
		device.Brand,
		string(device.State),
		device.CreationTime,
	)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return fmt.Errorf("device with ID %s already exists", device.ID)
		}
		return fmt.Errorf("failed to create device: %w", err)
	}
	return nil
}

// GetByID retrieves a device by its ID
func (r *PostgresDeviceRepository) GetByID(ctx context.Context, id string) (*models.Device, error) {
	query := `
		SELECT id, name, brand, state, creation_time
		FROM devices
		WHERE id = $1
	`

	var device models.Device
	var stateStr string

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&device.ID,
		&device.Name,
		&device.Brand,
		&stateStr,
		&device.CreationTime,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("device with ID %s not found", id)
		}
		return nil, fmt.Errorf("failed to get device by ID: %w", err)
	}
	device.State = models.DeviceState(stateStr)
	return &device, nil
}

// GetAll retrieves all devices
func (r *PostgresDeviceRepository) GetAll(ctx context.Context) ([]*models.Device, error) {
	query := `
		SELECT id, name, brand, state, creation_time
		FROM devices
		ORDER BY creation_time DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all devices: %w", err)
	}
	defer rows.Close()

	var devices []*models.Device

	for rows.Next() {
		var device models.Device
		var stateStr string

		err := rows.Scan(
			&device.ID,
			&device.Name,
			&device.Brand,
			&stateStr,
			&device.CreationTime,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan device: %w", err)
		}
		device.State = models.DeviceState(stateStr)
		devices = append(devices, &device)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over devices: %w", err)
	}
	return devices, nil
}

// GetByBrand retrieves devices by brand
func (r *PostgresDeviceRepository) GetByBrand(ctx context.Context, brand string) ([]*models.Device, error) {
	query := `
		SELECT id, name, brand, state, creation_time
		FROM devices
		WHERE brand = $1
		ORDER BY creation_time DESC
	`

	rows, err := r.db.QueryContext(ctx, query, brand)
	if err != nil {
		return nil, fmt.Errorf("failed to get devices by brand: %w", err)
	}
	defer rows.Close()

	var devices []*models.Device

	for rows.Next() {
		var device models.Device
		var stateStr string

		err := rows.Scan(
			&device.ID,
			&device.Name,
			&device.Brand,
			&stateStr,
			&device.CreationTime,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan device: %w", err)
		}
		device.State = models.DeviceState(stateStr)
		devices = append(devices, &device)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over devices: %w", err)
	}
	return devices, nil
}

// GetByState retrieves devices by state
func (r *PostgresDeviceRepository) GetByState(ctx context.Context, state models.DeviceState) ([]*models.Device, error) {
	query := `
		SELECT id, name, brand, state, creation_time
		FROM devices
		WHERE state = $1
		ORDER BY creation_time DESC
	`

	rows, err := r.db.QueryContext(ctx, query, string(state))
	if err != nil {
		return nil, fmt.Errorf("failed to get devices by state: %w", err)
	}
	defer rows.Close()

	var devices []*models.Device

	for rows.Next() {
		var device models.Device
		var stateStr string

		err := rows.Scan(
			&device.ID,
			&device.Name,
			&device.Brand,
			&stateStr,
			&device.CreationTime,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan device: %w", err)
		}
		device.State = models.DeviceState(stateStr)
		devices = append(devices, &device)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over devices: %w", err)
	}
	return devices, nil
}

// Update updates an existing device
func (r *PostgresDeviceRepository) Update(ctx context.Context, device *models.Device) error {
	query := `
		UPDATE devices
		SET name = $2, brand = $3, state = $4
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query,
		device.ID,
		device.Name,
		device.Brand,
		string(device.State),
	)
	if err != nil {
		return fmt.Errorf("failed to update device: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("device with ID %s not found", device.ID)
	}
	return nil
}

// Delete removes a device from the database
func (r *PostgresDeviceRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM devices WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete device by ID: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("device with ID %s not found", id)
	}

	return nil
}

// Exists checks if a device exists by ID
func (r *PostgresDeviceRepository) Exists(ctx context.Context, id string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM devices WHERE id = $1)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check device existence by ID: %w", err)
	}
	return exists, nil
}
