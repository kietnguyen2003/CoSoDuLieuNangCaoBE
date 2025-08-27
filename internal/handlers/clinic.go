package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

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
	query := `SELECT maPhongKham, tenPhongKham, diaChi, soDienThoai, email FROM PHONGKHAM`

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
	log.Printf("Fetching clinic with ID: %s", clinicID)

	var clinic models.Clinic
	query := `SELECT maPhongKham, tenPhongKham, diaChi, soDienThoai, email FROM PHONGKHAM WHERE maPhongKham = @p1`

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
	clinicID := c.Param("id")
	chuyenKhoa := c.Query("chuyen_khoa")
	userID, exists := c.Get("user_id")
	userType, _ := c.Get("user_type")

	log.Printf("GetDoctors called for clinic: %s", clinicID)

	// Check if user is a customer to show previously visited doctors first
	var previousDoctors []string
	if exists && userType == "CUSTOMER" {
		// Get doctors this customer has visited before at this clinic
		prevQuery := `
			SELECT DISTINCT l.maBacSi 
			FROM LICHKHAM l 
			WHERE l.maCustomer = @p1 AND l.maPhongKham = @p2 AND l.trangThai = 'COMPLETED'
			ORDER BY MAX(l.ngayGioKham) DESC
		`
		rows, err := h.db.Query(prevQuery, userID, clinicID)
		if err == nil {
			defer rows.Close()
			for rows.Next() {
				var doctorID string
				if rows.Scan(&doctorID) == nil {
					previousDoctors = append(previousDoctors, doctorID)
				}
			}
		}
	}

	// Main query to get active doctors who work at this clinic
	query := `
		SELECT DISTINCT u.userID, u.hoTen, u.soDienThoai, u.email, 
		       d.chuyenKhoa, d.namKinhNghiem, d.bangCap, d.maGiayPhep
		FROM [USER] u 
		JOIN BACSI d ON u.userID = d.maUser
		JOIN LICHLAMVIEC llv ON d.maUser = llv.maBacSi
		WHERE u.status = 'ACTIVE' AND llv.maPhongKham = @p1
	`
	args := []interface{}{clinicID}

	if chuyenKhoa != "" {
		query += " AND d.chuyenKhoa LIKE @p" + fmt.Sprintf("%d", len(args)+1)
		args = append(args, "%"+chuyenKhoa+"%")
	}

	log.Printf("Executing query: %s with args: %v", query, args)

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

	// Helper function to check if doctor was visited before
	wasPreviouslyVisited := func(doctorID string) bool {
		for _, prevID := range previousDoctors {
			if prevID == doctorID {
				return true
			}
		}
		return false
	}

	var allDoctors []models.Doctor
	var previouslyVisitedDoctors []models.Doctor
	var otherDoctors []models.Doctor

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

		if wasPreviouslyVisited(doctor.MaUser) {
			previouslyVisitedDoctors = append(previouslyVisitedDoctors, doctor)
		} else {
			otherDoctors = append(otherDoctors, doctor)
		}
		allDoctors = append(allDoctors, doctor)
	}

	// Sort previously visited doctors first, then other doctors
	sortedDoctors := append(previouslyVisitedDoctors, otherDoctors...)

	response := map[string]interface{}{
		"all_doctors": sortedDoctors,
		"previously_visited": previouslyVisitedDoctors,
		"other_doctors": otherDoctors,
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Doctors retrieved successfully", 
		Data:    response,
	})
}

func (h *ClinicHandler) GetSpecialties(c *gin.Context) {
	query := `
		SELECT DISTINCT d.chuyenKhoa 
		FROM BACSI d 
		JOIN [USER] u ON d.maUser = u.userID 
		WHERE u.status = 'ACTIVE' AND d.chuyenKhoa IS NOT NULL
		ORDER BY d.chuyenKhoa
	`
	
	rows, err := h.db.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to retrieve specialties",
			Error:   err.Error(),
		})
		return
	}
	defer rows.Close()

	var specialties []string
	for rows.Next() {
		var specialty string
		if err := rows.Scan(&specialty); err == nil {
			specialties = append(specialties, specialty)
		}
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Specialties retrieved successfully",
		Data:    specialties,
	})
}

func (h *ClinicHandler) GetSchedules(c *gin.Context) {
	clinicID := c.Param("id")
	doctorID := c.Query("doctor_id")
	date := c.Query("date")

	log.Printf("GetSchedules called with clinicID: %s, doctorID: %s, date: %s", clinicID, doctorID, date)

	if doctorID == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "doctor_id is required",
		})
		return
	}

	if date == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "date is required",
		})
		return
	}

	// Get doctor's work schedule for the date
	workScheduleQuery := `
		SELECT gioBatDau, gioKetThuc
		FROM LICHLAMVIEC
		WHERE maBacSi = @p1 AND ngayLamViec = @p2 AND status = 'AVAILABLE'
	`
	
	log.Printf("Executing query with doctorID: %s, date: %s", doctorID, date)
	
	var startTime, endTime string
	err := h.db.QueryRow(workScheduleQuery, doctorID, date).Scan(&startTime, &endTime)
	if err != nil {
		log.Printf("Work schedule query error: %s", err.Error())
		if err == sql.ErrNoRows {
			log.Printf("No work schedule found for doctor %s on date %s", doctorID, date)
			c.JSON(http.StatusOK, models.APIResponse{
				Success: true,
				Message: "No work schedule found for this date",
				Data:    []string{},
			})
		} else {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "Failed to retrieve work schedule",
				Error:   err.Error(),
			})
		}
		return
	}

	log.Printf("Found work schedule: %s - %s", startTime, endTime)

	// Get already booked appointments for this doctor on this date  
	bookedQuery := `
		SELECT ngayGioKham
		FROM LICHKHAM
		WHERE maBacSi = @p1 AND CAST(ngayGioKham AS DATE) = @p2 
		AND trangThai NOT IN ('CANCELLED', 'NO_SHOW')
	`
	
	rows, err := h.db.Query(bookedQuery, doctorID, date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to retrieve booked appointments",
			Error:   err.Error(),
		})
		return
	}
	defer rows.Close()

	var bookedTimes []string
	for rows.Next() {
		var bookedTime string
		if rows.Scan(&bookedTime) == nil {
			// Extract time part (HH:MM)
			if len(bookedTime) >= 16 { // "2025-08-30T10:00:00Z"
				timeStr := bookedTime[11:16] // Extract "10:00"
				bookedTimes = append(bookedTimes, timeStr)
			}
		}
	}

	// Generate available time slots (every hour between start and end time)
	availableSlots := generateTimeSlots(startTime, endTime, bookedTimes)

	log.Printf("Generated %d available slots: %v", len(availableSlots), availableSlots)

	response := map[string]interface{}{
		"work_schedule": map[string]string{
			"start_time": startTime,
			"end_time":   endTime,
		},
		"available_slots": availableSlots,
		"booked_times":    bookedTimes,
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Available time slots retrieved successfully",
		Data:    response,
	})
}

// Helper function to generate time slots
func generateTimeSlots(startTime, endTime string, bookedTimes []string) []string {
	var slots []string
	
	// Parse start and end times (handle both datetime and time formats)
	startHour, startMin := parseTime(startTime)
	endHour, endMin := parseTime(endTime)
	
	log.Printf("Parsed times: start %d:%d, end %d:%d", startHour, startMin, endHour, endMin)
	
	// Create booked times map for quick lookup
	bookedMap := make(map[string]bool)
	for _, booked := range bookedTimes {
		bookedMap[booked] = true
	}
	
	// Generate hourly slots
	for hour := startHour; hour < endHour; hour++ {
		timeSlot := fmt.Sprintf("%02d:00", hour)
		
		// Only add if not booked
		if !bookedMap[timeSlot] {
			slots = append(slots, timeSlot)
		}
	}
	
	log.Printf("Generated slots: %v", slots)
	return slots
}

// Helper function to parse time string to hour and minute
func parseTime(timeStr string) (int, int) {
	// Handle datetime format like "0001-01-01T07:00:00Z"
	if strings.Contains(timeStr, "T") {
		parts := strings.Split(timeStr, "T")
		if len(parts) > 1 {
			timePart := parts[1]
			if strings.Contains(timePart, "Z") {
				timePart = strings.Split(timePart, "Z")[0]
			}
			timeStr = timePart
		}
	}
	
	// Handle both "HH:MM:SS" and "HH:MM" formats
	parts := strings.Split(timeStr, ":")
	if len(parts) >= 2 {
		hour, _ := strconv.Atoi(parts[0])
		min, _ := strconv.Atoi(parts[1])
		return hour, min
	}
	return 9, 0 // default to 9:00 AM if parsing fails
}
