# Demo Build Modes: Development vs Production

File này hướng dẫn cách test sự khác biệt giữa **debug build** và **production build** của GoErrorKit.

## 🎯 Mục đích

Trong production, bạn không muốn debug/trace logs xuất hiện vì:
- 🚀 Giảm overhead về performance
- 📦 Giảm kích thước binary
- 🔒 Tránh expose thông tin nhạy cảm
- 💰 Tiết kiệm storage cho log files

GoErrorKit giải quyết vấn đề này bằng **build tags** - compile-time flags để bật/tắt debug/trace logging.

## 🧪 Test 1: Simple Logging Test

### Step 1: Tạo file test đơn giản

Tạo file `test.go`:

```go
package main

import (
	"github.com/techmaster-vietnam/goerrorkit"
)

func main() {
	// Init logger với trace level
	goerrorkit.InitLogger(goerrorkit.LoggerOptions{
		ConsoleOutput: true,
		FileOutput:    false,
		JSONFormat:    false,
		LogLevel:      "trace", // Thấp nhất - log mọi thứ
	})

	logger := goerrorkit.GetLogger()

	// Test tất cả log levels
	logger.Trace("🔍 TRACE: Chi tiết debug sâu nhất", nil)
	logger.Debug("🐛 DEBUG: Debug information", nil)
	logger.Info("ℹ️ INFO: General information", nil)
	logger.Warn("⚠️ WARN: Warning message", nil)
	logger.Error("❌ ERROR: Error occurred", nil)
}
```

### Step 2: Test Production Build (Mặc định)

```bash
# Run với production build (KHÔNG có tag debug)
go run test.go
```

**Expected Output:**
```
ℹ️ INFO: General information
⚠️ WARN: Warning message
❌ ERROR: Error occurred
```

⚡ **Chú ý:** TRACE và DEBUG **KHÔNG hiển thị** (zero overhead)!

### Step 3: Test Debug Build

```bash
# Run với debug build (có tag debug)
go run -tags=debug test.go
```

**Expected Output:**
```
🔍 TRACE: Chi tiết debug sâu nhất
🐛 DEBUG: Debug information
ℹ️ INFO: General information
⚠️ WARN: Warning message
❌ ERROR: Error occurred
```

✅ **Tất cả logs hiển thị**, bao gồm TRACE và DEBUG!

## 🧪 Test 2: Error Logging với Custom Level

### File test

```go
package main

import (
	"github.com/techmaster-vietnam/goerrorkit"
	fiberv2 "github.com/gofiber/fiber/v2"
)

func main() {
	goerrorkit.InitLogger(goerrorkit.LoggerOptions{
		ConsoleOutput: true,
		LogLevel:      "trace",
	})

	// Tạo một error với debug level
	err := goerrorkit.NewBusinessError(404, "Product not found").
		Level("debug"). // Custom log level = debug
		WithData(map[string]interface{}{
			"product_id": "123",
		})

	// Log error này
	goerrorkit.LogError(err, "/api/products/123")
}
```

**Production build:** Error **KHÔNG log** (vì level=debug)
**Debug build:** Error **được log** đầy đủ

## 🏗️ Test 3: Build Binary

### Production Binary

```bash
# Build production binary
go build -o app-prod

# Binary size nhỏ hơn, không có debug code
./app-prod
```

### Debug Binary

```bash
# Build debug binary
go build -tags=debug -o app-debug

# Binary size lớn hơn, có đầy đủ debug code
./app-debug
```

### So sánh kích thước

```bash
ls -lh app-*
```

Bạn sẽ thấy `app-debug` lớn hơn `app-prod` một chút (~5-10%).

## 📊 Test 4: Performance Benchmark

Tạo file `benchmark_test.go`:

```go
package main

import (
	"testing"
	"github.com/techmaster-vietnam/goerrorkit"
)

func init() {
	goerrorkit.InitLogger(goerrorkit.LoggerOptions{
		ConsoleOutput: false, // Disable console để chỉ test overhead
		FileOutput:    false,
		LogLevel:      "trace",
	})
}

func BenchmarkDebugLog(b *testing.B) {
	logger := goerrorkit.GetLogger()
	fields := map[string]interface{}{
		"user_id": 123,
		"action":  "test",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Debug("Debug message", fields)
	}
}

func BenchmarkInfoLog(b *testing.B) {
	logger := goerrorkit.GetLogger()
	fields := map[string]interface{}{
		"user_id": 123,
		"action":  "test",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("Info message", fields)
	}
}
```

### Chạy benchmark

```bash
# Production build
go test -bench=. -benchmem

# Debug build
go test -tags=debug -bench=. -benchmem
```

**Expected Results:**
- Production: `BenchmarkDebugLog` cực kỳ nhanh (no-op)
- Debug: `BenchmarkDebugLog` chậm hơn (thực sự log)

## 🐳 Test 5: Docker Build

### Production Dockerfile

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .

# Build KHÔNG có -tags=debug
RUN go build -o app .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/app .
CMD ["./app"]
```

### Debug Dockerfile

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .

# Build VỚI -tags=debug
RUN go build -tags=debug -o app .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/app .
CMD ["./app"]
```

### Build và test

```bash
# Production
docker build -f Dockerfile.prod -t myapp:prod .
docker run myapp:prod

# Debug
docker build -f Dockerfile.debug -t myapp:debug .
docker run myapp:debug
```

## ✅ Checklist Verification

Sau khi test, verify rằng:

- [ ] Production build: Debug/trace logs **KHÔNG hiển thị**
- [ ] Debug build: Debug/trace logs **hiển thị đầy đủ**
- [ ] Production build: Info/warn/error logs vẫn hoạt động bình thường
- [ ] Binary size: Debug build lớn hơn production build
- [ ] Performance: Debug logs trong production build có overhead = 0

## 💡 Best Practices Đã Verify

1. ✅ **Development**: Build với `-tags=debug`, set `LogLevel: "debug"`
2. ✅ **Staging**: Build với `-tags=debug`, set `LogLevel: "info"`
3. ✅ **Production**: Build KHÔNG có tag, set `LogLevel: "warn"` hoặc `"error"`

## 🔧 Troubleshooting

### Debug logs không hiển thị trong debug build?

Check:
```bash
# Đảm bảo bạn build với tag debug
go run -tags=debug main.go

# Hoặc
go build -tags=debug -o app
```

### Debug logs vẫn hiển thị trong production?

Check:
```bash
# Đảm bảo bạn build KHÔNG có tag debug
go build -o app  # ĐÚNG
go build -tags=debug -o app  # SAI (đây là debug build)
```

### Muốn verify build mode?

Thêm vào code:

```go
//go:build debug
// +build debug

package main

import "fmt"

func init() {
	fmt.Println("🐛 DEBUG BUILD MODE")
}
```

Tạo file tương tự cho production:

```go
//go:build !debug
// +build !debug

package main

import "fmt"

func init() {
	fmt.Println("🚀 PRODUCTION BUILD MODE")
}
```

## 📚 Tài liệu thêm

- [Build Modes Documentation](../docs/build-modes.md) - Chi tiết đầy đủ
- [Go Build Tags](https://pkg.go.dev/cmd/go#hdr-Build_constraints) - Official docs

