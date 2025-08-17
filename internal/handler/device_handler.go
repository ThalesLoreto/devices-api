package handler

import (
	"devices-api/internal/models"
	"devices-api/internal/service"
	"devices-api/internal/utils"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// DeviceHandler handles HTTP requests for device operations
type DeviceHandler struct {
	deviceService service.DeviceService
}

func NewDeviceHandler(deviceService service.DeviceService) *DeviceHandler {
	return &DeviceHandler{
		deviceService: deviceService,
	}
}

// CreateDevice handles POST /devices
func (h *DeviceHandler) CreateDevice(w http.ResponseWriter, r *http.Request) {
	var req service.CreateDeviceRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	device, err := h.deviceService.CreateDevice(r.Context(), req)
	if err != nil {
		if strings.Contains(err.Error(), "validation failed") {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		if strings.Contains(err.Error(), "already exists") {
			utils.WriteErrorResponse(w, http.StatusConflict, err.Error())
			return
		}
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to create device")
		return
	}
	utils.WriteJSONResponse(w, http.StatusCreated, device)
}

// GetDevice handles GET /devices/{id}
func (h *DeviceHandler) GetDevice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	device, err := h.deviceService.GetDevice(r.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.WriteErrorResponse(w, http.StatusNotFound, "Device not found")
			return
		}
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get device")
		return
	}
	utils.WriteJSONResponse(w, http.StatusOK, device)
}

// GetAllDevices handles GET /devices
func (h *DeviceHandler) GetAllDevices(w http.ResponseWriter, r *http.Request) {
	// Check for query parameters
	brand := r.URL.Query().Get("brand")
	state := r.URL.Query().Get("state")

	var devices []*models.Device
	var err error

	if brand != "" {
		devices, err = h.deviceService.GetDevicesByBrand(r.Context(), brand)
	} else if state != "" {
		deviceState := models.DeviceState(state)
		if !deviceState.IsValid() {
			utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid device state")
			return
		}
		devices, err = h.deviceService.GetDevicesByState(r.Context(), deviceState)
	} else {
		devices, err = h.deviceService.GetAllDevices(r.Context())
	}

	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get devices")
		return
	}
	utils.WriteJSONResponse(w, http.StatusOK, devices)
}

// UpdateDevice handles PUT /devices/{id} and PATCH /devices/{id}
func (h *DeviceHandler) UpdateDevice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req service.UpdateDeviceRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	device, err := h.deviceService.UpdateDevice(r.Context(), id, req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.WriteErrorResponse(w, http.StatusNotFound, "Device not found")
			return
		}
		if strings.Contains(err.Error(), "cannot update") || strings.Contains(err.Error(), "validation") {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to update device")
		return
	}
	utils.WriteJSONResponse(w, http.StatusOK, device)
}

// DeleteDevice handles DELETE /devices/{id}
func (h *DeviceHandler) DeleteDevice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := h.deviceService.DeleteDevice(r.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.WriteErrorResponse(w, http.StatusNotFound, "Device not found")
			return
		}
		if strings.Contains(err.Error(), "cannot delete") {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to delete device")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
