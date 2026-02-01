package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"wg-manager/backend/internal/wireguard"

	"github.com/skip2/go-qrcode"
)

type PeerHandler struct {
	Service wireguard.Service
}

func NewPeerHandler(service wireguard.Service) *PeerHandler {
	return &PeerHandler{Service: service}
}

func (h *PeerHandler) List(w http.ResponseWriter, r *http.Request) {
	peers, err := h.Service.ListPeers()
	if err != nil {
		slog.Error("Failed to list peers", "error", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(peers); err != nil {
		slog.Error("Failed to encode peers response", "error", err)
	}
}

type AddPeerRequest struct {
	Name                string   `json:"name"`
	PublicKey           string   `json:"publicKey"`
	AllowedIPs          []string `json:"allowedIPs"`
	DNS                 string   `json:"dns"`
	MTU                 int      `json:"mtu"`
	PersistentKeepalive int      `json:"persistentKeepalive"`
	PreSharedKey        bool     `json:"preSharedKey"`
	InterfaceAddress    string   `json:"interfaceAddress"`
}

func (h *PeerHandler) Add(w http.ResponseWriter, r *http.Request) {
	var req AddPeerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("Failed to decode add peer request", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Input Validation
	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	if len(req.AllowedIPs) == 0 {
		http.Error(w, "At least one AllowedIP is required", http.StatusBadRequest)
		return
	}

	for _, ip := range req.AllowedIPs {
		if _, _, err := net.ParseCIDR(ip); err != nil {
			http.Error(w, fmt.Sprintf("Invalid AllowedIP CIDR: %s", ip), http.StatusBadRequest)
			return
		}
	}

	opts := wireguard.AddPeerOptions{
		Name:                req.Name,
		PublicKey:           req.PublicKey,
		AllowedIPs:          req.AllowedIPs,
		DNS:                 req.DNS,
		MTU:                 req.MTU,
		PersistentKeepalive: req.PersistentKeepalive,
		PreSharedKey:        req.PreSharedKey,
		InterfaceAddress:    req.InterfaceAddress,
	}

	peer, err := h.Service.AddPeer(opts)
	if err != nil {
		slog.Error("Failed to add peer", "error", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(peer); err != nil {
		slog.Error("Failed to encode peer response", "error", err)
	}
}

func (h *PeerHandler) Remove(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	// id here is the public key of the peer
	if id == "" {
		http.Error(w, "Missing peer ID in path", http.StatusBadRequest)
		return
	}

	if err := h.Service.RemovePeer(id); err != nil {
		slog.Error("Failed to remove peer", "error", err, "id", id)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *PeerHandler) Regenerate(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Missing peer ID in path", http.StatusBadRequest)
		return
	}

	peer, err := h.Service.RegeneratePeer(id)
	if err != nil {
		slog.Error("Failed to regenerate peer keys", "error", err, "id", id)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(peer); err != nil {
		slog.Error("Failed to encode peer response", "error", err)
	}
}

type UpdatePeerRequest struct {
	Name                *string   `json:"name"`
	AllowedIPs          *[]string `json:"allowedIPs"`
	DNS                 *string   `json:"dns"`
	MTU                 *int      `json:"mtu"`
	PersistentKeepalive *int      `json:"persistentKeepalive"`
	InterfaceAddress    *string   `json:"interfaceAddress"`
}

func (h *PeerHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Missing peer ID in path", http.StatusBadRequest)
		return
	}

	var req UpdatePeerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("Failed to decode update peer request", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validation
	if req.AllowedIPs != nil {
		for _, ip := range *req.AllowedIPs {
			if _, _, err := net.ParseCIDR(ip); err != nil {
				http.Error(w, fmt.Sprintf("Invalid AllowedIP CIDR: %s", ip), http.StatusBadRequest)
				return
			}
		}
	}

	updates := wireguard.PeerUpdate{
		Name:                req.Name,
		AllowedIPs:          req.AllowedIPs,
		DNS:                 req.DNS,
		MTU:                 req.MTU,
		PersistentKeepalive: req.PersistentKeepalive,
		InterfaceAddress:    req.InterfaceAddress,
	}

	peer, err := h.Service.UpdatePeer(id, updates)
	if err != nil {
		slog.Error("Failed to update peer", "error", err, "id", id)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(peer); err != nil {
		slog.Error("Failed to encode peer response", "error", err)
	}
}

func (h *PeerHandler) Stats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.Service.GetStats()
	if err != nil {
		slog.Error("Failed to get stats", "error", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(stats); err != nil {
		slog.Error("Failed to encode stats response", "error", err)
	}
}

func (h *PeerHandler) GetConfig(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Missing peer ID in path", http.StatusBadRequest)
		return
	}

	config, err := h.Service.GetPeerConfig(id)
	if err != nil {
		slog.Error("Failed to get peer config", "error", err, "id", id)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.conf\"", id))
	w.Write([]byte(config))
}

func (h *PeerHandler) GetQR(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Missing peer ID in path", http.StatusBadRequest)
		return
	}

	config, err := h.Service.GetPeerConfig(id)
	if err != nil {
		slog.Error("Failed to get peer config for QR", "error", err, "id", id)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	png, err := qrcode.Encode(config, qrcode.High, 256)
	if err != nil {
		slog.Error("Failed to generate QR code", "error", err)
		http.Error(w, "Failed to generate QR code", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.Write(png)
}

func (h *PeerHandler) GetHistory(w http.ResponseWriter, r *http.Request) {
	history, err := h.Service.GetStatsHistory()
	if err != nil {
		slog.Error("Failed to get stats history", "error", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(history); err != nil {
		slog.Error("Failed to encode history response", "error", err)
	}
}

func (h *PeerHandler) GetSettings(w http.ResponseWriter, r *http.Request) {
	settings, err := h.Service.GetSettings()
	if err != nil {
		slog.Error("Failed to get settings", "error", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(settings); err != nil {
		slog.Error("Failed to encode settings response", "error", err)
	}
}

func (h *PeerHandler) UpdateSettings(w http.ResponseWriter, r *http.Request) {
	var settings wireguard.GlobalSettings
	if err := json.NewDecoder(r.Body).Decode(&settings); err != nil {
		slog.Error("Failed to decode update settings request", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.Service.UpdateSettings(settings); err != nil {
		slog.Error("Failed to update settings", "error", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
