package handlers

import (
	"database/sql"
	"net/http"

	"clinic-management/internal/models"

	"github.com/gin-gonic/gin"
)

type ClinicHandler struct {
	db *sql.DB
}

func NewClinicHandler(db *sql.DB) *ClinicHandler {
	return &ClinicHandler{db: db}
}

func (h *ClinicHandler) GetClinics(c *gin.Context) {
	query := `SELECT MaPhongKham, TenPhongKham, DiaChi, SoDienThoai, Email FROM PHONGKHAM`

	rows, err := h.db.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to retrieve clinics",
			Error:   err.Error(),
		})
		return
	}
	defer rows.Close()

	var clinics []models.Clinic
	for rows.Next() {
		var clinic models.Clinic
		err := rows.Scan(&clinic.MaPhongKham, &clinic.TenPhongKham, &clinic.DiaChi, &clinic.SoDienThoai, &clinic.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "Failed to scan clinic data",
				Error:   err.Error(),
			})
			return
		}
		clinics = append(clinics, clinic)
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Clinics retrieved successfully",
		Data:    clinics,
	})
}

func (h *ClinicHandler) GetClinic(c *gin.Context) {
	clinicID := c.Param("id")
	
	var clinic models.Clinic
	query := `SELECT MaPhongKham, TenPhongKham, DiaChi, SoDienThoai, Email FROM PHONGKHAM WHERE MaPhongKham = ?`

	err := h.db.QueryRow(query, clinicID).Scan(
		&clinic.MaPhongKham, &clinic.TenPhongKham, &clinic.DiaChi, &clinic.SoDienThoai, &clinic.Email,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, models.APIResponse{
				Success: false,
				Message: "Clinic not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "Failed to retrieve clinic",
				Error:   err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Clinic retrieved successfully",
		Data:    clinic,
	})
}

func (h *ClinicHandler) GetDoctors(c *gin.Context) {
	chuyenKhoa := c.Query("chuyen_khoa")

	query := `
		SELECT u.MaUser, u.HoTen, u.SoDienThoai, u.Email, 
		       d.ChuyenKhoa, d.NamKinhNghiem, d.BangCap, d.SoGiayPhepHanhNghe
		FROM [USER] u 
		JOIN BACSI d ON u.MaUser = d.MaUser
		WHERE u.TrangThai = 'ACTIVE'
	`
	args := []interface{}{}

	if chuyenKhoa != "" {
		query += " AND d.ChuyenKhoa LIKE ?"
		args = append(args, "%"+chuyenKhoa+"%")
	}

	rows, err := h.db.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to retrieve doctors",
			Error:   err.Error(),
		})
		return
	}
	defer rows.Close()

	var doctors []models.Doctor
	for rows.Next() {
		var doctor models.Doctor
		err := rows.Scan(
			&doctor.MaUser, &doctor.HoTen, &doctor.SoDienThoai, &doctor.Email,
			&doctor.ChuyenKhoa, &doctor.NamKinhNghiem, &doctor.BangCap, &doctor.SoGiayPhepHanhNghe,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "Failed to scan doctor data",
				Error:   err.Error(),
			})
			return
		}
		doctors = append(doctors, doctor)
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Doctors retrieved successfully",
		Data:    doctors,
	})
}

func (h *ClinicHandler) GetSchedules(c *gin.Context) {
	clinicID := c.Param("id")
	doctorID := c.Query("doctor_id")
	date := c.Query("date")

	query := `
		SELECT MaLichLamViec, MaBacSi, MaPhongKham, NgayLamViec, GioBatDau, GioKetThuc, TrangThai
		FROM LICHLAMVIEC
		WHERE MaPhongKham = ? AND TrangThai = 'SCHEDULED'
	`
	args := []interface{}{clinicID}

	if doctorID != "" {
		query += " AND MaBacSi = ?"
		args = append(args, doctorID)
	}

	if date != "" {
		query += " AND NgayLamViec = ?"
		args = append(args, date)
	}

	query += " ORDER BY NgayLamViec, GioBatDau"

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

	var schedules []models.WorkSchedule
	for rows.Next() {
		var schedule models.WorkSchedule
		err := rows.Scan(
			&schedule.MaLichLamViec, &schedule.MaBacSi, &schedule.MaPhongKham,
			&schedule.NgayLamViec, &schedule.GioBatDau, &schedule.GioKetThuc, &schedule.TrangThai,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "Failed to scan schedule data",
				Error:   err.Error(),
			})
			return
		}
		schedules = append(schedules, schedule)
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Schedules retrieved successfully",
		Data:    schedules,
	})
}