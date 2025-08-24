package handlers

import (
	"database/sql"
	"net/http"

	"clinic-management/internal/models"
	"clinic-management/internal/utils"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	db *sql.DB
}

func NewUserHandler(db *sql.DB) *UserHandler {
	return &UserHandler{db: db}
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	userType, _ := c.Get("user_type")

	switch userType.(string) {
	case "CUSTOMER":
		h.getCustomerProfile(c, userID.(string))
	case "DOCTOR":
		h.getDoctorProfile(c, userID.(string))
	case "RECEPTIONIST":
		h.getReceptionistProfile(c, userID.(string))
	case "ACCOUNTANT":
		h.getAccountantProfile(c, userID.(string))
	case "CLINIC_MANAGER":
		h.getClinicManagerProfile(c, userID.(string))
	case "OPERATION_MANAGER":
		h.getOperationManagerProfile(c, userID.(string))
	default:
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Unknown user type",
		})
	}
}

func (h *UserHandler) getCustomerProfile(c *gin.Context, userID string) {
	var customer models.Customer
	query := `
		SELECT u.MaUser, u.HoTen, u.SoDienThoai, u.Email, u.TenDangNhap, u.TrangThai, u.NgayTao, u.LanDangNhapCuoi,
		       c.NgaySinh, c.GioiTinh, c.DiaChi, c.NgayDangKy, c.MaBaoHiem
		FROM [USER] u 
		JOIN CUSTOMER c ON u.MaUser = c.MaUser 
		WHERE u.MaUser = ?
	`

	err := h.db.QueryRow(query, userID).Scan(
		&customer.MaUser, &customer.HoTen, &customer.SoDienThoai, &customer.Email,
		&customer.TenDangNhap, &customer.TrangThai, &customer.NgayTao, &customer.LanDangNhapCuoi,
		&customer.NgaySinh, &customer.GioiTinh, &customer.DiaChi, &customer.NgayDangKy, &customer.MaBaoHiem,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: "User not found",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Profile retrieved successfully",
		Data:    customer,
	})
}

func (h *UserHandler) getDoctorProfile(c *gin.Context, userID string) {
	var doctor models.Doctor
	query := `
		SELECT u.MaUser, u.HoTen, u.SoDienThoai, u.Email, u.TenDangNhap, u.TrangThai, u.NgayTao, u.LanDangNhapCuoi,
		       d.ChuyenKhoa, d.NamKinhNghiem, d.BangCap, d.SoGiayPhepHanhNghe
		FROM [USER] u 
		JOIN BACSI d ON u.MaUser = d.MaUser 
		WHERE u.MaUser = ?
	`

	err := h.db.QueryRow(query, userID).Scan(
		&doctor.MaUser, &doctor.HoTen, &doctor.SoDienThoai, &doctor.Email,
		&doctor.TenDangNhap, &doctor.TrangThai, &doctor.NgayTao, &doctor.LanDangNhapCuoi,
		&doctor.ChuyenKhoa, &doctor.NamKinhNghiem, &doctor.BangCap, &doctor.SoGiayPhepHanhNghe,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: "Doctor not found",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Profile retrieved successfully",
		Data:    doctor,
	})
}

func (h *UserHandler) getReceptionistProfile(c *gin.Context, userID string) {
	var receptionist models.Receptionist
	query := `
		SELECT u.MaUser, u.HoTen, u.SoDienThoai, u.Email, u.TenDangNhap, u.TrangThai, u.NgayTao, u.LanDangNhapCuoi,
		       r.MaPhongKham, r.LuongCoBan, r.NgayVaoLam
		FROM [USER] u 
		JOIN LETAN r ON u.MaUser = r.MaUser 
		WHERE u.MaUser = ?
	`

	err := h.db.QueryRow(query, userID).Scan(
		&receptionist.MaUser, &receptionist.HoTen, &receptionist.SoDienThoai, &receptionist.Email,
		&receptionist.TenDangNhap, &receptionist.TrangThai, &receptionist.NgayTao, &receptionist.LanDangNhapCuoi,
		&receptionist.MaPhongKham, &receptionist.LuongCoBan, &receptionist.NgayVaoLam,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: "Receptionist not found",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Profile retrieved successfully",
		Data:    receptionist,
	})
}

func (h *UserHandler) getAccountantProfile(c *gin.Context, userID string) {
	var accountant models.Accountant
	query := `
		SELECT u.MaUser, u.HoTen, u.SoDienThoai, u.Email, u.TenDangNhap, u.TrangThai, u.NgayTao, u.LanDangNhapCuoi,
		       a.LuongCoBan, a.NgayVaoLam, a.ChuyenMon
		FROM [USER] u 
		JOIN KETOAN a ON u.MaUser = a.MaUser 
		WHERE u.MaUser = ?
	`

	err := h.db.QueryRow(query, userID).Scan(
		&accountant.MaUser, &accountant.HoTen, &accountant.SoDienThoai, &accountant.Email,
		&accountant.TenDangNhap, &accountant.TrangThai, &accountant.NgayTao, &accountant.LanDangNhapCuoi,
		&accountant.LuongCoBan, &accountant.NgayVaoLam, &accountant.ChuyenMon,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: "Accountant not found",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Profile retrieved successfully",
		Data:    accountant,
	})
}

func (h *UserHandler) getClinicManagerProfile(c *gin.Context, userID string) {
	var manager models.ClinicManager
	query := `
		SELECT u.MaUser, u.HoTen, u.SoDienThoai, u.Email, u.TenDangNhap, u.TrangThai, u.NgayTao, u.LanDangNhapCuoi,
		       m.MaPhongKham, m.LuongCoBan, m.NgayVaoLam
		FROM [USER] u 
		JOIN QUANLYPHONGKHAM m ON u.MaUser = m.MaUser 
		WHERE u.MaUser = ?
	`

	err := h.db.QueryRow(query, userID).Scan(
		&manager.MaUser, &manager.HoTen, &manager.SoDienThoai, &manager.Email,
		&manager.TenDangNhap, &manager.TrangThai, &manager.NgayTao, &manager.LanDangNhapCuoi,
		&manager.MaPhongKham, &manager.LuongCoBan, &manager.NgayVaoLam,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: "Clinic manager not found",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Profile retrieved successfully",
		Data:    manager,
	})
}

func (h *UserHandler) getOperationManagerProfile(c *gin.Context, userID string) {
	var opManager models.OperationManager
	query := `
		SELECT u.MaUser, u.HoTen, u.SoDienThoai, u.Email, u.TenDangNhap, u.TrangThai, u.NgayTao, u.LanDangNhapCuoi,
		       o.ChucVu, o.KhuVucPhuTrach, o.LuongCoBan, o.NgayVaoLam
		FROM [USER] u 
		JOIN BANDIEUHANH o ON u.MaUser = o.MaUser 
		WHERE u.MaUser = ?
	`

	err := h.db.QueryRow(query, userID).Scan(
		&opManager.MaUser, &opManager.HoTen, &opManager.SoDienThoai, &opManager.Email,
		&opManager.TenDangNhap, &opManager.TrangThai, &opManager.NgayTao, &opManager.LanDangNhapCuoi,
		&opManager.ChucVu, &opManager.KhuVucPhuTrach, &opManager.LuongCoBan, &opManager.NgayVaoLam,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: "Operation manager not found",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Profile retrieved successfully",
		Data:    opManager,
	})
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	userType, _ := c.Get("user_type")

	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request format",
			Error:   err.Error(),
		})
		return
	}

	if len(updateData) == 0 {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "No data provided for update",
		})
		return
	}

	if email, exists := updateData["email"]; exists && email != "" {
		emailStr, ok := email.(string)
		if !ok || !utils.ValidateEmail(emailStr) {
			c.JSON(http.StatusBadRequest, models.APIResponse{
				Success: false,
				Message: "Invalid email format",
			})
			return
		}

		var existingUserID string
		err := h.db.QueryRow("SELECT MaUser FROM [USER] WHERE Email = ? AND MaUser != ?", emailStr, userID).Scan(&existingUserID)
		if err == nil {
			c.JSON(http.StatusConflict, models.APIResponse{
				Success: false,
				Message: "Email already exists",
			})
			return
		}
	}

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

	if hoTen, exists := updateData["ho_ten"]; exists && hoTen != "" {
		_, err = tx.Exec("UPDATE [USER] SET HoTen = ? WHERE MaUser = ?", hoTen, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "Failed to update name",
				Error:   err.Error(),
			})
			return
		}
	}

	if soDienThoai, exists := updateData["so_dien_thoai"]; exists {
		var phoneValue interface{}
		if soDienThoai == "" {
			phoneValue = nil
		} else {
			phoneValue = soDienThoai
		}
		_, err = tx.Exec("UPDATE [USER] SET SoDienThoai = ? WHERE MaUser = ?", phoneValue, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "Failed to update phone",
				Error:   err.Error(),
			})
			return
		}
	}

	if email, exists := updateData["email"]; exists {
		var emailValue interface{}
		if email == "" {
			emailValue = nil
		} else {
			emailValue = email
		}
		_, err = tx.Exec("UPDATE [USER] SET Email = ? WHERE MaUser = ?", emailValue, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "Failed to update email",
				Error:   err.Error(),
			})
			return
		}
	}

	switch userType.(string) {
	case "CUSTOMER":
		err = h.updateCustomerSpecificFields(tx, userID.(string), updateData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "Failed to update customer specific fields",
				Error:   err.Error(),
			})
			return
		}
	}

	if err = tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to update profile",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Profile updated successfully",
	})
}

func (h *UserHandler) updateCustomerSpecificFields(tx *sql.Tx, userID string, updateData map[string]interface{}) error {
	if ngaySinh, exists := updateData["ngay_sinh"]; exists {
		var dateValue interface{}
		if ngaySinh == "" {
			dateValue = nil
		} else {
			dateValue = ngaySinh
		}
		_, err := tx.Exec("UPDATE CUSTOMER SET NgaySinh = ? WHERE MaUser = ?", dateValue, userID)
		if err != nil {
			return err
		}
	}

	if gioiTinh, exists := updateData["gioi_tinh"]; exists {
		var genderValue interface{}
		if gioiTinh == "" {
			genderValue = nil
		} else {
			genderValue = gioiTinh
		}
		_, err := tx.Exec("UPDATE CUSTOMER SET GioiTinh = ? WHERE MaUser = ?", genderValue, userID)
		if err != nil {
			return err
		}
	}

	if diaChi, exists := updateData["dia_chi"]; exists {
		var addressValue interface{}
		if diaChi == "" {
			addressValue = nil
		} else {
			addressValue = diaChi
		}
		_, err := tx.Exec("UPDATE CUSTOMER SET DiaChi = ? WHERE MaUser = ?", addressValue, userID)
		if err != nil {
			return err
		}
	}

	if maBaoHiem, exists := updateData["ma_bao_hiem"]; exists {
		var insuranceValue interface{}
		if maBaoHiem == "" {
			insuranceValue = nil
		} else {
			insuranceValue = maBaoHiem
		}
		_, err := tx.Exec("UPDATE CUSTOMER SET MaBaoHiem = ? WHERE MaUser = ?", insuranceValue, userID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID, _ := c.Get("user_id")
	
	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request format",
			Error:   err.Error(),
		})
		return
	}

	if req.OldPassword == "" || req.NewPassword == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Old password and new password are required",
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

	if req.OldPassword == req.NewPassword {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "New password must be different from current password",
		})
		return
	}

	var currentHash string
	var userStatus string
	err := h.db.QueryRow("SELECT MatKhau, TrangThai FROM [USER] WHERE MaUser = ?", userID).Scan(&currentHash, &userStatus)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: "User not found",
		})
		return
	}

	if userStatus != "ACTIVE" {
		c.JSON(http.StatusForbidden, models.APIResponse{
			Success: false,
			Message: "Account is not active",
		})
		return
	}

	if !utils.CheckPasswordHash(req.OldPassword, currentHash) {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Current password is incorrect",
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

	_, err = h.db.Exec("UPDATE [USER] SET MatKhau = ? WHERE MaUser = ?", newHash, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to update password",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Password changed successfully",
	})
}