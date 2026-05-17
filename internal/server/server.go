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

	// Register handlers
	serverHandler.Handle("/", indexFileHandler())
	serverHandler.Handle("/healthz", healthHandler{})

	// Start server and log any failures
	slog.Info("Starting server", "address", "http://localhost"+string(srvr.Addr))
	err := srvr.ListenAndServe()
	if err != nil {
		slog.Error("Server Error", "error", err)
		os.Exit(1)
	}
}
