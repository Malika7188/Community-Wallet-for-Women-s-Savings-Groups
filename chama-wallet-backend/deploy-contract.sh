#!/bin/bash

# Script to deploy the Chama contract manually
# This should be run once to deploy the contract and get the contract ID

set -e

echo "🚀 Deploying Chama Contract..."

# Load environment variables
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
fi

# Check if required variables are set
if [ -z "$SOROBAN_SECRET_KEY" ]; then
    echo "❌ Error: SOROBAN_SECRET_KEY not set in .env file"
    exit 1
fi

# Set defaults
RPC_URL="${SOROBAN_RPC_URL:-https://soroban-testnet.stellar.org:443}"
NETWORK_PASSPHRASE="${SOROBAN_NETWORK_PASSPHRASE:-Test SDF Network ; September 2015}"
WASM_PATH="./chama_savings.wasm"

# Check if WASM file exists
if [ ! -f "$WASM_PATH" ]; then
    echo "❌ Error: WASM file not found at $WASM_PATH"
    echo "Please build the contract first"
    exit 1
fi

echo "📦 WASM file: $WASM_PATH"
echo "🌐 RPC URL: $RPC_URL"
echo "🔑 Using secret key from environment"

# Fund the account (in case it's not funded)
echo "💰 Funding account..."
soroban keys fund "$SOROBAN_SECRET_KEY" \
    --rpc-url "$RPC_URL" \
    --network-passphrase "$NETWORK_PASSPHRASE" || echo "⚠️ Funding skipped (account may already be funded)"

# Deploy the contract
echo "🚀 Deploying contract..."
CONTRACT_ID=$(soroban contract deploy \
    --wasm "$WASM_PATH" \
    --source "$SOROBAN_SECRET_KEY" \
    --rpc-url "$RPC_URL" \
    --network-passphrase "$NETWORK_PASSPHRASE")

echo ""
echo "✅ Contract deployed successfully!"
echo "📝 Contract ID: $CONTRACT_ID"
echo ""
echo "⚡ Next steps:"
echo "1. Add this line to your .env file:"
echo "   SOROBAN_CONTRACT_ID=$CONTRACT_ID"
echo ""
echo "2. Update the contract ID in Fly.io secrets:"
echo "   flyctl secrets set SOROBAN_CONTRACT_ID=$CONTRACT_ID"
echo ""
