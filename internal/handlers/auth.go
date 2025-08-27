package handlers

import (
	"database/sql"
	"net/http"

	"clinic-management/internal/middleware"
	"clinic-management/internal/models"
	"clinic-management/internal/utils"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	db        *sql.DB
	jwtSecret string
}

func NewAuthHandler(db *sql.DB, jwtSecret string) *AuthHandler {
	return &AuthHandler{db: db, jwtSecret: jwtSecret}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req models.AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request format",
			Error:   err.Error(),
		})
		return
	}

	if req.TenDangNhap == "" || req.MatKhau == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Username and password are required",
		})
		return
	}

	var user models.User
	query := `SELECT userID, HoTen, SoDienThoai, Email, username, password, status, createdAt, role 
			  FROM [USER] WHERE username = @p1 AND status = 'ACTIVE'`

	err := h.db.QueryRow(query, req.TenDangNhap).Scan(
		&user.MaUser, &user.HoTen, &user.SoDienThoai, &user.Email,
		&user.TenDangNhap, &user.MatKhau, &user.TrangThai,
		&user.NgayTao, &user.Role,
	)

	if err != nil {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Message: "Invalid credentials",
		})
		return
	}

	if !utils.CheckPasswordHash(req.MatKhau, user.MatKhau) {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Message: "Invalid credentials",
		})
		return
	}

	token, err := middleware.GenerateToken(user.MaUser, user.TenDangNhap, user.Role, h.jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to generate token",
			Error:   err.Error(),
		})
		return
	}

	refreshToken, err := middleware.GenerateRefreshToken(user.MaUser, h.jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to generate refresh token",
			Error:   err.Error(),
		})
		return
	}

	user.MatKhau = ""

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Login successful",
		Data: models.AuthResponse{
			Token:        token,
			RefreshToken: refreshToken,
			User:         user,
		},
	})
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request format",
			Error:   err.Error(),
		})
		return
	}

	if req.HoTen == "" || req.TenDangNhap == "" || req.MatKhau == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Name, username and password are required",
		})
		return
	}

	if req.Email != "" && !utils.ValidateEmail(req.Email) {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid email format",
		})
		return
	}

	var existingUser string
	err := h.db.QueryRow("SELECT userID FROM [USER] WHERE username = @p1", req.TenDangNhap).Scan(&existingUser)
	if err == nil {
		c.JSON(http.StatusConflict, models.APIResponse{
			Success: false,
			Message: "Username already exists",
		})
		return
	}

	if req.Email != "" {
		err = h.db.QueryRow("SELECT userID FROM [USER] WHERE Email = @p1", req.Email).Scan(&existingUser)
		if err == nil {
			c.JSON(http.StatusConflict, models.APIResponse{
				Success: false,
				Message: "Email already exists",
			})
			return
		}
	}

	hashedPassword, err := utils.HashPassword(req.MatKhau)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to process password",
			Error:   err.Error(),
		})
		return
	}

	userID := utils.GenerateUserID("CUSTOMER")

	tx, err := h.db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to begin transaction",
			Error:   err.Error(),
		})
		return
	}
	defer tx.Rollback()

	var email, phone interface{}
	if req.Email != "" {
		email = req.Email
	}
	if req.SoDienThoai != "" {
		phone = req.SoDienThoai
	}

	_, err = tx.Exec(`
		INSERT INTO [USER] (userID, HoTen, SoDienThoai, Email, username, password, status, createdAt, role)
		VALUES (@p1, @p2, @p3, @p4, @p5, @p6, 'ACTIVE', GETDATE(), 'CUSTOMER')
	`, userID, req.HoTen, phone, email, req.TenDangNhap, hashedPassword)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to create user account",
			Error:   err.Error(),
		})
		return
	}

	_, err = tx.Exec(`
		INSERT INTO CUSTOMER (MaUser, createdAt)
		VALUES (@p1, GETDATE())
	`, userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to create customer profile",
			Error:   err.Error(),
		})
		return
	}

	if err = tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to complete registration",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "Registration successful",
		Data: gin.H{
			"user_id":  userID,
			"username": req.TenDangNhap,
			"name":     req.HoTen,
		},
	})
}

func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req models.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request format",
			Error:   err.Error(),
		})
		return
	}

	if req.Email == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Email is required",
		})
		return
	}

	if !utils.ValidateEmail(req.Email) {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid email format",
		})
		return
	}

	var user models.User
	err := h.db.QueryRow("SELECT MaUser, HoTen, Email, TrangThai FROM [USER] WHERE Email = ? AND TrangThai = 'ACTIVE'", req.Email).Scan(
		&user.MaUser, &user.HoTen, &user.Email, &user.TrangThai,
	)

	if err != nil {
		c.JSON(http.StatusOK, models.APIResponse{
			Success: true,
			Message: "If the email exists in our system, a password reset code will be sent",
		})
		return
	}

	resetCode := utils.GenerateResetCode()
	resetID := utils.GeneratePasswordResetID()

	tx, err := h.db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to process request",
			Error:   err.Error(),
		})
		return
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
		UPDATE PASSWORD_RESET SET IsUsed = 1 
		WHERE UserID = ? AND IsUsed = 0 AND ExpiresAt > GETDATE()
	`, user.MaUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to invalidate previous reset codes",
			Error:   err.Error(),
		})
		return
	}

	_, err = tx.Exec(`
		INSERT INTO PASSWORD_RESET (ID, UserID, Email, ResetCode, IsUsed, ExpiresAt, CreatedAt)
		VALUES (?, ?, ?, ?, 0, DATEADD(HOUR, 1, GETDATE()), GETDATE())
	`, resetID, user.MaUser, req.Email, resetCode)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to create password reset request",
			Error:   err.Error(),
		})
		return
	}

	if err = tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to process password reset request",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Password reset code has been sent to your email",
		Data: gin.H{
			"reset_code": resetCode,
			"expires_in": "1 hour",
		},
	})
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req models.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request format",
			Error:   err.Error(),
		})
		return
	}

	if req.Email == "" || req.ResetCode == "" || req.NewPassword == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Email, reset code, and new password are required",
		})
		return
	}

	if !utils.ValidateEmail(req.Email) {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid email format",
		})
		return
	}

	if len(req.NewPassword) < 6 {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "New password must be at least 6 characters long",
		})
		return
	}

	var resetRecord models.PasswordReset
	err := h.db.QueryRow(`
		SELECT pr.ID, pr.UserID, pr.Email, pr.ResetCode, pr.IsUsed, pr.ExpiresAt, pr.CreatedAt
		FROM PASSWORD_RESET pr
		JOIN [USER] u ON pr.UserID = u.MaUser
		WHERE pr.Email = ? AND pr.ResetCode = ? AND pr.IsUsed = 0 
		AND pr.ExpiresAt > GETDATE() AND u.TrangThai = 'ACTIVE'
	`, req.Email, req.ResetCode).Scan(
		&resetRecord.ID, &resetRecord.UserID, &resetRecord.Email,
		&resetRecord.ResetCode, &resetRecord.IsUsed, &resetRecord.ExpiresAt,
		&resetRecord.CreatedAt,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid or expired reset code",
		})
		return
	}

	newHash, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to process new password",
			Error:   err.Error(),
		})
		return
	}

	tx, err := h.db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to process request",
			Error:   err.Error(),
		})
		return
	}
	defer tx.Rollback()

	_, err = tx.Exec("UPDATE [USER] SET MatKhau = ? WHERE MaUser = ?", newHash, resetRecord.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to update password",
			Error:   err.Error(),
		})
		return
	}

	_, err = tx.Exec("UPDATE PASSWORD_RESET SET IsUsed = 1 WHERE ID = ?", resetRecord.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to mark reset code as used",
			Error:   err.Error(),
		})
		return
	}

	if err = tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to reset password",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Password has been reset successfully",
	})
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, models.APIResponse{
		Success: false,
		Message: "Refresh token not implemented yet",
	})
}

func (h *AuthHandler) getUserType(userID string) string {
	var exists bool

	h.db.QueryRow("SELECT 1 FROM CUSTOMER WHERE MaUser = ?", userID).Scan(&exists)
	if exists {
		return "CUSTOMER"
	}

	h.db.QueryRow("SELECT 1 FROM BACSI WHERE MaUser = ?", userID).Scan(&exists)
	if exists {
		return "DOCTOR"
	}

	h.db.QueryRow("SELECT 1 FROM LETAN WHERE MaUser = ?", userID).Scan(&exists)
	if exists {
		return "RECEPTIONIST"
	}

	h.db.QueryRow("SELECT 1 FROM KETOAN WHERE MaUser = ?", userID).Scan(&exists)
	if exists {
		return "ACCOUNTANT"
	}

	h.db.QueryRow("SELECT 1 FROM QUANLYPHONGKHAM WHERE MaUser = ?", userID).Scan(&exists)
	if exists {
		return "CLINIC_MANAGER"
	}

	h.db.QueryRow("SELECT 1 FROM BANDIEUHANH WHERE MaUser = ?", userID).Scan(&exists)
	if exists {
		return "OPERATION_MANAGER"
	}

	return "UNKNOWN"
}
