package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"clinic-management/internal/models"
	"clinic-management/internal/utils"

	"github.com/gin-gonic/gin"
)

type PrescriptionHandler struct {
	db *sql.DB
}

func NewPrescriptionHandler(db *sql.DB) *PrescriptionHandler {
	return &PrescriptionHandler{db: db}
}

type PrescriptionRequest struct {
	MaHoSo      string                   `json:"ma_ho_so" binding:"required"`
	Medications []PrescriptionMedication `json:"medications" binding:"required,min=1"`
	GhiChu      string                   `json:"ghi_chu"`
}

type PrescriptionMedication struct {
	MaThuoc  string `json:"ma_thuoc" binding:"required"`
	TenThuoc string `json:"ten_thuoc" binding:"required"`
	SoLuong  int    `json:"so_luong" binding:"required,min=1"`
	CachDung string `json:"cach_dung" binding:"required"`
	GhiChu   string `json:"ghi_chu"`
}

type PrescriptionResponse struct {
	MaDonThuoc  string                   `json:"ma_don_thuoc"`
	MaHoSo      string                   `json:"ma_ho_so"`
	NgayKeDon   time.Time                `json:"ngay_ke_don"`
	GhiChu      string                   `json:"ghi_chu"`
	Medications []PrescriptionMedication `json:"medications"`
}

func (h *PrescriptionHandler) GetPrescriptions(c *gin.Context) {
	userID, _ := c.Get("user_id")
	userType, _ := c.Get("user_type")
	// status := c.Query("status")
	maHoSo := c.Query("ma_ho_so")

	var query string
	var args []interface{}

	// Base query with joins
	baseQuery := `
		SELECT DISTINCT dt.maDonThuoc, dt.maHoSo, dt.ngayHeHan as ngayKeDon, dt.ghiChu,
		       h.maCustomer, h.maBacSi, 
		       uc.hoTen as tenKhachHang, ud.hoTen as tenBacSi
		FROM DONTHUOC dt
		JOIN HOSO h ON dt.maHoSo = h.maHoSo
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
		query += " AND dt.maHoSo = ?"
		args = append(args, maHoSo)
	}

	query += " ORDER BY dt.ngayHeHan DESC"

	rows, err := h.db.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to retrieve prescriptions",
			Error:   err.Error(),
		})
		return
	}
	defer rows.Close()

	var prescriptions []map[string]interface{}
	for rows.Next() {
		var maDonThuoc, maHoSo, ghiChu, maCustomer, maBacSi, tenKhachHang, tenBacSi sql.NullString
		var ngayKeDon sql.NullTime

		err := rows.Scan(&maDonThuoc, &maHoSo, &ngayKeDon, &ghiChu,
			&maCustomer, &maBacSi, &tenKhachHang, &tenBacSi)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "Failed to scan prescription data",
				Error:   err.Error(),
			})
			return
		}

		// Get medications for this prescription
		medications, err := h.getPrescriptionMedications(maDonThuoc.String)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "Failed to get prescription medications",
				Error:   err.Error(),
			})
			return
		}

		prescription := map[string]interface{}{
			"ma_don_thuoc":   maDonThuoc.String,
			"ma_ho_so":       maHoSo.String,
			"ngay_ke_don":    ngayKeDon.Time,
			"ghi_chu":        ghiChu.String,
			"ma_customer":    maCustomer.String,
			"ma_bac_si":      maBacSi.String,
			"ten_khach_hang": tenKhachHang.String,
			"ten_bac_si":     tenBacSi.String,
			"medications":    medications,
		}

		prescriptions = append(prescriptions, prescription)
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Prescriptions retrieved successfully",
		Data:    prescriptions,
	})
}

func (h *PrescriptionHandler) getPrescriptionMedications(maDonThuoc string) ([]map[string]interface{}, error) {
	query := `
		SELECT ct.maThuoc, ct.soLuong, ct.cacDung, ct.ghiChu,
		       t.tenThuoc, t.gia, t.congDung, t.lieuLuong
		FROM CHITIETDONTHUOC ct
		JOIN THUOC t ON ct.maThuoc = t.maThuoc
		WHERE ct.maDonThuoc = ?
	`

	rows, err := h.db.Query(query, maDonThuoc)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var medications []map[string]interface{}
	for rows.Next() {
		var maThuoc, tenThuoc, cacDung, ghiChu, congDung, lieuLuong sql.NullString
		var soLuong sql.NullInt32
		var gia sql.NullFloat64

		err := rows.Scan(&maThuoc, &soLuong, &cacDung, &ghiChu,
			&tenThuoc, &gia, &congDung, &lieuLuong)
		if err != nil {
			return nil, err
		}

		medication := map[string]interface{}{
			"ma_thuoc":   maThuoc.String,
			"ten_thuoc":  tenThuoc.String,
			"so_luong":   soLuong.Int32,
			"cach_dung":  cacDung.String,
			"ghi_chu":    ghiChu.String,
			"gia":        gia.Float64,
			"cong_dung":  congDung.String,
			"lieu_luong": lieuLuong.String,
		}

		medications = append(medications, medication)
	}

	return medications, nil
}

func (h *PrescriptionHandler) GetPrescription(c *gin.Context) {
	prescriptionID := c.Param("id")
	userID, _ := c.Get("user_id")
	userType, _ := c.Get("user_type")

	query := `
		SELECT dt.maDonThuoc, dt.maHoSo, dt.ngayHeHan, dt.ghiChu,
		       h.maCustomer, h.maBacSi,
		       uc.hoTen as tenKhachHang, ud.hoTen as tenBacSi
		FROM DONTHUOC dt
		JOIN HOSO h ON dt.maHoSo = h.maHoSo
		JOIN [USER] uc ON h.maCustomer = uc.userID
		JOIN [USER] ud ON h.maBacSi = ud.userID
		WHERE dt.maDonThuoc = ?
	`
	args := []interface{}{prescriptionID}

	// Add role-based filtering
	if userType.(string) == "CUSTOMER" {
		query += " AND h.maCustomer = ?"
		args = append(args, userID)
	} else if userType.(string) == "DOCTOR" {
		query += " AND h.maBacSi = ?"
		args = append(args, userID)
	}

	var prescription map[string]interface{} = make(map[string]interface{})
	var maDonThuoc, maHoSo, ghiChu, maCustomer, maBacSi, tenKhachHang, tenBacSi sql.NullString
	var ngayKeDon sql.NullTime

	err := h.db.QueryRow(query, args...).Scan(
		&maDonThuoc, &maHoSo, &ngayKeDon, &ghiChu,
		&maCustomer, &maBacSi, &tenKhachHang, &tenBacSi,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, models.APIResponse{
				Success: false,
				Message: "Prescription not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "Failed to retrieve prescription",
				Error:   err.Error(),
			})
		}
		return
	}

	// Get medications
	medications, err := h.getPrescriptionMedications(maDonThuoc.String)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to get prescription medications",
			Error:   err.Error(),
		})
		return
	}

	prescription["ma_don_thuoc"] = maDonThuoc.String
	prescription["ma_ho_so"] = maHoSo.String
	prescription["ngay_ke_don"] = ngayKeDon.Time
	prescription["ghi_chu"] = ghiChu.String
	prescription["ma_customer"] = maCustomer.String
	prescription["ma_bac_si"] = maBacSi.String
	prescription["ten_khach_hang"] = tenKhachHang.String
	prescription["ten_bac_si"] = tenBacSi.String
	prescription["medications"] = medications

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Prescription retrieved successfully",
		Data:    prescription,
	})
}

func (h *PrescriptionHandler) CreatePrescription(c *gin.Context) {
	userID, _ := c.Get("user_id")
	userType, _ := c.Get("user_type")

	// Only doctors can create prescriptions
	if userType.(string) != "DOCTOR" {
		c.JSON(http.StatusForbidden, models.APIResponse{
			Success: false,
			Message: "Only doctors can create prescriptions",
		})
		return
	}

	var req PrescriptionRequest
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
			Message: "You can only create prescriptions for your own patients",
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

	// Generate prescription ID
	prescriptionID := utils.GeneratePrescriptionID()

	// Insert prescription
	_, err = tx.Exec(`
		INSERT INTO DONTHUOC (maDonThuoc, maHoSo, ngayHeHan, ghiChu)
		VALUES (?, ?, DATEADD(month, 3, GETDATE()), ?)
	`, prescriptionID, req.MaHoSo, req.GhiChu)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to create prescription",
			Error:   err.Error(),
		})
		return
	}

	// Insert medication details
	for _, med := range req.Medications {
		// Check if medication exists, if not create it
		var exists int
		err = tx.QueryRow("SELECT COUNT(*) FROM THUOC WHERE maThuoc = ?", med.MaThuoc).Scan(&exists)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "Failed to check medication",
				Error:   err.Error(),
			})
			return
		}

		if exists == 0 {
			// Create medication if it doesn't exist
			_, err = tx.Exec(`
				INSERT INTO THUOC (maThuoc, tenThuoc, soLuong, gia, congDung, lieuLuong)
				VALUES (?, ?, 100, 0, '', '')
			`, med.MaThuoc, med.TenThuoc)

			if err != nil {
				c.JSON(http.StatusInternalServerError, models.APIResponse{
					Success: false,
					Message: "Failed to create medication",
					Error:   err.Error(),
				})
				return
			}
		}

		// Insert prescription detail
		_, err = tx.Exec(`
			INSERT INTO CHITIETDONTHUOC (maDonThuoc, maThuoc, soLuong, cacDung, ghiChu)
			VALUES (?, ?, ?, ?, ?)
		`, prescriptionID, med.MaThuoc, med.SoLuong, med.CachDung, med.GhiChu)

		if err != nil {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "Failed to add medication to prescription",
				Error:   err.Error(),
			})
			return
		}
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to create prescription",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "Prescription created successfully",
		Data: gin.H{
			"ma_don_thuoc": prescriptionID,
		},
	})
}

func (h *PrescriptionHandler) UpdatePrescription(c *gin.Context) {
	prescriptionID := c.Param("id")
	userID, _ := c.Get("user_id")
	userType, _ := c.Get("user_type")

	// Only doctors can update prescriptions
	if userType.(string) != "DOCTOR" {
		c.JSON(http.StatusForbidden, models.APIResponse{
			Success: false,
			Message: "Only doctors can update prescriptions",
		})
		return
	}

	var req PrescriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request format",
			Error:   err.Error(),
		})
		return
	}

	// Verify prescription belongs to this doctor
	var doctorID string
	err := h.db.QueryRow(`
		SELECT h.maBacSi FROM DONTHUOC dt
		JOIN HOSO h ON dt.maHoSo = h.maHoSo
		WHERE dt.maDonThuoc = ?
	`, prescriptionID).Scan(&doctorID)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, models.APIResponse{
				Success: false,
				Message: "Prescription not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "Failed to verify prescription",
				Error:   err.Error(),
			})
		}
		return
	}

	if doctorID != userID.(string) {
		c.JSON(http.StatusForbidden, models.APIResponse{
			Success: false,
			Message: "You can only update your own prescriptions",
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

	// Update prescription
	_, err = tx.Exec("UPDATE DONTHUOC SET ghiChu = ? WHERE maDonThuoc = ?", req.GhiChu, prescriptionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to update prescription",
			Error:   err.Error(),
		})
		return
	}

	// Delete existing medication details
	_, err = tx.Exec("DELETE FROM CHITIETDONTHUOC WHERE maDonThuoc = ?", prescriptionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to update prescription medications",
			Error:   err.Error(),
		})
		return
	}

	// Insert updated medication details
	for _, med := range req.Medications {
		// Check if medication exists, if not create it
		var exists int
		err = tx.QueryRow("SELECT COUNT(*) FROM THUOC WHERE maThuoc = ?", med.MaThuoc).Scan(&exists)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "Failed to check medication",
				Error:   err.Error(),
			})
			return
		}

		if exists == 0 {
			_, err = tx.Exec(`
				INSERT INTO THUOC (maThuoc, tenThuoc, soLuong, gia, congDung, lieuLuong)
				VALUES (?, ?, 100, 0, '', '')
			`, med.MaThuoc, med.TenThuoc)

			if err != nil {
				c.JSON(http.StatusInternalServerError, models.APIResponse{
					Success: false,
					Message: "Failed to create medication",
					Error:   err.Error(),
				})
				return
			}
		}

		// Insert prescription detail
		_, err = tx.Exec(`
			INSERT INTO CHITIETDONTHUOC (maDonThuoc, maThuoc, soLuong, cacDung, ghiChu)
			VALUES (?, ?, ?, ?, ?)
		`, prescriptionID, med.MaThuoc, med.SoLuong, med.CachDung, med.GhiChu)

		if err != nil {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "Failed to update medication in prescription",
				Error:   err.Error(),
			})
			return
		}
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to update prescription",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Prescription updated successfully",
	})
}

func (h *PrescriptionHandler) GetMedications(c *gin.Context) {
	query := c.Query("q") // search query

	var sqlQuery string
	var args []interface{}

	if query != "" {
		sqlQuery = `
			SELECT maThuoc, tenThuoc, gia, congDung, lieuLuong
			FROM THUOC
			WHERE tenThuoc LIKE ? OR maThuoc LIKE ?
			ORDER BY tenThuoc
		`
		searchTerm := "%" + query + "%"
		args = []interface{}{searchTerm, searchTerm}
	} else {
		sqlQuery = `
			SELECT maThuoc, tenThuoc, gia, congDung, lieuLuong
			FROM THUOC
			ORDER BY tenThuoc
		`
	}

	rows, err := h.db.Query(sqlQuery, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to retrieve medications",
			Error:   err.Error(),
		})
		return
	}
	defer rows.Close()

	var medications []map[string]interface{}
	for rows.Next() {
		var maThuoc, tenThuoc, congDung, lieuLuong sql.NullString
		var gia sql.NullFloat64

		err := rows.Scan(&maThuoc, &tenThuoc, &gia, &congDung, &lieuLuong)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "Failed to scan medication data",
				Error:   err.Error(),
			})
			return
		}

		medication := map[string]interface{}{
			"ma_thuoc":   maThuoc.String,
			"ten_thuoc":  tenThuoc.String,
			"gia":        gia.Float64,
			"cong_dung":  congDung.String,
			"lieu_luong": lieuLuong.String,
		}

		medications = append(medications, medication)
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Medications retrieved successfully",
		Data:    medications,
	})
}
