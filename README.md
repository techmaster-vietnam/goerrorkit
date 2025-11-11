# GoErrorKit

🚀 **Framework-agnostic error handling library for Go** với khả năng capture chính xác **dòng code gây lỗi** và **stack trace chi tiết**.

## ✨ Features

- ✅ **Panic recovery tự động** - Capture chính xác dòng code gây panic (không phải dòng gọi)
- ✅ **Stack trace chi tiết** - Full call chain đến từng function
- ✅ **Framework agnostic** - Core logic hoàn toàn độc lập với web framework
- ✅ **Multiple framework support** - Adapters cho Fiber, Gin, Echo, Chi (coming soon)
- ✅ **Custom error types** - Business, System, Validation, Auth, External errors
- ✅ **Structured logging** - JSON format với full context
- ✅ **Tách biệt metadata và data** - Trường `data` riêng cho dữ liệu đặc thù, giúp log dễ đọc
- ✅ **File logging với rotation** - Tích hợp lumberjack
- ✅ **Caller info tracking** - Tự động capture file:line cho mọi error
- ✅ **Configurable** - Customize stack trace filtering, logger, etc.

## 📦 Installation

```bash
go get github.com/techmaster-vietnam/goerrorkit
```

## 🚀 Quick Start

### 1. Basic Setup với Fiber

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

    // 2. Cấu hình stack trace cho application
    goerrorkit.ConfigureForApplication("github.com/yourname/yourapp")

    // 3. Setup Fiber app
    app := fiberv2.New()
    
    // RequestID middleware (để track requests)
    app.Use(requestid.New())
    
    // GoErrorKit middleware (PHẢI sau requestid)
    app.Use(fiber.ErrorHandler())

    // 4. Define routes
    app.Get("/", homeHandler)
    app.Get("/panic", panicHandler)
    app.Get("/error", errorHandler)

    app.Listen(":3000")
}

func homeHandler(c *fiberv2.Ctx) error {
    return c.JSON(fiber.Map{"message": "Hello World"})
}

func panicHandler(c *fiberv2.Ctx) error {
    // Panic sẽ được tự động catch với CHÍNH XÁC location
    arr := []int{1, 2, 3}
    return c.JSON(fiber.Map{"value": arr[10]}) // ← Stack trace sẽ trỏ chính xác dòng này!
}

func errorHandler(c *fiberv2.Ctx) error {
    // Custom error với stack trace
    return goerrorkit.NewBusinessError(404, "Resource not found")
}
```

### 2. Custom Logger Configuration

```go
goerrorkit.InitLogger(goerrorkit.LoggerOptions{
    ConsoleOutput: true,           // Log ra console
    FileOutput:    true,            // Log ra file
    FilePath:      "logs/app.log", // Đường dẫn file
    JSONFormat:    true,            // JSON format
    MaxFileSize:   10,              // 10MB per file
    MaxBackups:    5,               // Keep 5 backup files
    MaxAge:        30,              // 30 days
    LogLevel:      "error",         // error, warn, info, debug
})
```

### 3. Stack Trace Configuration

```go
// Option 1: Auto-configure cho application package
goerrorkit.ConfigureForApplication("github.com/yourname/myapp")

// Option 2: Manual configuration
goerrorkit.SetStackTraceConfig(goerrorkit.StackTraceConfig{
    IncludePackages: []string{
        "github.com/yourname/myapp",
        "main", // for local development
    },
    SkipPackages: []string{
        "runtime",
        "runtime/debug",
    },
    ShowFullPath: false, // true: full path, false: short name
})
```

## 📝 Error Types & Usage

### Business Error (4xx)

```go
// Product không tồn tại (không cần data)
if product == nil {
    return goerrorkit.NewBusinessError(404, "Product not found")
}

// Hết hàng (với custom data)
if product.Stock == 0 {
    return goerrorkit.NewBusinessError(400, "Product out of stock").WithData(map[string]interface{}{
        "product_id": productID,
        "stock": 0,
    })
}
```

### System Error (5xx)

```go
// Database error (với custom data)
if err := db.Connect(); err != nil {
    return goerrorkit.NewSystemError(err).WithData(map[string]interface{}{
        "database": "postgres",
    })
}

// File system error (không cần data)
if err := os.ReadFile("config.json"); err != nil {
    return goerrorkit.NewSystemError(err)
}
```

### Validation Error (400)

```go
// Single field validation
if age < 18 {
    return goerrorkit.NewValidationError("Age must be >= 18", map[string]interface{}{
        "field":    "age",
        "min":      18,
        "received": age,
    })
}

// Multiple field validation
if user.Email == "" || user.Name == "" {
    return goerrorkit.NewValidationError("Missing required fields", map[string]interface{}{
        "required": []string{"email", "name"},
    })
}

// Thêm dữ liệu đặc thù với .WithData() (fluent API)
if stock < requested {
    return goerrorkit.NewBusinessError(400, "Insufficient stock").WithData(map[string]interface{}{
        "product_id": productID,
        "requested": requested,
        "available": stock,
    })
}
```

**Lưu ý:** 
- Validation error thường cần data → truyền trực tiếp vào parameter
- Các error khác thường không cần → dùng `.WithData()` khi cần
- Dữ liệu được log trong trường `data` riêng biệt, tách biệt với metadata hệ thống

### Auth Error (401, 403)

```go
// Missing token (không cần data)
if token == "" {
    return goerrorkit.NewAuthError(401, "Unauthorized: Missing token")
}

// Invalid token (với custom data)
if !isValidToken(token) {
    return goerrorkit.NewAuthError(401, "Unauthorized: Invalid token").WithData(map[string]interface{}{
        "token_type": getTokenType(token),
    })
}

// Insufficient permissions (với custom data)
if !hasPermission(user, "admin") {
    return goerrorkit.NewAuthError(403, "Forbidden: Insufficient permissions").WithData(map[string]interface{}{
        "user_id": user.ID,
        "required_role": "admin",
    })
}
```

### External Error (502-504)

```go
// Payment gateway error (với custom data)
if err := paymentGateway.Charge(amount); err != nil {
    return goerrorkit.NewExternalError(502, "Payment gateway unavailable", err).WithData(map[string]interface{}{
        "gateway": "stripe",
        "amount": amount,
    })
}

// Third-party API timeout (với custom data)
if err := apiClient.Call(); err != nil {
    return goerrorkit.NewExternalError(504, "External API timeout", err).WithData(map[string]interface{}{
        "api_endpoint": "/users",
        "timeout": "30s",
    })
}
```

## 📊 Log Output Examples

### Panic Log

Khi panic xảy ra, bạn sẽ nhận được log chi tiết như sau:

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

**Chú ý:** `file: "main.go:94"` là **CHÍNH XÁC** dòng code gây panic, không phải dòng gọi hàm!

### Validation Error với Data

Khi có validation error với custom data:

```json
{
  "timestamp": "2025-11-11T15:58:00+07:00",
  "level": "error",
  "message": "Không đủ hàng: yêu cầu 1, còn lại 0",
  "error_type": "VALIDATION",
  "status_code": 400,
  "path": "POST /order/create",
  "request_id": "c8e1aa21-9f08-4e73-809b-f3937266fe22",
  "function": "services.(*ProductService).ReserveProduct",
  "file": "product_service.go:70",
  "data": {
    "product_id": "123",
    "product_name": "iPhone 15",
    "requested": 1,
    "available_stock": 0
  }
}
```

**Ưu điểm:** Dữ liệu đặc thù được nhóm trong trường `data`, tách biệt với metadata hệ thống, giúp log dễ đọc và phân tích hơn!

## 🎯 Comparison với các thư viện khác

| Feature | GoErrorKit | pkg/errors | cockroachdb/errors | Sentry |
|---------|------------|------------|-------------------|--------|
| **Chính xác panic location** | ✅ main.go:94 | ❌ Capture tại wrap | ❌ Capture tại wrap | ✅ |
| **Call chain đầy đủ** | ✅ | ⚠️ Partial | ⚠️ Partial | ✅ |
| **Log vào file local** | ✅ JSON | ❌ | ❌ | ❌ |
| **Framework agnostic** | ✅ | ✅ | ✅ | ✅ |
| **Self-hosted** | ✅ | ✅ | ✅ | ⚠️ Optional |
| **Zero external service** | ✅ | ✅ | ✅ | ❌ |
| **Setup complexity** | Low | Low | Low | Medium |

## 🏗️ Architecture

```
goerrorkit/
├── *.go               # Core library (framework-agnostic)
│   ├── error.go       # Error types & factories
│   ├── handler.go     # Panic handling & conversion
│   ├── stacktrace.go  # Stack trace capture & filtering
│   ├── context.go     # HTTP context interface
│   ├── logger.go      # Logging interface
│   └── logrus_logger.go # Logrus logger implementation
│
├── adapters/          # Framework-specific adapters
│   └── fiber/         # Fiber v2 adapter
│       ├── middleware.go
│       └── context.go
│
└── examples/          # Example applications
    └── fiber-demo/
```

## 🔌 Adapters

### Currently Supported

- ✅ **Fiber v2** - `github.com/techmaster-vietnam/goerrorkit/adapters/fiber`

### Coming Soon

- 🚧 **Gin** - `github.com/techmaster-vietnam/goerrorkit/adapters/gin`
- 🚧 **Echo** - `github.com/techmaster-vietnam/goerrorkit/adapters/echo`
- 🚧 **Chi** - `github.com/techmaster-vietnam/goerrorkit/adapters/chi`

## 📚 Documentation

- [Getting Started](docs/getting-started.md)
- [Configuration Guide](docs/configuration.md)
- [Architecture Overview](docs/architecture.md)
- [Creating Custom Adapters](docs/custom-adapters.md)

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Inspired by the need for accurate panic location tracking in production Go applications
- Built with ❤️ for the Go community

## 📧 Contact

- GitHub: [@cuong](https://github.com/cuong)
- Email: your.email@example.com

---

⭐ If you find this library helpful, please consider giving it a star on GitHub!

