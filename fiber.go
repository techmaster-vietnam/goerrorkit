package goerrorkit

import (
	fiberv2 "github.com/gofiber/fiber/v2"
)

// FiberContext wrap Fiber's context để implement HTTPContext interface
type FiberContext struct {
	ctx *fiberv2.Ctx
}

// NewFiberContext tạo FiberContext từ fiber.Ctx
func NewFiberContext(c *fiberv2.Ctx) *FiberContext {
	return &FiberContext{ctx: c}
}

// Method implements HTTPContext
func (f *FiberContext) Method() string {
	return f.ctx.Method()
}

// Path implements HTTPContext
func (f *FiberContext) Path() string {
	return f.ctx.Path()
}

// GetLocal implements HTTPContext
func (f *FiberContext) GetLocal(key string) interface{} {
	return f.ctx.Locals(key)
}

// Status implements HTTPContext
func (f *FiberContext) Status(code int) HTTPContext {
	f.ctx.Status(code)
	return f
}

// JSON implements HTTPContext
func (f *FiberContext) JSON(data interface{}) error {
	return f.ctx.JSON(data)
}

// FiberErrorHandler là Fiber middleware để xử lý panic và errors
// Tự động recover panic và convert errors sang AppError với stack trace chi tiết
//
// Example:
//
//	import "github.com/techmaster-vietnam/goerrorkit"
//
//	app := fiber.New()
//	app.Use(goerrorkit.FiberErrorHandler())
//
//	app.Get("/test", func(c *fiber.Ctx) error {
//	    // Panic sẽ được tự động catch và log với chính xác location
//	    panic("something went wrong")
//	})
func FiberErrorHandler() fiberv2.Handler {
	return func(c *fiberv2.Ctx) error {
		// Wrap Fiber context
		ctx := NewFiberContext(c)

		requestPath := ctx.Method() + " " + ctx.Path()
		requestID := "unknown"
		if rid, ok := ctx.GetLocal("requestid").(string); ok {
			requestID = rid
		}

		// Panic recovery với chính xác panic location
		defer func() {
			r := recover()
			if r != nil {
				// Xử lý panic bằng core logic - capture chính xác dòng gây panic
				panicErr := HandlePanic(r, requestID)
				LogAndRespond(ctx, panicErr, requestPath)
			}
		}()

		// Thực thi handler
		err := c.Next()

		// Xử lý error nếu có
		if err != nil {
			// Convert sang AppError bằng core logic
			appErr := ConvertToAppError(err, requestID)
			LogAndRespond(ctx, appErr, requestPath)
			return nil
		}

		return nil
	}
}

