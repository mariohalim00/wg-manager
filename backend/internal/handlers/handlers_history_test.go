package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"wg-manager/backend/internal/wireguard"
)

func TestHistoryAndSettingsHandlers(t *testing.T) {
	mockWGService := wireguard.NewMockService()
	h := NewPeerHandler(mockWGService)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /stats/history", h.GetHistory)
	mux.HandleFunc("GET /settings", h.GetSettings)
	mux.HandleFunc("POST /settings", h.UpdateSettings)

	t.Run("GetHistory", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/stats/history", nil)
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rr.Code)
		}

		var history []wireguard.StatsHistoryItem
		if err := json.Unmarshal(rr.Body.Bytes(), &history); err != nil {
			t.Fatalf("unmarshal error: %v", err)
		}
		if len(history) == 0 {
			t.Error("expected history data, got empty")
		}
	})

	t.Run("GetSettings", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/settings", nil)
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rr.Code)
		}

		var settings wireguard.GlobalSettings
		if err := json.Unmarshal(rr.Body.Bytes(), &settings); err != nil {
			t.Fatalf("unmarshal error: %v", err)
		}
		if settings.MTU == 0 {
			t.Error("expected settings data, got empty/zero")
		}
	})

	t.Run("UpdateSettings", func(t *testing.T) {
		reqBody := `{"serverAddress":"10.0.0.1/24", "dns":"1.1.1.1", "mtu":1420, "keepalive":25, "endpoint":"vpn.example.com"}`
		req := httptest.NewRequest("POST", "/settings", strings.NewReader(reqBody))
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusNoContent {
			t.Errorf("expected 204, got %d", rr.Code)
		}
	})
}
