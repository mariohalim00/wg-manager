package wireguard

import (
	"fmt"
	"strings"
)

// PeerConfigInfo contains all data needed to generate a client .conf file.
type PeerConfigInfo struct {
	PrivateKey string
	Address    []string
	DNS        []string
	PublicKey  string   // Server's public key
	Endpoint   string   // Server's endpoint
	AllowedIPs []string // Typically "0.0.0.0/0, ::/0" for full tunnel
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
	sb.WriteString("\n")

	sb.WriteString("[Peer]\n")
	sb.WriteString(fmt.Sprintf("PublicKey = %s\n", info.PublicKey))
	sb.WriteString(fmt.Sprintf("Endpoint = %s\n", info.Endpoint))
	if len(info.AllowedIPs) > 0 {
		sb.WriteString(fmt.Sprintf("AllowedIPs = %s\n", strings.Join(info.AllowedIPs, ", ")))
	}

	return sb.String()
}
