package handlers

import (
	"net/http"

	"clinic-management/internal/services"
	"clinic-management/internal/utils"

	"github.com/gin-gonic/gin"
)

type MockClinicHandler struct {
	mockService *services.MockDataService
}

func NewMockClinicHandler(mockService *services.MockDataService) *MockClinicHandler {
	return &MockClinicHandler{
		mockService: mockService,
	}
}

func (h *MockClinicHandler) GetClinics(c *gin.Context) {
	clinics := h.mockService.GetAllClinics()
	utils.SuccessResponse(c, "Clinics retrieved successfully", clinics)
}

func (h *MockClinicHandler) GetClinic(c *gin.Context) {
	clinicID := c.Param("id")
	
	clinic, err := h.mockService.GetClinicByID(clinicID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Clinic not found", err.Error())
		return
	}
	
	utils.SuccessResponse(c, "Clinic retrieved successfully", clinic)
}

func (h *MockClinicHandler) GetDoctors(c *gin.Context) {
	clinicID := c.Param("id")
	
	// First check if clinic exists
	_, err := h.mockService.GetClinicByID(clinicID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Clinic not found", err.Error())
		return
	}
	
	doctors := h.mockService.GetDoctorsByClinic(clinicID)
	utils.SuccessResponse(c, "Doctors retrieved successfully", doctors)
}

func (h *MockClinicHandler) GetSchedules(c *gin.Context) {
	clinicID := c.Param("id")
	
	// First check if clinic exists
	_, err := h.mockService.GetClinicByID(clinicID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Clinic not found", err.Error())
		return
	}
	
	schedules := h.mockService.GetSchedulesByClinic(clinicID)
	utils.SuccessResponse(c, "Schedules retrieved successfully", schedules)
}