package services

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	"github.com/stellar/go/txnbuild"
)

var client = horizonclient.DefaultTestNetClient

// CreateWallet generates a new Stellar keypair
func CreateWallet() (string, string) {
	kp, err := keypair.Random()
	if err != nil {
		log.Fatal(err)
	}
	return kp.Address(), kp.Seed()
}

// FundWallet funds a testnet wallet using Friendbot
func FundWallet(address string) error {
	url := fmt.Sprintf("https://friendbot.stellar.org/?addr=%s", address)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("Friendbot response:", string(body))
	return nil
}

// GetBalance returns the XLM balance of a wallet
func GetBalance(address string) (string, error) {
	ar := horizonclient.AccountRequest{AccountID: address}
	account, err := client.AccountDetail(ar)
	if err != nil {
		return "", err
	}
	for _, b := range account.Balances {
		if b.Asset.Type == "native" {
			return b.Balance, nil
		}
	}
	return "0", nil
}

// SendXLM transfers XLM from sender to receiver
func SendXLM(senderSecret, receiverAddress, amount string) (string, error) {
	// Load sender keypair
	senderKP, err := keypair.ParseFull(senderSecret)
	if err != nil {
		return "", err
	}

	// Load sender account
	accountRequest := horizonclient.AccountRequest{AccountID: senderKP.Address()}
	sourceAccount, err := client.AccountDetail(accountRequest)
	if err != nil {
		return "", err
	}

	// Create payment operation
	payment := txnbuild.Payment{
		Destination: receiverAddress,
		Amount:      amount,
		Asset:       txnbuild.NativeAsset{},
	}

	// Build transaction
	txParams := txnbuild.TransactionParams{
		SourceAccount:        &sourceAccount,
		IncrementSequenceNum: true,
		Operations:           []txnbuild.Operation{&payment},
		BaseFee:              txnbuild.MinBaseFee,
		Preconditions: txnbuild.Preconditions{
			TimeBounds: txnbuild.NewInfiniteTimeout(),
		},
	}

	tx, err := txnbuild.NewTransaction(txParams)
	if err != nil {
		return "", err
	}

	// Sign
	tx, err = tx.Sign(network.TestNetworkPassphrase, senderKP)
	if err != nil {
		return "", err
	}

	// Submit
	resp, err := client.SubmitTransaction(tx)
	if err != nil {
		return "", err
	}

	fmt.Println("Transaction Successful! Hash:", resp.Hash)
	return resp.Hash, nil
}
