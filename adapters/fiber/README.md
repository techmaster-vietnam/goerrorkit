# Fiber Adapter

Adapter cho [Fiber v2](https://github.com/gofiber/fiber) web framework.

## Cài đặt

```bash
go get github.com/cuong/goerrorkit
go get github.com/gofiber/fiber/v2
```

## Sử dụng

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

    // 2. Cấu hình stack trace cho application của bạn
    core.ConfigureForApplication("github.com/yourname/yourapp")

    // 3. Tạo Fiber app
    app := fiberv2.New()

    // 4. Thêm middleware (RequestID phải trước ErrorHandler)
    app.Use(requestid.New())
    app.Use(fiber.ErrorHandler())

    // 5. Định nghĩa routes với error handling tự động
    app.Get("/panic", func(c *fiberv2.Ctx) error {
        // Panic sẽ được tự động catch với chính xác location
        panic("something went wrong")
    })

    app.Get("/error", func(c *fiberv2.Ctx) error {
        // Custom error với stack trace
        return core.NewBusinessError(404, "Resource not found")
    })

    app.Listen(":3000")
}
```

## Features

- ✅ Tự động recover panic với chính xác dòng code gây lỗi
- ✅ Stack trace chi tiết đến từng hàm trong call chain
- ✅ Tích hợp với Fiber's request ID
- ✅ JSON error response chuẩn
- ✅ Logging tự động với structured fields

