package handlers

import (
	"database/sql"
	"net/http"

	"clinic-management/internal/models"
	"clinic-management/internal/utils"

	"github.com/gin-gonic/gin"
)

type MedicalRecordHandler struct {
	db *sql.DB
}

func NewMedicalRecordHandler(db *sql.DB) *MedicalRecordHandler {
	return &MedicalRecordHandler{db: db}
}

func (h *MedicalRecordHandler) GetMedicalRecords(c *gin.Context) {
	userID, _ := c.Get("user_id")
	userType, _ := c.Get("user_type")
	customerID := c.Query("customer_id")

	var query string
	var args []interface{}

	switch userType.(string) {
	case "CUSTOMER":
		query = `
			SELECT h.MaHoSo, h.maCustomer, h.MaBacSi, h.MaPhongKham, 
			       h.NgayKham, h.TrieuChung, h.ChanDoan, h.Huongdan, 
			       h.MaICD10, h.NgayTaiKham,
			       u.HoTen as TenBacSi, p.TenPhongKham
			FROM HOSO h
			JOIN [USER] u ON h.MaBacSi = u.userId
			JOIN PHONGKHAM p ON h.MaPhongKham = p.MaPhongKham
			WHERE h.MaCustomer = @p1
			ORDER BY h.NgayKham DESC
		`
		args = append(args, userID)

	case "DOCTOR":
		query = `
			SELECT h.MaHoSo, h.MaCustomer, h.MaBacSi, h.MaPhongKham, 
			       h.NgayKham, h.TrieuChung, h.ChanDoan, h.huongdan, 
			       h.MaICD10, h.NgayTaiKham,
			       u.HoTen as TenKhachHang, p.TenPhongKham
			FROM HOSO h
			JOIN [USER] u ON h.MaCustomer = u.userID
			JOIN PHONGKHAM p ON h.MaPhongKham = p.MaPhongKham
			WHERE h.MaBacSi = @p1
		`
		args = append(args, userID)

		if customerID != "" {
			query += " AND h.MaCustomer = @p2"
			args = append(args, customerID)
		}

		query += " ORDER BY h.NgayKham DESC"

	default:
		query = `
			SELECT h.MaHoSo, h.MaCustomer, h.MaBacSi, h.MaPhongKham, 
			       h.NgayKham, h.TrieuChung, h.ChanDoan, h.huongdan, 
			       h.MaICD10, h.NgayTaiKham,
			       uc.HoTen as TenKhachHang, ud.HoTen as TenBacSi, p.TenPhongKham
			FROM HOSO h
			JOIN [USER] uc ON h.MaCustomer = uc.userID
			JOIN [USER] ud ON h.MaBacSi = ud.userID
			JOIN PHONGKHAM p ON h.MaPhongKham = p.MaPhongKham
			WHERE 1=1
		`

		if customerID != "" {
			query += " AND h.MaCustomer = @p1"
			args = append(args, customerID)
		}

		query += " ORDER BY h.NgayKham DESC"
	}

	rows, err := h.db.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to retrieve medical records",
			Error:   err.Error(),
		})
		return
	}
	defer rows.Close()

	var records []map[string]interface{}
	for rows.Next() {
		record := make(map[string]interface{})
		var maHoSo, maCustomer, maBacSi, maPhongKham string
		var ngayKham interface{}
		var trieuChung, chanDoan, huongDanDieuTri, maICD10, ngayTaiKham interface{}
		var tenPerson, tenPhongKham interface{}

		if userType.(string) == "CUSTOMER" || userType.(string) == "DOCTOR" {
			err := rows.Scan(&maHoSo, &maCustomer, &maBacSi, &maPhongKham,
				&ngayKham, &trieuChung, &chanDoan, &huongDanDieuTri, &maICD10, &ngayTaiKham,
				&tenPerson, &tenPhongKham)
			if err != nil {
				c.JSON(http.StatusInternalServerError, models.APIResponse{
					Success: false,
					Message: "Failed to scan medical record data",
					Error:   err.Error(),
				})
				return
			}
		} else {
			var tenKhachHang, tenBacSi interface{}
			err := rows.Scan(&maHoSo, &maCustomer, &maBacSi, &maPhongKham,
				&ngayKham, &trieuChung, &chanDoan, &huongDanDieuTri, &maICD10, &ngayTaiKham,
				&tenKhachHang, &tenBacSi, &tenPhongKham)
			if err != nil {
				c.JSON(http.StatusInternalServerError, models.APIResponse{
					Success: false,
					Message: "Failed to scan medical record data",
					Error:   err.Error(),
				})
				return
			}
			record["ten_khach_hang"] = tenKhachHang
			record["ten_bac_si"] = tenBacSi
		}

		record["ma_ho_so"] = maHoSo
		record["ma_customer"] = maCustomer
		record["ma_bac_si"] = maBacSi
		record["ma_phong_kham"] = maPhongKham
		record["ngay_kham"] = ngayKham
		record["trieu_chung"] = trieuChung
		record["chan_doan"] = chanDoan
		record["huong_dan_dieu_tri"] = huongDanDieuTri
		record["ma_icd10"] = maICD10
		record["ngay_tai_kham"] = ngayTaiKham
		record["ten_phong_kham"] = tenPhongKham

		if userType.(string) == "CUSTOMER" {
			record["ten_bac_si"] = tenPerson
		} else if userType.(string) == "DOCTOR" {
			record["ten_khach_hang"] = tenPerson
		}

		records = append(records, record)
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Medical records retrieved successfully",
		Data:    records,
	})
}

func (h *MedicalRecordHandler) GetMedicalRecord(c *gin.Context) {
	recordID := c.Param("id")
	userID, _ := c.Get("user_id")
	userType, _ := c.Get("user_type")

	query := `
		SELECT h.MaHoSo, h.MaCustomer, h.MaBacSi, h.MaPhongKham, 
		       h.NgayKham, h.TrieuChung, h.ChanDoan, h.HuongDanDieuTri, 
		       h.MaICD10, h.NgayTaiKham,
		       uc.HoTen as TenKhachHang, ud.HoTen as TenBacSi, p.TenPhongKham
		FROM HOSO h
		JOIN [USER] uc ON h.MaCustomer = uc.MaUser
		JOIN [USER] ud ON h.MaBacSi = ud.MaUser
		JOIN PHONGKHAM p ON h.MaPhongKham = p.MaPhongKham
		WHERE h.MaHoSo = @p1
	`
	args := []interface{}{recordID}

	if userType.(string) == "CUSTOMER" {
		query += " AND h.MaCustomer = @p1"
		args = append(args, userID)
	} else if userType.(string) == "DOCTOR" {
		query += " AND h.MaBacSi = @p1"
		args = append(args, userID)
	}

	var record map[string]interface{} = make(map[string]interface{})
	var maHoSo, maCustomer, maBacSi, maPhongKham string
	var ngayKham interface{}
	var trieuChung, chanDoan, huongDanDieuTri, maICD10, ngayTaiKham interface{}
	var tenKhachHang, tenBacSi, tenPhongKham interface{}

	err := h.db.QueryRow(query, args...).Scan(
		&maHoSo, &maCustomer, &maBacSi, &maPhongKham,
		&ngayKham, &trieuChung, &chanDoan, &huongDanDieuTri, &maICD10, &ngayTaiKham,
		&tenKhachHang, &tenBacSi, &tenPhongKham,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, models.APIResponse{
				Success: false,
				Message: "Medical record not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "Failed to retrieve medical record",
				Error:   err.Error(),
			})
		}
		return
	}

	record["ma_ho_so"] = maHoSo
	record["ma_customer"] = maCustomer
	record["ma_bac_si"] = maBacSi
	record["ma_phong_kham"] = maPhongKham
	record["ngay_kham"] = ngayKham
	record["trieu_chung"] = trieuChung
	record["chan_doan"] = chanDoan
	record["huong_dan_dieu_tri"] = huongDanDieuTri
	record["ma_icd10"] = maICD10
	record["ngay_tai_kham"] = ngayTaiKham
	record["ten_khach_hang"] = tenKhachHang
	record["ten_bac_si"] = tenBacSi
	record["ten_phong_kham"] = tenPhongKham

	prescriptionsQuery := `
		SELECT d.MaDonThuoc, d.NgayKeDon, d.GhiChu
		FROM DONTHUOC d
		WHERE d.MaHoSo = @p1
	`

	prescRows, err := h.db.Query(prescriptionsQuery, recordID)
	if err == nil {
		defer prescRows.Close()
		var prescriptions []map[string]interface{}

		for prescRows.Next() {
			var prescription map[string]interface{} = make(map[string]interface{})
			var maDonThuoc, ngayKeDon, ghiChu interface{}

			err := prescRows.Scan(&maDonThuoc, &ngayKeDon, &ghiChu)
			if err == nil {
				prescription["ma_don_thuoc"] = maDonThuoc
				prescription["ngay_ke_don"] = ngayKeDon
				prescription["ghi_chu"] = ghiChu

				medicinesQuery := `
					SELECT t.TenThuoc, ct.SoLuong, ct.CachDung, ct.GhiChu
					FROM CHITIETDONTHUOC ct
					JOIN THUOC t ON ct.MaThuoc = t.MaThuoc
					WHERE ct.MaDonThuoc = @p1
				`

				medRows, medErr := h.db.Query(medicinesQuery, maDonThuoc)
				if medErr == nil {
					defer medRows.Close()
					var medicines []map[string]interface{}

					for medRows.Next() {
						var medicine map[string]interface{} = make(map[string]interface{})
						var tenThuoc, soLuong, cachDung, ghiChuThuoc interface{}

						err := medRows.Scan(&tenThuoc, &soLuong, &cachDung, &ghiChuThuoc)
						if err == nil {
							medicine["ten_thuoc"] = tenThuoc
							medicine["so_luong"] = soLuong
							medicine["cach_dung"] = cachDung
							medicine["ghi_chu"] = ghiChuThuoc
							medicines = append(medicines, medicine)
						}
					}
					prescription["thuoc"] = medicines
				}

				prescriptions = append(prescriptions, prescription)
			}
		}
		record["don_thuoc"] = prescriptions
	}

	testResultsQuery := `
		SELECT MaXetNghiem, LoaiXetNghiem, NgayXetNghiem, KetQua, GhiChu, FileDinhKem
		FROM XETNGHIEM
		WHERE MaHoSo = @p1
	`

	testRows, err := h.db.Query(testResultsQuery, recordID)
	if err == nil {
		defer testRows.Close()
		var testResults []map[string]interface{}

		for testRows.Next() {
			var testResult map[string]interface{} = make(map[string]interface{})
			var maXetNghiem, loaiXetNghiem, ngayXetNghiem, ketQua, ghiChu, fileDinhKem interface{}

			err := testRows.Scan(&maXetNghiem, &loaiXetNghiem, &ngayXetNghiem, &ketQua, &ghiChu, &fileDinhKem)
			if err == nil {
				testResult["ma_xet_nghiem"] = maXetNghiem
				testResult["loai_xet_nghiem"] = loaiXetNghiem
				testResult["ngay_xet_nghiem"] = ngayXetNghiem
				testResult["ket_qua"] = ketQua
				testResult["ghi_chu"] = ghiChu
				testResult["file_dinh_kem"] = fileDinhKem
				testResults = append(testResults, testResult)
			}
		}
		record["ket_qua_xet_nghiem"] = testResults
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Medical record retrieved successfully",
		Data:    record,
	})
}

func (h *MedicalRecordHandler) CreateMedicalRecord(c *gin.Context) {
	userID, _ := c.Get("user_id")
	userType, _ := c.Get("user_type")

	if userType.(string) != "DOCTOR" {
		c.JSON(http.StatusForbidden, models.APIResponse{
			Success: false,
			Message: "Only doctors can create medical records",
		})
		return
	}

	var req struct {
		MaCustomer      string  `json:"ma_customer" binding:"required"`
		MaPhongKham     string  `json:"ma_phong_kham" binding:"required"`
		TrieuChung      *string `json:"trieu_chung"`
		ChanDoan        *string `json:"chan_doan"`
		HuongDanDieuTri *string `json:"huong_dan_dieu_tri"`
		MaICD10         *string `json:"ma_icd10"`
		NgayTaiKham     *string `json:"ngay_tai_kham"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request format",
			Error:   err.Error(),
		})
		return
	}

	recordID := utils.GenerateMedicalRecordID()

	_, err := h.db.Exec(`
		INSERT INTO HOSO (MaHoSo, MaCustomer, MaBacSi, MaPhongKham, NgayKham, 
		                      TrieuChung, ChanDoan, HuongDanDieuTri, MaICD10, NgayTaiKham)
		VALUES (?, ?, ?, ?, GETDATE(), ?, ?, ?, ?, ?)
	`, recordID, req.MaCustomer, userID, req.MaPhongKham,
		req.TrieuChung, req.ChanDoan, req.HuongDanDieuTri, req.MaICD10, req.NgayTaiKham)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to create medical record",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "Medical record created successfully",
		Data: gin.H{
			"record_id": recordID,
		},
	})
}

func (h *MedicalRecordHandler) UpdateMedicalRecord(c *gin.Context) {
	recordID := c.Param("id")
	userID, _ := c.Get("user_id")
	userType, _ := c.Get("user_type")

	if userType.(string) != "DOCTOR" {
		c.JSON(http.StatusForbidden, models.APIResponse{
			Success: false,
			Message: "Only doctors can update medical records",
		})
		return
	}

	var currentDoctorID string
	err := h.db.QueryRow("SELECT MaBacSi FROM HOSO WHERE MaHoSo = @p1", recordID).Scan(&currentDoctorID)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, models.APIResponse{
				Success: false,
				Message: "Medical record not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "Failed to find medical record",
				Error:   err.Error(),
			})
		}
		return
	}

	if currentDoctorID != userID.(string) {
		c.JSON(http.StatusForbidden, models.APIResponse{
			Success: false,
			Message: "You can only update your own medical records",
		})
		return
	}

	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request format",
			Error:   err.Error(),
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

	if trieuChung, exists := updateData["trieu_chung"]; exists {
		_, err = tx.Exec("UPDATE HOSO SET TrieuChung = @p1 WHERE MaHoSo = @p1", trieuChung, recordID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "Failed to update symptoms",
				Error:   err.Error(),
			})
			return
		}
	}

	if chanDoan, exists := updateData["chan_doan"]; exists {
		_, err = tx.Exec("UPDATE HOSO SET ChanDoan = @p1 WHERE MaHoSo = @p1", chanDoan, recordID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "Failed to update diagnosis",
				Error:   err.Error(),
			})
			return
		}
	}

	if huongDanDieuTri, exists := updateData["huong_dan_dieu_tri"]; exists {
		_, err = tx.Exec("UPDATE HOSO SET HuongDanDieuTri = @p1 WHERE MaHoSo = @p1", huongDanDieuTri, recordID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "Failed to update treatment instructions",
				Error:   err.Error(),
			})
			return
		}
	}

	if maICD10, exists := updateData["ma_icd10"]; exists {
		_, err = tx.Exec("UPDATE HOSO SET MaICD10 = @p1 WHERE MaHoSo = @p1", maICD10, recordID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "Failed to update ICD10 code",
				Error:   err.Error(),
			})
			return
		}
	}

	if ngayTaiKham, exists := updateData["ngay_tai_kham"]; exists {
		_, err = tx.Exec("UPDATE HOSO SET NgayTaiKham = @p1 WHERE MaHoSo = @p1", ngayTaiKham, recordID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "Failed to update follow-up date",
				Error:   err.Error(),
			})
			return
		}
	}

	if err = tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to update medical record",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Medical record updated successfully",
	})
}
