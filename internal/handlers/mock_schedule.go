package handlers

import (
	"fmt"
	"net/http"
	"time"

	"clinic-management/internal/services"
	"clinic-management/internal/utils"

	"github.com/gin-gonic/gin"
)

type MockScheduleHandler struct {
	mockService *services.MockDataService
}

func NewMockScheduleHandler(mockService *services.MockDataService) *MockScheduleHandler {
	return &MockScheduleHandler{
		mockService: mockService,
	}
}

type CreateScheduleRequest struct {
	MaBacSi      string `json:"maBacSi" binding:"required"`
	MaPhongKham  string `json:"maPhongKham" binding:"required"`
	NgayLamViec  string `json:"ngayLamViec" binding:"required"` // Format: 2024-01-15
	GioBatDau    string `json:"gioBatDau" binding:"required"`   // Format: 08:00
	GioKetThuc   string `json:"gioKetThuc" binding:"required"`  // Format: 17:00
}

type UpdateScheduleRequest struct {
	MaBacSi      *string `json:"maBacSi"`
	MaPhongKham  *string `json:"maPhongKham"`
	NgayLamViec  *string `json:"ngayLamViec"`
	GioBatDau    *string `json:"gioBatDau"`
	GioKetThuc   *string `json:"gioKetThuc"`
	TrangThai    *string `json:"trangThai"`
}

type AssignDoctorRequest struct {
	MaBacSi      string `json:"maBacSi" binding:"required"`
	MaPhongKham  string `json:"maPhongKham" binding:"required"`
	NgayLamViec  string `json:"ngayLamViec" binding:"required"`
	GioBatDau    string `json:"gioBatDau" binding:"required"`
	GioKetThuc   string `json:"gioKetThuc" binding:"required"`
}

// GET /api/v1/schedules - Lấy tất cả lịch làm việc
func (h *MockScheduleHandler) GetSchedules(c *gin.Context) {
	// Lọc theo query parameters
	doctorID := c.Query("doctorId")
	clinicID := c.Query("clinicId")
	date := c.Query("date")

	schedules, err := h.mockService.GetFilteredSchedules(doctorID, clinicID, date)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get schedules", err.Error())
		return
	}

	utils.SuccessResponse(c, "Schedules retrieved successfully", schedules)
}

// POST /api/v1/schedules - Tạo ca làm việc mới
func (h *MockScheduleHandler) CreateSchedule(c *gin.Context) {
	var req CreateScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request format", err.Error())
		return
	}

	// Validate time format
	if !isValidTimeFormat(req.GioBatDau) || !isValidTimeFormat(req.GioKetThuc) {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid time format. Use HH:MM", "")
		return
	}

	// Check if doctor exists
	if _, err := h.mockService.GetUserByID(req.MaBacSi); err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Doctor not found", err.Error())
		return
	}

	// Check if clinic exists
	if _, err := h.mockService.GetClinicByID(req.MaPhongKham); err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Clinic not found", err.Error())
		return
	}

	// Check for conflicts
	if h.mockService.HasScheduleConflict(req.MaBacSi, req.NgayLamViec, req.GioBatDau, req.GioKetThuc) {
		utils.ErrorResponse(c, http.StatusConflict, "Schedule conflict detected", "Doctor already has a schedule at this time")
		return
	}

	schedule, err := h.mockService.CreateSchedule(req.MaBacSi, req.MaPhongKham, req.NgayLamViec, req.GioBatDau, req.GioKetThuc)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create schedule", err.Error())
		return
	}

	utils.SuccessResponse(c, "Schedule created successfully", schedule)
}

// GET /api/v1/schedules/:id - Lấy chi tiết lịch làm việc
func (h *MockScheduleHandler) GetSchedule(c *gin.Context) {
	scheduleID := c.Param("id")

	schedule, err := h.mockService.GetScheduleByID(scheduleID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Schedule not found", err.Error())
		return
	}

	utils.SuccessResponse(c, "Schedule retrieved successfully", schedule)
}

// PUT /api/v1/schedules/:id - Cập nhật lịch làm việc
func (h *MockScheduleHandler) UpdateSchedule(c *gin.Context) {
	scheduleID := c.Param("id")
	var req UpdateScheduleRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request format", err.Error())
		return
	}

	// Validate time formats if provided
	if req.GioBatDau != nil && !isValidTimeFormat(*req.GioBatDau) {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid start time format. Use HH:MM", "")
		return
	}
	if req.GioKetThuc != nil && !isValidTimeFormat(*req.GioKetThuc) {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid end time format. Use HH:MM", "")
		return
	}

	schedule, err := h.mockService.UpdateSchedule(scheduleID, req.MaBacSi, req.MaPhongKham, req.NgayLamViec, req.GioBatDau, req.GioKetThuc, req.TrangThai)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update schedule", err.Error())
		return
	}

	utils.SuccessResponse(c, "Schedule updated successfully", schedule)
}

// DELETE /api/v1/schedules/:id - Xóa lịch làm việc
func (h *MockScheduleHandler) DeleteSchedule(c *gin.Context) {
	scheduleID := c.Param("id")

	err := h.mockService.DeleteSchedule(scheduleID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete schedule", err.Error())
		return
	}

	utils.SuccessResponse(c, "Schedule deleted successfully", nil)
}

// POST /api/v1/schedules/assign - Phân công bác sĩ vào ca (tạo mới hoặc ghi đè)
func (h *MockScheduleHandler) AssignDoctor(c *gin.Context) {
	var req AssignDoctorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request format", err.Error())
		return
	}

	// Validate inputs
	if !isValidTimeFormat(req.GioBatDau) || !isValidTimeFormat(req.GioKetThuc) {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid time format. Use HH:MM", "")
		return
	}

	// Check if doctor exists
	if _, err := h.mockService.GetUserByID(req.MaBacSi); err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Doctor not found", err.Error())
		return
	}

	// Check if clinic exists
	if _, err := h.mockService.GetClinicByID(req.MaPhongKham); err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Clinic not found", err.Error())
		return
	}

	schedule, err := h.mockService.AssignDoctorToSchedule(req.MaBacSi, req.MaPhongKham, req.NgayLamViec, req.GioBatDau, req.GioKetThuc)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to assign doctor", err.Error())
		return
	}

	utils.SuccessResponse(c, "Doctor assigned successfully", schedule)
}

// GET /api/v1/doctors/:id/schedules - Lấy lịch làm việc của bác sĩ cụ thể
func (h *MockScheduleHandler) GetDoctorSchedules(c *gin.Context) {
	doctorID := c.Param("id")
	
	// Validate doctor exists
	if _, err := h.mockService.GetUserByID(doctorID); err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Doctor not found", err.Error())
		return
	}

	schedules := h.mockService.GetSchedulesByDoctor(doctorID)
	utils.SuccessResponse(c, "Doctor schedules retrieved successfully", schedules)
}

// PUT /api/v1/schedules/:id/assign - Thay đổi phân công bác sĩ cho ca làm việc
func (h *MockScheduleHandler) ReassignDoctor(c *gin.Context) {
	scheduleID := c.Param("id")
	
	var req struct {
		MaBacSi string `json:"maBacSi" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request format", err.Error())
		return
	}

	// Check if doctor exists
	if _, err := h.mockService.GetUserByID(req.MaBacSi); err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Doctor not found", err.Error())
		return
	}

	schedule, err := h.mockService.ReassignDoctor(scheduleID, req.MaBacSi)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to reassign doctor", err.Error())
		return
	}

	utils.SuccessResponse(c, "Doctor reassigned successfully", schedule)
}

// Helper function to validate time format (HH:MM)
func isValidTimeFormat(timeStr string) bool {
	_, err := time.Parse("15:04", timeStr)
	return err == nil
}

// GET /api/v1/schedules/conflicts - Kiểm tra xung đột lịch
func (h *MockScheduleHandler) CheckConflicts(c *gin.Context) {
	doctorID := c.Query("doctorId")
	date := c.Query("date")
	startTime := c.Query("startTime")
	endTime := c.Query("endTime")

	if doctorID == "" || date == "" || startTime == "" || endTime == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Missing required parameters", "doctorId, date, startTime, endTime are required")
		return
	}

	hasConflict := h.mockService.HasScheduleConflict(doctorID, date, startTime, endTime)
	
	result := map[string]interface{}{
		"hasConflict": hasConflict,
		"doctorId":    doctorID,
		"date":        date,
		"timeSlot":    fmt.Sprintf("%s - %s", startTime, endTime),
	}

	if hasConflict {
		utils.SuccessResponse(c, "Schedule conflict detected", result)
	} else {
		utils.SuccessResponse(c, "No conflict found", result)
	}
}