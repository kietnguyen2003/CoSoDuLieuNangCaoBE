package handlers

import (
	"database/sql"
	"net/http"

	"clinic-management/internal/models"
	"clinic-management/internal/utils"

	"github.com/gin-gonic/gin"
)

type CustomerHandler struct {
	db *sql.DB
}

func NewCustomerHandler(db *sql.DB) *CustomerHandler {
	return &CustomerHandler{db: db}
}

// GetCustomers - Get all customers for receptionist
func (h *CustomerHandler) GetCustomers(c *gin.Context) {
	userRole, _ := c.Get("user_type")
	
	// Only allow receptionist and admin to view customers
	if userRole != "RECEPTIONIST" && userRole != "ADMIN" {
		c.JSON(http.StatusForbidden, models.APIResponse{
			Success: false,
			Message: "Access denied. Only receptionists can view customers",
		})
		return
	}

	search := c.Query("search")
	var query string
	var args []interface{}

	query = `
		SELECT u.userID, u.hoTen, u.soDienThoai, u.email, u.status,
		       c.ngaySinh, c.gioiTinh, c.diaChi, c.createdAt, c.maBaoHiem
		FROM [USER] u
		JOIN CUSTOMER c ON u.userID = c.maUser
		WHERE u.role = 'CUSTOMER' AND u.status = 'ACTIVE'
	`

	if search != "" {
		query += " AND (u.hoTen LIKE @p1 OR u.userID LIKE @p1 OR u.soDienThoai LIKE @p1)"
		args = append(args, "%"+search+"%")
	}

	query += " ORDER BY c.createdAt DESC"

	rows, err := h.db.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to retrieve customers",
			Error:   err.Error(),
		})
		return
	}
	defer rows.Close()

	var customers []map[string]interface{}
	for rows.Next() {
		var customer map[string]interface{} = make(map[string]interface{})
		var userID, hoTen, soDienThoai, email, status string
		var ngaySinh, gioiTinh, diaChi, createdAt, maBaoHiem sql.NullString

		err := rows.Scan(&userID, &hoTen, &soDienThoai, &email, &status,
			&ngaySinh, &gioiTinh, &diaChi, &createdAt, &maBaoHiem)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "Failed to scan customer data",
				Error:   err.Error(),
			})
			return
		}

		customer["user_id"] = userID
		customer["ho_ten"] = hoTen
		customer["so_dien_thoai"] = soDienThoai
		customer["email"] = email
		customer["status"] = status
		if ngaySinh.Valid {
			customer["ngay_sinh"] = ngaySinh.String
		}
		if gioiTinh.Valid {
			customer["gioi_tinh"] = gioiTinh.String
		}
		if diaChi.Valid {
			customer["dia_chi"] = diaChi.String
		}
		if createdAt.Valid {
			customer["ngay_dang_ky"] = createdAt.String
		}
		if maBaoHiem.Valid {
			customer["ma_bao_hiem"] = maBaoHiem.String
		}

		customers = append(customers, customer)
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Customers retrieved successfully",
		Data:    customers,
	})
}

// GetCustomer - Get customer details by ID
func (h *CustomerHandler) GetCustomer(c *gin.Context) {
	userRole, _ := c.Get("user_type")
	
	// Only allow receptionist and admin to view customers
	if userRole != "RECEPTIONIST" && userRole != "ADMIN" {
		c.JSON(http.StatusForbidden, models.APIResponse{
			Success: false,
			Message: "Access denied. Only receptionists can view customers",
		})
		return
	}

	customerID := c.Param("id")

	query := `
		SELECT u.userID, u.hoTen, u.soDienThoai, u.email, u.status,
		       c.ngaySinh, c.gioiTinh, c.diaChi, c.createdAt, c.maBaoHiem
		FROM [USER] u
		JOIN CUSTOMER c ON u.userID = c.maUser
		WHERE u.userID = @p1
	`

	var customer map[string]interface{} = make(map[string]interface{})
	var userID, hoTen, soDienThoai, email, status string
	var ngaySinh, gioiTinh, diaChi, createdAt, maBaoHiem sql.NullString

	err := h.db.QueryRow(query, customerID).Scan(
		&userID, &hoTen, &soDienThoai, &email, &status,
		&ngaySinh, &gioiTinh, &diaChi, &createdAt, &maBaoHiem,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, models.APIResponse{
				Success: false,
				Message: "Customer not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "Failed to retrieve customer",
				Error:   err.Error(),
			})
		}
		return
	}

	customer["user_id"] = userID
	customer["ho_ten"] = hoTen
	customer["so_dien_thoai"] = soDienThoai
	customer["email"] = email
	customer["status"] = status
	if ngaySinh.Valid {
		customer["ngay_sinh"] = ngaySinh.String
	}
	if gioiTinh.Valid {
		customer["gioi_tinh"] = gioiTinh.String
	}
	if diaChi.Valid {
		customer["dia_chi"] = diaChi.String
	}
	if createdAt.Valid {
		customer["ngay_dang_ky"] = createdAt.String
	}
	if maBaoHiem.Valid {
		customer["ma_bao_hiem"] = maBaoHiem.String
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Customer retrieved successfully",
		Data:    customer,
	})
}

// CreateCustomer - Create new customer (receptionist only)
func (h *CustomerHandler) CreateCustomer(c *gin.Context) {
	userRole, _ := c.Get("user_type")
	
	// Only allow receptionist and admin to create customers
	if userRole != "RECEPTIONIST" && userRole != "ADMIN" {
		c.JSON(http.StatusForbidden, models.APIResponse{
			Success: false,
			Message: "Access denied. Only receptionists can create customers",
		})
		return
	}

	var req struct {
		HoTen        string `json:"ho_ten" binding:"required"`
		TenDangNhap  string `json:"ten_dang_nhap" binding:"required"`
		MatKhau      string `json:"mat_khau" binding:"required"`
		SoDienThoai  string `json:"so_dien_thoai"`
		Email        string `json:"email"`
		NgaySinh     string `json:"ngay_sinh"`
		GioiTinh     string `json:"gioi_tinh"`
		DiaChi       string `json:"dia_chi"`
		MaBaoHiem    string `json:"ma_bao_hiem"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request data",
			Error:   err.Error(),
		})
		return
	}

	// Validate required fields
	if req.HoTen == "" || req.TenDangNhap == "" || req.MatKhau == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Name, username, and password are required",
		})
		return
	}

	// Check if username already exists
	var existingUser string
	err := h.db.QueryRow("SELECT userID FROM [USER] WHERE username = @p1", req.TenDangNhap).Scan(&existingUser)
	if err != sql.ErrNoRows {
		c.JSON(http.StatusConflict, models.APIResponse{
			Success: false,
			Message: "Username already exists",
		})
		return
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.MatKhau)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to hash password",
			Error:   err.Error(),
		})
		return
	}

	// Generate IDs
	userID := utils.GenerateCustomerID()
	
	// Start transaction
	tx, err := h.db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to start transaction",
			Error:   err.Error(),
		})
		return
	}
	defer tx.Rollback()

	// Insert into USER table
	_, err = tx.Exec(`
		INSERT INTO [USER] (userID, hoTen, username, password, soDienThoai, email, role, status, createdAt)
		VALUES (@p1, @p2, @p3, @p4, @p5, @p6, 'CUSTOMER', 'ACTIVE', GETDATE())
	`, userID, req.HoTen, req.TenDangNhap, hashedPassword, req.SoDienThoai, req.Email)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to create user",
			Error:   err.Error(),
		})
		return
	}

	// Insert into CUSTOMER table
	_, err = tx.Exec(`
		INSERT INTO CUSTOMER (maUser, ngaySinh, gioiTinh, diaChi, maBaoHiem)
		VALUES (@p1, @p2, @p3, @p4, @p5)
	`, userID, req.NgaySinh, req.GioiTinh, req.DiaChi, req.MaBaoHiem)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to create customer",
			Error:   err.Error(),
		})
		return
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to commit transaction",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "Customer created successfully",
		Data: map[string]interface{}{
			"user_id": userID,
			"ho_ten":  req.HoTen,
		},
	})
}