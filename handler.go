package goerrorkit

import (
	"fmt"
)

// HandlePanic xử lý panic và trả về AppError với stack trace chi tiết
// Đây là core function để capture panic location chính xác
//
// Example (internal use):
//
//	defer func() {
//	    if r := recover(); r != nil {
//	        panicErr := goerrorkit.HandlePanic(r, requestID)
//	        // Log and respond...
//	    }
//	}()
func HandlePanic(r interface{}, requestID string) *AppError {
	actualFile, actualLine, actualFunc := getActualPanicLocation()
	callChain := formatStackTraceArray()

	return &AppError{
		Type:      PanicError,
		Code:      500,
		Message:   fmt.Sprintf("Panic recovered: %v", r),
		RequestID: requestID,
		Details: map[string]interface{}{
			"panic_value": r,
			"function":    actualFunc,
			"file":        fmt.Sprintf("%s:%d", actualFile, actualLine),
			"call_chain":  callChain,
		},
	}
}

// ConvertToAppError chuyển đổi error thường thành AppError
// Nếu đã là AppError thì chỉ update RequestID
//
// Example (internal use):
//
//	err := someFunction()
//	if err != nil {
//	    appErr := goerrorkit.ConvertToAppError(err, requestID)
//	    return appErr
//	}
func ConvertToAppError(err error, requestID string) *AppError {
	// Check nếu đã là AppError
	if appErr, ok := err.(*AppError); ok {
		appErr.RequestID = requestID
		return appErr
	}

	// Convert error thường thành AppError
	return &AppError{
		Type:      SystemError,
		Code:      500,
		Message:   "Internal server error",
		Cause:     err,
		RequestID: requestID,
	}
}
