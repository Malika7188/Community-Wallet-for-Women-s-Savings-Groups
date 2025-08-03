package services

import (
	"fmt"
	"time"

	"github.com/stellar/go/clients/horizonclient"
)

func CheckBalance(address string) (string, error) {
	client := horizonclient.DefaultTestNetClient

	// First try to get account details
	account, err := client.AccountDetail(horizonclient.AccountRequest{AccountID: address})
	if err != nil {
		// Check if it's a "Resource Missing" error (account doesn't exist)
		if horizonError, ok := err.(*horizonclient.Error); ok {
			if horizonError.Problem.Status == 404 {
				fmt.Printf("⚠️ Account %s not found on network. Attempting to fund...\n", address)

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

	fmt.Printf("🔍 Balances for %s\n", address)
	var totalBalance string = "0"

	for _, b := range account.Balances {
		fmt.Printf(" - Type: %s | Balance: %s\n", b.Asset.Type, b.Balance)
		if b.Asset.Type == "native" {
			totalBalance = b.Balance
		}
	}

	return totalBalance, nil
}
