#!/bin/bash

# Debug script for group contribution issues
# Usage: ./debug-group-contribution.sh

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
API_URL="http://localhost:3000"
TEST_USER_EMAIL="test@example.com"
TEST_USER_PASSWORD="password123"
TEST_USER_NAME="Test User"

echo -e "${BLUE}🔧 Debugging Group Contribution Issues${NC}"
echo "======================================"
echo ""

# Check if API is running
echo -e "${BLUE}🔍 Checking API availability...${NC}"
if ! curl -s --max-time 3 "$API_URL" > /dev/null; then
    echo -e "${RED}❌ API is not responding. Please start your backend server.${NC}"
    exit 1
fi
echo -e "${GREEN}✅ API is running${NC}"
echo ""

# Step 1: Register/Login test user
echo -e "${YELLOW}Step 1: User Authentication${NC}"
echo "Attempting to register test user..."

REGISTER_DATA="{
    \"name\": \"$TEST_USER_NAME\",
    \"email\": \"$TEST_USER_EMAIL\",
    \"password\": \"$TEST_USER_PASSWORD\"
}"

REGISTER_RESPONSE=$(curl -s -X POST "$API_URL/auth/register" \
    -H "Content-Type: application/json" \
    -d "$REGISTER_DATA" \
    -w "\nHTTP_STATUS:%{http_code}")

REGISTER_STATUS=$(echo "$REGISTER_RESPONSE" | grep "HTTP_STATUS" | cut -d: -f2)
REGISTER_JSON=$(echo "$REGISTER_RESPONSE" | sed '/HTTP_STATUS/d')

if [[ "$REGISTER_STATUS" == "409" ]] || [[ "$REGISTER_STATUS" == "400" ]]; then
    echo "User already exists, attempting login..."
    
    LOGIN_DATA="{
        \"email\": \"$TEST_USER_EMAIL\",
        \"password\": \"$TEST_USER_PASSWORD\"
    }"
    
    LOGIN_RESPONSE=$(curl -s -X POST "$API_URL/auth/login" \
        -H "Content-Type: application/json" \
        -d "$LOGIN_DATA")
    
    TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.token // empty')
    USER_ID=$(echo "$LOGIN_RESPONSE" | jq -r '.user.id // empty')
    USER_WALLET=$(echo "$LOGIN_RESPONSE" | jq -r '.user.wallet // empty')
    USER_SECRET=$(echo "$LOGIN_RESPONSE" | jq -r '.user.secret_key // empty')
elif [[ "$REGISTER_STATUS" =~ ^2[0-9]{2}$ ]]; then
    echo -e "${GREEN}✅ User registered successfully${NC}"
    TOKEN=$(echo "$REGISTER_JSON" | jq -r '.token // empty')
    USER_ID=$(echo "$REGISTER_JSON" | jq -r '.user.id // empty')
    USER_WALLET=$(echo "$REGISTER_JSON" | jq -r '.user.wallet // empty')
    USER_SECRET=$(echo "$REGISTER_JSON" | jq -r '.user.secret_key // empty')
else
    echo -e "${RED}❌ Failed to register/login user${NC}"
    echo "$REGISTER_JSON"
    exit 1
fi

if [[ -z "$TOKEN" ]]; then
    echo -e "${RED}❌ Failed to get authentication token${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Authentication successful${NC}"
echo "User ID: $USER_ID"
echo "Wallet: $USER_WALLET"
echo ""

# Step 2: Fund user account
echo -e "${YELLOW}Step 2: Fund User Account${NC}"
FUND_RESPONSE=$(curl -s -X POST "$API_URL/fund/$USER_WALLET" \
    -H "Authorization: Bearer $TOKEN")
echo "$FUND_RESPONSE" | jq .
echo ""

# Wait for funding
sleep 3

# Step 3: Check user balance
echo -e "${YELLOW}Step 3: Check User Balance${NC}"
USER_BALANCE_RESPONSE=$(curl -s "$API_URL/balance/$USER_WALLET" \
    -H "Authorization: Bearer $TOKEN")
echo "$USER_BALANCE_RESPONSE" | jq .
echo ""

# Step 4: Create a test group
echo -e "${YELLOW}Step 4: Create Test Group${NC}"
GROUP_DATA="{
    \"name\": \"Debug Test Group $(date +%s)\",
    \"description\": \"Test group for debugging contributions\"
}"

GROUP_RESPONSE=$(curl -s -X POST "$API_URL/group/create" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d "$GROUP_DATA" \
    -w "\nHTTP_STATUS:%{http_code}")

GROUP_STATUS=$(echo "$GROUP_RESPONSE" | grep "HTTP_STATUS" | cut -d: -f2)
GROUP_JSON=$(echo "$GROUP_RESPONSE" | sed '/HTTP_STATUS/d')

if [[ "$GROUP_STATUS" =~ ^2[0-9]{2}$ ]]; then
    echo -e "${GREEN}✅ Group created successfully${NC}"
    GROUP_ID=$(echo "$GROUP_JSON" | jq -r '.group.id // .group.ID // empty')
    CONTRACT_ID=$(echo "$GROUP_JSON" | jq -r '.group.contract_id // .group.ContractID // empty')
    GROUP_WALLET=$(echo "$GROUP_JSON" | jq -r '.group.wallet // .group.Wallet // empty')
    
    echo "Group ID: $GROUP_ID"
    echo "Contract ID: $CONTRACT_ID"
    echo "Group Wallet: $GROUP_WALLET"
    echo "$GROUP_JSON" | jq .
else
    echo -e "${RED}❌ Failed to create group${NC}"
    echo "$GROUP_JSON"
    exit 1
fi
echo ""

# Step 5: Test direct contract contribution
echo -e "${YELLOW}Step 5: Test Direct Contract Contribution${NC}"
DIRECT_CONTRIB_DATA="{
    \"contract_id\": \"$CONTRACT_ID\",
    \"user_address\": \"$USER_WALLET\",
    \"amount\": \"50\",
    \"secret_key\": \"$USER_SECRET\"
}"

test_endpoint "POST" "$API_URL/contribute" "$DIRECT_CONTRIB_DATA" "Direct Contract Contribution"

# Step 6: Check contract balance
echo -e "${YELLOW}Step 6: Check Contract Balance${NC}"
CONTRACT_BALANCE_DATA="{
    \"contract_id\": \"$CONTRACT_ID\",
    \"user_address\": \"$USER_WALLET\"
}"

test_endpoint "POST" "$API_URL/balance" "$CONTRACT_BALANCE_DATA" "Check Contract Balance"

# Step 7: Test group contribution (if group is active)
echo -e "${YELLOW}Step 7: Test Group Contribution${NC}"
GROUP_CONTRIB_DATA="{
    \"from\": \"$USER_WALLET\",
    \"secret\": \"$USER_SECRET\",
    \"amount\": \"25\"
}"

GROUP_CONTRIB_RESPONSE=$(curl -s -X POST "$API_URL/group/$GROUP_ID/contribute" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d "$GROUP_CONTRIB_DATA" \
    -w "\nHTTP_STATUS:%{http_code}")

GROUP_CONTRIB_STATUS=$(echo "$GROUP_CONTRIB_RESPONSE" | grep "HTTP_STATUS" | cut -d: -f2)
GROUP_CONTRIB_JSON=$(echo "$GROUP_CONTRIB_RESPONSE" | sed '/HTTP_STATUS/d')

if [[ "$GROUP_CONTRIB_STATUS" =~ ^2[0-9]{2}$ ]]; then
    echo -e "${GREEN}✅ Group contribution successful${NC}"
    echo "$GROUP_CONTRIB_JSON" | jq .
else
    echo -e "${RED}❌ Group contribution failed (HTTP $GROUP_CONTRIB_STATUS)${NC}"
    echo "$GROUP_CONTRIB_JSON"
    
    # If group is not active, that's expected
    if echo "$GROUP_CONTRIB_JSON" | grep -q "not active"; then
        echo -e "${YELLOW}ℹ️  This is expected - group needs to be activated first${NC}"
    fi
fi
echo ""

# Step 8: Check final balances
echo -e "${YELLOW}Step 8: Final Balance Check${NC}"
echo "User wallet balance:"
curl -s "$API_URL/balance/$USER_WALLET" | jq .
echo ""

echo "Group wallet balance:"
curl -s "$API_URL/balance/$GROUP_WALLET" | jq .
echo ""

echo -e "${BLUE}🎯 Debug Summary${NC}"
echo "================"
echo -e "${GREEN}✅ Group contribution debugging completed${NC}"
echo ""
echo -e "${YELLOW}💡 Key Information:${NC}"
echo "Test User ID: $USER_ID"
echo "Test User Wallet: $USER_WALLET"
echo "Test Group ID: $GROUP_ID"
echo "Test Contract ID: $CONTRACT_ID"
echo "Test Group Wallet: $GROUP_WALLET"
echo ""
echo -e "${YELLOW}🔧 To activate the group for contributions:${NC}"
echo "1. The group needs to be approved first"
echo "2. Then it can be activated with contribution settings"
echo "3. Only then can members make contributions"
echo ""