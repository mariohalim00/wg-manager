package wireguard

import (
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"

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
	history        []StatsHistoryItem
	historyMu      sync.RWMutex
	stopChan       chan struct{}
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

	srv := &realService{
		client:         client,
		interfaceName:  interfaceName,
		storage:        storage,
		serverPubKey:   serverPubKey,
		serverEndpoint: serverEndpoint,
		vpnSubnet:      vpnSubnet,
		history:        make([]StatsHistoryItem, 0, 100),
		stopChan:       make(chan struct{}),
	}

	if err := srv.Sync(); err != nil {
		slog.Error("Failed to sync peers on startup", "error", err)
	}

	// Start background stats collector
	go srv.collectStats()

	return srv, nil
}

// Close releases resources held by the realService.
func (s *realService) Close() error {
	close(s.stopChan)
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
func (s *realService) AddPeer(opts AddPeerOptions) (PeerResponse, error) {
	var privateKey string
	var psk string
	var err error

	// If publicKey is empty, generate a new key pair
	if opts.PublicKey == "" {
		keys, err := GenerateKeyPair()
		if err != nil {
			return PeerResponse{}, fmt.Errorf("failed to generate key pair: %w", err)
		}
		opts.PublicKey = keys.PublicKey
		privateKey = keys.PrivateKey
	}

	if opts.PreSharedKey {
		psk, err = GeneratePresharedKey()
		if err != nil {
			return PeerResponse{}, fmt.Errorf("failed to generate preshared key: %w", err)
		}
	}

	pubKey, err := wgtypes.ParseKey(opts.PublicKey)
	if err != nil {
		return PeerResponse{}, fmt.Errorf("invalid public key: %w", err)
	}

	// Parse allowed IPs
	var allowedIPConfigs []net.IPNet
	for _, ipStr := range opts.AllowedIPs {
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

	if psk != "" {
		pskKey, err := wgtypes.ParseKey(psk)
		if err != nil {
			return PeerResponse{}, fmt.Errorf("invalid preshared key: %w", err)
		}
		peerConfig.PresharedKey = &pskKey
	}

	config := wgtypes.Config{
		ReplacePeers: false,
		Peers:        []wgtypes.PeerConfig{peerConfig},
	}

	if err := s.client.ConfigureDevice(s.interfaceName, config); err != nil {
		return PeerResponse{}, fmt.Errorf("failed to configure device: %w", err)
	}

	// Save metadata
	meta := PeerMetadata{
		PublicKey:           opts.PublicKey,
		PrivateKey:          privateKey,
		PresharedKey:        psk,
		Name:                opts.Name,
		AllowedIPs:          opts.AllowedIPs,
		DNS:                 opts.DNS,
		MTU:                 opts.MTU,
		PersistentKeepalive: opts.PersistentKeepalive,
		InterfaceAddress:    opts.InterfaceAddress,
	}

	// Generate config if we have a private key
	if privateKey != "" {
		settings := s.storage.GetSettings()
		dns := opts.DNS
		if dns == "" {
			dns = settings.DNS
		}
		mtu := opts.MTU
		if mtu == 0 {
			mtu = settings.MTU
		}
		keepalive := opts.PersistentKeepalive
		if keepalive == 0 {
			keepalive = settings.Keepalive
		}
		endpoint := settings.Endpoint
		if endpoint == "" {
			endpoint = s.serverEndpoint
		}

		dnsSplit := []string{}
		if dns != "" {
			for _, d := range strings.Split(dns, ",") {
				dnsSplit = append(dnsSplit, strings.TrimSpace(d))
			}
		}

		address := opts.AllowedIPs
		if opts.InterfaceAddress != "" {
			address = []string{opts.InterfaceAddress}
		}

		meta.Config = GenerateConfigString(PeerConfigInfo{
			PrivateKey:          privateKey,
			Address:             address,
			DNS:                 dnsSplit,
			MTU:                 mtu,
			PersistentKeepalive: keepalive,
			PublicKey:           s.serverPubKey,
			PresharedKey:        psk,
			Endpoint:            endpoint,
			AllowedIPs:          []string{"0.0.0.0/0", "::/0"},
		})
	}

	if err := s.storage.SetMetadata(opts.PublicKey, meta); err != nil {
		return PeerResponse{}, fmt.Errorf("failed to save metadata: %w", err)
	}

	response := PeerResponse{
		Peer: Peer{
			ID:         opts.PublicKey,
			PublicKey:  opts.PublicKey,
			Name:       opts.Name,
			AllowedIPs: opts.AllowedIPs,
		},
		PrivateKey:   meta.PrivateKey,
		PresharedKey: meta.PresharedKey,
		Config:       meta.Config,
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

// RegeneratePeer regenerates keys for a peer.
func (s *realService) RegeneratePeer(id string) (PeerResponse, error) {
	// 1. Fetch existing peer to get metadata and allowed IPs
	peers, err := s.ListPeers()
	if err != nil {
		return PeerResponse{}, fmt.Errorf("failed to list peers: %w", err)
	}

	var targetPeer *Peer
	for _, p := range peers {
		if p.ID == id {
			targetPeer = &p
			break
		}
	}

	if targetPeer == nil {
		return PeerResponse{}, fmt.Errorf("peer not found: %s", id)
	}

	// 2. Remove old peer
	if err := s.RemovePeer(id); err != nil {
		return PeerResponse{}, fmt.Errorf("failed to remove old peer: %w", err)
	}

	// 3. Add back with new keys
	// AddPeer generates new keys if publicKey is empty
	meta, _ := s.storage.GetMetadata(id)
	opts := AddPeerOptions{
		Name:                targetPeer.Name,
		AllowedIPs:          targetPeer.AllowedIPs,
		DNS:                 meta.DNS,
		MTU:                 meta.MTU,
		PersistentKeepalive: meta.PersistentKeepalive,
		PreSharedKey:        meta.PresharedKey != "",
		InterfaceAddress:    meta.InterfaceAddress,
	}

	response, err := s.AddPeer(opts)
	if err != nil {
		return PeerResponse{}, fmt.Errorf("failed to add peer with new keys: %w", err)
	}

	return response, nil
}

// UpdatePeer updates peer metadata or configuration.
func (s *realService) UpdatePeer(id string, updates PeerUpdate) (Peer, error) {
	pubKey, err := wgtypes.ParseKey(id)
	if err != nil {
		return Peer{}, fmt.Errorf("invalid public key: %w", err)
	}

	// Fetch existing metadata
	meta, ok := s.storage.GetMetadata(id)
	if !ok {
		return Peer{}, fmt.Errorf("peer metadata not found: %s", id)
	}

	// Update metadata
	metaChanged := false
	if updates.Name != nil {
		meta.Name = *updates.Name
		metaChanged = true
	}
	if updates.DNS != nil {
		meta.DNS = *updates.DNS
		metaChanged = true
	}
	if updates.MTU != nil {
		meta.MTU = *updates.MTU
		metaChanged = true
	}
	if updates.PersistentKeepalive != nil {
		meta.PersistentKeepalive = *updates.PersistentKeepalive
		metaChanged = true
	}
	if updates.InterfaceAddress != nil {
		meta.InterfaceAddress = *updates.InterfaceAddress
		metaChanged = true
	}

	if metaChanged {
		if err := s.storage.SetMetadata(id, meta); err != nil {
			return Peer{}, fmt.Errorf("failed to update metadata: %w", err)
		}
	}

	// Update WireGuard config if AllowedIPs changed
	if updates.AllowedIPs != nil {
		var allowedIPConfigs []net.IPNet
		for _, ipStr := range *updates.AllowedIPs {
			_, ipNet, err := net.ParseCIDR(ipStr)
			if err != nil {
				return Peer{}, fmt.Errorf("invalid allowed IP '%s': %w", ipStr, err)
			}
			allowedIPConfigs = append(allowedIPConfigs, *ipNet)
		}

		peerConfig := wgtypes.PeerConfig{
			PublicKey:         pubKey,
			UpdateOnly:        true,
			ReplaceAllowedIPs: true,
			AllowedIPs:        allowedIPConfigs,
		}

		config := wgtypes.Config{
			Peers: []wgtypes.PeerConfig{peerConfig},
		}

		if err := s.client.ConfigureDevice(s.interfaceName, config); err != nil {
			return Peer{}, fmt.Errorf("failed to update WireGuard peer config: %w", err)
		}
	}

	// Return updated peer
	peers, err := s.ListPeers()
	if err != nil {
		return Peer{}, fmt.Errorf("failed to list peers after update: %w", err)
	}

	for _, p := range peers {
		if p.ID == id {
			return p, nil
		}
	}

	return Peer{}, fmt.Errorf("peer not found after update: %s", id)
}

// Sync restores all peers from storage to the WireGuard interface.
func (s *realService) Sync() error {
	slog.Info("Syncing peers from storage to interface", "interface", s.interfaceName)
	s.storage.mu.RLock()
	defer s.storage.mu.RUnlock()

	var peerConfigs []wgtypes.PeerConfig
	for pubKeyStr, meta := range s.storage.data.Peers {
		pubKey, err := wgtypes.ParseKey(pubKeyStr)
		if err != nil {
			slog.Error("Invalid public key in storage", "key", pubKeyStr, "error", err)
			continue
		}

		var allowedIPConfigs []net.IPNet
		for _, ipStr := range meta.AllowedIPs {
			_, ipNet, err := net.ParseCIDR(ipStr)
			if err != nil {
				slog.Error("Invalid allowed IP in storage", "ip", ipStr, "error", err)
				continue
			}
			allowedIPConfigs = append(allowedIPConfigs, *ipNet)
		}

		peerConfigs = append(peerConfigs, wgtypes.PeerConfig{
			PublicKey:         pubKey,
			ReplaceAllowedIPs: true,
			AllowedIPs:        allowedIPConfigs,
		})
	}

	if len(peerConfigs) == 0 {
		return nil
	}

	config := wgtypes.Config{
		ReplacePeers: false, // We only want to ADD/UPDATE peers from our storage
		Peers:        peerConfigs,
	}

	if err := s.client.ConfigureDevice(s.interfaceName, config); err != nil {
		return fmt.Errorf("failed to sync peers to device: %w", err)
	}

	return nil
}

// GetPeerConfig returns the configuration string for a peer.
func (s *realService) GetPeerConfig(id string) (string, error) {
	meta, ok := s.storage.GetMetadata(id)
	if !ok {
		return "", fmt.Errorf("peer not found: %s", id)
	}
	if meta.Config == "" {
		return "", fmt.Errorf("config not available for peer (might need key regeneration): %s", id)
	}
	return meta.Config, nil
}

// GetPeerMetadata returns metadata for a peer.
func (s *realService) GetPeerMetadata(id string) (PeerMetadata, bool) {
	return s.storage.GetMetadata(id)
}

// GetStatsHistory returns historical statistics.
func (s *realService) GetStatsHistory() ([]StatsHistoryItem, error) {
	s.historyMu.RLock()
	defer s.historyMu.RUnlock()

	// Return a copy to avoid data races
	cp := make([]StatsHistoryItem, len(s.history))
	copy(cp, s.history)
	return cp, nil
}

// GetSettings returns application-wide settings.
func (s *realService) GetSettings() (GlobalSettings, error) {
	return s.storage.GetSettings(), nil
}

// UpdateSettings updates application-wide settings.
func (s *realService) UpdateSettings(settings GlobalSettings) error {
	return s.storage.UpdateSettings(settings)
}

func (s *realService) collectStats() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-s.stopChan:
			return
		case <-ticker.C:
			stats, err := s.GetStats()
			if err != nil {
				slog.Error("Failed to collect stats for history", "error", err)
				continue
			}

			s.historyMu.Lock()
			if len(s.history) >= 100 {
				s.history = s.history[1:]
			}
			s.history = append(s.history, StatsHistoryItem{
				Timestamp: time.Now().Unix(),
				TotalRX:   stats.TotalRX,
				TotalTX:   stats.TotalTX,
			})
			s.historyMu.Unlock()
		}
	}
}
