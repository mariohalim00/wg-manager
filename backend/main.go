package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"

	"wg-manager/backend/config"        // Import the config package
	"wg-manager/backend/middleware"
	"wg-manager/backend/wireguard"
)

// Application holds application-wide dependencies.
type Application struct {
	Config    *config.Config
	WireGuard wireguard.Service
}

// peersHandler returns a list of peers using the WireGuard service.
func (app *Application) peersHandler(w http.ResponseWriter, r *http.Request) {
	peers, err := app.WireGuard.ListPeers()
	if err != nil {
		slog.Error("Failed to list peers", "error", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(peers)
}

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("./config/config.json")
	if err != nil {
		slog.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	// Initialize application dependencies
	app := &Application{
		Config:    cfg,
		WireGuard: wireguard.NewMockService(), // Use the mock service for now
	}

	// Create a new ServeMux to apply middleware
	mux := http.NewServeMux()
	mux.HandleFunc("/peers", app.peersHandler)

	// Apply logging middleware to all routes
	wrappedMux := middleware.LoggingMiddleware(mux)

	slog.Info("Server starting", "port", app.Config.ServerPort)
	if err := http.ListenAndServe(app.Config.ServerPort, wrappedMux); err != nil {
		slog.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}