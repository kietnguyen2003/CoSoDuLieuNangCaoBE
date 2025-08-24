package handlers

import (
	"net/http"
	"strings"
	"time"

	"clinic-management/internal/services"
	"clinic-management/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type MockAuthHandler struct {
	mockService *services.MockDataService
}

func NewMockAuthHandler(mockService *services.MockDataService) *MockAuthHandler {
	return &MockAuthHandler{
		mockService: mockService,
	}
}

type LoginRequest struct {
	TenDangNhap string `json:"tenDangNhap" binding:"required"`
	MatKhau     string `json:"matKhau" binding:"required"`
}

type RegisterRequest struct {
	HoTen       string `json:"hoTen" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	SoDienThoai string `json:"soDienThoai"`
	TenDangNhap string `json:"tenDangNhap" binding:"required"`
	MatKhau     string `json:"matKhau" binding:"required,min=6"`
}

type LoginResponse struct {
	Token     string                 `json:"token"`
	User      map[string]interface{} `json:"user"`
	ExpiresAt time.Time              `json:"expiresAt"`
}

func (h *MockAuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request format", err.Error())
		return
	}

	// Find user by username
	user, err := h.mockService.GetUserByUsername(req.TenDangNhap)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid credentials", "User not found")
		return
	}

	// For mock purposes, we'll accept any password that matches the pattern
	// In real implementation, you would check bcrypt.CompareHashAndPassword
	if !strings.Contains(user.MatKhau, "hashed.password") {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid credentials", "Invalid password")
		return
	}

	// Generate JWT token
	token, expiresAt, err := generateToken(user.MaUser, user.Role)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate token", err.Error())
		return
	}

	// Prepare user response
	userResponse := map[string]interface{}{
		"maUser":      user.MaUser,
		"hoTen":       user.HoTen,
		"email":       user.Email,
		"soDienThoai": user.SoDienThoai,
		"role":        user.Role,
		"trangThai":   user.TrangThai,
	}

	response := LoginResponse{
		Token:     token,
		User:      userResponse,
		ExpiresAt: expiresAt,
	}

	utils.SuccessResponse(c, "Login successful", response)
}

func (h *MockAuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request format", err.Error())
		return
	}

	// Check if username already exists
	_, err := h.mockService.GetUserByUsername(req.TenDangNhap)
	if err == nil {
		utils.ErrorResponse(c, http.StatusConflict, "Username already exists", "")
		return
	}

	// For mock purposes, we'll return success without actually creating the user
	// In real implementation, you would hash password and save to database
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.MatKhau), bcrypt.DefaultCost)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to hash password", err.Error())
		return
	}

	// Mock response - in real app you'd generate a proper ID and save to DB
	userResponse := map[string]interface{}{
		"maUser":      "USER_NEW_" + req.TenDangNhap,
		"hoTen":       req.HoTen,
		"email":       req.Email,
		"soDienThoai": req.SoDienThoai,
		"role":        "CUSTOMER",
		"trangThai":   "ACTIVE",
		"message":     "User registered successfully (mock implementation)",
		"hashedPassword": string(hashedPassword),
	}

	utils.SuccessResponse(c, "Registration successful", userResponse)
}

func (h *MockAuthHandler) ForgotPassword(c *gin.Context) {
	type ForgotPasswordRequest struct {
		Email string `json:"email" binding:"required,email"`
	}

	var req ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request format", err.Error())
		return
	}

	// Mock implementation - in real app you'd send email
	utils.SuccessResponse(c, "Password reset email sent (mock)", map[string]interface{}{
		"email":   req.Email,
		"message": "If this email exists in our system, you will receive a password reset email",
	})
}

func (h *MockAuthHandler) ResetPassword(c *gin.Context) {
	type ResetPasswordRequest struct {
		Email    string `json:"email" binding:"required,email"`
		Code     string `json:"code" binding:"required"`
		Password string `json:"password" binding:"required,min=6"`
	}

	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request format", err.Error())
		return
	}

	// Mock implementation
	utils.SuccessResponse(c, "Password reset successful (mock)", map[string]interface{}{
		"email":   req.Email,
		"message": "Password has been reset successfully",
	})
}

func (h *MockAuthHandler) RefreshToken(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Authorization header required", "")
		return
	}

	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
	
	// Parse and validate existing token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("your-secret-key"), nil // Use your actual secret
	})

	if err != nil || !token.Valid {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid token", err.Error())
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		userID := claims["userID"].(string)
		role := claims["role"].(string)

		// Generate new token
		newToken, expiresAt, err := generateToken(userID, role)
		if err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate token", err.Error())
			return
		}

		response := map[string]interface{}{
			"token":     newToken,
			"expiresAt": expiresAt,
		}

		utils.SuccessResponse(c, "Token refreshed successfully", response)
	} else {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid token claims", "")
	}
}

func generateToken(userID, role string) (string, time.Time, error) {
	expiresAt := time.Now().Add(24 * time.Hour)
	
	claims := jwt.MapClaims{
		"userID": userID,
		"role":   role,
		"exp":    expiresAt.Unix(),
		"iat":    time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("your-secret-key"))
	
	return tokenString, expiresAt, err
}