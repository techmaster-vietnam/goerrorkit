# Configuration Guide

## Logger Configuration

### Default Configuration

```go
config.InitDefaultLogger()
```

Equivalent to:

```go
config.InitLogger(config.LoggerOptions{
    ConsoleOutput: true,
    FileOutput:    true,
    FilePath:      "logs/errors.log",
    JSONFormat:    true,
    MaxFileSize:   10,    // MB
    MaxBackups:    5,
    MaxAge:        30,    // days
    LogLevel:      "error",
})
```

### Custom Configuration

#### Console Only

```go
config.InitLogger(config.LoggerOptions{
    ConsoleOutput: true,
    FileOutput:    false,
    JSONFormat:    false, // Text format with colors
    LogLevel:      "debug",
})
```

#### File Only (Production)

```go
config.InitLogger(config.LoggerOptions{
    ConsoleOutput: false,
    FileOutput:    true,
    FilePath:      "/var/log/app/errors.log",
    JSONFormat:    true,
    MaxFileSize:   50,    // 50MB per file
    MaxBackups:    10,    // Keep 10 backups
    MaxAge:        90,    // Keep for 90 days
    LogLevel:      "error",
})
```

#### Both Console and File

```go
config.InitLogger(config.LoggerOptions{
    ConsoleOutput: true,  // Development: see logs in terminal
    FileOutput:    true,  // Production: persist to file
    FilePath:      "logs/app.log",
    JSONFormat:    true,
    LogLevel:      "info", // Log info and above
})
```

## Stack Trace Configuration

### Simple Configuration

```go
// Auto-configure for your application package
core.ConfigureForApplication("github.com/yourname/yourapp")
```

### Advanced Configuration

```go
core.SetStackTraceConfig(core.StackTraceConfig{
    // Only include these packages in stack trace
    IncludePackages: []string{
        "github.com/yourname/yourapp",
        "github.com/yourname/yourlib",
        "main", // For local development
    },
    
    // Skip these packages
    SkipPackages: []string{
        "runtime",
        "runtime/debug",
        "github.com/some/library/internal",
    },
    
    // Skip functions matching these names
    SkipFunctions: []string{
        "formatStackTraceArray",
        "HandlePanic",
        ".func", // Anonymous functions
    },
    
    // Show full path or short name
    // true:  github.com/yourname/yourapp.Handler
    // false: yourapp.Handler
    ShowFullPath: false,
})
```

### Multiple Application Packages

```go
core.SetStackTraceConfig(core.StackTraceConfig{
    IncludePackages: []string{
        "github.com/yourname/app",      // Main app
        "github.com/yourname/services", // Services package
        "github.com/yourname/models",   // Models package
    },
    ShowFullPath: false,
})
```

## Custom Logger Implementation

Bạn có thể implement interface `core.Logger` để dùng logger khác (zap, zerolog, etc.):

```go
import "github.com/cuong/goerrorkit/core"

type MyCustomLogger struct {
    // Your logger instance
}

func (l *MyCustomLogger) Error(msg string, fields map[string]interface{}) {
    // Your implementation
}

func (l *MyCustomLogger) Info(msg string, fields map[string]interface{}) {
    // Your implementation
}

func (l *MyCustomLogger) Debug(msg string, fields map[string]interface{}) {
    // Your implementation
}

func (l *MyCustomLogger) Warn(msg string, fields map[string]interface{}) {
    // Your implementation
}

// Usage
func main() {
    logger := &MyCustomLogger{}
    core.SetLogger(logger)
}
```

## Environment-based Configuration

### Development

```go
func initLogger() {
    if os.Getenv("ENV") == "development" {
        config.InitLogger(config.LoggerOptions{
            ConsoleOutput: true,
            FileOutput:    false,
            JSONFormat:    false, // Text format easier to read
            LogLevel:      "debug",
        })
    } else {
        // Production config
        config.InitLogger(config.LoggerOptions{
            ConsoleOutput: false,
            FileOutput:    true,
            FilePath:      "/var/log/app/errors.log",
            JSONFormat:    true,
            LogLevel:      "error",
        })
    }
}
```

### Using Config File

```go
import "encoding/json"

type Config struct {
    Logger config.LoggerOptions `json:"logger"`
}

func loadConfig() {
    data, _ := os.ReadFile("config.json")
    var cfg Config
    json.Unmarshal(data, &cfg)
    
    config.InitLogger(cfg.Logger)
}
```

config.json:
```json
{
  "logger": {
    "ConsoleOutput": true,
    "FileOutput": true,
    "FilePath": "logs/errors.log",
    "JSONFormat": true,
    "MaxFileSize": 10,
    "MaxBackups": 5,
    "MaxAge": 30,
    "LogLevel": "error"
  }
}
```

## Log Levels

Available log levels (from lowest to highest priority):

- `debug` - Detailed information for debugging
- `info` - General informational messages
- `warn` - Warning messages
- `error` - Error messages (default)

Set `LogLevel` to the minimum level you want to log. For example:
- `LogLevel: "error"` - Only log errors
- `LogLevel: "info"` - Log info, warn, and error
- `LogLevel: "debug"` - Log everything

## Best Practices

1. **Development**: Console output, text format, debug level
2. **Staging**: Both console and file, JSON format, info level
3. **Production**: File only, JSON format, error level
4. **Stack Trace**: Always configure with your application package
5. **File Rotation**: Use reasonable limits (10-50MB per file, keep 5-10 backups)

