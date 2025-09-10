package services

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"chama-wallet-backend/config"
)


func DeployChamaContract() (string, error) {
	if config.Config.IsMainnet {
		return "", fmt.Errorf("contract deployment should be done manually on mainnet for security. Use the configured SOROBAN_CONTRACT_ID instead")
	}

	// Load keys from environment
	source := os.Getenv("SOROBAN_PUBLIC_KEY")
	secret := os.Getenv("SOROBAN_SECRET_KEY")

	if source == "" || secret == "" {
		// Fallback to default test account
		source = "malika"
		secret = os.Getenv("SOROBAN_SECRET_KEY")
		if secret == "" {
			return "", fmt.Errorf("missing SOROBAN_SECRET_KEY in environment")
		}
	}

	// Check if WASM file exists
	wasmPath := "./chama_savings/target/wasm32-unknown-unknown/release/chama_savings.wasm"
	if _, err := os.Stat(wasmPath); os.IsNotExist(err) {
		// Try alternative path
		wasmPath = "./chama_savings.wasm"
		if _, err := os.Stat(wasmPath); os.IsNotExist(err) {
			return "", fmt.Errorf("WASM file not found. Please build the contract first with: cd chama_savings && stellar contract build")
		}
	}

	fmt.Printf("🔧 Deploying contract from WASM: %s on %s\n", wasmPath, config.Config.Network)
	fmt.Printf("🔧 Using source account: %s\n", source)

	network := config.GetSorobanNetwork()

	// Deploy using source account name (should be configured in soroban keys)
	cmd := exec.Command("soroban",
		"contract", "deploy",
		"--wasm", wasmPath,
		"--source-account", source,
		"--network", network,
	)

	fmt.Printf("🚀 Running deployment command on %s...\n", network)

	// Capture stdout and stderr
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	// Execute
	execErr := cmd.Run()
	if execErr != nil {
		fmt.Printf("❌ Deployment error: %v\n", execErr)
		fmt.Printf("❗ stderr: %s\n", stderr.String())
		fmt.Printf("❗ stdout: %s\n", out.String())

		// Try alternative method with temporary key storage
		return deployWithKeyStorage(source, secret)
	}

	output := strings.TrimSpace(out.String())
	fmt.Printf("✅ Contract deployed successfully on %s. Output: %s\n", network, output)

	// Extract contract address from output (usually the last line)
	lines := strings.Split(output, "\n")
	contractAddress := strings.TrimSpace(lines[len(lines)-1])

	// Validate contract address format
	if len(contractAddress) != 56 || !strings.HasPrefix(contractAddress, "C") {
		return "", fmt.Errorf("invalid contract address format: %s", contractAddress)
	}

	fmt.Printf("✅ Contract deployed at address: %s on %s\n", contractAddress, network)
	return contractAddress, nil
}

// deployWithKeyStorage: Alternative deployment method using temporary key storage
func deployWithKeyStorage(source, secret string) (string, error) {
	fmt.Printf("🔄 Trying alternative deployment method with key storage on %s...\n", config.Config.Network)

	keyName := fmt.Sprintf("temp-deploy-key-%d", time.Now().Unix())

	// Step 1: Add key to soroban keys
	addKeyCmd := exec.Command("soroban", "keys", "add", keyName, "--secret-key")
	addKeyCmd.Stdin = strings.NewReader(secret)

	var addKeyStderr bytes.Buffer
	addKeyCmd.Stderr = &addKeyStderr

	if err := addKeyCmd.Run(); err != nil {
		fmt.Printf("❌ Failed to add key: %v, stderr: %s\n", err, addKeyStderr.String())
		return "", fmt.Errorf("failed to add key: %v", err)
	}

	fmt.Printf("✅ Temporary key added: %s\n", keyName)

	// Ensure cleanup
	defer func() {
		fmt.Printf("🧹 Cleaning up temporary key: %s\n", keyName)
		cleanupCmd := exec.Command("soroban", "keys", "rm", keyName)
		if err := cleanupCmd.Run(); err != nil {
			fmt.Printf("⚠️ Warning: Failed to cleanup key: %v\n", err)
		}
	}()

	// Step 2: Deploy using the stored key
	wasmPath := "./chama_savings/target/wasm32-unknown-unknown/release/chama_savings.wasm"
	if _, err := os.Stat(wasmPath); os.IsNotExist(err) {
		wasmPath = "./chama_savings.wasm"
	}

	network := config.GetSorobanNetwork()

	deployCmd := exec.Command("soroban",
		"contract", "deploy",
		"--wasm", wasmPath,
		"--source-account", keyName,
		"--network", network,
	)

	var out bytes.Buffer
	var stderr bytes.Buffer
	deployCmd.Stdout = &out
	deployCmd.Stderr = &stderr

	if err := deployCmd.Run(); err != nil {
		fmt.Printf("❌ Deploy with key storage failed: %v, stderr: %s\n", err, stderr.String())
		return "", fmt.Errorf("deployment failed: %v", err)
	}

	output := strings.TrimSpace(out.String())
	fmt.Printf("✅ Contract deployed with key storage on %s. Output: %s\n", network, output)

	// Extract contract address from output
	lines := strings.Split(output, "\n")
	contractAddress := strings.TrimSpace(lines[len(lines)-1])

	// Validate contract address format
	if len(contractAddress) != 56 || !strings.HasPrefix(contractAddress, "C") {
		return "", fmt.Errorf("invalid contract address format: %s", contractAddress)
	}

	return contractAddress, nil
}

// Function to invoke contract methods
func InvokeContract(contractAddress, method string, args []string) (string, error) {
	if contractAddress == "" {
		return "", fmt.Errorf("contract address is required")
	}

	source := os.Getenv("SOROBAN_PUBLIC_KEY")
	secret := os.Getenv("SOROBAN_SECRET_KEY")

	if source == "" || secret == "" {
		return "", fmt.Errorf("missing SOROBAN_PUBLIC_KEY or SOROBAN_SECRET_KEY in environment")
	}

	keyName := "temp-invoke-key"

	// Add key temporarily
	addKeyCmd := exec.Command("soroban", "keys", "add", keyName, "--secret-key")
	addKeyCmd.Stdin = strings.NewReader(secret)

	if err := addKeyCmd.Run(); err != nil {
		return "", fmt.Errorf("failed to add key for invoke: %v", err)
	}

	defer func() {
		cleanupCmd := exec.Command("soroban", "keys", "rm", keyName)
		cleanupCmd.Run()
	}()

	network := config.GetSorobanNetwork()

	// Build invoke command
	cmdArgs := []string{
		"contract", "invoke",
		"--id", contractAddress,
		"--source-account", keyName,
		"--network", network,
		"--", method,
	}

	// Add method arguments
	cmdArgs = append(cmdArgs, args...)

	cmd := exec.Command("soroban", cmdArgs...)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("❌ Contract invoke failed: %v, stderr: %s\n", err, stderr.String())
		return "", fmt.Errorf("contract invoke failed: %v", err)
	}

	return strings.TrimSpace(out.String()), nil
}
