package wireguard

import (
	"fmt"
	"log/slog"
)

// Peer represents a WireGuard peer.
type Peer struct {
	ID               string   `json:"id"`
	PublicKey        string   `json:"publicKey"`
	Name             string   `json:"name"`
	Endpoint         string   `json:"endpoint"`
	AllowedIPs       []string `json:"allowedIPs"`
	LastHandshake    string   `json:"lastHandshake"`
	ReceiveBytes     int64    `json:"receiveBytes"`
	TransmitBytes    int64    `json:"transmitBytes"`
	InterfaceAddress string   `json:"interfaceAddress,omitempty"`
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
	Name                *string   `json:"name,omitempty"`
	AllowedIPs          *[]string `json:"allowedIPs,omitempty"`
	DNS                 *string   `json:"dns,omitempty"`
	MTU                 *int      `json:"mtu,omitempty"`
	PersistentKeepalive *int      `json:"persistentKeepalive,omitempty"`
	InterfaceAddress    *string   `json:"interfaceAddress,omitempty"`
}

// StatsHistoryItem represents a single data point in traffic history.
type StatsHistoryItem struct {
	Timestamp int64 `json:"timestamp"`
	TotalRX   int64 `json:"totalRx"`
	TotalTX   int64 `json:"totalTx"`
}

// PeerResponse represents a peer along with optional configuration details.
type PeerResponse struct {
	Peer
	Config       string `json:"config,omitempty"`
	PrivateKey   string `json:"privateKey,omitempty"`
	PresharedKey string `json:"presharedKey,omitempty"`
}

// AddPeerOptions represents the options for creating a new peer.
type AddPeerOptions struct {
	Name                string   `json:"name"`
	PublicKey           string   `json:"publicKey,omitempty"`
	AllowedIPs          []string `json:"allowedIPs"`
	DNS                 string   `json:"dns,omitempty"`
	MTU                 int      `json:"mtu,omitempty"`
	PersistentKeepalive int      `json:"persistentKeepalive,omitempty"`
	PreSharedKey        bool     `json:"preSharedKey,omitempty"`
	InterfaceAddress    string   `json:"interfaceAddress,omitempty"`
}

// Service defines the interface for WireGuard operations.
type Service interface {
	ListPeers() ([]Peer, error)
	AddPeer(options AddPeerOptions) (PeerResponse, error)
	RemovePeer(id string) error
	RegeneratePeer(id string) (PeerResponse, error)
	UpdatePeer(id string, updates PeerUpdate) (Peer, error)
	Sync() error
	GetPeerConfig(id string) (string, error)
	GetPeerMetadata(id string) (PeerMetadata, bool)
	GetStats() (Stats, error)
	GetStatsHistory() ([]StatsHistoryItem, error)
	GetSettings() (GlobalSettings, error)
	UpdateSettings(settings GlobalSettings) error
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
func (s *mockService) AddPeer(opts AddPeerOptions) (PeerResponse, error) {
	slog.Warn("Using mock WireGuard service for AddPeer")
	if opts.Name == "force-add-error" {
		return PeerResponse{}, fmt.Errorf("forced error")
	}
	peer := Peer{
		ID:         fmt.Sprintf("mock-peer-%d", len(s.peers)+1),
		PublicKey:  opts.PublicKey,
		Name:       opts.Name,
		AllowedIPs: opts.AllowedIPs,
	}
	if peer.PublicKey == "" {
		peer.PublicKey = "MOCK_PUBKEY_" + peer.ID
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

// Sync is a no-op for mockService.
func (s *mockService) Sync() error {
	slog.Warn("Using mock WireGuard service for Sync (no-op)")
	return nil
}

// GetPeerConfig returns a mock config.
func (s *mockService) GetPeerConfig(id string) (string, error) {
	slog.Warn("Using mock WireGuard service for GetPeerConfig")
	return "[Interface]\nPrivateKey = MOCK_KEY\n...", nil
}

// GetPeerMetadata returns mock metadata.
func (s *mockService) GetPeerMetadata(id string) (PeerMetadata, bool) {
	slog.Warn("Using mock WireGuard service for GetPeerMetadata")
	for _, p := range s.peers {
		if p.ID == id {
			return PeerMetadata{Name: p.Name, PublicKey: p.PublicKey, AllowedIPs: p.AllowedIPs}, true
		}
	}
	return PeerMetadata{}, false
}

// GetStatsHistory returns mock stats history.
func (s *mockService) GetStatsHistory() ([]StatsHistoryItem, error) {
	slog.Warn("Using mock WireGuard service for GetStatsHistory")
	return []StatsHistoryItem{
		{Timestamp: 1706745600, TotalRX: 1000, TotalTX: 500},
		{Timestamp: 1706745660, TotalRX: 1100, TotalTX: 550},
		{Timestamp: 1706745720, TotalRX: 1250, TotalTX: 600},
	}, nil
}

// GetSettings returns mock settings.
func (s *mockService) GetSettings() (GlobalSettings, error) {
	slog.Warn("Using mock WireGuard service for GetSettings")
	return GlobalSettings{
		ServerAddress: "10.0.0.1/24",
		DNS:           "1.1.1.1, 8.8.8.8",
		MTU:           1420,
		Keepalive:     25,
		Endpoint:      "vpn.example.com",
	}, nil
}

// UpdateSettings updates mock settings.
func (s *mockService) UpdateSettings(settings GlobalSettings) error {
	slog.Warn("Using mock WireGuard service for UpdateSettings")
	return nil
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
