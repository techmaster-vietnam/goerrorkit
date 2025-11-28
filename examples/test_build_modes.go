package main

import (
	"fmt"

	"github.com/techmaster-vietnam/goerrorkit"
)

// TestBuildModes demo sự khác biệt giữa debug build và production build
//
// Cách test:
//
//  1. Production build (mặc định): go run test_build_modes.go
//     → Debug/trace logs sẽ KHÔNG in ra
//
//  2. Debug build: go run -tags=debug test_build_modes.go
//     → Debug/trace logs sẽ in ra đầy đủ
func TestBuildModes() {
	fmt.Println("\n=== GoErrorKit Build Modes Demo ===")

	// Khởi tạo logger với debug level
	goerrorkit.InitLogger(goerrorkit.LoggerOptions{
		ConsoleOutput: true,
		FileOutput:    false,
		JSONFormat:    false,   // Text format dễ đọc hơn cho demo
		LogLevel:      "trace", // Set trace level (thấp nhất)
	})

	logger := goerrorkit.GetLogger()

	fmt.Println("📝 Testing all log levels:")

	// Test các log levels
	logger.Trace("TRACE level message", map[string]interface{}{
		"note":  "Chỉ hiển thị khi build với -tags=debug",
		"level": "trace",
	})

	logger.Debug("DEBUG level message", map[string]interface{}{
		"note":  "Chỉ hiển thị khi build với -tags=debug",
		"level": "debug",
	})

	logger.Info("INFO level message", map[string]interface{}{
		"note":  "Luôn hiển thị (production và debug)",
		"level": "info",
	})

	logger.Warn("WARN level message", map[string]interface{}{
		"note":  "Luôn hiển thị (production và debug)",
		"level": "warn",
	})

	logger.Error("ERROR level message", map[string]interface{}{
		"note":  "Luôn hiển thị (production và debug)",
		"level": "error",
	})

	fmt.Println("\n=== Kết luận ===")
	fmt.Println("✅ Production build: Chỉ thấy INFO, WARN, ERROR")
	fmt.Println("✅ Debug build (-tags=debug): Thấy tất cả TRACE, DEBUG, INFO, WARN, ERROR")
	fmt.Println("\n💡 Để test:")
	fmt.Println("   Production: go run test_build_modes.go")
	fmt.Println("   Debug:      go run -tags=debug test_build_modes.go")
}

// Uncomment dòng dưới nếu muốn chạy test này độc lập
// func main() {
//     TestBuildModes()
// }
