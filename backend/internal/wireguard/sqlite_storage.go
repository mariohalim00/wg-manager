package wireguard

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

// SQLiteStorage handles persistent storage of peer metadata and settings using SQLite.
type SQLiteStorage struct {
	db *sql.DB
	mu sync.RWMutex
}

// NewSQLiteStorage creates and initializes a new SQLite storage manager.
func NewSQLiteStorage(dbPath string) (*SQLiteStorage, error) {
	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite database: %w", err)
	}

	s := &SQLiteStorage{
		db: db,
	}

	if err := s.initSchema(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return s, nil
}

func (s *SQLiteStorage) initSchema() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Create peers table
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS peers (
			public_key TEXT PRIMARY KEY,
			private_key TEXT,
			preshared_key TEXT,
			name TEXT,
			allowed_ips TEXT,
			dns TEXT,
			mtu INTEGER,
			persistent_keepalive INTEGER,
			interface_address TEXT
		)
	`)
	if err != nil {
		return err
	}

	// Create settings table
	_, err = s.db.Exec(`
		CREATE TABLE IF NOT EXISTS settings (
			id INTEGER PRIMARY KEY CHECK (id = 1),
			server_address TEXT,
			dns TEXT,
			mtu INTEGER,
			keepalive INTEGER,
			endpoint TEXT
		)
	`)
	if err != nil {
		return err
	}

	// Insert default settings if they don't exist
	var count int
	err = s.db.QueryRow(`SELECT COUNT(*) FROM settings WHERE id = 1`).Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		_, err = s.db.Exec(`
			INSERT INTO settings (id, server_address, dns, mtu, keepalive, endpoint)
			VALUES (1, '', '1.1.1.1, 8.8.8.8', 1420, 25, '')
		`)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetMetadata returns metadata for a peer.
func (s *SQLiteStorage) GetMetadata(publicKey string) (PeerMetadata, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var meta PeerMetadata
	var allowedIPsJSON string

	err := s.db.QueryRow(`
		SELECT public_key, private_key, preshared_key, name, allowed_ips, 
		       dns, mtu, persistent_keepalive, interface_address
		FROM peers 
		WHERE public_key = ?
	`, publicKey).Scan(
		&meta.PublicKey,
		&meta.PrivateKey,
		&meta.PresharedKey,
		&meta.Name,
		&allowedIPsJSON,
		&meta.DNS,
		&meta.MTU,
		&meta.PersistentKeepalive,
		&meta.InterfaceAddress,
	)

	if err == sql.ErrNoRows {
		return PeerMetadata{}, false
	} else if err != nil {
		return PeerMetadata{}, false
	}

	json.Unmarshal([]byte(allowedIPsJSON), &meta.AllowedIPs)
	return meta, true
}

// SetMetadata updates metadata for a peer.
func (s *SQLiteStorage) SetMetadata(publicKey string, metadata PeerMetadata) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	allowedIPsJSON, err := json.Marshal(metadata.AllowedIPs)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(`
		INSERT INTO peers (
			public_key, private_key, preshared_key, name, allowed_ips, 
			dns, mtu, persistent_keepalive, interface_address
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(public_key) DO UPDATE SET
			private_key=excluded.private_key,
			preshared_key=excluded.preshared_key,
			name=excluded.name,
			allowed_ips=excluded.allowed_ips,
			dns=excluded.dns,
			mtu=excluded.mtu,
			persistent_keepalive=excluded.persistent_keepalive,
			interface_address=excluded.interface_address
	`,
		metadata.PublicKey,
		metadata.PrivateKey,
		metadata.PresharedKey,
		metadata.Name,
		string(allowedIPsJSON),
		metadata.DNS,
		metadata.MTU,
		metadata.PersistentKeepalive,
		metadata.InterfaceAddress,
	)

	return err
}

// DeleteMetadata removes metadata for a peer.
func (s *SQLiteStorage) DeleteMetadata(publicKey string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.db.Exec(`DELETE FROM peers WHERE public_key = ?`, publicKey)
	return err
}

// GetSettings returns application-wide settings.
func (s *SQLiteStorage) GetSettings() GlobalSettings {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var settings GlobalSettings
	err := s.db.QueryRow(`
		SELECT server_address, dns, mtu, keepalive, endpoint
		FROM settings 
		WHERE id = 1
	`).Scan(
		&settings.ServerAddress,
		&settings.DNS,
		&settings.MTU,
		&settings.Keepalive,
		&settings.Endpoint,
	)

	if err != nil {
		// Return safe defaults if retrieval fails
		return GlobalSettings{
			DNS: "1.1.1.1, 8.8.8.8",
			MTU: 1420,
		}
	}

	return settings
}

// UpdateSettings updates application-wide settings.
func (s *SQLiteStorage) UpdateSettings(settings GlobalSettings) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.db.Exec(`
		UPDATE settings SET
			server_address = ?,
			dns = ?,
			mtu = ?,
			keepalive = ?,
			endpoint = ?
		WHERE id = 1
	`,
		settings.ServerAddress,
		settings.DNS,
		settings.MTU,
		settings.Keepalive,
		settings.Endpoint,
	)

	return err
}

// GetAllPeers is a helper method to return all stored peer metadata, specifically for Sync() operations.
func (s *SQLiteStorage) GetAllPeers() (map[string]PeerMetadata, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	rows, err := s.db.Query(`
		SELECT public_key, private_key, preshared_key, name, allowed_ips, 
		       dns, mtu, persistent_keepalive, interface_address
		FROM peers
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	peers := make(map[string]PeerMetadata)
	for rows.Next() {
		var meta PeerMetadata
		var allowedIPsJSON string

		err := rows.Scan(
			&meta.PublicKey,
			&meta.PrivateKey,
			&meta.PresharedKey,
			&meta.Name,
			&allowedIPsJSON,
			&meta.DNS,
			&meta.MTU,
			&meta.PersistentKeepalive,
			&meta.InterfaceAddress,
		)
		if err != nil {
			return nil, err
		}

		json.Unmarshal([]byte(allowedIPsJSON), &meta.AllowedIPs)
		peers[meta.PublicKey] = meta
	}

	return peers, nil
}
