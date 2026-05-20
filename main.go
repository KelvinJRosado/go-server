package main

import (
	"github.com/KelvinJRosado/go-server/internal/server"
	"github.com/joho/godotenv"
)

func main() {
	// Load env vars
	godotenv.Load()

	// Run server
	server.Run()
}
