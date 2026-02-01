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

	t.Run("NonexistentPeer", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/peers/nonexistent", nil)
		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusInternalServerError)
		}
	})
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

	t.Run("InvalidJSON", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/peers", strings.NewReader(`{invalid json}`))
		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusBadRequest)
		}
	})
}

func TestRegeneratePeerHandler(t *testing.T) {
	mockWGService := wireguard.NewMockService()
	h := handlers.NewPeerHandler(mockWGService)
	mux := http.NewServeMux()
	mux.HandleFunc("POST /peers/{id}/regenerate-keys", h.Regenerate)

	req := httptest.NewRequest("POST", "/peers/mock-peer-1/regenerate-keys", nil)
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var resp wireguard.PeerResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("could not unmarshal response: %v", err)
	}

	if !strings.HasSuffix(resp.PublicKey, "-new") {
		t.Errorf("expected public key to be regenerated (suffix -new), got %s", resp.PublicKey)
	}

	t.Run("NonexistentPeer", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/peers/nonexistent/regenerate-keys", nil)
		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusInternalServerError)
		}
	})
}

func TestUpdatePeerHandler(t *testing.T) {
	mockWGService := wireguard.NewMockService()
	h := handlers.NewPeerHandler(mockWGService)
	mux := http.NewServeMux()
	mux.HandleFunc("PATCH /peers/{id}", h.Update)

	t.Run("UpdateName", func(t *testing.T) {
		reqBody := `{"name":"Updated Name"}`
		req := httptest.NewRequest("PATCH", "/peers/mock-peer-1", strings.NewReader(reqBody))
		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		var resp wireguard.Peer
		if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
			t.Fatalf("could not unmarshal response: %v", err)
		}

		if resp.Name != "Updated Name" {
			t.Errorf("expected name 'Updated Name', got '%s'", resp.Name)
		}
	})

	t.Run("UpdateAllowedIPs", func(t *testing.T) {
		reqBody := `{"allowedIPs":["10.0.0.10/32"]}`
		req := httptest.NewRequest("PATCH", "/peers/mock-peer-1", strings.NewReader(reqBody))
		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		var resp wireguard.Peer
		if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
			t.Fatalf("could not unmarshal response: %v", err)
		}

		if len(resp.AllowedIPs) != 1 || resp.AllowedIPs[0] != "10.0.0.10/32" {
			t.Errorf("expected AllowedIPs ['10.0.0.10/32'], got %v", resp.AllowedIPs)
		}
	})
	t.Run("InvalidCIDR", func(t *testing.T) {
		reqBody := `{"allowedIPs":["invalid"]}`
		req := httptest.NewRequest("PATCH", "/peers/mock-peer-1", strings.NewReader(reqBody))
		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusBadRequest)
		}
	})

	t.Run("NonexistentPeer", func(t *testing.T) {
		reqBody := `{"name":"New Name"}`
		req := httptest.NewRequest("PATCH", "/peers/nonexistent", strings.NewReader(reqBody))
		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusInternalServerError)
		}
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		req := httptest.NewRequest("PATCH", "/peers/mock-peer-1", strings.NewReader(`{invalid json}`))
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

	t.Run("ServiceError", func(t *testing.T) {
		mockWGService := wireguard.NewMockService()
		_, _ = mockWGService.AddPeer(wireguard.AddPeerOptions{Name: "force-stats-error", AllowedIPs: []string{"10.0.0.1/32"}})

		h := handlers.NewPeerHandler(mockWGService)
		mux := http.NewServeMux()
		mux.HandleFunc("GET /stats", h.Stats)

		req := httptest.NewRequest("GET", "/stats", nil)
		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusInternalServerError)
		}
	})
}

func TestListPeersHandlerError(t *testing.T) {
	mockWGService := wireguard.NewMockService()
	_, _ = mockWGService.AddPeer(wireguard.AddPeerOptions{Name: "force-list-error", AllowedIPs: []string{"10.0.0.1/32"}})

	h := handlers.NewPeerHandler(mockWGService)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /peers", h.List)

	req := httptest.NewRequest("GET", "/peers", nil)
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
}

func TestAddPeerHandlerServiceError(t *testing.T) {
	mockWGService := wireguard.NewMockService()
	h := handlers.NewPeerHandler(mockWGService)
	mux := http.NewServeMux()
	mux.HandleFunc("POST /peers", h.Add)

	reqBody := `{"name":"force-add-error", "allowedIPs":["10.0.0.5/32"]}`
	req := httptest.NewRequest("POST", "/peers", strings.NewReader(reqBody))
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
}

func TestRemovePeerHandlerForceError(t *testing.T) {
	mockWGService := wireguard.NewMockService()
	h := handlers.NewPeerHandler(mockWGService)
	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /peers/{id}", h.Remove)

	req := httptest.NewRequest("DELETE", "/peers/force-error", nil)
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
}

func TestRegeneratePeerHandlerForceError(t *testing.T) {
	mockWGService := wireguard.NewMockService()
	h := handlers.NewPeerHandler(mockWGService)
	mux := http.NewServeMux()
	mux.HandleFunc("POST /peers/{id}/regenerate-keys", h.Regenerate)

	req := httptest.NewRequest("POST", "/peers/force-error/regenerate-keys", nil)
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
}

func TestUpdatePeerHandlerForceError(t *testing.T) {
	mockWGService := wireguard.NewMockService()
	h := handlers.NewPeerHandler(mockWGService)
	mux := http.NewServeMux()
	mux.HandleFunc("PATCH /peers/{id}", h.Update)

	t.Run("UpdateIDForceError", func(t *testing.T) {
		reqBody := `{"name":"New Name"}`
		req := httptest.NewRequest("PATCH", "/peers/force-error", strings.NewReader(reqBody))
		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusInternalServerError)
		}
	})

	t.Run("UpdateNameForceError", func(t *testing.T) {
		// First add a normal peer
		_, _ = mockWGService.AddPeer(wireguard.AddPeerOptions{Name: "normal", AllowedIPs: []string{"10.0.0.1/32"}})

		reqBody := `{"name":"force-error"}`
		req := httptest.NewRequest("PATCH", "/peers/mock-peer-1", strings.NewReader(reqBody))
		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusInternalServerError)
		}
	})
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

func TestHandlerMissingID(t *testing.T) {
	mockWGService := wireguard.NewMockService()
	h := handlers.NewPeerHandler(mockWGService)
	mux := http.NewServeMux()

	// Register without {id} to test missing ID check in handler
	mux.HandleFunc("DELETE /peers/", h.Remove)
	mux.HandleFunc("POST /peers/regenerate-keys", h.Regenerate)
	mux.HandleFunc("PATCH /peers/", h.Update)

	t.Run("Remove", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/peers/", nil)
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", rr.Code)
		}
	})

	t.Run("Regenerate", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/peers/regenerate-keys", nil)
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", rr.Code)
		}
	})

	t.Run("Update", func(t *testing.T) {
		reqBody := `{"name":"test"}`
		req := httptest.NewRequest("PATCH", "/peers/", strings.NewReader(reqBody))
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", rr.Code)
		}
	})
}

func TestAddPeerHandlerEmptyAllowedIPs(t *testing.T) {
	mockWGService := wireguard.NewMockService()
	h := handlers.NewPeerHandler(mockWGService)
	mux := http.NewServeMux()
	mux.HandleFunc("POST /peers", h.Add)

	reqBody := `{"name":"Test", "allowedIPs":[]}`
	req := httptest.NewRequest("POST", "/peers", strings.NewReader(reqBody))
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}
