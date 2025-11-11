package goerrorkit

import (
	"fmt"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"
)

// StackTraceConfig cấu hình cho stack trace filtering
type StackTraceConfig struct {
	// SkipPackages - Danh sách packages cần bỏ qua (thư viện internal, runtime)
	SkipPackages []string

	// SkipFunctions - Danh sách functions cần bỏ qua (utility functions)
	SkipFunctions []string

	// IncludePackages - Danh sách packages cần include (application code)
	// Nếu empty, sẽ include tất cả trừ SkipPackages
	IncludePackages []string

	// ShowFullPath - Hiển thị full package path hay chỉ tên cuối
	// true: github.com/user/myapp.Handler
	// false: myapp.Handler
	ShowFullPath bool
}

// defaultConfig là cấu hình mặc định cho stack trace
var defaultConfig = StackTraceConfig{
	SkipPackages: []string{
		"runtime",
		"runtime/debug",
	},
	SkipFunctions: []string{
		"formatStackTraceArray",
		"getActualPanicLocation",
		"HandlePanic",
		"ErrorHandler",
		"middleware",
	},
	IncludePackages: []string{},
	ShowFullPath:    false,
}

// SetStackTraceConfig cho phép user customize stack trace behavior
//
// Example:
//
//	goerrorkit.SetStackTraceConfig(goerrorkit.StackTraceConfig{
//	    IncludePackages: []string{"github.com/yourname/myapp"},
//	    ShowFullPath: false,
//	})
func SetStackTraceConfig(config StackTraceConfig) {
	defaultConfig = config
}

// ConfigureForApplication là helper function để config nhanh cho application
//
// Example:
//
//	goerrorkit.ConfigureForApplication("github.com/yourname/myapp")
func ConfigureForApplication(appPackage string) {
	defaultConfig.IncludePackages = []string{appPackage}
	// Auto-skip thư viện goerrorkit
	defaultConfig.SkipPackages = append(defaultConfig.SkipPackages,
		"github.com/techmaster-vietnam/goerrorkit",
	)
}

// StackTraceConfigurator cung cấp fluent API để configure stack trace
type StackTraceConfigurator struct {
	config StackTraceConfig
}

// Configure trả về một configurator mới với config hiện tại
//
// Example:
//
//	goerrorkit.Configure().
//	    SkipPackage("mylib").
//	    SkipFunction("helper").
//	    SkipPattern(".RequestID.func").
//	    Apply()
func Configure() *StackTraceConfigurator {
	// Copy config hiện tại để tránh modify trực tiếp
	return &StackTraceConfigurator{
		config: StackTraceConfig{
			SkipPackages:    append([]string{}, defaultConfig.SkipPackages...),
			SkipFunctions:   append([]string{}, defaultConfig.SkipFunctions...),
			IncludePackages: append([]string{}, defaultConfig.IncludePackages...),
			ShowFullPath:    defaultConfig.ShowFullPath,
		},
	}
}

// SkipPackage thêm package cần bỏ qua vào danh sách
func (c *StackTraceConfigurator) SkipPackage(pkg string) *StackTraceConfigurator {
	c.config.SkipPackages = append(c.config.SkipPackages, pkg)
	return c
}

// SkipPackages thêm nhiều packages cần bỏ qua
func (c *StackTraceConfigurator) SkipPackages(pkgs ...string) *StackTraceConfigurator {
	c.config.SkipPackages = append(c.config.SkipPackages, pkgs...)
	return c
}

// SkipFunction thêm function pattern cần bỏ qua
func (c *StackTraceConfigurator) SkipFunction(fn string) *StackTraceConfigurator {
	c.config.SkipFunctions = append(c.config.SkipFunctions, fn)
	return c
}

// SkipFunctions thêm nhiều function patterns cần bỏ qua
func (c *StackTraceConfigurator) SkipFunctions(fns ...string) *StackTraceConfigurator {
	c.config.SkipFunctions = append(c.config.SkipFunctions, fns...)
	return c
}

// SkipPattern là alias cho SkipFunction (semantic clarity)
func (c *StackTraceConfigurator) SkipPattern(pattern string) *StackTraceConfigurator {
	return c.SkipFunction(pattern)
}

// SkipPatterns thêm nhiều patterns cần bỏ qua
func (c *StackTraceConfigurator) SkipPatterns(patterns ...string) *StackTraceConfigurator {
	return c.SkipFunctions(patterns...)
}

// IncludePackage set package cần include (application code)
func (c *StackTraceConfigurator) IncludePackage(pkg string) *StackTraceConfigurator {
	c.config.IncludePackages = append(c.config.IncludePackages, pkg)
	return c
}

// IncludePackages set nhiều packages cần include
func (c *StackTraceConfigurator) IncludePackages(pkgs ...string) *StackTraceConfigurator {
	c.config.IncludePackages = append(c.config.IncludePackages, pkgs...)
	return c
}

// ShowFullPath bật/tắt hiển thị full package path
func (c *StackTraceConfigurator) ShowFullPath(show bool) *StackTraceConfigurator {
	c.config.ShowFullPath = show
	return c
}

// Apply áp dụng configuration
func (c *StackTraceConfigurator) Apply() {
	defaultConfig = c.config
}

// AddSkipPatterns là shorthand function để nhanh chóng thêm skip patterns
// mà không cần tạo configurator
//
// Example:
//
//	goerrorkit.AddSkipPatterns(".RequestID.func", ".Logger.func", "telemetry")
func AddSkipPatterns(patterns ...string) {
	defaultConfig.SkipFunctions = append(defaultConfig.SkipFunctions, patterns...)
}

// AddSkipPackages là shorthand function để nhanh chóng thêm skip packages
//
// Example:
//
//	goerrorkit.AddSkipPackages("internal/telemetry", "vendor/monitoring")
func AddSkipPackages(packages ...string) {
	defaultConfig.SkipPackages = append(defaultConfig.SkipPackages, packages...)
}

// getActualPanicLocation lấy thông tin về dòng THỰC SỰ gây panic
// Đây là nơi thực sự phát sinh lỗi, không phải nơi gọi hàm
func getActualPanicLocation() (file string, line int, function string) {
	stack := string(debug.Stack())
	lines := strings.Split(stack, "\n")

	// Stack trace format:
	// goroutine X [running]:
	// runtime/debug.Stack()
	//     /path/to/runtime/debug/stack.go:24 +0x65
	// main.getActualPanicLocation()
	//     /path/to/main.go:73 +0x27
	// ...
	// main.GetElement(...)  <- Đây là nơi panic thực sự
	//     /path/to/main.go:117 +0x1f

	for i := 0; i < len(lines); i++ {
		l := strings.TrimSpace(lines[i])

		// Bỏ qua dòng trống
		if l == "" {
			continue
		}

		// Tìm function đầu tiên của user code (không phải runtime/debug và không phải goerrorkit)
		if isUserFunction(l) && !shouldSkipFunction(l) {
			function = l
			// Lấy tên function (bỏ phần parameter)
			if idx := strings.Index(function, "("); idx > 0 {
				function = function[:idx]
			}

			// Format function name
			function = formatFunctionName(function)

			// Dòng tiếp theo chứa file:line
			if i+1 < len(lines) {
				locationLine := strings.TrimSpace(lines[i+1])
				parts := strings.Fields(locationLine)
				if len(parts) > 0 {
					fileAndLine := parts[0]
					if idx := strings.LastIndex(fileAndLine, ":"); idx > 0 {
						file = fileAndLine[:idx]
						fmt.Sscanf(fileAndLine[idx+1:], "%d", &line)
						// Chỉ lấy tên file, bỏ đường dẫn
						file = filepath.Base(file)
					}
				}
			}
			break
		}
	}

	if file == "" {
		return "unknown", 0, "unknown"
	}

	return file, line, function
}

// formatStackTraceArray format stack trace thành array dễ đọc
// Tự động lọc các hàm utility và chỉ lấy application code
func formatStackTraceArray() []string {
	stack := string(debug.Stack())
	lines := strings.Split(stack, "\n")

	var callChain []string
	skipNext := false

	for i := 0; i < len(lines); i++ {
		l := strings.TrimSpace(lines[i])

		if skipNext {
			skipNext = false
			continue
		}

		// Chỉ lấy user functions, bỏ qua utility và runtime
		if isUserFunction(l) && !shouldSkipFunction(l) {
			funcName := l
			// Lấy tên function (bỏ phần parameter)
			if idx := strings.Index(funcName, "("); idx > 0 {
				funcName = funcName[:idx]
			}

			// Format function name
			funcName = formatFunctionName(funcName)

			// Dòng tiếp theo chứa file:line
			if i+1 < len(lines) {
				locationLine := strings.TrimSpace(lines[i+1])
				parts := strings.Fields(locationLine)
				if len(parts) > 0 {
					fileAndLine := parts[0]

					// Chỉ lấy tên file, bỏ đường dẫn đầy đủ
					if idx := strings.LastIndex(fileAndLine, "/"); idx >= 0 {
						fileAndLine = fileAndLine[idx+1:]
					}

					callChain = append(callChain, fmt.Sprintf("%s (%s)", funcName, fileAndLine))
				}
			}
			skipNext = true
		}
	}

	return callChain
}

// getCallerInfo lấy thông tin về nơi gọi factory function
// skip = 1: hàm gọi trực tiếp (default)
// skip = 2: hàm gọi hàm gọi factory function
func getCallerInfo(skip int) (file string, line int, function string) {
	// skip + 1 để bỏ qua chính hàm getCallerInfo này
	pc, file, line, ok := runtime.Caller(skip + 1)
	if !ok {
		return "unknown", 0, "unknown"
	}

	// Lấy tên function
	fn := runtime.FuncForPC(pc)
	if fn != nil {
		function = fn.Name()
	} else {
		function = "unknown"
	}

	// Format function name
	function = formatFunctionName(function)

	// Chỉ lấy tên file, bỏ đường dẫn đầy đủ
	file = filepath.Base(file)

	return file, line, function
}

// isUserFunction kiểm tra xem có phải user code không
func isUserFunction(line string) bool {
	// Bỏ qua dòng trống và các dòng không phải function
	if line == "" || !strings.Contains(line, ".") {
		return false
	}

	// Bỏ qua runtime internal
	if strings.HasPrefix(line, "runtime.") ||
		strings.HasPrefix(line, "runtime/debug.") ||
		strings.Contains(line, "runtime.") ||
		strings.Contains(line, "runtime/debug.") {
		return false
	}

	// Bỏ qua các packages trong SkipPackages
	for _, pkg := range defaultConfig.SkipPackages {
		if strings.HasPrefix(line, pkg+".") || strings.Contains(line, pkg+".") || strings.Contains(line, "/"+pkg+".") {
			return false
		}
	}

	// Nếu có config IncludePackages, chỉ lấy những packages đó
	if len(defaultConfig.IncludePackages) > 0 {
		for _, pkg := range defaultConfig.IncludePackages {
			// Hỗ trợ cả package name và full path
			if strings.HasPrefix(line, pkg+".") ||
				strings.Contains(line, pkg+".") ||
				strings.Contains(line, "/"+pkg+".") {
				return true
			}
		}
		return false
	}

	return true
}

// shouldSkipFunction kiểm tra xem có cần skip function này không
func shouldSkipFunction(line string) bool {
	for _, skipFunc := range defaultConfig.SkipFunctions {
		if strings.Contains(line, skipFunc) {
			return true
		}
	}

	// Skip tất cả anonymous functions có pattern như "package.function.Method.func1"
	// Ví dụ: "main.main.New.func1", "github.com/gofiber/fiber.(*App).Next.func1"
	// Đây thường là middleware wrappers, không phải business logic
	if isMiddlewareAnonymousFunc(line) {
		return true
	}

	return false
}

// isMiddlewareAnonymousFunc kiểm tra xem có phải anonymous function từ middleware không
func isMiddlewareAnonymousFunc(line string) bool {
	// Kiểm tra xem có chứa ".func" (anonymous function) không
	if !strings.Contains(line, ".func") {
		return false
	}

	// Pattern 1: package.package.Method.func (ví dụ: main.main.New.func1)
	// Đếm số lần xuất hiện của package name lặp lại
	parts := strings.Split(line, ".")
	if len(parts) >= 4 {
		// Nếu có >= 4 parts và chứa ".func", nhiều khả năng là middleware
		// Ví dụ: ["main", "main", "New", "func1"] hoặc ["github", "com/gofiber/fiber", "App", "Next", "func1"]
		for i := 0; i < len(parts)-1; i++ {
			if strings.Contains(parts[i+1], "func") {
				return true
			}
		}
	}

	// Pattern 2: Các pattern middleware phổ biến
	middlewarePatterns := []string{
		".main.New.func",   // Fiber middleware setup
		".main.Use.func",   // Middleware use
		".Handler.func",    // Handler wrapper
		".Middleware.func", // Generic middleware
		".Next.func",       // Middleware chain
		".middleware.func", // Lowercase middleware
		".recover.func",    // Recovery middleware
		".logger.func",     // Logger middleware
	}

	for _, pattern := range middlewarePatterns {
		if strings.Contains(line, pattern) {
			return true
		}
	}

	return false
}

// formatFunctionName format function name theo config
func formatFunctionName(fullName string) string {
	if defaultConfig.ShowFullPath {
		return fullName
	}

	// Chỉ lấy package.Function (bỏ full path)
	// github.com/user/app.MyFunc → app.MyFunc
	parts := strings.Split(fullName, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}

	return fullName
}
