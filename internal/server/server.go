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

	// Instantiate metric counter
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
	}

	// Create server
	srvr := &http.Server{
		Handler: apiCfg.middlewareRequestLogging(serverHandler),
		Addr:    serverAddress,
	}

	// Register handlers
	serverHandler.Handle("/app/", apiCfg.middlewareMetricsInc(apiCfg.indexFileHandler()))
	serverHandler.Handle("GET /healthz", apiCfg.healthHandler())
	serverHandler.Handle("GET /metrics", apiCfg.metricsHandler())
	serverHandler.Handle("POST /reset", apiCfg.resetMetricsHandler())

	// Start server and log any failures
	slog.Info("Starting server", "address", "http://localhost"+string(srvr.Addr)+"/app")
	err := srvr.ListenAndServe()
	if err != nil {
		slog.Error("Server Error", "error", err)
		os.Exit(1)
	}
}
