#!/bin/bash

# Script to deploy Chama Savings contract to Stellar Mainnet
# Usage: ./deploy-mainnet-contract.sh

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🚀 Deploying Chama Savings Contract to Mainnet${NC}"
echo "=============================================="
echo ""

# Check if we're in the right directory
if [ ! -d "chama_savings" ]; then
    echo -e "${RED}❌ chama_savings directory not found${NC}"
    echo "Please run this script from the chama-wallet-backend directory"
    exit 1
fi

# Check if soroban CLI is installed
if ! command -v soroban &> /dev/null; then
    echo -e "${RED}❌ Soroban CLI not found${NC}"
    echo "Please install Soroban CLI first"
    exit 1
fi

# Load environment variables
if [ -f ".env" ]; then
    source .env
else
    echo -e "${RED}❌ .env file not found${NC}"
    echo "Please create .env file with mainnet configuration"
    exit 1
fi

# Validate environment variables
if [ -z "$SOROBAN_PUBLIC_KEY" ] || [ -z "$SOROBAN_SECRET_KEY" ]; then
    echo -e "${RED}❌ Missing SOROBAN_PUBLIC_KEY or SOROBAN_SECRET_KEY${NC}"
    echo "Please set these in your .env file"
    exit 1
fi

if [ "$STELLAR_NETWORK" != "mainnet" ]; then
    echo -e "${YELLOW}⚠️  STELLAR_NETWORK is not set to mainnet${NC}"
    read -p "Continue anyway? (y/N): " confirm
    if [[ ! $confirm =~ ^[Yy]$ ]]; then
        exit 0
    fi
fi

echo -e "${YELLOW}📋 Deployment Configuration:${NC}"
echo "Network: mainnet"
echo "Source Account: $SOROBAN_PUBLIC_KEY"
echo "Contract Directory: chama_savings"
echo ""

# Confirm deployment
echo -e "${RED}⚠️  WARNING: This will deploy to MAINNET using real XLM!${NC}"
read -p "Are you sure you want to continue? (y/N): " confirm
if [[ ! $confirm =~ ^[Yy]$ ]]; then
    echo "Deployment cancelled"
    exit 0
fi

# Step 1: Build the contract
echo -e "${BLUE}Step 1: Building contract...${NC}"
cd chama_savings
if ! stellar contract build; then
    echo -e "${RED}❌ Contract build failed${NC}"
    exit 1
fi
echo -e "${GREEN}✅ Contract built successfully${NC}"
cd ..

# Step 2: Check if WASM file exists
WASM_FILE="chama_savings/target/wasm32-unknown-unknown/release/chama_savings.wasm"
if [ ! -f "$WASM_FILE" ]; then
    echo -e "${RED}❌ WASM file not found: $WASM_FILE${NC}"
    exit 1
fi

echo -e "${GREEN}✅ WASM file found: $WASM_FILE${NC}"

# Step 3: Add mainnet network configuration
echo -e "${BLUE}Step 2: Configuring mainnet network...${NC}"
soroban network add \
  --global mainnet \
  --rpc-url https://soroban-rpc.mainnet.stellar.org:443 \
  --network-passphrase "Public Global Stellar Network ; September 2015"

echo -e "${GREEN}✅ Mainnet network configured${NC}"

# Step 4: Add deployment key
echo -e "${BLUE}Step 3: Adding deployment key...${NC}"
TEMP_KEY_NAME="mainnet-deploy-$(date +%s)"

# Add key securely
echo "$SOROBAN_SECRET_KEY" | soroban keys add "$TEMP_KEY_NAME" --secret-key

if [ $? -ne 0 ]; then
    echo -e "${RED}❌ Failed to add deployment key${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Deployment key added: $TEMP_KEY_NAME${NC}"

# Cleanup function
cleanup() {
    echo -e "${BLUE}🧹 Cleaning up temporary key...${NC}"
    soroban keys rm "$TEMP_KEY_NAME" 2>/dev/null
}

# Ensure cleanup on exit
trap cleanup EXIT

# Step 5: Deploy contract
echo -e "${BLUE}Step 4: Deploying contract to mainnet...${NC}"
echo -e "${YELLOW}This may take a few minutes...${NC}"

CONTRACT_ID=$(soroban contract deploy \
  --wasm "$WASM_FILE" \
  --source-account "$TEMP_KEY_NAME" \
  --network mainnet 2>&1)

if [ $? -ne 0 ]; then
    echo -e "${RED}❌ Contract deployment failed${NC}"
    echo "$CONTRACT_ID"
    exit 1
fi

# Extract contract ID (remove any extra output)
CONTRACT_ID=$(echo "$CONTRACT_ID" | tail -n 1 | tr -d '\n\r')

echo -e "${GREEN}✅ Contract deployed successfully!${NC}"
echo ""
echo -e "${BLUE}📋 Deployment Results:${NC}"
echo "Contract ID: $CONTRACT_ID"
echo "Network: mainnet"
echo "Explorer: https://stellar.expert/explorer/public/contract/$CONTRACT_ID"
echo ""

# Step 6: Initialize contract
echo -e "${BLUE}Step 5: Initializing contract...${NC}"
INIT_RESULT=$(soroban contract invoke \
  --id "$CONTRACT_ID" \
  --source-account "$TEMP_KEY_NAME" \
  --network mainnet \
  -- initialize 2>&1)

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ Contract initialized successfully${NC}"
else
    echo -e "${YELLOW}⚠️  Contract initialization result: $INIT_RESULT${NC}"
    echo "This might be normal if the contract was already initialized"
fi

# Step 7: Update .env file
echo -e "${BLUE}Step 6: Updating .env file...${NC}"
if grep -q "SOROBAN_CONTRACT_ID=" .env; then
    # Update existing line
    sed -i "s/SOROBAN_CONTRACT_ID=.*/SOROBAN_CONTRACT_ID=$CONTRACT_ID/" .env
else
    # Add new line
    echo "SOROBAN_CONTRACT_ID=$CONTRACT_ID" >> .env
fi

echo -e "${GREEN}✅ .env file updated with contract ID${NC}"

echo ""
echo -e "${GREEN}🎉 Mainnet Deployment Complete!${NC}"
echo "================================"
echo ""
echo -e "${BLUE}📋 Next Steps:${NC}"
echo "1. Update your .env file with STELLAR_NETWORK=mainnet"
echo "2. Restart your backend server"
echo "3. Test with small amounts first"
echo "4. Monitor transactions on Stellar Explorer"
echo ""
echo -e "${BLUE}📋 Important Information:${NC}"
echo "Contract ID: $CONTRACT_ID"
echo "Network: mainnet"
echo "Explorer: https://stellar.expert/explorer/public/contract/$CONTRACT_ID"
echo ""
echo -e "${RED}⚠️  Security Reminder:${NC}"
echo "- Keep your secret keys secure"
echo "- Test thoroughly before using large amounts"
echo "- Monitor all transactions"
echo "- Have a backup plan for key recovery"