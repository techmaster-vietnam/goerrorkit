# Fiber Demo Application

Ứng dụng demo cho GoErrorKit với Fiber framework.

## Chạy Demo

```bash
cd examples/fiber-demo
go run main.go
```

Server sẽ chạy tại `http://localhost:8081`

## Features Demo

### 🎯 Dual-Level Logging (NEW!)

GoErrorKit hỗ trợ **phân cấp log level** giữa console và file:

- **Console**: Log tất cả errors từ `warn` trở lên (để developer debug)
- **File**: Chỉ log errors nghiêm trọng (`error`, `panic`) để dễ phân tích production issues

**Kết quả:**
- ✅ ValidationError (level: `warn`) → Console: ✓, File: ✗
- ✅ AuthError (level: `warn`) → Console: ✓, File: ✗
- ✅ SystemError (level: `error`) → Console: ✓, File: ✓
- ✅ PanicError (level: `error`) → Console: ✓, File: ✓

### 📝 Test Endpoints

#### 1. Panic Demos (Tự động recovered)

```bash
# Division by zero panic
curl http://localhost:8081/panic/division

# Index out of range panic
curl http://localhost:8081/panic/index

# Deep call stack panic
curl http://localhost:8081/panic/stack
```

#### 2. Wrap Error Demos

```bash
# Wrap() - Đơn giản nhất
curl http://localhost:8081/error/wrap?type=json

# WrapWithMessage() - Thêm context message
curl http://localhost:8081/error/wrap-message?scenario=database
```

#### 3. Custom Error Demos

```bash
# Business error (404)
curl http://localhost:8081/error/business?product_id=123

# System error (500) → Logs to FILE
curl http://localhost:8081/error/system

# Validation error (400) → Console only, NOT in file
curl http://localhost:8081/error/validation?age=15

# Auth error (401) → Console only, NOT in file
curl http://localhost:8081/error/auth

# External error (502)
curl http://localhost:8081/error/external?service=payment
```

#### 4. Log Level Override Demo (NEW! ⭐)

Demo fluent API `.Level()` để override log level:

```bash
# ValidationError với warn level (default)
# → Console: ✓, File: ✗
curl http://localhost:8081/error/log-level?level=warn&scenario=validation

# ValidationError với error level (override)
# → Console: ✓, File: ✓
curl http://localhost:8081/error/log-level?level=error&scenario=validation

# AuthError với warn level (default)
curl http://localhost:8081/error/log-level?level=warn&scenario=auth

# AuthError với error level (multiple failed attempts)
curl http://localhost:8081/error/log-level?level=error&scenario=auth

# BusinessError với error level (default)
curl http://localhost:8081/error/log-level?level=error&scenario=business

# BusinessError downgrade to warn
curl http://localhost:8081/error/log-level?level=warn&scenario=business
```

**💡 Tip**: Sau khi gọi các endpoints, check:
1. **Console output** - Sẽ thấy TẤT CẢ errors
2. **logs/errors.log** - Chỉ thấy errors nghiêm trọng (error, panic)

## Kiểm tra Logs

### Console Logs

Console sẽ hiển thị tất cả errors (warn, error, panic):

```bash
# Watch console output khi chạy server
go run main.go
```

### File Logs

File `logs/errors.log` chỉ chứa errors nghiêm trọng:

```bash
# View formatted JSON logs
cat logs/errors.log | jq

# Count errors in file
wc -l logs/errors.log

# Clear logs để test lại
rm logs/errors.log
```

### So sánh Console vs File

```bash
# Terminal 1: Run server
go run main.go

# Terminal 2: Trigger validation error
curl http://localhost:8081/error/validation?age=15
# → Check console: ✓ (có log)
# → Check file: ✗ (không có)

# Terminal 2: Trigger system error
curl http://localhost:8081/error/system
# → Check console: ✓ (có log)
# → Check file: ✓ (có log)
```

## Code Examples

### Default Log Levels

```go
// ValidationError → level: "warn" (console only)
goerrorkit.NewValidationError("Email không hợp lệ", nil)

// AuthError → level: "warn" (console only)
goerrorkit.NewAuthError(401, "Unauthorized")

// SystemError → level: "error" (console + file)
goerrorkit.NewSystemError(err)
```

### Override Log Level với .Level()

```go
// Force validation error vào file
goerrorkit.NewValidationError("Suspicious input", nil).
    Level("error")  // Override: warn → error

// Multiple failed login attempts
goerrorkit.NewAuthError(401, "Brute force detected").
    Level("error").  // Override: warn → error
    WithData(map[string]interface{}{
        "attempts": 5,
        "ip": "192.168.1.100",
    })

// Chain với các methods khác
goerrorkit.NewSystemError(err).
    WithData(map[string]interface{}{"db": "postgres"}).
    Level("error").
    WithCallChain()
```

