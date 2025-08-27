package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
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

	var query string
	var args []interface{}

	switch userType.(string) {
	case "CUSTOMER":
		query = `
			SELECT l.maLichKham, l.maCustomer, l.maBacSi, l.maPhongKham
			       , l.trangThai, l.ghiChu, l.createdAt, l.ngaygiokham,
			       u.hoTen as TenBacSi, p.tenPhongKham
			FROM LICHKHAM l
			JOIN [USER] u ON l.maBacSi = u.userID
			JOIN PHONGKHAM p ON l.maPhongKham = p.maPhongKham
			WHERE l.maCustomer = @p1
		`
		args = append(args, userID)

	case "DOCTOR":
		query = `
			SELECT l.maLichKham, l.maCustomer, l.maBacSi, l.maPhongKham
			       , l.trangThai, l.ghiChu, l.createdAt, l.ngaygiokham,
			       u.hoTen as TenBacSi, p.tenPhongKham
			FROM LICHKHAM l
			JOIN [USER] u ON l.maBacSi = u.userID
			JOIN PHONGKHAM p ON l.maPhongKham = p.maPhongKham
			WHERE l.maBacSi = @p1
		`
		args = append(args, userID)

	default:
		query = `
			SELECT l.maLichKham, l.maCustomer, l.maBacSi, l.maPhongKham
			       , l.trangThai, l.ghiChu, l.createdAt, l.ngaygiokham,
			       uc.hoTen as TenKhachHang, ud.hoTen as TenBacSi, p.tenPhongKham
			FROM LICHKHAM l
			JOIN [USER] uc ON l.maCustomer = uc.userID
			JOIN [USER] ud ON l.maBacSi = ud.userID
			JOIN PHONGKHAM p ON l.maPhongKham = p.maPhongKham
			WHERE 1=1
		`
	}

	if status != "" {
		query += " AND l.trangThai = @p" + fmt.Sprintf("%d", len(args)+1)
		args = append(args, status)
	}

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
		var ngayDat, ngaygiokham interface{}
		var ghiChu, tenPerson, tenPhongKham interface{}

		if userType.(string) == "CUSTOMER" || userType.(string) == "DOCTOR" {
			err := rows.Scan(&maLichKham, &maCustomer, &maBacSi, &maPhongKham,
				&trangThai, &ghiChu, &ngayDat, &ngaygiokham, &tenPerson, &tenPhongKham)
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
				&trangThai, &ghiChu, &ngayDat, &ngaygiokham, &tenKhachHang, &tenBacSi, &tenPhongKham)
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
		appointment["trang_thai"] = trangThai
		appointment["ghi_chu"] = ghiChu
		appointment["ngay_dat"] = ngayDat
		appointment["ngay_gio_kham"] = ngaygiokham
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

	// Use a map instead of struct for debugging
	var req map[string]interface{}

	// Debug: Log raw request body first
	rawData, _ := c.GetRawData()
	log.Printf("DEBUG: Raw JSON received: %s", string(rawData))
	
	// Manual JSON parsing using map
	if err := json.Unmarshal(rawData, &req); err != nil {
		log.Printf("DEBUG: JSON unmarshal error: %s", err.Error())
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid JSON format",
			Error:   err.Error(),
		})
		return
	}

	log.Printf("DEBUG: Parsed map: %+v", req)

	// Extract values from map
	maBacSi, _ := req["ma_bac_si"].(string)
	maPhongKham, _ := req["ma_phong_kham"].(string)
	ngayGioKham, _ := req["ngay_gio_kham"].(string)
	ghiChu, _ := req["ghi_chu"].(string)

	// Debug logging
	log.Printf("DEBUG: Creating appointment - UserID: %v, UserType: %v, DoctorID: '%s', ClinicID: '%s', DateTime: '%s'", 
		userID, userType, maBacSi, maPhongKham, ngayGioKham)

	// Manual validation
	if maBacSi == "" || maPhongKham == "" || ngayGioKham == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Missing required fields: ma_bac_si, ma_phong_kham, and ngay_gio_kham are required",
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
		WHERE maBacSi = @p1 AND ngayGioKham = @p2 AND trangThai NOT IN ('CANCELLED', 'COMPLETED')
	`, maBacSi, ngayGioKham).Scan(&count)

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
		INSERT INTO LICHKHAM (maLichKham, maCustomer, maBacSi, maPhongKham, ngayGioKham, trangThai, ghiChu, createdAt)
		VALUES (@p1, @p2, @p3, @p4, @p5, 'SCHEDULED', @p6, GETDATE())
	`, appointmentID, customerID, maBacSi, maPhongKham, ngayGioKham, ghiChu)

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
		SELECT l.maLichKham, l.maCustomer, l.maBacSi, l.maPhongKham, 
		       l.ngayGioKham, l.trangThai, l.ghiChu, l.createdAt,
		       uc.hoTen as TenKhachHang, ud.hoTen as TenBacSi, p.tenPhongKham
		FROM LICHKHAM l
		JOIN [USER] uc ON l.maCustomer = uc.userID
		JOIN [USER] ud ON l.maBacSi = ud.userID
		JOIN PHONGKHAM p ON l.maPhongKham = p.maPhongKham
		WHERE l.maLichKham = @p1
	`
	args := []interface{}{appointmentID}

	if userType.(string) == "CUSTOMER" {
		query += " AND l.maCustomer = @p2"
		args = append(args, userID)
	} else if userType.(string) == "DOCTOR" {
		query += " AND l.maBacSi = @p2"
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
	err := h.db.QueryRow("SELECT maCustomer, maBacSi FROM LICHKHAM WHERE maLichKham = @p1", appointmentID).
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
		_, err = tx.Exec("UPDATE LICHKHAM SET ngayGioKham = @p1 WHERE maLichKham = @p2", ngayGioKham, appointmentID)
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
		_, err = tx.Exec("UPDATE LICHKHAM SET trangThai = @p1 WHERE maLichKham = @p2", trangThai, appointmentID)
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
		_, err = tx.Exec("UPDATE LICHKHAM SET ghiChu = @p1 WHERE maLichKham = @p2", ghiChu, appointmentID)
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
	err := h.db.QueryRow("SELECT maCustomer FROM LICHKHAM WHERE maLichKham = @p1", appointmentID).Scan(&customerID)

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

	_, err = h.db.Exec("UPDATE LICHKHAM SET trangThai = 'CANCELLED' WHERE maLichKham = @p1", appointmentID)
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
