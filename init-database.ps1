# PowerShell script to initialize SQL Server database

Write-Host "Waiting for SQL Server to be ready..." -ForegroundColor Green

# Chờ SQL Server sẵn sáng
$maxAttempts = 60
$attempt = 1

while ($attempt -le $maxAttempts) {
    try {
        $result = docker exec clinic_sqlserver /opt/mssql-tools18/bin/sqlcmd -S localhost -U sa -P "StrongPassword123!" -C -Q "SELECT 1" -h -1
        if ($result -eq "1") {
            Write-Host "SQL Server is ready!" -ForegroundColor Green
            break
        }
    }
    catch {
        Write-Host "Waiting for SQL Server... ($attempt/$maxAttempts)" -ForegroundColor Yellow
    }
    
    Start-Sleep -Seconds 1
    $attempt++
}

if ($attempt -gt $maxAttempts) {
    Write-Host "SQL Server failed to start within timeout period" -ForegroundColor Red
    exit 1
}

# Tạo database
Write-Host "Creating clinic_management database..." -ForegroundColor Green
docker exec clinic_sqlserver /opt/mssql-tools18/bin/sqlcmd -S localhost -U sa -P "StrongPassword123!" -C -Q "
IF NOT EXISTS (SELECT name FROM sys.databases WHERE name = 'clinic_management')
BEGIN
    CREATE DATABASE clinic_management;
    PRINT 'Database clinic_management created successfully!';
END
ELSE
    PRINT 'Database clinic_management already exists.';
"

# Tạo tables
Write-Host "Creating database tables..." -ForegroundColor Green
$sqlContent = Get-Content -Path ".\init-db.sql" -Raw
$sqlCommands = $sqlContent -split "GO"

foreach ($command in $sqlCommands) {
    if ($command.Trim() -ne "") {
        docker exec clinic_sqlserver /opt/mssql-tools18/bin/sqlcmd -S localhost -U sa -P "StrongPassword123!" -C -d clinic_management -Q $command.Trim()
    }
}

Write-Host "Database initialization completed!" -ForegroundColor Green