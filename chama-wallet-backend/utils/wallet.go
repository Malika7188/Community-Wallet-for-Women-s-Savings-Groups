package utils

import (
	"fmt"

	"github.com/stellar/go/keypair"
)

type StellarWallet struct {
	PublicKey string
	SecretKey string
}

// GenerateStellarWallet creates a new Stellar keypair
func GenerateStellarWallet() (*StellarWallet, error) {
	// Generate a new keypair
	pair, err := keypair.Random()
	if err != nil {
		return nil, fmt.Errorf("failed to generate keypair: %w", err)
	}

	// Validate the generated keys
	if pair.Address() == "" {
		return nil, fmt.Errorf("generated public key is empty")
	}

	if pair.Seed() == "" {
		return nil, fmt.Errorf("generated secret key is empty")
	}

	wallet := &StellarWallet{
		PublicKey: pair.Address(),
		SecretKey: pair.Seed(),
	}

	fmt.Printf("✅ Generated wallet: %s\n", wallet.PublicKey)
	return wallet, nil
}
