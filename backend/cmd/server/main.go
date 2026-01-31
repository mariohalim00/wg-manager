package main

import (
	"log/slog"
	"net/http"
	"os"

	"wg-manager/backend/internal/config" // Import the config package
	"wg-manager/backend/internal/handlers"
	"wg-manager/backend/internal/middleware"
	"wg-manager/backend/internal/wireguard"
)

// Application holds application-wide dependencies.
type Application struct {
	Config    *config.Config
	WireGuard wireguard.Service
}

func main() {
	// Set default logger to JSON mode
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Load configuration
	cfg, err := config.LoadConfig("internal/config/config.json")
	if err != nil {
		slog.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	// Initialize application dependencies
	var wgService wireguard.Service
	wgService, err = wireguard.NewRealService(
		cfg.InterfaceName,
		cfg.StoragePath,
		cfg.ServerEndpoint,
		cfg.ServerPubKey,
	)
	if err != nil {
		slog.Warn("Failed to initialize native WireGuard service, falling back to mock", "error", err)
		wgService = wireguard.NewMockService()
	}

	app := &Application{
		Config:    cfg,
		WireGuard: wgService,
	}

	peerHandler := handlers.NewPeerHandler(app.WireGuard)

	// Create a new ServeMux and register routes using modern syntax
	mux := http.NewServeMux()
	mux.HandleFunc("GET /peers", peerHandler.List)
	mux.HandleFunc("POST /peers", peerHandler.Add)
	mux.HandleFunc("DELETE /peers/{id}", peerHandler.Remove)
	mux.HandleFunc("GET /stats", peerHandler.Stats)

	// Apply middleware to all routes
	wrappedMux := middleware.LoggingMiddleware(mux)
	wrappedMux = middleware.CORSMiddleware(wrappedMux)

	slog.Info("Server starting", "port", app.Config.ServerPort)
	if err := http.ListenAndServe(app.Config.ServerPort, wrappedMux); err != nil {
		slog.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}
