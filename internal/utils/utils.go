package utils

import (
	"crypto/rand"
	"fmt"
	"regexp"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Password hashing functions
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// ID generation functions following database patterns

// Counter for sequential ID generation (in production, this should be managed by database)
var idCounters = make(map[string]int)

// Generate sequential ID with prefix and padding
func generateSequentialID(prefix string, padding int) string {
	idCounters[prefix]++
	return fmt.Sprintf("%s%0*d", prefix, padding, idCounters[prefix])
}

// User ID generators based on roles
func GenerateUserID(role string) string {
	switch strings.ToUpper(role) {
	case "CUSTOMER":
		return generateSequentialID("CUS", 6) // CUS000001
	case "DOCTOR":
		return generateSequentialID("DOC", 3) // DOC001
	case "RECEPTIONIST":
		return generateSequentialID("REC", 3) // REC001
	case "ACCOUNTANT":
		return generateSequentialID("ACC", 3) // ACC001
	case "CLINIC_MANAGER":
		return generateSequentialID("CLM", 3) // CLM001
	case "OPERATION_MANAGER":
		return generateSequentialID("OPM", 3) // OPM001
	default:
		return generateSequentialID("USR", 6) // USR000001
	}
}

// Legacy function for backward compatibility
func GenerateCustomerID() string {
	return GenerateUserID("CUSTOMER")
}

// Clinic and facility ID generators
func GenerateClinicID() string {
	return generateSequentialID("PK", 3) // PK001 (PhongKham)
}

func GenerateWorkScheduleID() string {
	return generateSequentialID("LLV", 6) // LLV000001 (LichLamViec)
}

// Medical record and appointment ID generators
func GenerateAppointmentID() string {
	return generateSequentialID("LK", 6) // LK000001 (LichKham)
}

func GenerateMedicalRecordID() string {
	return generateSequentialID("HS", 6) // HS000001 (HoSo)
}

func GeneratePrescriptionID() string {
	return generateSequentialID("DT", 6) // DT000001 (DonThuoc)
}

func GenerateTestResultID() string {
	return generateSequentialID("XN", 6) // XN000001 (XetNghiem)
}

func GenerateLabTestID() string {
	return generateSequentialID("XN", 6) // XN000001 (XetNghiem)
}

func GenerateScheduleID() string {
	return generateSequentialID("LLV", 6) // LLV000001 (LichLamViec)
}

func GenerateMedicalImageID() string {
	return generateSequentialID("HA", 6) // HA000001 (HinhAnhKham)
}

// Medicine ID generator
func GenerateMedicineID() string {
	return generateSequentialID("MED", 3) // MED001
}

// Financial ID generators
func GeneratePaymentID() string {
	return generateSequentialID("TT", 6) // TT000001 (ThanhToan)
}

func GenerateSalaryID() string {
	return generateSequentialID("LG", 6) // LG000001 (Luong)
}

func GenerateReportID() string {
	return generateSequentialID("BC", 6) // BC000001 (BaoCao)
}

// Time slot ID generator (for appointment scheduling)
func GenerateTimeSlotID() string {
	return generateSequentialID("GK", 3) // GK001 (GioKham)
}

// Validation functions
func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func ValidatePhoneNumber(phone string) bool {
	// Vietnamese phone number validation
	phoneRegex := regexp.MustCompile(`^(0[3|5|7|8|9])+([0-9]{8})$`)
	return phoneRegex.MatchString(phone)
}

func ValidateUserRole(role string) bool {
	validRoles := []string{"CUSTOMER", "DOCTOR", "RECEPTIONIST", "ACCOUNTANT", "CLINIC_MANAGER", "OPERATION_MANAGER"}
	roleUpper := strings.ToUpper(role)
	for _, validRole := range validRoles {
		if roleUpper == validRole {
			return true
		}
	}
	return false
}

func ValidateGender(gender string) bool {
	validGenders := []string{"Nam", "Nữ", "Khác"}
	for _, validGender := range validGenders {
		if gender == validGender {
			return true
		}
	}
	return false
}

func ValidateAppointmentStatus(status string) bool {
	validStatuses := []string{"SCHEDULED", "COMPLETED", "CANCELLED", "NO_SHOW"}
	statusUpper := strings.ToUpper(status)
	for _, validStatus := range validStatuses {
		if statusUpper == validStatus {
			return true
		}
	}
	return false
}

func ValidatePaymentStatus(status string) bool {
	validStatuses := []string{"PENDING", "COMPLETED", "FAILED", "CANCELLED"}
	statusUpper := strings.ToUpper(status)
	for _, validStatus := range validStatuses {
		if statusUpper == validStatus {
			return true
		}
	}
	return false
}

func ValidatePaymentMethod(method string) bool {
	validMethods := []string{"Tiền mặt", "Thẻ ATM", "Chuyển khoản", "Ví điện tử"}
	for _, validMethod := range validMethods {
		if method == validMethod {
			return true
		}
	}
	return false
}

// Date and time formatting functions
func FormatVietnameseDate(t time.Time) string {
	return t.Format("02/01/2006 15:04")
}

func FormatDateOnly(t time.Time) string {
	return t.Format("02/01/2006")
}

func FormatTimeOnly(t time.Time) string {
	return t.Format("15:04")
}

func ParseVietnameseDate(dateStr string) (time.Time, error) {
	return time.Parse("02/01/2006 15:04", dateStr)
}

func ParseDateOnly(dateStr string) (time.Time, error) {
	return time.Parse("02/01/2006", dateStr)
}

func ParseTimeOnly(timeStr string) (time.Time, error) {
	return time.Parse("15:04", timeStr)
}

// Business logic helper functions
func IsWorkingDay(t time.Time) bool {
	weekday := t.Weekday()
	return weekday >= time.Monday && weekday <= time.Friday
}

func IsWorkingHour(t time.Time) bool {
	hour := t.Hour()
	return (hour >= 7 && hour < 12) || (hour >= 13 && hour < 22)
}

func CalculateAge(birthDate time.Time) int {
	now := time.Now()
	age := now.Year() - birthDate.Year()
	if now.YearDay() < birthDate.YearDay() {
		age--
	}
	return age
}

// Generate reset code for password recovery
func GenerateResetCode() string {
	randomBytes := make([]byte, 3)
	rand.Read(randomBytes)

	code := 0
	for _, b := range randomBytes {
		code = (code * 256) + int(b)
	}

	return fmt.Sprintf("%06d", code%1000000)
}

func GeneratePasswordResetID() string {
	return generateSequentialID("PWR", 6) // PWR000001
}

func FormatCurrency(amount float64) string {
	return fmt.Sprintf("%.0f VND", amount)
}

// Report period helpers
func GetMonthRange(year int, month int) (time.Time, time.Time) {
	start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, -1).Add(23*time.Hour + 59*time.Minute + 59*time.Second)
	return start, end
}

func GetYearRange(year int) (time.Time, time.Time) {
	start := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(year, 12, 31, 23, 59, 59, 0, time.UTC)
	return start, end
}

// Initialize counters (should be called once at application startup)
func InitializeCounters() {
	// In production, these values should be loaded from database
	// This is just for demonstration
	idCounters["CUS"] = 50000
	idCounters["DOC"] = 50
	idCounters["REC"] = 20
	idCounters["ACC"] = 5
	idCounters["CLM"] = 0
	idCounters["OPM"] = 0
	idCounters["PK"] = 10
	idCounters["LK"] = 40000 // Set to higher than existing data
	idCounters["HS"] = 20000 // Set to higher than existing data
	idCounters["DT"] = 20000 // Set to higher than existing data
	idCounters["XN"] = 4000  // Set to higher than existing data
	idCounters["HA"] = 2000  // Set to higher than existing data
	idCounters["MED"] = 100
	idCounters["TT"] = 25000 // Set to higher than existing data
	idCounters["LG"] = 5000  // Set to higher than existing data
	idCounters["BC"] = 100   // Set to higher than existing data
	idCounters["GK"] = 0
	idCounters["LLV"] = 50000 // Set to higher than existing data
	idCounters["PWR"] = 0
}
