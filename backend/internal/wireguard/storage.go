package wireguard

import (
	"encoding/json"
	"os"
	"sync"
)

// PeerMetadata stores persistent information about a peer.
type PeerMetadata struct {
	PublicKey  string   `json:"publicKey"`
	PrivateKey string   `json:"privateKey,omitempty"`
	Name       string   `json:"name"`
	AllowedIPs []string `json:"allowedIPs"`
	Config     string   `json:"config,omitempty"`
}

// Storage handles persistent storage of peer metadata.
type Storage struct {
	path string
	mu   sync.RWMutex
	data map[string]PeerMetadata
}

// NewStorage creates a new storage manager.
func NewStorage(path string) (*Storage, error) {
	s := &Storage{
		path: path,
		data: make(map[string]PeerMetadata),
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
	m, ok := s.data[publicKey]
	return m, ok
}

// SetMetadata updates metadata for a peer.
func (s *Storage) SetMetadata(publicKey string, metadata PeerMetadata) error {
	s.mu.Lock()
	s.data[publicKey] = metadata
	s.mu.Unlock()

	return s.save()
}

// DeleteMetadata removes metadata for a peer.
func (s *Storage) DeleteMetadata(publicKey string) error {
	s.mu.Lock()
	delete(s.data, publicKey)
	s.mu.Unlock()

	return s.save()
}
