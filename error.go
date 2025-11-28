package goerrorkit

import (
	"fmt"
)

// ErrorType định nghĩa các loại lỗi trong hệ thống
type ErrorType string

const (
	BusinessError   ErrorType = "BUSINESS"   // Lỗi business logic (4xx)
	SystemError     ErrorType = "SYSTEM"     // Lỗi hệ thống (5xx)
	ValidationError ErrorType = "VALIDATION" // Lỗi validation (400)
	AuthError       ErrorType = "AUTH"       // Lỗi authentication/authorization (401-403)
	ExternalError   ErrorType = "EXTERNAL"   // Lỗi từ external service (502-504)
	PanicError      ErrorType = "PANIC"      // Recovered panic
)

// AppError là cấu trúc error chính của thư viện
// Chứa đầy đủ thông tin về lỗi bao gồm type, code, message, stack trace, etc.
type AppError struct {
	Type      ErrorType              // Loại lỗi
	Code      int                    // HTTP status code
	Message   string                 // Message hiển thị
	Details   map[string]interface{} // Thông tin metadata hệ thống (file, line, function, stack trace)
	Data      map[string]interface{} // Dữ liệu đặc thù của tình huống (product_id, user_id, etc.)
	Cause     error                  // Lỗi gốc (nếu có)
	RequestID string                 // Request ID để trace
	logLevel  string                 // Custom log level (warn, error, panic) - private field
}

// Error implements error interface
func (e *AppError) Error() string {
	return e.Message
}

// Unwrap implements errors.Unwrap interface để support errors.Is và errors.As
func (e *AppError) Unwrap() error {
	return e.Cause
}

// WithData thêm dữ liệu đặc thù của tình huống vào error
// Dữ liệu này sẽ được log trong trường "data" riêng biệt
//
// Example:
//
//	return goerrorkit.NewValidationError("Không đủ hàng", nil).WithData(map[string]interface{}{
//	    "product_id": "123",
//	    "product_name": "iPhone 15",
//	    "requested": 1,
//	    "available_stock": 0,
//	})
func (e *AppError) WithData(data map[string]interface{}) *AppError {
	e.Data = data
	return e
}

// WithCallChain thêm full call chain (stack trace) vào error
// Hữu ích khi cần debug chi tiết hoặc trace flow phức tạp
// Lưu ý: Có overhead performance nên chỉ dùng khi cần thiết
//
// Example:
//
//	// Lỗi phức tạp cần trace chi tiết
//	if err := complexOperation(); err != nil {
//	    return goerrorkit.NewSystemError(err).WithCallChain()
//	}
//
//	// Chain với WithData
//	return goerrorkit.NewBusinessError(404, "Product not found").
//	    WithData(map[string]interface{}{"product_id": id}).
//	    WithCallChain()
func (e *AppError) WithCallChain() *AppError {
	callChain := formatStackTraceArray()
	if e.Details == nil {
		e.Details = make(map[string]interface{})
	}
	e.Details["call_chain"] = callChain
	return e
}

// Level thiết lập custom log level cho error
// Hỗ trợ fluent API và cho phép override log level mặc định
// Valid levels: "trace", "debug", "info", "warn", "error", "panic"
//
// Example:
//
//	// ValidationError với log level warn (thay vì error mặc định)
//	return goerrorkit.NewValidationError("Email không hợp lệ", nil).Level("warn")
//
//	// BusinessError nghiêm trọng với log level panic
//	return goerrorkit.NewBusinessError(500, "Data corruption detected").Level("panic")
//
//	// Chain với các methods khác
//	return goerrorkit.NewSystemError(err).
//	    WithData(map[string]interface{}{"db": "postgres"}).
//	    Level("panic").
//	    WithCallChain()
func (e *AppError) Level(level string) *AppError {
	e.logLevel = level
	return e
}

// GetLogLevel trả về log level của error
// Nếu không có custom level, trả về level mặc định dựa trên ErrorType
func (e *AppError) GetLogLevel() string {
	// Nếu có custom level, dùng custom level
	if e.logLevel != "" {
		return e.logLevel
	}

	// Ngược lại, dùng log level mặc định theo error type
	switch e.Type {
	case ValidationError, AuthError:
		return "warn"
	case PanicError, SystemError:
		return "error"
	case BusinessError, ExternalError:
		return "error"
	default:
		return "error"
	}
}

// ============================================================================
// Factory Functions - Tạo Error Dễ Dàng
// ============================================================================

// Wrap đóng gói một Go error thành SystemError với stack trace tự động
// Đây là cách nhanh nhất để wrap error với thông tin chi tiết về vị trí phát sinh
//
// Example:
//
//	if err := db.Query(); err != nil {
//	    return goerrorkit.Wrap(err)
//	}
//
//	// Với custom data
//	if err := json.Unmarshal(data, &result); err != nil {
//	    return goerrorkit.Wrap(err).WithData(map[string]interface{}{
//	        "input_size": len(data),
//	        "expected_type": "User",
//	    })
//	}
func Wrap(err error) *AppError {
	if err == nil {
		return nil
	}
	file, line, function := getCallerInfo(1)
	return &AppError{
		Type:    SystemError,
		Code:    500,
		Message: err.Error(),
		Cause:   err,
		Details: map[string]interface{}{
			"function": function,
			"file":     fmt.Sprintf("%s:%d", file, line),
		},
	}
}

// WrapWithMessage đóng gói Go error với custom message để thêm context
// Message mô tả rõ hơn về ngữ cảnh lỗi, error gốc vẫn được giữ trong Cause
//
// Example:
//
//	if err := redis.Get(key); err != nil {
//	    return goerrorkit.WrapWithMessage(err, "Failed to get user session from cache")
//	}
//
//	// Với custom data
//	if err := os.ReadFile(path); err != nil {
//	    return goerrorkit.WrapWithMessage(err, "Failed to read config file").WithData(map[string]interface{}{
//	        "path": path,
//	        "attempted_at": time.Now(),
//	    })
//	}
func WrapWithMessage(err error, message string) *AppError {
	if err == nil {
		return nil
	}
	file, line, function := getCallerInfo(1)
	return &AppError{
		Type:    SystemError,
		Code:    500,
		Message: message,
		Cause:   err,
		Details: map[string]interface{}{
			"function": function,
			"file":     fmt.Sprintf("%s:%d", file, line),
		},
	}
}

// NewBusinessError tạo lỗi business logic với stack trace chính xác
// Sử dụng .WithData() để thêm dữ liệu đặc thù nếu cần
//
// Example:
//
//	// Không có data
//	if product.Stock == 0 {
//	    return goerrorkit.NewBusinessError(404, "Product out of stock")
//	}
//
//	// Với custom data
//	return goerrorkit.NewBusinessError(404, "Product out of stock").WithData(map[string]interface{}{
//	    "product_id": "123",
//	    "stock": 0,
//	})
func NewBusinessError(code int, msg string) *AppError {
	file, line, function := getCallerInfo(1)
	return &AppError{
		Type:    BusinessError,
		Code:    code,
		Message: msg,
		Details: map[string]interface{}{
			"function": function,
			"file":     fmt.Sprintf("%s:%d", file, line),
		},
	}
}

// NewSystemError tạo lỗi hệ thống với cause và stack trace
// Sử dụng .WithData() để thêm dữ liệu đặc thù nếu cần
//
// Example:
//
//	if err := db.Connect(); err != nil {
//	    return goerrorkit.NewSystemError(err)
//	}
//
//	// Với custom data
//	return goerrorkit.NewSystemError(err).WithData(map[string]interface{}{
//	    "database": "postgres",
//	    "host": "localhost:5432",
//	})
func NewSystemError(err error) *AppError {
	file, line, function := getCallerInfo(1)
	return &AppError{
		Type:    SystemError,
		Code:    500,
		Message: "Internal server error",
		Cause:   err,
		Details: map[string]interface{}{
			"function": function,
			"file":     fmt.Sprintf("%s:%d", file, line),
		},
	}
}

// NewValidationError tạo lỗi validation với custom data
// Data sẽ được log trong trường "data" riêng biệt, tách biệt với metadata hệ thống
//
// Example:
//
//	if age < 18 {
//	    return goerrorkit.NewValidationError("Age must be >= 18", map[string]interface{}{
//	        "field": "age",
//	        "min": 18,
//	        "received": age,
//	    })
//	}
func NewValidationError(msg string, data map[string]interface{}) *AppError {
	file, line, function := getCallerInfo(1)
	return &AppError{
		Type:    ValidationError,
		Code:    400,
		Message: msg,
		Details: map[string]interface{}{
			"function": function,
			"file":     fmt.Sprintf("%s:%d", file, line),
		},
		Data: data,
	}
}

// NewAuthError tạo lỗi authentication/authorization với stack trace
// Sử dụng .WithData() để thêm dữ liệu đặc thù nếu cần
//
// Example:
//
//	if token == "" {
//	    return goerrorkit.NewAuthError(401, "Unauthorized: Missing token")
//	}
//
//	// Với custom data
//	return goerrorkit.NewAuthError(401, "Invalid token").WithData(map[string]interface{}{
//	    "user_id": "123",
//	    "token_expired": true,
//	})
func NewAuthError(code int, msg string) *AppError {
	file, line, function := getCallerInfo(1)
	return &AppError{
		Type:    AuthError,
		Code:    code,
		Message: msg,
		Details: map[string]interface{}{
			"function": function,
			"file":     fmt.Sprintf("%s:%d", file, line),
		},
	}
}

// NewExternalError tạo lỗi từ external service với cause
// Sử dụng .WithData() để thêm dữ liệu đặc thù nếu cần
//
// Example:
//
//	if err := paymentGateway.Charge(); err != nil {
//	    return goerrorkit.NewExternalError(502, "Payment gateway unavailable", err)
//	}
//
//	// Với custom data
//	return goerrorkit.NewExternalError(502, "Payment failed", err).WithData(map[string]interface{}{
//	    "gateway": "stripe",
//	    "amount": 1000,
//	})
func NewExternalError(code int, msg string, cause error) *AppError {
	file, line, function := getCallerInfo(1)
	return &AppError{
		Type:    ExternalError,
		Code:    code,
		Message: msg,
		Cause:   cause,
		Details: map[string]interface{}{
			"function": function,
			"file":     fmt.Sprintf("%s:%d", file, line),
		},
	}
}
