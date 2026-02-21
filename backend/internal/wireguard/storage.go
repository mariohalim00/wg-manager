package wireguard

import (
	"encoding/json"
	"os"
	"sync"
)

// PeerMetadata stores persistent information about a peer.
type PeerMetadata struct {
	PublicKey           string   `json:"publicKey"`
	PrivateKey          string   `json:"privateKey,omitempty"`
	PresharedKey        string   `json:"presharedKey,omitempty"`
	Name                string   `json:"name"`
	AllowedIPs          []string `json:"allowedIPs"`
	DNS                 string   `json:"dns,omitempty"`
	MTU                 int      `json:"mtu,omitempty"`
	PersistentKeepalive int      `json:"persistentKeepalive,omitempty"`
	InterfaceAddress    string   `json:"interfaceAddress,omitempty"`
	Config              string   `json:"config,omitempty"`
}

// GlobalSettings stores application-wide WireGuard settings.
type GlobalSettings struct {
	ServerAddress string `json:"serverAddress"`
	DNS           string `json:"dns"`
	MTU           int    `json:"mtu"`
	Keepalive     int    `json:"keepalive"`
	Endpoint      string `json:"endpoint"`
}

// storageContainer is used for JSON marshaling/unmarshaling of all persistent data.
type storageContainer struct {
	Peers    map[string]PeerMetadata `json:"peers"`
	Settings GlobalSettings          `json:"settings"`
}

// Storage handles persistent storage of peer metadata and settings.
type Storage struct {
	path string
	mu   sync.RWMutex
	data storageContainer
}

// NewStorage creates a new storage manager.
func NewStorage(path string) (*Storage, error) {
	s := &Storage{
		path: path,
		data: storageContainer{
			Peers: make(map[string]PeerMetadata),
			Settings: GlobalSettings{
				DNS: "1.1.1.1, 8.8.8.8",
				MTU: 1420,
			},
		},
	}

	if err := s.load(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	return s, nil
}

func (s *Storage) load() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	file, err := os.Open(s.path)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewDecoder(file).Decode(&s.data)
}

func (s *Storage) save() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	file, err := os.Create(s.path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(s.data)
}

// GetMetadata returns metadata for a peer.
func (s *Storage) GetMetadata(publicKey string) (PeerMetadata, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	m, ok := s.data.Peers[publicKey]
	return m, ok
}

// SetMetadata updates metadata for a peer.
func (s *Storage) SetMetadata(publicKey string, metadata PeerMetadata) error {
	s.mu.Lock()
	s.data.Peers[publicKey] = metadata
	s.mu.Unlock()

	return s.save()
}

// DeleteMetadata removes metadata for a peer.
func (s *Storage) DeleteMetadata(publicKey string) error {
	s.mu.Lock()
	delete(s.data.Peers, publicKey)
	s.mu.Unlock()

	return s.save()
}

// GetSettings returns application-wide settings.
func (s *Storage) GetSettings() GlobalSettings {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.data.Settings
}

// UpdateSettings updates application-wide settings.
func (s *Storage) UpdateSettings(settings GlobalSettings) error {
	s.mu.Lock()
	s.data.Settings = settings
	s.mu.Unlock()

	return s.save()
}
