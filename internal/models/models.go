package models

import (
	"time"
)

type User struct {
	MaUser      string    `json:"ma_user" db:"userID"`
	HoTen       string    `json:"ho_ten" db:"hoTen"`
	SoDienThoai *string   `json:"so_dien_thoai" db:"soDienThoai"`
	Email       *string   `json:"email" db:"email"`
	TenDangNhap string    `json:"ten_dang_nhap" db:"username"`
	MatKhau     string    `json:"-" db:"password"`
	TrangThai   string    `json:"trang_thai" db:"status"`
	NgayTao     time.Time `json:"ngay_tao" db:"createdAt"`
	Role        string    `json:"role" db:"role"`
}

type Customer struct {
	User
	NgaySinh   *time.Time `json:"ngay_sinh" db:"ngaySinh"`
	GioiTinh   *string    `json:"gioi_tinh" db:"gioiTinh"`
	DiaChi     *string    `json:"dia_chi" db:"diaChi"`
	NgayDangKy time.Time  `json:"ngay_dang_ky" db:"createdAt"`
	MaBaoHiem  *string    `json:"ma_bao_hiem" db:"maBaoHiem"`
}

type Doctor struct {
	User
	ChuyenKhoa         *string `json:"chuyen_khoa" db:"chuyenKhoa"`
	NamKinhNghiem      *int    `json:"nam_kinh_nghiem" db:"namKinhNghiem"`
	BangCap            *string `json:"bang_cap" db:"bangCap"`
	SoGiayPhepHanhNghe *string `json:"so_giay_phep_hanh_nghe" db:"maGiayPhep"`
}

type Receptionist struct {
	User
	MaPhongKham string     `json:"ma_phong_kham" db:"maPhongKham"`
	LuongCoBan  float64    `json:"luong_co_ban" db:"luongCoBan"`
	NgayVaoLam  *time.Time `json:"ngay_vao_lam" db:"ngayVaoLam"`
}

type Accountant struct {
	User
	LuongCoBan float64    `json:"luong_co_ban" db:"luongCoBan"`
	NgayVaoLam *time.Time `json:"ngay_vao_lam" db:"ngayVaoLam"`
	ChuyenMon  *string    `json:"chuyen_mon" db:"chuyenMon"`
}

type ClinicManager struct {
	User
	MaPhongKham string     `json:"ma_phong_kham" db:"maPhongKham"`
	LuongCoBan  float64    `json:"luong_co_ban" db:"luongCoBan"`
	NgayVaoLam  *time.Time `json:"ngay_vao_lam" db:"ngayVaoLam"`
}

type OperationManager struct {
	User
	ChucVu         *string    `json:"chuc_vu" db:"chucVu"`
	KhuVucPhuTrach *string    `json:"khu_vuc_phu_trach" db:"khuVucPhuTrach"`
	LuongCoBan     float64    `json:"luong_co_ban" db:"luongCoBan"`
	NgayVaoLam     *time.Time `json:"ngay_vao_lam" db:"ngayVaoLam"`
}

type Clinic struct {
	MaPhongKham  string  `json:"ma_phong_kham" db:"maPhongKham"`
	TenPhongKham string  `json:"ten_phong_kham" db:"tenPhongKham"`
	DiaChi       *string `json:"dia_chi" db:"diaChi"`
	SoDienThoai  *string `json:"so_dien_thoai" db:"soDienThoai"`
	Email        *string `json:"email" db:"email"`
}

type Appointment struct {
	MaLichKham  string    `json:"ma_lich_kham" db:"maLichKham"`
	MaCustomer  string    `json:"ma_customer" db:"maCustomer"`
	MaBacSi     string    `json:"ma_bac_si" db:"maBacSi"`
	MaPhongKham string    `json:"ma_phong_kham" db:"maPhongKham"`
	NgayGioKham time.Time `json:"ngay_gio_kham" db:"maGioKham"`
	TrangThai   string    `json:"trang_thai" db:"trangThai"`
	GhiChu      *string   `json:"ghi_chu" db:"ghiChu"`
	NgayDat     time.Time `json:"ngay_dat" db:"createdAt"`
}

type WorkSchedule struct {
	MaLichLamViec string    `json:"ma_lich_lam_viec" db:"maLichLamViec"`
	MaBacSi       string    `json:"ma_bac_si" db:"maBacSi"`
	MaPhongKham   string    `json:"ma_phong_kham" db:"maPhongKham"`
	NgayLamViec   time.Time `json:"ngay_lam_viec" db:"ngayLamViec"`
	GioBatDau     string    `json:"gio_bat_dau" db:"gioBatDau"`
	GioKetThuc    string    `json:"gio_ket_thuc" db:"gioKetThuc"`
	TrangThai     string    `json:"trang_thai" db:"trangThai"`
}

type MedicalRecord struct {
	MaHoSo          string     `json:"ma_ho_so" db:"maHoSo"`
	MaCustomer      string     `json:"ma_customer" db:"maCustomer"`
	MaBacSi         string     `json:"ma_bac_si" db:"maBacSi"`
	MaPhongKham     string     `json:"ma_phong_kham" db:"maPhongKham"`
	NgayKham        time.Time  `json:"ngay_kham" db:"ngayKham"`
	TrieuChung      *string    `json:"trieu_chung" db:"trieuChung"`
	ChanDoan        *string    `json:"chan_doan" db:"ChanDoan"`
	HuongDanDieuTri *string    `json:"huong_dan_dieu_tri" db:"huongDanDieuTri"`
	MaICD10         *string    `json:"ma_icd10" db:"maICD10"`
	NgayTaiKham     *time.Time `json:"ngay_tai_kham" db:"ngayTaiKham"`
}

type Medicine struct {
	MaThuoc   string  `json:"ma_thuoc" db:"maThuoc"`
	TenThuoc  string  `json:"ten_thuoc" db:"tenThuoc"`
	DonVi     *string `json:"don_vi" db:"donVi"`
	Gia       float64 `json:"gia" db:"Gia"`
	CongDung  *string `json:"cong_dung" db:"congDung"`
	LieuLuong *string `json:"lieu_luong" db:"lieuLuong"`
}

type Prescription struct {
	MaDonThuoc string    `json:"ma_don_thuoc" db:"maDonThuoc"`
	MaHoSo     string    `json:"ma_ho_so" db:"maHoSo"`
	NgayKeDon  time.Time `json:"ngay_ke_don" db:"ngayKeDon"`
	GhiChu     *string   `json:"ghi_chu" db:"ghiChu"`
}

type PrescriptionDetail struct {
	MaDonThuoc string  `json:"ma_don_thuoc" db:"maDonThuoc"`
	MaThuoc    string  `json:"ma_thuoc" db:"maThuoc"`
	SoLuong    int     `json:"so_luong" db:"soLuong"`
	CachDung   *string `json:"cach_dung" db:"cachDung"`
	GhiChu     *string `json:"ghi_chu" db:"ghiChu"`
}

type TestResult struct {
	MaXetNghiem   string     `json:"ma_xet_nghiem" db:"maXetNghiem"`
	MaHoSo        string     `json:"ma_ho_so" db:"maHoSo"`
	LoaiXetNghiem *string    `json:"loai_xet_nghiem" db:"loaiXetNghiem"`
	NgayXetNghiem *time.Time `json:"ngay_xet_nghiem" db:"ngayXetNghiem"`
	KetQua        *string    `json:"ket_qua" db:"ketQua"`
	GhiChu        *string    `json:"ghi_chu" db:"ghiChu"`
	FileDinhKem   *string    `json:"file_dinh_kem" db:"FileDinhKem"`
}

type MedicalImage struct {
	MaHinhAnh    string    `json:"ma_hinh_anh" db:"maHinhAnh"`
	MaHoSo       string    `json:"ma_ho_so" db:"maHoSo"`
	TenFile      *string   `json:"ten_file" db:"tenFile"`
	DuongDanFile *string   `json:"duong_dan_file" db:"DuongDanFile"`
	MoTa         *string   `json:"mo_ta" db:"MoTa"`
	NgayChup     time.Time `json:"ngay_chup" db:"ngayChup"`
}

type Payment struct {
	MaThanhToan         string    `json:"ma_thanh_toan" db:"maThanhToan"`
	MaLichKham          string    `json:"ma_lich_kham" db:"maLichKham"`
	TongTien            float64   `json:"tong_tien" db:"tongTien"`
	NgayThanhToan       time.Time `json:"ngay_thanh_toan" db:"ngayThanhToan"`
	PhuongThucThanhToan *string   `json:"phuong_thuc_thanh_toan" db:"phuongThucThanhToan"`
	TrangThai           string    `json:"trang_thai" db:"trangThai"`
}

type Salary struct {
	MaLuong       string    `json:"ma_luong" db:"maLuong"`
	MaUser        string    `json:"ma_user" db:"maUser"`
	Thang         int       `json:"thang" db:"thang"`
	Nam           int       `json:"nam" db:"nam"`
	LuongCoBan    float64   `json:"luong_co_ban" db:"luongCoBan"`
	ThuLao        float64   `json:"thu_lao" db:"thuLao"`
	Thuong        float64   `json:"thuong" db:"thuong"`
	TongLuong     *float64  `json:"tong_luong" db:"tongLuong"`
	NgayTinhLuong time.Time `json:"ngay_tinh_luong" db:"ngayTinhLuong"`
}

type Report struct {
	MaBaoCao      string     `json:"ma_bao_cao" db:"maBaoCao"`
	MaUser        string     `json:"ma_user" db:"maUser"`
	LoaiBaoCao    *string    `json:"loai_bao_cao" db:"loaiBaoCao"`
	TuNgay        *time.Time `json:"tu_ngay" db:"tuNgay"`
	DenNgay       *time.Time `json:"den_ngay" db:"DenNgay"`
	NoiDung       *string    `json:"noi_dung" db:"noiDung"`
	NgayTaoBaoCao time.Time  `json:"ngay_tao_bao_cao" db:"ngayTaoBaoCao"`
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
	ID        string    `json:"id" db:"ID"`
	UserID    string    `json:"user_id" db:"UserID"`
	Email     string    `json:"email" db:"email"`
	ResetCode string    `json:"reset_code" db:"ResetCode"`
	IsUsed    bool      `json:"is_used" db:"IsUsed"`
	ExpiresAt time.Time `json:"expires_at" db:"ExpiresAt"`
	CreatedAt time.Time `json:"created_at" db:"CreatedAt"`
}
