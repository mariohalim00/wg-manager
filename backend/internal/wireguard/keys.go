package wireguard

import (
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

// Keys represents a WireGuard key pair.
type Keys struct {
	PrivateKey string `json:"privateKey"`
	PublicKey  string `json:"publicKey"`
}

// GenerateKeyPair generates a new WireGuard private/public key pair.
func GenerateKeyPair() (Keys, error) {
	key, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		return Keys{}, err
	}

	return Keys{
		PrivateKey: key.String(),
		PublicKey:  key.PublicKey().String(),
	}, nil
}

// GeneratePresharedKey generates a new WireGuard preshared key.
func GeneratePresharedKey() (string, error) {
	key, err := wgtypes.GenerateKey()
	if err != nil {
		return "", err
	}
	return key.String(), nil
}
