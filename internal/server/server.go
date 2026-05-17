package server

import (
	"log/slog"
	"net/http"
	"os"
	"sync/atomic"
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

	// Instantiate metric counter
	appState := apiConfig{
		fileserverHits: atomic.Int32{},
	}

	// Register handlers
	serverHandler.Handle("/app/", appState.middlewareMetricsInc(indexFileHandler()))
	serverHandler.Handle("/healthz", healthHandler{})
	serverHandler.Handle("/metrics", &appState)
	serverHandler.Handle("/reset", appState.resetMetric())

	// Start server and log any failures
	slog.Info("Starting server", "address", "http://localhost"+string(srvr.Addr)+"/app")
	err := srvr.ListenAndServe()
	if err != nil {
		slog.Error("Server Error", "error", err)
		os.Exit(1)
	}
}
