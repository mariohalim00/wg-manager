package wireguard

import (
	"fmt"

	"net"

	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

type realService struct {
	client         *wgctrl.Client
	interfaceName  string
	storage        *Storage
	serverPubKey   string
	serverEndpoint string
	vpnSubnet      string
}

// NewRealService creates and returns a new native WireGuard service.
func NewRealService(interfaceName string, storagePath string, serverEndpoint string, serverPubKey string, vpnSubnet string) (Service, error) {
	client, err := wgctrl.New()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize wgctrl: %w", err)
	}

	// Initialize storage
	storage, err := NewStorage(storagePath)
	if err != nil {
		client.Close()
		return nil, fmt.Errorf("failed to initialize storage: %w", err)
	}

	// Verify we can access the device (checks permissions and existence)
	_, err = client.Device(interfaceName)
	if err != nil {
		client.Close()
		return nil, fmt.Errorf("failed to access device %s: %w", interfaceName, err)
	}

	return &realService{
		client:         client,
		interfaceName:  interfaceName,
		storage:        storage,
		serverPubKey:   serverPubKey,
		serverEndpoint: serverEndpoint,
		vpnSubnet:      vpnSubnet,
	}, nil
}

// Close releases resources held by the realService.
func (s *realService) Close() error {
	if s.client == nil {
		return nil
	}
	return s.client.Close()
}

// ListPeers returns the current list of peers from the WireGuard interface.
func (s *realService) ListPeers() ([]Peer, error) {
	device, err := s.client.Device(s.interfaceName)
	if err != nil {
		return nil, fmt.Errorf("failed to get device %s: %w", s.interfaceName, err)
	}

	peers := make([]Peer, 0, len(device.Peers))
	for _, p := range device.Peers {
		allowedIPs := make([]string, len(p.AllowedIPs))
		for i, ip := range p.AllowedIPs {
			allowedIPs[i] = ip.String()
		}

		name := ""
		if meta, ok := s.storage.GetMetadata(p.PublicKey.String()); ok {
			name = meta.Name
		}

		endpoint := ""
		if p.Endpoint != nil {
			endpoint = p.Endpoint.String()
		}

		peers = append(peers, Peer{
			ID:            p.PublicKey.String(),
			PublicKey:     p.PublicKey.String(),
			Name:          name,
			Endpoint:      endpoint,
			AllowedIPs:    allowedIPs,
			LastHandshake: p.LastHandshakeTime.String(),
			ReceiveBytes:  p.ReceiveBytes,
			TransmitBytes: p.TransmitBytes,
		})
	}
	return peers, nil
}

// AddPeer adds a new peer to the WireGuard interface.
func (s *realService) AddPeer(name string, publicKey string, allowedIPs []string) (PeerResponse, error) {
	var privateKey string
	var err error

	// If publicKey is empty, generate a new key pair
	if publicKey == "" {
		keys, err := GenerateKeyPair()
		if err != nil {
			return PeerResponse{}, fmt.Errorf("failed to generate key pair: %w", err)
		}
		publicKey = keys.PublicKey
		privateKey = keys.PrivateKey
	}

	pubKey, err := wgtypes.ParseKey(publicKey)
	if err != nil {
		return PeerResponse{}, fmt.Errorf("invalid public key: %w", err)
	}

	// Parse allowed IPs
	var allowedIPConfigs []net.IPNet
	for _, ipStr := range allowedIPs {
		_, ipNet, err := net.ParseCIDR(ipStr)
		if err != nil {
			return PeerResponse{}, fmt.Errorf("invalid allowed IP '%s': %w", ipStr, err)
		}
		allowedIPConfigs = append(allowedIPConfigs, *ipNet)
	}

	peerConfig := wgtypes.PeerConfig{
		PublicKey:         pubKey,
		ReplaceAllowedIPs: true,
		AllowedIPs:        allowedIPConfigs,
	}

	config := wgtypes.Config{
		ReplacePeers: false,
		Peers:        []wgtypes.PeerConfig{peerConfig},
	}

	if err := s.client.ConfigureDevice(s.interfaceName, config); err != nil {
		return PeerResponse{}, fmt.Errorf("failed to configure device: %w", err)
	}

	// Save metadata
	if err := s.storage.SetMetadata(publicKey, PeerMetadata{Name: name}); err != nil {
		return PeerResponse{}, fmt.Errorf("failed to save metadata: %w", err)
	}

	response := PeerResponse{
		Peer: Peer{
			ID:         publicKey,
			PublicKey:  publicKey,
			Name:       name,
			AllowedIPs: allowedIPs,
		},
	}

	// Generate config if we have a private key
	if privateKey != "" {
		response.PrivateKey = privateKey
		response.Config = GenerateConfigString(PeerConfigInfo{
			PrivateKey: privateKey,
			Address:    allowedIPs, // Defaulting client address to its allowed IPs on server-side
			PublicKey:  s.serverPubKey,
			Endpoint:   s.serverEndpoint,
			AllowedIPs: []string{"0.0.0.0/0", "::/0"},
		})
	}

	return response, nil
}

// RemovePeer removes a peer from the WireGuard interface.
func (s *realService) RemovePeer(id string) error {
	pubKey, err := wgtypes.ParseKey(id)
	if err != nil {
		return fmt.Errorf("invalid public key: %w", err)
	}

	peerConfig := wgtypes.PeerConfig{
		PublicKey: pubKey,
		Remove:    true,
	}

	config := wgtypes.Config{
		Peers: []wgtypes.PeerConfig{peerConfig},
	}

	if err := s.client.ConfigureDevice(s.interfaceName, config); err != nil {
		return fmt.Errorf("failed to remove peer: %w", err)
	}

	// Delete metadata from storage
	if err := s.storage.DeleteMetadata(id); err != nil {
		return fmt.Errorf("failed to delete metadata: %w", err)
	}
	return nil
}

// GetStats returns the current statistics for the WireGuard interface.
func (s *realService) GetStats() (Stats, error) {
	device, err := s.client.Device(s.interfaceName)
	if err != nil {
		return Stats{}, fmt.Errorf("failed to get device %s: %w", s.interfaceName, err)
	}

	var totalRX, totalTX int64
	for _, p := range device.Peers {
		totalRX += p.ReceiveBytes
		totalTX += p.TransmitBytes
	}

	return Stats{
		InterfaceName: s.interfaceName,
		PublicKey:     device.PublicKey.String(),
		ListenPort:    device.ListenPort,
		Subnet:        s.vpnSubnet,
		PeerCount:     len(device.Peers),
		TotalRX:       totalRX,
		TotalTX:       totalTX,
	}, nil
}
