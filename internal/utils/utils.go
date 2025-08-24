package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateID(prefix string) string {
	timestamp := time.Now().UnixNano()
	randomBytes := make([]byte, 4)
	rand.Read(randomBytes)
	randomHex := hex.EncodeToString(randomBytes)
	return fmt.Sprintf("%s_%d_%s", prefix, timestamp, randomHex)
}

func GenerateUserID() string {
	return GenerateID("USR")
}

func GenerateClinicID() string {
	return GenerateID("CLN")
}

func GenerateAppointmentID() string {
	return GenerateID("APT")
}

func GenerateMedicalRecordID() string {
	return GenerateID("MRC")
}

func GeneratePrescriptionID() string {
	return GenerateID("PRE")
}

func GenerateTestResultID() string {
	return GenerateID("TST")
}

func GeneratePaymentID() string {
	return GenerateID("PAY")
}

func GenerateSalaryID() string {
	return GenerateID("SAL")
}

func GenerateReportID() string {
	return GenerateID("RPT")
}

func ValidateEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func FormatVietnameseDate(t time.Time) string {
	return t.Format("02/01/2006 15:04")
}

func ParseVietnameseDate(dateStr string) (time.Time, error) {
	return time.Parse("02/01/2006 15:04", dateStr)
}

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
	return GenerateID("PWR")
}