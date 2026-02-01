package wireguard

import (
	"fmt"
	"log/slog"
)

// Peer represents a WireGuard peer.
type Peer struct {
	ID            string   `json:"id"`
	PublicKey     string   `json:"publicKey"`
	Name          string   `json:"name"`
	Endpoint      string   `json:"endpoint"`
	AllowedIPs    []string `json:"allowedIPs"`
	LastHandshake string   `json:"lastHandshake"`
	ReceiveBytes  int64    `json:"receiveBytes"`
	TransmitBytes int64    `json:"transmitBytes"`
}

// Stats represents interface-level statistics.
type Stats struct {
	InterfaceName string `json:"interfaceName"`
	PublicKey     string `json:"publicKey"`
	ListenPort    int    `json:"listenPort"`
	Subnet        string `json:"subnet"`
	PeerCount     int    `json:"peerCount"`
	TotalRX       int64  `json:"totalRx"`
	TotalTX       int64  `json:"totalTx"`
}

// PeerUpdate represents optional updates for a peer.
type PeerUpdate struct {
	Name       *string   `json:"name,omitempty"`
	AllowedIPs *[]string `json:"allowedIPs,omitempty"`
}

// PeerResponse represents a peer along with optional configuration details.
type PeerResponse struct {
	Peer
	Config     string `json:"config,omitempty"`
	PrivateKey string `json:"privateKey,omitempty"`
}

// Service defines the interface for WireGuard operations.
type Service interface {
	ListPeers() ([]Peer, error)
	AddPeer(name string, publicKey string, allowedIPs []string) (PeerResponse, error)
	RemovePeer(id string) error
	RegeneratePeer(id string) (PeerResponse, error)
	UpdatePeer(id string, updates PeerUpdate) (Peer, error)
	GetStats() (Stats, error)
	Close() error
}

// mockService is a mock implementation of the WireGuard service for development.
type mockService struct {
	peers []Peer
}

// NewMockService creates and returns a new mock WireGuard service.
func NewMockService() Service {
	return &mockService{
		peers: []Peer{
			{
				ID:            "mock-peer-1",
				PublicKey:     "ABC...",
				Name:          "Primary Server",
				Endpoint:      "192.168.1.1:51820",
				AllowedIPs:    []string{"10.0.0.2/32"},
				LastHandshake: "2026-01-31 02:00:00",
				ReceiveBytes:  1024,
				TransmitBytes: 2048,
			},
			{
				ID:            "mock-peer-2",
				PublicKey:     "XYZ...",
				Name:          "Mobile Client",
				Endpoint:      "192.168.1.2:51820",
				AllowedIPs:    []string{"10.0.0.3/32"},
				LastHandshake: "2026-01-31 02:05:00",
				ReceiveBytes:  512,
				TransmitBytes: 256,
			},
		},
	}
}

// ListPeers returns a list of mock WireGuard peers.
func (s *mockService) ListPeers() ([]Peer, error) {
	slog.Warn("Using mock WireGuard service for ListPeers")
	for _, p := range s.peers {
		if p.Name == "force-list-error" {
			return nil, fmt.Errorf("forced error")
		}
	}
	return s.peers, nil
}

// AddPeer adds a mock WireGuard peer.
func (s *mockService) AddPeer(name string, publicKey string, allowedIPs []string) (PeerResponse, error) {
	slog.Warn("Using mock WireGuard service for AddPeer")
	if name == "force-add-error" {
		return PeerResponse{}, fmt.Errorf("forced error")
	}
	peer := Peer{
		ID:         fmt.Sprintf("mock-peer-%d", len(s.peers)+1),
		PublicKey:  publicKey,
		Name:       name,
		AllowedIPs: allowedIPs,
	}
	s.peers = append(s.peers, peer)
	return PeerResponse{
		Peer:   peer,
		Config: "[Interface]\nPrivateKey = MOCK_KEY\n...",
	}, nil
}

// RemovePeer removes a mock WireGuard peer.
func (s *mockService) RemovePeer(id string) error {
	slog.Warn("Using mock WireGuard service for RemovePeer")
	if id == "force-error" {
		return fmt.Errorf("forced error")
	}
	for i, p := range s.peers {
		if p.ID == id {
			s.peers = append(s.peers[:i], s.peers[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("mock peer with ID %s not found", id)
}

// RegeneratePeer regenerates keys for a mock WireGuard peer.
func (s *mockService) RegeneratePeer(id string) (PeerResponse, error) {
	slog.Warn("Using mock WireGuard service for RegeneratePeer")
	if id == "force-error" {
		return PeerResponse{}, fmt.Errorf("forced error")
	}
	for i, p := range s.peers {
		if p.ID == id {
			// In a real implementation we'd generate new keys
			// For mock, just append "-new" to the public key to simulate change
			p.PublicKey = p.PublicKey + "-new"
			p.ID = p.PublicKey
			s.peers[i] = p
			return PeerResponse{
				Peer:   p,
				Config: "[Interface]\nPrivateKey = MOCK_REGENERATED_KEY\n...",
			}, nil
		}
	}
	return PeerResponse{}, fmt.Errorf("mock peer with ID %s not found", id)
}

// UpdatePeer updates a mock WireGuard peer.
func (s *mockService) UpdatePeer(id string, updates PeerUpdate) (Peer, error) {
	slog.Warn("Using mock WireGuard service for UpdatePeer")
	if id == "force-error" {
		return Peer{}, fmt.Errorf("forced error")
	}
	for i, p := range s.peers {
		if p.ID == id {
			if updates.Name != nil {
				if *updates.Name == "force-error" {
					return Peer{}, fmt.Errorf("forced error")
				}
				p.Name = *updates.Name
			}
			if updates.AllowedIPs != nil {
				p.AllowedIPs = *updates.AllowedIPs
			}
			s.peers[i] = p
			return p, nil
		}
	}
	return Peer{}, fmt.Errorf("mock peer with ID %s not found", id)
}

// GetStats returns mock interface-level statistics.
func (s *mockService) GetStats() (Stats, error) {
	slog.Warn("Using mock WireGuard service for GetStats")
	for _, p := range s.peers {
		if p.Name == "force-stats-error" {
			return Stats{}, fmt.Errorf("forced error")
		}
	}
	return Stats{
		InterfaceName: "mock-wg0",
		PublicKey:     "MOCK_SERVER_PUBKEY",
		ListenPort:    51820,
		Subnet:        "10.0.0.0/24",
		PeerCount:     len(s.peers),
		TotalRX:       1536,
		TotalTX:       2304,
	}, nil
}

// Close is a no-op for mockService.
func (s *mockService) Close() error {
	slog.Warn("Using mock WireGuard service for Close (no-op)")
	return nil
}
