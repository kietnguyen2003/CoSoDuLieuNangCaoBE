# Database Mock Data

Thư mục này chứa dữ liệu mẫu (mock data) cho hệ thống quản lý phòng khám.

## Cấu trúc

```
database/
├── mockdata/              # Dữ liệu JSON mock
│   ├── users.json         # Người dùng cơ bản
│   ├── customers.json     # Thông tin bệnh nhân
│   ├── doctors.json       # Thông tin bác sĩ
│   ├── clinics.json       # Phòng khám
│   ├── receptionists.json # Lễ tân
│   ├── accountants.json   # Kế toán
│   ├── clinic_managers.json    # Quản lý phòng khám
│   ├── operation_managers.json # Ban điều hành
│   ├── medicines.json     # Thuốc
│   ├── appointments.json  # Lịch khám
│   ├── work_schedules.json # Lịch làm việc bác sĩ
│   └── medical_records.json # Hồ sơ bệnh án
├── seed.go               # Utility để load mock data
└── README.md            # File này
```

## Cách sử dụng

### 1. Tự động load mock data khi khởi động server

Trong `main.go`, thêm code để seed database:

```go
import "clinic-management/database"

func main() {
    // ... existing code ...
    
    db, err := database.Connect(cfg.DatabaseURL)
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    defer db.Close()
    
    // Seed mock data (optional - chỉ chạy khi database trống)
    if cfg.SeedDatabase {
        if err := database.SeedDatabase(db); err != nil {
            log.Printf("Warning: Failed to seed database: %v", err)
        }
    }
    
    // ... rest of code ...
}
```

### 2. Manual seeding

Tạo một CLI command riêng để seed data:

```bash
go run cmd/seed/main.go
```

## Dữ liệu mẫu

### Người dùng (Users)
- 8 users với các role khác nhau
- Mật khẩu đã được hash (sử dụng bcrypt)
- Bao gồm: 2 bệnh nhân, 2 bác sĩ, 1 lễ tân, 1 kế toán, 1 quản lý, 1 điều hành

### Phòng khám (Clinics)
- 5 phòng khám chuyên khoa khác nhau
- Tim mạch, Sản phụ khoa, Nhi khoa, Da liễu, Đa khoa

### Lịch khám (Appointments)
- 5 cuộc hẹn mẫu với các trạng thái khác nhau
- Bao gồm cả lịch đã hoàn thành và sắp tới

### Thuốc (Medicines)
- 6 loại thuốc phổ biến
- Bao gồm giá, liều lượng và công dụng

## Lưu ý

1. **Không kết nối server**: Mock data được thiết kế để chạy độc lập, không cần kết nối SQL Server thực tế
2. **Dữ liệu quan hệ**: Các foreign key đã được thiết kế để liên kết chính xác
3. **Định dạng ngày tháng**: Sử dụng ISO 8601 format (YYYY-MM-DDTHH:MM:SSZ)
4. **Mật khẩu**: Tất cả mật khẩu đều đã được hash, không lưu plain text

## Mở rộng

Để thêm dữ liệu mock:

1. Tạo file JSON mới trong `/mockdata`
2. Định nghĩa struct tương ứng trong `seed.go`
3. Thêm function `seedXXX()` 
4. Gọi function trong `SeedDatabase()`