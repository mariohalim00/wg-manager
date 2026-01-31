package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"wg-manager/backend/internal/wireguard"
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
	Name       string   `json:"name"`
	PublicKey  string   `json:"publicKey"`
	AllowedIPs []string `json:"allowedIPs"`
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

	peer, err := h.Service.AddPeer(req.Name, req.PublicKey, req.AllowedIPs)
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
