package main

import (
	"fmt"

	"calorie-counter/internal/config"
	"calorie-counter/internal/server"
)

// main is the process entrypoint. It only handles startup: load configuration
// into the package-wide global, then hand off to the server. It does not build
// controllers or services directly.
func main() {
	fmt.Println("Starting.")
	config.LoadConfig()
	server.Serve()
}
