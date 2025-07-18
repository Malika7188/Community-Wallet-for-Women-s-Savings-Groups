package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	"github.com/stellar/go/txnbuild"

	"chama-wallet-backend/services"
)

var horizonClient = horizonclient.DefaultTestNetClient

// wallet handlers
// Creates and returns a wallet
func CreateWallet(c *fiber.Ctx) error {
	address, seed := services.CreateWallet()
	return c.JSON(fiber.Map{
		"address": address,
		"seed":    seed,
	})
}

func GetBalance(c *fiber.Ctx) error {
	address := c.Params("address")

	accountRequest := horizonclient.AccountRequest{AccountID: address}
	account, err := horizonClient.AccountDetail(accountRequest)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var balances []string
	for _, b := range account.Balances {
		balances = append(balances, fmt.Sprintf("%s: %s", b.Asset.Type, b.Balance))
	}

	return c.JSON(fiber.Map{
		"balances": balances,
	})
}

type TransferRequest struct {
	FromSeed  string `json:"from_seed"`
	ToAddress string `json:"to_address"`
	Amount    string `json:"amount"`
}

func TransferFunds(c *fiber.Ctx) error {
	var req TransferRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	sourceKP, err := keypair.ParseFull(req.FromSeed)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid secret key"})
	}

	client := horizonclient.DefaultTestNetClient
	ar := horizonclient.AccountRequest{AccountID: sourceKP.Address()}
	sourceAccount, err := client.AccountDetail(ar)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot load source account"})
	}

	op := txnbuild.Payment{
		Destination: req.ToAddress,
		Amount:      req.Amount,
		Asset:       txnbuild.NativeAsset{},
	}

	txParams := txnbuild.TransactionParams{
		SourceAccount:        &sourceAccount,
		IncrementSequenceNum: true,
		BaseFee:              txnbuild.MinBaseFee,
		// Timebounds:           txnbuild.NewTimeout(300),
		Operations: []txnbuild.Operation{&op},
	}

	tx, err := txnbuild.NewTransaction(txParams)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to build transaction"})
	}

	signedTx, err := tx.Sign(network.TestNetworkPassphrase, sourceKP)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to sign transaction"})
	}

	resp, err := client.SubmitTransaction(signedTx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"transaction_hash": resp.Hash,
	})
}

func GenerateKeypair(c *fiber.Ctx) error {
	kp, err := keypair.Random()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"public_key":  kp.Address(),
		"secret_seed": kp.Seed(),
	})
}

func FundAccount(c *fiber.Ctx) error {
	address := c.Params("address")
	if address == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Address is required"})
	}

	err := services.FundTestAccount(address)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Account funded",
		"address": address,
	})
}

func GetTransactionHistory(c *fiber.Ctx) error {
	address := c.Params("address")
	if address == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Address is required"})
	}

	client := horizonclient.DefaultTestNetClient
	txRequest := horizonclient.TransactionRequest{
		ForAccount: address,
		Limit:      10,
		Order:      horizonclient.OrderDesc,
	}

	txPage, err := client.Transactions(txRequest)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	var txs []fiber.Map
	for _, tx := range txPage.Embedded.Records {
		txs = append(txs, fiber.Map{
			"hash":        tx.Hash,
			"ledger":      tx.Ledger,
			"memo":        tx.Memo,
			"successful":  tx.Successful,
			"created_at":  tx.LedgerCloseTime.Format("2006-01-02 15:04:05"),
			"fee_charged": tx.FeeCharged,
		})
	}

	return c.JSON(fiber.Map{"transactions": txs})
}
