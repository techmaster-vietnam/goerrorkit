package fiber

import (
	fiberv2 "github.com/gofiber/fiber/v2"
	"github.com/techmaster-vietnam/goerrorkit"
)

// FiberContext wrap Fiber's context để implement goerrorkit.HTTPContext interface
type FiberContext struct {
	ctx *fiberv2.Ctx
}

// NewFiberContext tạo FiberContext từ fiber.Ctx
func NewFiberContext(c *fiberv2.Ctx) *FiberContext {
	return &FiberContext{ctx: c}
}

// Method implements goerrorkit.HTTPContext
func (f *FiberContext) Method() string {
	return f.ctx.Method()
}

// Path implements goerrorkit.HTTPContext
func (f *FiberContext) Path() string {
	return f.ctx.Path()
}

// GetLocal implements goerrorkit.HTTPContext
func (f *FiberContext) GetLocal(key string) interface{} {
	return f.ctx.Locals(key)
}

// Status implements goerrorkit.HTTPContext
func (f *FiberContext) Status(code int) goerrorkit.HTTPContext {
	f.ctx.Status(code)
	return f
}

// JSON implements goerrorkit.HTTPContext
func (f *FiberContext) JSON(data interface{}) error {
	return f.ctx.JSON(data)
}
