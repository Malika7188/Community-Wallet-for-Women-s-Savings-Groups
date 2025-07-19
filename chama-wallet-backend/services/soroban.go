// services/soroban.go
package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"

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

func CallSorobanFunction(contractID, functionName string, args []string) (string, error) {
	cmdArgs := []string{
		"contract", "invoke",
		"--rpc-url", "https://soroban-testnet.stellar.org:443",
		"--network-passphrase", "Test SDF Network ; September 2015",
		"--id", contractID,
		"--source", "malika",
		"--fn", functionName,
	}

	// Append args as --arg <value>
	for _, arg := range args {
		cmdArgs = append(cmdArgs, "--arg", arg)
	}

	cmd := exec.Command("soroban", cmdArgs...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("soroban error: %s", stderr.String())
	}

	return stdout.String(), nil
}

func SorobanContribute(userAddress string, amount string) ([]byte, error) {
	sorobanRPC := "https://rpc-futurenet.stellar.org"
	contractID := "CADHKUC557DJ2F2XGEO4BGHFIYQ6O5QDVNG637ANRAGPBSWXMXXPMOI4"

	payload := SorobanInvokeRequest{
		ContractID: contractID,
		Function:   "contribute",
		Args:       []interface{}{userAddress, amount},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal payload: %w", err)
	}

	resp, err := http.Post(sorobanRPC+"/functions/invoke", "application/json", bytes.NewReader(jsonData))
	if err != nil {
		return nil, fmt.Errorf("http post failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("non-200 from RPC: %s", string(body))
	}

	return body, nil
}
