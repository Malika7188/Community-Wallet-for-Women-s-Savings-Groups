package services

import (
	"fmt"
	"log"

	"github.com/stellar/go/clients/horizonclient"
)

func CheckBalance(address string) (string, error) {
	client := horizonclient.DefaultTestNetClient

	account, err := client.AccountDetail(horizonclient.AccountRequest{AccountID: address})
	if err != nil {
		log.Fatalf("Failed to get account details: %v", err)
	}

	fmt.Println("🔍 Balances for", address)
	for _, b := range account.Balances {
		fmt.Printf(" - Type: %s | Balance: %s\n", b.Asset.Type, b.Balance)
	}
	return "0", nil
}
