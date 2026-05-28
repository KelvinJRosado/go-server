package server

import (
	"database/sql"
	"log/slog"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/KelvinJRosado/go-server/internal/database"
	_ "github.com/lib/pq"
)

func Run() {
	// Get DB Connection
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		slog.Error("Failed to connect to DB", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	dbQueries := database.New(db)

	// Get JWT secret
	jwtSecret := os.Getenv("JWT_SECRET")

	// Create inputs for server
	serverHandler := http.NewServeMux()
	serverAddress := ":8080"

	// Instantiate metric counter
	apiCfg := apiConfig{
		fileserverHits:  atomic.Int32{},
		databaseQueries: dbQueries,
		platform:        os.Getenv("PLATFORM"),
		jwtSecret:       jwtSecret,
	}

	// Create server
	srvr := &http.Server{
		Handler: apiCfg.middlewareRequestLogging(serverHandler), // Apply request logging middleware to all handlers
		Addr:    serverAddress,
	}

	// Register handlers
	serverHandler.Handle("/app/", apiCfg.middlewareMetricsInc(apiCfg.indexFileHandler()))

	serverHandler.Handle("GET /admin/metrics", apiCfg.metricsHandler())
	serverHandler.Handle("POST /admin/reset", apiCfg.adminResetHandler())

	serverHandler.Handle("POST /api/users", apiCfg.createUserHandler())
	serverHandler.Handle("PUT /api/users", apiCfg.updateUserHandler())

	serverHandler.Handle("POST /api/login", apiCfg.loginHandler())
	serverHandler.Handle("POST /api/refresh", apiCfg.refreshHandler())
	serverHandler.Handle("POST /api/revoke", apiCfg.revokeHandler())

	serverHandler.Handle("GET /api/healthz", apiCfg.healthHandler())

	serverHandler.Handle("GET /api/chirps", apiCfg.getChirpsHandler())
	serverHandler.Handle("GET /api/chirps/{id}", apiCfg.getChirpByIdHandler())
	serverHandler.Handle("DELETE /api/chirps/{id}", apiCfg.deleteChirpByIdHandler())
	serverHandler.Handle("POST /api/chirps", apiCfg.createChirpHandler())

	serverHandler.Handle("POST /api/polka/webhooks", apiCfg.upgradeUserToRedHandler())

	// Start server and log any failures
	slog.Info("Starting server", "address", "http://localhost"+string(srvr.Addr)+"/app")
	err = srvr.ListenAndServe()
	if err != nil {
		slog.Error("Server Error", "error", err)
		os.Exit(1)
	}
}
