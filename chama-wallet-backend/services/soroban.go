// services/soroban.go
package services

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// var (
// sorobanRPC = "https://soroban-testnet.stellar.org"
// contractID = "CADHKUC557DJ2F2XGEO4BGHFIYQ6O5QDVNG637ANRAGPBSWXMXXPMOI4"
// )

// var (
// 	rpcURL            = "https://soroban-testnet.stellar.org:443"
// 	contractID        = "CADHKUC557DJ2F2XGEO4BGHFIYQ6O5QDVNG637ANRAGPBSWXMXXPMOI4" // Your contract
// 	networkPassphrase = network.TestNetworkPassphrase
// 	signerSecret      = os.Getenv("SIGNER_SECRET") // Or hardcode temporarily
// sorobanClient     = sorobanrpc.NewClient("https://soroban-testnet.stellar.org:443")
// client = horizonclient.DefaultTestNetClient
// )

type SorobanInvokeRequest struct {
	ContractID string        `json:"contract_id"`
	Function   string        `json:"function"`
	Args       []interface{} `json:"args"`
}

// CallSorobanFunction executes a Soroban contract function
func CallSorobanFunction(contractID, functionName string, args []string) (string, error) {
	cmd := []string{
		"contract", "invoke",
		"--id", contractID,
		"--network", "testnet",
		"--",
		functionName,
	}

	// For contribute function, convert args to named parameters
	if functionName == "contribute" && len(args) >= 2 {
		cmd = append(cmd, "--user", args[0], "--amount", args[1])
	} else if functionName == "get_balance" && len(args) >= 1 {
		cmd = append(cmd, "--user", args[0])
	} else if functionName == "withdraw" && len(args) >= 2 {
		cmd = append(cmd, "--user", args[0], "--amount", args[1])
	} else if functionName == "get_contribution_history" && len(args) >= 1 {
		cmd = append(cmd, "--user", args[0])
	} else {
		// For other functions, add args as-is
		cmd = append(cmd, args...)
	}

	// Execute the command
	out, err := exec.Command("soroban", cmd...).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("soroban error: %s", string(out))
	}

	return strings.TrimSpace(string(out)), nil
}

// Alternative version if source account is needed
func CallSorobanFunctionWithAuth(contractID, functionName, userSecretKey string, args []string) (string, error) {
	if contractID == "" {
		return "", fmt.Errorf("contract ID is required")
	}

	if userSecretKey == "" {
		return "", fmt.Errorf("user secret key is required")
	}

	keyName := fmt.Sprintf("temp-user-key-%d", os.Getpid())

	// Add user's key temporarily
	addKeyCmd := exec.Command("soroban", "keys", "add", keyName, "--secret-key")
	addKeyCmd.Stdin = strings.NewReader(userSecretKey)

	var addKeyStderr bytes.Buffer
	addKeyCmd.Stderr = &addKeyStderr

	if err := addKeyCmd.Run(); err != nil {
		return "", fmt.Errorf("failed to add user key: %v, stderr: %s", err, addKeyStderr.String())
	}

	// Ensure cleanup
	defer func() {
		cleanupCmd := exec.Command("soroban", "keys", "rm", keyName)
		cleanupCmd.Run() // Ignore errors in cleanup
	}()

	// Build command arguments
	cmdArgs := []string{
		"contract", "invoke",
		"--id", contractID,
		"--source-account", keyName, // Use the temporary key
		"--network", "testnet",
		"--", functionName,
	}

	// Add function-specific arguments - keeping your existing logic
	if functionName == "contribute" && len(args) >= 2 {
		cmdArgs = append(cmdArgs, "--user", args[0], "--amount", args[1])
	} else if functionName == "get_balance" && len(args) >= 1 {
		cmdArgs = append(cmdArgs, "--user", args[0])
	} else if functionName == "withdraw" && len(args) >= 2 {
		cmdArgs = append(cmdArgs, "--user", args[0], "--amount", args[1])
	} else if functionName == "get_contribution_history" && len(args) >= 1 {
		cmdArgs = append(cmdArgs, "--user", args[0])
	} else {
		// For other functions, add args as-is
		cmdArgs = append(cmdArgs, args...)
	}

	// Execute the command
	cmd := exec.Command("soroban", cmdArgs...)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("soroban invoke failed: %v, stderr: %s", err, stderr.String())
	}

	return strings.TrimSpace(out.String()), nil
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
