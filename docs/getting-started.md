# Getting Started với GoErrorKit

## Installation

```bash
go get github.com/techmaster-vietnam/goerrorkit
```

## Basic Setup

### Bước 1: Khởi tạo Logger

```go
package main

import (
    "github.com/techmaster-vietnam/goerrorkit"
)

func main() {
    // Option 1: Sử dụng config mặc định
    goerrorkit.InitDefaultLogger()
    
    // Option 2: Custom config
    goerrorkit.InitLogger(goerrorkit.LoggerOptions{
        ConsoleOutput: true,
        FileOutput:    true,
        FilePath:      "logs/errors.log",
        JSONFormat:    true,
        MaxFileSize:   10,  // MB
        MaxBackups:    5,
        MaxAge:        30,  // days
        LogLevel:      "error",
    })
}
```

### Bước 2: Cấu hình Stack Trace

```go
import "github.com/techmaster-vietnam/goerrorkit"

func main() {
    // ...
    
    // Configure stack trace để chỉ show application code
    goerrorkit.ConfigureForApplication("github.com/yourname/yourapp")
}
```

### Bước 3: Setup Framework Middleware

#### Với Fiber

```go
import (
    "github.com/techmaster-vietnam/goerrorkit/adapters/fiber"
    fiberv2 "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/requestid"
)

func main() {
    app := fiberv2.New()
    
    // RequestID middleware (must be before ErrorHandler)
    app.Use(requestid.New())
    
    // GoErrorKit ErrorHandler middleware
    app.Use(fiber.ErrorHandler())
    
    // Your routes...
    app.Get("/", handler)
    
    app.Listen(":3000")
}
```

## Sử dụng Error Types

### Business Error

```go
func handler(c *fiber.Ctx) error {
    product := findProduct(id)
    if product == nil {
        return goerrorkit.NewBusinessError(404, "Product not found")
    }
    return c.JSON(product)
}
```

### System Error

```go
func handler(c *fiber.Ctx) error {
    err := db.Connect()
    if err != nil {
        return goerrorkit.NewSystemError(err)
    }
    // ...
}
```

### Validation Error

```go
func handler(c *fiber.Ctx) error {
    age := c.Query("age")
    if age < 18 {
        return goerrorkit.NewValidationError("Age must be >= 18", map[string]interface{}{
            "field": "age",
            "min": 18,
            "received": age,
        })
    }
    // ...
}
```

## Complete Example

```go
package main

import (
    "github.com/techmaster-vietnam/goerrorkit"
    "github.com/techmaster-vietnam/goerrorkit/adapters/fiber"
    fiberv2 "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/requestid"
)

func main() {
    // 1. Init logger
    goerrorkit.InitDefaultLogger()
    
    // 2. Configure stack trace
    goerrorkit.ConfigureForApplication("main")
    
    // 3. Setup Fiber
    app := fiberv2.New()
    app.Use(requestid.New())
    app.Use(fiber.ErrorHandler())
    
    // 4. Routes
    app.Get("/", func(c *fiberv2.Ctx) error {
        return c.JSON(fiber.Map{"status": "ok"})
    })
    
    app.Get("/panic", func(c *fiberv2.Ctx) error {
        panic("test panic") // Will be auto-recovered
    })
    
    app.Get("/error", func(c *fiberv2.Ctx) error {
        return goerrorkit.NewBusinessError(404, "Not found")
    })
    
    app.Listen(":3000")
}
```

## Next Steps

- [Configuration Guide](configuration.md) - Chi tiết về cấu hình
- [Architecture Overview](architecture.md) - Hiểu cách thư viện hoạt động
- [Creating Custom Adapters](custom-adapters.md) - Tạo adapter cho framework khác

