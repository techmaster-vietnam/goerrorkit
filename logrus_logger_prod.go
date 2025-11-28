//go:build !debug
// +build !debug

package goerrorkit

// Debug implements Logger - PRODUCTION MODE: No-op
// Không làm gì cả trong production build để tối ưu performance
// Code này sẽ được compiler optimize away hoàn toàn
func (l *LogrusLogger) Debug(msg string, fields map[string]interface{}) {
	// No-op: Hoàn toàn không log trong production
}

// Trace implements Logger - PRODUCTION MODE: No-op
// Không làm gì cả trong production build để tối ưu performance
// Code này sẽ được compiler optimize away hoàn toàn
func (l *LogrusLogger) Trace(msg string, fields map[string]interface{}) {
	// No-op: Hoàn toàn không log trong production
}
