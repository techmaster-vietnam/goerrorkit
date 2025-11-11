# GoErrorKit

🚀 **Framework-agnostic error handling library for Go** với khả năng capture chính xác **dòng code gây lỗi** và **stack trace chi tiết**.

## ✨ Features

- ✅ **Panic recovery tự động** - Capture chính xác dòng code gây panic (không phải dòng gọi)
- ✅ **Stack trace chi tiết** - Full call chain đến từng function
- ✅ **Framework agnostic** - Core logic hoàn toàn độc lập với web framework
- ✅ **Multiple framework support** - Adapters cho Fiber, Gin, Echo, Chi (coming soon)
- ✅ **Custom error types** - Business, System, Validation, Auth, External errors
- ✅ **Structured logging** - JSON format với full context
- ✅ **File logging với rotation** - Tích hợp lumberjack
- ✅ **Caller info tracking** - Tự động capture file:line cho mọi error
- ✅ **Configurable** - Customize stack trace filtering, logger, etc.

## 📦 Installation

```bash
go get github.com/cuong/goerrorkit
```

## 🚀 Quick Start

### 1. Basic Setup với Fiber

```go
package main

import (
    "github.com/cuong/goerrorkit/adapters/fiber"
    "github.com/cuong/goerrorkit/config"
    "github.com/cuong/goerrorkit/core"
    fiberv2 "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/requestid"
)

func main() {
    // 1. Khởi tạo logger
    config.InitDefaultLogger()

    // 2. Cấu hình stack trace cho application
    core.ConfigureForApplication("github.com/yourname/yourapp")

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
    return core.NewBusinessError(404, "Resource not found")
}
```

### 2. Custom Logger Configuration

```go
config.InitLogger(config.LoggerOptions{
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
core.ConfigureForApplication("github.com/yourname/myapp")

// Option 2: Manual configuration
core.SetStackTraceConfig(core.StackTraceConfig{
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
// Product không tồn tại
if product == nil {
    return core.NewBusinessError(404, "Product not found")
}

// Hết hàng
if product.Stock == 0 {
    return core.NewBusinessError(400, "Product out of stock")
}
```

### System Error (5xx)

```go
// Database error
if err := db.Connect(); err != nil {
    return core.NewSystemError(err)
}

// File system error
if err := os.ReadFile("config.json"); err != nil {
    return core.NewSystemError(err)
}
```

### Validation Error (400)

```go
// Single field validation
if age < 18 {
    return core.NewValidationError("Age must be >= 18", map[string]interface{}{
        "field":    "age",
        "min":      18,
        "received": age,
    })
}

// Multiple field validation
if user.Email == "" || user.Name == "" {
    return core.NewValidationError("Missing required fields", map[string]interface{}{
        "required": []string{"email", "name"},
    })
}
```

### Auth Error (401, 403)

```go
// Missing token
if token == "" {
    return core.NewAuthError(401, "Unauthorized: Missing token")
}

// Invalid token
if !isValidToken(token) {
    return core.NewAuthError(401, "Unauthorized: Invalid token")
}

// Insufficient permissions
if !hasPermission(user, "admin") {
    return core.NewAuthError(403, "Forbidden: Insufficient permissions")
}
```

### External Error (502-504)

```go
// Payment gateway error
if err := paymentGateway.Charge(amount); err != nil {
    return core.NewExternalError(502, "Payment gateway unavailable", err)
}

// Third-party API timeout
if err := apiClient.Call(); err != nil {
    return core.NewExternalError(504, "External API timeout", err)
}
```

## 📊 Log Output Example

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
├── core/              # Framework-agnostic core logic
│   ├── error.go       # Error types & factories
│   ├── handler.go     # Panic handling & conversion
│   ├── stacktrace.go  # Stack trace capture & filtering
│   ├── context.go     # HTTP context interface
│   └── logger.go      # Logging interface
│
├── adapters/          # Framework-specific adapters
│   └── fiber/         # Fiber v2 adapter
│       ├── middleware.go
│       └── context.go
│
├── config/            # Configuration
│   └── logger.go      # Logger setup (logrus implementation)
│
└── examples/          # Example applications
    └── fiber-demo/
```

## 🔌 Adapters

### Currently Supported

- ✅ **Fiber v2** - `github.com/cuong/goerrorkit/adapters/fiber`

### Coming Soon

- 🚧 **Gin** - `github.com/cuong/goerrorkit/adapters/gin`
- 🚧 **Echo** - `github.com/cuong/goerrorkit/adapters/echo`
- 🚧 **Chi** - `github.com/cuong/goerrorkit/adapters/chi`

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

