# GoErrorKit

🚀 Thư viện xử lý lỗi cho Go với khả năng **capture chính xác dòng code gây lỗi** và **stack trace chi tiết**.

## ✨ Tính Năng Chính

- ✅ **Panic recovery tự động** - Capture chính xác dòng code gây panic (không phải dòng gọi hàm)
- ✅ **Wrap error dễ dàng** - `Wrap(err)` và `WrapWithMessage(err, msg)` tự động capture stack trace
- ✅ **Stack trace chi tiết** - Full call chain để debug dễ dàng
- ✅ **Framework agnostic** - Hỗ trợ Fiber, Gin, Echo, Chi (adapters)
- ✅ **Nhiều loại error** - Business, System, Validation, Auth, External
- ✅ **Structured logging** - JSON format với full context
- ✅ **Fluent API** - Chain methods dễ dùng: `.WithData().WithCallChain()`

## 📦 Cài Đặt

```bash
go get github.com/techmaster-vietnam/goerrorkit
```

## 🚀 Quick Start

```go
package main

import (
    "github.com/techmaster-vietnam/goerrorkit"
    "github.com/techmaster-vietnam/goerrorkit/adapters/fiber"
    fiberv2 "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/requestid"
)

func main() {
    // 1. Khởi tạo logger
    goerrorkit.InitDefaultLogger()

    // 2. Cấu hình stack trace
    goerrorkit.ConfigureForApplication("github.com/yourname/yourapp")

    // 3. Setup Fiber
    app := fiberv2.New()
    app.Use(requestid.New())
    app.Use(fiber.ErrorHandler()) // Middleware xử lý error

    // 4. Routes
    app.Get("/", homeHandler)
    app.Listen(":3000")
}

func homeHandler(c *fiberv2.Ctx) error {
    return c.JSON(fiberv2.Map{"message": "Hello"})
}
```

## ⚙️ Cấu Hình

### 1. Cấu Hình Logger

```go
goerrorkit.InitLogger(goerrorkit.LoggerOptions{
    ConsoleOutput: true,           // Log ra console (development)
    FileOutput:    true,            // Log ra file (production)
    FilePath:      "logs/app.log", // Đường dẫn file log
    JSONFormat:    true,            // JSON format (dễ parse, search)
    MaxFileSize:   10,              // 10MB/file (tự động rotate)
    MaxBackups:    5,               // Giữ 5 file backup
    MaxAge:        30,              // Giữ log 30 ngày
    LogLevel:      "error",         // Mức log: error, warn, info, debug
})
```

**Giải thích:**
- `ConsoleOutput`: Hiển thị log trên terminal (tốt cho dev)
- `FileOutput`: Lưu log vào file (cần thiết cho production để trace bugs)
- `JSONFormat`: Format JSON giúp dễ parse bằng ELK, Splunk, hoặc grep
- `MaxFileSize`: Kích thước tối đa mỗi file trước khi rotate (tránh file quá lớn)
- `MaxBackups`: Số lượng file backup giữ lại (cân bằng giữa storage và history)
- `MaxAge`: Số ngày giữ log (tự động xóa log cũ)

### 2. Cấu Hình Stack Trace

#### Option 1: Tự động (Khuyên dùng)

```go
// Tự động lọc stack trace CHỈ HIỂN THỊ code của BẠN
goerrorkit.ConfigureForApplication("main")
// hoặc
goerrorkit.ConfigureForApplication("github.com/yourname/myapp")
```

**Giải thích:**
- Tự động include TẤT CẢ packages bắt đầu với `github.com/yourname/myapp`
- Tự động skip runtime code và thư viện bên thứ 3
- Stack trace ngắn gọn, chỉ 5-10 dòng thay vì 50+ dòng

#### Option 2: Thủ công (Advanced)

```go
goerrorkit.SetStackTraceConfig(goerrorkit.StackTraceConfig{
    IncludePackages: []string{
        "github.com/yourname/myapp",  // Chỉ hiện code của app
        "main",                       // Include main package
    },
    SkipPackages: []string{
        "runtime",                    // Bỏ qua Go runtime
        "github.com/gofiber/fiber",   // Bỏ qua Fiber framework
    },
    ShowFullPath: false,              // false: myapp.Handler, true: github.com/user/myapp.Handler
})
```

**Giải thích:**
- `IncludePackages`: Chỉ hiển thị các packages này trong stack trace
- `SkipPackages`: Bỏ qua các packages này (runtime, framework)
- `ShowFullPath`: 
  - `false`: Ngắn gọn → `myapp.Handler`
  - `true`: Đầy đủ → `github.com/user/myapp.Handler`

#### Option 3: Fluent API (Dynamic)

```go
goerrorkit.Configure().
    SkipPackage("internal/metrics").
    SkipPattern(".RequestID.func").
    SkipPattern(".Logger.func").
    ShowFullPath(false).
    Apply()
```

**Giải thích:**
- Dùng khi cần thêm skip patterns động (middleware, telemetry)
- Chain nhiều cấu hình một lúc
- `.Apply()` để áp dụng

## 📝 Các Loại Error & Tình Huống Sử Dụng

### 1. Panic Error - Tự Động Recovery

**Khi nào xảy ra:** Nil pointer, index out of range, type assertion failed, divide by zero

**Ưu điểm:** Middleware tự động bắt và capture **CHÍNH XÁC** dòng code gây panic!

```go
func panicHandler(c *fiber.Ctx) error {
    numbers := []int{1, 2, 3}
    // ⚠️ Panic sẽ xảy ra ở dòng này
    value := numbers[10] // panic: index out of range
    return c.JSON(fiber.Map{"value": value})
}
```

**Log output:**

```json
{
  "level": "error",
  "message": "runtime error: index out of range [10] with length 3",
  "error_type": "PANIC",
  "status_code": 500,
  "function": "main.panicHandler",
  "file": "main.go:87",  // ⭐ CHÍNH XÁC dòng gây panic!
  "call_chain": [
    "main.panicHandler (main.go:87)",
    "github.com/gofiber/fiber/v2.(*App).next (app.go:512)"
  ]
}
```

**Không cần làm gì:** Middleware tự động handle!

---

### 2. Wrap Error - Đóng Gói Go Error

**Khi nào dùng:** Khi có Go error từ thư viện chuẩn hoặc third-party, cần thêm context và stack trace

#### 2.1. `Wrap(err)` - Đơn giản nhất

**Use case:** Wrap nhanh error với message gốc

```go
// Database error
if err := db.Query("SELECT * FROM users"); err != nil {
    return goerrorkit.Wrap(err)
    // → Message: "sql: connection refused"
    // → Tự động capture: file, line, function
}

// JSON parsing error
if err := json.Unmarshal(data, &result); err != nil {
    return goerrorkit.Wrap(err)
    // → Message: "invalid character '}' looking for beginning of value"
}

// File operation error  
if err := os.ReadFile(path); err != nil {
    return goerrorkit.Wrap(err)
    // → Message: "open config.json: no such file or directory"
}
```

#### 2.2. `WrapWithMessage(err, msg)` - Thêm Context

**Use case:** Thêm message mô tả rõ hơn, giữ error gốc trong `cause`

```go
// Database với custom message
if err := db.Query(query); err != nil {
    return goerrorkit.WrapWithMessage(err, "Failed to fetch user list from database")
    // → Message: "Failed to fetch user list from database"
    // → Cause: "sql: connection refused"
}

// Redis với context
if err := redis.Get(key); err != nil {
    return goerrorkit.WrapWithMessage(err, "Failed to get user session from cache")
}

// External API với context
if err := paymentAPI.Charge(amount); err != nil {
    return goerrorkit.WrapWithMessage(err, "Payment processing failed")
}
```

#### 2.3. Wrap + WithData - Thêm Metadata

**Use case:** Cần thêm dữ liệu đặc thù để debug

```go
// Database query với data
if err := db.Query(query); err != nil {
    return goerrorkit.Wrap(err).WithData(map[string]interface{}{
        "query": query,
        "table": "users",
        "timeout": "5s",
    })
}

// Redis với data
if err := redis.Get(key); err != nil {
    return goerrorkit.WrapWithMessage(err, "Cache miss").WithData(map[string]interface{}{
        "key": key,
        "ttl": 3600,
    })
}

// File operation với data
if err := os.ReadFile(path); err != nil {
    return goerrorkit.Wrap(err).WithData(map[string]interface{}{
        "path": path,
        "size": fileSize,
    })
}
```

#### 2.4. Wrap + WithCallChain - Debug Phức Tạp

**Use case:** Lỗi phức tạp, cần trace flow qua nhiều tầng

```go
// Complex operation với full trace
if err := complexDatabaseOperation(); err != nil {
    return goerrorkit.WrapWithMessage(err, "Complex operation failed").
        WithData(map[string]interface{}{
            "operation": "bulk_insert",
            "records": 1000,
        }).
        WithCallChain()
}
```

**Log output:**

```json
{
  "level": "error",
  "message": "Failed to fetch user list from database",
  "error_type": "SYSTEM",
  "status_code": 500,
  "function": "services.GetUsers",
  "file": "user_service.go:45",
  "cause": "sql: connection refused",
  "data": {
    "query": "SELECT * FROM users WHERE active = true",
    "table": "users"
  }
}
```

---

### 3. Business Error (4xx)

**Khi nào dùng:** Lỗi business logic, user có thể fix được

```go
// Tình huống 1: Product không tồn tại
if product == nil {
    return goerrorkit.NewBusinessError(404, "Product not found")
}

// Tình huống 2: Hết hàng (có thêm thông tin chi tiết)
if product.Stock == 0 {
    return goerrorkit.NewBusinessError(400, "Product out of stock").WithData(map[string]interface{}{
        "product_id": productID,
        "stock": 0,
    })
}
```

---

### 4. System Error (5xx)

**Khi nào dùng:** Lỗi hệ thống không mong muốn (khuyến nghị dùng `Wrap()` thay thế)

```go
// Cách cũ (vẫn dùng được)
if err := db.Connect(); err != nil {
    return goerrorkit.NewSystemError(err).WithData(map[string]interface{}{
        "database": "postgres",
        "host": "localhost:5432",
    })
}

// ⭐ Cách mới (khuyên dùng)
if err := db.Connect(); err != nil {
    return goerrorkit.WrapWithMessage(err, "Database connection failed").WithData(map[string]interface{}{
        "database": "postgres",
        "host": "localhost:5432",
    })
}
```

---

### 5. Validation Error (400)

**Khi nào dùng:** Input không hợp lệ, missing fields, wrong format

```go
// Tình huống 1: Single field validation
if age < 18 {
    return goerrorkit.NewValidationError("Age must be >= 18", map[string]interface{}{
        "field": "age",
        "min": 18,
        "received": age,
    })
}

// Tình huống 2: Multiple fields validation
if user.Email == "" || user.Name == "" {
    return goerrorkit.NewValidationError("Missing required fields", map[string]interface{}{
        "required": []string{"email", "name"},
    })
}
```

---

### 6. Auth Error (401, 403)

**Khi nào dùng:** Authentication, authorization issues

```go
// Tình huống 1: Missing token
if token == "" {
    return goerrorkit.NewAuthError(401, "Unauthorized: Missing token")
}

// Tình huống 2: Invalid token
if !isValidToken(token) {
    return goerrorkit.NewAuthError(401, "Invalid token").WithData(map[string]interface{}{
        "token_type": getTokenType(token),
    })
}

// Tình huống 3: Insufficient permissions
if !hasPermission(user, "admin") {
    return goerrorkit.NewAuthError(403, "Forbidden").WithData(map[string]interface{}{
        "required_role": "admin",
        "user_role": user.Role,
    })
}
```

---

### 7. External Error (502-504)

**Khi nào dùng:** Lỗi từ third-party services (khuyến nghị dùng `WrapWithMessage()` thay thế)

```go
// Chỉ định ra đây là ExternalError
if err := paymentGateway.Charge(amount); err != nil {
    return goerrorkit.NewExternalError(502, "Payment gateway unavailable", err).WithData(map[string]interface{}{
        "gateway": "stripe",
        "amount": amount,
    })
}

// Tự động bọc lấy error trả về kiểu SystemError
if err := paymentAPI.Charge(amount); err != nil {
    return goerrorkit.WrapWithMessage(err, "Payment gateway unavailable").WithData(map[string]interface{}{
        "gateway": "stripe",
        "amount": amount,
    })
}
```

## 🔍 WithCallChain() - Debug Chi Tiết

**Mặc định:** Chỉ **panic errors** có full call chain.

**Khi nào dùng `.WithCallChain()`:**
- ✅ Debug lỗi phức tạp qua nhiều tầng function
- ✅ Trace flow trong microservices
- ✅ Investigate production issues
- ✅ Deep call stack cần chi tiết

**Khi nào KHÔNG cần:**
- ❌ Lỗi đơn giản, rõ ràng
- ❌ Performance critical code
- ❌ Log volume quá lớn

### Ví Dụ

```go
func processOrder(orderID string) error {
    if err := validateOrder(orderID); err != nil {
        return err // err đã có WithCallChain()
    }
    
    if err := checkInventory(orderID); err != nil {
        return err // err đã có WithCallChain()
    }
    
    return nil
}

func validateOrder(orderID string) error {
    if orderID == "" {
        // ⭐ Thêm WithCallChain() để trace flow đầy đủ
        return goerrorkit.NewValidationError("Invalid order", map[string]interface{}{
            "reason": "empty_order_id",
        }).WithCallChain()
    }
    return nil
}

func checkInventory(orderID string) error {
    stock := getStock(orderID)
    if stock == 0 {
        // ⭐ Chain với WithData()
        return goerrorkit.NewBusinessError(422, "Out of stock").
            WithData(map[string]interface{}{
                "order_id": orderID,
                "stock": 0,
            }).
            WithCallChain()
    }
    return nil
}
```

### Output So Sánh

**Không có `.WithCallChain()`:**

```json
{
  "level": "error",
  "message": "Order validation failed",
  "function": "main.validateOrder",
  "file": "order.go:45"
}
```

**Có `.WithCallChain()`:**

```json
{
  "level": "error",
  "message": "Order validation failed",
  "function": "main.validateOrder",
  "file": "order.go:45",
  "call_chain": [
    "main.validateOrder (order.go:45)",
    "main.processOrder (order.go:23)",
    "main.handleOrderRequest (handler.go:78)"
  ]
}
```

## 📊 Log Output Examples

### Panic Log (Tự động capture chính xác)

```json
{
  "timestamp": "2025-11-11T10:30:45+07:00",
  "level": "error",
  "message": "Panic recovered: runtime error: index out of range [10] with length 3",
  "error_type": "PANIC",
  "status_code": 500,
  "path": "GET /panic",
  "request_id": "abc-123-def-456",
  "function": "main.GetElement",
  "file": "main.go:94",
  "call_chain": [
    "main.panicHandler (main.go:87)",
    "main.errorHandler (main.go:102)"
  ],
  "panic_value": "runtime error: index out of range [10] with length 3"
}
```

**Lưu ý:** `file: "main.go:94"` là **CHÍNH XÁC** dòng gây panic!

### Validation Error với Data

```json
{
  "timestamp": "2025-11-11T15:58:00+07:00",
  "level": "error",
  "message": "Insufficient stock",
  "error_type": "VALIDATION",
  "status_code": 400,
  "path": "POST /order/create",
  "request_id": "c8e1aa21-9f08-4e73-809b",
  "function": "services.ReserveProduct",
  "file": "product_service.go:70",
  "data": {
    "product_id": "123",
    "product_name": "iPhone 15",
    "requested": 1,
    "available_stock": 0
  }
}
```

**Ưu điểm:** Dữ liệu đặc thù nằm trong trường `data` riêng biệt, dễ đọc và phân tích!

## 🎯 So Sánh Với Các Thư Viện Khác

| Feature | GoErrorKit | pkg/errors | cockroachdb/errors | Sentry |
|---------|------------|------------|-------------------|--------|
| Chính xác panic location | ✅ main.go:94 | ❌ Tại wrap | ❌ Tại wrap | ✅ |
| Call chain đầy đủ | ✅ | ⚠️ Partial | ⚠️ Partial | ✅ |
| Log vào file local | ✅ JSON | ❌ | ❌ | ❌ |
| Framework agnostic | ✅ | ✅ | ✅ | ✅ |
| Self-hosted | ✅ | ✅ | ✅ | ⚠️ Optional |
| Zero external service | ✅ | ✅ | ✅ | ❌ |

## 🏗️ Architecture

```
goerrorkit/
├── *.go               # Core library (framework-agnostic)
│   ├── error.go       # Error types & factories
│   ├── handler.go     # Panic handling & conversion
│   ├── stacktrace.go  # Stack trace capture & filtering
│   ├── logger.go      # Logging interface
│   └── context.go     # HTTP context interface
│
├── adapters/          # Framework adapters
│   └── fiber/         # Fiber v2 adapter
│
└── examples/          # Demo apps
    └── fiber-demo/
```

## 🔌 Framework Adapters

**Supported:**
- ✅ **Fiber v2** - `github.com/techmaster-vietnam/goerrorkit/adapters/fiber`

**Coming Soon:**
- 🚧 **Gin**
- 🚧 **Echo**
- 🚧 **Chi**

## 📚 Documentation

- [Getting Started](docs/getting-started.md)
- [Configuration Guide](docs/configuration.md)
- [Stack Trace Configuration](docs/stack-trace-configuration.md)

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details.

