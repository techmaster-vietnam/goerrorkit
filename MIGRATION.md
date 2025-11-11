# Migration Guide

## Breaking Changes - v0.2.0

### 🎉 Không có Breaking Changes!

**Chữ ký hàm giữ nguyên**, không cần sửa code hiện tại. Chỉ có thêm tính năng mới.

### ✨ Tính Năng Mới: Trường `data` Riêng Biệt

#### Không thay đổi chữ ký hàm

Tất cả các hàm giữ nguyên như cũ:

```go
// Vẫn dùng được như cũ - KHÔNG CẦN SỬA
goerrorkit.NewBusinessError(404, "Product not found")
goerrorkit.NewSystemError(err)
goerrorkit.NewAuthError(401, "Unauthorized")
goerrorkit.NewExternalError(502, "Gateway error", err)
```

#### Thêm tính năng: .WithData()

Giờ bạn có thể thêm dữ liệu đặc thù với fluent API `.WithData()`:

```go
// Không có data - clean và simple (majority case)
return goerrorkit.NewBusinessError(404, "Product not found")

// Có data - dùng .WithData() khi cần (minority case)
return goerrorkit.NewBusinessError(404, "Product not found").WithData(map[string]interface{}{
    "product_id": productID,
})
```

### Lợi Ích

1. **Clean code**: Không cần viết `, nil` cho ~80% trường hợp không cần data
2. **Self-documenting**: `.WithData()` rõ ràng về mục đích  
3. **Go idioms**: Giống pattern của stdlib (`context.WithTimeout()`, `grpc.WithInsecure()`)
4. **Backward compatible**: Code cũ chạy ngon không cần sửa

### 📊 Log Output Format - Dễ Đọc Hơn

**Trước (v0.1.x):**
```json
{
  "level": "error",
  "message": "Không đủ hàng",
  "error_type": "VALIDATION",
  "function": "services.ReserveProduct",
  "file": "product_service.go:70",
  "product_id": "123",
  "product_name": "iPhone 15",
  "requested": 1,
  "available_stock": 0
}
```

**Sau (v0.2.0):**
```json
{
  "level": "error",
  "message": "Không đủ hàng",
  "error_type": "VALIDATION",
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

**Dễ đọc hơn rất nhiều!** Metadata hệ thống tách biệt với dữ liệu tình huống.

### 🗑️ Request ID Đã Bị Loại Bỏ

Request ID đã được loại bỏ khỏi response trả về client (vẫn có trong log).

**Trước:**
```json
{
  "error": "Not found",
  "type": "BUSINESS",
  "request_id": "abc-123"
}
```

**Sau:**
```json
{
  "error": "Not found",
  "type": "BUSINESS"
}
```

**Lý do:** Request ID là thông tin internal, không nên expose ra client.

### 💡 Ví Dụ Sử Dụng

#### Trường hợp đơn giản (không cần data)

```go
// Clean và concise!
return goerrorkit.NewBusinessError(404, "Product not found")
return goerrorkit.NewSystemError(err)
return goerrorkit.NewAuthError(401, "Unauthorized")
```

#### Trường hợp cần data đặc thù

```go
// Dùng .WithData() khi cần
return goerrorkit.NewBusinessError(404, "Product not found").WithData(map[string]interface{}{
    "product_id": productID,
    "category": "electronics",
})

return goerrorkit.NewSystemError(dbErr).WithData(map[string]interface{}{
    "database": "postgres",
    "host": "localhost:5432",
})

return goerrorkit.NewAuthError(403, "Insufficient permissions").WithData(map[string]interface{}{
    "user_id": userID,
    "required_role": "admin",
    "user_role": currentRole,
})
```

#### Validation Error (có parameter data)

```go
// Validation thường cần data → truyền trực tiếp
return goerrorkit.NewValidationError("Age must be >= 18", map[string]interface{}{
    "field": "age",
    "min": 18,
    "received": age,
})
```

### 🎯 Best Practices

1. **Majority case (80%)**: Không cần data → code ngắn gọn
   ```go
   return goerrorkit.NewBusinessError(404, "Not found")
   ```

2. **Minority case (20%)**: Cần data → dùng `.WithData()`
   ```go
   return goerrorkit.NewBusinessError(404, "Not found").WithData(data)
   ```

3. **Validation**: Hầu hết cần data → parameter
   ```go
   return goerrorkit.NewValidationError("Invalid", data)
   ```

