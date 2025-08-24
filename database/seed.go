package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
)

type User struct {
	MaUser             string  `json:"MaUser"`
	HoTen              string  `json:"HoTen"`
	SoDienThoai        *string `json:"SoDienThoai"`
	Email              *string `json:"Email"`
	TenDangNhap        string  `json:"TenDangNhap"`
	MatKhau            string  `json:"MatKhau"`
	Role               string  `json:"Role"`
	TrangThai          string  `json:"TrangThai"`
	NgayTao            string  `json:"NgayTao"`
	LanDangNhapCuoi    *string `json:"LanDangNhapCuoi"`
}

type Customer struct {
	MaUser      string  `json:"MaUser"`
	NgaySinh    *string `json:"NgaySinh"`
	GioiTinh    *string `json:"GioiTinh"`
	DiaChi      *string `json:"DiaChi"`
	NgayDangKy  string  `json:"NgayDangKy"`
	MaBaoHiem   *string `json:"MaBaoHiem"`
}

type Doctor struct {
	MaUser              string  `json:"MaUser"`
	ChuyenKhoa          *string `json:"ChuyenKhoa"`
	NamKinhNghiem       *int    `json:"NamKinhNghiem"`
	BangCap             *string `json:"BangCap"`
	SoGiayPhepHanhNghe  *string `json:"SoGiayPhepHanhNghe"`
}

type Clinic struct {
	MaPhongKham   string  `json:"MaPhongKham"`
	TenPhongKham  string  `json:"TenPhongKham"`
	DiaChi        *string `json:"DiaChi"`
	SoDienThoai   *string `json:"SoDienThoai"`
	Email         *string `json:"Email"`
}

type Receptionist struct {
	MaUser       string  `json:"MaUser"`
	MaPhongKham  *string `json:"MaPhongKham"`
	LuongCoBan   *float64 `json:"LuongCoBan"`
	NgayVaoLam   *string `json:"NgayVaoLam"`
}

type Accountant struct {
	MaUser       string  `json:"MaUser"`
	LuongCoBan   *float64 `json:"LuongCoBan"`
	NgayVaoLam   *string `json:"NgayVaoLam"`
	ChuyenMon    *string `json:"ChuyenMon"`
}

type ClinicManager struct {
	MaUser       string  `json:"MaUser"`
	MaPhongKham  *string `json:"MaPhongKham"`
	LuongCoBan   *float64 `json:"LuongCoBan"`
	NgayVaoLam   *string `json:"NgayVaoLam"`
}

type OperationManager struct {
	MaUser         string  `json:"MaUser"`
	ChucVu         *string `json:"ChucVu"`
	KhuVucPhuTrach *string `json:"KhuVucPhuTrach"`
	LuongCoBan     *float64 `json:"LuongCoBan"`
	NgayVaoLam     *string `json:"NgayVaoLam"`
}

type Medicine struct {
	MaThuoc   string   `json:"MaThuoc"`
	TenThuoc  string   `json:"TenThuoc"`
	DonVi     *string  `json:"DonVi"`
	Gia       *float64 `json:"Gia"`
	CongDung  *string  `json:"CongDung"`
	LieuLuong *string  `json:"LieuLuong"`
}

type Appointment struct {
	MaLichKham   string  `json:"MaLichKham"`
	MaCustomer   string  `json:"MaCustomer"`
	MaBacSi      string  `json:"MaBacSi"`
	MaPhongKham  string  `json:"MaPhongKham"`
	NgayGioKham  string  `json:"NgayGioKham"`
	TrangThai    string  `json:"TrangThai"`
	GhiChu       *string `json:"GhiChu"`
	NgayDat      string  `json:"NgayDat"`
}

type WorkSchedule struct {
	MaLichLamViec string `json:"MaLichLamViec"`
	MaBacSi       string `json:"MaBacSi"`
	MaPhongKham   string `json:"MaPhongKham"`
	NgayLamViec   string `json:"NgayLamViec"`
	GioBatDau     string `json:"GioBatDau"`
	GioKetThuc    string `json:"GioKetThuc"`
	TrangThai     string `json:"TrangThai"`
}

type MedicalRecord struct {
	MaHoSo           string  `json:"MaHoSo"`
	MaCustomer       string  `json:"MaCustomer"`
	MaBacSi          string  `json:"MaBacSi"`
	MaPhongKham      string  `json:"MaPhongKham"`
	NgayKham         string  `json:"NgayKham"`
	TrieuChung       *string `json:"TrieuChung"`
	ChanDoan         *string `json:"ChanDoan"`
	HuongDanDieuTri  *string `json:"HuongDanDieuTri"`
	MaICD10          *string `json:"MaICD10"`
	NgayTaiKham      *string `json:"NgayTaiKham"`
}

func SeedDatabase(db *sql.DB) error {
	log.Println("Starting database seeding...")
	
	// Seed in order based on foreign key dependencies
	if err := seedClinics(db); err != nil {
		return fmt.Errorf("failed to seed clinics: %w", err)
	}
	
	if err := seedUsers(db); err != nil {
		return fmt.Errorf("failed to seed users: %w", err)
	}
	
	if err := seedCustomers(db); err != nil {
		return fmt.Errorf("failed to seed customers: %w", err)
	}
	
	if err := seedDoctors(db); err != nil {
		return fmt.Errorf("failed to seed doctors: %w", err)
	}
	
	if err := seedReceptionists(db); err != nil {
		return fmt.Errorf("failed to seed receptionists: %w", err)
	}
	
	if err := seedAccountants(db); err != nil {
		return fmt.Errorf("failed to seed accountants: %w", err)
	}
	
	if err := seedClinicManagers(db); err != nil {
		return fmt.Errorf("failed to seed clinic managers: %w", err)
	}
	
	if err := seedOperationManagers(db); err != nil {
		return fmt.Errorf("failed to seed operation managers: %w", err)
	}
	
	if err := seedMedicines(db); err != nil {
		return fmt.Errorf("failed to seed medicines: %w", err)
	}
	
	if err := seedAppointments(db); err != nil {
		return fmt.Errorf("failed to seed appointments: %w", err)
	}
	
	if err := seedWorkSchedules(db); err != nil {
		return fmt.Errorf("failed to seed work schedules: %w", err)
	}
	
	if err := seedMedicalRecords(db); err != nil {
		return fmt.Errorf("failed to seed medical records: %w", err)
	}
	
	log.Println("Database seeding completed successfully!")
	return nil
}

func loadJSONData(filename string, data interface{}) error {
	filePath := filepath.Join("database", "mockdata", filename)
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", filename, err)
	}
	
	if err := json.Unmarshal(content, data); err != nil {
		return fmt.Errorf("failed to unmarshal JSON from %s: %w", filename, err)
	}
	
	return nil
}

func seedClinics(db *sql.DB) error {
	var clinics []Clinic
	if err := loadJSONData("clinics.json", &clinics); err != nil {
		return err
	}
	
	for _, clinic := range clinics {
		query := `INSERT INTO PHONGKHAM (MaPhongKham, TenPhongKham, DiaChi, SoDienThoai, Email) 
				  VALUES (?, ?, ?, ?, ?)`
		_, err := db.Exec(query, clinic.MaPhongKham, clinic.TenPhongKham, clinic.DiaChi, clinic.SoDienThoai, clinic.Email)
		if err != nil {
			return fmt.Errorf("failed to insert clinic %s: %w", clinic.MaPhongKham, err)
		}
	}
	
	log.Printf("Seeded %d clinics", len(clinics))
	return nil
}

func seedUsers(db *sql.DB) error {
	var users []User
	if err := loadJSONData("users.json", &users); err != nil {
		return err
	}
	
	for _, user := range users {
		query := `INSERT INTO [USER] (MaUser, HoTen, SoDienThoai, Email, TenDangNhap, MatKhau, Role, TrangThai, NgayTao, LanDangNhapCuoi) 
				  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
		_, err := db.Exec(query, user.MaUser, user.HoTen, user.SoDienThoai, user.Email, user.TenDangNhap, user.MatKhau, user.Role, user.TrangThai, user.NgayTao, user.LanDangNhapCuoi)
		if err != nil {
			return fmt.Errorf("failed to insert user %s: %w", user.MaUser, err)
		}
	}
	
	log.Printf("Seeded %d users", len(users))
	return nil
}

func seedCustomers(db *sql.DB) error {
	var customers []Customer
	if err := loadJSONData("customers.json", &customers); err != nil {
		return err
	}
	
	for _, customer := range customers {
		query := `INSERT INTO CUSTOMER (MaUser, NgaySinh, GioiTinh, DiaChi, NgayDangKy, MaBaoHiem) 
				  VALUES (?, ?, ?, ?, ?, ?)`
		_, err := db.Exec(query, customer.MaUser, customer.NgaySinh, customer.GioiTinh, customer.DiaChi, customer.NgayDangKy, customer.MaBaoHiem)
		if err != nil {
			return fmt.Errorf("failed to insert customer %s: %w", customer.MaUser, err)
		}
	}
	
	log.Printf("Seeded %d customers", len(customers))
	return nil
}

func seedDoctors(db *sql.DB) error {
	var doctors []Doctor
	if err := loadJSONData("doctors.json", &doctors); err != nil {
		return err
	}
	
	for _, doctor := range doctors {
		query := `INSERT INTO BACSI (MaUser, ChuyenKhoa, NamKinhNghiem, BangCap, SoGiayPhepHanhNghe) 
				  VALUES (?, ?, ?, ?, ?)`
		_, err := db.Exec(query, doctor.MaUser, doctor.ChuyenKhoa, doctor.NamKinhNghiem, doctor.BangCap, doctor.SoGiayPhepHanhNghe)
		if err != nil {
			return fmt.Errorf("failed to insert doctor %s: %w", doctor.MaUser, err)
		}
	}
	
	log.Printf("Seeded %d doctors", len(doctors))
	return nil
}

func seedReceptionists(db *sql.DB) error {
	var receptionists []Receptionist
	if err := loadJSONData("receptionists.json", &receptionists); err != nil {
		return err
	}
	
	for _, receptionist := range receptionists {
		query := `INSERT INTO LETAN (MaUser, MaPhongKham, LuongCoBan, NgayVaoLam) 
				  VALUES (?, ?, ?, ?)`
		_, err := db.Exec(query, receptionist.MaUser, receptionist.MaPhongKham, receptionist.LuongCoBan, receptionist.NgayVaoLam)
		if err != nil {
			return fmt.Errorf("failed to insert receptionist %s: %w", receptionist.MaUser, err)
		}
	}
	
	log.Printf("Seeded %d receptionists", len(receptionists))
	return nil
}

func seedAccountants(db *sql.DB) error {
	var accountants []Accountant
	if err := loadJSONData("accountants.json", &accountants); err != nil {
		return err
	}
	
	for _, accountant := range accountants {
		query := `INSERT INTO KETOAN (MaUser, LuongCoBan, NgayVaoLam, ChuyenMon) 
				  VALUES (?, ?, ?, ?)`
		_, err := db.Exec(query, accountant.MaUser, accountant.LuongCoBan, accountant.NgayVaoLam, accountant.ChuyenMon)
		if err != nil {
			return fmt.Errorf("failed to insert accountant %s: %w", accountant.MaUser, err)
		}
	}
	
	log.Printf("Seeded %d accountants", len(accountants))
	return nil
}

func seedClinicManagers(db *sql.DB) error {
	var managers []ClinicManager
	if err := loadJSONData("clinic_managers.json", &managers); err != nil {
		return err
	}
	
	for _, manager := range managers {
		query := `INSERT INTO QUANLYPHONGKHAM (MaUser, MaPhongKham, LuongCoBan, NgayVaoLam) 
				  VALUES (?, ?, ?, ?)`
		_, err := db.Exec(query, manager.MaUser, manager.MaPhongKham, manager.LuongCoBan, manager.NgayVaoLam)
		if err != nil {
			return fmt.Errorf("failed to insert clinic manager %s: %w", manager.MaUser, err)
		}
	}
	
	log.Printf("Seeded %d clinic managers", len(managers))
	return nil
}

func seedOperationManagers(db *sql.DB) error {
	var managers []OperationManager
	if err := loadJSONData("operation_managers.json", &managers); err != nil {
		return err
	}
	
	for _, manager := range managers {
		query := `INSERT INTO BANDIEUHANH (MaUser, ChucVu, KhuVucPhuTrach, LuongCoBan, NgayVaoLam) 
				  VALUES (?, ?, ?, ?, ?)`
		_, err := db.Exec(query, manager.MaUser, manager.ChucVu, manager.KhuVucPhuTrach, manager.LuongCoBan, manager.NgayVaoLam)
		if err != nil {
			return fmt.Errorf("failed to insert operation manager %s: %w", manager.MaUser, err)
		}
	}
	
	log.Printf("Seeded %d operation managers", len(managers))
	return nil
}

func seedMedicines(db *sql.DB) error {
	var medicines []Medicine
	if err := loadJSONData("medicines.json", &medicines); err != nil {
		return err
	}
	
	for _, medicine := range medicines {
		query := `INSERT INTO THUOC (MaThuoc, TenThuoc, DonVi, Gia, CongDung, LieuLuong) 
				  VALUES (?, ?, ?, ?, ?, ?)`
		_, err := db.Exec(query, medicine.MaThuoc, medicine.TenThuoc, medicine.DonVi, medicine.Gia, medicine.CongDung, medicine.LieuLuong)
		if err != nil {
			return fmt.Errorf("failed to insert medicine %s: %w", medicine.MaThuoc, err)
		}
	}
	
	log.Printf("Seeded %d medicines", len(medicines))
	return nil
}

func seedAppointments(db *sql.DB) error {
	var appointments []Appointment
	if err := loadJSONData("appointments.json", &appointments); err != nil {
		return err
	}
	
	for _, appointment := range appointments {
		query := `INSERT INTO LICHKHAM (MaLichKham, MaCustomer, MaBacSi, MaPhongKham, NgayGioKham, TrangThai, GhiChu, NgayDat) 
				  VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
		_, err := db.Exec(query, appointment.MaLichKham, appointment.MaCustomer, appointment.MaBacSi, appointment.MaPhongKham, appointment.NgayGioKham, appointment.TrangThai, appointment.GhiChu, appointment.NgayDat)
		if err != nil {
			return fmt.Errorf("failed to insert appointment %s: %w", appointment.MaLichKham, err)
		}
	}
	
	log.Printf("Seeded %d appointments", len(appointments))
	return nil
}

func seedWorkSchedules(db *sql.DB) error {
	var schedules []WorkSchedule
	if err := loadJSONData("work_schedules.json", &schedules); err != nil {
		return err
	}
	
	for _, schedule := range schedules {
		query := `INSERT INTO LICHLAMVIEC (MaLichLamViec, MaBacSi, MaPhongKham, NgayLamViec, GioBatDau, GioKetThuc, TrangThai) 
				  VALUES (?, ?, ?, ?, ?, ?, ?)`
		_, err := db.Exec(query, schedule.MaLichLamViec, schedule.MaBacSi, schedule.MaPhongKham, schedule.NgayLamViec, schedule.GioBatDau, schedule.GioKetThuc, schedule.TrangThai)
		if err != nil {
			return fmt.Errorf("failed to insert work schedule %s: %w", schedule.MaLichLamViec, err)
		}
	}
	
	log.Printf("Seeded %d work schedules", len(schedules))
	return nil
}

func seedMedicalRecords(db *sql.DB) error {
	var records []MedicalRecord
	if err := loadJSONData("medical_records.json", &records); err != nil {
		return err
	}
	
	for _, record := range records {
		query := `INSERT INTO HOSOBENH (MaHoSo, MaCustomer, MaBacSi, MaPhongKham, NgayKham, TrieuChung, ChanDoan, HuongDanDieuTri, MaICD10, NgayTaiKham) 
				  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
		_, err := db.Exec(query, record.MaHoSo, record.MaCustomer, record.MaBacSi, record.MaPhongKham, record.NgayKham, record.TrieuChung, record.ChanDoan, record.HuongDanDieuTri, record.MaICD10, record.NgayTaiKham)
		if err != nil {
			return fmt.Errorf("failed to insert medical record %s: %w", record.MaHoSo, err)
		}
	}
	
	log.Printf("Seeded %d medical records", len(records))
	return nil
}