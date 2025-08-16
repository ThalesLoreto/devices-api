package models

import "time"

type DeviceState string

const (
	StateAvailable DeviceState = "available"
	StateInUse     DeviceState = "in-use"
	StateInactive  DeviceState = "inactive"
)

func (ds DeviceState) isValid() bool {
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
