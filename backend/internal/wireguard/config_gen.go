package wireguard

import (
	"fmt"
	"strings"
)

// PeerConfigInfo contains all data needed to generate a client .conf file.
type PeerConfigInfo struct {
	PrivateKey          string
	Address             []string
	DNS                 []string
	MTU                 int
	PersistentKeepalive int
	PublicKey           string // Server's public key
	PresharedKey        string
	Endpoint            string   // Server's endpoint
	AllowedIPs          []string // Typically "0.0.0.0/0, ::/0" for full tunnel
}

// GenerateConfigString creates a WireGuard .conf content.
func GenerateConfigString(info PeerConfigInfo) string {
	var sb strings.Builder

	sb.WriteString("[Interface]\n")
	sb.WriteString(fmt.Sprintf("PrivateKey = %s\n", info.PrivateKey))
	if len(info.Address) > 0 {
		sb.WriteString(fmt.Sprintf("Address = %s\n", strings.Join(info.Address, ", ")))
	}
	if len(info.DNS) > 0 {
		sb.WriteString(fmt.Sprintf("DNS = %s\n", strings.Join(info.DNS, ", ")))
	}
	if info.MTU > 0 {
		sb.WriteString(fmt.Sprintf("MTU = %d\n", info.MTU))
	}
	sb.WriteString("\n")

	sb.WriteString("[Peer]\n")
	sb.WriteString(fmt.Sprintf("PublicKey = %s\n", info.PublicKey))
	if info.PresharedKey != "" {
		sb.WriteString(fmt.Sprintf("PresharedKey = %s\n", info.PresharedKey))
	}
	sb.WriteString(fmt.Sprintf("Endpoint = %s\n", info.Endpoint))
	if len(info.AllowedIPs) > 0 {
		sb.WriteString(fmt.Sprintf("AllowedIPs = %s\n", strings.Join(info.AllowedIPs, ", ")))
	}
	if info.PersistentKeepalive > 0 {
		sb.WriteString(fmt.Sprintf("PersistentKeepalive = %d\n", info.PersistentKeepalive))
	}

	return sb.String()
}
