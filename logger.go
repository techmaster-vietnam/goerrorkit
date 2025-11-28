package goerrorkit

// Logger interface cho phép user tùy chỉnh logging implementation
// Default implementation sẽ dùng logrus, nhưng user có thể dùng zap, zerolog, etc.
type Logger interface {
	// Error logs error level message với fields
	Error(msg string, fields map[string]interface{})

	// Info logs info level message với fields
	Info(msg string, fields map[string]interface{})

	// Debug logs debug level message với fields
	// Lưu ý: Chỉ hoạt động khi build với tag -tags=debug
	// Production build sẽ bỏ qua hoàn toàn (zero overhead)
	Debug(msg string, fields map[string]interface{})

	// Trace logs trace level message với fields
	// Lưu ý: Chỉ hoạt động khi build với tag -tags=debug
	// Production build sẽ bỏ qua hoàn toàn (zero overhead)
	Trace(msg string, fields map[string]interface{})

	// Warn logs warning level message với fields
	Warn(msg string, fields map[string]interface{})

	// Panic logs panic level message với fields
	Panic(msg string, fields map[string]interface{})
}

// defaultLogger là logger mặc định (sẽ được set từ config package)
var defaultLogger Logger

// SetLogger cho phép user set custom logger implementation
//
// Example:
//
//	logger := myCustomLogger{} // implements goerrorkit.Logger interface
//	goerrorkit.SetLogger(logger)
func SetLogger(l Logger) {
	defaultLogger = l
}

// GetLogger trả về logger hiện tại
func GetLogger() Logger {
	return defaultLogger
}

// LogError xử lý logging cho AppError
// Sử dụng appropriate log level dựa trên error.GetLogLevel()
func LogError(appErr *AppError, requestPath string) {
	if defaultLogger == nil {
		// Nếu chưa set logger, skip logging
		return
	}

	// Chuẩn bị log fields với metadata cơ bản
	fields := map[string]interface{}{
		"error_type": string(appErr.Type),
		"path":       requestPath,
	}

	// Thêm metadata hệ thống từ Details (function, file, stack trace)
	for k, v := range appErr.Details {
		fields[k] = v
	}

	// Thêm dữ liệu đặc thù vào trường "data" riêng biệt (nếu có)
	if len(appErr.Data) > 0 {
		fields["data"] = appErr.Data
	}

	// Thêm cause nếu có
	if appErr.Cause != nil {
		fields["cause"] = appErr.Cause.Error()
	}

	// Log với level phù hợp (trace, debug, info, warn, error, panic)
	logLevel := appErr.GetLogLevel()
	switch logLevel {
	case "panic":
		defaultLogger.Panic(appErr.Message, fields)
	case "error":
		defaultLogger.Error(appErr.Message, fields)
	case "warn":
		defaultLogger.Warn(appErr.Message, fields)
	case "info":
		defaultLogger.Info(appErr.Message, fields)
	case "debug":
		defaultLogger.Debug(appErr.Message, fields)
	case "trace":
		defaultLogger.Trace(appErr.Message, fields)
	default:
		// Default fallback to error
		defaultLogger.Error(appErr.Message, fields)
	}
}

// FormatErrorResponse tạo response data cho client
// Chỉ trả về thông tin cần thiết, không expose internal details
func FormatErrorResponse(appErr *AppError) map[string]interface{} {
	return map[string]interface{}{
		"error": appErr.Message,
		"type":  string(appErr.Type),
	}
}

// LogAndRespond xử lý logging và gửi response (framework agnostic)
// Đây là helper function cho adapters
func LogAndRespond(ctx HTTPContext, appErr *AppError, requestPath string) {
	// 1. Log error
	LogError(appErr, requestPath)

	// 2. Send response
	ctx.Status(appErr.Code).JSON(FormatErrorResponse(appErr))
}
