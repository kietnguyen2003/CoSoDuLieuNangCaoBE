package handlers

import (
	"database/sql"
	"net/http"

	"clinic-management/internal/models"
	"clinic-management/internal/utils"

	"github.com/gin-gonic/gin"
)

type AppointmentHandler struct {
	db *sql.DB
}

func NewAppointmentHandler(db *sql.DB) *AppointmentHandler {
	return &AppointmentHandler{db: db}
}

func (h *AppointmentHandler) GetAppointments(c *gin.Context) {
	userID, _ := c.Get("user_id")
	userType, _ := c.Get("user_type")
	status := c.Query("status")
	date := c.Query("date")

	var query string
	var args []interface{}

	switch userType.(string) {
	case "CUSTOMER":
		query = `
			SELECT l.MaLichKham, l.MaCustomer, l.MaBacSi, l.MaPhongKham, 
			       l.NgayGioKham, l.TrangThai, l.GhiChu, l.NgayDat,
			       u.HoTen as TenBacSi, p.TenPhongKham
			FROM LICHKHAM l
			JOIN [USER] u ON l.MaBacSi = u.MaUser
			JOIN PHONGKHAM p ON l.MaPhongKham = p.MaPhongKham
			WHERE l.MaCustomer = ?
		`
		args = append(args, userID)

	case "DOCTOR":
		query = `
			SELECT l.MaLichKham, l.MaCustomer, l.MaBacSi, l.MaPhongKham, 
			       l.NgayGioKham, l.TrangThai, l.GhiChu, l.NgayDat,
			       u.HoTen as TenKhachHang, p.TenPhongKham
			FROM LICHKHAM l
			JOIN [USER] u ON l.MaCustomer = u.MaUser
			JOIN PHONGKHAM p ON l.MaPhongKham = p.MaPhongKham
			WHERE l.MaBacSi = ?
		`
		args = append(args, userID)

	default:
		query = `
			SELECT l.MaLichKham, l.MaCustomer, l.MaBacSi, l.MaPhongKham, 
			       l.NgayGioKham, l.TrangThai, l.GhiChu, l.NgayDat,
			       uc.HoTen as TenKhachHang, ud.HoTen as TenBacSi, p.TenPhongKham
			FROM LICHKHAM l
			JOIN [USER] uc ON l.MaCustomer = uc.MaUser
			JOIN [USER] ud ON l.MaBacSi = ud.MaUser
			JOIN PHONGKHAM p ON l.MaPhongKham = p.MaPhongKham
			WHERE 1=1
		`
	}

	if status != "" {
		query += " AND l.TrangThai = ?"
		args = append(args, status)
	}

	if date != "" {
		query += " AND CAST(l.NgayGioKham AS DATE) = ?"
		args = append(args, date)
	}

	query += " ORDER BY l.NgayGioKham DESC"

	rows, err := h.db.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to retrieve appointments",
			Error:   err.Error(),
		})
		return
	}
	defer rows.Close()

	var appointments []map[string]interface{}
	for rows.Next() {
		appointment := make(map[string]interface{})
		var maLichKham, maCustomer, maBacSi, maPhongKham, trangThai string
		var ngayGioKham, ngayDat interface{}
		var ghiChu, tenPerson, tenPhongKham interface{}

		if userType.(string) == "CUSTOMER" || userType.(string) == "DOCTOR" {
			err := rows.Scan(&maLichKham, &maCustomer, &maBacSi, &maPhongKham,
				&ngayGioKham, &trangThai, &ghiChu, &ngayDat, &tenPerson, &tenPhongKham)
			if err != nil {
				c.JSON(http.StatusInternalServerError, models.APIResponse{
					Success: false,
					Message: "Failed to scan appointment data",
					Error:   err.Error(),
				})
				return
			}
		} else {
			var tenKhachHang, tenBacSi interface{}
			err := rows.Scan(&maLichKham, &maCustomer, &maBacSi, &maPhongKham,
				&ngayGioKham, &trangThai, &ghiChu, &ngayDat, &tenKhachHang, &tenBacSi, &tenPhongKham)
			if err != nil {
				c.JSON(http.StatusInternalServerError, models.APIResponse{
					Success: false,
					Message: "Failed to scan appointment data",
					Error:   err.Error(),
				})
				return
			}
			appointment["ten_khach_hang"] = tenKhachHang
			appointment["ten_bac_si"] = tenBacSi
		}

		appointment["ma_lich_kham"] = maLichKham
		appointment["ma_customer"] = maCustomer
		appointment["ma_bac_si"] = maBacSi
		appointment["ma_phong_kham"] = maPhongKham
		appointment["ngay_gio_kham"] = ngayGioKham
		appointment["trang_thai"] = trangThai
		appointment["ghi_chu"] = ghiChu
		appointment["ngay_dat"] = ngayDat
		appointment["ten_phong_kham"] = tenPhongKham

		if userType.(string) == "CUSTOMER" {
			appointment["ten_bac_si"] = tenPerson
		} else if userType.(string) == "DOCTOR" {
			appointment["ten_khach_hang"] = tenPerson
		}

		appointments = append(appointments, appointment)
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Appointments retrieved successfully",
		Data:    appointments,
	})
}

func (h *AppointmentHandler) CreateAppointment(c *gin.Context) {
	userID, _ := c.Get("user_id")
	userType, _ := c.Get("user_type")

	var req struct {
		MaBacSi     string `json:"ma_bac_si" binding:"required"`
		MaPhongKham string `json:"ma_phong_kham" binding:"required"`
		NgayGioKham string `json:"ngay_gio_kham" binding:"required"`
		GhiChu      string `json:"ghi_chu"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request format",
			Error:   err.Error(),
		})
		return
	}

	var customerID string
	if userType.(string) == "CUSTOMER" {
		customerID = userID.(string)
	} else {
		c.JSON(http.StatusForbidden, models.APIResponse{
			Success: false,
			Message: "Only customers can book appointments",
		})
		return
	}

	var count int
	err := h.db.QueryRow(`
		SELECT COUNT(*) FROM LICHKHAM 
		WHERE MaBacSi = ? AND NgayGioKham = ? AND TrangThai NOT IN ('HUY_LICH', 'HOAN_THANH')
	`, req.MaBacSi, req.NgayGioKham).Scan(&count)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to check availability",
			Error:   err.Error(),
		})
		return
	}

	if count > 0 {
		c.JSON(http.StatusConflict, models.APIResponse{
			Success: false,
			Message: "Time slot is not available",
		})
		return
	}

	appointmentID := utils.GenerateAppointmentID()
	
	_, err = h.db.Exec(`
		INSERT INTO LICHKHAM (MaLichKham, MaCustomer, MaBacSi, MaPhongKham, NgayGioKham, TrangThai, GhiChu, NgayDat)
		VALUES (?, ?, ?, ?, ?, 'DAT_LICH', ?, GETDATE())
	`, appointmentID, customerID, req.MaBacSi, req.MaPhongKham, req.NgayGioKham, req.GhiChu)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to create appointment",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "Appointment created successfully",
		Data: gin.H{
			"appointment_id": appointmentID,
		},
	})
}

func (h *AppointmentHandler) GetAppointment(c *gin.Context) {
	appointmentID := c.Param("id")
	userID, _ := c.Get("user_id")
	userType, _ := c.Get("user_type")

	query := `
		SELECT l.MaLichKham, l.MaCustomer, l.MaBacSi, l.MaPhongKham, 
		       l.NgayGioKham, l.TrangThai, l.GhiChu, l.NgayDat,
		       uc.HoTen as TenKhachHang, ud.HoTen as TenBacSi, p.TenPhongKham
		FROM LICHKHAM l
		JOIN [USER] uc ON l.MaCustomer = uc.MaUser
		JOIN [USER] ud ON l.MaBacSi = ud.MaUser
		JOIN PHONGKHAM p ON l.MaPhongKham = p.MaPhongKham
		WHERE l.MaLichKham = ?
	`
	args := []interface{}{appointmentID}

	if userType.(string) == "CUSTOMER" {
		query += " AND l.MaCustomer = ?"
		args = append(args, userID)
	} else if userType.(string) == "DOCTOR" {
		query += " AND l.MaBacSi = ?"
		args = append(args, userID)
	}

	var appointment map[string]interface{} = make(map[string]interface{})
	var maLichKham, maCustomer, maBacSi, maPhongKham, trangThai string
	var ngayGioKham, ngayDat interface{}
	var ghiChu, tenKhachHang, tenBacSi, tenPhongKham interface{}

	err := h.db.QueryRow(query, args...).Scan(
		&maLichKham, &maCustomer, &maBacSi, &maPhongKham,
		&ngayGioKham, &trangThai, &ghiChu, &ngayDat,
		&tenKhachHang, &tenBacSi, &tenPhongKham,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, models.APIResponse{
				Success: false,
				Message: "Appointment not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "Failed to retrieve appointment",
				Error:   err.Error(),
			})
		}
		return
	}

	appointment["ma_lich_kham"] = maLichKham
	appointment["ma_customer"] = maCustomer
	appointment["ma_bac_si"] = maBacSi
	appointment["ma_phong_kham"] = maPhongKham
	appointment["ngay_gio_kham"] = ngayGioKham
	appointment["trang_thai"] = trangThai
	appointment["ghi_chu"] = ghiChu
	appointment["ngay_dat"] = ngayDat
	appointment["ten_khach_hang"] = tenKhachHang
	appointment["ten_bac_si"] = tenBacSi
	appointment["ten_phong_kham"] = tenPhongKham

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Appointment retrieved successfully",
		Data:    appointment,
	})
}

func (h *AppointmentHandler) UpdateAppointment(c *gin.Context) {
	appointmentID := c.Param("id")
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

	var currentCustomerID, currentDoctorID string
	err := h.db.QueryRow("SELECT MaCustomer, MaBacSi FROM LICHKHAM WHERE MaLichKham = ?", appointmentID).
		Scan(&currentCustomerID, &currentDoctorID)
	
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, models.APIResponse{
				Success: false,
				Message: "Appointment not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "Failed to find appointment",
				Error:   err.Error(),
			})
		}
		return
	}

	if userType.(string) == "CUSTOMER" && currentCustomerID != userID.(string) {
		c.JSON(http.StatusForbidden, models.APIResponse{
			Success: false,
			Message: "You can only update your own appointments",
		})
		return
	}

	if userType.(string) == "DOCTOR" && currentDoctorID != userID.(string) {
		c.JSON(http.StatusForbidden, models.APIResponse{
			Success: false,
			Message: "You can only update your own appointments",
		})
		return
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

	if ngayGioKham, exists := updateData["ngay_gio_kham"]; exists {
		_, err = tx.Exec("UPDATE LICHKHAM SET NgayGioKham = ? WHERE MaLichKham = ?", ngayGioKham, appointmentID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "Failed to update appointment time",
				Error:   err.Error(),
			})
			return
		}
	}

	if trangThai, exists := updateData["trang_thai"]; exists {
		_, err = tx.Exec("UPDATE LICHKHAM SET TrangThai = ? WHERE MaLichKham = ?", trangThai, appointmentID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "Failed to update appointment status",
				Error:   err.Error(),
			})
			return
		}
	}

	if ghiChu, exists := updateData["ghi_chu"]; exists {
		_, err = tx.Exec("UPDATE LICHKHAM SET GhiChu = ? WHERE MaLichKham = ?", ghiChu, appointmentID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "Failed to update appointment note",
				Error:   err.Error(),
			})
			return
		}
	}

	if err = tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to update appointment",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Appointment updated successfully",
	})
}

func (h *AppointmentHandler) CancelAppointment(c *gin.Context) {
	appointmentID := c.Param("id")
	userID, _ := c.Get("user_id")
	userType, _ := c.Get("user_type")

	var customerID string
	err := h.db.QueryRow("SELECT MaCustomer FROM LICHKHAM WHERE MaLichKham = ?", appointmentID).Scan(&customerID)
	
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, models.APIResponse{
				Success: false,
				Message: "Appointment not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "Failed to find appointment",
				Error:   err.Error(),
			})
		}
		return
	}

	if userType.(string) == "CUSTOMER" && customerID != userID.(string) {
		c.JSON(http.StatusForbidden, models.APIResponse{
			Success: false,
			Message: "You can only cancel your own appointments",
		})
		return
	}

	_, err = h.db.Exec("UPDATE LICHKHAM SET TrangThai = 'HUY_LICH' WHERE MaLichKham = ?", appointmentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to cancel appointment",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Appointment cancelled successfully",
	})
}