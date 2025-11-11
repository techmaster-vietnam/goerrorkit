package core

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
	Details   map[string]interface{} // Thông tin chi tiết (file, line, function, stack trace)
	Cause     error                  // Lỗi gốc (nếu có)
	RequestID string                 // Request ID để trace
}

// Error implements error interface
func (e *AppError) Error() string {
	return e.Message
}

// Unwrap implements errors.Unwrap interface để support errors.Is và errors.As
func (e *AppError) Unwrap() error {
	return e.Cause
}

// ============================================================================
// Factory Functions - Tạo Error Dễ Dàng
// ============================================================================

// NewBusinessError tạo lỗi business logic với stack trace chính xác
//
// Example:
//
//	if product.Stock == 0 {
//	    return core.NewBusinessError(404, "Product out of stock")
//	}
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
//
// Example:
//
//	if err := db.Connect(); err != nil {
//	    return core.NewSystemError(err)
//	}
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

// NewValidationError tạo lỗi validation với custom details
//
// Example:
//
//	if age < 18 {
//	    return core.NewValidationError("Age must be >= 18", map[string]interface{}{
//	        "field": "age",
//	        "min": 18,
//	        "received": age,
//	    })
//	}
func NewValidationError(msg string, details map[string]interface{}) *AppError {
	file, line, function := getCallerInfo(1)
	if details == nil {
		details = make(map[string]interface{})
	}
	details["function"] = function
	details["file"] = fmt.Sprintf("%s:%d", file, line)
	return &AppError{
		Type:    ValidationError,
		Code:    400,
		Message: msg,
		Details: details,
	}
}

// NewAuthError tạo lỗi authentication/authorization với stack trace
//
// Example:
//
//	if token == "" {
//	    return core.NewAuthError(401, "Unauthorized: Missing token")
//	}
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
//
// Example:
//
//	if err := paymentGateway.Charge(); err != nil {
//	    return core.NewExternalError(502, "Payment gateway unavailable", err)
//	}
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
