package services

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/joho/godotenv"
)

func init() {
	// Load .env file and handle potential errors
	if err := godotenv.Load(".env"); err != nil {
		fmt.Printf("Warning: Error loading .env file: %v\n", err)
	}

	// Print all relevant environment variables for debugging (but not the full secret key)
	fmt.Println("Environment Variables:")
	fmt.Printf("Public Key: %s\n", os.Getenv("SOROBAN_PUBLIC_KEY"))

	secretKey := os.Getenv("SOROBAN_SECRET_KEY")
	if len(secretKey) > 4 {
		fmt.Printf("Secret Key (first 4 chars): %s...\n", secretKey[:4])
	}

	fmt.Printf("Network: %s\n", os.Getenv("SOROBAN_NETWORK"))
	fmt.Printf("Contract ID: %s\n", os.Getenv("SOROBAN_CONTRACT_ID"))

	// Check if soroban CLI is available
	if _, err := exec.LookPath("soroban"); err != nil {
		fmt.Printf("Warning: soroban CLI not found in PATH: %v\n", err)
	} else {
		// Get soroban version
		if out, err := exec.Command("soroban", "--version").Output(); err == nil {
			fmt.Printf("Soroban CLI version: %s\n", string(out))
		}
	}
}

func DeployChamaContract() (string, error) {
	// Load keys from environment
	source := os.Getenv("SOROBAN_PUBLIC_KEY")
	secret := os.Getenv("SOROBAN_SECRET_KEY")

	if source == "" || secret == "" {
		return "", fmt.Errorf("missing SOROBAN_PUBLIC_KEY or SOROBAN_SECRET_KEY in environment")
	}

	// Check if WASM file exists
	wasmPath := "./chama_savings.wasm"
	if _, err := os.Stat(wasmPath); os.IsNotExist(err) {
		return "", fmt.Errorf("WASM file not found at %s", wasmPath)
	}

	// Method 1: Try using the secret key directly with SOROBAN_SECRET_KEY environment variable
	// Set the secret key in environment for the command
	cmd := exec.Command("soroban",
		"contract", "deploy",
		"--wasm", wasmPath,
		"--source-account", source,
		"--network", "testnet",
	)

	// Set environment variables for the command
	cmd.Env = append(os.Environ(), fmt.Sprintf("SOROBAN_SECRET_KEY=%s", secret))

	fmt.Println("🚀 Running deployment command...")

	// Capture stdout and stderr
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	// Execute
	execErr := cmd.Run()
	if execErr != nil {
		fmt.Println("❌ Deployment error:", execErr)
		fmt.Println("❗ stderr:", stderr.String())

		// Try alternative method with soroban keys
		return deployWithKeyStorage(source, secret)
	}

	output := strings.TrimSpace(out.String())
	fmt.Println("✅ Contract deployed. Output:", output)

	// Extract contract address from output (usually the last line)
	lines := strings.Split(output, "\n")
	contractAddress := strings.TrimSpace(lines[len(lines)-1])

	return contractAddress, nil
}

// Alternative method: Store key temporarily in soroban keys
func deployWithKeyStorage(source, secret string) (string, error) {
	fmt.Println("🔄 Trying alternative deployment method with key storage...")

	keyName := "temp-deploy-key"

	// Step 1: Add key to soroban keys
	addKeyCmd := exec.Command("soroban", "keys", "add", keyName, "--secret-key")
	addKeyCmd.Stdin = strings.NewReader(secret)

	var addKeyStderr bytes.Buffer
	addKeyCmd.Stderr = &addKeyStderr

	if err := addKeyCmd.Run(); err != nil {
		fmt.Printf("❌ Failed to add key: %v, stderr: %s\n", err, addKeyStderr.String())
		return "", fmt.Errorf("failed to add key: %v", err)
	}

	// Ensure cleanup
	defer func() {
		cleanupCmd := exec.Command("soroban", "keys", "rm", keyName)
		cleanupCmd.Run() // Ignore errors in cleanup
	}()

	// Step 2: Deploy using the stored key
	deployCmd := exec.Command("soroban",
		"contract", "deploy",
		"--wasm", "./chama_savings.wasm",
		"--source-account", keyName, // Use the key name instead of address
		"--network", "testnet",
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
	fmt.Println("✅ Contract deployed with key storage. Output:", output)

	// Extract contract address from output
	lines := strings.Split(output, "\n")
	contractAddress := strings.TrimSpace(lines[len(lines)-1])

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

	// Build invoke command
	cmdArgs := []string{
		"contract", "invoke",
		"--id", contractAddress,
		"--source-account", keyName,
		"--network", "testnet",
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
