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
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
	}

	// Register handlers
	serverHandler.Handle("/app/", apiCfg.middlewareMetricsInc(apiCfg.indexFileHandler()))
	serverHandler.Handle("/healthz", apiCfg.healthHandler())
	serverHandler.Handle("/metrics", apiCfg.metricsHandler())
	serverHandler.Handle("/reset", apiCfg.resetMetricsHandler())

	// Start server and log any failures
	slog.Info("Starting server", "address", "http://localhost"+string(srvr.Addr)+"/app")
	err := srvr.ListenAndServe()
	if err != nil {
		slog.Error("Server Error", "error", err)
		os.Exit(1)
	}
}
