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
		".func", // anonymous functions
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

// getActualPanicLocation lấy thông tin về dòng THỰC SỰ gây panic
// Đây là nơi thực sự phát sinh lỗi, không phải nơi gọi hàm
func getActualPanicLocation() (file string, line int, function string) {
	stack := string(debug.Stack())
	lines := strings.Split(stack, "\n")

	panicFound := false
	for i := 0; i < len(lines); i++ {
		l := strings.TrimSpace(lines[i])

		// Tìm dòng "panic"
		if !panicFound && strings.HasPrefix(l, "panic(") {
			panicFound = true
			continue
		}

		// Sau khi tìm thấy panic, tìm function thực sự gây panic
		if panicFound && isUserFunction(l) {
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
	// Bỏ qua runtime internal
	if strings.HasPrefix(line, "runtime.") ||
		strings.HasPrefix(line, "runtime/debug.") {
		return false
	}

	// Nếu có config IncludePackages, chỉ lấy những packages đó
	if len(defaultConfig.IncludePackages) > 0 {
		for _, pkg := range defaultConfig.IncludePackages {
			if strings.HasPrefix(line, pkg+".") || strings.Contains(line, pkg+".") {
				return true
			}
		}
		return false
	}

	// Không trong SkipPackages
	for _, pkg := range defaultConfig.SkipPackages {
		if strings.HasPrefix(line, pkg+".") || strings.Contains(line, pkg+".") {
			return false
		}
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
