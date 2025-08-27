package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"clinic-management/internal/models"
	"clinic-management/internal/utils"

	"github.com/gin-gonic/gin"
)

type ScheduleHandler struct {
	db *sql.DB
}

func NewScheduleHandler(db *sql.DB) *ScheduleHandler {
	return &ScheduleHandler{db: db}
}

type ScheduleRequest struct {
	MaBacSi     string `json:"ma_bac_si" binding:"required"`
	MaPhongKham string `json:"ma_phong_kham" binding:"required"`
	NgayLamViec string `json:"ngay_lam_viec" binding:"required"` // YYYY-MM-DD format
	GioBatDau   string `json:"gio_bat_dau" binding:"required"`   // HH:MM format
	GioKetThuc  string `json:"gio_ket_thuc" binding:"required"`  // HH:MM format
	Status      string `json:"status"`                           // AVAILABLE, UNAVAILABLE
}

type ScheduleResponse struct {
	MaLichLamViec string    `json:"ma_lich_lam_viec"`
	MaBacSi       string    `json:"ma_bac_si"`
	MaPhongKham   string    `json:"ma_phong_kham"`
	NgayLamViec   time.Time `json:"ngay_lam_viec"`
	GioBatDau     string    `json:"gio_bat_dau"`
	GioKetThuc    string    `json:"gio_ket_thuc"`
	Status        string    `json:"status"`
	TenBacSi      string    `json:"ten_bac_si"`
	TenPhongKham  string    `json:"ten_phong_kham"`
}

func (h *ScheduleHandler) GetSchedules(c *gin.Context) {
	userID, _ := c.Get("user_id")
	userType, _ := c.Get("user_type")
	doctorID := c.Query("doctor_id")
	clinicID := c.Query("clinic_id")
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	var query string
	var args []interface{}

	// Base query
	baseQuery := `
		SELECT ll.maLichLamViec, ll.maBacSi, ll.maPhongKham, ll.ngayLamViec,
		       ll.gioBatDau, ll.gioKetThuc, ll.status,
		       u.hoTen as tenBacSi, p.tenPhongKham
		FROM LICHLAMVIEC ll
		JOIN [USER] u ON ll.maBacSi = u.userID
		JOIN PHONGKHAM p ON ll.maPhongKham = p.maPhongKham
		WHERE 1=1
	`

	// Filter by user role
	switch userType.(string) {
	case "DOCTOR":
		query = baseQuery + " AND ll.maBacSi = @p1"
		args = append(args, userID)
	case "CLINIC_MANAGER":
		// Clinic managers can see schedules for their clinic
		// First get the clinic manager's clinic
		var managerClinic string
		err := h.db.QueryRow("SELECT maPhongKham FROM QUANLYPHONGKHAM WHERE maUser = @p1", userID).Scan(&managerClinic)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "Failed to get manager's clinic",
				Error:   err.Error(),
			})
			return
		}
		query = baseQuery + " AND ll.maPhongKham = @p1"
		args = append(args, managerClinic)
	default:
		// Admin and other roles can see all schedules
		query = baseQuery
	}

	// Additional filters
	if doctorID != "" {
		query += fmt.Sprintf(" AND ll.maBacSi = @p%d", len(args)+1)
		args = append(args, doctorID)
	}

	if clinicID != "" {
		query += fmt.Sprintf(" AND ll.maPhongKham = @p%d", len(args)+1)
		args = append(args, clinicID)
	}

	if dateFrom != "" {
		query += fmt.Sprintf(" AND ll.ngayLamViec >= @p%d", len(args)+1)
		args = append(args, dateFrom)
	}

	if dateTo != "" {
		query += fmt.Sprintf(" AND ll.ngayLamViec <= @p%d", len(args)+1)
		args = append(args, dateTo)
	}

	query += " ORDER BY ll.ngayLamViec, ll.gioBatDau"

	rows, err := h.db.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to retrieve schedules",
			Error:   err.Error(),
		})
		return
	}
	defer rows.Close()

	var schedules []map[string]interface{}
	for rows.Next() {
		var maLichLamViec, maBacSi, maPhongKham, status sql.NullString
		var gioBatDau, gioKetThuc, tenBacSi, tenPhongKham sql.NullString
		var ngayLamViec sql.NullTime

		err := rows.Scan(&maLichLamViec, &maBacSi, &maPhongKham, &ngayLamViec,
			&gioBatDau, &gioKetThuc, &status, &tenBacSi, &tenPhongKham)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "Failed to scan schedule data",
				Error:   err.Error(),
			})
			return
		}

		schedule := map[string]interface{}{
			"ma_lich_lam_viec": maLichLamViec.String,
			"ma_bac_si":        maBacSi.String,
			"ma_phong_kham":    maPhongKham.String,
			"ngay_lam_viec":    ngayLamViec.Time,
			"gio_bat_dau":      gioBatDau.String,
			"gio_ket_thuc":     gioKetThuc.String,
			"status":           status.String,
			"ten_bac_si":       tenBacSi.String,
			"ten_phong_kham":   tenPhongKham.String,
		}

		schedules = append(schedules, schedule)
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Schedules retrieved successfully",
		Data:    schedules,
	})
}

func (h *ScheduleHandler) GetSchedule(c *gin.Context) {
	scheduleID := c.Param("id")
	userID, _ := c.Get("user_id")
	userType, _ := c.Get("user_type")

	query := `
		SELECT ll.maLichLamViec, ll.maBacSi, ll.maPhongKham, ll.ngayLamViec,
		       ll.gioBatDau, ll.gioKetThuc, ll.status,
		       u.hoTen as tenBacSi, p.tenPhongKham
		FROM LICHLAMVIEC ll
		JOIN [USER] u ON ll.maBacSi = u.userID
		JOIN PHONGKHAM p ON ll.maPhongKham = p.maPhongKham
		WHERE ll.maLichLamViec = ?
	`
	args := []interface{}{scheduleID}

	// Add role-based filtering
	if userType.(string) == "DOCTOR" {
		query += " AND ll.maBacSi = ?"
		args = append(args, userID)
	} else if userType.(string) == "CLINIC_MANAGER" {
		// Get manager's clinic
		var managerClinic string
		err := h.db.QueryRow("SELECT maPhongKham FROM QUANLYPHONGKHAM WHERE maUser = ?", userID).Scan(&managerClinic)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "Failed to get manager's clinic",
				Error:   err.Error(),
			})
			return
		}
		query += " AND ll.maPhongKham = ?"
		args = append(args, managerClinic)
	}

	var schedule map[string]interface{} = make(map[string]interface{})
	var maLichLamViec, maBacSi, maPhongKham, status sql.NullString
	var gioBatDau, gioKetThuc, tenBacSi, tenPhongKham sql.NullString
	var ngayLamViec sql.NullTime

	err := h.db.QueryRow(query, args...).Scan(
		&maLichLamViec, &maBacSi, &maPhongKham, &ngayLamViec,
		&gioBatDau, &gioKetThuc, &status, &tenBacSi, &tenPhongKham,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, models.APIResponse{
				Success: false,
				Message: "Schedule not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "Failed to retrieve schedule",
				Error:   err.Error(),
			})
		}
		return
	}

	schedule["ma_lich_lam_viec"] = maLichLamViec.String
	schedule["ma_bac_si"] = maBacSi.String
	schedule["ma_phong_kham"] = maPhongKham.String
	schedule["ngay_lam_viec"] = ngayLamViec.Time
	schedule["gio_bat_dau"] = gioBatDau.String
	schedule["gio_ket_thuc"] = gioKetThuc.String
	schedule["status"] = status.String
	schedule["ten_bac_si"] = tenBacSi.String
	schedule["ten_phong_kham"] = tenPhongKham.String

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Schedule retrieved successfully",
		Data:    schedule,
	})
}

func (h *ScheduleHandler) CreateSchedule(c *gin.Context) {
	userID, _ := c.Get("user_id")
	userType, _ := c.Get("user_type")

	var req ScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request format",
			Error:   err.Error(),
		})
		return
	}

	// Validate permissions
	if userType.(string) == "DOCTOR" {
		// Doctors can only create schedules for themselves
		if req.MaBacSi != userID.(string) {
			c.JSON(http.StatusForbidden, models.APIResponse{
				Success: false,
				Message: "You can only create schedules for yourself",
			})
			return
		}
	} else if userType.(string) != "CLINIC_MANAGER" && userType.(string) != "OPERATION_MANAGER" {
		c.JSON(http.StatusForbidden, models.APIResponse{
			Success: false,
			Message: "Only doctors, clinic managers, or operation managers can create schedules",
		})
		return
	}

	// Parse work date
	workDate, err := time.Parse("2006-01-02", req.NgayLamViec)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid date format. Use YYYY-MM-DD",
			Error:   err.Error(),
		})
		return
	}

	// Validate time format
	if !isValidTimeFormat(req.GioBatDau) || !isValidTimeFormat(req.GioKetThuc) {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid time format. Use HH:MM",
		})
		return
	}

	// Check for schedule conflicts
	var conflicts int
	err = h.db.QueryRow(`
		SELECT COUNT(*) FROM LICHLAMVIEC
		WHERE maBacSi = ? AND ngayLamViec = ? 
		AND ((gioBatDau <= ? AND gioKetThuc > ?) OR (gioBatDau < ? AND gioKetThuc >= ?))
	`, req.MaBacSi, workDate, req.GioBatDau, req.GioBatDau, req.GioKetThuc, req.GioKetThuc).Scan(&conflicts)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to check for schedule conflicts",
			Error:   err.Error(),
		})
		return
	}

	if conflicts > 0 {
		c.JSON(http.StatusConflict, models.APIResponse{
			Success: false,
			Message: "Schedule conflicts with existing schedule",
		})
		return
	}

	// Set default status
	if req.Status == "" {
		req.Status = "AVAILABLE"
	}

	// Generate schedule ID
	scheduleID := utils.GenerateScheduleID()

	// Insert schedule
	_, err = h.db.Exec(`
		INSERT INTO LICHLAMVIEC (maLichLamViec, maBacSi, maPhongKham, ngayLamViec, gioBatDau, gioKetThuc, status)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, scheduleID, req.MaBacSi, req.MaPhongKham, workDate, req.GioBatDau, req.GioKetThuc, req.Status)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to create schedule",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "Schedule created successfully",
		Data: gin.H{
			"ma_lich_lam_viec": scheduleID,
		},
	})
}

func (h *ScheduleHandler) UpdateSchedule(c *gin.Context) {
	scheduleID := c.Param("id")
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

	// Get current schedule to verify permissions
	var currentDoctorID, currentClinicID string
	err := h.db.QueryRow("SELECT maBacSi, maPhongKham FROM LICHLAMVIEC WHERE maLichLamViec = ?", scheduleID).
		Scan(&currentDoctorID, &currentClinicID)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, models.APIResponse{
				Success: false,
				Message: "Schedule not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "Failed to find schedule",
				Error:   err.Error(),
			})
		}
		return
	}

	// Validate permissions
	if userType.(string) == "DOCTOR" && currentDoctorID != userID.(string) {
		c.JSON(http.StatusForbidden, models.APIResponse{
			Success: false,
			Message: "You can only update your own schedules",
		})
		return
	} else if userType.(string) == "CLINIC_MANAGER" {
		// Verify clinic manager can update this schedule
		var managerClinic string
		err := h.db.QueryRow("SELECT maPhongKham FROM QUANLYPHONGKHAM WHERE maUser = ?", userID).Scan(&managerClinic)
		if err != nil || managerClinic != currentClinicID {
			c.JSON(http.StatusForbidden, models.APIResponse{
				Success: false,
				Message: "You can only update schedules for your clinic",
			})
			return
		}
	} else if userType.(string) != "OPERATION_MANAGER" {
		c.JSON(http.StatusForbidden, models.APIResponse{
			Success: false,
			Message: "Insufficient permissions to update schedule",
		})
		return
	}

	// Start transaction
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

	// Update fields
	if ngayLamViec, exists := updateData["ngay_lam_viec"]; exists {
		if dateStr, ok := ngayLamViec.(string); ok {
			workDate, err := time.Parse("2006-01-02", dateStr)
			if err != nil {
				c.JSON(http.StatusBadRequest, models.APIResponse{
					Success: false,
					Message: "Invalid date format. Use YYYY-MM-DD",
					Error:   err.Error(),
				})
				return
			}
			_, err = tx.Exec("UPDATE LICHLAMVIEC SET ngayLamViec = ? WHERE maLichLamViec = ?", workDate, scheduleID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, models.APIResponse{
					Success: false,
					Message: "Failed to update schedule date",
					Error:   err.Error(),
				})
				return
			}
		}
	}

	if gioBatDau, exists := updateData["gio_bat_dau"]; exists {
		if timeStr, ok := gioBatDau.(string); ok && isValidTimeFormat(timeStr) {
			_, err = tx.Exec("UPDATE LICHLAMVIEC SET gioBatDau = ? WHERE maLichLamViec = ?", timeStr, scheduleID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, models.APIResponse{
					Success: false,
					Message: "Failed to update start time",
					Error:   err.Error(),
				})
				return
			}
		} else {
			c.JSON(http.StatusBadRequest, models.APIResponse{
				Success: false,
				Message: "Invalid start time format. Use HH:MM",
			})
			return
		}
	}

	if gioKetThuc, exists := updateData["gio_ket_thuc"]; exists {
		if timeStr, ok := gioKetThuc.(string); ok && isValidTimeFormat(timeStr) {
			_, err = tx.Exec("UPDATE LICHLAMVIEC SET gioKetThuc = ? WHERE maLichLamViec = ?", timeStr, scheduleID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, models.APIResponse{
					Success: false,
					Message: "Failed to update end time",
					Error:   err.Error(),
				})
				return
			}
		} else {
			c.JSON(http.StatusBadRequest, models.APIResponse{
				Success: false,
				Message: "Invalid end time format. Use HH:MM",
			})
			return
		}
	}

	if status, exists := updateData["status"]; exists {
		if statusStr, ok := status.(string); ok {
			if statusStr != "AVAILABLE" && statusStr != "UNAVAILABLE" {
				c.JSON(http.StatusBadRequest, models.APIResponse{
					Success: false,
					Message: "Invalid status. Use AVAILABLE or UNAVAILABLE",
				})
				return
			}
			_, err = tx.Exec("UPDATE LICHLAMVIEC SET status = ? WHERE maLichLamViec = ?", statusStr, scheduleID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, models.APIResponse{
					Success: false,
					Message: "Failed to update schedule status",
					Error:   err.Error(),
				})
				return
			}
		}
	}

	if maPhongKham, exists := updateData["ma_phong_kham"]; exists {
		_, err = tx.Exec("UPDATE LICHLAMVIEC SET maPhongKham = ? WHERE maLichLamViec = ?", maPhongKham, scheduleID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "Failed to update clinic",
				Error:   err.Error(),
			})
			return
		}
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to update schedule",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Schedule updated successfully",
	})
}

func (h *ScheduleHandler) DeleteSchedule(c *gin.Context) {
	scheduleID := c.Param("id")
	userID, _ := c.Get("user_id")
	userType, _ := c.Get("user_type")

	// Get current schedule to verify permissions
	var currentDoctorID, currentClinicID string
	err := h.db.QueryRow("SELECT maBacSi, maPhongKham FROM LICHLAMVIEC WHERE maLichLamViec = ?", scheduleID).
		Scan(&currentDoctorID, &currentClinicID)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, models.APIResponse{
				Success: false,
				Message: "Schedule not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "Failed to find schedule",
				Error:   err.Error(),
			})
		}
		return
	}

	// Validate permissions
	if userType.(string) == "DOCTOR" && currentDoctorID != userID.(string) {
		c.JSON(http.StatusForbidden, models.APIResponse{
			Success: false,
			Message: "You can only delete your own schedules",
		})
		return
	} else if userType.(string) == "CLINIC_MANAGER" {
		// Verify clinic manager can delete this schedule
		var managerClinic string
		err := h.db.QueryRow("SELECT maPhongKham FROM QUANLYPHONGKHAM WHERE maUser = ?", userID).Scan(&managerClinic)
		if err != nil || managerClinic != currentClinicID {
			c.JSON(http.StatusForbidden, models.APIResponse{
				Success: false,
				Message: "You can only delete schedules for your clinic",
			})
			return
		}
	} else if userType.(string) != "OPERATION_MANAGER" {
		c.JSON(http.StatusForbidden, models.APIResponse{
			Success: false,
			Message: "Insufficient permissions to delete schedule",
		})
		return
	}

	// Check if there are appointments for this schedule
	var appointmentCount int
	err = h.db.QueryRow(`
		SELECT COUNT(*) FROM LICHKHAM l
		JOIN LICHLAMVIEC ll ON l.maBacSi = ll.maBacSi 
		WHERE ll.maLichLamViec = ? AND l.trangThai NOT IN ('CANCELLED', 'COMPLETED')
	`, scheduleID).Scan(&appointmentCount)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to check for existing appointments",
			Error:   err.Error(),
		})
		return
	}

	if appointmentCount > 0 {
		c.JSON(http.StatusConflict, models.APIResponse{
			Success: false,
			Message: "Cannot delete schedule with active appointments",
		})
		return
	}

	// Delete schedule
	_, err = h.db.Exec("DELETE FROM LICHLAMVIEC WHERE maLichLamViec = ?", scheduleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to delete schedule",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Schedule deleted successfully",
	})
}

// Helper function to validate time format HH:MM
func isValidTimeFormat(timeStr string) bool {
	_, err := time.Parse("15:04", timeStr)
	return err == nil
}