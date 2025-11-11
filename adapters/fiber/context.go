package fiber

import (
	"github.com/cuong/goerrorkit/core"
	fiberv2 "github.com/gofiber/fiber/v2"
)

// FiberContext wrap Fiber's context để implement core.HTTPContext interface
type FiberContext struct {
	ctx *fiberv2.Ctx
}

// NewFiberContext tạo FiberContext từ fiber.Ctx
func NewFiberContext(c *fiberv2.Ctx) *FiberContext {
	return &FiberContext{ctx: c}
}

// Method implements core.HTTPContext
func (f *FiberContext) Method() string {
	return f.ctx.Method()
}

// Path implements core.HTTPContext
func (f *FiberContext) Path() string {
	return f.ctx.Path()
}

// GetLocal implements core.HTTPContext
func (f *FiberContext) GetLocal(key string) interface{} {
	return f.ctx.Locals(key)
}

// Status implements core.HTTPContext
func (f *FiberContext) Status(code int) core.HTTPContext {
	f.ctx.Status(code)
	return f
}

// JSON implements core.HTTPContext
func (f *FiberContext) JSON(data interface{}) error {
	return f.ctx.JSON(data)
}
