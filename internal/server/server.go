package server

import (
	"log/slog"
	"net/http"
	"os"
)

func Run() {
	// Create inputs for server
	serverHandler := http.NewServeMux()
	serverAddress := ":8080"

	// Create server
	server := http.Server{
		Handler: serverHandler,
		Addr:    serverAddress,
	}

	// Start server and log any failures
	slog.Info("Starting server", "address", serverAddress)
	err := server.ListenAndServe()
	if err != nil {
		slog.Error("Server Error", "error", err)
		os.Exit(1)
	}
}
