package services

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	"github.com/stellar/go/protocols/horizon"
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
func SendXLM(seed, destination, amount string) (horizon.Transaction, error) {
	client := horizonclient.DefaultTestNetClient

	// Load source account
	kp, err := keypair.ParseFull(seed)
	if err != nil {
		return horizon.Transaction{}, err
	}

	ar := horizonclient.AccountRequest{AccountID: kp.Address()}
	sourceAccount, err := client.AccountDetail(ar)
	if err != nil {
		return horizon.Transaction{}, err
	}

	// Build the transaction
	op := txnbuild.Payment{
		Destination: destination,
		Amount:      amount,
		Asset:       txnbuild.NativeAsset{},
	}
	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &sourceAccount,
			IncrementSequenceNum: true,
			Operations:           []txnbuild.Operation{&op},
			BaseFee:              txnbuild.MinBaseFee,
			// Timebounds:           txnbuild.NewInfiniteTimeout(),
		},
	)
	if err != nil {
		return horizon.Transaction{}, err
	}

	tx, err = tx.Sign(network.TestNetworkPassphrase, kp)
	if err != nil {
		return horizon.Transaction{}, err
	}

	txeBase64, err := tx.Base64()
	if err != nil {
		return horizon.Transaction{}, err
	}

	resp, err := client.SubmitTransactionXDR(txeBase64)
	if err != nil {
		return horizon.Transaction{}, err
	}

	return resp, nil
}
