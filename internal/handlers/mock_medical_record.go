package handlers

import (
	"net/http"
	"time"

	"clinic-management/internal/services"
	"clinic-management/internal/utils"

	"github.com/gin-gonic/gin"
)

type MockMedicalRecordHandler struct {
	mockService *services.MockDataService
}

func NewMockMedicalRecordHandler(mockService *services.MockDataService) *MockMedicalRecordHandler {
	return &MockMedicalRecordHandler{
		mockService: mockService,
	}
}

func (h *MockMedicalRecordHandler) GetMedicalRecords(c *gin.Context) {
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

	var medicalRecords []services.MedicalRecord

	switch role {
	case "CUSTOMER":
		medicalRecords = h.mockService.GetMedicalRecordsByCustomer(userID.(string))
	case "DOCTOR":
		// Doctors can see all medical records they created
		allRecords := h.mockService.GetAllMedicalRecords()
		for _, record := range allRecords {
			if record.MaBacSi == userID.(string) {
				medicalRecords = append(medicalRecords, record)
			}
		}
	case "RECEPTIONIST", "CLINIC_MANAGER", "OPERATION_MANAGER":
		medicalRecords = h.mockService.GetAllMedicalRecords()
	default:
		utils.ErrorResponse(c, http.StatusForbidden, "Access denied", "")
		return
	}

	utils.SuccessResponse(c, "Medical records retrieved successfully", medicalRecords)
}

func (h *MockMedicalRecordHandler) GetMedicalRecord(c *gin.Context) {
	recordID := c.Param("id")
	
	medicalRecord, err := h.mockService.GetMedicalRecordByID(recordID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Medical record not found", err.Error())
		return
	}

	// Check if user has permission to view this medical record
	userID, _ := c.Get("userID")
	role, _ := c.Get("role")

	switch role {
	case "CUSTOMER":
		if medicalRecord.MaCustomer != userID.(string) {
			utils.ErrorResponse(c, http.StatusForbidden, "Access denied", "")
			return
		}
	case "DOCTOR":
		if medicalRecord.MaBacSi != userID.(string) {
			utils.ErrorResponse(c, http.StatusForbidden, "Access denied", "")
			return
		}
	case "RECEPTIONIST", "CLINIC_MANAGER", "OPERATION_MANAGER":
		// These roles can view all medical records
	default:
		utils.ErrorResponse(c, http.StatusForbidden, "Access denied", "")
		return
	}

	utils.SuccessResponse(c, "Medical record retrieved successfully", medicalRecord)
}

func (h *MockMedicalRecordHandler) CreateMedicalRecord(c *gin.Context) {
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

	// Only doctors can create medical records
	if role != "DOCTOR" {
		utils.ErrorResponse(c, http.StatusForbidden, "Only doctors can create medical records", "")
		return
	}

	type CreateMedicalRecordRequest struct {
		MaCustomer       string  `json:"maCustomer" binding:"required"`
		MaPhongKham      string  `json:"maPhongKham" binding:"required"`
		TrieuChung       *string `json:"trieuChung"`
		ChanDoan         *string `json:"chanDoan"`
		HuongDanDieuTri  *string `json:"huongDanDieuTri"`
		MaICD10          *string `json:"maICD10"`
		NgayTaiKham      *string `json:"ngayTaiKham"`
	}

	var req CreateMedicalRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request format", err.Error())
		return
	}

	// Validate that customer exists
	_, err := h.mockService.GetCustomerByUserID(req.MaCustomer)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Customer not found", err.Error())
		return
	}

	// Validate that clinic exists
	_, err = h.mockService.GetClinicByID(req.MaPhongKham)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Clinic not found", err.Error())
		return
	}

	// Mock implementation - generate medical record ID and return success
	newMedicalRecord := map[string]interface{}{
		"maHoSo":           "MR_" + time.Now().Format("20060102150405"),
		"maCustomer":       req.MaCustomer,
		"maBacSi":          userID.(string),
		"maPhongKham":      req.MaPhongKham,
		"ngayKham":         time.Now().Format(time.RFC3339),
		"trieuChung":       req.TrieuChung,
		"chanDoan":         req.ChanDoan,
		"huongDanDieuTri":  req.HuongDanDieuTri,
		"maICD10":          req.MaICD10,
		"ngayTaiKham":      req.NgayTaiKham,
		"message":          "Medical record created successfully (mock implementation)",
	}

	utils.SuccessResponse(c, "Medical record created successfully", newMedicalRecord)
}

func (h *MockMedicalRecordHandler) UpdateMedicalRecord(c *gin.Context) {
	recordID := c.Param("id")
	
	medicalRecord, err := h.mockService.GetMedicalRecordByID(recordID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Medical record not found", err.Error())
		return
	}

	userID, _ := c.Get("userID")
	role, _ := c.Get("role")

	// Only doctors can update medical records, and only their own records
	if role != "DOCTOR" {
		utils.ErrorResponse(c, http.StatusForbidden, "Only doctors can update medical records", "")
		return
	}

	if medicalRecord.MaBacSi != userID.(string) {
		utils.ErrorResponse(c, http.StatusForbidden, "You can only update your own medical records", "")
		return
	}

	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request format", err.Error())
		return
	}

	// Mock implementation
	updatedMedicalRecord := map[string]interface{}{
		"maHoSo":  recordID,
		"message": "Medical record updated successfully (mock implementation)",
		"updated": updateData,
	}

	utils.SuccessResponse(c, "Medical record updated successfully", updatedMedicalRecord)
}