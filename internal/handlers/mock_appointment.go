package handlers

import (
	"net/http"
	"time"

	"clinic-management/internal/services"
	"clinic-management/internal/utils"

	"github.com/gin-gonic/gin"
)

type MockAppointmentHandler struct {
	mockService *services.MockDataService
}

func NewMockAppointmentHandler(mockService *services.MockDataService) *MockAppointmentHandler {
	return &MockAppointmentHandler{
		mockService: mockService,
	}
}

func (h *MockAppointmentHandler) GetAppointments(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User ID not found in token", "")
		return
	}

	role, exists := c.Get("role")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Role not found in token", "")
		return
	}

	var appointments []services.Appointment

	switch role {
	case "CUSTOMER":
		appointments = h.mockService.GetAppointmentsByCustomer(userID.(string))
	case "DOCTOR":
		appointments = h.mockService.GetAppointmentsByDoctor(userID.(string))
	case "RECEPTIONIST", "CLINIC_MANAGER", "OPERATION_MANAGER":
		appointments = h.mockService.GetAllAppointments()
	default:
		utils.ErrorResponse(c, http.StatusForbidden, "Access denied", "")
		return
	}

	utils.SuccessResponse(c, "Appointments retrieved successfully", appointments)
}

func (h *MockAppointmentHandler) GetAppointment(c *gin.Context) {
	appointmentID := c.Param("id")
	
	appointment, err := h.mockService.GetAppointmentByID(appointmentID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Appointment not found", err.Error())
		return
	}

	// Check if user has permission to view this appointment
	userID, _ := c.Get("userID")
	role, _ := c.Get("role")

	switch role {
	case "CUSTOMER":
		if appointment.MaCustomer != userID.(string) {
			utils.ErrorResponse(c, http.StatusForbidden, "Access denied", "")
			return
		}
	case "DOCTOR":
		if appointment.MaBacSi != userID.(string) {
			utils.ErrorResponse(c, http.StatusForbidden, "Access denied", "")
			return
		}
	case "RECEPTIONIST", "CLINIC_MANAGER", "OPERATION_MANAGER":
		// These roles can view all appointments
	default:
		utils.ErrorResponse(c, http.StatusForbidden, "Access denied", "")
		return
	}

	utils.SuccessResponse(c, "Appointment retrieved successfully", appointment)
}

func (h *MockAppointmentHandler) CreateAppointment(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User ID not found in token", "")
		return
	}

	role, exists := c.Get("role")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Role not found in token", "")
		return
	}

	type CreateAppointmentRequest struct {
		MaBacSi      string  `json:"maBacSi" binding:"required"`
		MaPhongKham  string  `json:"maPhongKham" binding:"required"`
		NgayGioKham  string  `json:"ngayGioKham" binding:"required"`
		GhiChu       *string `json:"ghiChu"`
	}

	var req CreateAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request format", err.Error())
		return
	}

	// Validate that doctor exists
	_, err := h.mockService.GetDoctorByUserID(req.MaBacSi)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Doctor not found", err.Error())
		return
	}

	// Validate that clinic exists
	_, err = h.mockService.GetClinicByID(req.MaPhongKham)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Clinic not found", err.Error())
		return
	}

	// For customers, they can only create appointments for themselves
	// For staff, they can create appointments for any customer
	var customerID string
	if role == "CUSTOMER" {
		customerID = userID.(string)
	} else {
		// For staff creating appointments, you might get customer ID from request
		customerID = userID.(string) // Mock - in real app you'd get this from request
	}

	// Mock implementation - generate appointment ID and return success
	newAppointment := map[string]interface{}{
		"maLichKham":   "APT_" + time.Now().Format("20060102150405"),
		"maCustomer":   customerID,
		"maBacSi":      req.MaBacSi,
		"maPhongKham":  req.MaPhongKham,
		"ngayGioKham":  req.NgayGioKham,
		"trangThai":    "DAT_LICH",
		"ghiChu":       req.GhiChu,
		"ngayDat":      time.Now().Format(time.RFC3339),
		"message":      "Appointment created successfully (mock implementation)",
	}

	utils.SuccessResponse(c, "Appointment created successfully", newAppointment)
}

func (h *MockAppointmentHandler) UpdateAppointment(c *gin.Context) {
	appointmentID := c.Param("id")
	
	appointment, err := h.mockService.GetAppointmentByID(appointmentID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Appointment not found", err.Error())
		return
	}

	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request format", err.Error())
		return
	}

	// Check permissions
	userID, _ := c.Get("userID")
	role, _ := c.Get("role")

	switch role {
	case "CUSTOMER":
		if appointment.MaCustomer != userID.(string) {
			utils.ErrorResponse(c, http.StatusForbidden, "Access denied", "")
			return
		}
	case "DOCTOR":
		if appointment.MaBacSi != userID.(string) {
			utils.ErrorResponse(c, http.StatusForbidden, "Access denied", "")
			return
		}
	case "RECEPTIONIST", "CLINIC_MANAGER", "OPERATION_MANAGER":
		// These roles can update all appointments
	default:
		utils.ErrorResponse(c, http.StatusForbidden, "Access denied", "")
		return
	}

	// Mock implementation
	updatedAppointment := map[string]interface{}{
		"maLichKham": appointmentID,
		"message":    "Appointment updated successfully (mock implementation)",
		"updated":    updateData,
	}

	utils.SuccessResponse(c, "Appointment updated successfully", updatedAppointment)
}

func (h *MockAppointmentHandler) CancelAppointment(c *gin.Context) {
	appointmentID := c.Param("id")
	
	appointment, err := h.mockService.GetAppointmentByID(appointmentID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Appointment not found", err.Error())
		return
	}

	// Check permissions
	userID, _ := c.Get("userID")
	role, _ := c.Get("role")

	switch role {
	case "CUSTOMER":
		if appointment.MaCustomer != userID.(string) {
			utils.ErrorResponse(c, http.StatusForbidden, "Access denied", "")
			return
		}
	case "DOCTOR":
		if appointment.MaBacSi != userID.(string) {
			utils.ErrorResponse(c, http.StatusForbidden, "Access denied", "")
			return
		}
	case "RECEPTIONIST", "CLINIC_MANAGER", "OPERATION_MANAGER":
		// These roles can cancel all appointments
	default:
		utils.ErrorResponse(c, http.StatusForbidden, "Access denied", "")
		return
	}

	// Mock implementation
	cancelledAppointment := map[string]interface{}{
		"maLichKham": appointmentID,
		"trangThai":  "DA_HUY",
		"message":    "Appointment cancelled successfully (mock implementation)",
	}

	utils.SuccessResponse(c, "Appointment cancelled successfully", cancelledAppointment)
}