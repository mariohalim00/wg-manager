package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"

	"wg-manager/backend/internal/config" // Import the config package
	"wg-manager/backend/internal/middleware"
	"wg-manager/backend/internal/wireguard"
)

// Application holds application-wide dependencies.
type Application struct {
	Config    *config.Config
	WireGuard wireguard.Service
}

// peersHandler handles GET, POST, and DELETE requests for peers.
func (app *Application) peersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		app.listPeers(w, r)
	case http.MethodPost:
		app.addPeer(w, r)
	case http.MethodDelete:
		app.removePeer(w, r)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (app *Application) listPeers(w http.ResponseWriter, r *http.Request) {
	peers, err := app.WireGuard.ListPeers()
	if err != nil {
		slog.Error("Failed to list peers", "error", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(peers)
}

type AddPeerRequest struct {
	Name       string   `json:"name"`
	PublicKey  string   `json:"publicKey"`
	AllowedIPs []string `json:"allowedIPs"`
}

func (app *Application) addPeer(w http.ResponseWriter, r *http.Request) {
	var req AddPeerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("Failed to decode add peer request", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	peer, err := app.WireGuard.AddPeer(req.Name, req.PublicKey, req.AllowedIPs)
	if err != nil {
		slog.Error("Failed to add peer", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(peer)
}

func (app *Application) removePeer(w http.ResponseWriter, r *http.Request) {
	publicKey := r.URL.Query().Get("publicKey")
	if publicKey == "" {
		http.Error(w, "Missing publicKey query parameter", http.StatusBadRequest)
		return
	}

	if err := app.WireGuard.RemovePeer(publicKey); err != nil {
		slog.Error("Failed to remove peer", "error", err, "publicKey", publicKey)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// statsHandler returns interface statistics using the WireGuard service.
func (app *Application) statsHandler(w http.ResponseWriter, r *http.Request) {
	stats, err := app.WireGuard.GetStats()
	if err != nil {
		slog.Error("Failed to get stats", "error", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
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

	// Create a new ServeMux to apply middleware
	mux := http.NewServeMux()
	mux.HandleFunc("/peers", app.peersHandler)
	mux.HandleFunc("/stats", app.statsHandler)

	// Apply middleware to all routes
	wrappedMux := middleware.LoggingMiddleware(mux)
	wrappedMux = middleware.CORSMiddleware(wrappedMux)

	slog.Info("Server starting", "port", app.Config.ServerPort)
	if err := http.ListenAndServe(app.Config.ServerPort, wrappedMux); err != nil {
		slog.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}
