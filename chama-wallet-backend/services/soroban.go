// services/soroban.go
package services

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/network"
)

const (
	sorobanRPC = "https://soroban-testnet.stellar.org"
	// contractID = "CADHKUC557DJ2F2XGEO4BGHFIYQ6O5QDVNG637ANRAGPBSWXMXXPMOI4"
)

var (
	rpcURL            = "https://soroban-testnet.stellar.org:443"
	contractID        = "CADHKUC557DJ2F2XGEO4BGHFIYQ6O5QDVNG637ANRAGPBSWXMXXPMOI4" // Your contract
	networkPassphrase = network.TestNetworkPassphrase
	signerSecret      = os.Getenv("SIGNER_SECRET") // Or hardcode temporarily
	// sorobanClient     = sorobanrpc.NewClient("https://soroban-testnet.stellar.org:443")
	client = horizonclient.DefaultTestNetClient
)

type SorobanInvokeRequest struct {
	ContractID string        `json:"contract_id"`
	Function   string        `json:"function"`
	Args       []interface{} `json:"args"`
}

// CallSorobanFunction executes a Soroban contract function
func CallSorobanFunction(contractID, functionName string, args []string) (string, error) {
	// Updated command structure based on the new Soroban CLI syntax
	cmd := []string{
		"contract", "invoke",
		"--id", contractID,
		"--source-account", "<WALLET_ADDRESS>", // Replace with actual source account
		"--network", "testnet",
		"--",
		functionName,
	}

	// Add the function arguments
	cmd = append(cmd, args...)

	// Execute the command
	out, err := exec.Command("soroban", cmd...).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("soroban error: %s", string(out))
	}

	return strings.TrimSpace(string(out)), nil
}

// Contribute function wrapper
func Contribute(contractID, userAddress, amount string) (string, error) {
	args := []string{userAddress, amount}
	return CallSorobanFunction(contractID, "contribute", args)
}

// GetBalance function wrapper
func GetBalance(contractID, userAddress string) (string, error) {
	args := []string{userAddress}
	return CallSorobanFunction(contractID, "get_balance", args)
}

// Withdraw function wrapper
func Withdraw(contractID, userAddress, amount string) (string, error) {
	args := []string{userAddress, amount}
	return CallSorobanFunction(contractID, "withdraw", args)
}

// GetContributionHistory function wrapper
func GetContributionHistory(contractID, userAddress string) (string, error) {
	args := []string{userAddress}
	return CallSorobanFunction(contractID, "get_contribution_history", args)
}
