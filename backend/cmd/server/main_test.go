package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"wg-manager/backend/internal/config"
	"wg-manager/backend/internal/handlers"
	"wg-manager/backend/internal/wireguard"
)

func TestPeersHandler(t *testing.T) {
	mockWGService := wireguard.NewMockService()
	h := handlers.NewPeerHandler(mockWGService)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /peers", h.List)

	req := httptest.NewRequest("GET", "/peers", nil)
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var peers []wireguard.Peer
	if err := json.Unmarshal(rr.Body.Bytes(), &peers); err != nil {
		t.Fatalf("could not unmarshal response: %v", err)
	}

	if len(peers) != 2 {
		t.Errorf("handler returned unexpected number of peers: got %d want %d",
			len(peers), 2)
	}
}

func TestRemovePeerHandler(t *testing.T) {
	mockWGService := wireguard.NewMockService()
	h := handlers.NewPeerHandler(mockWGService)
	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /peers/{id}", h.Remove)

	// Test successful removal
	req := httptest.NewRequest("DELETE", "/peers/mock-peer-1", nil)
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNoContent)
	}
}

func TestStatsHandler(t *testing.T) {
	mockWGService := wireguard.NewMockService()
	h := handlers.NewPeerHandler(mockWGService)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /stats", h.Stats)

	req := httptest.NewRequest("GET", "/stats", nil)
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var stats wireguard.Stats
	if err := json.Unmarshal(rr.Body.Bytes(), &stats); err != nil {
		t.Fatalf("could not unmarshal response: %v", err)
	}

	if stats.PeerCount != 2 {
		t.Errorf("handler returned unexpected peer count: got %d want %d",
			stats.PeerCount, 2)
	}
}

func TestLoadConfig(t *testing.T) {
	cfg, err := config.LoadConfig("../../internal/config/config.json")
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if cfg.ServerPort != ":8080" {
		t.Errorf("Expected ServerPort to be ':8080', got %s", cfg.ServerPort)
	}
}
