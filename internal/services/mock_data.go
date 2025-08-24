package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type MockDataService struct {
	dataPath string
	users    []User
	customers []Customer
	doctors   []Doctor
	clinics   []Clinic
	receptionists []Receptionist
	accountants   []Accountant
	clinicManagers []ClinicManager
	operationManagers []OperationManager
	medicines     []Medicine
	appointments  []Appointment
	workSchedules []WorkSchedule
	medicalRecords []MedicalRecord
}

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
	MaUser       string   `json:"MaUser"`
	MaPhongKham  *string  `json:"MaPhongKham"`
	LuongCoBan   *float64 `json:"LuongCoBan"`
	NgayVaoLam   *string  `json:"NgayVaoLam"`
}

type Accountant struct {
	MaUser       string   `json:"MaUser"`
	LuongCoBan   *float64 `json:"LuongCoBan"`
	NgayVaoLam   *string  `json:"NgayVaoLam"`
	ChuyenMon    *string  `json:"ChuyenMon"`
}

type ClinicManager struct {
	MaUser       string   `json:"MaUser"`
	MaPhongKham  *string  `json:"MaPhongKham"`
	LuongCoBan   *float64 `json:"LuongCoBan"`
	NgayVaoLam   *string  `json:"NgayVaoLam"`
}

type OperationManager struct {
	MaUser         string   `json:"MaUser"`
	ChucVu         *string  `json:"ChucVu"`
	KhuVucPhuTrach *string  `json:"KhuVucPhuTrach"`
	LuongCoBan     *float64 `json:"LuongCoBan"`
	NgayVaoLam     *string  `json:"NgayVaoLam"`
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

func NewMockDataService() (*MockDataService, error) {
	service := &MockDataService{
		dataPath: filepath.Join("database", "mockdata"),
	}
	
	if err := service.loadAllData(); err != nil {
		return nil, fmt.Errorf("failed to load mock data: %w", err)
	}
	
	return service, nil
}

func (s *MockDataService) loadAllData() error {
	if err := s.loadJSONData("users.json", &s.users); err != nil {
		return err
	}
	if err := s.loadJSONData("customers.json", &s.customers); err != nil {
		return err
	}
	if err := s.loadJSONData("doctors.json", &s.doctors); err != nil {
		return err
	}
	if err := s.loadJSONData("clinics.json", &s.clinics); err != nil {
		return err
	}
	if err := s.loadJSONData("receptionists.json", &s.receptionists); err != nil {
		return err
	}
	if err := s.loadJSONData("accountants.json", &s.accountants); err != nil {
		return err
	}
	if err := s.loadJSONData("clinic_managers.json", &s.clinicManagers); err != nil {
		return err
	}
	if err := s.loadJSONData("operation_managers.json", &s.operationManagers); err != nil {
		return err
	}
	if err := s.loadJSONData("medicines.json", &s.medicines); err != nil {
		return err
	}
	if err := s.loadJSONData("appointments.json", &s.appointments); err != nil {
		return err
	}
	if err := s.loadJSONData("work_schedules.json", &s.workSchedules); err != nil {
		return err
	}
	if err := s.loadJSONData("medical_records.json", &s.medicalRecords); err != nil {
		return err
	}
	
	return nil
}

func (s *MockDataService) loadJSONData(filename string, data interface{}) error {
	filePath := filepath.Join(s.dataPath, filename)
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", filename, err)
	}
	
	if err := json.Unmarshal(content, data); err != nil {
		return fmt.Errorf("failed to unmarshal JSON from %s: %w", filename, err)
	}
	
	return nil
}

// User operations
func (s *MockDataService) GetUserByUsername(username string) (*User, error) {
	for _, user := range s.users {
		if user.TenDangNhap == username {
			return &user, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}

func (s *MockDataService) GetUserByID(userID string) (*User, error) {
	for _, user := range s.users {
		if user.MaUser == userID {
			return &user, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}

func (s *MockDataService) GetAllUsers() []User {
	return s.users
}

// Customer operations
func (s *MockDataService) GetCustomerByUserID(userID string) (*Customer, error) {
	for _, customer := range s.customers {
		if customer.MaUser == userID {
			return &customer, nil
		}
	}
	return nil, fmt.Errorf("customer not found")
}

func (s *MockDataService) GetAllCustomers() []Customer {
	return s.customers
}

// Doctor operations
func (s *MockDataService) GetDoctorByUserID(userID string) (*Doctor, error) {
	for _, doctor := range s.doctors {
		if doctor.MaUser == userID {
			return &doctor, nil
		}
	}
	return nil, fmt.Errorf("doctor not found")
}

func (s *MockDataService) GetAllDoctors() []Doctor {
	return s.doctors
}

func (s *MockDataService) GetDoctorsByClinic(clinicID string) []Doctor {
	var result []Doctor
	for _, schedule := range s.workSchedules {
		if schedule.MaPhongKham == clinicID {
			for _, doctor := range s.doctors {
				if doctor.MaUser == schedule.MaBacSi {
					result = append(result, doctor)
					break
				}
			}
		}
	}
	return result
}

// Clinic operations
func (s *MockDataService) GetClinicByID(clinicID string) (*Clinic, error) {
	for _, clinic := range s.clinics {
		if clinic.MaPhongKham == clinicID {
			return &clinic, nil
		}
	}
	return nil, fmt.Errorf("clinic not found")
}

func (s *MockDataService) GetAllClinics() []Clinic {
	return s.clinics
}

// Appointment operations
func (s *MockDataService) GetAppointmentsByCustomer(customerID string) []Appointment {
	var result []Appointment
	for _, appointment := range s.appointments {
		if appointment.MaCustomer == customerID {
			result = append(result, appointment)
		}
	}
	return result
}

func (s *MockDataService) GetAppointmentsByDoctor(doctorID string) []Appointment {
	var result []Appointment
	for _, appointment := range s.appointments {
		if appointment.MaBacSi == doctorID {
			result = append(result, appointment)
		}
	}
	return result
}

func (s *MockDataService) GetAppointmentByID(appointmentID string) (*Appointment, error) {
	for _, appointment := range s.appointments {
		if appointment.MaLichKham == appointmentID {
			return &appointment, nil
		}
	}
	return nil, fmt.Errorf("appointment not found")
}

func (s *MockDataService) GetAllAppointments() []Appointment {
	return s.appointments
}

// Work Schedule operations
func (s *MockDataService) GetSchedulesByDoctor(doctorID string) []WorkSchedule {
	var result []WorkSchedule
	for _, schedule := range s.workSchedules {
		if schedule.MaBacSi == doctorID {
			result = append(result, schedule)
		}
	}
	return result
}

func (s *MockDataService) GetSchedulesByClinic(clinicID string) []WorkSchedule {
	var result []WorkSchedule
	for _, schedule := range s.workSchedules {
		if schedule.MaPhongKham == clinicID {
			result = append(result, schedule)
		}
	}
	return result
}

// Medical Record operations
func (s *MockDataService) GetMedicalRecordsByCustomer(customerID string) []MedicalRecord {
	var result []MedicalRecord
	for _, record := range s.medicalRecords {
		if record.MaCustomer == customerID {
			result = append(result, record)
		}
	}
	return result
}

func (s *MockDataService) GetMedicalRecordByID(recordID string) (*MedicalRecord, error) {
	for _, record := range s.medicalRecords {
		if record.MaHoSo == recordID {
			return &record, nil
		}
	}
	return nil, fmt.Errorf("medical record not found")
}

func (s *MockDataService) GetAllMedicalRecords() []MedicalRecord {
	return s.medicalRecords
}

// Medicine operations
func (s *MockDataService) GetMedicineByID(medicineID string) (*Medicine, error) {
	for _, medicine := range s.medicines {
		if medicine.MaThuoc == medicineID {
			return &medicine, nil
		}
	}
	return nil, fmt.Errorf("medicine not found")
}

func (s *MockDataService) GetAllMedicines() []Medicine {
	return s.medicines
}

func (s *MockDataService) SearchMedicines(keyword string) []Medicine {
	var result []Medicine
	keyword = strings.ToLower(keyword)
	
	for _, medicine := range s.medicines {
		if strings.Contains(strings.ToLower(medicine.TenThuoc), keyword) ||
		   (medicine.CongDung != nil && strings.Contains(strings.ToLower(*medicine.CongDung), keyword)) {
			result = append(result, medicine)
		}
	}
	return result
}

// Enhanced Schedule Management Operations
func (s *MockDataService) GetFilteredSchedules(doctorID, clinicID, date string) ([]WorkSchedule, error) {
	var result []WorkSchedule
	
	for _, schedule := range s.workSchedules {
		match := true
		
		if doctorID != "" && schedule.MaBacSi != doctorID {
			match = false
		}
		if clinicID != "" && schedule.MaPhongKham != clinicID {
			match = false
		}
		if date != "" && schedule.NgayLamViec != date {
			match = false
		}
		
		if match {
			result = append(result, schedule)
		}
	}
	
	return result, nil
}

func (s *MockDataService) GetScheduleByID(scheduleID string) (*WorkSchedule, error) {
	for i := range s.workSchedules {
		if s.workSchedules[i].MaLichLamViec == scheduleID {
			return &s.workSchedules[i], nil
		}
	}
	return nil, fmt.Errorf("schedule not found")
}

func (s *MockDataService) CreateSchedule(doctorID, clinicID, date, startTime, endTime string) (*WorkSchedule, error) {
	// Generate new schedule ID
	newID := fmt.Sprintf("SCH%03d", len(s.workSchedules)+1)
	
	newSchedule := WorkSchedule{
		MaLichLamViec: newID,
		MaBacSi:       doctorID,
		MaPhongKham:   clinicID,
		NgayLamViec:   date,
		GioBatDau:     startTime,
		GioKetThuc:    endTime,
		TrangThai:     "ACTIVE",
	}
	
	s.workSchedules = append(s.workSchedules, newSchedule)
	return &newSchedule, nil
}

func (s *MockDataService) UpdateSchedule(scheduleID string, doctorID, clinicID, date, startTime, endTime, status *string) (*WorkSchedule, error) {
	for i := range s.workSchedules {
		if s.workSchedules[i].MaLichLamViec == scheduleID {
			if doctorID != nil {
				s.workSchedules[i].MaBacSi = *doctorID
			}
			if clinicID != nil {
				s.workSchedules[i].MaPhongKham = *clinicID
			}
			if date != nil {
				s.workSchedules[i].NgayLamViec = *date
			}
			if startTime != nil {
				s.workSchedules[i].GioBatDau = *startTime
			}
			if endTime != nil {
				s.workSchedules[i].GioKetThuc = *endTime
			}
			if status != nil {
				s.workSchedules[i].TrangThai = *status
			}
			return &s.workSchedules[i], nil
		}
	}
	return nil, fmt.Errorf("schedule not found")
}

func (s *MockDataService) DeleteSchedule(scheduleID string) error {
	for i, schedule := range s.workSchedules {
		if schedule.MaLichLamViec == scheduleID {
			// Remove from slice
			s.workSchedules = append(s.workSchedules[:i], s.workSchedules[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("schedule not found")
}

func (s *MockDataService) AssignDoctorToSchedule(doctorID, clinicID, date, startTime, endTime string) (*WorkSchedule, error) {
	// Check if a schedule already exists for this time slot at this clinic
	for i := range s.workSchedules {
		if s.workSchedules[i].MaPhongKham == clinicID && 
		   s.workSchedules[i].NgayLamViec == date &&
		   s.workSchedules[i].GioBatDau == startTime &&
		   s.workSchedules[i].GioKetThuc == endTime {
			// Update existing schedule
			s.workSchedules[i].MaBacSi = doctorID
			s.workSchedules[i].TrangThai = "ACTIVE"
			return &s.workSchedules[i], nil
		}
	}
	
	// Create new schedule
	return s.CreateSchedule(doctorID, clinicID, date, startTime, endTime)
}

func (s *MockDataService) ReassignDoctor(scheduleID, newDoctorID string) (*WorkSchedule, error) {
	for i := range s.workSchedules {
		if s.workSchedules[i].MaLichLamViec == scheduleID {
			s.workSchedules[i].MaBacSi = newDoctorID
			return &s.workSchedules[i], nil
		}
	}
	return nil, fmt.Errorf("schedule not found")
}

func (s *MockDataService) HasScheduleConflict(doctorID, date, startTime, endTime string) bool {
	for _, schedule := range s.workSchedules {
		if schedule.MaBacSi == doctorID && schedule.NgayLamViec == date && schedule.TrangThai == "ACTIVE" {
			// Check time overlap
			if timeOverlaps(schedule.GioBatDau, schedule.GioKetThuc, startTime, endTime) {
				return true
			}
		}
	}
	return false
}

// Helper function to check if two time ranges overlap
func timeOverlaps(start1, end1, start2, end2 string) bool {
	// Simple string comparison for HH:MM format
	// This assumes times are in proper format
	if start2 >= end1 || start1 >= end2 {
		return false // No overlap
	}
	return true // Overlap detected
}