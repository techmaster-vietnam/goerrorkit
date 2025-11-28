//go:build debug
// +build debug

package goerrorkit

// Debug implements Logger - CHỈ hoạt động khi build với -tags=debug
// Logs debug level message với fields
func (l *LogrusLogger) Debug(msg string, fields map[string]interface{}) {
	if l.consoleLogger != nil {
		l.consoleLogger.WithFields(fields).Debug(msg)
	}
	if l.fileLogger != nil {
		l.fileLogger.WithFields(fields).Debug(msg)
	}
}

// Trace implements Logger - CHỈ hoạt động khi build với -tags=debug
// Logs trace level message với fields (chi tiết nhất, dùng cho deep debugging)
func (l *LogrusLogger) Trace(msg string, fields map[string]interface{}) {
	if l.consoleLogger != nil {
		l.consoleLogger.WithFields(fields).Trace(msg)
	}
	if l.fileLogger != nil {
		l.fileLogger.WithFields(fields).Trace(msg)
	}
}
