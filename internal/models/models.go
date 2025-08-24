package models

import (
	"time"
)

type User struct {
	MaUser           string     `json:"ma_user" db:"MaUser"`
	HoTen            string     `json:"ho_ten" db:"HoTen"`
	SoDienThoai      *string    `json:"so_dien_thoai" db:"SoDienThoai"`
	Email            *string    `json:"email" db:"Email"`
	TenDangNhap      string     `json:"ten_dang_nhap" db:"TenDangNhap"`
	MatKhau          string     `json:"-" db:"MatKhau"`
	TrangThai        string     `json:"trang_thai" db:"TrangThai"`
	NgayTao          time.Time  `json:"ngay_tao" db:"NgayTao"`
	LanDangNhapCuoi  *time.Time `json:"lan_dang_nhap_cuoi" db:"LanDangNhapCuoi"`
}

type Customer struct {
	User
	NgaySinh    *time.Time `json:"ngay_sinh" db:"NgaySinh"`
	GioiTinh    *string    `json:"gioi_tinh" db:"GioiTinh"`
	DiaChi      *string    `json:"dia_chi" db:"DiaChi"`
	NgayDangKy  time.Time  `json:"ngay_dang_ky" db:"NgayDangKy"`
	MaBaoHiem   *string    `json:"ma_bao_hiem" db:"MaBaoHiem"`
}

type Doctor struct {
	User
	ChuyenKhoa           *string `json:"chuyen_khoa" db:"ChuyenKhoa"`
	NamKinhNghiem        *int    `json:"nam_kinh_nghiem" db:"NamKinhNghiem"`
	BangCap              *string `json:"bang_cap" db:"BangCap"`
	SoGiayPhepHanhNghe   *string `json:"so_giay_phep_hanh_nghe" db:"SoGiayPhepHanhNghe"`
}

type Receptionist struct {
	User
	MaPhongKham string  `json:"ma_phong_kham" db:"MaPhongKham"`
	LuongCoBan  float64 `json:"luong_co_ban" db:"LuongCoBan"`
	NgayVaoLam  *time.Time `json:"ngay_vao_lam" db:"NgayVaoLam"`
}

type Accountant struct {
	User
	LuongCoBan  float64 `json:"luong_co_ban" db:"LuongCoBan"`
	NgayVaoLam  *time.Time `json:"ngay_vao_lam" db:"NgayVaoLam"`
	ChuyenMon   *string `json:"chuyen_mon" db:"ChuyenMon"`
}

type ClinicManager struct {
	User
	MaPhongKham string  `json:"ma_phong_kham" db:"MaPhongKham"`
	LuongCoBan  float64 `json:"luong_co_ban" db:"LuongCoBan"`
	NgayVaoLam  *time.Time `json:"ngay_vao_lam" db:"NgayVaoLam"`
}

type OperationManager struct {
	User
	ChucVu         *string `json:"chuc_vu" db:"ChucVu"`
	KhuVucPhuTrach *string `json:"khu_vuc_phu_trach" db:"KhuVucPhuTrach"`
	LuongCoBan     float64 `json:"luong_co_ban" db:"LuongCoBan"`
	NgayVaoLam     *time.Time `json:"ngay_vao_lam" db:"NgayVaoLam"`
}

type Clinic struct {
	MaPhongKham   string  `json:"ma_phong_kham" db:"MaPhongKham"`
	TenPhongKham  string  `json:"ten_phong_kham" db:"TenPhongKham"`
	DiaChi        *string `json:"dia_chi" db:"DiaChi"`
	SoDienThoai   *string `json:"so_dien_thoai" db:"SoDienThoai"`
	Email         *string `json:"email" db:"Email"`
}

type Appointment struct {
	MaLichKham   string    `json:"ma_lich_kham" db:"MaLichKham"`
	MaCustomer   string    `json:"ma_customer" db:"MaCustomer"`
	MaBacSi      string    `json:"ma_bac_si" db:"MaBacSi"`
	MaPhongKham  string    `json:"ma_phong_kham" db:"MaPhongKham"`
	NgayGioKham  time.Time `json:"ngay_gio_kham" db:"NgayGioKham"`
	TrangThai    string    `json:"trang_thai" db:"TrangThai"`
	GhiChu       *string   `json:"ghi_chu" db:"GhiChu"`
	NgayDat      time.Time `json:"ngay_dat" db:"NgayDat"`
}

type WorkSchedule struct {
	MaLichLamViec string    `json:"ma_lich_lam_viec" db:"MaLichLamViec"`
	MaBacSi       string    `json:"ma_bac_si" db:"MaBacSi"`
	MaPhongKham   string    `json:"ma_phong_kham" db:"MaPhongKham"`
	NgayLamViec   time.Time `json:"ngay_lam_viec" db:"NgayLamViec"`
	GioBatDau     string    `json:"gio_bat_dau" db:"GioBatDau"`
	GioKetThuc    string    `json:"gio_ket_thuc" db:"GioKetThuc"`
	TrangThai     string    `json:"trang_thai" db:"TrangThai"`
}

type MedicalRecord struct {
	MaHoSo              string     `json:"ma_ho_so" db:"MaHoSo"`
	MaCustomer          string     `json:"ma_customer" db:"MaCustomer"`
	MaBacSi             string     `json:"ma_bac_si" db:"MaBacSi"`
	MaPhongKham         string     `json:"ma_phong_kham" db:"MaPhongKham"`
	NgayKham            time.Time  `json:"ngay_kham" db:"NgayKham"`
	TrieuChung          *string    `json:"trieu_chung" db:"TrieuChung"`
	ChanDoan            *string    `json:"chan_doan" db:"ChanDoan"`
	HuongDanDieuTri     *string    `json:"huong_dan_dieu_tri" db:"HuongDanDieuTri"`
	MaICD10             *string    `json:"ma_icd10" db:"MaICD10"`
	NgayTaiKham         *time.Time `json:"ngay_tai_kham" db:"NgayTaiKham"`
}

type Medicine struct {
	MaThuoc   string  `json:"ma_thuoc" db:"MaThuoc"`
	TenThuoc  string  `json:"ten_thuoc" db:"TenThuoc"`
	DonVi     *string `json:"don_vi" db:"DonVi"`
	Gia       float64 `json:"gia" db:"Gia"`
	CongDung  *string `json:"cong_dung" db:"CongDung"`
	LieuLuong *string `json:"lieu_luong" db:"LieuLuong"`
}

type Prescription struct {
	MaDonThuoc string     `json:"ma_don_thuoc" db:"MaDonThuoc"`
	MaHoSo     string     `json:"ma_ho_so" db:"MaHoSo"`
	NgayKeDon  time.Time  `json:"ngay_ke_don" db:"NgayKeDon"`
	GhiChu     *string    `json:"ghi_chu" db:"GhiChu"`
}

type PrescriptionDetail struct {
	MaDonThuoc string  `json:"ma_don_thuoc" db:"MaDonThuoc"`
	MaThuoc    string  `json:"ma_thuoc" db:"MaThuoc"`
	SoLuong    int     `json:"so_luong" db:"SoLuong"`
	CachDung   *string `json:"cach_dung" db:"CachDung"`
	GhiChu     *string `json:"ghi_chu" db:"GhiChu"`
}

type TestResult struct {
	MaXetNghiem    string     `json:"ma_xet_nghiem" db:"MaXetNghiem"`
	MaHoSo         string     `json:"ma_ho_so" db:"MaHoSo"`
	LoaiXetNghiem  *string    `json:"loai_xet_nghiem" db:"LoaiXetNghiem"`
	NgayXetNghiem  *time.Time `json:"ngay_xet_nghiem" db:"NgayXetNghiem"`
	KetQua         *string    `json:"ket_qua" db:"KetQua"`
	GhiChu         *string    `json:"ghi_chu" db:"GhiChu"`
	FileDinhKem    *string    `json:"file_dinh_kem" db:"FileDinhKem"`
}

type MedicalImage struct {
	MaHinhAnh     string    `json:"ma_hinh_anh" db:"MaHinhAnh"`
	MaHoSo        string    `json:"ma_ho_so" db:"MaHoSo"`
	TenFile       *string   `json:"ten_file" db:"TenFile"`
	DuongDanFile  *string   `json:"duong_dan_file" db:"DuongDanFile"`
	MoTa          *string   `json:"mo_ta" db:"MoTa"`
	NgayChup      time.Time `json:"ngay_chup" db:"NgayChup"`
}

type Payment struct {
	MaThanhToan           string    `json:"ma_thanh_toan" db:"MaThanhToan"`
	MaLichKham            string    `json:"ma_lich_kham" db:"MaLichKham"`
	TongTien              float64   `json:"tong_tien" db:"TongTien"`
	NgayThanhToan         time.Time `json:"ngay_thanh_toan" db:"NgayThanhToan"`
	PhuongThucThanhToan   *string   `json:"phuong_thuc_thanh_toan" db:"PhuongThucThanhToan"`
	TrangThai             string    `json:"trang_thai" db:"TrangThai"`
}

type Salary struct {
	MaLuong        string    `json:"ma_luong" db:"MaLuong"`
	MaUser         string    `json:"ma_user" db:"MaUser"`
	Thang          int       `json:"thang" db:"Thang"`
	Nam            int       `json:"nam" db:"Nam"`
	LuongCoBan     float64   `json:"luong_co_ban" db:"LuongCoBan"`
	ThuLao         float64   `json:"thu_lao" db:"ThuLao"`
	Thuong         float64   `json:"thuong" db:"Thuong"`
	TongLuong      *float64  `json:"tong_luong" db:"TongLuong"`
	NgayTinhLuong  time.Time `json:"ngay_tinh_luong" db:"NgayTinhLuong"`
}

type Report struct {
	MaBaoCao      string     `json:"ma_bao_cao" db:"MaBaoCao"`
	MaUser        string     `json:"ma_user" db:"MaUser"`
	LoaiBaoCao    *string    `json:"loai_bao_cao" db:"LoaiBaoCao"`
	TuNgay        *time.Time `json:"tu_ngay" db:"TuNgay"`
	DenNgay       *time.Time `json:"den_ngay" db:"DenNgay"`
	NoiDung       *string    `json:"noi_dung" db:"NoiDung"`
	NgayTaoBaoCao time.Time  `json:"ngay_tao_bao_cao" db:"NgayTaoBaoCao"`
}

type AuthRequest struct {
	TenDangNhap string `json:"ten_dang_nhap" binding:"required"`
	MatKhau     string `json:"mat_khau" binding:"required"`
}

type RegisterRequest struct {
	HoTen       string `json:"ho_ten" binding:"required"`
	TenDangNhap string `json:"ten_dang_nhap" binding:"required"`
	MatKhau     string `json:"mat_khau" binding:"required"`
	Email       string `json:"email"`
	SoDienThoai string `json:"so_dien_thoai"`
}

type AuthResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	User         User   `json:"user"`
}

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required"`
}

type ResetPasswordRequest struct {
	Email       string `json:"email" binding:"required"`
	ResetCode   string `json:"reset_code" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

type PasswordReset struct {
	ID          string    `json:"id" db:"ID"`
	UserID      string    `json:"user_id" db:"UserID"`
	Email       string    `json:"email" db:"Email"`
	ResetCode   string    `json:"reset_code" db:"ResetCode"`
	IsUsed      bool      `json:"is_used" db:"IsUsed"`
	ExpiresAt   time.Time `json:"expires_at" db:"ExpiresAt"`
	CreatedAt   time.Time `json:"created_at" db:"CreatedAt"`
}