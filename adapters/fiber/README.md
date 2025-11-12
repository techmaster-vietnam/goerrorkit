# Fiber Adapter

Adapter cho [Fiber v2](https://github.com/gofiber/fiber) web framework.

## Cài đặt

```bash
go get github.com/techmaster-vietnam/goerrorkit
go get github.com/gofiber/fiber/v2
```

## Sử dụng

### Cách 1: Import trực tiếp từ package chính (Khuyên dùng)

```go
package main

import (
    "github.com/techmaster-vietnam/goerrorkit"
    fiberv2 "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/requestid"
)

func main() {
    // 1. Khởi tạo logger
    goerrorkit.InitDefaultLogger()

    // 2. Cấu hình stack trace cho application của bạn
    goerrorkit.ConfigureForApplication("github.com/yourname/yourapp")

    // 3. Tạo Fiber app
    app := fiberv2.New()

    // 4. Thêm middleware (RequestID phải trước ErrorHandler)
    app.Use(requestid.New())
    app.Use(goerrorkit.FiberErrorHandler()) // ⭐ Sử dụng trực tiếp từ goerrorkit

    // 5. Định nghĩa routes với error handling tự động
    app.Get("/panic", func(c *fiberv2.Ctx) error {
        // Panic sẽ được tự động catch với chính xác location
        panic("something went wrong")
    })

    app.Get("/error", func(c *fiberv2.Ctx) error {
        // Custom error với stack trace (clean!)
        return goerrorkit.NewBusinessError(404, "Resource not found")
    })
    
    app.Get("/error-with-data", func(c *fiberv2.Ctx) error {
        // Với custom data (dùng .WithData() khi cần)
        return goerrorkit.NewBusinessError(404, "Product not found").WithData(map[string]interface{}{
            "product_id": "123",
        })
    })

    app.Listen(":3000")
}
```

### Cách 2: Import từ adapter (Vẫn hỗ trợ)

```go
package main

import (
    "github.com/techmaster-vietnam/goerrorkit"
    "github.com/techmaster-vietnam/goerrorkit/adapters/fiber"
    fiberv2 "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/requestid"
)

func main() {
    // ...
    app.Use(fiber.ErrorHandler()) // Cách cũ vẫn hoạt động
    // ...
}
```

## Features

- ✅ Tự động recover panic với chính xác dòng code gây lỗi
- ✅ Stack trace chi tiết đến từng hàm trong call chain
- ✅ Tích hợp với Fiber's request ID
- ✅ JSON error response chuẩn
- ✅ Logging tự động với structured fields

