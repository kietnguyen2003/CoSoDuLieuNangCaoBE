# Chạy Server với Mock Data

## Cách sử dụng

### 1. Chạy server với mock data (không cần database)
```bash
go run main_mock.go
```

### 2. Chạy server với database thật (cần SQL Server)
```bash
go run main.go
```

## Test API với Mock Data

### Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "tenDangNhap": "nguyenvanan",
    "matKhau": "123456"
  }'
```

### Lấy danh sách phòng khám
```bash
curl -X GET http://localhost:8080/api/v1/clinics \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Lấy profile user
```bash
curl -X GET http://localhost:8080/api/v1/users/profile \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## Mock Users có thể login

| Username | Role | MaUser |
|----------|------|--------|
| nguyenvanan | CUSTOMER | USER001 |
| tranthibinh | CUSTOMER | USER002 |
| bsleminhcuong | DOCTOR | DOC001 |
| bsphamthudung | DOCTOR | DOC002 |
| hoangthien | RECEPTIONIST | REC001 |
| dovanphat | ACCOUNTANT | ACC001 |
| vuthigiang | CLINIC_MANAGER | MAN001 |
| ngovanhai | OPERATION_MANAGER | OPE001 |

**Mật khẩu**: Bất kỳ password nào (mock implementation)

## Tính năng Mock

✅ **Hoạt động**:
- Authentication (login, register)
- User profile management
- Clinics listing
- Appointments CRUD
- Medical records CRUD
- Role-based authorization

❌ **Chưa implement**:
- Thực sự lưu data mới
- Password reset email
- File upload
- Database persistence

## Cấu trúc Files

```
internal/
├── services/
│   └── mock_data.go        # Mock data service
├── handlers/
│   ├── mock_auth.go        # Authentication handlers
│   ├── mock_user.go        # User management
│   ├── mock_clinic.go      # Clinic management
│   ├── mock_appointment.go # Appointment management
│   └── mock_medical_record.go # Medical record management
└── routes/
    └── mock_routes.go      # Routes for mock mode
```

## Lưu ý

1. **Không lưu data**: Tất cả thao tác CREATE/UPDATE chỉ trả về success message
2. **Password**: Accept bất kỳ password nào trong mock mode
3. **JWT**: Sử dụng secret key cố định cho demo
4. **File paths**: Mock data đọc từ `database/mockdata/`