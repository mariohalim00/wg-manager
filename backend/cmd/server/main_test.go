package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
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

func TestAddPeerHandler(t *testing.T) {
	mockWGService := wireguard.NewMockService()
	h := handlers.NewPeerHandler(mockWGService)
	mux := http.NewServeMux()
	mux.HandleFunc("POST /peers", h.Add)

	t.Run("Success", func(t *testing.T) {
		reqBody := `{"name":"New Peer", "allowedIPs":["10.0.0.5/32"]}`
		req := httptest.NewRequest("POST", "/peers", strings.NewReader(reqBody))
		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusCreated)
		}

		var resp wireguard.PeerResponse
		if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
			t.Fatalf("could not unmarshal response: %v", err)
		}

		if resp.Name != "New Peer" {
			t.Errorf("expected name 'New Peer', got '%s'", resp.Name)
		}
	})

	t.Run("MissingName", func(t *testing.T) {
		reqBody := `{"name":"", "allowedIPs":["10.0.0.5/32"]}`
		req := httptest.NewRequest("POST", "/peers", strings.NewReader(reqBody))
		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusBadRequest)
		}
	})

	t.Run("InvalidCIDR", func(t *testing.T) {
		reqBody := `{"name":"Test", "allowedIPs":["invalid-cidr"]}`
		req := httptest.NewRequest("POST", "/peers", strings.NewReader(reqBody))
		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusBadRequest)
		}
	})
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

	if stats.PublicKey != "MOCK_SERVER_PUBKEY" {
		t.Errorf("expected PublicKey 'MOCK_SERVER_PUBKEY', got '%s'", stats.PublicKey)
	}

	if stats.ListenPort != 51820 {
		t.Errorf("expected ListenPort 51820, got %d", stats.ListenPort)
	}

	if stats.Subnet != "10.0.0.0/24" {
		t.Errorf("expected Subnet '10.0.0.0/24', got '%s'", stats.Subnet)
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
