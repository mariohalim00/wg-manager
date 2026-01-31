package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"wg-manager/backend/internal/config"
	"wg-manager/backend/internal/wireguard"
)

func TestPeersHandler(t *testing.T) {
	// Initialize a mock WireGuard service for testing
	mockWGService := wireguard.NewMockService()
	app := &Application{WireGuard: mockWGService}

	req, err := http.NewRequest("GET", "/peers", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.peersHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var peers []wireguard.Peer
	err = json.Unmarshal(rr.Body.Bytes(), &peers)
	if err != nil {
		t.Fatalf("could not unmarshal response: %v", err)
	}

	if len(peers) != 2 {
		t.Errorf("handler returned unexpected number of peers: got %d want %d",
			len(peers), 2)
	}

	if contentType := rr.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("handler returned wrong Content-Type: got %v want %v",
			contentType, "application/json")
	}
}

func TestStatsHandler(t *testing.T) {
	mockWGService := wireguard.NewMockService()
	app := &Application{WireGuard: mockWGService}

	req, err := http.NewRequest("GET", "/stats", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.statsHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var stats wireguard.Stats
	err = json.Unmarshal(rr.Body.Bytes(), &stats)
	if err != nil {
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
