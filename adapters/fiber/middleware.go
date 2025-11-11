package fiber

import (
	fiberv2 "github.com/gofiber/fiber/v2"
	"github.com/techmaster-vietnam/goerrorkit/core"
)

// ErrorHandler là Fiber middleware để xử lý panic và errors
// Tự động recover panic và convert errors sang AppError với stack trace chi tiết
//
// Example:
//
//	app := fiber.New()
//	app.Use(fiber.ErrorHandler())
//
//	app.Get("/test", func(c *fiber.Ctx) error {
//	    // Panic sẽ được tự động catch và log với chính xác location
//	    panic("something went wrong")
//	})
func ErrorHandler() fiberv2.Handler {
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
				panicErr := core.HandlePanic(r, requestID)
				core.LogAndRespond(ctx, panicErr, requestPath)
			}
		}()

		// Thực thi handler
		err := c.Next()

		// Xử lý error nếu có
		if err != nil {
			// Convert sang AppError bằng core logic
			appErr := core.ConvertToAppError(err, requestID)
			core.LogAndRespond(ctx, appErr, requestPath)
			return nil
		}

		return nil
	}
}
