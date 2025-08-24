package handlers

import (
	"net/http"

	"clinic-management/internal/services"
	"clinic-management/internal/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type MockUserHandler struct {
	mockService *services.MockDataService
}

func NewMockUserHandler(mockService *services.MockDataService) *MockUserHandler {
	return &MockUserHandler{
		mockService: mockService,
	}
}

func (h *MockUserHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User ID not found in token", "")
		return
	}

	user, err := h.mockService.GetUserByID(userID.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "User not found", err.Error())
		return
	}

	// Get additional user details based on role
	var profileData map[string]interface{}
	
	baseProfile := map[string]interface{}{
		"maUser":      user.MaUser,
		"hoTen":       user.HoTen,
		"email":       user.Email,
		"soDienThoai": user.SoDienThoai,
		"role":        user.Role,
		"trangThai":   user.TrangThai,
		"ngayTao":     user.NgayTao,
	}

	switch user.Role {
	case "CUSTOMER":
		customer, err := h.mockService.GetCustomerByUserID(user.MaUser)
		if err == nil {
			baseProfile["ngaySinh"] = customer.NgaySinh
			baseProfile["gioiTinh"] = customer.GioiTinh
			baseProfile["diaChi"] = customer.DiaChi
			baseProfile["maBaoHiem"] = customer.MaBaoHiem
		}
	case "DOCTOR":
		doctor, err := h.mockService.GetDoctorByUserID(user.MaUser)
		if err == nil {
			baseProfile["chuyenKhoa"] = doctor.ChuyenKhoa
			baseProfile["namKinhNghiem"] = doctor.NamKinhNghiem
			baseProfile["bangCap"] = doctor.BangCap
			baseProfile["soGiayPhepHanhNghe"] = doctor.SoGiayPhepHanhNghe
		}
	}

	profileData = baseProfile
	utils.SuccessResponse(c, "Profile retrieved successfully", profileData)
}

func (h *MockUserHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User ID not found in token", "")
		return
	}

	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request format", err.Error())
		return
	}

	// Check if user exists
	user, err := h.mockService.GetUserByID(userID.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "User not found", err.Error())
		return
	}

	// Mock implementation - in real app you would update the database
	updatedProfile := map[string]interface{}{
		"maUser":  user.MaUser,
		"message": "Profile updated successfully (mock implementation)",
		"updated": updateData,
	}

	utils.SuccessResponse(c, "Profile updated successfully", updatedProfile)
}

func (h *MockUserHandler) ChangePassword(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User ID not found in token", "")
		return
	}

	type ChangePasswordRequest struct {
		CurrentPassword string `json:"currentPassword" binding:"required"`
		NewPassword     string `json:"newPassword" binding:"required,min=6"`
	}

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request format", err.Error())
		return
	}

	// Check if user exists
	user, err := h.mockService.GetUserByID(userID.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "User not found", err.Error())
		return
	}

	// Mock password verification - in real app you'd use bcrypt.CompareHashAndPassword
	if !checkMockPassword(req.CurrentPassword, user.MatKhau) {
		utils.ErrorResponse(c, http.StatusBadRequest, "Current password is incorrect", "")
		return
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to hash password", err.Error())
		return
	}

	// Mock implementation - in real app you would update the database
	response := map[string]interface{}{
		"message":       "Password changed successfully (mock implementation)",
		"hashedPassword": string(hashedPassword),
	}

	utils.SuccessResponse(c, "Password changed successfully", response)
}

func checkMockPassword(plaintext, hashed string) bool {
	// Mock implementation - in real app you'd use bcrypt.CompareHashAndPassword
	// For mock purposes, accept if the hashed password contains "hashed.password"
	return len(plaintext) >= 6 && len(hashed) > 0
}