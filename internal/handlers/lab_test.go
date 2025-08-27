package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"clinic-management/internal/models"
	"clinic-management/internal/utils"

	"github.com/gin-gonic/gin"
)

type LabTestHandler struct {
	db *sql.DB
}

func NewLabTestHandler(db *sql.DB) *LabTestHandler {
	return &LabTestHandler{db: db}
}

type LabTestRequest struct {
	MaHoSo          string `json:"ma_ho_so" binding:"required"`
	LoaiXetNghiem   string `json:"loai_xet_nghiem" binding:"required"`
	GhiChu          string `json:"ghi_chu"`
	NgayXetNghiem   string `json:"ngay_xet_nghiem"` // optional, defaults to today
}

type LabTestResponse struct {
	MaXetNghiem   string    `json:"ma_xet_nghiem"`
	MaHoSo        string    `json:"ma_ho_so"`
	LoaiXetNghiem string    `json:"loai_xet_nghiem"`
	NgayXetNghiem time.Time `json:"ngay_xet_nghiem"`
	KetQua        string    `json:"ket_qua"`
	GhiChu        string    `json:"ghi_chu"`
	TenKhachHang  string    `json:"ten_khach_hang"`
	TenBacSi      string    `json:"ten_bac_si"`
	Status        string    `json:"status"` // ordered, collected, processing, completed
}

func (h *LabTestHandler) GetLabTests(c *gin.Context) {
	userID, _ := c.Get("user_id")
	userType, _ := c.Get("user_type")
	status := c.Query("status")
	maHoSo := c.Query("ma_ho_so")

	var query string
	var args []interface{}

	// Base query
	baseQuery := `
		SELECT xn.maXetNghiem, xn.maHoSo, xn.loaiXetNghiem, xn.ngayXetNghiem, 
		       xn.ketQua, xn.ghiChu,
		       h.maCustomer, h.maBacSi, h.ngayKham,
		       uc.hoTen as tenKhachHang, ud.hoTen as tenBacSi
		FROM XETNGHIEM xn
		JOIN HOSO h ON xn.maHoSo = h.maHoSo
		JOIN [USER] uc ON h.maCustomer = uc.userID
		JOIN [USER] ud ON h.maBacSi = ud.userID
		WHERE 1=1
	`

	// Filter by user role
	switch userType.(string) {
	case "CUSTOMER":
		query = baseQuery + " AND h.maCustomer = ?"
		args = append(args, userID)
	case "DOCTOR":
		query = baseQuery + " AND h.maBacSi = ?"
		args = append(args, userID)
	default:
		query = baseQuery
	}

	// Additional filters
	if maHoSo != "" {
		query += " AND xn.maHoSo = ?"
		args = append(args, maHoSo)
	}

	query += " ORDER BY xn.ngayXetNghiem DESC"

	rows, err := h.db.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to retrieve lab tests",
			Error:   err.Error(),
		})
		return
	}
	defer rows.Close()

	var labTests []map[string]interface{}
	for rows.Next() {
		var maXetNghiem, maHoSo, loaiXetNghiem, ketQua, ghiChu sql.NullString
		var maCustomer, maBacSi, tenKhachHang, tenBacSi sql.NullString
		var ngayXetNghiem, ngayKham sql.NullTime

		err := rows.Scan(&maXetNghiem, &maHoSo, &loaiXetNghiem, &ngayXetNghiem,
			&ketQua, &ghiChu, &maCustomer, &maBacSi, &ngayKham,
			&tenKhachHang, &tenBacSi)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "Failed to scan lab test data",
				Error:   err.Error(),
			})
			return
		}

		// Determine status based on data
		testStatus := "ordered"
		if ketQua.String != "" {
			testStatus = "completed"
		} else if ngayXetNghiem.Valid && ngayXetNghiem.Time.Before(time.Now()) {
			testStatus = "processing"
		}

		labTest := map[string]interface{}{
			"ma_xet_nghiem":   maXetNghiem.String,
			"ma_ho_so":        maHoSo.String,
			"loai_xet_nghiem": loaiXetNghiem.String,
			"ngay_xet_nghiem": ngayXetNghiem.Time,
			"ket_qua":         ketQua.String,
			"ghi_chu":         ghiChu.String,
			"ma_customer":     maCustomer.String,
			"ma_bac_si":       maBacSi.String,
			"ten_khach_hang":  tenKhachHang.String,
			"ten_bac_si":      tenBacSi.String,
			"status":          testStatus,
			"ngay_kham":       ngayKham.Time,
		}

		labTests = append(labTests, labTest)
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Lab tests retrieved successfully",
		Data:    labTests,
	})
}

func (h *LabTestHandler) GetLabTest(c *gin.Context) {
	labTestID := c.Param("id")
	userID, _ := c.Get("user_id")
	userType, _ := c.Get("user_type")

	query := `
		SELECT xn.maXetNghiem, xn.maHoSo, xn.loaiXetNghiem, xn.ngayXetNghiem,
		       xn.ketQua, xn.ghiChu, xn.fileName,
		       h.maCustomer, h.maBacSi, h.ngayKham,
		       uc.hoTen as tenKhachHang, ud.hoTen as tenBacSi
		FROM XETNGHIEM xn
		JOIN HOSO h ON xn.maHoSo = h.maHoSo
		JOIN [USER] uc ON h.maCustomer = uc.userID
		JOIN [USER] ud ON h.maBacSi = ud.userID
		WHERE xn.maXetNghiem = ?
	`
	args := []interface{}{labTestID}

	// Add role-based filtering
	if userType.(string) == "CUSTOMER" {
		query += " AND h.maCustomer = ?"
		args = append(args, userID)
	} else if userType.(string) == "DOCTOR" {
		query += " AND h.maBacSi = ?"
		args = append(args, userID)
	}

	var labTest map[string]interface{} = make(map[string]interface{})
	var maXetNghiem, maHoSo, loaiXetNghiem, ketQua, ghiChu, fileName sql.NullString
	var maCustomer, maBacSi, tenKhachHang, tenBacSi sql.NullString
	var ngayXetNghiem, ngayKham sql.NullTime

	err := h.db.QueryRow(query, args...).Scan(
		&maXetNghiem, &maHoSo, &loaiXetNghiem, &ngayXetNghiem,
		&ketQua, &ghiChu, &fileName,
		&maCustomer, &maBacSi, &ngayKham,
		&tenKhachHang, &tenBacSi,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, models.APIResponse{
				Success: false,
				Message: "Lab test not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "Failed to retrieve lab test",
				Error:   err.Error(),
			})
		}
		return
	}

	// Determine status
	testStatus := "ordered"
	if ketQua.String != "" {
		testStatus = "completed"
	} else if ngayXetNghiem.Valid && ngayXetNghiem.Time.Before(time.Now()) {
		testStatus = "processing"
	}

	labTest["ma_xet_nghiem"] = maXetNghiem.String
	labTest["ma_ho_so"] = maHoSo.String
	labTest["loai_xet_nghiem"] = loaiXetNghiem.String
	labTest["ngay_xet_nghiem"] = ngayXetNghiem.Time
	labTest["ket_qua"] = ketQua.String
	labTest["ghi_chu"] = ghiChu.String
	labTest["file_name"] = fileName.String
	labTest["ma_customer"] = maCustomer.String
	labTest["ma_bac_si"] = maBacSi.String
	labTest["ten_khach_hang"] = tenKhachHang.String
	labTest["ten_bac_si"] = tenBacSi.String
	labTest["status"] = testStatus
	labTest["ngay_kham"] = ngayKham.Time

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Lab test retrieved successfully",
		Data:    labTest,
	})
}

func (h *LabTestHandler) CreateLabOrder(c *gin.Context) {
	userID, _ := c.Get("user_id")
	userType, _ := c.Get("user_type")

	// Only doctors can create lab orders
	if userType.(string) != "DOCTOR" {
		c.JSON(http.StatusForbidden, models.APIResponse{
			Success: false,
			Message: "Only doctors can create lab orders",
		})
		return
	}

	var req LabTestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request format",
			Error:   err.Error(),
		})
		return
	}

	// Verify the medical record belongs to this doctor
	var doctorID string
	err := h.db.QueryRow("SELECT maBacSi FROM HOSO WHERE maHoSo = ?", req.MaHoSo).Scan(&doctorID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, models.APIResponse{
				Success: false,
				Message: "Medical record not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "Failed to verify medical record",
				Error:   err.Error(),
			})
		}
		return
	}

	if doctorID != userID.(string) {
		c.JSON(http.StatusForbidden, models.APIResponse{
			Success: false,
			Message: "You can only create lab orders for your own patients",
		})
		return
	}

	// Generate lab test ID
	labTestID := utils.GenerateLabTestID()

	// Parse test date
	var testDate time.Time
	if req.NgayXetNghiem != "" {
		parsedDate, err := time.Parse("2006-01-02", req.NgayXetNghiem)
		if err != nil {
			c.JSON(http.StatusBadRequest, models.APIResponse{
				Success: false,
				Message: "Invalid date format. Use YYYY-MM-DD",
				Error:   err.Error(),
			})
			return
		}
		testDate = parsedDate
	} else {
		testDate = time.Now()
	}

	// Insert lab test order
	_, err = h.db.Exec(`
		INSERT INTO XETNGHIEM (maXetNghiem, maHoSo, loaiXetNghiem, ngayXetNghiem, ghiChu)
		VALUES (?, ?, ?, ?, ?)
	`, labTestID, req.MaHoSo, req.LoaiXetNghiem, testDate, req.GhiChu)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to create lab order",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "Lab order created successfully",
		Data: gin.H{
			"ma_xet_nghiem": labTestID,
		},
	})
}

func (h *LabTestHandler) UpdateLabTest(c *gin.Context) {
	labTestID := c.Param("id")
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

	// Verify lab test belongs to this doctor (for updates) or lab technician
	var doctorID string
	err := h.db.QueryRow(`
		SELECT h.maBacSi FROM XETNGHIEM xn
		JOIN HOSO h ON xn.maHoSo = h.maHoSo
		WHERE xn.maXetNghiem = ?
	`, labTestID).Scan(&doctorID)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, models.APIResponse{
				Success: false,
				Message: "Lab test not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "Failed to verify lab test",
				Error:   err.Error(),
			})
		}
		return
	}

	// Only doctors can update their own lab orders, or lab technicians can update results
	if userType.(string) == "DOCTOR" && doctorID != userID.(string) {
		c.JSON(http.StatusForbidden, models.APIResponse{
			Success: false,
			Message: "You can only update your own lab orders",
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
	if ketQua, exists := updateData["ket_qua"]; exists {
		_, err = tx.Exec("UPDATE XETNGHIEM SET ketQua = ? WHERE maXetNghiem = ?", ketQua, labTestID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "Failed to update lab test results",
				Error:   err.Error(),
			})
			return
		}
	}

	if ghiChu, exists := updateData["ghi_chu"]; exists {
		_, err = tx.Exec("UPDATE XETNGHIEM SET ghiChu = ? WHERE maXetNghiem = ?", ghiChu, labTestID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "Failed to update lab test notes",
				Error:   err.Error(),
			})
			return
		}
	}

	if loaiXetNghiem, exists := updateData["loai_xet_nghiem"]; exists {
		_, err = tx.Exec("UPDATE XETNGHIEM SET loaiXetNghiem = ? WHERE maXetNghiem = ?", loaiXetNghiem, labTestID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "Failed to update lab test type",
				Error:   err.Error(),
			})
			return
		}
	}

	if ngayXetNghiem, exists := updateData["ngay_xet_nghiem"]; exists {
		if dateStr, ok := ngayXetNghiem.(string); ok {
			parsedDate, err := time.Parse("2006-01-02", dateStr)
			if err != nil {
				c.JSON(http.StatusBadRequest, models.APIResponse{
					Success: false,
					Message: "Invalid date format. Use YYYY-MM-DD",
					Error:   err.Error(),
				})
				return
			}
			_, err = tx.Exec("UPDATE XETNGHIEM SET ngayXetNghiem = ? WHERE maXetNghiem = ?", parsedDate, labTestID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, models.APIResponse{
					Success: false,
					Message: "Failed to update lab test date",
					Error:   err.Error(),
				})
				return
			}
		}
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to update lab test",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Lab test updated successfully",
	})
}

func (h *LabTestHandler) DeleteLabTest(c *gin.Context) {
	labTestID := c.Param("id")
	userID, _ := c.Get("user_id")
	userType, _ := c.Get("user_type")

	// Only doctors can delete lab orders
	if userType.(string) != "DOCTOR" {
		c.JSON(http.StatusForbidden, models.APIResponse{
			Success: false,
			Message: "Only doctors can delete lab orders",
		})
		return
	}

	// Verify lab test belongs to this doctor
	var doctorID string
	err := h.db.QueryRow(`
		SELECT h.maBacSi FROM XETNGHIEM xn
		JOIN HOSO h ON xn.maHoSo = h.maHoSo
		WHERE xn.maXetNghiem = ?
	`, labTestID).Scan(&doctorID)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, models.APIResponse{
				Success: false,
				Message: "Lab test not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "Failed to verify lab test",
				Error:   err.Error(),
			})
		}
		return
	}

	if doctorID != userID.(string) {
		c.JSON(http.StatusForbidden, models.APIResponse{
			Success: false,
			Message: "You can only delete your own lab orders",
		})
		return
	}

	// Delete lab test
	_, err = h.db.Exec("DELETE FROM XETNGHIEM WHERE maXetNghiem = ?", labTestID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to delete lab test",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Lab test deleted successfully",
	})
}

// Get common lab test types
func (h *LabTestHandler) GetLabTestTypes(c *gin.Context) {
	query := `
		SELECT DISTINCT loaiXetNghiem
		FROM XETNGHIEM
		WHERE loaiXetNghiem IS NOT NULL AND loaiXetNghiem != ''
		ORDER BY loaiXetNghiem
	`

	rows, err := h.db.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to retrieve lab test types",
			Error:   err.Error(),
		})
		return
	}
	defer rows.Close()

	var testTypes []string
	for rows.Next() {
		var testType sql.NullString
		err := rows.Scan(&testType)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "Failed to scan lab test types",
				Error:   err.Error(),
			})
			return
		}
		if testType.Valid {
			testTypes = append(testTypes, testType.String)
		}
	}

	// Add common lab test types if none exist
	if len(testTypes) == 0 {
		testTypes = []string{
			"Complete Blood Count (CBC)",
			"Basic Metabolic Panel (BMP)",
			"Hemoglobin A1c",
			"Lipid Panel",
			"Liver Function Tests",
			"Thyroid Function Tests",
			"Troponin I",
			"C-Reactive Protein",
			"Blood Culture",
			"Urine Culture",
		}
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Lab test types retrieved successfully",
		Data:    testTypes,
	})
}