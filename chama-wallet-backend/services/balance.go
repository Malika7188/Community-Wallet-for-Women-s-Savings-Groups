package services

import (
	"fmt"
	"time"

	"github.com/stellar/go/clients/horizonclient"

	"chama-wallet-backend/config"
)

func CheckBalance(address string) (string, error) {
	client := config.GetHorizonClient()

	// First try to get account details
	account, err := client.AccountDetail(horizonclient.AccountRequest{AccountID: address})
	if err != nil {
		// Check if it's a "Resource Missing" error (account doesn't exist)
		if horizonError, ok := err.(*horizonclient.Error); ok {
			if horizonError.Problem.Status == 404 {
				if config.Config.IsMainnet {
					return "0", fmt.Errorf("account not found on mainnet - account needs to be funded with real XLM first")
				}

				fmt.Printf("⚠️ Account %s not found on testnet. Attempting to fund...\n", address)

				// Try to fund the account
				if fundErr := FundTestAccount(address); fundErr != nil {
					return "0", fmt.Errorf("account not found and funding failed: %w", fundErr)
				}

				// Wait a moment for the funding to process
				time.Sleep(2 * time.Second)

				// Try again to get account details
				account, err = client.AccountDetail(horizonclient.AccountRequest{AccountID: address})
				if err != nil {
					return "0", fmt.Errorf("account still not found after funding: %w", err)
				}
			} else {
				return "0", fmt.Errorf("failed to get account details: %w", err)
			}
		} else {
			return "0", fmt.Errorf("failed to get account details: %w", err)
		}
	}

	fmt.Printf("🔍 Balances for %s (%s)\n", address, config.Config.Network)
	var totalBalance string = "0"

	for _, b := range account.Balances {
		assetInfo := "native"
		if b.Asset.Type != "native" {
			assetInfo = fmt.Sprintf("%s:%s", b.Asset.Code, b.Asset.Issuer)
		}
		fmt.Printf(" - Asset: %s | Balance: %s\n", assetInfo, b.Balance)

		if b.Asset.Type == "native" {
			totalBalance = b.Balance
		}
	}

	return totalBalance, nil
}

// CheckUSDCBalance returns the USDC balance of a wallet (mainnet only)
func CheckUSDCBalance(address string) (string, error) {
	if !config.Config.IsMainnet {
		return "0", fmt.Errorf("USDC balance checking only available on mainnet")
	}

	client := config.GetHorizonClient()
	account, err := client.AccountDetail(horizonclient.AccountRequest{AccountID: address})
	if err != nil {
		return "0", fmt.Errorf("failed to get account details: %w", err)
	}

	for _, b := range account.Balances {
		if b.Asset.Type != "native" &&
			b.Asset.Code == config.Config.USDCAssetCode &&
			b.Asset.Issuer == config.Config.USDCAssetIssuer {
			return b.Balance, nil
		}
	}

	return "0", nil // No USDC balance found
}
