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
	srvr := &http.Server{
		Handler: serverHandler,
		Addr:    serverAddress,
	}

	// Add file handler
	currentDir := http.Dir("./internal/server/")
	serverHandler.Handle("/", http.FileServer(currentDir))

	// Start server and log any failures
	slog.Info("Starting server", "address", serverAddress)
	err := srvr.ListenAndServe()
	if err != nil {
		slog.Error("Server Error", "error", err)
		os.Exit(1)
	}
}
