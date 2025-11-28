# GoErrorKit

🚀 Thư viện xử lý lỗi cho Go với khả năng **capture chính xác dòng code gây lỗi** và **stack trace chi tiết**.

## ✨ Tính Năng Chính

- ✅ **Panic recovery tự động** - Capture chính xác dòng code gây panic
- ✅ **Wrap error dễ dàng** - `Wrap(err)` và `WrapWithMessage(err, msg)` 
- ✅ **Stack trace thông minh** - Lọc chỉ hiện code của bạn
- ✅ **Framework agnostic** - Hỗ trợ Fiber, Gin, Echo, Chi
- ✅ **Dual-level logging** - Console và File với mức độ khác nhau
- 🚀 **Build modes** - Debug/trace logs tự động loại bỏ trong production (zero overhead)

## 📦 Cài Đặt

```bash
go get github.com/techmaster-vietnam/goerrorkit
```

## 🚀 Quick Start

```go
package main

import (
    "github.com/techmaster-vietnam/goerrorkit"
    fiberv2 "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/requestid"
)

func main() {
    // 1. Khởi tạo logger
    goerrorkit.InitLogger(goerrorkit.LoggerOptions{
        ConsoleOutput: true,
        FileOutput:    true,
        FilePath:      "logs/app.log",
        JSONFormat:    true,
        LogLevel:      "warn",   // Console: warn trở lên
        FileLogLevel:  "error",  // File: chỉ error và panic
    })

    // 2. Cấu hình stack trace (tự động lọc runtime & thư viện)
    goerrorkit.ConfigureForApplication("github.com/yourname/yourapp")

    // 3. Setup Fiber với error handler
    app := fiberv2.New()
    app.Use(requestid.New())
    app.Use(goerrorkit.FiberErrorHandler())

    // 4. Routes
    app.Get("/users/:id", getUserHandler)
    app.Listen(":3000")
}

func getUserHandler(c *fiberv2.Ctx) error {
    // Validation error
    if c.Params("id") == "" {
        return goerrorkit.NewValidationError("Missing user ID", nil)
    }
    
    // Wrap database error
    user, err := db.GetUser(id)
    if err != nil {
        return goerrorkit.WrapWithMessage(err, "Failed to fetch user")
    }
    
    return c.JSON(user)
}
```

## 📊 Phân Loại Error Types

| Error Type | HTTP Code | Default Log Level | Console | File | Khi Nào Dùng |
|------------|-----------|-------------------|---------|------|--------------|
| **ValidationError** | 400 | `warn` | ✅ | ❌ | Input không hợp lệ, missing fields |
| **AuthError** | 401, 403 | `warn` | ✅ | ❌ | Authentication/authorization failed |
| **BusinessError** | 404, 422 | `error` | ✅ | ✅ | Business logic errors |
| **SystemError** | 500 | `error` | ✅ | ✅ | Database, file system errors |
| **ExternalError** | 502-504 | `error` | ✅ | ✅ | Third-party service errors |
| **PanicError** | 500 | `panic` | ✅ | ✅ | Runtime panic (tự động bắt) |

**💡 Giải thích:**
- **ValidationError/AuthError**: Không nghiêm trọng → Console only (không làm nhiễu file log)
- **Business/System/External**: Nghiêm trọng → Ghi cả console và file
- **PanicError**: Tự động recovery bởi middleware

## 📈 Phân Loại Log Levels

```
trace < debug < info < warn < error < panic
  ↓       ↓      ↓      ↓       ↓       ↓
(dev)   (dev)  (all)  (all)   (all)   (all)
```

| Level | Khi Nào Dùng | Console | File | Build Mode |
|-------|--------------|---------|------|------------|
| **trace** | Track flow trong dev, chi tiết nhất | ✅ | ❌ | Chỉ `-tags=debug` |
| **debug** | Debug biến, state trong dev | ✅ | ❌ | Chỉ `-tags=debug` |
| **info** | Thông tin chung, normal operations | ✅ | ✅ | Mọi build |
| **warn** | Cảnh báo không nghiêm trọng | ✅ | ❌ | Mọi build |
| **error** | Lỗi nghiêm trọng cần investigate | ✅ | ✅ | Mọi build |
| **panic** | Critical errors, system crash | ✅ | ✅ | Mọi build |

### Build Modes

```bash
# Development (trace/debug hoạt động)
go run -tags=debug main.go

# Production (trace/debug bị loại bỏ - zero overhead)
go run main.go
```

**💡 Lưu ý:**
- Trace/Debug chỉ hoạt động khi build với `-tags=debug`
- Production build: Trace/Debug là **no-op** (zero overhead)
- Cấu hình `LogLevel: "trace"` trong dev + `-tags=debug` để log tất cả

## 🎯 Cú Pháp Sử Dụng

### 1. Tạo Error Types

```go
// ValidationError (400, log level: warn)
return goerrorkit.NewValidationError("Age must be >= 18", map[string]interface{}{
    "field": "age",
    "min": 18,
    "received": 15,
})

// AuthError (401/403, log level: warn)
return goerrorkit.NewAuthError(401, "Unauthorized: Invalid token")

// BusinessError (404/422, log level: error)
return goerrorkit.NewBusinessError(404, "Product not found")

// SystemError (500, log level: error) - DEPRECATED, dùng Wrap() thay thế
return goerrorkit.NewSystemError(err)

// ExternalError (502-504, log level: error) - DEPRECATED, dùng WrapWithMessage() thay thế
return goerrorkit.NewExternalError(502, "Payment gateway unavailable", err)
```

### 2. Wrap Error (Khuyên Dùng)

```go
// Wrap() - Giữ nguyên message gốc
if err := db.Query(); err != nil {
    return goerrorkit.Wrap(err)
    // Message: "sql: connection refused"
    // Tự động capture: file, line, function
}

// WrapWithMessage() - Thêm context message
if err := redis.Get(key); err != nil {
    return goerrorkit.WrapWithMessage(err, "Failed to get user session")
    // Message: "Failed to get user session"
    // Cause: "redis: connection timeout"
}
```

### 3. Chain Methods - Bổ Sung Metadata

```go
// WithData() - Thêm dữ liệu debug
return goerrorkit.Wrap(err).WithData(map[string]interface{}{
    "user_id": 123,
    "query": "SELECT * FROM users",
})

// WithCallChain() - Thêm full call stack (dùng cho debug phức tạp)
return goerrorkit.NewBusinessError(422, "Out of stock").
    WithData(map[string]interface{}{
        "product_id": "PROD-123",
        "stock": 0,
    }).
    WithCallChain()

// Level() - Override log level
return goerrorkit.NewValidationError("Suspicious input", nil).
    Level("error")  // Warn → Error (ghi vào file)

// Chain tất cả
return goerrorkit.WrapWithMessage(err, "Complex operation failed").
    WithData(map[string]interface{}{"operation": "bulk_insert"}).
    WithCallChain().
    Level("error")
```

### 4. Logging Trực Tiếp

```go
// Error logging
goerrorkit.Error("Database query failed", map[string]interface{}{
    "query": sql,
    "duration": "5s",
})

// Warning
goerrorkit.Warn("Slow query detected", map[string]interface{}{
    "duration": "2.5s",
})

// Info
goerrorkit.Info("User logged in", map[string]interface{}{
    "user_id": 123,
})

// Debug (chỉ hoạt động với -tags=debug)
goerrorkit.Debug("Processing payment", map[string]interface{}{
    "amount": 10000,
    "gateway": "stripe",
})

// Trace (chỉ hoạt động với -tags=debug)
goerrorkit.Trace("Fetching user from database", map[string]interface{}{
    "user_id": 123,
})
```

## 📋 Bảng Tổng Hợp Cú Pháp

### Error Creation

| Cú Pháp | Use Case | HTTP Code | Log Level |
|---------|----------|-----------|-----------|
| `NewValidationError(msg, data)` | Input không hợp lệ | 400 | warn |
| `NewAuthError(code, msg)` | Auth failed | 401/403 | warn |
| `NewBusinessError(code, msg)` | Business logic | 4xx | error |
| `Wrap(err)` | ⭐ Wrap Go error | 500 | error |
| `WrapWithMessage(err, msg)` | ⭐ Wrap + context | 500 | error |

### Error Enhancement

| Method | Mục Đích | Example |
|--------|----------|---------|
| `.WithData(map)` | Thêm debug data | `.WithData(map[string]interface{}{"user_id": 123})` |
| `.WithCallChain()` | Thêm full stack trace | `.WithCallChain()` |
| `.Level(level)` | Override log level | `.Level("error")` |

### Direct Logging

| Method | Log Level | Build Mode | File Output |
|--------|-----------|------------|-------------|
| `goerrorkit.Trace(msg, fields)` | trace | `-tags=debug` only | ❌ |
| `goerrorkit.Debug(msg, fields)` | debug | `-tags=debug` only | ❌ |
| `goerrorkit.Info(msg, fields)` | info | All | ✅ (if FileLogLevel <= info) |
| `goerrorkit.Warn(msg, fields)` | warn | All | ❌ (nếu FileLogLevel=error) |
| `goerrorkit.Error(msg, fields)` | error | All | ✅ |
| `goerrorkit.Panic(msg, fields)` | panic | All | ✅ |

## ⚙️ Cấu Hình Logger

### Dual-Level Logging

```go
goerrorkit.InitLogger(goerrorkit.LoggerOptions{
    ConsoleOutput: true,            // Log ra console (development)
    FileOutput:    true,            // Log ra file (production)
    FilePath:      "logs/app.log",  // Đường dẫn file log
    JSONFormat:    true,            // JSON format (dễ parse)
    MaxFileSize:   10,              // 10MB/file (auto rotate)
    MaxBackups:    5,               // Giữ 5 file backup
    MaxAge:        30,              // Giữ log 30 ngày
    LogLevel:      "warn",          // Console: log từ warn trở lên
    FileLogLevel:  "error",         // File: CHỈ log error/panic
})
```

**Ưu điểm Dual-Level:**
- Console: Log tất cả (warn, error) để developer debug
- File: Chỉ log nghiêm trọng (error, panic) → File log sạch sẽ, dễ phân tích

### Stack Trace Configuration

```go
// Tự động lọc stack trace (khuyên dùng)
goerrorkit.ConfigureForApplication("github.com/yourname/yourapp")
// → Chỉ hiện code của bạn, bỏ qua runtime & thư viện

// Hoặc fluent API
goerrorkit.Configure().
    SkipPackage("internal/metrics").
    SkipPattern(".RequestID.func").
    SkipPattern(".Logger.func").
    ShowFullPath(false).
    Apply()
```

## 📝 Ví Dụ Chi Tiết

### Example 1: Validation với Override Level

```go
func validateAge(age int) error {
    // Normal validation (log level: warn)
    if age < 18 {
        return goerrorkit.NewValidationError("Age must be >= 18", map[string]interface{}{
            "field": "age",
            "min": 18,
            "received": age,
        })
    }
    
    // Suspicious input (override to error level)
    if age > 150 {
        return goerrorkit.NewValidationError("Suspicious age detected", map[string]interface{}{
            "field": "age",
            "received": age,
            "reason": "possible_attack",
        }).Level("error")  // ⭐ Ghi vào file
    }
    
    return nil
}
```

### Example 2: Wrap Database Error

```go
func getUser(id string) (*User, error) {
    user := &User{}
    
    // Wrap với context message
    if err := db.Get(user, id); err != nil {
        return nil, goerrorkit.WrapWithMessage(err, "Failed to fetch user").
            WithData(map[string]interface{}{
                "user_id": id,
                "table": "users",
            })
    }
    
    return user, nil
}
```

### Example 3: Complex Flow với Call Chain

```go
func processOrder(orderID string) error {
    // Validate
    if err := validateOrder(orderID); err != nil {
        return err  // err đã có WithCallChain()
    }
    
    // Check inventory
    if err := checkInventory(orderID); err != nil {
        return err  // err đã có WithCallChain()
    }
    
    return nil
}

func validateOrder(orderID string) error {
    if orderID == "" {
        return goerrorkit.NewValidationError("Invalid order", nil).
            WithCallChain()  // ⭐ Thêm full call stack
    }
    return nil
}

func checkInventory(orderID string) error {
    stock := getStock(orderID)
    if stock == 0 {
        return goerrorkit.NewBusinessError(422, "Out of stock").
            WithData(map[string]interface{}{
                "order_id": orderID,
                "stock": 0,
            }).
            WithCallChain()  // ⭐ Trace full flow
    }
    return nil
}
```

### Example 4: Debug Logging (Development Only)

```go
func processPayment(amount int) error {
    // Trace flow (chỉ hoạt động với -tags=debug)
    goerrorkit.Trace("Payment processing started", map[string]interface{}{
        "amount": amount,
        "gateway": "stripe",
    })
    
    // Debug detailed state
    goerrorkit.Debug("Validating payment", map[string]interface{}{
        "amount": amount,
        "currency": "VND",
        "customer_id": "cust_123",
    })
    
    // Process payment...
    
    goerrorkit.Trace("Payment completed", map[string]interface{}{
        "status": "success",
        "transaction_id": "txn_456",
    })
    
    return nil
}
```

## 📊 Log Output

### Panic Error (Tự Động Capture)

```json
{
  "timestamp": "2025-11-28T10:30:45+07:00",
  "level": "error",
  "message": "Panic recovered: index out of range [10] with length 3",
  "error_type": "PANIC",
  "status_code": 500,
  "path": "/users/123",
  "function": "main.GetElement",
  "file": "main.go:94",
  "call_chain": [
    "main.GetElement (main.go:94)",
    "main.getUserHandler (main.go:87)"
  ]
}
```

### Wrapped Error với Data

```json
{
  "timestamp": "2025-11-28T10:30:45+07:00",
  "level": "error",
  "message": "Failed to fetch user",
  "error_type": "SYSTEM",
  "status_code": 500,
  "function": "services.GetUser",
  "file": "user_service.go:45",
  "cause": "sql: connection refused",
  "data": {
    "user_id": "123",
    "table": "users"
  }
}
```

## 🎯 Best Practices

### ✅ DO

```go
// 1. Dùng Wrap() cho Go errors
if err := db.Query(); err != nil {
    return goerrorkit.Wrap(err)
}

// 2. Thêm context với WrapWithMessage()
if err := redis.Get(key); err != nil {
    return goerrorkit.WrapWithMessage(err, "Failed to get cache")
}

// 3. Thêm debug data với WithData()
return goerrorkit.Wrap(err).WithData(map[string]interface{}{
    "query": sql,
})

// 4. Dùng WithCallChain() cho debug phức tạp
return goerrorkit.NewBusinessError(422, "Out of stock").
    WithData(data).
    WithCallChain()

// 5. Override level khi cần
return goerrorkit.NewValidationError("Suspicious input", nil).
    Level("error")
```

### ❌ DON'T

```go
// 1. KHÔNG tạo SystemError khi có thể dùng Wrap()
// BAD
return goerrorkit.NewSystemError(err)
// GOOD
return goerrorkit.Wrap(err)

// 2. KHÔNG dùng WithCallChain() cho mọi error (overhead)
// BAD
return goerrorkit.NewValidationError("Invalid email", nil).WithCallChain()
// GOOD
return goerrorkit.NewValidationError("Invalid email", nil)

// 3. KHÔNG quên cấu hình stack trace
// BAD
// goerrorkit.InitLogger(...) only
// GOOD
goerrorkit.InitLogger(...)
goerrorkit.ConfigureForApplication("yourapp")

// 4. KHÔNG set LogLevel quá thấp trong production
// BAD
LogLevel: "debug"  // Quá nhiều log
// GOOD
LogLevel: "error"  // Chỉ log errors
```

## 🏗️ Architecture

```
goerrorkit/
├── error.go            # Error types & factories
├── handler.go          # Panic handling & conversion
├── stacktrace.go       # Stack trace capture & filtering
├── logger.go           # Logging interface & wrappers
├── context.go          # HTTP context interface
├── adapters/
│   └── fiber/          # Fiber v2 adapter
└── examples/           # Demo apps
```

## 🔌 Framework Adapters

**Supported:**
- ✅ **Fiber v2** - `goerrorkit.FiberErrorHandler()`

**Coming Soon:**
- 🚧 **Gin**
- 🚧 **Echo**
- 🚧 **Chi**

## 📚 Documentation

- [Getting Started](docs/getting-started.md)
- [Configuration Guide](docs/configuration.md)
- [Stack Trace Configuration](docs/stack-trace-configuration.md)
- [Build Modes](docs/build-modes.md)

## 🎯 So Sánh Với Các Thư Viện Khác

| Feature | GoErrorKit | pkg/errors | cockroachdb/errors | Sentry |
|---------|------------|------------|-------------------|--------|
| Panic location chính xác | ✅ | ❌ | ❌ | ✅ |
| Dual-level logging | ✅ | ❌ | ❌ | ❌ |
| Build modes (debug/prod) | ✅ | ❌ | ❌ | ❌ |
| Stack trace filtering | ✅ | ⚠️ | ⚠️ | ✅ |
| Log vào file JSON | ✅ | ❌ | ❌ | ❌ |
| Zero external service | ✅ | ✅ | ✅ | ❌ |
| Self-hosted | ✅ | ✅ | ✅ | ⚠️ |

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details.
