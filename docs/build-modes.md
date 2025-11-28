# Build Modes: Development vs Production

GoErrorKit hỗ trợ 2 chế độ build với behavior khác nhau cho debug/trace logging:

## 🔍 Development Mode (Debug Build)

**Khi nào dùng:**
- Khi đang phát triển ứng dụng
- Khi cần debug chi tiết với trace/debug logs
- Trong môi trường testing/staging

**Cách build:**
```bash
# Build binary với debug mode
go build -tags=debug -o app

# Hoặc run trực tiếp
go run -tags=debug main.go

# Test với debug mode
go test -tags=debug ./...
```

**Behavior:**
- ✅ Debug logs hoạt động bình thường
- ✅ Trace logs hoạt động bình thường
- ✅ Tất cả log levels (trace, debug, info, warn, error, panic) đều hoạt động
- 📊 Output sẽ nhiều hơn, chi tiết hơn

**Example:**
```go
logger := goerrorkit.NewLogrusLogger(goerrorkit.LoggerOptions{
    LogLevel: "debug",  // Sẽ log debug messages
})

// Sẽ in ra console
logger.Debug("Fetching user from database", map[string]interface{}{
    "user_id": 123,
})
```

## 🚀 Production Mode (Default Build)

**Khi nào dùng:**
- Khi build binary cho production
- Khi muốn performance tốt nhất (zero overhead)
- Khi không cần debug/trace logs

**Cách build:**
```bash
# Build binary production (KHÔNG có tag debug)
go build -o app

# Hoặc run
go run main.go

# Test production mode
go test ./...
```

**Behavior:**
- ❌ Debug logs là **no-op** (không làm gì cả, zero overhead)
- ❌ Trace logs là **no-op** (không làm gì cả, zero overhead)
- ✅ Info, warn, error, panic logs vẫn hoạt động bình thường
- 🚀 Performance tốt hơn (không có overhead của debug/trace logging)
- 📦 Binary size nhỏ hơn

**Example:**
```go
logger := goerrorkit.NewLogrusLogger(goerrorkit.LoggerOptions{
    LogLevel: "debug",  // Sẽ KHÔNG log gì (no-op)
})

// KHÔNG in ra gì trong production
logger.Debug("Fetching user from database", map[string]interface{}{
    "user_id": 123,
})
```

## 🎯 Best Practices

### 1. Sử dụng Log Levels đúng mục đích

```go
// ❌ SAI: Dùng debug cho lỗi quan trọng
logger.Debug("Payment failed", fields)

// ✅ ĐÚNG: Dùng error cho lỗi quan trọng
logger.Error("Payment failed", fields)
```

### 2. Cấu hình LogLevel phù hợp với môi trường

```go
// Development
goerrorkit.InitLogger(goerrorkit.LoggerOptions{
    LogLevel: "debug",      // Log mọi thứ
    FileLogLevel: "debug",  // File cũng log debug
})

// Production
goerrorkit.InitLogger(goerrorkit.LoggerOptions{
    LogLevel: "warn",       // Chỉ log warn/error/panic
    FileLogLevel: "error",  // File chỉ log error/panic
})
```

### 3. Sử dụng biến môi trường để switch

```go
package main

import (
    "os"
    "github.com/techmaster-vietnam/goerrorkit"
)

func main() {
    // Đọc environment
    env := os.Getenv("APP_ENV")
    
    var logLevel string
    if env == "production" {
        logLevel = "warn"  // Production: chỉ warn trở lên
    } else {
        logLevel = "debug" // Development: log debug
                           // Nhưng chỉ hoạt động nếu build với -tags=debug
    }
    
    goerrorkit.InitLogger(goerrorkit.LoggerOptions{
        LogLevel: logLevel,
    })
}
```

### 4. Dockerfile cho Production

```dockerfile
# Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .

# Build KHÔNG có -tags=debug để loại bỏ debug logs
RUN go build -o app .

# Runtime stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/app .

# Set environment
ENV APP_ENV=production

CMD ["./app"]
```

### 5. Dockerfile cho Development/Staging

```dockerfile
# Build với -tags=debug
FROM golang:1.21-alpine
WORKDIR /app
COPY . .

# Build với debug mode
RUN go build -tags=debug -o app .

# Set environment
ENV APP_ENV=development

CMD ["./app"]
```

## 📊 So sánh Performance

| Metric | Production Build | Debug Build |
|--------|-----------------|-------------|
| Binary Size | Nhỏ hơn | Lớn hơn ~5-10% |
| Runtime Overhead | Zero (debug/trace là no-op) | Có overhead khi log |
| Memory Usage | Thấp hơn | Cao hơn khi log nhiều |
| Log Output | Chỉ warn/error/panic | Đầy đủ trace/debug/info/warn/error/panic |

## 🔧 CI/CD Integration

### GitHub Actions Example

```yaml
name: Build and Deploy

on:
  push:
    branches: [main, develop]

jobs:
  build-production:
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      # Build production (không có debug tag)
      - name: Build Production
        run: go build -o app-prod .
      
      - name: Deploy to Production
        run: ./deploy-prod.sh

  build-staging:
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/develop'
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      # Build staging với debug mode
      - name: Build Staging with Debug
        run: go build -tags=debug -o app-staging .
      
      - name: Deploy to Staging
        run: ./deploy-staging.sh
```

## ❓ FAQs

**Q: Tại sao debug logs không hoạt động trong production?**
A: Đây là design có chỉ định. Debug/trace logs chỉ hoạt động khi build với `-tags=debug`. Production build (không có tag) sẽ loại bỏ hoàn toàn code debug/trace để tối ưu performance.

**Q: Làm sao biết app đang chạy ở chế độ nào?**
A: Bạn có thể thêm log khi khởi động:
```go
func main() {
    fmt.Println("Build mode: Production") // Hoặc check build tag
    goerrorkit.InitLogger(...)
}
```

**Q: Có thể enable debug logs trong production không?**
A: Không, nếu bạn build production (không có `-tags=debug`), debug/trace logs sẽ là no-op và không thể enable lại. Bạn cần rebuild với `-tags=debug`.

**Q: Performance impact của debug mode?**
A: Debug mode có overhead khi log messages (I/O, formatting, etc.). Production mode không có overhead vì debug/trace là no-op (compiler optimize away).

**Q: Có cần set environment variable không?**
A: Không bắt buộc. Build tags hoạt động ở compile-time. Nhưng bạn có thể dùng env vars để điều chỉnh LogLevel như trong best practices trên.

