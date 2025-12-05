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
	// Check if a contract is already deployed
	existingContractID := os.Getenv("SOROBAN_CONTRACT_ID")
	if existingContractID != "" {
		fmt.Printf("✅ Using existing contract: %s\n", existingContractID)
		return existingContractID, nil
	}

	fmt.Println("📝 No existing contract found, attempting to deploy a new one...")

	// Load secret key from environment
	secret := os.Getenv("SOROBAN_SECRET_KEY")
	if secret == "" {
		return "", fmt.Errorf("missing SOROBAN_SECRET_KEY in environment")
	}

	// Get network configuration from environment
	rpcURL := os.Getenv("SOROBAN_RPC_URL")
	networkPassphrase := os.Getenv("SOROBAN_NETWORK_PASSPHRASE")

	if rpcURL == "" {
		rpcURL = "https://soroban-testnet.stellar.org:443"
	}
	if networkPassphrase == "" {
		networkPassphrase = "Test SDF Network ; September 2015"
	}

	// First, ensure the account is funded
	fmt.Println("💰 Checking account funding...")
	fundCmd := exec.Command("soroban",
		"keys", "fund", secret,
		"--rpc-url", rpcURL,
		"--network-passphrase", networkPassphrase,
	)

	var fundOut bytes.Buffer
	var fundErr bytes.Buffer
	fundCmd.Stdout = &fundOut
	fundCmd.Stderr = &fundErr

	if err := fundCmd.Run(); err != nil {
		fmt.Printf("⚠️ Account funding warning: %v\n", err)
		fmt.Printf("Fund stderr: %s\n", fundErr.String())
		// Don't fail here - account might already be funded
	} else {
		fmt.Printf("✅ Account funded: %s\n", fundOut.String())
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

	fmt.Printf("🔧 Deploying contract from WASM: %s\n", wasmPath)
	fmt.Printf("🌐 RPC URL: %s\n", rpcURL)
	fmt.Printf("🌐 Network Passphrase: %s\n", networkPassphrase)

	// Deploy using secret key directly via stdin
	cmd := exec.Command("soroban",
		"contract", "deploy",
		"--wasm", wasmPath,
		"--source", secret,
		"--rpc-url", rpcURL,
		"--network-passphrase", networkPassphrase,
	)

	fmt.Println("🚀 Running deployment command...")

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
		return "", fmt.Errorf("contract deployment failed: %v", execErr)
	}

	output := strings.TrimSpace(out.String())
	fmt.Printf("✅ Contract deployed successfully. Output: %s\n", output)

	// Extract contract address from output (usually the last line)
	lines := strings.Split(output, "\n")
	contractAddress := strings.TrimSpace(lines[len(lines)-1])

	// Validate contract address format
	if len(contractAddress) != 56 || !strings.HasPrefix(contractAddress, "C") {
		return "", fmt.Errorf("invalid contract address format: %s", contractAddress)
	}

	fmt.Printf("✅ Contract deployed at address: %s\n", contractAddress)
	return contractAddress, nil
}

// Function to invoke contract methods
func InvokeContract(contractAddress, method string, args []string) (string, error) {
	if contractAddress == "" {
		return "", fmt.Errorf("contract address is required")
	}

	secret := os.Getenv("SOROBAN_SECRET_KEY")
	if secret == "" {
		return "", fmt.Errorf("missing SOROBAN_SECRET_KEY in environment")
	}

	// Get network configuration from environment
	rpcURL := os.Getenv("SOROBAN_RPC_URL")
	networkPassphrase := os.Getenv("SOROBAN_NETWORK_PASSPHRASE")

	if rpcURL == "" {
		rpcURL = "https://soroban-testnet.stellar.org:443"
	}
	if networkPassphrase == "" {
		networkPassphrase = "Test SDF Network ; September 2015"
	}

	// Build invoke command
	cmdArgs := []string{
		"contract", "invoke",
		"--id", contractAddress,
		"--source", secret,
		"--rpc-url", rpcURL,
		"--network-passphrase", networkPassphrase,
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
