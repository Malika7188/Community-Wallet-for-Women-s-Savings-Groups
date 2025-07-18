package utils

import (
	"github.com/stellar/go/keypair"
)

type StellarWallet struct {
	PublicKey string
	SecretKey string
}

// GenerateStellarWallet creates a new Stellar keypair
func GenerateStellarWallet() (StellarWallet, error) {
	kp, err := keypair.Random()
	if err != nil {
		return StellarWallet{}, err
	}
	return StellarWallet{
		PublicKey: kp.Address(),
		SecretKey: kp.Seed(),
	}, nil
}
