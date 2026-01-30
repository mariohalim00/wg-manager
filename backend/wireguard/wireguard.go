package wireguard

import (
	"fmt"
	"log/slog"
)

// Peer represents a WireGuard peer.
type Peer struct {
	ID        string `json:"id"`
	PublicKey string `json:"publicKey"`
	Endpoint  string `json:"endpoint"`
	AllowedIPs []string `json:"allowedIPs"`
	// Add more fields as needed for a WireGuard peer
}

// Service defines the interface for WireGuard operations.
type Service interface {
	ListPeers() ([]Peer, error)
	AddPeer(peer Peer) (Peer, error)
	RemovePeer(id string) error
	// More WireGuard related operations can be added here
}

// mockService is a mock implementation of the WireGuard service for development.
type mockService struct{}

// NewMockService creates and returns a new mock WireGuard service.
func NewMockService() Service {
	return &mockService{}
}

// ListPeers returns a list of mock WireGuard peers.
func (s *mockService) ListPeers() ([]Peer, error) {
	slog.Warn("Using mock WireGuard service for ListPeers")
	return []Peer{
		{
			ID:        "mock-peer-1",
			PublicKey: "ABC...",
			Endpoint:  "192.168.1.1:51820",
			AllowedIPs: []string{"10.0.0.2/32"},
		},
		{
			ID:        "mock-peer-2",
			PublicKey: "XYZ...",
			Endpoint:  "192.168.1.2:51820",
			AllowedIPs: []string{"10.0.0.3/32"},
		},
	}, nil
}

// AddPeer adds a mock WireGuard peer.
func (s *mockService) AddPeer(peer Peer) (Peer, error) {
	slog.Warn("Using mock WireGuard service for AddPeer")
	// In a real implementation, you'd interact with WireGuard
	// For mock, we just assign an ID and add to our in-memory slice
	peer.ID = fmt.Sprintf("mock-peer-%d", len(mockPeers)+1) // Simple ID generation
	mockPeers = append(mockPeers, peer)
	return peer, nil
}

// RemovePeer removes a mock WireGuard peer.
func (s *mockService) RemovePeer(id string) error {
	slog.Warn("Using mock WireGuard service for RemovePeer")
	for i, p := range mockPeers {
		if p.ID == id {
			mockPeers = append(mockPeers[:i], mockPeers[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("mock peer with ID %s not found", id)
}

var mockPeers []Peer // Simple in-memory store for mock peers
