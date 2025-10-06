package handlers

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/txnbuild"

	"chama-wallet-backend/config"
	"chama-wallet-backend/services"
)

// wallet handlers
// Creates and returns a wallet
func CreateWallet(c *fiber.Ctx) error {
	address, seed := services.CreateWallet()
	return c.JSON(fiber.Map{
		"address": address,
		"seed":    seed,
		"network": config.Config.Network,
	})
}

func GetBalance(c *fiber.Ctx) error {
	address := c.Params("address")

	client := config.GetHorizonClient()
	accountRequest := horizonclient.AccountRequest{AccountID: address}
	account, err := client.AccountDetail(accountRequest)
	if err != nil {
		// For mainnet, provide more helpful error messages
		if config.Config.IsMainnet {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   "Account not found on mainnet. Please ensure the account is funded with real XLM first.",
				"network": config.Config.Network,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var balances []string
	for _, b := range account.Balances {
		assetInfo := "XLM"
		if b.Asset.Type != "native" {
			assetInfo = fmt.Sprintf("%s:%s", b.Asset.Code, b.Asset.Issuer)
		}
		balances = append(balances, fmt.Sprintf("%s: %s", assetInfo, b.Balance))
	}

	return c.JSON(fiber.Map{
		"balances": balances,
		"network":  config.Config.Network,
	})
}

type TransferRequest struct {
	FromSeed  string `json:"from_seed"`
	ToAddress string `json:"to_address"`
	Amount    string `json:"amount"`
	AssetType string `json:"asset_type,omitempty"` // "XLM" or "USDC"
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

	// Validate transfer limits for mainnet
	if config.Config.IsMainnet {
		amount, _ := strconv.ParseFloat(req.Amount, 64)
		minAmount, _ := strconv.ParseFloat(os.Getenv("MIN_TRANSFER_AMOUNT"), 64)
		maxAmount, _ := strconv.ParseFloat(os.Getenv("MAX_TRANSFER_AMOUNT"), 64)

		if minAmount > 0 && amount < minAmount {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": fmt.Sprintf("Amount below minimum transfer limit of %f", minAmount),
			})
		}

		if maxAmount > 0 && amount > maxAmount {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": fmt.Sprintf("Amount exceeds maximum transfer limit of %f", maxAmount),
			})
		}
	}

	assetType := req.AssetType
	if assetType == "" {
		assetType = "XLM" // Default to XLM
	}

	fmt.Printf("🔄 Processing %s transfer: %s to %s on %s\n", assetType, req.Amount, req.ToAddress, config.Config.Network)
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

	client := config.GetHorizonClient()
	ar := horizonclient.AccountRequest{AccountID: sourceKP.Address()}
	sourceAccount, err := client.AccountDetail(ar)
	if err != nil {
		fmt.Printf("❌ Cannot load source account: %v\n", err)

		// Check if account doesn't exist and try to fund it
		if horizonError, ok := err.(*horizonclient.Error); ok && horizonError.Problem.Status == 404 {
			if config.Config.IsMainnet {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Source account not found on mainnet. Please fund the account with real XLM first.",
				})
			}

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
			if config.Config.IsMainnet {
				fmt.Printf("⚠️ Destination account not found on mainnet: %s\n", req.ToAddress)
				// On mainnet, we can still send to non-existent accounts (they'll be created)
			} else {
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
	}

	// Determine asset type and create appropriate payment operation
	var asset txnbuild.Asset
	var op txnbuild.Operation

	if assetType == "USDC" && config.Config.IsMainnet {
		// USDC payment (mainnet only)
		asset = txnbuild.CreditAsset{
			Code:   config.Config.USDCAssetCode,
			Issuer: config.Config.USDCAssetIssuer,
		}
	} else {
		// XLM payment (default)
		asset = txnbuild.NativeAsset{}
	}

	op = &txnbuild.Payment{
		Destination: req.ToAddress,
		Amount:      req.Amount,
		Asset:       asset,
	}

	// Add memo for mainnet compliance
	var memo txnbuild.Memo
	if config.Config.IsMainnet {
		memo = txnbuild.MemoText(fmt.Sprintf("Chama Wallet %s Transfer", assetType))
	}

	txParams := txnbuild.TransactionParams{
		SourceAccount:        &sourceAccount,
		IncrementSequenceNum: true,
		BaseFee:              txnbuild.MinBaseFee,
		Memo:                 memo,
		Preconditions: txnbuild.Preconditions{
			TimeBounds: txnbuild.NewTimeout(300), // 5 minute timeout
		},
		Operations: []txnbuild.Operation{op},
	}

	tx, err := txnbuild.NewTransaction(txParams)
	if err != nil {
		fmt.Printf("❌ Failed to build transaction: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to build transaction"})
	}

	signedTx, err := tx.Sign(config.GetNetworkPassphrase(), sourceKP)
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

	fmt.Printf("✅ %s transfer successful on %s: %s\n", assetType, config.Config.Network, resp.Hash)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":          "Transfer completed successfully",
		"transaction_hash": resp.Hash,
		"from":             sourceKP.Address(),
		"to":               req.ToAddress,
		"amount":           req.Amount,
		"asset_type":       assetType,
		"network":          config.Config.Network,
		"ledger":           resp.Ledger,
		"explorer_url":     getExplorerURL(resp.Hash),
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
		"network":     config.Config.Network,
		"warning":     getNetworkWarning(),
	})
}

func FundAccount(c *fiber.Ctx) error {
	address := c.Params("address")
	if address == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Address is required"})
	}

	if config.Config.IsMainnet {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Account funding not available on mainnet. Please deposit real XLM to fund your account.",
			"network": config.Config.Network,
		})
	}
	err := services.FundTestAccount(address)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Account funded",
		"address": address,
		"network": config.Config.Network,
	})
}

func GetTransactionHistory(c *fiber.Ctx) error {
	address := c.Params("address")
	if address == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Address is required"})
	}

	client := config.GetHorizonClient()
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
			"hash":         tx.Hash,
			"ledger":       tx.Ledger,
			"memo":         tx.Memo,
			"successful":   tx.Successful,
			"created_at":   tx.LedgerCloseTime.Format("2006-01-02 15:04:05"),
			"fee_charged":  tx.FeeCharged,
			"explorer_url": getExplorerURL(tx.Hash),
		})
	}

	return c.JSON(fiber.Map{
		"transactions": txs,
		"network":      config.Config.Network,
	})
}

// Helper functions
func getExplorerURL(txHash string) string {
	if config.Config.IsMainnet {
		return fmt.Sprintf("https://stellar.expert/explorer/public/tx/%s", txHash)
	}
	return fmt.Sprintf("https://stellar.expert/explorer/testnet/tx/%s", txHash)
}

func getNetworkWarning() string {
	if config.Config.IsMainnet {
		return "⚠️ MAINNET: This keypair can control real funds. Keep the secret key secure!"
	}
	return "ℹ️ TESTNET: This is a test keypair for development purposes only."
}
