module github.com/techmaster-vietnam/goerrorkit/examples/fiber-demo

go 1.21

require (
	github.com/techmaster-vietnam/goerrorkit v0.1.0
	github.com/gofiber/fiber/v2 v2.52.9
)

// For local development, use replace directive
replace github.com/techmaster-vietnam/goerrorkit => ../..

