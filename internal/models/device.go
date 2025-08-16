package models

import (
	"errors"
	"time"
)

type DeviceState string

const (
	StateAvailable DeviceState = "available"
	StateInUse     DeviceState = "in-use"
	StateInactive  DeviceState = "inactive"
)

func (ds DeviceState) IsValid() bool {
	switch ds {
	case StateAvailable, StateInUse, StateInactive:
		return true
	default:
		return false
	}
}

type Device struct {
	ID           string      `json:"id"`
	Name         string      `json:"name"`
	Brand        string      `json:"brand"`
	State        DeviceState `json:"state"`
	CreationTime time.Time   `json:"creation_time"`
}

func NewDevice(id, name, brand string, state DeviceState) (*Device, error) {
	if !state.IsValid() {
		return nil, errors.New("invalid device state")
	}
	if id == "" {
		return nil, errors.New("id cannot be empty")
	}
	if name == "" {
		return nil, errors.New("name cannot be empty")
	}
	if brand == "" {
		return nil, errors.New("brand cannot be empty")
	}
	return &Device{
		ID:           id,
		Name:         name,
		Brand:        brand,
		State:        state,
		CreationTime: time.Now(),
	}, nil
}

func (d *Device) CanUpdateNameAndBrand() bool {
	return d.State != StateInUse
}

func (d *Device) CanDelete() bool {
	return d.State != StateInUse
}

func (d *Device) UpdateState(newState DeviceState) error {
	if !newState.IsValid() {
		return errors.New("invalid device state")
	}
	d.State = newState
	return nil
}

func (d *Device) UpdateNameAndBrand(newName, newBrand string) error {
	if !d.CanUpdateNameAndBrand() {
		return errors.New("cannot update name and brand while device is in use")
	}
	if newName == "" {
		return errors.New("name cannot be empty")
	}
	if newBrand == "" {
		return errors.New("brand cannot be empty")
	}
	d.Name = newName
	d.Brand = newBrand
	return nil
}
