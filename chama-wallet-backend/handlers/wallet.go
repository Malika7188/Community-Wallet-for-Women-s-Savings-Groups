package handlers

import (
	"fmt"
	"strconv"
	"strings"
	"time"

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
		fmt.Printf("❌ Failed to parse transfer request: %v\n", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	// Validate required fields
	if req.FromSeed == "" || req.ToAddress == "" || req.Amount == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing required fields: from_seed, to_address, and amount are required",
		})
	}

	// Validate amount is positive
	if amount, err := strconv.ParseFloat(req.Amount, 64); err != nil || amount <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Amount must be a positive number",
		})
	}

	fmt.Printf("🔄 Processing transfer: %s XLM to %s\n", req.Amount, req.ToAddress)

	sourceKP, err := keypair.ParseFull(req.FromSeed)
	if err != nil {
		fmt.Printf("❌ Invalid secret key: %v\n", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid secret key"})
	}

	// Validate destination address format
	if len(req.ToAddress) != 56 || !strings.HasPrefix(req.ToAddress, "G") {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid destination address format",
		})
	}

	client := horizonclient.DefaultTestNetClient
	ar := horizonclient.AccountRequest{AccountID: sourceKP.Address()}
	sourceAccount, err := client.AccountDetail(ar)
	if err != nil {
		fmt.Printf("❌ Cannot load source account: %v\n", err)
		
		// Check if account doesn't exist and try to fund it
		if horizonError, ok := err.(*horizonclient.Error); ok && horizonError.Problem.Status == 404 {
			fmt.Printf("🔄 Source account not found, attempting to fund: %s\n", sourceKP.Address())
			if fundErr := services.FundTestAccount(sourceKP.Address()); fundErr != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Source account not found and funding failed",
				})
			}
			
			// Wait for funding to process
			time.Sleep(3 * time.Second)
			
			// Try to load account again
			sourceAccount, err = client.AccountDetail(ar)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Cannot load source account after funding",
				})
			}
		} else {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Cannot load source account",
			})
		}
	}

	// Check if destination account exists, if not try to create it
	destAccountRequest := horizonclient.AccountRequest{AccountID: req.ToAddress}
	_, err = client.AccountDetail(destAccountRequest)
	if err != nil {
		if horizonError, ok := err.(*horizonclient.Error); ok && horizonError.Problem.Status == 404 {
			fmt.Printf("🔄 Destination account not found, attempting to fund: %s\n", req.ToAddress)
			if fundErr := services.FundTestAccount(req.ToAddress); fundErr != nil {
				fmt.Printf("⚠️ Warning: Could not fund destination account: %v\n", fundErr)
				// Continue with transfer anyway - account will be created
			} else {
				// Wait for funding to process
				time.Sleep(2 * time.Second)
			}
		}
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
		Preconditions: txnbuild.Preconditions{
			TimeBounds: txnbuild.NewTimeout(300), // 5 minute timeout
		},
		Operations: []txnbuild.Operation{&op},
	}

	tx, err := txnbuild.NewTransaction(txParams)
	if err != nil {
		fmt.Printf("❌ Failed to build transaction: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to build transaction"})
	}

	signedTx, err := tx.Sign(network.TestNetworkPassphrase, sourceKP)
	if err != nil {
		fmt.Printf("❌ Failed to sign transaction: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to sign transaction"})
	}

	resp, err := client.SubmitTransaction(signedTx)
	if err != nil {
		fmt.Printf("❌ Failed to submit transaction: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Transaction submission failed: %v", err),
		})
	}

	fmt.Printf("✅ Transfer successful: %s\n", resp.Hash)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":          "Transfer completed successfully",
		"transaction_hash": resp.Hash,
		"from":            sourceKP.Address(),
		"to":              req.ToAddress,
		"amount":          req.Amount,
		"ledger":          resp.Ledger,
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
