# GoErrorKit Demo Application

Ứng dụng demo toàn diện cho **GoErrorKit** với Fiber framework, showcase tất cả tính năng chính của thư viện.

## 📋 Table of Contents

- [Prerequisites](#prerequisites)
- [Chạy Demo](#chạy-demo)
- [Build Modes](#build-modes)
- [Features Overview](#features-overview)
- [Development Tools](#development-tools)
- [Kiểm tra Logs](#kiểm-tra-logs)
- [Code Examples](#code-examples)

## Prerequisites

- Go 1.21+
- Terminal với hỗ trợ UTF-8 (để hiển thị emoji và Vietnamese)
- `jq` (optional, để format JSON logs): `brew install jq`

## Chạy Demo

### Production Mode (Default)

```bash
# Chạy ở production mode - trace/debug logging bị tắt
cd examples
go run main.go
```

### Development Mode (với Trace/Debug)

```bash
# Chạy ở development mode - bật trace/debug logging
cd examples
go run -tags=debug main.go
```

Server sẽ chạy tại `http://localhost:8081`

💡 **Lưu ý**: Khi server khởi động, bạn sẽ thấy demo migration error được log ra console và file.

## Build Modes

GoErrorKit hỗ trợ 2 build modes:

| Mode | Command | Trace/Debug Logs | Use Case |
|------|---------|------------------|----------|
| **Production** | `go run main.go` | ❌ Tắt (no-op) | Production servers |
| **Development** | `go run -tags=debug main.go` | ✅ Bật | Local development, debugging |

📖 **Chi tiết**: Xem [BUILD_MODES_DEMO.md](BUILD_MODES_DEMO.md)

## Features Overview

### 🎯 1. Dual-Level Logging

GoErrorKit hỗ trợ **phân cấp log level** giữa console và file:

- **Console** (`LogLevel: "trace"`): Log tất cả errors từ trace/debug (cần `-tags=debug`), info, warn, error
- **File** (`FileLogLevel: "error"`): Chỉ log errors nghiêm trọng (`error`, `panic`)

**Kết quả:**
- ✅ ValidationError (level: `warn`) → Console: ✓, File: ✗
- ✅ AuthError (level: `warn`) → Console: ✓, File: ✗
- ✅ SystemError (level: `error`) → Console: ✓, File: ✓
- ✅ PanicError (level: `error`) → Console: ✓, File: ✓

### 🔍 2. Smart Stack Trace Filtering

Tự động lọc stack trace để chỉ hiển thị **code của bạn**, bỏ qua:
- Go runtime code (`runtime.*`)
- Thư viện bên thứ 3 (fiber, goerrorkit, etc.)

```go
goerrorkit.ConfigureForApplication("main")
```

### 🎚️ 3. Log Level Override với `.Level()`

Override log level của error bằng fluent API:

```go
// Force validation error vào file
goerrorkit.NewValidationError("Suspicious input", nil).Level("error")

// Downgrade business error
goerrorkit.NewBusinessError(404, "Not found").Level("warn")
```

### 🎁 4. Wrap Go Errors

Wrap standard Go errors thành GoErrorKit errors:

```go
// Giữ nguyên message gốc
goerrorkit.Wrap(err)

// Thêm context message
goerrorkit.WrapWithMessage(err, "Failed to connect database")
```

### 🔗 5. Call Chain Tracking

Thêm full call chain cho non-panic errors:

```go
goerrorkit.NewValidationError("Invalid", nil).WithCallChain()
```

### 🐛 6. Development Logging (Trace & Debug)

Log chi tiết cho development (chỉ với `-tags=debug`):

```go
goerrorkit.Trace("Operation started", data)
goerrorkit.Debug("Current state", data)
```

## Development Tools

⚠️ **QUAN TRỌNG**: Các endpoints này CHỈ hoạt động khi build với `-tags=debug`!

```bash
# ✅ Development mode - trace/debug hoạt động
go run -tags=debug main.go

# ❌ Production mode - trace/debug là no-op
go run main.go
```

### 🔍 5. Trace Logging

Track operations trong development để hiểu flow execution:

```bash
# Trace single operation - Fetch user
curl "http://localhost:8081/dev/trace?op=fetch_user&user_id=12345"

# Trace cache miss event
curl "http://localhost:8081/dev/trace?op=cache_miss&key=user:12345"

# Trace slow query warning
curl "http://localhost:8081/dev/trace?op=slow_query&query=SELECT+*+FROM+users"
```

**Use case**: Track các events không phải errors (cache miss, slow queries, etc.)

### 🐛 6. Debug Logging

Log chi tiết variable states và object properties:

```bash
# Debug user login flow
curl "http://localhost:8081/dev/debug?scenario=user_login&username=john@example.com"

# Debug payment processing
curl "http://localhost:8081/dev/debug?scenario=payment_process&amount=100000&currency=VND"

# Debug external API calls
curl "http://localhost:8081/dev/debug?scenario=api_request&service=user-service"
```

**Use case**: Log detailed context để debug complex flows

### 📊 7. Trace Complex Flow

Trace toàn bộ multi-step operation với timing cho từng step:

```bash
# Trace complete order processing flow (6 steps)
curl "http://localhost:8081/dev/trace-complex?order_id=ORD-12345"
```

**Kết quả console** (khi chạy với `-tags=debug`):
```
[TRACE] Order Processing Flow Started
[TRACE] Step 1: Validating order (50ms)
[TRACE] Step 2: Checking inventory (120ms)  
[TRACE] Step 3: Reserving inventory (80ms)
[TRACE] Step 4: Processing payment (450ms)
[TRACE] Step 5: Creating shipment (200ms)
[TRACE] Step 6: Sending confirmation email (300ms)
[TRACE] Order Processing Flow Completed (total: 1200ms)
```

**Use case**: Performance profiling và identify bottlenecks

### 🎭 Demo Trace/Debug Behavior

```bash
# Terminal 1: Chạy production mode
go run main.go
# → Trace/Debug = no-op (không log gì)

# Terminal 2: Test trace endpoint
curl http://localhost:8081/dev/trace?op=fetch_user
# → Response OK nhưng KHÔNG CÓ log trong console

# Terminal 1: Stop server, chạy development mode
go run -tags=debug main.go  
# → Trace/Debug enabled

# Terminal 2: Test lại
curl http://localhost:8081/dev/trace?op=fetch_user
# → CÓ log chi tiết trong console!
```

## Kiểm tra Logs

### 📺 Console Logs

Console hiển thị:
- **Production mode**: info, warn, error, panic
- **Development mode** (`-tags=debug`): trace, debug, info, warn, error, panic

```bash
# Watch console output khi chạy server
go run main.go                    # Production
go run -tags=debug main.go        # Development
```

### 📄 File Logs

File `logs/errors.log` chỉ chứa errors nghiêm trọng (`error`, `panic`):

```bash
# View formatted JSON logs
cat logs/errors.log | jq

# View specific fields
cat logs/errors.log | jq '{level, message, file, line}'

# Count errors in file
wc -l logs/errors.log

# Watch logs real-time
tail -f logs/errors.log | jq

# Clear logs để test lại
rm logs/errors.log
```

### 🔄 Database Migration Demo

Khi server khởi động, bạn sẽ thấy demo migration error:

```
⚠️  Migration failed but server will continue...
```

Error này được log ra cả console và file (`logs/errors.log`) với đầy đủ thông tin:
- Error message: "Failed to run database migrations"
- Cause: "dial tcp 127.0.0.1:5432: connect: connection refused"
- Metadata: database, host, migration_files, versions
- Stack trace (chỉ code của bạn!)

### 🔍 So sánh Console vs File

```bash
# Terminal 1: Run server
go run main.go

# Terminal 2: Trigger validation error (warn level)
curl "http://localhost:8081/error/validation?age=15"
# → Console: ✓ (có log)
# → File: ✗ (không có - vì warn < error)

# Terminal 2: Trigger system error (error level)
curl http://localhost:8081/error/system
# → Console: ✓ (có log)
# → File: ✓ (có log - vì error >= error)

# Terminal 2: Check file
cat logs/errors.log | jq '.level'
# Chỉ thấy "error", không thấy "warn"
```

### 📊 Log Analysis Examples

```bash
# Count errors by level
cat logs/errors.log | jq -r '.level' | sort | uniq -c

# Find all SystemErrors
cat logs/errors.log | jq 'select(.error_type == "SystemError")'

# Find errors in specific file
cat logs/errors.log | jq 'select(.file | contains("main.go"))'

# Extract error messages only
cat logs/errors.log | jq -r '.message'

# Find errors with specific data field
cat logs/errors.log | jq 'select(.data.database != null)'
```

## Code Examples

### 1. Khởi tạo Logger với Dual-Level Logging

```go
goerrorkit.InitLogger(goerrorkit.LoggerOptions{
    ConsoleOutput: true,
    FileOutput:    true,
    FilePath:      "logs/errors.log",
    JSONFormat:    true,
    LogLevel:      "trace",  // Console: trace (cần -tags=debug), info, warn, error
    FileLogLevel:  "error",  // File: chỉ error và panic
    MaxFileSize:   10,       // MB
    MaxBackups:    5,
    MaxAge:        30,       // days
})
```

### 2. Cấu hình Stack Trace Filtering

```go
// App đơn giản (1 file main.go)
goerrorkit.ConfigureForApplication("main")

// App với nhiều packages
goerrorkit.ConfigureForApplication("github.com/yourname/project")
// → Tự động include TẤT CẢ sub-packages!

// Thêm custom skip patterns (optional)
goerrorkit.AddSkipPatterns(".RequestID.func", ".Logger.func")

// Hoặc dùng Fluent API chi tiết
goerrorkit.Configure().
    SkipPattern(".CustomMiddleware.func").
    SkipPackage("internal/metrics").
    ShowFullPath(false).
    Apply()
```

### 3. Wrap Standard Go Errors

```go
// Wrap() - Giữ nguyên message gốc
err := fmt.Errorf("connection refused")
return goerrorkit.Wrap(err)

// WrapWithMessage() - Thêm context
err := fmt.Errorf("connection refused")  
return goerrorkit.WrapWithMessage(err, "Failed to connect database")

// Wrap với metadata
return goerrorkit.Wrap(err).WithData(map[string]interface{}{
    "host": "localhost:5432",
    "retries": 3,
})
```

### 4. Custom Error Types với Default Log Levels

```go
// ValidationError → level: "warn" (console only, không vào file)
goerrorkit.NewValidationError("Email không hợp lệ", map[string]interface{}{
    "field": "email",
    "value": "invalid@",
})

// AuthError → level: "warn" (console only)
goerrorkit.NewAuthError(401, "Unauthorized")

// BusinessError → level: "error" (console + file)
goerrorkit.NewBusinessError(404, "Product not found")

// SystemError → level: "error" (console + file)
goerrorkit.NewSystemError(err)

// ExternalError → level: "error" (console + file)
goerrorkit.NewExternalError(502, "Service unavailable", err)
```

### 5. Override Log Level với `.Level()`

```go
// Force validation error vào file (suspicious input)
goerrorkit.NewValidationError("SQL injection attempt", nil).
    Level("error")  // Override: warn → error

// Multiple failed login attempts
goerrorkit.NewAuthError(401, "Brute force detected").
    Level("error").  // Override: warn → error
    WithData(map[string]interface{}{
        "attempts": 5,
        "ip": "192.168.1.100",
    })

// Downgrade business error (optional, rare case)
goerrorkit.NewBusinessError(404, "Temporarily unavailable").
    Level("warn")  // Override: error → warn

// Chain với nhiều methods
goerrorkit.NewSystemError(err).
    WithData(map[string]interface{}{"db": "postgres"}).
    Level("error").
    WithCallChain()
```

### 6. Add Call Chain Tracking

```go
// Thêm full call chain cho non-panic errors
func validateOrder() error {
    if !isValid {
        return goerrorkit.NewValidationError("Invalid order", nil).
            WithCallChain()  // ⭐ Thêm call_chain!
    }
    return nil
}

// Chain với methods khác
return goerrorkit.NewBusinessError(422, "Insufficient inventory").
    WithData(map[string]interface{}{
        "product_id": "PROD-123",
        "available": 0,
    }).
    WithCallChain()  // ⭐ Thêm call_chain để trace flow
```

### 7. Development Logging (Trace & Debug)

⚠️ **CHỈ hoạt động với**: `go run -tags=debug main.go`

```go
// Trace logging - Track operations
goerrorkit.Trace("Fetching user from database", map[string]interface{}{
    "user_id": userID,
})

// Debug logging - Detailed context
goerrorkit.Debug("User login attempt", map[string]interface{}{
    "username": username,
    "ip_address": ipAddr,
    "user_agent": userAgent,
})

// Trace complex flow
goerrorkit.Trace("Step 1: Validating order", map[string]interface{}{
    "order_id": orderID,
    "duration_ms": 50,
})
goerrorkit.Trace("Step 2: Processing payment", map[string]interface{}{
    "amount": 10000,
    "duration_ms": 450,
})
```

### 8. Manual Error Logging

```go
// Log error manually (thay vì return)
if err := someOperation(); err != nil {
    appErr := goerrorkit.NewSystemError(err)
    goerrorkit.LogError(appErr, "/path/to/operation")
    // Server tiếp tục chạy...
}
```

### 9. Integration với Fiber Middleware

```go
app := fiber.New()

// Add RequestID middleware (must be before ErrorHandler)
app.Use(requestid.New())

// Add GoErrorKit error handler
app.Use(goerrorkit.FiberErrorHandler())

// Route handlers
app.Get("/api/users", func(c *fiber.Ctx) error {
    // Return AppError, middleware sẽ tự động xử lý
    return goerrorkit.NewBusinessError(404, "User not found")
})
```

### 10. Database Migration Example

```go
func runDatabaseMigrations(simulateError bool) error {
    if simulateError {
        dbErr := fmt.Errorf("connection refused")
        return goerrorkit.WrapWithMessage(dbErr, "Failed to run migrations").
            WithData(map[string]interface{}{
                "database": "postgresql",
                "host": "127.0.0.1:5432",
                "migration_files": []string{"001_users.sql", "002_products.sql"},
            })
    }
    return nil
}

// Usage in main()
if err := runDatabaseMigrations(true); err != nil {
    goerrorkit.LogError(err.(*goerrorkit.AppError), "/startup/migrations")
    fmt.Println("⚠️  Migration failed but server continues...")
}
```

## 📚 Additional Resources

- **Main README**: [../README.md](../README.md) - Full documentation
- **Build Modes**: [BUILD_MODES_DEMO.md](BUILD_MODES_DEMO.md) - Chi tiết về production vs development modes
- **Configuration Guide**: [../docs/configuration.md](../docs/configuration.md) - Cấu hình nâng cao
- **Stack Trace Guide**: [../docs/stack-trace-configuration.md](../docs/stack-trace-configuration.md) - Stack trace filtering

## 🐛 Troubleshooting

### Trace/Debug không hoạt động

```bash
# ❌ Không hoạt động
go run main.go

# ✅ Phải có -tags=debug
go run -tags=debug main.go
```

### File logs rỗng

Kiểm tra `FileLogLevel` trong config. Nếu set `"error"`, chỉ có errors nghiêm trọng mới được ghi.

```bash
# Trigger error để test
curl http://localhost:8081/error/system
cat logs/errors.log | jq
```

### Stack trace quá dài

Đảm bảo đã gọi `ConfigureForApplication()`:

```go
goerrorkit.ConfigureForApplication("main")  // Hoặc package path của bạn
```

## 🎯 Testing Checklist

- [ ] **Production mode**: `go run main.go` - Trace/debug tắt
- [ ] **Development mode**: `go run -tags=debug main.go` - Trace/debug bật
- [ ] **Panic recovery**: Test `/panic/*` endpoints
- [ ] **Dual-level logging**: Validation errors chỉ ở console, system errors ở cả file
- [ ] **Log level override**: Test `.Level()` API
- [ ] **Wrap errors**: Test `Wrap()` và `WrapWithMessage()`
- [ ] **Call chain**: Test `WithCallChain()`
- [ ] **Trace/Debug**: Test `/dev/*` endpoints với `-tags=debug`
- [ ] **Migration demo**: Xem log khi server khởi động

