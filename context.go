package goerrorkit

// HTTPContext là interface trừu tượng cho HTTP context
// Cho phép thư viện hoạt động với bất kỳ web framework nào
// Framework-specific adapters sẽ implement interface này
type HTTPContext interface {
	// Method trả về HTTP method (GET, POST, etc.)
	Method() string

	// Path trả về request path
	Path() string

	// GetLocal lấy giá trị từ context locals (như request ID)
	GetLocal(key string) interface{}

	// Status set HTTP status code cho response
	Status(code int) HTTPContext

	// JSON gửi JSON response
	JSON(data interface{}) error
}
