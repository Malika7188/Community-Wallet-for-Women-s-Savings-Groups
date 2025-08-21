// services/soroban.go
package services

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

type SorobanInvokeRequest struct {
	ContractID string        `json:"contract_id"`
	Function   string        `json:"function"`
	Args       []interface{} `json:"args"`
}

// validateContractID ensures the contract ID is valid
func validateContractID(contractID string) error {
	if contractID == "" {
		return fmt.Errorf("contract ID cannot be empty")
	}
	if len(contractID) != 56 {
		return fmt.Errorf("invalid contract ID length: expected 56 characters, got %d", len(contractID))
	}
	if !strings.HasPrefix(contractID, "C") {
		return fmt.Errorf("contract ID must start with 'C'")
	}
	return nil
}

// checkContractExists verifies the contract exists on the network
func checkContractExists(contractID string) error {
	cmd := exec.Command("soroban", "contract", "inspect", "--id", contractID, "--network", "testnet")
	
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("contract does not exist or is not accessible: %v, stderr: %s", err, stderr.String())
	}
	return nil
}

// CallSorobanFunction executes a Soroban contract function
func CallSorobanFunction(contractID, functionName string, args []string) (string, error) {
	// Validate inputs
	if err := validateContractID(contractID); err != nil {
		return "", fmt.Errorf("invalid contract ID: %w", err)
	}
	
	if functionName == "" {
		return "", fmt.Errorf("function name cannot be empty")
	}

	// Check if contract exists
	if err := checkContractExists(contractID); err != nil {
		return "", fmt.Errorf("contract validation failed: %w", err)
	}

	cmd := []string{
		"contract", "invoke",
		"--id", contractID,
		"--network", "testnet",
		"--source-account", "malika", // Use the configured source account
		"--",
		functionName,
	}

	// Convert function arguments to proper format
	if functionName == "contribute" && len(args) >= 2 {
		cmd = append(cmd, "--user", args[0], "--amount", args[1])
	} else if functionName == "get_balance" && len(args) >= 1 {
		cmd = append(cmd, "--user", args[0])
	} else if functionName == "withdraw" && len(args) >= 2 {
		cmd = append(cmd, "--user", args[0], "--amount", args[1])
	} else if functionName == "get_contribution_history" && len(args) >= 1 {
		cmd = append(cmd, "--user", args[0])
	} else {
		cmd = append(cmd, args...)
	}

	fmt.Printf("🔧 Executing Soroban command: soroban %s\n", strings.Join(cmd, " "))
	
	// Execute the command with timeout
	execCmd := exec.Command("soroban", cmd...)
	
	var out bytes.Buffer
	var stderr bytes.Buffer
	execCmd.Stdout = &out
	execCmd.Stderr = &stderr
	
	// Set a timeout for the command
	done := make(chan error, 1)
	go func() {
		done <- execCmd.Run()
	}()
	
	select {
	case err := <-done:
		if err != nil {
			fmt.Printf("❌ Soroban command failed: %v\n", err)
			fmt.Printf("❌ Stderr: %s\n", stderr.String())
			fmt.Printf("❌ Stdout: %s\n", out.String())
			return "", fmt.Errorf("soroban invoke failed: %v, stderr: %s", err, stderr.String())
		}
	case <-time.After(30 * time.Second):
		execCmd.Process.Kill()
		return "", fmt.Errorf("soroban command timed out after 30 seconds")
	}
	
	result := strings.TrimSpace(out.String())
	fmt.Printf("✅ Soroban result: %s\n", result)
	
	if err != nil {
		return "", fmt.Errorf("soroban error: %s", string(out))
	}

	return result, nil
}

// CallSorobanFunctionWithAuth executes a Soroban contract function with user authentication
func CallSorobanFunctionWithAuth(contractID, functionName, userSecretKey string, args []string) (string, error) {
	// Validate inputs
	if err := validateContractID(contractID); err != nil {
		return "", fmt.Errorf("invalid contract ID: %w", err)
	}
	
	if functionName == "" {
		return "", fmt.Errorf("function name cannot be empty")
	}
	
	if userSecretKey == "" {
		return "", fmt.Errorf("user secret key is required")
	}

	// Check if contract exists
	if err := checkContractExists(contractID); err != nil {
		return "", fmt.Errorf("contract validation failed: %w", err)
	}

	keyName := fmt.Sprintf("temp-user-key-%d-%d", os.Getpid(), time.Now().Unix())
	
	fmt.Printf("🔑 Adding temporary key: %s\n", keyName)

	// Add user's key temporarily with better error handling
	addKeyCmd := exec.Command("soroban", "keys", "add", keyName, "--secret-key")
	addKeyCmd.Stdin = strings.NewReader(userSecretKey)

	var addKeyStderr bytes.Buffer
	var addKeyStdout bytes.Buffer
	addKeyCmd.Stderr = &addKeyStderr
	addKeyCmd.Stdout = &addKeyStdout

	if err := addKeyCmd.Run(); err != nil {
		fmt.Printf("❌ Failed to add key: %v\n", err)
		fmt.Printf("❌ Stderr: %s\n", addKeyStderr.String())
		fmt.Printf("❌ Stdout: %s\n", addKeyStdout.String())
		return "", fmt.Errorf("failed to add user key: %v, stderr: %s, stdout: %s", err, addKeyStderr.String(), addKeyStdout.String())
	}
	
	fmt.Printf("✅ Key added successfully\n")

	// Ensure cleanup
	defer func() {
		fmt.Printf("🧹 Cleaning up temporary key: %s\n", keyName)
		cleanupCmd := exec.Command("soroban", "keys", "rm", keyName)
		if err := cleanupCmd.Run(); err != nil {
			fmt.Printf("⚠️ Warning: Failed to cleanup key %s: %v\n", keyName, err)
		}
	}()

	// Build command arguments
	cmdArgs := []string{
		"contract", "invoke",
		"--id", contractID,
		"--source-account", keyName,
		"--network", "testnet",
		"--",
		functionName,
	}

	// Add function-specific arguments
	if functionName == "contribute" && len(args) >= 2 {
		cmdArgs = append(cmdArgs, "--user", args[0], "--amount", args[1])
	} else if functionName == "get_balance" && len(args) >= 1 {
		cmdArgs = append(cmdArgs, "--user", args[0])
	} else if functionName == "withdraw" && len(args) >= 2 {
		cmdArgs = append(cmdArgs, "--user", args[0], "--amount", args[1])
	} else if functionName == "get_contribution_history" && len(args) >= 1 {
		cmdArgs = append(cmdArgs, "--user", args[0])
	} else {
		cmdArgs = append(cmdArgs, args...)
	}

	fmt.Printf("🔧 Executing authenticated Soroban command: soroban %s\n", strings.Join(cmdArgs, " "))
	
	// Execute the command with timeout
	cmd := exec.Command("soroban", cmdArgs...)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	// Set a timeout for the command
	done := make(chan error, 1)
	go func() {
		done <- cmd.Run()
	}()
	
	select {
	case err := <-done:
		if err != nil {
			fmt.Printf("❌ Soroban invoke failed: %v\n", err)
			fmt.Printf("❌ Stderr: %s\n", stderr.String())
			fmt.Printf("❌ Stdout: %s\n", out.String())
			return "", fmt.Errorf("soroban invoke failed: %v, stderr: %s", err, stderr.String())
		}
	case <-time.After(45 * time.Second):
		cmd.Process.Kill()
		return "", fmt.Errorf("soroban command timed out after 45 seconds")
	}

	result := strings.TrimSpace(out.String())
	fmt.Printf("✅ Authenticated Soroban result: %s\n", result)
	
	return result, nil
}

// Wrapper functions with improved error handling
func Contribute(contractID, userAddress, amount string) (string, error) {
	fmt.Printf("🔄 Contributing %s XLM from %s to contract %s\n", amount, userAddress, contractID)
	args := []string{userAddress, amount}
	return CallSorobanFunction(contractID, "contribute", args)
}

func GetBalance(contractID, userAddress string) (string, error) {
	fmt.Printf("🔍 Getting balance for %s from contract %s\n", userAddress, contractID)
	args := []string{userAddress}
	return CallSorobanFunction(contractID, "get_balance", args)
}

func Withdraw(contractID, userAddress, amount string) (string, error) {
	fmt.Printf("💸 Withdrawing %s XLM for %s from contract %s\n", amount, userAddress, contractID)
	args := []string{userAddress, amount}
	return CallSorobanFunction(contractID, "withdraw", args)
}

func GetContributionHistory(contractID, userAddress string) (string, error) {
	fmt.Printf("📊 Getting contribution history for %s from contract %s\n", userAddress, contractID)
	args := []string{userAddress}
	return CallSorobanFunction(contractID, "get_contribution_history", args)
}

// ContributeWithAuth - wrapper for authenticated contributions
func ContributeWithAuth(contractID, userAddress, amount, secretKey string) (string, error) {
	fmt.Printf("🔐 Making authenticated contribution: %s XLM from %s\n", amount, userAddress)
	args := []string{userAddress, amount}
	return CallSorobanFunctionWithAuth(contractID, "contribute", secretKey, args)
}
