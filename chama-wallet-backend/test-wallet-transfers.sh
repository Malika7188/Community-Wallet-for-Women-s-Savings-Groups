#!/bin/bash

# Test script for wallet-to-wallet transfers
# Usage: ./test-wallet-transfers.sh

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
API_URL="http://localhost:3000"

echo -e "${BLUE}💸 Testing Wallet-to-Wallet Transfers${NC}"
echo "====================================="
echo ""

# Check if API is running
echo -e "${BLUE}🔍 Checking API availability...${NC}"
if ! curl -s --max-time 3 "$API_URL" > /dev/null; then
    echo -e "${RED}❌ API is not responding. Please start your backend server.${NC}"
    exit 1
fi
echo -e "${GREEN}✅ API is running${NC}"
echo ""

# Test 1: Generate two keypairs for testing
echo -e "${YELLOW}Test 1: Generate Test Keypairs${NC}"
echo "Generating sender keypair..."
SENDER_RESPONSE=$(curl -s "$API_URL/generate-keypair")
SENDER_PUBLIC=$(echo "$SENDER_RESPONSE" | jq -r '.public_key')
SENDER_SECRET=$(echo "$SENDER_RESPONSE" | jq -r '.secret_seed')

echo "Generating receiver keypair..."
RECEIVER_RESPONSE=$(curl -s "$API_URL/generate-keypair")
RECEIVER_PUBLIC=$(echo "$RECEIVER_RESPONSE" | jq -r '.public_key')
RECEIVER_SECRET=$(echo "$RECEIVER_RESPONSE" | jq -r '.secret_seed')

echo -e "${GREEN}✅ Keypairs generated${NC}"
echo "Sender: $SENDER_PUBLIC"
echo "Receiver: $RECEIVER_PUBLIC"
echo ""

# Test 2: Fund sender account
echo -e "${YELLOW}Test 2: Fund Sender Account${NC}"
FUND_RESPONSE=$(curl -s -X POST "$API_URL/fund/$SENDER_PUBLIC")
echo "$FUND_RESPONSE" | jq .
echo ""

# Wait for funding to process
echo "Waiting 3 seconds for funding to process..."
sleep 3

# Test 3: Check sender balance
echo -e "${YELLOW}Test 3: Check Sender Balance${NC}"
BALANCE_RESPONSE=$(curl -s "$API_URL/balance/$SENDER_PUBLIC")
echo "$BALANCE_RESPONSE" | jq .
echo ""

# Test 4: Perform transfer
echo -e "${YELLOW}Test 4: Perform Transfer (10 XLM)${NC}"
TRANSFER_DATA="{
    \"from_seed\": \"$SENDER_SECRET\",
    \"to_address\": \"$RECEIVER_PUBLIC\",
    \"amount\": \"10\"
}"

echo "Transfer data:"
echo "$TRANSFER_DATA" | jq .

TRANSFER_RESPONSE=$(curl -s -X POST "$API_URL/transfer" \
    -H "Content-Type: application/json" \
    -d "$TRANSFER_DATA" \
    -w "\nHTTP_STATUS:%{http_code}")

HTTP_STATUS=$(echo "$TRANSFER_RESPONSE" | grep "HTTP_STATUS" | cut -d: -f2)
JSON_RESPONSE=$(echo "$TRANSFER_RESPONSE" | sed '/HTTP_STATUS/d')

if [[ "$HTTP_STATUS" =~ ^2[0-9]{2}$ ]]; then
    echo -e "${GREEN}✅ Transfer successful (HTTP $HTTP_STATUS)${NC}"
    echo "$JSON_RESPONSE" | jq .
    
    # Extract transaction hash
    TX_HASH=$(echo "$JSON_RESPONSE" | jq -r '.transaction_hash // empty')
    if [[ -n "$TX_HASH" ]]; then
        echo -e "${GREEN}🔗 Transaction Hash: $TX_HASH${NC}"
        echo -e "${BLUE}🌐 View on Stellar Explorer: https://stellar.expert/explorer/testnet/tx/$TX_HASH${NC}"
    fi
else
    echo -e "${RED}❌ Transfer failed (HTTP $HTTP_STATUS)${NC}"
    echo "$JSON_RESPONSE"
fi
echo ""

# Test 5: Check balances after transfer
echo -e "${YELLOW}Test 5: Check Balances After Transfer${NC}"
echo "Sender balance:"
curl -s "$API_URL/balance/$SENDER_PUBLIC" | jq .
echo ""
echo "Receiver balance:"
curl -s "$API_URL/balance/$RECEIVER_PUBLIC" | jq .
echo ""

# Test 6: Get transaction history
echo -e "${YELLOW}Test 6: Get Transaction History${NC}"
echo "Sender transaction history:"
curl -s "$API_URL/transactions/$SENDER_PUBLIC" | jq .
echo ""

echo -e "${BLUE}🎯 Transfer Test Summary${NC}"
echo "======================="
echo -e "${GREEN}✅ Wallet transfer tests completed${NC}"
echo ""
echo -e "${YELLOW}💡 Test accounts created:${NC}"
echo "Sender: $SENDER_PUBLIC"
echo "Receiver: $RECEIVER_PUBLIC"
echo ""