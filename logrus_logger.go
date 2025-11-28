package goerrorkit

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// LogrusLogger implement Logger interface sử dụng logrus
// Hỗ trợ dual-level logging: console và file có thể có log level khác nhau
type LogrusLogger struct {
	consoleLogger *logrus.Logger // Logger cho console
	fileLogger    *logrus.Logger // Logger cho file (có thể nil nếu không dùng file)
}

// Error implements Logger
func (l *LogrusLogger) Error(msg string, fields map[string]interface{}) {
	if l.consoleLogger != nil {
		l.consoleLogger.WithFields(fields).Error(msg)
	}
	if l.fileLogger != nil {
		l.fileLogger.WithFields(fields).Error(msg)
	}
}

// Info implements Logger
func (l *LogrusLogger) Info(msg string, fields map[string]interface{}) {
	if l.consoleLogger != nil {
		l.consoleLogger.WithFields(fields).Info(msg)
	}
	if l.fileLogger != nil {
		l.fileLogger.WithFields(fields).Info(msg)
	}
}

// Debug và Trace methods được implement trong:
// - logrus_logger_debug.go (với build tag 'debug') - có đầy đủ logging
// - logrus_logger_prod.go (mặc định, không có tag) - no-op cho performance

// Warn implements Logger
func (l *LogrusLogger) Warn(msg string, fields map[string]interface{}) {
	if l.consoleLogger != nil {
		l.consoleLogger.WithFields(fields).Warn(msg)
	}
	if l.fileLogger != nil {
		l.fileLogger.WithFields(fields).Warn(msg)
	}
}

// Panic implements Logger
func (l *LogrusLogger) Panic(msg string, fields map[string]interface{}) {
	if l.consoleLogger != nil {
		l.consoleLogger.WithFields(fields).Error(msg) // Log as Error, not Panic (không muốn panic thật)
	}
	if l.fileLogger != nil {
		l.fileLogger.WithFields(fields).Error(msg)
	}
}

// LoggerOptions cấu hình cho logger
type LoggerOptions struct {
	// ConsoleOutput - Log ra console hay không
	ConsoleOutput bool

	// FileOutput - Log ra file hay không
	FileOutput bool

	// FilePath - Đường dẫn file log
	FilePath string

	// JSONFormat - Dùng JSON format hay text format
	JSONFormat bool

	// MaxFileSize - Kích thước tối đa của file log (MB) trước khi rotate
	MaxFileSize int

	// MaxBackups - Số lượng file backup giữ lại
	MaxBackups int

	// MaxAge - Số ngày giữ file log cũ
	MaxAge int

	// LogLevel - Level tối thiểu để log ra console (trace, debug, info, warn, error, panic)
	// LƯU Ý: trace và debug chỉ hoạt động khi build với -tags=debug
	//        Production build sẽ bỏ qua hoàn toàn (zero overhead)
	LogLevel string

	// FileLogLevel - Level tối thiểu để log ra file (trace, debug, info, warn, error, panic)
	// Mặc định sẽ dùng "error" để chỉ log các lỗi nghiêm trọng vào file
	// VD: FileLogLevel = "error" -> chỉ log error và panic vào file, bỏ qua warn
	// LƯU Ý: trace và debug chỉ hoạt động khi build với -tags=debug
	FileLogLevel string
}

// DefaultLoggerOptions trả về cấu hình mặc định
func DefaultLoggerOptions() LoggerOptions {
	return LoggerOptions{
		ConsoleOutput: true,
		FileOutput:    true,
		FilePath:      "logs/errors.log",
		JSONFormat:    true,
		MaxFileSize:   10,
		MaxBackups:    5,
		MaxAge:        30,
		LogLevel:      "warn",  // Console log tất cả từ warn trở lên
		FileLogLevel:  "error", // File chỉ log error và panic (bỏ qua warn)
	}
}

// InitLogger khởi tạo logger với custom options
// Hỗ trợ dual-level logging: console và file có thể có log level khác nhau
//
// LƯU Ý VỀ DEBUG/TRACE LOGS:
//   - Debug và trace logs CHỈ hoạt động khi build với -tags=debug
//   - Production build (không có tag): debug/trace là no-op (zero overhead)
//   - Cách build: go build -tags=debug    (development)
//     go build                  (production)
//
// Example:
//
//	// Development: LogLevel="debug" sẽ log debug messages
//	// Production: LogLevel="debug" sẽ KHÔNG log gì (no-op)
//	goerrorkit.InitLogger(goerrorkit.LoggerOptions{
//	    ConsoleOutput: true,
//	    FileOutput: true,
//	    FilePath: "logs/app.log",
//	    JSONFormat: true,
//	    LogLevel: "debug",     // Development: log debug, Production: no-op
//	    FileLogLevel: "error", // File chỉ log error và panic
//	})
func InitLogger(opts LoggerOptions) {
	var consoleLogger *logrus.Logger
	var fileLogger *logrus.Logger

	// Khởi tạo console logger
	if opts.ConsoleOutput {
		consoleLogger = logrus.New()
		consoleLogger.SetOutput(os.Stdout)

		// Cấu hình formatter cho console
		if opts.JSONFormat {
			consoleLogger.SetFormatter(&logrus.JSONFormatter{
				TimestampFormat: time.RFC3339,
				PrettyPrint:     true,
				FieldMap: logrus.FieldMap{
					logrus.FieldKeyTime:  "timestamp",
					logrus.FieldKeyLevel: "level",
					logrus.FieldKeyMsg:   "message",
				},
			})
		} else {
			consoleLogger.SetFormatter(&logrus.TextFormatter{
				ForceColors:     true,
				FullTimestamp:   true,
				TimestampFormat: "2006-01-02 15:04:05",
			})
		}

		// Set log level cho console
		level, err := logrus.ParseLevel(opts.LogLevel)
		if err != nil {
			level = logrus.WarnLevel // Default warn cho console
		}
		consoleLogger.SetLevel(level)
	}

	// Khởi tạo file logger
	if opts.FileOutput {
		// Tạo thư mục logs nếu chưa có
		if err := os.MkdirAll("logs", 0755); err != nil {
			if consoleLogger != nil {
				consoleLogger.Errorf("Cannot create logs directory: %v", err)
			}
		}

		fileLogger = logrus.New()
		logFile := &lumberjack.Logger{
			Filename:   opts.FilePath,
			MaxSize:    opts.MaxFileSize,
			MaxBackups: opts.MaxBackups,
			MaxAge:     opts.MaxAge,
			Compress:   true,
			LocalTime:  true,
		}
		fileLogger.SetOutput(logFile)

		// Cấu hình formatter cho file (luôn dùng JSON)
		fileLogger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
			PrettyPrint:     true,
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "timestamp",
				logrus.FieldKeyLevel: "level",
				logrus.FieldKeyMsg:   "message",
			},
		})

		// Set log level cho file
		fileLogLevel := opts.FileLogLevel
		if fileLogLevel == "" {
			fileLogLevel = "error" // Default error cho file
		}
		fileLevel, err := logrus.ParseLevel(fileLogLevel)
		if err != nil {
			fileLevel = logrus.ErrorLevel
		}
		fileLogger.SetLevel(fileLevel)
	}

	// Wrap và set vào goerrorkit
	logrusLogger := &LogrusLogger{
		consoleLogger: consoleLogger,
		fileLogger:    fileLogger,
	}
	SetLogger(logrusLogger)

	if consoleLogger != nil {
		consoleLogger.Info("✓ GoErrorKit logger initialized")
	}
}

// InitDefaultLogger khởi tạo logger với cấu hình mặc định
//
// Example:
//
//	goerrorkit.InitDefaultLogger()
func InitDefaultLogger() {
	InitLogger(DefaultLoggerOptions())
}
