package main

import (
	"fmt"

	"github.com/cuong/goerrorkit/adapters/fiber"
	"github.com/cuong/goerrorkit/config"
	"github.com/cuong/goerrorkit/core"
	fiberv2 "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func main() {
	// 1. Initialize logger với custom options
	config.InitLogger(config.LoggerOptions{
		ConsoleOutput: true,
		FileOutput:    true,
		FilePath:      "logs/errors.log",
		JSONFormat:    true,
		MaxFileSize:   10,
		MaxBackups:    5,
		MaxAge:        30,
		LogLevel:      "info",
	})

	// 2. Configure stack trace for this application
	// Replace "github.com/cuong/goerrorkit/examples/fiber-demo" with your app package
	core.ConfigureForApplication("main")

	// 3. Create Fiber app
	app := fiberv2.New(fiberv2.Config{
		AppName: "GoErrorKit Demo",
	})

	// 4. Add middlewares (RequestID must be before ErrorHandler)
	app.Use(requestid.New())
	app.Use(logger.New())
	app.Use(fiber.ErrorHandler())

	// 5. Routes - Demo different error types
	app.Get("/", homeHandler)

	// Panic demos
	app.Get("/panic/division", panicDivisionHandler)
	app.Get("/panic/index", panicIndexHandler)
	app.Get("/panic/stack", panicStackHandler)

	// Custom error demos
	app.Get("/error/business", businessErrorHandler)
	app.Get("/error/system", systemErrorHandler)
	app.Get("/error/validation", validationErrorHandler)
	app.Get("/error/auth", authErrorHandler)
	app.Get("/error/external", externalErrorHandler)

	// Start server
	fmt.Println("🚀 Server starting on http://localhost:8081")
	fmt.Println("\n📝 Try these endpoints:")
	fmt.Println("  GET  /                     - Home")
	fmt.Println("\n  🔥 Panic Demos (auto-recovered with exact location):")
	fmt.Println("  GET  /panic/division       - Division by zero panic")
	fmt.Println("  GET  /panic/index          - Index out of range panic")
	fmt.Println("  GET  /panic/stack          - Deep call stack panic")
	fmt.Println("\n  ⚠️  Custom Error Demos:")
	fmt.Println("  GET  /error/business       - Business logic error (404)")
	fmt.Println("  GET  /error/system         - System error (500)")
	fmt.Println("  GET  /error/validation     - Validation error (400)")
	fmt.Println("  GET  /error/auth           - Auth error (401)")
	fmt.Println("  GET  /error/external       - External service error (502)")
	fmt.Println("\n📄 Check logs/errors.log for detailed error logs\n")

	if err := app.Listen(":8081"); err != nil {
		panic(err)
	}
}

func homeHandler(c *fiberv2.Ctx) error {
	return c.JSON(fiberv2.Map{
		"message": "Welcome to GoErrorKit Demo!",
		"status":  "ok",
	})
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
// Custom Error Handlers - Demonstrate different error types
// ============================================================================

func businessErrorHandler(c *fiberv2.Ctx) error {
	productID := c.Query("product_id", "unknown")

	// Simulate product not found
	if productID == "123" {
		return core.NewBusinessError(404, fmt.Sprintf("Product ID=%s not found", productID))
	}

	return c.JSON(fiberv2.Map{
		"message":    "Product available",
		"product_id": productID,
	})
}

func systemErrorHandler(c *fiberv2.Ctx) error {
	// Simulate database connection error
	err := fmt.Errorf("connection refused: database is down")
	return core.NewSystemError(err)
}

func validationErrorHandler(c *fiberv2.Ctx) error {
	age := c.Query("age", "")

	if age == "" {
		return core.NewValidationError("Missing parameter 'age'", map[string]interface{}{
			"field":    "age",
			"required": true,
		})
	}

	// Check if age is a number
	var ageInt int
	if _, err := fmt.Sscanf(age, "%d", &ageInt); err != nil {
		return core.NewValidationError("Parameter 'age' must be an integer", map[string]interface{}{
			"field":    "age",
			"type":     "integer",
			"received": age,
		})
	}

	if ageInt < 18 {
		return core.NewValidationError("Age must be >= 18", map[string]interface{}{
			"field":    "age",
			"min":      18,
			"received": ageInt,
		})
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
		return core.NewAuthError(401, "Unauthorized: Missing authorization token")
	}

	// Simulate invalid token
	if token != "Bearer valid-token-123" {
		return core.NewAuthError(401, "Unauthorized: Invalid token")
	}

	// Simulate permission check
	role := c.Get("X-User-Role")
	if role != "admin" {
		return core.NewAuthError(403, "Forbidden: Insufficient permissions")
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

	return core.NewExternalError(statusCode, message, err)
}
