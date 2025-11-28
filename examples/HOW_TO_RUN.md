# 🚀 Hướng Dẫn Chạy Demo

## ⚠️  QUAN TRỌNG: Debug vs Production Mode

GoErrorKit có 2 chế độ build khác nhau:

### 1️⃣ **Development Mode** (với Trace & Debug logging)

```bash
cd examples
go run -tags=debug main.go
```

**Kích hoạt:**
- ✅ Trace logging
- ✅ Debug logging
- ✅ Info, Warn, Error logging

**Use case:**
- Development và debugging
- Chi tiết flow execution
- Performance profiling
- Troubleshooting

### 2️⃣ **Production Mode** (không có Trace & Debug)

```bash
cd examples
go run main.go
```

**Kích hoạt:**
- ❌ Trace logging (no-op)
- ❌ Debug logging (no-op)
- ✅ Info, Warn, Error logging

**Lợi ích:**
- 🚀 Zero overhead cho debug/trace
- 📦 Binary nhỏ hơn
- 🔒 Không expose debug info
- 💰 Tiết kiệm storage

---

## 📋 Endpoints để Test

### Trace & Debug Endpoints
**⚠️  Chỉ hoạt động với `-tags=debug`**

```bash
# Trace single operation
curl http://localhost:8081/dev/trace?op=fetch_user

# Debug with detailed context
curl http://localhost:8081/dev/debug?scenario=user_login

# Trace complex multi-step flow
curl http://localhost:8081/dev/trace-complex?order_id=ORD-12345
```

### Error Demo Endpoints
**Hoạt động ở cả 2 modes**

```bash
# Panic demos
curl http://localhost:8081/panic/division
curl http://localhost:8081/panic/index

# Error demos
curl http://localhost:8081/error/business
curl http://localhost:8081/error/system
curl http://localhost:8081/error/validation

# Log level demo
curl "http://localhost:8081/error/log-level?level=warn"
curl "http://localhost:8081/error/log-level?level=error"
```

---

## 🧪 Test Build Modes

### Test 1: Kiểm tra Trace không hoạt động ở Production

**Production Mode:**
```bash
go run main.go
# Truy cập: http://localhost:8081/dev/trace
# Kết quả: Không có trace log nào trong console ❌
```

**Development Mode:**
```bash
go run -tags=debug main.go
# Truy cập: http://localhost:8081/dev/trace
# Kết quả: Thấy trace logs trong console ✅
```

### Test 2: Kiểm tra Debug không hoạt động ở Production

**Production Mode:**
```bash
go run main.go
# Truy cập: http://localhost:8081/dev/debug
# Kết quả: Không có debug log nào trong console ❌
```

**Development Mode:**
```bash
go run -tags=debug main.go
# Truy cập: http://localhost:8081/dev/debug
# Kết quả: Thấy debug logs với chi tiết trong console ✅
```

---

## 💡 FAQ

### Q: Tại sao cần 2 modes?

**A:** Trong production:
- Trace/debug logs thường có quá nhiều thông tin không cần thiết
- Tốn performance và storage
- Có thể expose thông tin nhạy cảm
- GoErrorKit dùng build tags để **compile-time disable** debug/trace → zero overhead

### Q: LogLevel="trace" nhưng không log ra?

**A:** Bạn đang chạy ở production mode!
- ❌ `go run main.go` → trace/debug = no-op
- ✅ `go run -tags=debug main.go` → trace/debug hoạt động

### Q: Có nên gán `goerrorkit.GetLogger()` vào biến?

**A:** Không nên!
```go
// ❌ SAI - Không cache logger
logger := goerrorkit.GetLogger()
logger.Trace("...", nil)
logger.Debug("...", nil)

// ✅ ĐÚNG - Gọi GetLogger() mỗi lần
goerrorkit.GetLogger().Trace("...", nil)
goerrorkit.GetLogger().Debug("...", nil)
```

**Lý do:**
- `GetLogger()` chỉ trả về biến global, không tốn performance
- Cache có thể miss updates nếu logger thay đổi
- Best practice: luôn gọi `GetLogger()` mỗi lần dùng

---

## 📖 Chi Tiết Thêm

- **Build modes chi tiết:** Xem [BUILD_MODES_DEMO.md](./BUILD_MODES_DEMO.md)
- **Architecture:** Xem [../BUILD_MODES_SUMMARY.md](../BUILD_MODES_SUMMARY.md)
- **Getting started:** Xem [../docs/getting-started.md](../docs/getting-started.md)

