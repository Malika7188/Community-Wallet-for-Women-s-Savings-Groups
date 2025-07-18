package services

import (
	"fmt"

	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	"github.com/stellar/go/txnbuild"
)

func SendPayment(fromSecret, toAddress, amount string) error {
	senderKP, err := keypair.ParseFull(fromSecret)
	if err != nil {
		return fmt.Errorf("invalid secret key: %w", err)
	}

	client := horizonclient.DefaultTestNetClient
	ar := horizonclient.AccountRequest{AccountID: senderKP.Address()}
	sourceAccount, err := client.AccountDetail(ar)
	if err != nil {
		return fmt.Errorf("could not load source account: %w", err)
	}

	op := txnbuild.Payment{
		Destination: toAddress,
		Amount:      amount,
		Asset:       txnbuild.NativeAsset{},
	}
	txParams := txnbuild.TransactionParams{
		SourceAccount: &sourceAccount,
		Operations:    []txnbuild.Operation{&op},
		BaseFee:       txnbuild.MinBaseFee,
		Preconditions: txnbuild.Preconditions{
			TimeBounds: txnbuild.NewInfiniteTimeout(), // ✅ This returns *Timebounds
		},
		IncrementSequenceNum: true,
	}

	tx, err := txnbuild.NewTransaction(txParams)
	if err != nil {
		return fmt.Errorf("cannot build tx: %w", err)
	}

	tx, err = tx.Sign(network.TestNetworkPassphrase, senderKP)
	if err != nil {
		return fmt.Errorf("cannot sign tx: %w", err)
	}

	_, err = client.SubmitTransaction(tx)
	if err != nil {
		return fmt.Errorf("tx failed: %w", err)
	}

	return nil
}
