# Stack Trace Configuration

GoErrorKit cung cấp các công cụ mạnh mẽ để lọc và tùy chỉnh stack trace, giúp bạn chỉ hiển thị code quan trọng và loại bỏ các hàm "rác" không cần thiết.

## 📋 Mục Lục

- [Tại Sao Cần Configuration?](#tại-sao-cần-configuration)
- [Cấu Hình Cơ Bản](#cấu-hình-cơ-bản)
- [Fluent API](#fluent-api)
- [Shorthand Functions](#shorthand-functions)
- [Advanced Configuration](#advanced-configuration)
- [Best Practices](#best-practices)

## Tại Sao Cần Configuration?

### ❌ Không có configuration:
```json
{
  "call_chain": [
    "runtime/debug.Stack (stack.go:24)",
    "github.com/techmaster-vietnam/goerrorkit.formatStackTraceArray (stacktrace.go:133)",
    "github.com/techmaster-vietnam/goerrorkit.HandlePanic (handler.go:45)",
    "main.main.New.func1 (requestid.go:31)",
    "github.com/gofiber/fiber/v2.(*App).Next.func1 (app.go:532)",
    "main.panicDivisionHandler (main.go:105)",
    "... 40+ dòng khác ..."
  ]
}
```

### ✅ Với configuration:
```json
{
  "call_chain": [
    "main.panicDivisionHandler (main.go:105)"
  ]
}
```

**Kết quả:** Stack trace ngắn gọn, chỉ hiển thị business logic của bạn!

---

## Cấu Hình Cơ Bản

### 1. ConfigureForApplication (Recommended)

Cách nhanh nhất để bắt đầu:

```go
func main() {
    // Cho application đơn giản (1 package main)
    goerrorkit.ConfigureForApplication("main")
    
    // Hoặc cho application với nhiều packages
    goerrorkit.ConfigureForApplication("github.com/yourname/myapp")
    // → Tự động include TẤT CẢ sub-packages (services/, handlers/, models/...)
}
```

### 2. SetStackTraceConfig (Full Control)

Configuration đầy đủ:

```go
goerrorkit.SetStackTraceConfig(goerrorkit.StackTraceConfig{
    SkipPackages: []string{
        "runtime",
        "runtime/debug",
        "github.com/techmaster-vietnam/goerrorkit",
    },
    SkipFunctions: []string{
        "middleware",
        "wrapper",
        "helper",
    },
    IncludePackages: []string{
        "github.com/yourname/myapp",
    },
    ShowFullPath: false, // true: github.com/user/app.Handler, false: app.Handler
})
```

---

## Fluent API

### Cú Pháp

```go
goerrorkit.Configure().
    SkipPackage("package1").
    SkipFunction("function1").
    Apply()
```

### Ví Dụ Chi Tiết

#### 1. Skip Middleware Anonymous Functions

```go
goerrorkit.Configure().
    SkipPattern(".RequestID.func").      // Fiber RequestID middleware
    SkipPattern(".Logger.func").         // Fiber Logger middleware
    SkipPattern(".Recover.func").        // Fiber Recover middleware
    SkipPattern(".Cors.func").           // Fiber CORS middleware
    Apply()
```

#### 2. Skip Custom Packages

```go
goerrorkit.Configure().
    SkipPackage("internal/telemetry").   // Monitoring code
    SkipPackage("internal/metrics").     // Metrics collection
    SkipPackage("pkg/cache").            // Cache utilities
    Apply()
```

#### 3. Skip Helper Functions

```go
goerrorkit.Configure().
    SkipFunction("wrapper").             // Wrapper functions
    SkipFunction("helper").              // Helper utilities
    SkipFunction("transform").           // Data transformers
    Apply()
```

#### 4. Configuration Phức Tạp

```go
goerrorkit.Configure().
    IncludePackage("github.com/mycompany/myapp").
    SkipPackages("internal/telemetry", "internal/metrics").
    SkipFunctions("wrapper", "helper", "transform").
    SkipPatterns(".RequestID.func", ".Logger.func").
    ShowFullPath(false).
    Apply()
```

### Các Methods Có Sẵn

| Method | Mô Tả | Ví Dụ |
|--------|-------|-------|
| `SkipPackage(pkg)` | Bỏ qua 1 package | `.SkipPackage("runtime")` |
| `SkipPackages(pkgs...)` | Bỏ qua nhiều packages | `.SkipPackages("runtime", "debug")` |
| `SkipFunction(fn)` | Bỏ qua 1 function pattern | `.SkipFunction("middleware")` |
| `SkipFunctions(fns...)` | Bỏ qua nhiều functions | `.SkipFunctions("helper", "wrapper")` |
| `SkipPattern(pattern)` | Alias cho SkipFunction | `.SkipPattern(".func")` |
| `SkipPatterns(patterns...)` | Bỏ qua nhiều patterns | `.SkipPatterns(".func1", ".func2")` |
| `IncludePackage(pkg)` | Include 1 package | `.IncludePackage("main")` |
| `IncludePackages(pkgs...)` | Include nhiều packages | `.IncludePackages("main", "myapp")` |
| `ShowFullPath(bool)` | Hiển thị full path | `.ShowFullPath(true)` |
| `Apply()` | Áp dụng configuration | `.Apply()` |

---

## Shorthand Functions

### AddSkipPatterns

Thêm nhanh các skip patterns mà không cần fluent API:

```go
func main() {
    goerrorkit.ConfigureForApplication("main")
    
    // Thêm các patterns tùy chỉnh
    goerrorkit.AddSkipPatterns(
        ".RequestID.func",
        ".Logger.func",
        ".Telemetry.func",
    )
}
```

### AddSkipPackages

Thêm nhanh các skip packages:

```go
goerrorkit.AddSkipPackages(
    "internal/telemetry",
    "internal/metrics",
    "vendor/monitoring",
)
```

---

## Advanced Configuration

### 1. Middleware-Specific Patterns

Các pattern phổ biến cho Fiber middleware:

```go
goerrorkit.AddSkipPatterns(
    ".main.New.func",        // Fiber app.New() setup
    ".main.Use.func",        // Fiber app.Use() middleware
    ".Next.func",            // Middleware chain
    ".recover.func",         // Recovery middleware
    ".logger.func",          // Logger middleware
    ".requestid.func",       // RequestID middleware
    ".cors.func",            // CORS middleware
)
```

### 2. Third-Party Library Filtering

```go
goerrorkit.Configure().
    SkipPackage("github.com/gofiber/fiber").
    SkipPackage("github.com/sirupsen/logrus").
    SkipPackage("go.uber.org/zap").
    Apply()
```

### 3. Dynamic Configuration

Có thể thay đổi configuration runtime:

```go
func init() {
    goerrorkit.ConfigureForApplication("main")
}

func enableDetailedStackTrace() {
    goerrorkit.Configure().
        ShowFullPath(true).
        Apply()
}

func enableProductionMode() {
    goerrorkit.Configure().
        SkipPackages("debug", "testing").
        ShowFullPath(false).
        Apply()
}
```

### 4. Environment-Based Configuration

```go
import "os"

func setupStackTrace() {
    config := goerrorkit.Configure()
    
    if os.Getenv("ENV") == "production" {
        config.ShowFullPath(false).
               SkipPatterns("debug", "testing")
    } else {
        config.ShowFullPath(true)
    }
    
    config.Apply()
}
```

---

## Best Practices

### ✅ DO

1. **Luôn gọi configuration trước khi app chạy:**
```go
func main() {
    goerrorkit.ConfigureForApplication("main")
    // ... rest of app setup
}
```

2. **Sử dụng ConfigureForApplication cho hầu hết trường hợp:**
```go
// Simple & effective
goerrorkit.ConfigureForApplication("github.com/mycompany/myapp")
```

3. **Thêm patterns cụ thể khi cần:**
```go
goerrorkit.ConfigureForApplication("main")
goerrorkit.AddSkipPatterns(".CustomMiddleware.func")
```

4. **Test stack trace trong development:**
```go
// Trigger một panic để xem stack trace output
panic("test")
```

### ❌ DON'T

1. **Đừng skip quá nhiều - có thể mất thông tin quan trọng:**
```go
// ❌ Too aggressive
goerrorkit.Configure().
    SkipPackages("main", "app", "services"). // Oops! Skip cả business logic
    Apply()
```

2. **Đừng dùng quá nhiều patterns chung chung:**
```go
// ❌ Too broad
goerrorkit.AddSkipPatterns("func", "handler") // Skip quá nhiều
```

3. **Đừng quên gọi Apply():**
```go
// ❌ Configuration không được áp dụng!
goerrorkit.Configure().
    SkipPattern(".middleware.func")
    // Missing .Apply()
```

---

## Examples

### Example 1: Simple Web Server

```go
package main

import "github.com/techmaster-vietnam/goerrorkit"

func main() {
    goerrorkit.ConfigureForApplication("main")
    
    // Your web server code here...
}
```

### Example 2: Microservice với Nhiều Middleware

```go
package main

import "github.com/techmaster-vietnam/goerrorkit"

func main() {
    // Base configuration
    goerrorkit.ConfigureForApplication("github.com/mycompany/myservice")
    
    // Skip common middleware patterns
    goerrorkit.AddSkipPatterns(
        ".RequestID.func",
        ".Logger.func",
        ".Metrics.func",
        ".Tracing.func",
    )
    
    // Skip monitoring packages
    goerrorkit.AddSkipPackages(
        "internal/telemetry",
        "internal/monitoring",
    )
}
```

### Example 3: Advanced với Environment Variables

```go
package main

import (
    "os"
    "github.com/techmaster-vietnam/goerrorkit"
)

func init() {
    setupStackTrace()
}

func setupStackTrace() {
    env := os.Getenv("APP_ENV")
    
    config := goerrorkit.Configure().
        IncludePackage("github.com/mycompany/myapp")
    
    switch env {
    case "production":
        config.ShowFullPath(false).
               SkipPackages("debug", "testing", "internal/dev")
    case "development":
        config.ShowFullPath(true)
    default:
        config.ShowFullPath(false)
    }
    
    config.Apply()
}

func main() {
    // App code...
}
```

---

## Troubleshooting

### Stack trace vẫn quá dài?

1. Kiểm tra xem đã gọi configuration chưa:
```go
goerrorkit.ConfigureForApplication("main")
```

2. Thêm các skip patterns cụ thể:
```go
goerrorkit.AddSkipPatterns(".YourMiddleware.func")
```

3. Log ra stack trace để debug:
```go
// Trong development, tạm thời enable full path
goerrorkit.Configure().ShowFullPath(true).Apply()
```

### Stack trace bị mất thông tin quan trọng?

1. Giảm số lượng skip patterns
2. Kiểm tra IncludePackages có đúng không
3. Tạm thời disable configuration để xem full stack:
```go
// Comment out để xem full stack trace
// goerrorkit.ConfigureForApplication("main")
```

---

## Tham Khảo

- [Getting Started Guide](./getting-started.md)
- [Configuration Guide](./configuration.md)
- [API Reference](../README.md)

