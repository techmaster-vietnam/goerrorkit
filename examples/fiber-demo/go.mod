module github.com/cuong/goerrorkit/examples/fiber-demo

go 1.21

require (
	github.com/cuong/goerrorkit v0.1.0
	github.com/gofiber/fiber/v2 v2.52.0
)

// For local development, use replace directive
replace github.com/cuong/goerrorkit => ../..

