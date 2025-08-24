# Clinic Management System Backend

Hệ thống quản lý phòng khám được xây dựng bằng Go Gin framework và SQL Server.

## Tính năng chính

### Khách hàng (Bệnh nhân)
- CN01-CN05: Quản lý tài khoản (đăng ký, đăng nhập, cập nhật thông tin, đổi mật khẩu)
- CN06-CN13: Đặt lịch khám (xem bác sĩ, lịch làm việc, đặt/hủy/thay đổi lịch hẹn)
- CN14-CN19: Quản lý hồ sơ bệnh án (xem lịch sử khám, chi tiết khám bệnh, tải PDF)

### Bác sĩ
- CN36-CN50: Quản lý lịch làm việc, khám bệnh, ghi chép hồ sơ, kê đơn thuốc

### Lễ tân
- CN23-CN35: Quản lý bệnh nhân, hỗ trợ đặt lịch, thu ngân

### Kế toán
- CN51-CN59: Quản lý lương thù lao, báo cáo tài chính

### Quản lý phòng khám
- CN60-CN72: Quản lý lịch khám, nhân sự, cơ sở vật chất

### Ban điều hành
- CN73-CN82: Báo cáo thống kê tổng thể, quản lý chiến lược

## API Endpoints

### Authentication
- `POST /api/v1/auth/login` - Đăng nhập
- `POST /api/v1/auth/register` - Đăng ký tài khoản khách hàng
- `POST /api/v1/auth/forgot-password` - Quên mật khẩu
- `POST /api/v1/auth/refresh` - Làm mới token

### User Management
- `GET /api/v1/users/profile` - Xem thông tin cá nhân
- `PUT /api/v1/users/profile` - Cập nhật thông tin cá nhân
- `PUT /api/v1/users/password` - Đổi mật khẩu

### Clinics
- `GET /api/v1/clinics` - Danh sách phòng khám
- `GET /api/v1/clinics/:id` - Thông tin phòng khám
- `GET /api/v1/clinics/:id/doctors` - Danh sách bác sĩ theo phòng khám
- `GET /api/v1/clinics/:id/schedules` - Lịch làm việc theo phòng khám

### Appointments
- `GET /api/v1/appointments` - Danh sách lịch khám
- `POST /api/v1/appointments` - Đặt lịch khám
- `GET /api/v1/appointments/:id` - Chi tiết lịch khám
- `PUT /api/v1/appointments/:id` - Cập nhật lịch khám
- `DELETE /api/v1/appointments/:id` - Hủy lịch khám

### Medical Records
- `GET /api/v1/medical-records` - Danh sách hồ sơ bệnh án
- `GET /api/v1/medical-records/:id` - Chi tiết hồ sơ bệnh án
- `POST /api/v1/medical-records` - Tạo hồ sơ bệnh án (chỉ bác sĩ)
- `PUT /api/v1/medical-records/:id` - Cập nhật hồ sơ bệnh án (chỉ bác sĩ)

## Cài đặt và chạy

1. **Cài đặt dependencies:**
   ```bash
   go mod tidy
   ```

2. **Cấu hình database:**
   - Tạo database SQL Server với tên `clinic_management`
   - Chạy script SQL trong file `server.sql` để tạo bảng
   - Cấu hình connection string trong file `.env`

3. **Tạo file .env:**
   ```bash
   cp .env.example .env
   ```
   Sau đó chỉnh sửa các thông số phù hợp với môi trường của bạn.

4. **Chạy server:**
   ```bash
   go run main.go
   ```

Server sẽ chạy trên port 8080 (hoặc port được cấu hình trong .env).

## Cấu trúc thư mục

```
clinic-management/
├── main.go
├── internal/
│   ├── config/          # Cấu hình ứng dụng
│   ├── database/        # Kết nối database
│   ├── handlers/        # Xử lý HTTP requests
│   ├── middleware/      # Middleware (auth, CORS, etc.)
│   ├── models/          # Data models
│   ├── routes/          # Định tuyến API
│   ├── services/        # Business logic
│   └── utils/           # Tiện ích chung
├── server.sql          # Database schema
├── .env.example        # Cấu hình mẫu
└── README.md
```

## Database Schema

Hệ thống sử dụng SQL Server với các bảng chính:
- `USER` - Thông tin người dùng cơ bản
- `CUSTOMER`, `BACSI`, `LETAN`, `KETOAN`, etc. - Thông tin chi tiết theo từng loại người dùng
- `PHONGKHAM` - Thông tin phòng khám
- `LICHKHAM` - Lịch khám bệnh
- `LICHLAMVIEC` - Lịch làm việc của bác sĩ
- `HOSOBENH` - Hồ sơ bệnh án
- `DONTHUOC`, `CHITIETDONTHUOC` - Đơn thuốc và chi tiết
- `XETNGHIEM` - Kết quả xét nghiệm
- `THANHTOAN` - Thông tin thanh toán
- `LUONGTHULAO` - Lương và thù lao
- `BAOCAO` - Báo cáo hệ thống

## Authentication & Authorization

- Hệ thống sử dụng JWT tokens cho authentication
- Phân quyền theo từng loại người dùng: CUSTOMER, DOCTOR, RECEPTIONIST, ACCOUNTANT, CLINIC_MANAGER, OPERATION_MANAGER
- Middleware bảo vệ các endpoints yêu cầu đăng nhập

## API Response Format

Tất cả API responses đều có format:
```json
{
  "success": true,
  "message": "Success message",
  "data": { ... },
  "error": "Error message (nếu có)"
}
```

## Tính năng sẽ phát triển

- Tích hợp SMS/Email notifications
- File upload cho hình ảnh và kết quả xét nghiệm
- Báo cáo và thống kê nâng cao
- API cho mobile app
- Tích hợp cổng thanh toán
- Backup và phục hồi dữ liệu tự động"# CoSoDuLieuNangCaoBE" 
