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
	serverHandler.Handle("/app/", appState.middlewareMetricsInc(appState.indexFileHandler()))
	serverHandler.Handle("/healthz", appState.healthHandler())
	serverHandler.Handle("/metrics", appState.metricsHandler())
	serverHandler.Handle("/reset", appState.resetMetricsHandler())

	// Start server and log any failures
	slog.Info("Starting server", "address", "http://localhost"+string(srvr.Addr)+"/app")
	err := srvr.ListenAndServe()
	if err != nil {
		slog.Error("Server Error", "error", err)
		os.Exit(1)
	}
}
