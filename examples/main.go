package main

import (
	"fmt"

	fiberv2 "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/techmaster-vietnam/goerrorkit"
)

func main() {
	// 1. Initialize logger với dual-level logging
	// 🎯 MỤC ĐÍCH: Phân tách log levels cho console và file
	//    - Console: Log tất cả errors từ warn trở lên (để developer debug)
	//    - File: Chỉ log errors nghiêm trọng (error, panic) để dễ phân tích production issues
	//
	// ✅ KẾT QUẢ:
	//    - ValidationError (level: warn) → Console: ✓, File: ✗
	//    - AuthError (level: warn)       → Console: ✓, File: ✗
	//    - SystemError (level: error)    → Console: ✓, File: ✓
	//    - PanicError (level: error)     → Console: ✓, File: ✓
	//
	// ⚠️  QUAN TRỌNG - TRACE & DEBUG LOGGING:
	//    - Trace/Debug CHỈ hoạt động khi build với tag: -tags=debug
	//    - Development: go run -tags=debug main.go
	//    - Production:  go run main.go (trace/debug = no-op)
	//
	// 📌 LogLevel="trace" trong config có nghĩa là:
	//    - Với -tags=debug:  Log TRACE, DEBUG, INFO, WARN, ERROR
	//    - Không tag debug:  Log INFO, WARN, ERROR (trace/debug bị tắt)
	goerrorkit.InitLogger(goerrorkit.LoggerOptions{
		ConsoleOutput: true,
		FileOutput:    true,
		FilePath:      "logs/errors.log",
		JSONFormat:    true,
		MaxFileSize:   10,
		MaxBackups:    5,
		MaxAge:        30,
		LogLevel:      "trace", // Console log từ trace (cần -tags=debug), warn, error trở lên
		FileLogLevel:  "error", // File chỉ log error và panic (bỏ qua warn)
	})

	// 2. Configure stack trace for this application
	// 🎯 MỤC ĐÍCH: Lọc stack trace để CHỈ HIỂN THỊ code của BẠN, bỏ qua:
	//    - Go runtime code (runtime.*, runtime/debug.*)
	//    - Thư viện bên thứ 3 (fiber, goerrorkit, etc.)
	//
	// ✅ CÁCH DÙNG:
	//    - App đơn giản (1 file main.go):
	//      goerrorkit.ConfigureForApplication("main")
	//
	//    - App với nhiều package (services/, handlers/, models/...):
	//      goerrorkit.ConfigureForApplication("github.com/techmaster-vietnam/goerrorkit/examples/fiber-demo")
	//      → Tự động include TẤT CẢ sub-packages!
	//
	// 📊 KẾT QUẢ:
	//    KHÔNG cấu hình: Stack trace dài 50+ dòng (runtime, fiber, goerrorkit...)
	//    CÓ cấu hình:    Stack trace ngắn gọn, chỉ 5-10 dòng CODE CỦA BẠN!
	//
	goerrorkit.ConfigureForApplication("main")

	// 🔧 FLUENT API: Nếu cần thêm các patterns tùy chỉnh, có thể dùng:
	//
	// Cách 1: Shorthand - Nhanh chóng thêm skip patterns
	// goerrorkit.AddSkipPatterns(".RequestID.func", ".Logger.func", "telemetry")
	//
	// Cách 2: Fluent API - Configuration chi tiết hơn
	// goerrorkit.Configure().
	//     SkipPattern(".CustomMiddleware.func").
	//     SkipPackage("internal/metrics").
	//     SkipFunctions("helper", "wrapper").
	//     ShowFullPath(false).
	//     Apply()

	// 3. Create Fiber app
	app := fiberv2.New(fiberv2.Config{
		AppName: "GoErrorKit Demo",
	})

	// 🗄️ Run database migrations
	// Giả lập lỗi migration có cause (lỗi gốc)
	// Set simulateError = true để demo error logging
	if err := runDatabaseMigrations(true); err != nil {
		// ⭐ Log error ra console và file sử dụng GoErrorKit
		// LogError tự động log ra cả console và file (đã config ở dòng 14)
		goerrorkit.LogError(err.(*goerrorkit.AppError), "/startup/migrations")

		// In thông báo warning nhưng vẫn tiếp tục chạy server
		fmt.Println("⚠️  Migration failed but server will continue...")
	}

	// 4. Add middlewares (RequestID must be before ErrorHandler)
	app.Use(requestid.New())
	app.Use(logger.New())
	app.Use(goerrorkit.FiberErrorHandler())

	// 5. Routes - Demo different error types
	app.Get("/", homeHandler)
	app.Get("/favicon.ico", faviconHandler) // Serve favicon
	app.Get("/favicon.svg", faviconHandler) // Modern SVG favicon

	// Panic demos
	app.Get("/panic/division", panicDivisionHandler)
	app.Get("/panic/index", panicIndexHandler)
	app.Get("/panic/stack", panicStackHandler)

	// Wrap error demos (NEW!)
	app.Get("/error/wrap", wrapErrorHandler)
	app.Get("/error/wrap-message", wrapWithMessageHandler)

	// Custom error demos
	app.Get("/error/business", businessErrorHandler)
	app.Get("/error/system", systemErrorHandler)
	app.Get("/error/validation", validationErrorHandler)
	app.Get("/error/auth", authErrorHandler)
	app.Get("/error/external", externalErrorHandler)
	app.Get("/error/complex", complexErrorWithCallChainHandler)

	// Log level override demo (NEW!)
	app.Get("/error/log-level", logLevelDemoHandler)

	// Development tools - Trace & Debug demos (NEW!)
	// ⚠️  QUAN TRỌNG: Chỉ hoạt động khi build với: go run -tags=debug main.go
	app.Get("/dev/trace", traceHandler)
	app.Get("/dev/debug", debugHandler)
	app.Get("/dev/trace-complex", traceComplexFlowHandler)

	// Start server
	fmt.Println("🚀 Server starting on http://localhost:8081")
	fmt.Println("\n📝 Try these endpoints:")
	fmt.Println("  GET  /                     - Home")
	fmt.Println("\n  🔥 Panic Demos (auto-recovered with exact location):")
	fmt.Println("  GET  /panic/division       - Division by zero panic")
	fmt.Println("  GET  /panic/index          - Index out of range panic")
	fmt.Println("  GET  /panic/stack          - Deep call stack panic")
	fmt.Println("\n  🎁 Wrap Error Demos:")
	fmt.Println("  GET  /error/wrap           - Wrap(err) - Đơn giản nhất")
	fmt.Println("  GET  /error/wrap-message   - WrapWithMessage(err, msg) - Thêm context")
	fmt.Println("\n  ⚠️  Custom Error Demos:")
	fmt.Println("  GET  /error/business       - Business logic error (404)")
	fmt.Println("  GET  /error/system         - System error (500)")
	fmt.Println("  GET  /error/validation     - Validation error (400)")
	fmt.Println("  GET  /error/auth           - Auth error (401)")
	fmt.Println("  GET  /error/external       - External service error (502)")
	fmt.Println("  GET  /error/complex        - Complex error WITH call_chain")
	fmt.Println("\n  🎯 Log Level Demo (NEW!):")
	fmt.Println("  GET  /error/log-level?level=warn  - Demo log level override ⭐")
	fmt.Println("       ?level=warn   → Console: ✓, File: ✗")
	fmt.Println("       ?level=error  → Console: ✓, File: ✓")
	fmt.Println("\n  🔧 Development Tools - Trace & Debug:")
	fmt.Println("  ⚠️  CHÚ Ý: Cần build với -tags=debug để trace/debug hoạt động!")
	fmt.Println("     ✅ go run -tags=debug main.go    (development)")
	fmt.Println("     ❌ go run main.go                (trace/debug không hoạt động)")
	fmt.Println("  GET  /dev/trace          - Trace single operation")
	fmt.Println("  GET  /dev/debug          - Debug with detailed context")
	fmt.Println("  GET  /dev/trace-complex  - Trace complex multi-step flow")
	fmt.Println("\n📄 Logs:")
	fmt.Println("  - Console: Shows ALL errors (warn, error, panic)")
	fmt.Println("  - File (logs/errors.log): Only SERIOUS errors (error, panic)")
	fmt.Println("  💡 Try validation/auth errors → see console log, but NOT in file!")

	if err := app.Listen(":8081"); err != nil {
		panic(err)
	}
}

func homeHandler(c *fiberv2.Ctx) error {
	// Serve the index.html file
	return c.SendFile("./index.html")
}

func faviconHandler(c *fiberv2.Ctx) error {
	// Serve favicon.svg (modern browsers support SVG favicons)
	return c.SendFile("./favicon.svg")
}

// ============================================================================
// Panic Handlers - Demonstrate automatic panic recovery
// ============================================================================

func panicDivisionHandler(c *fiberv2.Ctx) error {
	// This will panic with "integer divide by zero"
	denominator := 0
	result := 100 / denominator // ← Panic location will be captured HERE!
	return c.JSON(fiberv2.Map{"result": result})
}

func panicIndexHandler(c *fiberv2.Ctx) error {
	// This will panic with "index out of range"
	element := GetElement() // Panic happens inside GetElement()
	return c.JSON(fiberv2.Map{"element": element})
}

func GetElement() int {
	arr := []int{1, 2, 3}
	return arr[10] // ← Panic location will be captured HERE!
}

func panicStackHandler(c *fiberv2.Ctx) error {
	// Deep call stack demo
	result := callX()
	return c.JSON(fiberv2.Map{"result": result})
}

func callX() int {
	return callY()
}

func callY() int {
	return callZ()
}

func callZ() int {
	return callW()
}

func callW() int {
	return GetElement() // Panic happens here, full call chain will be logged
}

// ============================================================================
// Wrap Error Handlers - Demonstrate Wrap() and WrapWithMessage()
// ============================================================================

// wrapErrorHandler demonstrates goerrorkit.Wrap(err)
// ✅ Use case: Wrap nhanh Go error với message gốc, tự động capture stack trace
func wrapErrorHandler(c *fiberv2.Ctx) error {
	errorType := c.Query("type", "json")

	switch errorType {
	case "json":
		// Simulate JSON parsing error
		err := fmt.Errorf("json: invalid character '}' looking for beginning of value")
		// ⭐ Wrap() - Đơn giản nhất, giữ nguyên message gốc
		// → Message: error message gốc
		// → Type: SystemError, Code: 500
		// → Tự động capture: file, line, function
		return goerrorkit.Wrap(err)

	case "database":
		// Simulate database connection error
		err := fmt.Errorf("sql: connection refused")
		// ⭐ Wrap() với error database
		return goerrorkit.Wrap(err)

	case "file":
		// Simulate file not found error
		err := fmt.Errorf("open config.json: no such file or directory")
		// ⭐ Wrap() với .WithData() - Thêm metadata
		return goerrorkit.Wrap(err).WithData(map[string]interface{}{
			"path":      "config.json",
			"operation": "read",
		})

	case "network":
		// Simulate network timeout
		err := fmt.Errorf("net/http: request timeout after 30s")
		// ⭐ Wrap() + WithData() + WithCallChain()
		return goerrorkit.Wrap(err).
			WithData(map[string]interface{}{
				"url":     "https://api.example.com/users",
				"timeout": "30s",
				"retries": 3,
			}).
			WithCallChain()
	}

	return c.JSON(fiberv2.Map{"message": "No error"})
}

// wrapWithMessageHandler demonstrates goerrorkit.WrapWithMessage(err, msg)
// ✅ Use case: Wrap error với custom message để thêm context, giữ error gốc trong Cause
func wrapWithMessageHandler(c *fiberv2.Ctx) error {
	scenario := c.Query("scenario", "database")

	switch scenario {
	case "database":
		// Simulate database query error
		err := fmt.Errorf("connection refused")
		// ⭐ WrapWithMessage() - Thêm context message
		// → Message: "Failed to fetch user list from database"
		// → Cause: "connection refused"
		// → Type: SystemError, Code: 500
		return goerrorkit.WrapWithMessage(err, "Failed to fetch user list from database")

	case "redis":
		// Simulate Redis cache error
		err := fmt.Errorf("redis: connection timeout")
		// ⭐ WrapWithMessage() với .WithData()
		return goerrorkit.WrapWithMessage(err, "Failed to get user session from cache").WithData(map[string]interface{}{
			"key": "user:session:12345",
			"ttl": 3600,
		})

	case "payment":
		// Simulate payment API error
		err := fmt.Errorf("stripe: card declined")
		// ⭐ WrapWithMessage() + WithData() - Detailed context
		return goerrorkit.WrapWithMessage(err, "Payment processing failed").WithData(map[string]interface{}{
			"gateway":    "stripe",
			"amount":     10000,
			"currency":   "VND",
			"payment_id": "pay_123456",
		})

	case "email":
		// Simulate email service error
		err := fmt.Errorf("smtp: authentication failed")
		// ⭐ WrapWithMessage() + WithData() + WithCallChain()
		return goerrorkit.WrapWithMessage(err, "Failed to send verification email").
			WithData(map[string]interface{}{
				"to":       "user@example.com",
				"template": "email_verification",
				"smtp":     "smtp.gmail.com:587",
			}).
			WithCallChain()
	}

	return c.JSON(fiberv2.Map{"message": "No error"})
}

// ============================================================================
// Custom Error Handlers - Demonstrate different error types
// ============================================================================

func businessErrorHandler(c *fiberv2.Ctx) error {
	productID := c.Query("product_id", "unknown")

	// Simulate product not found (normal business case)
	if productID == "123" {
		// ⭐ BusinessError với default log level "error"
		// → Console: ✓, File: ✓
		return goerrorkit.NewBusinessError(404, fmt.Sprintf("Product ID=%s not found", productID)).WithData(map[string]interface{}{
			"product_id": productID,
		})
	}

	// 🎯 DEMO: BusinessError nghiêm trọng với .Level("error")
	// Nếu stock < 0, đây là lỗi nghiêm trọng cần investigate
	if productID == "corrupted" {
		return goerrorkit.NewBusinessError(500, "Data corruption: Negative stock detected").
			WithData(map[string]interface{}{
				"product_id": productID,
				"stock":      -10,
				"warehouse":  "WH-01",
			}).
			Level("error") // ⭐ Đảm bảo ghi vào file (đã là error rồi, nhưng làm rõ intent)
	}

	return c.JSON(fiberv2.Map{
		"message":    "Product available",
		"product_id": productID,
	})
}

func systemErrorHandler(c *fiberv2.Ctx) error {
	// Simulate database connection error
	err := fmt.Errorf("connection refused: database is down")
	// ⭐ SystemError với default log level "error"
	// → Console: ✓ (log), File: ✓ (log vào file vì error >= FileLogLevel)
	return goerrorkit.NewSystemError(err).WithData(map[string]interface{}{
		"database": "postgres",
		"host":     "localhost:5432",
	})
}

func validationErrorHandler(c *fiberv2.Ctx) error {
	age := c.Query("age", "")

	if age == "" {
		// ⭐ ValidationError với default log level "warn"
		// → Console: ✓ (log), File: ✗ (bỏ qua vì FileLogLevel = "error")
		return goerrorkit.NewValidationError("Missing parameter 'age'", map[string]interface{}{
			"field":    "age",
			"required": true,
		})
	}

	// Check if age is a number
	var ageInt int
	if _, err := fmt.Sscanf(age, "%d", &ageInt); err != nil {
		// ⭐ ValidationError với default log level "warn"
		return goerrorkit.NewValidationError("Parameter 'age' must be an integer", map[string]interface{}{
			"field":    "age",
			"type":     "integer",
			"received": age,
		})
	}

	if ageInt < 18 {
		// ⭐ ValidationError với default log level "warn"
		return goerrorkit.NewValidationError("Age must be >= 18", map[string]interface{}{
			"field":    "age",
			"min":      18,
			"received": ageInt,
		})
	}

	// 🎯 DEMO: Override log level cho suspicious input
	// Nếu age quá lớn (>150), coi là suspicious và force log vào file
	if ageInt > 150 {
		return goerrorkit.NewValidationError("Suspicious age value detected", map[string]interface{}{
			"field":    "age",
			"received": ageInt,
			"max":      150,
			"reason":   "possible_attack",
		}).Level("error") // ⭐ Override: warn → error (ghi vào file)
	}

	return c.JSON(fiberv2.Map{
		"message": "Validation successful",
		"age":     ageInt,
	})
}

func authErrorHandler(c *fiberv2.Ctx) error {
	token := c.Get("Authorization")

	// Check if token exists
	if token == "" {
		// ⭐ AuthError với default log level "warn"
		// → Console: ✓ (log), File: ✗ (bỏ qua vì FileLogLevel = "error")
		return goerrorkit.NewAuthError(401, "Unauthorized: Missing authorization token")
	}

	// Simulate invalid token
	if token != "Bearer valid-token-123" {
		// ⭐ AuthError với default log level "warn"
		return goerrorkit.NewAuthError(401, "Unauthorized: Invalid token").WithData(map[string]interface{}{
			"token_length": len(token),
		})
	}

	// Simulate permission check
	role := c.Get("X-User-Role")
	if role != "admin" {
		// ⭐ AuthError với default log level "warn"
		return goerrorkit.NewAuthError(403, "Forbidden: Insufficient permissions").WithData(map[string]interface{}{
			"required_role": "admin",
			"user_role":     role,
		})
	}

	return c.JSON(fiberv2.Map{
		"message": "Authentication successful",
		"role":    role,
	})
}

func externalErrorHandler(c *fiberv2.Ctx) error {
	// Simulate external API call failure
	service := c.Query("service", "payment")

	err := fmt.Errorf("timeout after 30s")

	var statusCode int
	var message string

	switch service {
	case "payment":
		statusCode = 502
		message = "Payment gateway not responding"
	case "shipping":
		statusCode = 503
		message = "Shipping service under maintenance"
	case "notification":
		statusCode = 504
		message = "Notification service timeout"
	default:
		statusCode = 502
		message = "External service unavailable"
	}

	return goerrorkit.NewExternalError(statusCode, message, err).WithData(map[string]interface{}{
		"service": service,
		"timeout": "30s",
	})
}

// ============================================================================
// Complex Error Handler - Demonstrate WithCallChain()
// ============================================================================

// complexErrorWithCallChainHandler demonstrates using .WithCallChain()
// to add full call chain to non-panic errors for better debugging
func complexErrorWithCallChainHandler(c *fiberv2.Ctx) error {
	// Simulate a complex operation with multiple function calls
	result, err := processOrder()
	if err != nil {
		return err
	}

	return c.JSON(fiberv2.Map{
		"message": "Order processed",
		"result":  result,
	})
}

func processOrder() (string, error) {
	// Call validation
	if err := validateOrder(); err != nil {
		return "", err
	}

	// Call inventory check
	if err := checkInventory(); err != nil {
		return "", err
	}

	return "success", nil
}

func validateOrder() error {
	// Simulate validation
	isValid := false

	if !isValid {
		// ⭐ Sử dụng .WithCallChain() để thêm full call chain
		// Giúp trace được: complexErrorWithCallChainHandler → processOrder → validateOrder
		return goerrorkit.NewValidationError("Order validation failed", map[string]interface{}{
			"reason": "invalid_order_data",
		}).WithCallChain() // ⭐ Thêm call_chain vào error!
	}

	return nil
}

func checkInventory() error {
	// Simulate inventory check
	stockAvailable := 0

	if stockAvailable == 0 {
		// ⭐ Chain nhiều methods: WithData() + WithCallChain()
		return goerrorkit.NewBusinessError(422, "Insufficient inventory").
			WithData(map[string]interface{}{
				"product_id": "PROD-123",
				"requested":  10,
				"available":  0,
				"warehouse":  "WH-01",
			}).
			WithCallChain() // ⭐ Thêm call_chain để trace flow
	}

	return nil
}

// ============================================================================
// Log Level Demo Handler - Showcase .Level() fluent API (NEW!)
// ============================================================================

// logLevelDemoHandler demonstrates how to override log level with .Level()
// 🎯 MỤC ĐÍCH: Show sự khác biệt giữa warn và error log levels
//
// Use cases:
// 1. ?level=warn   → Log ra console, KHÔNG log vào file (vì FileLogLevel = "error")
// 2. ?level=error  → Log ra cả console VÀ file
// 3. ?level=panic  → Log ra cả console VÀ file (treated as error trong logrus)
func logLevelDemoHandler(c *fiberv2.Ctx) error {
	level := c.Query("level", "warn")
	scenario := c.Query("scenario", "validation")

	switch scenario {
	case "validation":
		// ValidationError mặc định có log level = "warn"
		// → Console: ✓, File: ✗
		if level == "warn" {
			return goerrorkit.NewValidationError("Email format invalid", map[string]interface{}{
				"field":    "email",
				"received": "invalid@",
				"reason":   "missing_domain",
			}) // Default level = "warn"
		}

		// Override để log vào file (suspicious input pattern)
		// → Console: ✓, File: ✓
		if level == "error" {
			return goerrorkit.NewValidationError("Suspicious input pattern detected", map[string]interface{}{
				"field":    "email",
				"received": "'; DROP TABLE users; --",
				"reason":   "sql_injection_attempt",
			}).Level("error") // ⭐ Override: warn → error
		}

	case "auth":
		// AuthError mặc định có log level = "warn"
		if level == "warn" {
			return goerrorkit.NewAuthError(401, "Invalid credentials").WithData(map[string]interface{}{
				"username":      "john@example.com",
				"failed_at":     "2025-11-28T10:30:00Z",
				"attempt_count": 1,
			}) // Default level = "warn" → Console: ✓, File: ✗
		}

		// Multiple failed attempts → upgrade to error level
		if level == "error" {
			return goerrorkit.NewAuthError(401, "Multiple failed login attempts").
				WithData(map[string]interface{}{
					"username":      "john@example.com",
					"attempt_count": 5,
					"ip_address":    "192.168.1.100",
					"reason":        "possible_brute_force",
				}).
				Level("error") // ⭐ Override: warn → error (cần investigate)
		}

	case "business":
		// BusinessError mặc định có log level = "error"
		if level == "warn" {
			// Downgrade từ error → warn (optional, rare case)
			return goerrorkit.NewBusinessError(404, "Product temporarily unavailable").
				WithData(map[string]interface{}{
					"product_id": "PROD-456",
					"status":     "out_of_stock",
				}).
				Level("warn") // ⭐ Override: error → warn
		}

		// Giữ nguyên error level (default)
		return goerrorkit.NewBusinessError(500, "Critical business error").
			WithData(map[string]interface{}{
				"product_id": "PROD-789",
				"stock":      -5, // Negative stock!
				"reason":     "data_corruption",
			}) // Default level = "error" → Console: ✓, File: ✓
	}

	return c.JSON(fiberv2.Map{
		"message": "No error triggered",
		"hint":    "Try ?level=warn or ?level=error with &scenario=validation/auth/business",
	})
}

// ============================================================================
// Development Tools - Trace & Debug Handlers (NEW!)
// ============================================================================

// traceHandler demonstrates simple trace logging for a single operation
// 🎯 USE CASE: Track một operation đơn giản trong development
// ⭐ Trace level thường chỉ dùng trong dev, không nên log vào file production
//
// ⚠️  QUAN TRỌNG: Handler này CHỈ hoạt động khi build với -tags=debug
//   - Chạy: go run -tags=debug main.go
//   - Nếu chạy không có tag, Trace() sẽ là no-op (không làm gì)
func traceHandler(c *fiberv2.Ctx) error {
	operation := c.Query("op", "fetch_user")

	switch operation {
	case "fetch_user":
		// Giả lập fetch user từ database
		userID := c.Query("user_id", "12345")

		// ⭐ Trace log - Không phải error, chỉ để track flow
		// Level: "info" hoặc "debug" (tùy implementation)
		goerrorkit.Trace("Fetching user from database", map[string]interface{}{
			"user_id": userID,
		})

		// Simulate successful fetch
		return c.JSON(fiberv2.Map{
			"message": "User fetched successfully",
			"user_id": userID,
			"trace":   "Check console for trace log",
		})

	case "cache_miss":
		// Giả lập cache miss scenario
		key := c.Query("key", "user:12345")

		// ⭐ Trace cache miss (not an error, just tracking)
		goerrorkit.Trace("Cache miss", map[string]interface{}{
			"key":    key,
			"action": "fetch_from_db",
		})

		return c.JSON(fiberv2.Map{
			"message": "Cache miss - fetched from database",
			"key":     key,
			"trace":   "Cache miss event traced in console",
		})

	case "slow_query":
		// Giả lập slow query warning
		query := c.Query("query", "SELECT * FROM users")
		duration := "2.5s"

		// ⭐ Trace slow query (warning, not error)
		goerrorkit.Trace("Slow query detected", map[string]interface{}{
			"duration": duration,
			"query":    query,
		})

		return c.JSON(fiberv2.Map{
			"message":  "Query executed but slow",
			"duration": duration,
			"trace":    "Slow query traced in console",
		})
	}

	return c.JSON(fiberv2.Map{
		"message": "Unknown operation",
		"hint":    "Try ?op=fetch_user, ?op=cache_miss, or ?op=slow_query",
	})
}

// debugHandler demonstrates debug logging with detailed context
// 🎯 USE CASE: Log chi tiết variable states, object properties trong development
// ⭐ Debug logs giúp hiểu rõ state của application tại một thời điểm
//
// ⚠️  QUAN TRỌNG: Handler này CHỈ hoạt động khi build với -tags=debug
//   - Chạy: go run -tags=debug main.go
//   - Nếu chạy không có tag, Debug() sẽ là no-op (không làm gì)
func debugHandler(c *fiberv2.Ctx) error {
	scenario := c.Query("scenario", "user_login")

	switch scenario {
	case "user_login":
		// Giả lập user login flow với debug info
		username := c.Query("username", "john@example.com")

		// ⭐ Debug log - Log detailed state
		goerrorkit.Debug("User login attempt", map[string]interface{}{
			"username":   username,
			"ip_address": c.IP(),
			"user_agent": c.Get("User-Agent"),
			"timestamp":  "2025-11-28T10:30:00Z",
		})

		return c.JSON(fiberv2.Map{
			"message": "Login successful",
			"debug":   "Check console for detailed debug logs",
		})

	case "payment_process":
		// Giả lập payment processing với debug info
		amount := c.Query("amount", "100000")
		currency := c.Query("currency", "VND")

		// ⭐ Debug log - Track payment state
		goerrorkit.Debug("Processing payment", map[string]interface{}{
			"amount":         amount,
			"currency":       currency,
			"gateway":        "stripe",
			"customer_id":    "cust_123456",
			"payment_method": "card_****1234",
			"state":          "validating → processing → completed",
		})

		return c.JSON(fiberv2.Map{
			"message": "Payment processed",
			"debug":   "Check console for payment flow debug logs",
		})

	case "api_request":
		// Giả lập external API request với debug info
		service := c.Query("service", "user-service")

		// ⭐ Debug log - Track API request/response
		goerrorkit.Debug("External API call", map[string]interface{}{
			"service":       service,
			"endpoint":      "https://api.example.com/users/123",
			"method":        "GET",
			"headers":       "Authorization: Bearer ***, Content-Type: application/json",
			"request_id":    "req_abc123",
			"response_time": "150ms",
			"status_code":   200,
		})

		return c.JSON(fiberv2.Map{
			"message": "API call successful",
			"debug":   "Check console for API request/response debug logs",
		})
	}

	return c.JSON(fiberv2.Map{
		"message": "Unknown scenario",
		"hint":    "Try ?scenario=user_login, ?scenario=payment_process, or ?scenario=api_request",
	})
}

// traceComplexFlowHandler demonstrates tracing a complex multi-step operation
// 🎯 USE CASE: Trace toàn bộ flow của một operation phức tạp với nhiều steps
// ⭐ Giúp hiểu rõ flow execution và identify performance bottlenecks
//
// ⚠️  QUAN TRỌNG: Handler này CHỈ hoạt động khi build với -tags=debug
//   - Chạy: go run -tags=debug main.go
//   - Nếu chạy không có tag, Trace() sẽ là no-op (không làm gì)
func traceComplexFlowHandler(c *fiberv2.Ctx) error {
	orderID := c.Query("order_id", "ORD-12345")

	// ⭐ Start trace
	goerrorkit.Trace("Order Processing Flow Started", map[string]interface{}{
		"order_id":  orderID,
		"timestamp": "2025-11-28T10:30:00Z",
	})

	// Step 1: Validate order
	goerrorkit.Trace("Step 1: Validating order", map[string]interface{}{
		"order_exists":   true,
		"customer_id":    "CUST-456",
		"payment_method": "valid",
		"duration_ms":    50,
	})

	// Step 2: Check inventory
	goerrorkit.Trace("Step 2: Checking inventory", map[string]interface{}{
		"product_id":    "PROD-789",
		"requested_qty": 2,
		"available_qty": 10,
		"stock_status":  "available",
		"duration_ms":   120,
	})

	// Step 3: Reserve inventory
	goerrorkit.Trace("Step 3: Reserving inventory", map[string]interface{}{
		"warehouse":      "WH-01",
		"reservation_id": "RES-999",
		"status":         "reserved",
		"duration_ms":    80,
	})

	// Step 4: Process payment
	goerrorkit.Trace("Step 4: Processing payment", map[string]interface{}{
		"amount":         "200,000 VND",
		"gateway":        "stripe",
		"transaction_id": "TXN-111",
		"status":         "captured",
		"duration_ms":    450,
	})

	// Step 5: Create shipment
	goerrorkit.Trace("Step 5: Creating shipment", map[string]interface{}{
		"carrier":            "DHL",
		"tracking_number":    "DHL123456789",
		"estimated_delivery": "2025-12-02",
		"status":             "created",
		"duration_ms":        200,
	})

	// Step 6: Send confirmation
	goerrorkit.Trace("Step 6: Sending confirmation email", map[string]interface{}{
		"to":          "customer@example.com",
		"template":    "order_confirmation",
		"status":      "sent",
		"duration_ms": 300,
	})

	// ⭐ End trace with summary
	goerrorkit.Trace("Order Processing Flow Completed", map[string]interface{}{
		"total_duration_ms": 1200,
		"order_status":      "confirmed",
		"all_steps":         "successful",
	})

	return c.JSON(fiberv2.Map{
		"message":         "Order processed successfully",
		"order_id":        orderID,
		"status":          "confirmed",
		"tracking_number": "DHL123456789",
		"trace":           "Check console for detailed flow trace (6 steps)",
		"total_duration":  "1,200ms",
	})
}

// ============================================================================
// Database Migration Helper - Demo error logging với cause
// ============================================================================

// runDatabaseMigrations giả lập database migration với khả năng thành công hoặc thất bại
// simulateError = true → trả về lỗi migration với cause để demo logging
// simulateError = false → migration thành công
func runDatabaseMigrations(simulateError bool) error {
	if !simulateError {
		// Migration thành công
		fmt.Println("✅ Database migrations completed successfully")
		return nil
	}

	// Giả lập lỗi kết nối database (lỗi gốc/cause)
	dbConnectionErr := fmt.Errorf("dial tcp 127.0.0.1:5432: connect: connection refused")

	// ⭐ WrapWithMessage() - Wrap error gốc với message mô tả context
	// → Message: "Failed to run database migrations"
	// → Cause: "dial tcp 127.0.0.1:5432: connect: connection refused"
	// → Type: SystemError, Code: 500
	// → Tự động capture: file, line, function, stack trace
	migrationErr := goerrorkit.WrapWithMessage(dbConnectionErr, "Failed to run database migrations").
		WithData(map[string]interface{}{
			"database":        "postgresql",
			"host":            "127.0.0.1:5432",
			"migration_files": []string{"001_create_users.sql", "002_create_products.sql"},
			"last_version":    0,
			"target_version":  2,
		})

	return migrationErr
}
