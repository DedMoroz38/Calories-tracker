package main

import (
	"fmt"
	"log"

	"calorie-counter/internal/config"
	"calorie-counter/internal/server"
	"calorie-counter/internal/storage"
)

// main is the process entrypoint. It only handles startup: load configuration
// into the package-wide global, initialise the S3 storage handle, then hand off
// to the server. It does not build controllers or services directly.
func main() {
	fmt.Println("Starting.")
	config.LoadConfig()
	if err := storage.Init(); err != nil {
		log.Fatal("failed to initialise S3 storage: ", err)
	}
	if !storage.IsEnabled() {
		log.Println("WARNING: S3 not configured (AWS_* env vars missing) — photo endpoints will return errors")
	}
	server.Serve()
}
